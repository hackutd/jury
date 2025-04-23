package funcs

import (
	"errors"
	"server/database"
	"server/models"
	"sort"

	"go.mongodb.org/mongo-driver/mongo"
)

// ReassignNums assigns project numbers in order.
func ReassignNums(db *mongo.Database) error {
	err := database.WithTransaction(db, func(sc mongo.SessionContext) error {
		// Get all the projects from the database
		projects, err := database.FindAllProjects(db, sc)
		if err != nil {
			return errors.New("error getting projects from database: " + err.Error())
		}

		// If projects is empty, send OK
		if len(projects) == 0 {
			return nil
		}

		// Sort projets by table num
		sort.Sort(models.ByTableNumber(projects))

		// Set init table num to 0
		tableNum := int64(0)

		// Loop through all projects
		for _, project := range projects {
			tableNum++
			project.Location = tableNum
		}

		// Update all projects in the database
		err = database.UpdateProjects(db, sc, projects)
		if err != nil {
			return errors.New("error updating projects in database: " + err.Error())
		}
		return nil
	})

	return err
}

// IncrementJudgeGroupNum increments every single judges' group number.
func IncrementJudgeGroupNum(db *mongo.Database) error {
	return database.WithTransaction(db, func(sc mongo.SessionContext) error {
		// Get the options from the database
		options, err := database.GetOptions(db, sc)
		if err != nil {
			return errors.New("error getting options from database: " + err.Error())
		}

		// Get the judges from the database
		judges, err := database.FindAllJudges(db, sc)
		if err != nil {
			return errors.New("error getting judges from database: " + err.Error())
		}

		// Increment the group number for each judge
		for _, judge := range judges {
			judge.Group = (judge.Group + 1) % options.NumGroups
		}

		// Update all judges in the database
		// TODO: Can we write a function that will only update that field instead of have to pass ALL the judge data to the db :(
		err = database.UpdateJudgesWithTx(db, sc, judges)
		if err != nil {
			return errors.New("error updating judges in database: " + err.Error())
		}

		// Increment the manual switch count in the database
		err = database.IncrementManualSwitches(db, sc)
		if err != nil {
			return errors.New("error incrementing manual switches in database: " + err.Error())
		}

		return nil
	})
}
