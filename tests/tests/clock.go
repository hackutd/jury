package tests

import (
	"fmt"
	"tests/util"
)

// --- Clock Tests ---

// ClockPauseUnpause verifies the clock can be paused and then resumed
func ClockPauseUnpause(context *util.Context) util.Result {
	// Pause
	pauseRes := util.PostRequest(context.Logger, "/admin/clock/pause", nil, util.AdminAuth())
	if !util.IsOk(pauseRes) {
		return util.NewResult(false, "Failed to pause clock: "+pauseRes)
	}

	clockRes := util.GetRequest(context.Logger, "/admin/clock", util.AdminAuth())
	if !util.IsValue(clockRes, "running", util.BoolType, false) {
		return util.NewResult(false, "Clock should be paused (running=false) after pause")
	}

	// Unpause
	unpauseRes := util.PostRequest(context.Logger, "/admin/clock/unpause", nil, util.AdminAuth())
	if !util.IsOk(unpauseRes) {
		return util.NewResult(false, "Failed to unpause clock: "+unpauseRes)
	}

	clockRes2 := util.GetRequest(context.Logger, "/admin/clock", util.AdminAuth())
	if !util.IsValue(clockRes2, "running", util.BoolType, true) {
		return util.NewResult(false, "Clock should be running (running=true) after unpause")
	}

	return util.ResultOk()
}

// ClockReset verifies the clock returns to a stopped zero state after reset
func ClockReset(context *util.Context) util.Result {
	resetRes := util.PostRequest(context.Logger, "/admin/clock/reset", nil, util.AdminAuth())
	if !util.IsOk(resetRes) {
		return util.NewResult(false, "Failed to reset clock: "+resetRes)
	}

	clockRes := util.GetRequest(context.Logger, "/admin/clock", util.AdminAuth())
	if !util.IsValue(clockRes, "running", util.BoolType, false) {
		return util.NewResult(false, "Clock should not be running after reset")
	}
	if !util.IsValue(clockRes, "time", util.Float64Type, 0.0) {
		return util.NewResult(false, "Clock time should be 0 after reset: "+clockRes)
	}

	return util.ResultOk()
}

// AdminStartedReflectsClockState checks /admin/started matches actual clock running state
func AdminStartedReflectsClockState(context *util.Context) util.Result {
	// Reset clock to stopped state
	util.PostRequest(context.Logger, "/admin/clock/reset", nil, util.AdminAuth())

	startedRes := util.GetRequest(context.Logger, "/admin/started", util.AdminAuth())
	// When clock is stopped, ok should be 0
	if util.IsOk(startedRes) {
		return util.NewResult(false, "/admin/started should return ok=0 when clock is not running")
	}

	// Start the clock
	util.PostRequest(context.Logger, "/admin/clock/unpause", nil, util.AdminAuth())

	startedRes2 := util.GetRequest(context.Logger, "/admin/started", util.AdminAuth())
	if !util.IsOk(startedRes2) {
		return util.NewResult(false, "/admin/started should return ok=1 when clock is running")
	}

	return util.ResultOk()
}

// --- Admin Options / Settings Tests ---

// GetOptionsReturnsValidStructure verifies /admin/options returns a well-formed response
func GetOptionsReturnsValidStructure(context *util.Context) util.Result {
	res := util.GetRequest(context.Logger, "/admin/options", util.AdminAuth())

	// Check for a few known fields
	if util.ExtractString(res, "switching_mode") == "" && !util.IsValue(res, "min_views", util.IntType, 0) {
		return util.NewResult(false, "GET /admin/options response missing expected fields: "+res)
	}

	return util.ResultOk()
}

// SetAndGetJudgingTimer verifies that setting the judging timer persists
func SetAndGetJudgingTimer(context *util.Context) util.Result {
	setRes := util.PostRequest(context.Logger, "/admin/options", util.H{
		"judging_timer": 300,
	}, util.AdminAuth())
	if !util.IsOk(setRes) {
		return util.NewResult(false, "Failed to set judging timer: "+setRes)
	}

	// Create a judge for judge auth
	judgeToken, result := createNamedJudge(context, "testForTimer@example.com", "Test For Timer")
	if !result.Success {
		return result
	}

	timerRes := util.GetRequest(context.Logger, "/admin/timer", util.JudgeAuth(judgeToken))
	timer := util.ExtractInt(timerRes, "judging_timer")
	if timer != 300 {
		return util.NewResult(false, fmt.Sprintf("Judging timer should be 300, got %d", timer))
	}

	return util.ResultOk()
}

// SetMinViews verifies that the min_views setting can be updated
func SetMinViews(context *util.Context) util.Result {
	setRes := util.PostRequest(context.Logger, "/admin/options", util.H{
		"min_views": 5,
	}, util.AdminAuth())
	if !util.IsOk(setRes) {
		return util.NewResult(false, "Failed to set min_views: "+setRes)
	}

	optRes := util.GetRequest(context.Logger, "/admin/options", util.AdminAuth())
	if !util.IsValue(optRes, "min_views", util.IntType, 5) {
		return util.NewResult(false, "min_views should be 5 after setting: "+optRes)
	}

	return util.ResultOk()
}

// DeliberationToggle verifies deliberation mode can be toggled on and off
func DeliberationToggle(context *util.Context) util.Result {
	// Enable deliberation
	onRes := util.PostRequest(context.Logger, "/admin/deliberation", util.H{"start": true}, util.AdminAuth())
	if !util.IsOk(onRes) {
		return util.NewResult(false, "Failed to enable deliberation: "+onRes)
	}

	// Judges should see deliberation mode on
	token, result := createNamedJudge(context, "deliberation_test@example.com", "Deliberation Test Judge")
	if !result.Success {
		return result
	}
	delibRes := util.GetRequest(context.Logger, "/judge/deliberation", util.JudgeAuth(token))
	if !util.IsOk(delibRes) {
		return util.NewResult(false, "Judge should see deliberation mode as active")
	}

	// Disable deliberation
	offRes := util.PostRequest(context.Logger, "/admin/deliberation", util.H{"start": false}, util.AdminAuth())
	if !util.IsOk(offRes) {
		return util.NewResult(false, "Failed to disable deliberation: "+offRes)
	}

	delibRes2 := util.GetRequest(context.Logger, "/judge/deliberation", util.JudgeAuth(token))
	if util.IsOk(delibRes2) {
		return util.NewResult(false, "Deliberation mode should be off after disabling")
	}

	return util.ResultOk()
}

// BlockReqsToggle verifies the block_reqs setting can be toggled
func BlockReqsToggle(context *util.Context) util.Result {
	// Enable blocking
	enableRes := util.PostRequest(context.Logger, "/admin/block-reqs", util.H{"block_reqs": true}, util.AdminAuth())
	if !util.IsOk(enableRes) {
		return util.NewResult(false, "Failed to enable block_reqs: "+enableRes)
	}

	optRes := util.GetRequest(context.Logger, "/admin/options", util.AdminAuth())
	if !util.IsValue(optRes, "block_reqs", util.BoolType, true) {
		return util.NewResult(false, "block_reqs should be true after enabling")
	}

	// Disable blocking (restore default)
	disableRes := util.PostRequest(context.Logger, "/admin/block-reqs", util.H{"block_reqs": false}, util.AdminAuth())
	if !util.IsOk(disableRes) {
		return util.NewResult(false, "Failed to disable block_reqs: "+disableRes)
	}

	return util.ResultOk()
}