package tests

import "tests/src"

type TestGroup struct {
	Name  string
	Tests map[string]func(*src.Context) src.Result
}

func RunTests(context *src.Context) {
	// ADD TESTS AND GROUPS HERE
	tests := []TestGroup{
		{
			Name: "Endpoint Unit Tests",
			Tests: map[string]func(*src.Context) src.Result{
				"Heartbeat":            Heartbeat,
				"Invalid Admin Log In": InvalidAdminLogin,
				"Valid Admin Log In":   AdminLogin,
				"Invalid Admin Auth":   InvalidAdminAuth,
				"Valid Admin Auth":     AdminAuth,
				"Add Judge":            AddJudge,
				"Get Clock":            GetClock,
			},
		},
	}

	// Loop through each test group and run each test in the group
	for _, group := range tests {
		context.Logger.Log(src.Info, "\n%s\n------------------\n", group.Name)
		context.Logger.Log(src.Info, "Date/time: %s\n", src.GetDateTime())
		context.Logger.Log(src.Info, "Found %d tests\n", len(group.Tests))
		context.Logger.Log(src.Info, "Running tests...\n")
		context.Logger.Log(src.Info, "------------------\n")

		errors := 0

		for name, test := range group.Tests {
			context.Logger.Log(src.Info, "Running test: %s\n", name)
			res := test(context)
			if !res.Success {
				context.Logger.Log(src.Error, "\tFAILED: %s\n", res.Message)
				errors++
			} else {
				context.Logger.LogLn(src.Info, "\tPassed")
			}
		}

		context.Logger.LogLn(src.Info, "------------------")
		if errors == 0 {
			context.Logger.LogLn(src.Info, "All tests passed!")
		} else {
			context.Logger.Log(src.Error, "%d tests failed\n", errors)
		}
		context.Logger.LogLn(src.Info, "------------------\n")
	}
}
