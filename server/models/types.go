package models

type JudgeStats struct {
	Num       int64   `json:"num"`
	AvgVotes  float64 `json:"avg_votes"`
	NumActive int64   `json:"num_active"`
}

type ProjectStats struct {
	Num      int64   `json:"num"`
	AvgVotes float64 `json:"avg_votes"`
	AvgSeen  float64 `json:"avg_seen"`
}

type Stats struct {
	Projects int64   `json:"projects"`
	Judges   int64   `json:"judges"`
	Seen     int64   `json:"seen"`
	Votes    int64   `json:"votes"`
	AvgMu    float64 `json:"avg_mu"`
	AvgSigma float64 `json:"avg_sigma"`
}
