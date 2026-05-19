package tests

import (
	"encoding/json"
	"fmt"
	"tests/util"
)

// --- Judge CRUD Tests ---

// AddAndDeleteJudge verifies a judge can be created and then deleted, disappearing from the list
func AddAndDeleteJudge(context *util.Context) util.Result {
	// Add a unique judge
	email := "delete_test@example.com"
	addRes := util.PostRequest(context.Logger, "/judge/new", util.H{
		"name":    "Delete Test Judge",
		"email":   email,
		"track":   "",
		"notes":   "",
		"no_send": true,
	}, util.AdminAuth())
	if !util.IsOk(addRes) {
		return util.NewResult(false, "Failed to add judge: "+addRes)
	}

	// Find their ID
	id, result := findJudgeIDByEmail(context, email)
	if !result.Success {
		return result
	}

	// Delete the judge
	delRes := util.DeleteRequest(context.Logger, "/judge/"+id, util.AdminAuth())
	if !util.IsOk(delRes) {
		return util.NewResult(false, "Failed to delete judge: "+delRes)
	}

	// Confirm they no longer appear in the list
	_, result = findJudgeIDByEmail(context, email)
	if result.Success {
		return util.NewResult(false, "Judge still appears in list after deletion")
	}

	return util.ResultOk()
}

// EditJudge verifies that judge info can be updated and the changes persist
func EditJudge(context *util.Context) util.Result {
	email := "edit_test@example.com"
	addRes := util.PostRequest(context.Logger, "/judge/new", util.H{
		"name":    "Edit Test Judge",
		"email":   email,
		"track":   "",
		"notes":   "original notes",
		"no_send": true,
	}, util.AdminAuth())
	if !util.IsOk(addRes) {
		return util.NewResult(false, "Failed to add judge for edit test: "+addRes)
	}

	id, result := findJudgeIDByEmail(context, email)
	if !result.Success {
		return result
	}

	// Edit the judge
	editRes := util.PutRequest(context.Logger, "/judge/"+id, util.H{
		"name":  "Edited Judge Name",
		"email": email,
		"notes": "updated notes",
	}, util.AdminAuth())
	if !util.IsOk(editRes) {
		return util.NewResult(false, "Failed to edit judge: "+editRes)
	}

	// Verify the change
	listRes := util.GetRequest(context.Logger, "/judge/list", util.AdminAuth())
	name := findJudgeFieldByEmail(listRes, email, "name")
	if name != "Edited Judge Name" {
		return util.NewResult(false, fmt.Sprintf("Judge name not updated: expected 'Edited Judge Name', got '%s'", name))
	}

	return util.ResultOk()
}

// JudgeLoginWithValidCode verifies a judge can log in and receive a token
func JudgeLoginWithValidCode(context *util.Context) util.Result {
	email := "login_test@example.com"
	addRes := util.PostRequest(context.Logger, "/judge/new", util.H{
		"name":    "Login Test Judge",
		"email":   email,
		"track":   "",
		"notes":   "",
		"no_send": true,
	}, util.AdminAuth())
	if !util.IsOk(addRes) {
		return util.NewResult(false, "Failed to add judge: "+addRes)
	}

	listRes := util.GetRequest(context.Logger, "/judge/list", util.AdminAuth())
	code := extractJudgeCode(listRes, email)
	if code == "" {
		return util.NewResult(false, "Could not find judge code in list")
	}

	loginRes := util.PostRequest(context.Logger, "/judge/login", util.H{"code": code}, util.DefaultAuth())
	token := util.ExtractString(loginRes, "token")
	if token == "" {
		return util.NewResult(false, "Judge login did not return a token: "+loginRes)
	}

	return util.ResultOk()
}

// JudgeLoginWithInvalidCode verifies that a bad login code is rejected
func JudgeLoginWithInvalidCode(context *util.Context) util.Result {
	res := util.PostRequest(context.Logger, "/judge/login", util.H{"code": "DEFINITELY-NOT-A-REAL-CODE"}, util.DefaultAuth())
	return util.AssertNotOk(res, "Login with invalid code should not succeed")
}

// JudgeWelcomeFlow verifies the read_welcome flag can be set and queried
func JudgeWelcomeFlow(context *util.Context) util.Result {
	token, result := createNamedJudge(context, "welcome_test@example.com", "Welcome Test Judge")
	if !result.Success {
		return result
	}

	auth := util.JudgeAuth(token)

	// Initially read_welcome should be false
	getRes := util.GetRequest(context.Logger, "/judge/welcome", auth)
	if util.IsOk(getRes) {
		return util.NewResult(false, "read_welcome should initially be false (not ok)")
	}

	// Mark welcome as read
	putRes := util.PostRequest(context.Logger, "/judge/welcome", util.H{}, auth)
	if !util.IsOk(putRes) {
		return util.NewResult(false, "Failed to set read_welcome: "+putRes)
	}

	// Now it should return ok
	getRes2 := util.GetRequest(context.Logger, "/judge/welcome", auth)
	if !util.IsOk(getRes2) {
		return util.NewResult(false, "read_welcome should be true after PUT /judge/welcome")
	}

	return util.ResultOk()
}

// HideJudge verifies that hiding a judge marks them as inactive
func HideJudge(context *util.Context) util.Result {
	email := "hide_test@example.com"
	addRes := util.PostRequest(context.Logger, "/judge/new", util.H{
		"name":    "Hide Test Judge",
		"email":   email,
		"track":   "",
		"notes":   "",
		"no_send": true,
	}, util.AdminAuth())
	if !util.IsOk(addRes) {
		return util.NewResult(false, "Failed to add judge: "+addRes)
	}

	id, result := findJudgeIDByEmail(context, email)
	if !result.Success {
		return result
	}

	hideRes := util.PutRequest(context.Logger, "/judge/hide/"+id, util.H{"hide": true}, util.AdminAuth())
	if !util.IsOk(hideRes) {
		return util.NewResult(false, "Failed to hide judge: "+hideRes)
	}

	// Confirm active is now false
	listRes := util.GetRequest(context.Logger, "/judge/list", util.AdminAuth())
	active := findJudgeFieldByEmail(listRes, email, "active")
	if active != "false" {
		return util.NewResult(false, fmt.Sprintf("Judge should be inactive after hiding, got active=%s", active))
	}

	return util.ResultOk()
}

// JudgeStatsReflectAdditions checks that judge stats update after adding judges
func JudgeStatsReflectAdditions(context *util.Context) util.Result {
	// Get current count
	statsBefore := util.GetRequest(context.Logger, "/judge/stats", util.AdminAuth())
	numBefore := util.ExtractInt(statsBefore, "num")

	// Add a judge
	addRes := util.PostRequest(context.Logger, "/judge/new", util.H{
		"name":    "Stats Test Judge",
		"email":   "stats_judge@example.com",
		"track":   "",
		"notes":   "",
		"no_send": true,
	}, util.AdminAuth())
	if !util.IsOk(addRes) {
		return util.NewResult(false, "Failed to add judge: "+addRes)
	}

	statsAfter := util.GetRequest(context.Logger, "/judge/stats", util.AdminAuth())
	numAfter := util.ExtractInt(statsAfter, "num")

	if numAfter != numBefore+1 {
		return util.NewResult(false, fmt.Sprintf("Judge count should increase by 1: before=%d, after=%d", numBefore, numAfter))
	}

	return util.ResultOk()
}

// --- Helpers ---

// createNamedJudge creates a judge and returns their token
func createNamedJudge(context *util.Context, email string, name string) (string, util.Result) {
	addRes := util.PostRequest(context.Logger, "/judge/new", util.H{
		"name":    name,
		"email":   email,
		"track":   "",
		"notes":   "",
		"no_send": true,
	}, util.AdminAuth())
	if !util.IsOk(addRes) {
		return "", util.NewResult(false, "Failed to create judge '"+name+"': "+addRes)
	}

	listRes := util.GetRequest(context.Logger, "/judge/list", util.AdminAuth())
	code := extractJudgeCode(listRes, email)
	if code == "" {
		return "", util.NewResult(false, "Could not find code for judge: "+email)
	}

	loginRes := util.PostRequest(context.Logger, "/judge/login", util.H{"code": code}, util.DefaultAuth())
	token := util.ExtractString(loginRes, "token")
	if token == "" {
		return "", util.NewResult(false, "Judge login did not return a token for: "+email)
	}

	return token, util.ResultOk()
}

// findJudgeIDByEmail scans the judge list for a judge with the given email and returns their ID
func findJudgeIDByEmail(context *util.Context, email string) (string, util.Result) {
	listRes := util.GetRequest(context.Logger, "/judge/list", util.AdminAuth())
	id := findJudgeFieldByEmail(listRes, email, "id")
	if id == "" {
		return "", util.NewResult(false, "Could not find judge with email '"+email+"' in judge list")
	}
	return id, util.ResultOk()
}

// findJudgeFieldByEmail scans a judge list JSON body and returns the value of 'field' for the judge with the given email
func findJudgeFieldByEmail(body string, email string, field string) string {
	var judges []map[string]any
	if err := json.Unmarshal([]byte(body), &judges); err != nil {
		return ""
	}
	for _, judge := range judges {
		if judge["email"] == email {
			if val, ok := judge[field]; ok {
				return fmt.Sprintf("%v", val)
			}
		}
	}
	return ""
}

// extractJudgeCode returns the login code for a judge identified by email
func extractJudgeCode(body string, email string) string {
	return findJudgeFieldByEmail(body, email, "code")
}

// QRCheckEmptyCodeRejected verifies that /qr/check rejects an empty code.
// Before the fix, options.QRCode defaulted to "" so submitting "" satisfied
// the equality check on a fresh instance with no QR generated yet.
func QRCheckEmptyCodeRejected(context *util.Context) util.Result {
	res := util.PostRequest(context.Logger, "/qr/check", util.H{"code": ""}, util.DefaultAuth())
	return util.AssertNotOk(res, "Empty string should not pass QR code check")
}

// QRCheckTrackEmptyCodeRejected verifies that /qr/check/:track rejects an
// empty code. Track QR codes also default to "" and were equally bypassable.
func QRCheckTrackEmptyCodeRejected(context *util.Context) util.Result {
	res := util.PostRequest(context.Logger, "/qr/check/unset-track", util.H{"code": ""}, util.DefaultAuth())
	return util.AssertNotOk(res, "Empty string should not pass track QR code check")
}

// QRAddEmptyCodeDoesNotCreateJudge verifies that /qr/add rejects an empty
// code and does NOT create a judge. Checks judge count before and after to
// catch a silently-created judge even if the response body looks benign.
func QRAddEmptyCodeDoesNotCreateJudge(context *util.Context) util.Result {
	listBefore := util.GetRequest(context.Logger, "/judge/list", util.AdminAuth())
	countBefore := countJudgesInList(listBefore)

	res := util.PostRequest(context.Logger, "/qr/add", util.H{
		"name":  "Attacker",
		"email": "attacker@evil.com",
		"notes": "",
		"code":  "",
	}, util.DefaultAuth())

	if util.IsOk(res) {
		return util.NewResult(false, "POST /qr/add with empty code should be rejected")
	}

	listAfter := util.GetRequest(context.Logger, "/judge/list", util.AdminAuth())
	countAfter := countJudgesInList(listAfter)
	if countAfter > countBefore {
		return util.NewResult(false, "POST /qr/add with empty code must not create a judge — judge count increased")
	}

	return util.ResultOk()
}

// QRAddGarbageCodeDoesNotCreateJudge verifies that a non-empty but invalid
// code is also rejected — guards against a fix that only special-cases "".
func QRAddGarbageCodeDoesNotCreateJudge(context *util.Context) util.Result {
	listBefore := util.GetRequest(context.Logger, "/judge/list", util.AdminAuth())
	countBefore := countJudgesInList(listBefore)

	res := util.PostRequest(context.Logger, "/qr/add", util.H{
		"name":  "Attacker",
		"email": "attacker2@evil.com",
		"notes": "",
		"code":  "THIS-IS-NOT-A-REAL-QR-CODE",
	}, util.DefaultAuth())

	if util.IsOk(res) {
		return util.NewResult(false, "POST /qr/add with invalid code should be rejected")
	}

	listAfter := util.GetRequest(context.Logger, "/judge/list", util.AdminAuth())
	countAfter := countJudgesInList(listAfter)
	if countAfter > countBefore {
		return util.NewResult(false, "POST /qr/add with invalid code must not create a judge — judge count increased")
	}

	return util.ResultOk()
}

// QRValidFlowStillWorks verifies the legitimate QR signup flow still works
// after the fix: generate a code as admin, verify it, then use it to add a judge.
func QRValidFlowStillWorks(context *util.Context) util.Result {
	// POST /admin/qr returns {"qr_code":"..."} directly, not {"ok":1}
	genRes := util.PostRequest(context.Logger, "/admin/qr", nil, util.AdminAuth())
	qrCode := util.ExtractString(genRes, "qr_code")
	if qrCode == "" {
		return util.NewResult(false, "POST /admin/qr did not return a qr_code: "+genRes)
	}

	checkRes := util.PostRequest(context.Logger, "/qr/check", util.H{"code": qrCode}, util.DefaultAuth())
	if !util.IsOk(checkRes) {
		return util.NewResult(false, "Valid QR code should pass /qr/check: "+checkRes)
	}

	addRes := util.PostRequest(context.Logger, "/qr/add", util.H{
		"name":  "Legitimate QR Judge",
		"email": "qr_judge@example.com",
		"notes": "",
		"code":  qrCode,
	}, util.DefaultAuth())
	if !util.IsOk(addRes) {
		return util.NewResult(false, "POST /qr/add with valid code should succeed: "+addRes)
	}

	return util.ResultOk()
}

func countJudgesInList(body string) int {
	var judges []map[string]any
	if err := jsonUnmarshalList(body, &judges); err != nil {
		return 0
	}
	return len(judges)
}