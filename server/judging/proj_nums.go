package judging

import (
	"errors"
	"server/database"
	"server/models"
	"sort"

	"go.mongodb.org/mongo-driver/mongo"
)

// ReassignNumsInOrder assigns project numbers in order.
func ReassignNumsInOrder(db *mongo.Database) error {
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

		// Get the options from the database
		options, err := database.GetOptions(db, sc)
		if err != nil {
			return errors.New("error getting options from database: " + err.Error())
		}

		// Set init table num to 0
		options.CurrTableNum = 0

		// Loop through all projects
		for _, project := range projects {
			project.Location = options.GetNextIncrTableNum()
		}

		// Update the options in the database
		err = database.UpdateCurrTableNum(db, sc, options.CurrTableNum)
		if err != nil {
			return errors.New("error updating options in database: " + err.Error())
		}

		// Update all projects in the database
		err = database.UpdateProjects(db, projects)
		if err != nil {
			return errors.New("error updating projects in database: " + err.Error())
		}
		return nil
	})

	return err
}

func ReassignNumsByGroup(db *mongo.Database) error {
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

		// Get the options from the database
		options, err := database.GetOptions(db, sc)
		if err != nil {
			return errors.New("error getting options from database: " + err.Error())
		}

		// Set init table num to 0
		options.CurrTableNum = 0

		// Loop through all projects
		for _, project := range projects {
			project.Location = options.GetNextGroupTableNum()
		}

		// Update the options in the database
		err = database.UpdateCurrTableNum(db, sc, options.CurrTableNum)
		if err != nil {
			return errors.New("error updating options in database: " + err.Error())
		}

		// Update all projects in the database
		err = database.UpdateProjects(db, projects)
		if err != nil {
			return errors.New("error updating projects in database: " + err.Error())
		}
		return nil
	})

	return err
}
