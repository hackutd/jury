package util

import (
	"encoding/csv"
	"fmt"
	"io"
	"server/models"
	"strings"
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
