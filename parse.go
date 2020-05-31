package main

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// cronFields describes the number of fields expected in a single line of cron.
const cronFields = 6

const (
	minutes = iota
	hours
	dom
	mon
	dow

	unused
	repeat
	single
	list
	notimpl
)

var unusedEntryRE = `^\*$`
var repeatEntryRE = `^\*\/\d{1,2}$`
var singleEntryRE = `^\d{1,2}$`
var listEntryRE = `^\d{1,2},`

// Not implemented, e.g. 1-2, 9-12
var simpleRangeRE = `^\d{1,2}-\d{1,2}$`

// I don't know if complexRange is likely/possible...
var complexRangeRE = `^\d{1,2}-\d{1,2}\/\d{1,2}$`

func parseRepeatVal(timeunit string) (int, bool) {
	// parseRepeatVal returns the integer value associated with a repeating
	// (*/n) style cron entry.
	repvalOrig := strings.SplitAfterN(timeunit, "/", 2)
	repval, _ := strconv.Atoi(repvalOrig[1])
	return repval, true
}

func parseListVal(timeunit string) []int {
	// parseList value returns the complete list of integers associated with
	// a comma-separated style cron entry.
	listvalOrig := strings.Split(timeunit, ",")
	var listvals = []int{}
	for _, i := range listvalOrig {
		l, _ := strconv.Atoi(i)
		listvals = append(listvals, l)
	}
	return listvals
}

func regexMatch(timeunit string) int {
	// regexMatch allows us to identify what type of cron entry we're working
	// with.
	for _, value := range [...]string{simpleRangeRE, complexRangeRE} {
		m, _ := regexp.Match(value, []byte(timeunit))
		if m == true {
			return notimpl
		}
	}
	m, _ := regexp.Match(unusedEntryRE, []byte(timeunit))
	if m == true {
		return unused
	}
	m, _ = regexp.Match(singleEntryRE, []byte(timeunit))
	if m == true {
		return single
	}
	m, _ = regexp.Match(repeatEntryRE, []byte(timeunit))
	if m == true {
		return repeat
	}
	m, _ = regexp.Match(listEntryRE, []byte(timeunit))
	if m == true {
		return list
	}
	return notimpl
}

func parseTimeUnit(timeunit string, typ int) ([]int, bool, error) {
	// Value might look like as follows:
	//
	// IMPLEMENTED
	//
	// 1. If asterisk on own, then ignore.
	// 2. If integer on its own, then return as value.
	// 3. If */int then we have a repeating value.
	// 4. If we receive int,int etc. we have list of values to parse.
	//
	// TODO: NOT IMPLEMENTED
	//
	// 5. If we have int-int we have a range to contend with
	//
	var t int // Temporary value to hold cron integers.
	var timeval = []int{}
	var repeatFlag = false
	var err error
	unitType := regexMatch(timeunit)
	switch unitType {
	case notimpl:
		return timeval, false, fmt.Errorf("Time pattern is not implemented: %s", timeunit)
	case unused:
		// Because '0' can be used for Sunday as well as '7' we return -1 to
		// say that this field is unused. We do this for all fields that are
		// unused.
		timeval = append(timeval, -1)
		return timeval, false, nil
	case single:
		t, err = strconv.Atoi(timeunit)
		if err != nil {
			return timeval, repeatFlag, err
		}
		timeval = append(timeval, t)
	case repeat:
		t, repeatFlag = parseRepeatVal(timeunit)
		timeval = append(timeval, t)
	case list:
		timeval = parseListVal(timeunit)
		return timeval, repeatFlag, nil
	}
	return timeval, repeatFlag, nil
}

func validateCron(entry []string) error {
	// validateCron will perform rudimentary validation on the cron entry being
	// parsed. Not all validation is performed in the module. Some validation
	// will be performed when the cron is converted to dates.
	if len(entry) < cronFields {
		return fmt.Errorf(
			"Cron entry length: '%d' is less than expected: '%d'",
			len(entry),
			cronFields,
		)
	}
	if entry[0] == "*" &&
		entry[1] == "*" &&
		entry[2] == "*" &&
		entry[3] == "*" &&
		entry[4] == "*" {
		return fmt.Errorf("Cron entry: '%s' is null", entry[0:5])
	}
	return nil
}

func createEntries(
	cron []Cron,
	mins []int,
	hrs []int,
	dom []int,
	mon []int,
	dow []int,
) []Cron {
	// createEntries returns a cron{} per permutation of entry supplied.
	//
	// All of the permutations of a cron entry are laid out in the various
	// arrays given to us.
	//
	// Each slice implies a new permutation and so the possibilities are
	// multiplicative. two slices, with length two requires 2 x 2 new cron
	// entries to be generated, a third slice of 3 == 2 x 2 x 3.
	//
	// Process:
	//
	// If the slice is 1 in length, add the value to all cron entries.
	//
	// If the slice is > 1, copy all cron entries so far. For all existing
	// entries then add one of the slice values to each. Clone the cron a
	// second time and then add all the next slice values entries to that too.
	// Repeat clone + add per slice value.
	var field = mins
	var newCron = cron
	if len(field) > 1 {
		var clone []Cron
		for _, val := range field {
			for idx := range newCron {
				newCron[idx].Mins = val
			}
			clone = append(clone, newCron...)
		}
		cron = clone
	} else if len(field) == 1 {
		for idx := 0; idx < len(cron); idx++ {
			cron[idx].Mins = field[0]
		}
	}

	// hrs

	field = hrs
	newCron = cron
	if len(field) > 1 {
		var clone []Cron
		for _, val := range field {
			for idx := range newCron {
				newCron[idx].Hrs = val
			}
			clone = append(clone, newCron...)
		}
		cron = clone
	} else if len(field) == 1 {
		for idx := 0; idx < len(cron); idx++ {
			cron[idx].Hrs = field[0]
		}
	}

	// dom

	field = dom
	newCron = cron
	if len(field) > 1 {
		var clone []Cron
		for _, val := range field {
			for idx := range newCron {
				newCron[idx].Dom = val
			}
			clone = append(clone, newCron...)
		}
		cron = clone
	} else if len(field) == 1 {
		for idx := 0; idx < len(cron); idx++ {
			cron[idx].Dom = field[0]
		}
	}

	// mon

	field = mon
	newCron = cron
	if len(field) > 1 {
		var clone []Cron
		for _, val := range field {
			for idx := range newCron {
				newCron[idx].Mon = val
			}
			clone = append(clone, newCron...)
		}
		cron = clone
	} else if len(field) == 1 {
		for idx := 0; idx < len(cron); idx++ {
			cron[idx].Mon = field[0]
		}
	}

	// dow

	field = dow
	newCron = cron
	if len(field) > 1 {
		var clone []Cron
		for _, val := range field {
			for idx := range newCron {
				newCron[idx].Dow = val
			}
			clone = append(clone, newCron...)
		}
		cron = clone
	} else if len(field) == 1 {
		for idx := 0; idx < len(cron); idx++ {
			cron[idx].Dow = field[0]
		}
	}

	return cron
}

func entryToCron(entry []string) ([]Cron, error) {
	// entryToCron takes an entry and splits it into all of its different
	// permutations and returns that as a slice to the caller.

	var cron = []Cron{}
	cron = append(cron, Cron{})
	cron[0].Command = strings.Join(entry[5:], " ")

	err := validateCron(entry)
	if err != nil {
		return cron, err
	}

	mins, minsR, err := parseTimeUnit(entry[0], minutes)
	if err != nil {
		return cron, err
	}
	cron[0].MinsR = minsR

	hrs, hrsR, err := parseTimeUnit(entry[1], hours)
	if err != nil {
		return cron, err
	}
	cron[0].HrsR = hrsR

	dom, domR, err := parseTimeUnit(entry[2], dom)
	if err != nil {
		return cron, err
	}
	cron[0].DomR = domR

	mon, monR, err := parseTimeUnit(entry[3], mon)
	if err != nil {
		return cron, err
	}
	cron[0].MonR = monR

	dow, dowR, err := parseTimeUnit(entry[4], dow)
	if err != nil {
		return cron, err
	}
	cron[0].DowR = dowR

	cron = createEntries(cron, mins, hrs, dom, mon, dow)
	return cron, nil
}

// ParseCronEntry will parse a single cron string and return a single cron
// type to the caller.
func ParseCronEntry(cronEntry string) ([]Cron, error) {
	entry := strings.Split(cronEntry, " ")
	cron, err := entryToCron(entry)
	if err != nil {
		return cron, err
	}
	return cron, nil
}

func parseCronEntries(entries string) ([]Cron, []error) {
	// parseCronEntries will parse all cron entries in a supplied text extract.
	// e.g. an entire crontab file.
	var cronlist = []Cron{}
	var errorlist = []error{}
	scanner := bufio.NewScanner(strings.NewReader(entries))
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") || scanner.Text() == "" {
			continue
		}
		cron, err := ParseCronEntry(scanner.Text())
		if err != nil {
			errorlist = append(errorlist, err)
			continue
		}
		cronlist = append(cronlist, cron...)
	}
	return cronlist, errorlist
}

// ParseCrontab will parse an entire crontab file provided as a string
// argument and return a list of cron types to the caller.
func ParseCrontab(cronOutput string) ([]Cron, []error) {
	return parseCronEntries(cronOutput)
}
