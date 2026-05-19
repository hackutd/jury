package tests

import (
	"encoding/json"
	"fmt"
	"tests/util"
)

// --- Project CRUD Tests ---

// AddProject verifies a project can be added successfully
func AddProject(context *util.Context) util.Result {
	res := util.PostRequest(context.Logger, "/project/new", util.H{
		"name":           "Test Project",
		"description":    "A test project",
		"url":            "https://example.com",
		"try_link":       "",
		"video_link":     "",
		"challenge_list": "",
	}, util.AdminAuth())
	return util.AssertOk(res, "Failed to add project: "+res)
}

// AddAndDeleteProject verifies a project can be created and then deleted
func AddAndDeleteProject(context *util.Context) util.Result {
	name := "Delete Me Project"
	addRes := util.PostRequest(context.Logger, "/project/new", util.H{
		"name":           name,
		"description":    "Project to be deleted",
		"url":            "https://example.com",
		"try_link":       "",
		"video_link":     "",
		"challenge_list": "",
	}, util.AdminAuth())
	if !util.IsOk(addRes) {
		return util.NewResult(false, "Failed to add project: "+addRes)
	}

	id, result := findProjectIDByName(context, name)
	if !result.Success {
		return result
	}

	delRes := util.DeleteRequest(context.Logger, "/project/"+id, util.AdminAuth())
	if !util.IsOk(delRes) {
		return util.NewResult(false, "Failed to delete project: "+delRes)
	}

	// Confirm it's gone
	_, result = findProjectIDByName(context, name)
	if result.Success {
		return util.NewResult(false, "Project still appears in list after deletion")
	}

	return util.ResultOk()
}

// DeleteNonExistentProject verifies that deleting a bogus ID doesn't cause a 500
func DeleteNonExistentProject(context *util.Context) util.Result {
	status, _ := util.DeleteRequestWithStatus(context.Logger, "/project/000000000000000000000000", util.AdminAuth())
	if status == 500 {
		return util.NewResult(false, "Deleting a non-existent project should not return 500")
	}
	return util.ResultOk()
}

// EditProject verifies project fields can be updated
func EditProject(context *util.Context) util.Result {
	name := "Edit Me Project"
	addRes := util.PostRequest(context.Logger, "/project/new", util.H{
		"name":           name,
		"description":    "Original description",
		"url":            "https://example.com",
		"try_link":       "",
		"video_link":     "",
		"challenge_list": "",
	}, util.AdminAuth())
	if !util.IsOk(addRes) {
		return util.NewResult(false, "Failed to add project: "+addRes)
	}

	id, result := findProjectIDByName(context, name)
	if !result.Success {
		return result
	}

	editRes := util.PutRequest(context.Logger, "/project/"+id, util.H{
		"name":           "Edited Project Name",
		"description":    "Updated description",
		"url":            "https://updated.example.com",
		"try_link":       "",
		"video_link":     "",
		"challenge_list": "",
	}, util.AdminAuth())
	if !util.IsOk(editRes) {
		return util.NewResult(false, "Failed to edit project: "+editRes)
	}

	// Verify the change
	listRes := util.GetRequest(context.Logger, "/project/list", util.AdminAuth())
	desc := findProjectFieldByName(listRes, "Edited Project Name", "description")
	if desc != "Updated description" {
		return util.NewResult(false, fmt.Sprintf("Project description not updated: got '%s'", desc))
	}

	return util.ResultOk()
}

// HideProject verifies that hiding a project marks it inactive and removes it from the public list
func HideProject(context *util.Context) util.Result {
	name := "Hide Me Project"
	addRes := util.PostRequest(context.Logger, "/project/new", util.H{
		"name":           name,
		"description":    "Will be hidden",
		"url":            "https://example.com",
		"try_link":       "",
		"video_link":     "",
		"challenge_list": "",
	}, util.AdminAuth())
	if !util.IsOk(addRes) {
		return util.NewResult(false, "Failed to add project: "+addRes)
	}

	id, result := findProjectIDByName(context, name)
	if !result.Success {
		return result
	}

	hideRes := util.PutRequest(context.Logger, "/project/hide/"+id, util.H{"hide": true}, util.AdminAuth())
	if !util.IsOk(hideRes) {
		return util.NewResult(false, "Failed to hide project: "+hideRes)
	}

	// Check that it appears in the project list as hidden
	projectListRes := util.GetRequest(context.Logger, "/project/list", util.AdminAuth())
	projectActive := findProjectFieldByName(projectListRes, name, "active")
	if projectActive != "false" {
		return util.NewResult(false, "Hidden project should have active field set to false")
	}

	return util.ResultOk()
}

// ProjectStatsReflectAdditions verifies that project stats update after adding a project
func ProjectStatsReflectAdditions(context *util.Context) util.Result {
	statsBefore := util.GetRequest(context.Logger, "/project/stats", util.AdminAuth())
	numBefore := util.ExtractInt(statsBefore, "num")

	addRes := util.PostRequest(context.Logger, "/project/new", util.H{
		"name":           "Stats Count Project",
		"description":    "Counting project",
		"url":            "https://example.com",
		"try_link":       "",
		"video_link":     "",
		"challenge_list": "",
	}, util.AdminAuth())
	if !util.IsOk(addRes) {
		return util.NewResult(false, "Failed to add project: "+addRes)
	}

	statsAfter := util.GetRequest(context.Logger, "/project/stats", util.AdminAuth())
	numAfter := util.ExtractInt(statsAfter, "num")

	if numAfter != numBefore+1 {
		return util.NewResult(false, fmt.Sprintf("Project count should increase by 1: before=%d, after=%d", numBefore, numAfter))
	}

	return util.ResultOk()
}

// PrioritizeProject verifies a project can be prioritized and the flag persists
func PrioritizeProject(context *util.Context) util.Result {
	name := "Prioritize Me Project"
	addRes := util.PostRequest(context.Logger, "/project/new", util.H{
		"name":           name,
		"description":    "Should be prioritized",
		"url":            "https://example.com",
		"try_link":       "",
		"video_link":     "",
		"challenge_list": "",
	}, util.AdminAuth())
	if !util.IsOk(addRes) {
		return util.NewResult(false, "Failed to add project: "+addRes)
	}

	id, result := findProjectIDByName(context, name)
	if !result.Success {
		return result
	}

	prioRes := util.PutRequest(context.Logger, "/project/prioritize/"+id, util.H{"prioritize": true}, util.AdminAuth())
	if !util.IsOk(prioRes) {
		return util.NewResult(false, "Failed to prioritize project: "+prioRes)
	}

	// Confirm prioritized is true in the admin list
	listRes := util.GetRequest(context.Logger, "/project/list", util.AdminAuth())
	prioritized := findProjectFieldByName(listRes, name, "prioritized")
	if prioritized != "true" {
		return util.NewResult(false, fmt.Sprintf("Project should be prioritized, got '%s'", prioritized))
	}

	return util.ResultOk()
}

// MoveProject verifies a project can be moved to a specific table number
func MoveProject(context *util.Context) util.Result {
	name := "Move Me Project"
	addRes := util.PostRequest(context.Logger, "/project/new", util.H{
		"name":           name,
		"description":    "Will be moved",
		"url":            "https://example.com",
		"try_link":       "",
		"video_link":     "",
		"challenge_list": "",
	}, util.AdminAuth())
	if !util.IsOk(addRes) {
		return util.NewResult(false, "Failed to add project: "+addRes)
	}

	id, result := findProjectIDByName(context, name)
	if !result.Success {
		return result
	}

	moveRes := util.PutRequest(context.Logger, "/project/move/"+id, util.H{"location": 999}, util.AdminAuth())
	if !util.IsOk(moveRes) {
		return util.NewResult(false, "Failed to move project: "+moveRes)
	}

	listRes := util.GetRequest(context.Logger, "/project/list", util.AdminAuth())
	location := findProjectFieldByName(listRes, name, "location")
	if location != "999" {
		return util.NewResult(false, fmt.Sprintf("Project location should be 999, got '%s'", location))
	}

	return util.ResultOk()
}

// PublicProjectListIsAccessible verifies the public endpoint works without auth
func PublicProjectListIsAccessible(context *util.Context) util.Result {
	status, _ := util.GetRequestWithStatus(context.Logger, "/project/list/public", util.DefaultAuth())
	if status != 200 {
		return util.NewResult(false, fmt.Sprintf("Public project list should be accessible without auth, got %d", status))
	}
	return util.ResultOk()
}

// --- Helpers ---

func findProjectIDByName(context *util.Context, name string) (string, util.Result) {
	listRes := util.GetRequest(context.Logger, "/project/list", util.AdminAuth())
	id := findProjectFieldByName(listRes, name, "id")
	if id == "" {
		return "", util.NewResult(false, "Could not find project with name '"+name+"'")
	}
	return id, util.ResultOk()
}

func findProjectFieldByName(body string, name string, field string) string {
var projects []map[string]any
	if err := json.Unmarshal([]byte(body), &projects); err != nil {
		return ""
	}
	for _, project := range projects {
		if project["name"] == name {
			if val, ok := project[field]; ok {
				return fmt.Sprintf("%v", val)
			}
		}
	}
	return ""
}