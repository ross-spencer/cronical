package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

const fixtures = "crontests"
const crontabFile = "cron-entries"

var crontabFixture string

// Errors anticipated from parsing the cron fixtures.
var errResults = [...]string{
	"Time pattern is not implemented: π",
	"Cron entry: '[* * * * *]' is null",
	// TODO: NOT IMPLEMENTED
	"Time pattern is not implemented: JAN",
	"Time pattern is not implemented: SUN",
}

// cron results anticipated from parsing the cron fixtures.
//
// Fields: m h dom mon dow command
//
var cronResult01 = []Cron{
	Cron{10, false, 20, false, 1, false, -1, false, -1, false, "echo \"test one (8 tasks)\""},
	Cron{15, false, 20, false, 1, false, -1, false, -1, false, "echo \"test one (8 tasks)\""},
	Cron{10, false, 16, false, 1, false, -1, false, -1, false, "echo \"test one (8 tasks)\""},
	Cron{15, false, 16, false, 1, false, -1, false, -1, false, "echo \"test one (8 tasks)\""},
	Cron{10, false, 20, false, 20, false, -1, false, -1, false, "echo \"test one (8 tasks)\""},
	Cron{15, false, 20, false, 20, false, -1, false, -1, false, "echo \"test one (8 tasks)\""},
	Cron{10, false, 16, false, 20, false, -1, false, -1, false, "echo \"test one (8 tasks)\""},
	Cron{15, false, 16, false, 20, false, -1, false, -1, false, "echo \"test one (8 tasks)\""},
}

var cronResult02 = []Cron{
	Cron{0, false, 2, false, -1, false, -1, false, 5, false, "echo \"test two (2 tasks)\""},
	Cron{0, false, 2, false, -1, false, -1, false, 6, false, "echo \"test two (2 tasks)\""},
}

var cronResult03 = []Cron{
	Cron{0, false, 0, false, 1, false, 1, false, -1, false, "echo \"test three happy new year! (1 task)\""},
}

var cronResult04 = []Cron{
	Cron{1, false, -1, false, -1, false, -1, false, -1, false, "echo \"test four (1 task)\""},
}

var cronResult05 = []Cron{
	Cron{15, true, -1, false, -1, false, -1, false, -1, false, "echo \"test five (1 task)\""},
}

var cronResult06 = []Cron{
	Cron{30, true, -1, false, -1, false, -1, false, -1, false, "echo \"test seven (1 task)\""},
}

var cronResult07 = []Cron{
	Cron{10, true, -1, false, -1, false, -1, false, -1, false, "echo \"test eight (1 task)\""},
}

var cronResult08 = []Cron{
	Cron{-1, false, 0, false, -1, false, -1, false, -1, false, "echo \"test nine (3 tasks)\""},
	Cron{-1, false, 2, false, -1, false, -1, false, -1, false, "echo \"test nine (3 tasks)\""},
	Cron{-1, false, 4, false, -1, false, -1, false, -1, false, "echo \"test nine (3 tasks)\""},
}

var cronResult09 = []Cron{
	Cron{-1, false, 1, true, -1, false, -1, false, -1, false, "echo \"test ten (1 task)\""},
}

var cronResult10 = []Cron{
	Cron{-1, false, -1, false, 31, false, 8, false, -1, false, "echo \"test eleven (6 tasks / 2 invalid)\""},
	Cron{-1, false, -1, false, 30, false, 8, false, -1, false, "echo \"test eleven (6 tasks / 2 invalid)\""},
	Cron{-1, false, -1, false, 31, false, 1, false, -1, false, "echo \"test eleven (6 tasks / 2 invalid)\""},
	Cron{-1, false, -1, false, 30, false, 1, false, -1, false, "echo \"test eleven (6 tasks / 2 invalid)\""},
	Cron{-1, false, -1, false, 31, false, 2, false, -1, false, "echo \"test eleven (6 tasks / 2 invalid)\""},
	Cron{-1, false, -1, false, 30, false, 2, false, -1, false, "echo \"test eleven (6 tasks / 2 invalid)\""},
}

var cronResult11 = []Cron{
	Cron{10, false, -1, false, -1, false, 10, false, -1, false, "echo \"test twelve (1 tasks)\""},
}

// TODO: NOT IMPLEMENTED

var cronResult12 = []Cron{
	Cron{0, false, 0, false, 1, false, 1, false, 1, false, "echo \"test fourteen happy new year! (1 task)\""},
}

var cronResult13 = []Cron{
	Cron{0, false, 0, false, 1, false, 1, false, 0, false, "echo \"test fifteen happy new year! (1 task)\""},
}

var allCron = [][]Cron{
	cronResult01,
	cronResult02,
	cronResult03,
	cronResult04,
	cronResult05,
	cronResult06,
	cronResult07,
	cronResult08,
	cronResult09,
	cronResult10,
	cronResult11,
}
var cronResults = []Cron{}

func init() {
	crontabFixture = path.Join(fixtures, crontabFile)
	for _, cron := range allCron {
		cronResults = append(cronResults, cron...)
	}

}

func TestParseAll(t *testing.T) {
	testCrontab, err := os.Open(crontabFixture)
	if err != nil {
		t.Errorf(fmt.Sprintf("%s\n", err))
	}
	defer testCrontab.Close()

	crontab, err := ioutil.ReadAll(testCrontab)
	if err != nil {
		t.Errorf(fmt.Sprintf("%s\n", err))
	}

	cronList, errList := ParseCrontab(string(crontab))

	// TEST RESULTS

	// Test that the number of results returned is correct.
	if len(cronList) != len(cronResults) {
		t.Errorf(
			"Entries to parse: '%d' is different than expected: '%d'",
			len(cronList),
			len(cronResults),
		)
	}

	for _, cronL := range cronList {
		present := false
		for _, cronR := range cronResults {
			if cronL == cronR {
				present = true
				break
			}
		}
		if !present {
			t.Errorf("Result not found in set: %+v", cronL)
		}
	}

	// TEST ERRORS

	// Test that the number of errors returned is correct.
	if len(errList) != len(errResults) {
		t.Errorf(
			"Expected errors: '%d' are less than expected: '%d'\n",
			len(errList),
			len(errResults),
		)
	}

	// Test that the errors are all the anticipated errors.
	for idx, err := range errList {
		if err.Error() != errResults[idx] {
			t.Errorf(
				"Error '%s', was not returned as anticipated: '%s'",
				err.Error(),
				errResults[idx],
			)
		}
	}
}

/*
func TestToIcal(t *testing.T) {
	for _, cron := range cronResults {
		cron.ToIcal()
		break
	}
}
*/

/*
func TestToDates(t *testing.T) {
	for _, cron := range cronResults {
		val, err := cron.ToDates()
		fmt.Println(val, err)
	}
}
*/
