package router

import (
	"context"
	"server/database"
	"server/judging"
	"server/logging"
	"server/models"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// cachedStats holds a time-limited copy of AggregateStats results.
type cachedStats struct {
	value     *models.Stats
	expiresAt time.Time
}

// cachedScores holds a time-limited copy of AggregateScores results.
type cachedScores struct {
	value     map[primitive.ObjectID]judging.ProjectScores
	expiresAt time.Time
}

const aggCacheTTL = 10 * time.Second

// State is the shared application state attached to every request via middleware.
type State struct {
	Db      *mongo.Database
	Clock   *models.SafeClock
	Comps   *judging.Comparisons
	Logger  *logging.Logger
	Limiter *Limiter

	// In-memory options cache — avoids a DB round-trip for every judge action.
	optsMu sync.RWMutex
	opts   *models.Options

	// Short-lived caches for the expensive aggregation pipelines that power
	// the admin dashboard.  A 10-second TTL means concurrent admin sessions
	// share a single aggregation result instead of each triggering their own.
	statsCacheMu sync.RWMutex
	statsCache   map[string]*cachedStats // keyed by track ("" = general)

	scoresCacheMu sync.RWMutex
	scoresCache   *cachedScores
}

func NewState(db *mongo.Database, clock *models.SafeClock, comps *judging.Comparisons, logger *logging.Logger, limiter *Limiter, opts *models.Options) *State {
	return &State{
		Db:         db,
		Clock:      clock,
		Comps:      comps,
		Logger:     logger,
		Limiter:    limiter,
		opts:       opts,
		statsCache: make(map[string]*cachedStats),
	}
}

// GetCachedOptions returns the in-memory options without hitting the database.
func (s *State) GetCachedOptions() *models.Options {
	s.optsMu.RLock()
	defer s.optsMu.RUnlock()
	return s.opts
}

// SetCachedOptions replaces the in-memory options cache.
func (s *State) SetCachedOptions(o *models.Options) {
	s.optsMu.Lock()
	defer s.optsMu.Unlock()
	s.opts = o
}

// ReloadOptions re-reads options from the database and refreshes the cache.
// Call this after any operation that modifies the options document.
func (s *State) ReloadOptions(ctx context.Context) error {
	opts, err := database.GetOptions(s.Db, ctx)
	if err != nil {
		return err
	}
	s.SetCachedOptions(opts)
	return nil
}

// GetCachedStats returns AggregateStats for the given track, using a 10-second
// TTL cache to prevent redundant aggregation pipelines when multiple admins poll
// simultaneously.
func (s *State) GetCachedStats(track string) (*models.Stats, error) {
	s.statsCacheMu.RLock()
	entry, ok := s.statsCache[track]
	s.statsCacheMu.RUnlock()

	if ok && time.Now().Before(entry.expiresAt) {
		return entry.value, nil
	}

	stats, err := database.AggregateStats(s.Db, track)
	if err != nil {
		return nil, err
	}

	s.statsCacheMu.Lock()
	s.statsCache[track] = &cachedStats{value: stats, expiresAt: time.Now().Add(aggCacheTTL)}
	s.statsCacheMu.Unlock()

	return stats, nil
}

// GetCachedScores returns AggregateScores, using a 10-second TTL cache to
// avoid running the multi-stage judges→projects aggregation pipeline for
// every admin refresh cycle.
func (s *State) GetCachedScores(ctx context.Context) (map[primitive.ObjectID]judging.ProjectScores, error) {
	s.scoresCacheMu.RLock()
	entry := s.scoresCache
	s.scoresCacheMu.RUnlock()

	if entry != nil && time.Now().Before(entry.expiresAt) {
		return entry.value, nil
	}

	scores, err := judging.AggregateScores(s.Db, ctx)
	if err != nil {
		return nil, err
	}

	s.scoresCacheMu.Lock()
	s.scoresCache = &cachedScores{value: scores, expiresAt: time.Now().Add(aggCacheTTL)}
	s.scoresCacheMu.Unlock()

	return scores, nil
}

// InvalidateScoresCache clears the scores cache, forcing the next request to
// recompute.  Call this after judging actions that change scores.
func (s *State) InvalidateScoresCache() {
	s.scoresCacheMu.Lock()
	s.scoresCache = nil
	s.scoresCacheMu.Unlock()
}

func GetState(ctx *gin.Context) *State {
	state := ctx.MustGet("state").(*State)
	return state
}
