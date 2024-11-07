package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Options struct {
	Id             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Ref            int64              `bson:"ref" json:"ref"`
	Clock          ClockState         `bson:"clock" json:"clock"`
	JudgingTimer   int64              `bson:"judging_timer" json:"judging_timer"`
	MinViews       int64              `bson:"min_views" json:"min_views"`
	ClockSync      bool               `bson:"clock_sync" json:"clock_sync"`
	JudgeTracks    bool               `bson:"judge_tracks" json:"judge_tracks"`
	Tracks         []string           `bson:"tracks" json:"tracks"`
	MultiGroup     bool               `bson:"multi_group" json:"multi_group"`
	NumGroups      int64              `bson:"num_groups" json:"num_groups"`             // Number of groups to split projects into
	GroupSizes     []int64            `bson:"group_sizes" json:"group_sizes"`           // Number of projects in each group except last (last group will be remainder, size = numGroups - 1)
	SwitchingMode  string             `bson:"switching_mode" json:"switching_mode"`     // "auto" or "manual"
	AutoSwitchProp float64            `bson:"auto_switch_prop" json:"auto_switch_prop"` // Proportion of projects to view in each group
	ManualSwitches int64              `bson:"manual_switches" json:"manual_switches"`   // Number of manual switches that have happened
}

func NewOptions() *Options {
	return &Options{
		Ref:            0,
		JudgingTimer:   300,
		MinViews:       3,
		Clock:          *NewClockState(),
		ClockSync:      false,
		JudgeTracks:    false,
		Tracks:         []string{},
		MultiGroup:     false,
		NumGroups:      3,
		GroupSizes:     []int64{30, 30, 30},
		SwitchingMode:  "auto",
		AutoSwitchProp: 0.1,
		ManualSwitches: 0,
	}
}

type OptionalOptions struct {
	JudgingTimer   *int64    `bson:"judging_timer,omitempty" json:"judging_timer,omitempty"`
	MinViews       *int64    `bson:"min_views,omitempty" json:"min_views,omitempty"`
	ClockSync      *bool     `bson:"clock_sync,omitempty" json:"clock_sync,omitempty"`
	JudgeTracks    *bool     `bson:"judge_tracks,omitempty" json:"judge_tracks,omitempty"`
	Tracks         *[]string `bson:"tracks,omitempty" json:"tracks,omitempty"`
	MultiGroup     *bool     `bson:"multi_group,omitempty" json:"multi_group,omitempty"`
	NumGroups      *int64    `bson:"num_groups,omitempty" json:"num_groups,omitempty"`
	GroupSizes     *[]int64  `bson:"group_sizes,omitempty" json:"group_sizes,omitempty"`
	SwitchingMode  *string   `bson:"switching_mode,omitempty" json:"switching_mode,omitempty"`
	AutoSwitchProp *float64  `bson:"auto_switch_prop,omitempty" json:"auto_switch_prop,omitempty"`
}
