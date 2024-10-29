package funcs

import (
	"server/models"
	"testing"

	"github.com/joho/godotenv"
)

func TestSendJudgeEmail(t *testing.T) {
	godotenv.Load("../.env")

	judge := models.NewJudge("Michael Zhao", "michaelzhao314@gmail.com", "notes here", -1)
	err := SendJudgeEmail(judge, "http://localhost:3000")
	if err != nil {
		t.Errorf("%s\n", err.Error())
		t.FailNow()
	}
}
