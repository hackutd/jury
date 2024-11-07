package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Options struct {
	Id             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Ref            int64              `bson:"ref" json:"ref"`
	Clock          ClockState         `bson:"clock" json:"clock"`
	JudgingTimer   int64              `bson:"judging_timer" json:"judging_timer"`
	MinViews       int64              `bson:"min_views" json:"min_views"`
	ClockSync      bool               `bson:"clock_sync" json:"clock_sync"`
	Categories     []string           `bson:"categories" json:"categories"`
	JudgeTracks    bool               `bson:"judge_tracks" json:"judge_tracks"`
	Tracks         []string           `bson:"tracks" json:"tracks"`
	MultiGroup     bool               `bson:"multi_group" json:"multi_group"`
	NumGroups      int64              `bson:"num_groups" json:"num_groups"`             // Number of groups to split projects into
	GroupSizes     []int64            `bson:"group_sizes" json:"group_sizes"`           // Number of projects in each group except last (last group will be remainder, size = numGroups - 1)
	GroupTableNums []int64            `bson:"group_table_nums" json:"group_table_nums"` // Table number assignments to each group -- this will be round-robin due to the way groups are split
	MainGroup      GroupSwitchOps     `bson:"main_group" json:"main_group"`             // Group options for the general judging track
}

type GroupSwitchOps struct {
	SwitchingMode    string  `bson:"switching_mode" json:"switching_mode"`         // "auto" or "manual"
	AutoSwitchMethod string  `bson:"auto_switch_method" json:"auto_switch_method"` // "count" or "proportion"
	AutoSwitchCount  int64   `bson:"auto_switch_count" json:"auto_switch_count"`   // Number of projects to view in each group
	AutoSwitchProp   float64 `bson:"auto_switch_prop" json:"auto_switch_prop"`     // Proportion of projects to view in each group
	ManualSwitches   int64   `bson:"manual_switches" json:"manual_switches"`       // Number of manual switches that have happened
}

func NewOptions() *Options {
	return &Options{
		Ref:            0,
		JudgingTimer:   300,
		MinViews:       3,
		Clock:          *NewClockState(),
		ClockSync:      false,
		Categories:     []string{"Creativity/Innovation", "Technical Competence/Execution", "Research/Design", "Presentation"},
		JudgeTracks:    false,
		Tracks:         []string{},
		MultiGroup:     false,
		NumGroups:      3,
		GroupSizes:     []int64{30, 30, 30},
		GroupTableNums: []int64{0, 30, 60},
		MainGroup:      *NewGroupOptions(),
	}
}

func NewGroupOptions() *GroupSwitchOps {
	return &GroupSwitchOps{
		SwitchingMode:    "auto",
		AutoSwitchMethod: "count",
		AutoSwitchCount:  3,
		AutoSwitchProp:   0.1,
		ManualSwitches:   0,
	}
}

// OptionalGroupOptions is a struct that will be used to update the group options.
// This is used to deserialize the JSON body of the request to update the group options,
// where all fields that are not filled in will be nil.
type OptionalGroupOptions struct {
	SwitchingMode    *string  `json:"switching_mode"`
	AutoSwitchMethod *string  `json:"auto_switch_method"`
	AutoSwitchCount  *int64   `json:"auto_switch_count"`
	AutoSwitchProp   *float64 `json:"auto_switch_prop"`
}
