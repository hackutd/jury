package models

type JudgeStats struct {
	Num       int64   `json:"num"`
	AvgSeen   float64 `json:"avg_seen"`
	NumActive int64   `json:"num_active"`
}

type ProjectStats struct {
	Num       int64   `json:"num"`
	AvgSeen   float64 `json:"avg_seen"`
	NumActive int64   `json:"num_active"`
}

type Stats struct {
	Projects       int64   `json:"projects"`
	Judges         int64   `json:"judges"`
	AvgProjectSeen float64 `json:"avg_project_seen"`
	AvgJudgeSeen   float64 `json:"avg_judge_seen"`
}

type JudgeVote struct {
	CurrWinner bool `json:"curr_winner"`
}

type IdRequest struct {
	Id string `json:"id"`
}

type AddJudgeRequest struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Track  string `json:"track"`
	Notes  string `json:"notes"`
	NoSend bool   `json:"no_send"`
}

type AddProjectRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Url           string `json:"url"`
	TryLink       string `json:"try_link"`
	VideoLink     string `json:"video_link"`
	ChallengeList string `json:"challenge_list"`
}
