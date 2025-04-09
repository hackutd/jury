package tests

import (
	"fmt"
	"tests/src"
)

func Heartbeat(context *src.Context) src.Result {
	res := src.GetRequest(context.Logger, "/", src.DefaultAuth())
	return src.AssertOk(res, "Error with heartbeat endpoint")
}

func InvalidAdminLogin(context *src.Context) src.Result {
	res := src.PostRequest(context.Logger, "/admin/login", src.H{"password": "THIS IS DEFINITELY THE WRONG PASSWORD"}, src.DefaultAuth())
	return src.AssertNotOk(res, "Invalid login to admin account should not be successful")
}

func AdminLogin(context *src.Context) src.Result {
	password := src.GetEnv("ADMIN_PASSWORD")
	res := src.PostRequest(context.Logger, "/admin/login", src.H{"password": password}, src.DefaultAuth())
	return src.AssertOk(res, "Error logging in as admin")
}

func InvalidAdminAuth(context *src.Context) src.Result {
	res := src.PostRequest(context.Logger, "/admin/auth", nil, "Bearer INVALID_TOKEN")
	fmt.Println(res)
	return src.AssertNotOk(res, "Invalid auth token should not be successful")
}

func AdminAuth(context *src.Context) src.Result {
	res := src.PostRequest(context.Logger, "/admin/auth", nil, src.AdminAuth())
	return src.AssertOk(res, "Error authenticating as admin")
}

func AddJudge(context *src.Context) src.Result {
	res := src.PostRequest(context.Logger, "/judge/new", src.H{"name": "Test Judge", "email": "test@gmail.com", "track": "", "notes": "test", "no_send": true}, src.AdminAuth())
	return src.AssertOk(res, "Error adding a judge")
}

func GetClock(context *src.Context) src.Result {
	res := src.GetRequest(context.Logger, "/admin/clock", src.AdminAuth())
	if !src.IsValue(res, "running", src.BoolType, false) || !src.IsValue(res, "time", src.Float64Type, 0.0) {
		return src.NewResult(false, "Error getting the clock")
	}
	return src.ResultOk()
}
