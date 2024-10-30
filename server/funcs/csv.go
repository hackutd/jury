package funcs

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"server/database"
	"server/models"
	"server/util"
	"strings"

	"github.com/gin-gonic/gin"
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

		// Make sure the record has 2-3 elements (name, email, notes [optional])
		notes := ""
		if len(record) >= 3 {
			notes = record[2]
		} else if len(record) < 2 {
			return nil, fmt.Errorf("record does not contain 2-3 (name, email, notes [optional]) elements: '%s'", strings.Join(record, ","))
		}

		// Add judge to slice
		judges = append(judges, models.NewJudge(record[0], record[1], notes, -1))
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
	options, err := database.GetOptions(db, context.Background())
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
		group, num := GetNextTableNum(db, options)

		// Add project to slice
		projects = append(projects, models.NewProject(record[0], num, group, record[1], record[2], tryLink, videoLink, challengeList))
	}

	// Update the options table number in the database
	err = database.UpdateCurrTableNum(db, context.Background(), options)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

// TODO: After event, devpost will add a column between 0 and 1, the "auto assigned table numbers" TT - idk what to do abt this
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
	options, err := database.GetOptions(db, context.Background())
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
		group, num := GetNextTableNum(db, options)

		// Add project to slice
		projects = append(projects, models.NewProject(
			record[0],
			num,
			group,
			record[6],
			record[1],
			record[7],
			record[8],
			challengeList,
		))
	}

	// Update the options table number in the database
	err = database.UpdateCurrTableNum(db, context.Background(), options)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

// AddCSVData adds a CSV file to the response
func AddCsvData(name string, content []byte, ctx *gin.Context) {
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", name))
	ctx.Header("Content-Type", "text/csv")
	ctx.Data(http.StatusOK, "text/csv", content)
}

// AddZipFile adds a zip file to the response
func AddZipFile(name string, content []byte, ctx *gin.Context) {
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", name))
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Data(http.StatusOK, "application/octet-stream", content)
}

// Create a CSV file from a list of judges
func CreateJudgeCSV(judges []*models.Judge) []byte {
	csvBuffer := &bytes.Buffer{}

	// Create a new CSV writer
	w := csv.NewWriter(csvBuffer)

	// Write the header
	// TODO: Add judge rankings to output
	w.Write([]string{"Name", "Email", "Notes", "Code", "Active", "ReadWelcome", "Seen", "LastActivity"})

	// Write each judge
	for _, judge := range judges {
		w.Write([]string{judge.Name, judge.Email, judge.Notes, judge.Code, fmt.Sprintf("%t", judge.Active), fmt.Sprintf("%t", judge.ReadWelcome), fmt.Sprintf("%d", judge.Seen), fmt.Sprintf("%d", judge.LastActivity)})
	}

	// Flush the writer
	w.Flush()

	return csvBuffer.Bytes()
}

// Create a CSV file from the judges but only the rankings
func CreateJudgeRankingCSV(judges []*models.Judge) []byte {
	csvBuffer := &bytes.Buffer{}

	// Create a new CSV writer
	w := csv.NewWriter(csvBuffer)

	// Write the header
	w.Write([]string{"Name", "Code", "Ranked", "Unranked"})

	// Write each judge
	for _, judge := range judges {
		// Don't include if their rankings are empty :/
		if len(judge.Rankings) == 0 {
			continue
		}

		// Create a list of all ranked projects (just their location)
		ranked := make([]int64, 0, len(judge.Rankings))
		for _, projId := range judge.Rankings {
			idx := util.IndexFunc(judge.SeenProjects, func(p models.JudgedProject) bool {
				return p.ProjectId == projId
			})
			if idx == -1 {
				continue
			}

			ranked = append(ranked, judge.SeenProjects[idx].Location)
		}

		// Create a list of all unranked projects (filter using ranked projects)
		unranked := make([]int64, 0, len(judge.SeenProjects)-len(judge.Rankings))
		for _, proj := range judge.SeenProjects {
			if util.ContainsFunc(ranked, func(table int64) bool { return table == proj.Location }) {
				unranked = append(unranked, proj.Location)
			}
		}

		// Convert arrays to strings
		rankedStr := util.IntToString(ranked)
		unrankedStr := util.IntToString(unranked)

		// Write line to CSV
		w.Write([]string{judge.Name, judge.Code, strings.Join(rankedStr, ","), strings.Join(unrankedStr, ",")})
	}

	// Flush the writer
	w.Flush()

	return csvBuffer.Bytes()
}

// Create a CSV file from a list of projects
func CreateProjectCSV(projects []*models.Project) []byte {
	csvBuffer := &bytes.Buffer{}

	// Create a new CSV writer
	w := csv.NewWriter(csvBuffer)

	// Write the header
	w.Write([]string{"Name", "Table", "Description", "URL", "TryLink", "VideoLink", "ChallengeList", "Seen", "Active", "LastActivity"})

	// Write each project
	for _, project := range projects {
		w.Write([]string{project.Name, fmt.Sprintf("Table %d", project.Location), project.Description, project.Url, project.TryLink, project.VideoLink, strings.Join(project.ChallengeList, ","), fmt.Sprintf("%d", project.Seen), fmt.Sprintf("%t", project.Active), fmt.Sprintf("%d", project.LastActivity)})
	}

	// Flush the writer
	w.Flush()

	return csvBuffer.Bytes()
}

// CreateProjectChallengeZip creates a zip file with a CSV for each challenge
func CreateProjectChallengeZip(projects []*models.Project) ([]byte, error) {
	csvList := [][]byte{}

	// Get list of challenges
	challengeList := []string{}
	for _, project := range projects {
		for _, challenge := range project.ChallengeList {
			if !contains(challengeList, challenge) {
				challengeList = append(challengeList, challenge)
			}
		}
	}

	// Create a CSV for each challenge
	for _, challenge := range challengeList {
		currChallengeProjects := []*models.Project{}
		for _, project := range projects {
			if contains(project.ChallengeList, challenge) {
				currChallengeProjects = append(currChallengeProjects, project)
			}
		}

		// Create CSV for the challenge
		csv := CreateProjectCSV(currChallengeProjects)
		csvList = append(csvList, csv)
	}

	// Create buffer for zip file
	zipBuffer := &bytes.Buffer{}

	// Create a new zip writer
	w := zip.NewWriter(zipBuffer)

	// Write each CSV to the zip file
	for i, csv := range csvList {
		f, err := w.Create(fmt.Sprintf("%s.csv", challengeList[i]))
		if err != nil {
			return nil, err
		}

		_, err = f.Write(csv)
		if err != nil {
			return nil, err
		}
	}

	// Close the zip writer
	err := w.Close()
	if err != nil {
		return nil, err
	}

	return zipBuffer.Bytes(), nil
}

// contains checks if a string is in a list of strings
func contains(list []string, str string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}
