package funcs

import (
	"server/models"
	"testing"

	"github.com/joho/godotenv"
)

func TestSendJudgeEmail(t *testing.T) {
	godotenv.Load("../.env")

	judge := models.NewJudge("Michael Zhao", "michaelzhao314@gmail.com", "notes here")
	err := SendJudgeEmail(judge, "http://localhost:3000")
	if err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}
}
