package tests

import (
	"encoding/json"
	"tests/util"
	"time"
)

// jsonUnmarshalList is a shared helper used across test files to parse JSON arrays
func jsonUnmarshalList(body string, v any) error {
	return json.Unmarshal([]byte(body), v)
}

type TestGroup struct {
	Name  string
	Setup func(*util.Context) util.Result
	Tests map[string]func(*util.Context) util.Result
}

func RunTests(context *util.Context) {
	tests := []TestGroup{
		{
			Name: "Smoke Tests",
			Tests: map[string]func(*util.Context) util.Result{
				"Heartbeat": Heartbeat,
				"Admin Login":     SmokeAdminLogin,
				"Admin Auth":      SmokeAdminAuth,
			},
		},
		{
			Name: "Auth & Access Control",
			Tests: map[string]func(*util.Context) util.Result{
				"Unauthenticated Requests Rejected":             UnauthenticatedRequestsRejected,
				"Invalid Judge Token Rejected":                  InvalidJudgeAuth,
				"Judge Token Cannot Access Admin Routes":        JudgeTokenCannotAccessAdminRoutes,
				"Admin Credentials Cannot Act As Judge":         AdminAuthWithJudgeEndpoint,
				"Empty Bearer Token Rejected":                   EmptyBearerTokenRejected,
				"Empty Bearer Token Cannot Access Judge Routes": EmptyBearerTokenCannotAccessJudgeRoutes,
				"Bearer Prefix Only Rejected":                   BearerPrefixOnlyRejected,
				"Uninitialized Judge Cannot Be Impersonated":    UninitializedJudgeCannotBeImpersonated,
			},
		},
		{
			Name: "Judge CRUD",
			Tests: map[string]func(*util.Context) util.Result{
				"Judge Login With Valid Code":               JudgeLoginWithValidCode,
				"Judge Login With Invalid Code":             JudgeLoginWithInvalidCode,
				"Add And Delete Judge":                      AddAndDeleteJudge,
				"Edit Judge":                               EditJudge,
				"Hide Judge":                               HideJudge,
				"Judge Welcome Flow":                       JudgeWelcomeFlow,
				"Judge Stats Reflect Additions":            JudgeStatsReflectAdditions,
				"QR Check Empty Code Rejected":              QRCheckEmptyCodeRejected,
				"QR Check Track Empty Code Rejected":        QRCheckTrackEmptyCodeRejected,
				"QR Add Empty Code Does Not Create Judge":   QRAddEmptyCodeDoesNotCreateJudge,
				"QR Add Garbage Code Does Not Create Judge": QRAddGarbageCodeDoesNotCreateJudge,
				"QR Valid Flow Still Works":                 QRValidFlowStillWorks,
			},
		},
		{
			Name: "Project CRUD",
			Tests: map[string]func(*util.Context) util.Result{
				"Add Project":                     AddProject,
				"Add And Delete Project":          AddAndDeleteProject,
				"Delete Non-Existent Project":     DeleteNonExistentProject,
				"Edit Project":                    EditProject,
				"Hide Project":                    HideProject,
				"Prioritize Project":              PrioritizeProject,
				"Move Project To Table":           MoveProject,
				"Public Project List Accessible":  PublicProjectListIsAccessible,
				"Project Stats Reflect Additions": ProjectStatsReflectAdditions,
			},
		},
		{
			Name: "Judging Workflow",
			Setup: JudgingTestSetup,
			Tests: map[string]func(*util.Context) util.Result{
				"Judging Standard Path":          JudgingStandardPath,
				"Judge Does Not Repeat Projects": JudgeDoesNotRepeatProjects,
				"Skip Project Creates Flag":      SkipProjectCreatesFlag,
				"Judge Rank Projects":            JudgeRankProjects,
				"Star Project":                   StarProject,
				"Judge Notes Update":             JudgeNotesUpdate,
				"Judge Next With No Projects":    JudgeNextWithNoActiveProjects,
			},
		},
		{
			Name: "Clock & Admin Settings",
			Tests: map[string]func(*util.Context) util.Result{
				"Clock Pause And Unpause":             ClockPauseUnpause,
				"Clock Reset":                         ClockReset,
				"Admin Started Reflects Clock":        AdminStartedReflectsClockState,
				"Get Options Returns Valid Structure": GetOptionsReturnsValidStructure,
				"Set And Get Judging Timer":           SetAndGetJudgingTimer,
				"Set Min Views":                       SetMinViews,
				"Deliberation Toggle":                 DeliberationToggle,
				"Block Reqs Toggle":                   BlockReqsToggle,
			},
		},
	}

	totalErrors := 0

	for _, group := range tests {
		context.Logger.Log(util.Info, "\n%s\n------------------\n", group.Name)
		context.Logger.Log(util.Info, "Date/time: %s\n", util.GetDateTime())
		context.Logger.Log(util.Info, "Found %d tests\n", len(group.Tests))
		context.Logger.Log(util.Info, "Running tests...\n")
		context.Logger.Log(util.Info, "------------------\n")

		errors := 0

		// Run the setup script if it exists
		// If it fails, assume we cannot continue this test set
		if group.Setup != nil {
			context.Logger.LogLn(util.Info, "Running setup")
			res := group.Setup(context)
			if !res.Success {
				context.Logger.Log(util.Error, "\tFAILED: %s\n", res.Message)
				context.Logger.LogLn(util.Info, "------------------")
				context.Logger.LogLn(util.Error, "Could not run tests due to failure of setup script")
				context.Logger.LogLn(util.Info, "------------------\n")
				continue
			}
			context.Logger.LogLn(util.Info, "\tSuccess")
			time.Sleep(100 * time.Millisecond)
		}

		for name, test := range group.Tests {
			context.Logger.Log(util.Info, "Running test: %s\n", name)
			res := test(context)
			if !res.Success {
				context.Logger.Log(util.Error, "\tFAILED: %s\n", res.Message)
				errors++
			} else {
				context.Logger.LogLn(util.Info, "\tPassed")
			}

			// Add delay between tests
			time.Sleep(100 * time.Millisecond)
		}

		context.Logger.LogLn(util.Info, "------------------")
		if errors == 0 {
			context.Logger.LogLn(util.Info, "All tests passed!")
		} else {
			context.Logger.Log(util.Error, "%d tests failed\n", errors)
		}
		context.Logger.LogLn(util.Info, "------------------\n")

		totalErrors += errors
	}

	if totalErrors == 0 {
		context.Logger.LogLn(util.Info, "Entire test suite passed!")
	} else {
		context.Logger.Log(util.Error, "There are %d failing tests across the test suite\n", totalErrors)
	}
}