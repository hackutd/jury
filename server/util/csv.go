package util

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"server/database"
	"server/models"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

// Read CSV file and return a slice of judge structs
func ParseJudgeCSV(content string, hasHeader bool) ([]*models.Judge, error) {
	r := csv.NewReader(strings.NewReader(content))

	// Empty CSV file
	if content == "" {
		return []*models.Judge{}, nil
	}

	// If the CSV file has a header, skip the first line
	if hasHeader {
		r.Read()
	}

	// Read the CSV file, looping through each record
	var judges []*models.Judge
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Make sure the record has 3 elements (name, email, notes)
		if len(record) != 3 {
			return nil, fmt.Errorf("record does not contain 3 elements: '%s'", strings.Join(record, ","))
		}

		// Add judge to slice
		judges = append(judges, models.NewJudge(record[0], record[1], record[2]))
	}

	return judges, nil
}

// Read CSV file and return a slice of project structs
func ParseProjectCsv(content string, hasHeader bool, db *mongo.Database) ([]*models.Project, error) {
	r := csv.NewReader(strings.NewReader(content))

	// Empty CSV file
	if content == "" {
		return []*models.Project{}, nil
	}

	// Get the options from the database
	options, err := database.GetOptions(db)
	if err != nil {
		return nil, err
	}

	// If the CSV file has a header, skip the first line
	if hasHeader {
		r.Read()
	}

	// Read the CSV file, looping through each record
	var projects []*models.Project
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Make sure the record has at least 3 elements (name, description, URL)
		if len(record) < 3 {
			return nil, fmt.Errorf("record contains less than 3 elements: '%s'", strings.Join(record, ","))
		}

		// Get the challenge list
		challengeList := []string{}
		if len(record) > 5 && record[5] != "" {
			challengeList = strings.Split(record[5], ",")
		}
		for i := range challengeList {
			challengeList[i] = strings.TrimSpace(challengeList[i])
		}

		// Optional fields
		var tryLink string
		if len(record) > 3 && record[3] != "" {
			tryLink = record[3]
		}
		var videoLink string
		if len(record) > 4 && record[4] != "" {
			videoLink = record[4]
		}

		// Increment the table number
		database.GetNextTableNum(options)

		// Add project to slice
		projects = append(projects, models.NewProject(record[0], options.CurrTableNum, record[1], record[2], tryLink, videoLink, challengeList))
	}

	// Update the options table number in the database
	err = database.UpdateCurrTableNum(db, context.Background(), options.CurrTableNum)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

// Generate a workable CSV for Jury based on the output CSV from Devpost
// Columns:
//  0. Project Title - title
//  1. Submission Url - url
//  2. Project Status - Draft or Submitted (ignore drafts)
//  3. Judging Status - ignore
//  4. Highest Step Completed - ignore
//  5. Project Created At - ignore
//  6. About The Project - description
//  7. "Try it out" Links" - try_link
//  8. Video Demo Link - video_link
//  9. Opt-In Prizes - challenge_list
//  10. Built With - ignore
//  11. Notes - ignore
//  12. Team Colleges/Universities - ignore
//  13. Additional Team Member Count - ignore
//  14. (and remiaining rows) Custom questions - custom_questions (ignore for now)
func ParseDevpostCSV(content string, db *mongo.Database) ([]*models.Project, error) {
	r := csv.NewReader(strings.NewReader(content))

	// Empty CSV file
	if content == "" {
		return []*models.Project{}, nil
	}

	// Skip the first line
	r.Read()

	// Get the options from the database
	options, err := database.GetOptions(db)
	if err != nil {
		return nil, err
	}

	// Read the CSV file, looping through each record
	var projects []*models.Project
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Make sure the record has 14 or more elements (see above)
		if len(record) < 13 {
			return nil, fmt.Errorf("record does not contain 14 or more elements (invalid devpost csv): '%s'", strings.Join(record, ","))
		}

		// If the project is a Draft, skip it
		if record[2] == "Draft" {
			continue
		}

		// Split challenge list into a slice and trim them
		challengeList := strings.Split(record[9], ",")
		if record[9] == "" {
			challengeList = []string{}
		}
		for i := range challengeList {
			challengeList[i] = strings.TrimSpace(challengeList[i])
		}

		// Increment table number
		database.GetNextTableNum(options)

		// Add project to slice
		projects = append(projects, models.NewProject(
			record[0],
			options.CurrTableNum,
			record[6],
			record[1],
			record[7],
			record[8],
			challengeList,
		))
	}

	// Update the options table number in the database
	err = database.UpdateCurrTableNum(db, context.Background(), options.CurrTableNum)
	if err != nil {
		return nil, err
	}

	return projects, nil
}
