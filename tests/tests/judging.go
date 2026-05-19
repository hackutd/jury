package tests

import (
	"fmt"
	"tests/util"
)

// --- Judging Workflow Tests ---

// JudgingTestSetup will set the setting that may break a standard workflow
// This includes actually starting judging and disabling groups + tracks
func JudgingTestSetup(context *util.Context) util.Result {
	// Disable group/track judging
	setRes := util.PostRequest(context.Logger, "/admin/options", util.H{
		"judge_tracks": false,
		"multi_group": false,
	}, util.AdminAuth())
	if !util.IsOk(setRes) {
		return util.NewResult(false, "Failed to set judge-tracks and multi_group: "+setRes)
	}

	// Unpause clock (to make sure judging starts)
	unpauseRes := util.PostRequest(context.Logger, "/admin/clock/unpause", nil, util.AdminAuth())
	if !util.IsOk(unpauseRes) {
		return util.NewResult(false, "Failed to unpause clock: "+unpauseRes)
	}

	return util.ResultOk()
}

// JudgingStandardPath runs through the core judging loop:
// create judge + projects → login → get next → finish → verify seen counts increment
func JudgingStandardPath(context *util.Context) util.Result {
	// Delete all projects
	setRes := util.PostRequest(context.Logger, "/admin/reset", util.H{
		"type": "projects",
	}, util.AdminAuth())
	if !util.IsOk(setRes) {
		return util.NewResult(false, "Failed to set judge-tracks and multi_group: "+setRes)
	}

	// Add two projects
	for i := 1; i <= 2; i++ {
		addRes := util.PostRequest(context.Logger, "/project/new", util.H{
			"name":           fmt.Sprintf("Judging Project %d", i),
			"description":    "Test project for judging flow",
			"url":            "https://example.com",
			"try_link":       "",
			"video_link":     "",
			"challenge_list": "",
		}, util.AdminAuth())
		if !util.IsOk(addRes) {
			return util.NewResult(false, fmt.Sprintf("Failed to add judging project %d: %s", i, addRes))
		}
	}

	// Create a judge and log in
	token, result := createNamedJudge(context, "judging_flow@example.com", "Judging Flow Judge")
	if !result.Success {
		return result
	}
	auth := util.JudgeAuth(token)

	// Get project count before
	statsBefore := util.GetRequest(context.Logger, "/project/stats", util.AdminAuth())
	avgSeenBefore := util.ExtractInt(statsBefore, "avg_seen")

	// Get next project
	nextRes := util.PostRequest(context.Logger, "/judge/next", nil, auth)
	projectID := util.ExtractString(nextRes, "project_id")
	if projectID == "" {
		return util.NewResult(false, "GET /judge/next did not return a project_id: "+nextRes)
	}

	// Finish judging that project
	finishRes := util.PostRequest(context.Logger, "/judge/finish", util.H{
		"notes":   "Looks great",
		"starred": false,
	}, auth)
	if !util.IsOk(finishRes) {
		return util.NewResult(false, "POST /judge/finish failed: "+finishRes)
	}

	// Verify avg_seen increased (or at minimum, seen count on the judge went up)
	judgeRes := util.GetRequest(context.Logger, "/judge", auth)
	seen := util.ExtractInt(judgeRes, "seen")
	if seen < 1 {
		return util.NewResult(false, fmt.Sprintf("Judge 'seen' count should be at least 1 after judging, got %d", seen))
	}

	_ = avgSeenBefore // avoid unused variable error; avg_seen may not change with just one judge

	return util.ResultOk()
}

// JudgeDoesNotRepeatProjects verifies a judge never gets the same project twice
func JudgeDoesNotRepeatProjects(context *util.Context) util.Result {
	// Delete all projects
	setRes := util.PostRequest(context.Logger, "/admin/reset", util.H{
		"type": "projects",
	}, util.AdminAuth())
	if !util.IsOk(setRes) {
		return util.NewResult(false, "Failed to set judge-tracks and multi_group: "+setRes)
	}

	// Add 3 projects
	for i := 1; i <= 3; i++ {
		util.PostRequest(context.Logger, "/project/new", util.H{
			"name":           fmt.Sprintf("No Repeat Project %d", i),
			"description":    "Uniqueness test",
			"url":            "https://example.com",
			"try_link":       "",
			"video_link":     "",
			"challenge_list": "",
		}, util.AdminAuth())
	}

	token, result := createNamedJudge(context, "no_repeat@example.com", "No Repeat Judge")
	if !result.Success {
		return result
	}
	auth := util.JudgeAuth(token)

	seen := map[string]bool{}

	for i := 0; i < 3; i++ {
		nextRes := util.PostRequest(context.Logger, "/judge/next", nil, auth)
		projectID := util.ExtractString(nextRes, "project_id")
		if projectID == "" {
			// No more projects — that's fine, stop early
			break
		}
		if seen[projectID] {
			return util.NewResult(false, fmt.Sprintf("Judge was assigned the same project twice: %s", projectID))
		}
		seen[projectID] = true

		finishRes := util.PostRequest(context.Logger, "/judge/finish", util.H{
			"notes":   "",
			"starred": false,
		}, auth)
		if !util.IsOk(finishRes) {
			return util.NewResult(false, "POST /judge/finish failed on iteration "+fmt.Sprint(i))
		}
	}

	return util.ResultOk()
}

// SkipProjectCreatesFlag verifies that skipping a project records a flag visible to admin
func SkipProjectCreatesFlag(context *util.Context) util.Result {
	// Delete all projects
	setRes := util.PostRequest(context.Logger, "/admin/reset", util.H{
		"type": "projects",
	}, util.AdminAuth())
	if !util.IsOk(setRes) {
		return util.NewResult(false, "Failed to set judge-tracks and multi_group: "+setRes)
	}

	// Add a project to skip
	util.PostRequest(context.Logger, "/project/new", util.H{
		"name":           "Skip Test Project",
		"description":    "Will be skipped",
		"url":            "https://example.com",
		"try_link":       "",
		"video_link":     "",
		"challenge_list": "",
	}, util.AdminAuth())

	token, result := createNamedJudge(context, "skip_test@example.com", "Skip Test Judge")
	if !result.Success {
		return result
	}
	auth := util.JudgeAuth(token)

	// Get count of flags before
	flagsBefore := util.GetRequest(context.Logger, "/admin/flags", util.AdminAuth())
	countBefore := countFlags(flagsBefore)

	// Get next project
	nextRes := util.PostRequest(context.Logger, "/judge/next", nil, auth)
	projectID := util.ExtractString(nextRes, "project_id")
	if projectID == "" {
		return util.NewResult(false, "No project returned for skip test: "+nextRes)
	}

	// Skip the project
	skipRes := util.PostRequest(context.Logger, "/judge/skip", util.H{
		"reason": "absent",
	}, auth)
	if !util.IsOk(skipRes) {
		return util.NewResult(false, "POST /judge/skip failed: "+skipRes)
	}

	// Verify flag count increased
	flagsAfter := util.GetRequest(context.Logger, "/admin/flags", util.AdminAuth())
	countAfter := countFlags(flagsAfter)

	if countAfter <= countBefore {
		return util.NewResult(false, fmt.Sprintf("Flag count should increase after skip: before=%d, after=%d", countBefore, countAfter))
	}

	return util.ResultOk()
}

// JudgeRankProjects verifies that submitting a ranking succeeds
func JudgeRankProjects(context *util.Context) util.Result {
	// Delete all projects
	setRes := util.PostRequest(context.Logger, "/admin/reset", util.H{
		"type": "projects",
	}, util.AdminAuth())
	if !util.IsOk(setRes) {
		return util.NewResult(false, "Failed to set judge-tracks and multi_group: "+setRes)
	}

	// Add projects
	var projectIDs []string
	for i := 1; i <= 2; i++ {
		name := fmt.Sprintf("Rank Test Project %d", i)
		util.PostRequest(context.Logger, "/project/new", util.H{
			"name":           name,
			"description":    "For ranking",
			"url":            "https://example.com",
			"try_link":       "",
			"video_link":     "",
			"challenge_list": "",
		}, util.AdminAuth())

		id, result := findProjectIDByName(context, name)
		if !result.Success {
			return result
		}
		projectIDs = append(projectIDs, id)
	}

	token, result := createNamedJudge(context, "rank_test@example.com", "Rank Test Judge")
	if !result.Success {
		return result
	}
	auth := util.JudgeAuth(token)

	// Judge must see projects before ranking — do 2 finish cycles
	for i := 0; i < 2; i++ {
		nextRes := util.PostRequest(context.Logger, "/judge/next", nil, auth)
		if util.ExtractString(nextRes, "project_id") == "" {
			break
		}
		util.PostRequest(context.Logger, "/judge/finish", util.H{"notes": "", "starred": false}, auth)
	}

	// Submit a ranking
	rankRes := util.PostRequest(context.Logger, "/judge/rank", util.H{
		"ranking": projectIDs,
	}, auth)
	if !util.IsOk(rankRes) {
		return util.NewResult(false, "POST /judge/rank failed: "+rankRes)
	}

	return util.ResultOk()
}

// StarProject verifies a judge can star and unstar a project
func StarProject(context *util.Context) util.Result {
	// Delete all projects
	setRes := util.PostRequest(context.Logger, "/admin/reset", util.H{
		"type": "projects",
	}, util.AdminAuth())
	if !util.IsOk(setRes) {
		return util.NewResult(false, "Failed to set judge-tracks and multi_group: "+setRes)
	}

	util.PostRequest(context.Logger, "/project/new", util.H{
		"name":           "Star Test Project",
		"description":    "Will be starred",
		"url":            "https://example.com",
		"try_link":       "",
		"video_link":     "",
		"challenge_list": "",
	}, util.AdminAuth())

	token, result := createNamedJudge(context, "star_test@example.com", "Star Test Judge")
	if !result.Success {
		return result
	}
	auth := util.JudgeAuth(token)

	nextRes := util.PostRequest(context.Logger, "/judge/next", nil, auth)
	projectID := util.ExtractString(nextRes, "project_id")
	if projectID == "" {
		return util.NewResult(false, "No project returned for star test")
	}

	// Finish with starred = true
	finishRes := util.PostRequest(context.Logger, "/judge/finish", util.H{
		"notes":   "",
		"starred": true,
	}, auth)
	if !util.IsOk(finishRes) {
		return util.NewResult(false, "POST /judge/finish (starred) failed: "+finishRes)
	}

	// Verify via /judge/projects that starred is true
	projsRes := util.GetRequest(context.Logger, "/judge/projects", auth)
	starred := findSeenProjectField(projsRes, projectID, "starred")
	if starred != "true" {
		return util.NewResult(false, fmt.Sprintf("Project should be starred, got '%s'", starred))
	}

	return util.ResultOk()
}

// JudgeNotesUpdate verifies that notes can be updated for a seen project
func JudgeNotesUpdate(context *util.Context) util.Result {
	// Delete all projects
	setRes := util.PostRequest(context.Logger, "/admin/reset", util.H{
		"type": "projects",
	}, util.AdminAuth())
	if !util.IsOk(setRes) {
		return util.NewResult(false, "Failed to set judge-tracks and multi_group: "+setRes)
	}

	util.PostRequest(context.Logger, "/project/new", util.H{
		"name":           "Notes Test Project",
		"description":    "Will have notes",
		"url":            "https://example.com",
		"try_link":       "",
		"video_link":     "",
		"challenge_list": "",
	}, util.AdminAuth())

	token, result := createNamedJudge(context, "notes_test@example.com", "Notes Test Judge")
	if !result.Success {
		return result
	}
	auth := util.JudgeAuth(token)

	nextRes := util.PostRequest(context.Logger, "/judge/next", nil, auth)
	projectID := util.ExtractString(nextRes, "project_id")
	if projectID == "" {
		return util.NewResult(false, "No project returned for notes test")
	}

	util.PostRequest(context.Logger, "/judge/finish", util.H{"notes": "initial note", "starred": false}, auth)

	// Update notes
	notesRes := util.PutRequest(context.Logger, "/judge/notes/"+projectID, util.H{"notes": "updated note"}, auth)
	if !util.IsOk(notesRes) {
		return util.NewResult(false, "PUT /judge/notes/:id failed: "+notesRes)
	}

	// Verify
	projsRes := util.GetRequest(context.Logger, "/judge/projects", auth)
	notes := findSeenProjectField(projsRes, projectID, "notes")
	if notes != "updated note" {
		return util.NewResult(false, fmt.Sprintf("Notes not updated: got '%s'", notes))
	}

	return util.ResultOk()
}

// JudgeNextWithNoProjects verifies the API handles the case gracefully (no 500)
func JudgeNextWithNoActiveProjects(context *util.Context) util.Result {
	// Delete all projects
	setRes := util.PostRequest(context.Logger, "/admin/reset", util.H{
		"type": "projects",
	}, util.AdminAuth())
	if !util.IsOk(setRes) {
		return util.NewResult(false, "Failed to set judge-tracks and multi_group: "+setRes)
	}

	// Create a fresh judge
	token, result := createNamedJudge(context, "empty_test@example.com", "Empty Test Judge")
	if !result.Success {
		return result
	}
	auth := util.JudgeAuth(token)

	status, _ := util.PostRequestWithStatus(context.Logger, "/judge/next", nil, auth)
	if status == 500 {
		return util.NewResult(false, "GET /judge/next with no projects should not return 500")
	}

	return util.ResultOk()
}

// --- Helpers ---

// countFlags counts the number of flag objects in a JSON array response
func countFlags(body string) int {
	var flags []map[string]any
	if err := jsonUnmarshalList(body, &flags); err != nil {
		return 0
	}
	return len(flags)
}

// findSeenProjectField finds a field in the /judge/projects list by project_id
func findSeenProjectField(body string, projectID string, field string) string {
	var projects []map[string]any
	if err := jsonUnmarshalList(body, &projects); err != nil {
		return ""
	}
	for _, p := range projects {
		id, _ := p["project_id"].(string)
		if id == projectID {
			if val, ok := p[field]; ok {
				return fmt.Sprintf("%v", val)
			}
		}
	}
	return ""
}