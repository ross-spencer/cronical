package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

const fixtures = "crontests"
const crontabFile = "cron-entries"

var crontabFixture string

// Errors anticipated from parsing the cron fixtures.
var errResults = [...]string{
	"Time pattern is not implemented: π",

	// TODO: These are only temporarily not implemented. To fix we need to:
	//
	// 1. Create regex patterns to match MON, DAY.
	// 2. Map the values back into integers for easy of use in code.
	//
	// NB. That would make conversion from cron -> cronical non-reversible, but
	// that might be okay, e.g. MON -> cronical ;; cronical -> cron would
	// return 1.
	//
	"Time pattern is not implemented: JAN",
	"Time pattern is not implemented: SUN",
}

// cron results anticipated from parsing the cron fixtures.
//
// Fields: m h dom mon dow command
//
var cronResults = [...]Cron{
	Cron{[]int{10, 15}, false, []int{16}, false, []int{}, false, []int{}, false, []int{}, false, "echo \"test one\""},
	Cron{[]int{0}, false, []int{2}, false, []int{}, false, []int{}, false, []int{5, 6}, false, "echo \"test two\""},
	Cron{[]int{0}, false, []int{0}, false, []int{1}, false, []int{1}, false, []int{}, false, "echo \"happy new year!\""},
	Cron{[]int{0}, false, []int{}, false, []int{}, false, []int{}, false, []int{}, false, "echo \"test three\""},
	Cron{[]int{15}, true, []int{}, false, []int{}, false, []int{}, false, []int{}, false, "echo \"test four\""},
	Cron{[]int{30}, true, []int{}, false, []int{}, false, []int{}, false, []int{}, false, "echo \"do not echo\" 1&2> /dev/null"},
	Cron{[]int{10}, true, []int{}, false, []int{}, false, []int{}, false, []int{}, false, "echo \"test five\""},
	Cron{[]int{}, false, []int{0, 2, 4}, false, []int{}, false, []int{}, false, []int{}, false, "echo \"test six\""},
	Cron{[]int{}, false, []int{1}, true, []int{}, false, []int{}, false, []int{}, false, "echo \"test seven\""},
	Cron{[]int{}, false, []int{}, false, []int{31}, false, []int{2, 8}, false, []int{}, false, "echo \"test eight\""},
	Cron{[]int{10}, false, []int{}, false, []int{}, false, []int{10}, false, []int{}, false, "echo \"test nine\""},
	Cron{[]int{}, false, []int{}, false, []int{}, false, []int{}, false, []int{}, false, "echo \"test ten\""},
}

func init() {
	crontabFixture = path.Join(fixtures, crontabFile)
}

func TestParse(t *testing.T) {
	testCrontab, err := os.Open(crontabFixture)
	if err != nil {
		t.Errorf(fmt.Sprintf("%s\n", err))
	}
	defer testCrontab.Close()

	crontab, err := ioutil.ReadAll(testCrontab)
	if err != nil {
		t.Errorf(fmt.Sprintf("%s\n", err))
	}

	cronlist, errlist := ParseCrontab(string(crontab))
	if len(errlist) != len(errResults) {
		t.Errorf(
			"Expected errors: '%d' are less than expected: '%d'\n",
			len(errlist),
			len(errResults),
		)
	}

	for idx, err := range errlist {
		if err.Error() != errResults[idx] {
			t.Errorf(
				"Error '%s', was not returned as anticipated: '%s'",
				err.Error(),
				errResults[idx],
			)
		}
	}

	if len(cronlist) != len(cronResults) {
		t.Errorf(
			"Entries to parse: '%d' is different than expected: '%d'",
			len(cronlist),
			len(cronResults),
		)
	}

	for idx, cron := range cronlist {
		if !reflect.DeepEqual(cron, cronResults[idx]) {
			t.Errorf("Cron results aren't equivalent:\n%+v,\n%+v", cron, cronResults[idx])
		}
	}

}

func TestToIcal(t *testing.T) {
	for _, cron := range cronResults {
		cron.ToIcal()
		break
	}
}

func TestToDates(t *testing.T) {
	for _, cron := range cronResults {
		fmt.Println(cron.ToDates())
	}
}
