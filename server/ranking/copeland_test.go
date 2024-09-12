package ranking_test

import (
	"fmt"
	"server/ranking"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCalcRanking(t *testing.T) {
	// Objects
	obj1 := primitive.NewObjectID()
	obj2 := primitive.NewObjectID()
	obj3 := primitive.NewObjectID()
	obj4 := primitive.NewObjectID()
	obj5 := primitive.NewObjectID()

	jrs := []ranking.JudgeRanking{
		{
			Rankings: []primitive.ObjectID{obj1, obj5, obj3},
			Unranked: []primitive.ObjectID{obj4, obj2},
		},
		{
			Rankings: []primitive.ObjectID{obj5, obj4, obj2},
			Unranked: []primitive.ObjectID{obj3, obj1},
		},
		{
			Rankings: []primitive.ObjectID{obj1, obj5, obj2},
			Unranked: []primitive.ObjectID{obj3, obj4},
		},
		{
			Rankings: []primitive.ObjectID{obj2, obj1, obj5},
			Unranked: []primitive.ObjectID{obj5, obj3},
		},
	}

	projects := []primitive.ObjectID{obj1, obj2, obj3, obj4, obj5}

	rankings := ranking.CalcCopelandRanking(jrs, projects)

	fmt.Println(rankings)
	fmt.Println(projects)
}
