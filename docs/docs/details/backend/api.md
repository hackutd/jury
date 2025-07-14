---
sidebar_position: 2
title: API
description: The private API that the frontend calls
---

# API

The API uses the Gin framework. The `server/router/init.go` file is the main powerhouse that defines all routes, creates shared resources, and attaches any middleware.

## Router Groups

Gin has a concept of Router groups, which let us attach different middleware and group routes together:

```go
judgeRouter := router.Group("/api", AuthenticateJudge())
adminRouter := router.Group("/api", AuthenticateAdmin())
defaultRouter := router.Group("/api")
```

The `judgeRouter` requires a judge token and will attach a `judge` object to the context. This middleware will also deny requests with invalid judge tokens. You can use the judge object like this example (`server/router/judge.go`):

```go
func CheckJudgeReadWelcome(ctx *gin.Context) {
	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Send OK
	if judge.ReadWelcome {
		ctx.JSON(http.StatusOK, gin.H{"ok": 1})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"ok": 0})
	}
}
```

Note that we use `ctx.MustGet` and cast it to a `Judge` object.

The `adminRouter` requires an admin password basic auth. This will not attach anything to the context but will deny requests that do not have the right admin password.

## Middleware

The other non-auth middleware consists of:

- **Shared State**: allows each route to use resources such as the database connnection and logger.
- **Limiter**: limits login requests (set in [admin options](/docs/usage/admin/configuration))
- **Logger**: adds function to state that lets you log requests to a global audit log
- **CORS**: uses to avoid pesky CORS issues

## Route Handlers

All handler functions for routes should be placed in the `server/router` package, in one of the files based on the most relevant resource it's modifying. The handler can then be called from the `server/router/init.go` file.

An example of a handler can be seen with `CheckJudgeReadWelcome`:

```go
// POST /judge/welcome - Endpoint to set a judge's readWelcome field to true
func SetJudgeReadWelcome(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Update judge in database
	err := database.UpdateJudgeReadWelcome(state.Db, ctx, &judge.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating judge in database: " + err.Error()})
		return
	}

	// Send OK
	state.Logger.JudgeLogf(judge, "Read welcome message")
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}
```

There are a few key elements to go over:

- All handler functions will have ONE parameter: `ctx *gin.Context`
- The state ctx variable contains all the useful parameters such as the database connection and the logger, which can be retrieved with `GetState(ctx)`
- Here we get the judge from the context, but this is only allowed in judge routes (using the `judgeRouter`)
- We call a database function (see the relavent function on the [database page](/docs/details/backend/database#inserting-intoupdating-the-database)) instead of defining it in this function--this allows for us to abstract DB code to the database module
- Before returning, we log [using the logger](/docs/details/backend/logging)
- Finally, we return the OK Response (`{"ok": 1}`), which can be 0 if we are say checking if a value is true or not
- If there is an error, we return an error object along with the relevant status code: `{"error": "error message"}`
