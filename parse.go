package main

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const cronlen = 6

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

// Not implemented.
var simpleRangeRE = `^\d{1,2}-\d{1,2}$`
var complexRangeRE = `^\d{1,2}-\d{1,2}$\/\d{1,2}$`

func parseRepeatVal(timeunit string) (int, bool) {
	repvalOrig := strings.SplitAfterN(timeunit, "/", 2)
	repval, _ := strconv.Atoi(repvalOrig[1])
	return repval, true
}

func parseListVal(timeunit string) []int {
	listvalOrig := strings.Split(timeunit, ",")
	var listvals = []int{}
	for _, i := range listvalOrig {
		l, _ := strconv.Atoi(i)
		listvals = append(listvals, l)
	}
	return listvals
}

func regexMatch(timeunit string) int {
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
	// 1. If asterisk on own, then ignore.
	// 2. If integer on its own, then return as value.
	// 3. If */int then we have a repeating value.
	// 4. If we receive int,int etc. we have list of values to parse.
	// 5. If we have int-int we have a range to contend with.
	//
	// I think we can implement 1-4 easily. I don't know if we've captured
	// all cases.

	var t int // Temporary value to hold cron integers.

	var timeval = []int{}
	var repeatFlag = false

	var err error

	unitType := regexMatch(timeunit)
	switch unitType {
	case notimpl:
		return timeval, false, fmt.Errorf("Time pattern is not implemented: %s", timeunit)
	case unused:
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
	// Additional validation per time unit.
	switch typ {
	case minutes:
		//
	case hours:
		//
	case dom:
		//
	case mon:
		//
	case dow:
		//
	}
	return timeval, repeatFlag, nil
}

func validateCron(entry []string) error {
	if len(entry) < cronlen {
		return fmt.Errorf("Cron entry length: '%d' is less than expected: '%d'", len(entry), cronlen)
	}
	return nil
}

func entryToCron(entry []string) (Cron, error) {
	var cron Cron
	err := validateCron(entry)
	if err != nil {
		return cron, err
	}
	cron.Mins, cron.MinsR, err = parseTimeUnit(entry[0], minutes)
	if err != nil {
		return cron, err
	}
	cron.Hrs, cron.HrsR, err = parseTimeUnit(entry[1], hours)
	if err != nil {
		return cron, err
	}
	cron.Dom, cron.DomR, err = parseTimeUnit(entry[2], dom)
	if err != nil {
		return cron, err
	}
	cron.Mon, cron.MonR, err = parseTimeUnit(entry[3], mon)
	if err != nil {
		return cron, err
	}
	cron.Dow, cron.DowR, err = parseTimeUnit(entry[4], dow)
	if err != nil {
		return cron, err
	}
	cron.Command = strings.Join(entry[5:], " ")
	return cron, nil
}

// ParseCronEntry will parse a single cron string and return a single cron
// type to the caller.
func ParseCronEntry(cronEntry string) (Cron, error) {
	entry := strings.Split(cronEntry, " ")
	cron, err := entryToCron(entry)
	if err != nil {
		return cron, err
	}
	return cron, nil
}

func parseCronEntries(entries string) ([]Cron, []error) {
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
		cronlist = append(cronlist, cron)
	}
	return cronlist, errorlist
}

// ParseCrontab will parse an entire crontab file provided as a string
// argument and return a list of cron types to the caller.
func ParseCrontab(cronOutput string) ([]Cron, []error) {
	return parseCronEntries(cronOutput)
}
