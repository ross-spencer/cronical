package main

import (
	"fmt"
	"strings"
	"time"
)

/*Cron describes the structure of a cron line in crontab.

Fields: m h dom mon dow command

 * * * * *  command to execute
 │ │ │ │ │
 │ │ │ │ │
 │ │ │ │ └───── day of week (0 - 6) (0 to 6 are Sunday to Saturday, or use names; 7 is Sunday, the same as 0)
 │ │ │ └────────── month (1 - 12)
 │ │ └─────────────── day of month (1 - 31)
 │ └──────────────────── hour (0 - 23)
 └───────────────────────── min (0 - 59)

Fields can also be comma separated. They can also have ranges, but ranges are
not currently implemented.
*/
type Cron struct {
	Mins    []int  // Number of minutes on which to run the command.
	MinsR   bool   // Are the minutes repeated?
	Hrs     []int  // Hour to run the command.
	HrsR    bool   // Is the number of hours repeated?
	Dom     []int  // Day of the month to run the command?
	DomR    bool   // Is the day of the month repeated?
	Mon     []int  // Month to run the command.
	MonR    bool   // Is the month repeated?
	Dow     []int  // Day of the week to run the command.
	DowR    bool   // Are the days of the week repeated?
	Command string // Command for cron to run.
}

var limit = 10

func (c *Cron) null() bool {
	if len(c.Mon) == 0 && len(c.Dom) == 0 && len(c.Hrs) == 0 && len(c.Mins) == 0 && len(c.Dow) == 0 {
		return true
	}
	return false
}

func (c *Cron) anyR() bool {
	// If any repeat field is true, return true, else false.
	//
	if !c.MinsR && !c.HrsR && !c.DomR && !c.MonR && !c.DowR {
		return false
	}
	return true
}

func (c *Cron) setDefaults() {
	if len(c.Mon) == 0 {
		c.Mon = append(c.Mon, 0)
	}
	if len(c.Dom) == 0 {
		c.Dom = append(c.Dom, 0)
	}
	if len(c.Hrs) == 0 {
		c.Hrs = append(c.Hrs, 0)
	}
	if len(c.Mins) == 0 {
		c.Mins = append(c.Mins, 0)
	}
	if len(c.Dow) == 0 {
		c.Dow = append(c.Dow, 0)
	}
}

func (c *Cron) toRepeatingDates() string {
	// Output dates based on repeating fields.
	return "Return repeating entries..."
}

func (c *Cron) toSpecificDates() string {
	// Process little to big-endian adding data as we go...
	//
	// Fields: m h dom mon dow command
	//
	c.setDefaults()
	date := time.Now()
	for i := 0; i < limit; i++ {
		// Create X entries...
	}
	return fmt.Sprintf("Return specific dates: %s", date)
}

// ToDates ...
func (c *Cron) ToDates() string {
	if c.null() {
		return "Do nothing, we have a nil entry..."
	}
	if c.anyR() {
		// Fields repeat, take special action...
		return c.toRepeatingDates()
	}
	// Dates are specific.
	return c.toSpecificDates()
}

// ToIcal ...
func (c *Cron) ToIcal() {
	// Output to ICAL.
}

// ToCron ...
func (c *Cron) ToCron() {
	// Output to cron again.
}

func joinInt(is []int, delim string) string {
	// One-liner: https://stackoverflow.com/a/37533144
	//
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(is)), delim), "[]")
}

func (c Cron) String() string {
	return fmt.Sprintf(
		"Minutes: %s (%t)\nHours: %s (%t)\nDay of Month: %s (%t)\nMonth: %s (%t)\nDay of Week: %s (%t)\nCommand: %s",
		joinInt(c.Mins, ","), c.MinsR,
		joinInt(c.Hrs, ","), c.HrsR,
		joinInt(c.Dom, ","), c.DomR,
		joinInt(c.Mon, ","), c.MonR,
		joinInt(c.Dow, ","), c.DowR,
		c.Command,
	)
}
