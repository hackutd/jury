package tests

import "tests/util"

// Heartbeat verifies the backend is reachable and returning a healthy response.
func Heartbeat(context *util.Context) util.Result {
	res := util.GetRequest(context.Logger, "/", util.DefaultAuth())
	return util.AssertOk(res, "Error with heartbeat endpoint")
}

// SmokeAdminLogin verifies the admin login endpoint accepts the correct password.
// If this fails, every test that depends on admin auth will also fail.
func SmokeAdminLogin(context *util.Context) util.Result {
	res := util.PostRequest(context.Logger, "/admin/login", util.H{"password": util.GetEnv("ADMIN_PASSWORD")}, util.DefaultAuth())
	return util.AssertOk(res, "Admin login failed — check ADMIN_PASSWORD env var and backend config")
}

// SmokeAdminAuth verifies the Basic auth middleware is functioning correctly.
func SmokeAdminAuth(context *util.Context) util.Result {
	res := util.PostRequest(context.Logger, "/admin/auth", nil, util.AdminAuth())
	return util.AssertOk(res, "Admin Basic auth failed — middleware may be broken")
}