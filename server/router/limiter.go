package router

import "time"

type Limiter struct {
	MaxReqPerMin int            `json:"max_requests_per_minute"` // The maximum number of requests per minute, stored here and in the DB
	IpMap        map[string]int `json:"ip_map"`                  // A map of IP addresses to the number of requests they have made
	LastReset    int64          `json:"last_reset"`              // The last time the IP map was reset
	Block        bool           `json:"block"`                   // Whether or not to completely block requests
}

// CreateLimiter creates a new Limiter struct
func CreateLimiter(maxReqPerMin int, block bool) *Limiter {
	return &Limiter{
		MaxReqPerMin: maxReqPerMin,
		IpMap:        make(map[string]int),
		LastReset:    time.Now().Unix(),
		Block:        block,
	}
}

// CheckNewRequest checks if a new request is allowed
func (l *Limiter) CheckNewRequest(ip string) bool {
	// Block all requests
	if l.Block {
		return false
	}

	// If the IP map was last reset more than a minute ago, reset it
	if time.Now().Unix()-l.LastReset > 60 {
		l.IpMap = make(map[string]int)
		l.LastReset = time.Now().Unix()
	}

	// Create entry in IpMap if not exists
	if _, ok := l.IpMap[ip]; !ok {
		l.IpMap[ip] = 0
	}

	// If IP has exceeded the max requests per min, block it
	if l.IpMap[ip] >= l.MaxReqPerMin {
		return false
	}

	// Increment the request count for the IP
	l.IpMap[ip]++
	return true
}
