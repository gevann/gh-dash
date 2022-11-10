package data

import (
	"strconv"
	"strings"

	"time"

	"github.com/cli/go-gh"
)

type RunData struct {
	Status   string
	Title    string
	Workflow string
	Branch   string
	Event    string
	Id       string
	Elapsed  string
	Age      string
}

func (data RunData) GetRepoNameWithOwner() string {
	return "foo/bar"
}

func (data RunData) GetNumber() int {
	return 0
}

func (data RunData) GetUrl() string {
	return "foo/bar"
}

func (data RunData) GetUpdatedAt() time.Time {
	return time.Now()
}

func ListRuns(limit *int) ([]RunData, error) {
	// Execute `gh run list`, and print the output.
	args := []string{"run", "list", "--limit", strconv.FormatInt(int64(*limit), 10)}
	stdOut, _, err := gh.Exec(args...)

	// convert the bytes.Buffer to a string
	output := stdOut.String()

	// split the string into a slice of strings
	lines := strings.Split(output, "\n")

	// remove the first line, which is the header
	lines = lines[1:]

	// remove the last line, which is empty
	lines = lines[:len(lines)-1]

	// create a slice of RunData
	runs := make([]RunData, len(lines))

	// loop through the lines
	for i, line := range lines {
		// split the line into a slice of strings
		fields := strings.Split(line, "\t")

		// create a RunData struct
		run := RunData{
			Status:   fields[0],
			Title:    fields[1],
			Workflow: fields[2],
			Branch:   fields[3],
			Event:    fields[4],
			Id:       fields[5],
			Elapsed:  fields[6],
			Age:      fields[7],
		}

		// append the RunData struct to the slice
		runs[i] = run
	}

	return runs, err
}
