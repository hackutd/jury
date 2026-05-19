package tests

import "tests/util"

// --- Auth / Access Control Tests ---

// JudgeTokenCannotAccessAdminRoutes verifies that a valid judge token is rejected on admin-only endpoints
func JudgeTokenCannotAccessAdminRoutes(context *util.Context) util.Result {
	// Add a judge and get their token
	token, result := createJudgeAndGetToken(context)
	if !result.Success {
		return result
	}

	judgeAuth := util.JudgeAuth(token)

	// Try hitting several admin-only endpoints with the judge token
	adminRoutes := []struct {
		method string
		url    string
	}{
		{"GET", "/judge/list"},
		{"GET", "/admin/stats"},
		{"GET", "/admin/flags"},
		{"GET", "/admin/options"},
	}

	for _, route := range adminRoutes {
		var status int
		if route.method == "GET" {
			status, _ = util.GetRequestWithStatus(context.Logger, route.url, judgeAuth)
		} else {
			status, _ = util.PostRequestWithStatus(context.Logger, route.url, nil, judgeAuth)
		}
		if status != 401 && status != 403 {
			return util.NewResult(false, "Judge token should not be able to access admin route "+route.url+" (got status "+statusStr(status)+")")
		}
	}

	return util.ResultOk()
}

// UnauthenticatedRequestsRejected verifies that admin routes reject requests with no auth at all
func UnauthenticatedRequestsRejected(context *util.Context) util.Result {
	adminRoutes := []string{
		"/judge/list",
		"/admin/stats",
		"/admin/flags",
		"/admin/options",
		"/admin/clock",
	}

	for _, url := range adminRoutes {
		status, _ := util.GetRequestWithStatus(context.Logger, url, util.DefaultAuth())
		if status != 401 && status != 403 {
			return util.NewResult(false, "Unauthenticated request should be rejected on "+url+" (got status "+statusStr(status)+")")
		}
	}

	return util.ResultOk()
}

// InvalidJudgeAuth verifies that a garbage judge token is rejected
func InvalidJudgeAuth(context *util.Context) util.Result {
	status, _ := util.GetRequestWithStatus(context.Logger, "/judge", util.JudgeAuth("totally-invalid-token-xyz"))
	if status == 200 {
		return util.NewResult(false, "Invalid judge token should not authenticate successfully")
	}
	return util.ResultOk()
}

// AdminAuthWithJudgeEndpoint verifies admin credentials work on /admin/auth but not as a judge token
func AdminAuthWithJudgeEndpoint(context *util.Context) util.Result {
	// Admin auth should work on /admin/auth
	res := util.PostRequest(context.Logger, "/admin/auth", nil, util.AdminAuth())
	if !util.IsOk(res) {
		return util.NewResult(false, "Admin auth should succeed on /admin/auth")
	}

	// Admin credentials should NOT work as a judge token on judge routes
	status, _ := util.GetRequestWithStatus(context.Logger, "/judge", util.AdminAuth())
	if status == 200 {
		return util.NewResult(false, "Admin credentials should not authenticate as a judge")
	}

	return util.ResultOk()
}

// helper: creates a throw-away judge via the admin API and returns their login token
func createJudgeAndGetToken(context *util.Context) (string, util.Result) {
	addRes := util.PostRequest(context.Logger, "/judge/new", util.H{
		"name":    "Auth Test Judge",
		"email":   "authtest@example.com",
		"track":   "",
		"notes":   "",
		"no_send": true,
	}, util.AdminAuth())
	if !util.IsOk(addRes) {
		return "", util.NewResult(false, "Failed to create judge for auth test: "+addRes)
	}

	// Get the judge list to find the login code
	listRes := util.GetRequest(context.Logger, "/judge/list", util.AdminAuth())
	code := extractJudgeCode(listRes, "authtest@example.com")
	if code == "" {
		return "", util.NewResult(false, "Could not find judge code in judge list")
	}

	loginRes := util.PostRequest(context.Logger, "/judge/login", util.H{"code": code}, util.DefaultAuth())
	token := util.ExtractString(loginRes, "token")
	if token == "" {
		return "", util.NewResult(false, "Judge login did not return a token")
	}

	return token, util.ResultOk()
}

// helper: converts a status int to a string for error messages
func statusStr(status int) string {
	switch status {
	case 200:
		return "200 OK"
	case 401:
		return "401 Unauthorized"
	case 403:
		return "403 Forbidden"
	case 404:
		return "404 Not Found"
	default:
		return "unknown"
	}
}

// EmptyBearerTokenRejected verifies that "Authorization: Bearer " (empty token
// after the prefix) is rejected outright, even when uninitialized judges exist
// in the database with token "".
func EmptyBearerTokenRejected(context *util.Context) util.Result {
	// Ensure at least one judge exists with an uninitialized (empty) token
	// by adding a judge but NOT logging them in.
	util.PostRequest(context.Logger, "/judge/new", util.H{
		"name":    "Uninitialized Token Judge",
		"email":   "uninit_token@example.com",
		"track":   "",
		"notes":   "",
		"no_send": true,
	}, util.AdminAuth())

	// "Bearer " — exactly 7 characters, empty string after slicing prefix
	status, _ := util.GetRequestWithStatus(context.Logger, "/judge", "Bearer ")
	if status == 200 {
		return util.NewResult(false, "\"Bearer \" (empty token) must not authenticate — matched uninitialized judge")
	}
	return util.ResultOk()
}

// EmptyBearerTokenCannotAccessJudgeRoutes checks the bypass across multiple
// judge-authenticated endpoints, not just GET /judge.
// Basically a check for passing an empty token into the Bearer auth.
func EmptyBearerTokenCannotAccessJudgeRoutes(context *util.Context) util.Result {
	// Add an uninitialized judge to guarantee token "" exists in Mongo
	util.PostRequest(context.Logger, "/judge/new", util.H{
		"name":    "Uninitialized Token Judge 2",
		"email":   "uninit_token2@example.com",
		"track":   "",
		"notes":   "",
		"no_send": true,
	}, util.AdminAuth())

	emptyBearerAuth := "Bearer "

	judgeRoutes := []struct {
		method string
		url    string
	}{
		{"GET", "/judge"},
		{"GET", "/judge/welcome"},
		{"GET", "/judge/projects"},
		{"POST", "/judge/next"},
		{"GET", "/judge/deliberation"},
	}

	for _, route := range judgeRoutes {
		var status int
		if route.method == "GET" {
			status, _ = util.GetRequestWithStatus(context.Logger, route.url, emptyBearerAuth)
		} else {
			status, _ = util.PostRequestWithStatus(context.Logger, route.url, nil, emptyBearerAuth)
		}
		if status == 200 {
			return util.NewResult(false, "Empty Bearer token must not grant access to "+route.url+"")
		}
	}

	return util.ResultOk()
}

// BearerPrefixOnlyRejected checks that "Bearer" with no space or token is also
// rejected — guards against off-by-one variants of the same class of bug.
func BearerPrefixOnlyRejected(context *util.Context) util.Result {
	status, _ := util.GetRequestWithStatus(context.Logger, "/judge", "Bearer")
	if status == 200 {
		return util.NewResult(false, "\"Bearer\" with no token must not authenticate")
	}
	return util.ResultOk()
}

// UninitializedJudgeCannotBeImpersonated verifies that a judge added but never
// logged in cannot be accessed by any near-empty token variant.
func UninitializedJudgeCannotBeImpersonated(context *util.Context) util.Result {
	util.PostRequest(context.Logger, "/judge/new", util.H{
		"name":    "Never Logged In Judge",
		"email":   "never_login@example.com",
		"track":   "",
		"notes":   "",
		"no_send": true,
	}, util.AdminAuth())

	suspiciousTokens := []string{
		"Bearer ",
		"Bearer  ",    // two spaces
		"Bearer\t",    // tab
		"Bearer null",
		"Bearer undefined",
	}

	for _, authHeader := range suspiciousTokens {
		status, _ := util.GetRequestWithStatus(context.Logger, "/judge", authHeader)
		if status == 200 {
			return util.NewResult(false, "Token \""+authHeader+"\" must not authenticate as an uninitialized judge")
		}
	}

	return util.ResultOk()
}