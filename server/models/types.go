package models

type JudgeStats struct {
	Num       int64   `json:"num"`
	AvgVotes  float64 `json:"avg_votes"`
	NumActive int64   `json:"num_active"`
}
