package main

import (
	"fmt"
	"time"

	"github.com/arran4/golang-ical"
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
	Mins    int    // Number of minutes on which to run the command.
	MinsR   bool   // Are the minutes repeated?
	Hrs     int    // Hour to run the command.
	HrsR    bool   // Is the number of hours repeated?
	Dom     int    // Day of the month to run the command?
	DomR    bool   // Is the day of the month repeated?
	Mon     int    // Month to run the command.
	MonR    bool   // Is the month repeated?
	Dow     int    // Day of the week to run the command.
	DowR    bool   // Are the days of the week repeated?
	Command string // Command for cron to run.
}

const timeFormat = "2006-01-02T15:04"

var limit = 3

func (c *Cron) null() bool {
	if c.Mon == -1 &&
		c.Dom == -1 &&
		c.Hrs == -1 &&
		c.Mins == -1 &&
		c.Dow == -1 {
		return true
	}
	return false
}

func (c *Cron) anyR() bool {
	// If any repeat field is true, return true, else false.
	if !c.MinsR && !c.HrsR && !c.DomR && !c.MonR && !c.DowR {
		return false
	}
	return true
}

func (c *Cron) toRepeatingDates() string {
	// Output dates based on repeating fields.
	return "Return repeating entries..."
}

func (c *Cron) toSpecificDates() (string, error) {
	// Process little to big-endian adding data as we go...
	//
	// Fields: m h dom mon dow command
	//
	// Work backwards from command, and calculate increments (big-endian)
	// processing.

	t := getAndResetTime()
	ts := []time.Time{}

	for i := 0; i < limit; i++ {
		ts = append(ts, t)
	}

	if c.Dow != -1 {
		t = setDow(c.Dow, t)
	}

	if c.Mon != -1 {
		t = setMon(c.Mon, t)
	}

	if c.Dom != -1 {
		t = setDom(c.Dom, t)
		if c.Mon != -1 && t.Month() != time.Month(c.Mon) {
			return "",
				fmt.Errorf(
					"Cron entry is invalid, month: '%d', day: '%d'",
					c.Mon,
					c.Dom,
				)
		}
	}

	if c.Hrs != -1 {
		t = setHours(c.Hrs, t)
	}

	if c.Mins != -1 {
		t = setMins(c.Mins, t)
	}

	fmt.Println("---")

	fmt.Printf(
		"Specific date: '%s' cmd: '%s' \n",
		t.Format(timeFormat),
		c.Command,
	)

	return fmt.Sprintf("---"), nil
}

// ToIcal will convert cron entries to ical formatted events.
func (c *Cron) ToIcal() {
	cal := ics.NewCalendar()
	event := cal.AddEvent("cronical-cron-entry")
	event.SetCreatedTime(time.Now())
	event.SetDtStampTime(time.Now())
	event.SetModifiedAt(time.Now())
	event.SetStartAt(time.Now())
	event.SetEndAt(time.Now())
	event.SetSummary("cron")
	event.SetLocation("ASRV-01")
	event.SetDescription("Execute command...")
	event.SetOrganizer("sender@domain", ics.WithCN("This Machine"))
	fmt.Println(event.Serialize())
}

// ToDates converts cron structures into the next possible date entries in a
// standard calendar.
func (c *Cron) ToDates() (string, error) {
	if c.null() {
		return "Do nothing, we have a nil entry...", nil
	}
	if c.anyR() {
		// Fields repeat, we need to handle this differently.
		return c.toRepeatingDates(), nil
	}
	return c.toSpecificDates()
}

// ToCron ...
func (c *Cron) ToCron() {
	// Output Cron{} to cron.
}

func (c Cron) String() string {
	return fmt.Sprintf(
		"Cron entry (Command: '%s'): %d (%t), %d (%t), %d (%t), %d (%t), %d (%t)",
		c.Command,
		c.Mins, c.MinsR,
		c.Hrs, c.HrsR,
		c.Dom, c.DomR,
		c.Mon, c.MonR,
		c.Dow, c.DowR,
	)
}
