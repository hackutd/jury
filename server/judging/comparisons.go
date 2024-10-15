package judging

import (
	"server/database"
	"server/models"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Comparisons struct {
	Arr      [][]int                    `json:"arr"`
	IdNumMap map[primitive.ObjectID]int `json:"id_num_map"`
	Mutex    sync.Mutex                 `json:"mutex"`
}

// CreateComparisons will create the comparisons array from a list of
// judges and projects
func CreateComparisons(projects []*models.Project, judges []*models.Judge) *Comparisons {
	// Create comps object
	comps := Comparisons{
		Arr:      make([][]int, len(projects)),
		IdNumMap: make(map[primitive.ObjectID]int),
		Mutex:    sync.Mutex{},
	}

	// Fill array
	for i := range comps.Arr {
		comps.Arr[i] = make([]int, len(projects))
	}

	// Loop through projects and create the map
	for i, p := range projects {
		comps.IdNumMap[p.Id] = i
	}

	// For each judge, add their comparisons
	for _, j := range judges {
		// Loop through list of judged projects and add each comparison to the matrix
		for i, ap := range j.SeenProjects {
			for _, bp := range j.SeenProjects[i+1:] {
				aid := comps.IdNumMap[ap.ProjectId]
				bid := comps.IdNumMap[bp.ProjectId]
				comps.Arr[aid][bid]++
				comps.Arr[bid][aid]++
			}
		}
	}

	return &comps
}

// LoadComparisons will create the comparisons from the database
func LoadComparisons(db *mongo.Database) (*Comparisons, error) {
	// Get all judges
	judges, err := database.FindAllJudges(db)
	if err != nil {
		return nil, err
	}

	// Get all project
	projects, err := database.FindAllProjects(db)
	if err != nil {
		return nil, err
	}

	return CreateComparisons(projects, judges), nil
}

// AddProjectToComparisons will increase the size of the comparisons array when
// adding a project to the database
func (c *Comparisons) AddProjectToComparisons(p *models.Project) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	// Add project to map
	num := len(c.IdNumMap)
	c.IdNumMap[p.Id] = num

	// Increase size of matrix
	for i := range c.Arr {
		c.Arr[i] = append(c.Arr[i], 0)
	}
	c.Arr = append(c.Arr, make([]int, num+1))

	return nil
}

// RemoveProjectFromComparison removes a project from the comparisons array.
// This is a really expensive function, but should only be run when projects are deleted,
// which really shouldn't happen all that often.
func (c *Comparisons) RemoveProjectFromComparison(id primitive.ObjectID) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	// Get project to remove
	removeId := c.IdNumMap[id]

	// Remove project from each row
	for i, v := range c.Arr {
		c.Arr[i] = removeSliceItem(v, removeId)
	}

	// Remove the row itself
	c.Arr = removeSliceItem(c.Arr, removeId)

	return nil
}

// removeSliceItem removes an item from a slice
func removeSliceItem[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

// Once a judge has judged a project, add all pairs of that project to
// the projects that the judge has judged to the array
func (c *Comparisons) UpdateProjectComparisonCount(prevSeen []models.JudgedProject, newProj primitive.ObjectID) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	// Increment the indices of the current project + all prev projects
	newProjId := c.IdNumMap[newProj]
	for _, v := range prevSeen {
		vId := c.IdNumMap[v.ProjectId]
		c.Arr[vId][newProjId]++
		c.Arr[newProjId][vId]++
	}
}

// FindLeastCompared finds the project that has been compared the LEAST
// to all other projects. Projects param MUST not be empty.
func (c *Comparisons) FindLeastCompared(projects []*models.Project, prevSeen []models.JudgedProject) *models.Project {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	minProj := projects[0]
	min := 0x7FFFFFFF
	// Loop through all potential projects and find the one with the most seen
	for _, v := range projects {
		curr := 0
		vId := c.IdNumMap[v.Id]
		for _, p := range prevSeen {
			pId := c.IdNumMap[p.ProjectId]
			curr += c.Arr[vId][pId]
		}
		if curr < min {
			min = curr
			minProj = v
		}
	}

	return minProj
}
