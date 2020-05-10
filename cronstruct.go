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

const timeFormat = "2006-01-02T15:04"

var limit = 2

func (c *Cron) null() bool {
	if len(c.Mon) == 0 && len(c.Dom) == 0 && len(c.Hrs) == 0 && len(c.Mins) == 0 && len(c.Dow) == 0 {
		return true
	}
	return false
}

func (c *Cron) getCronLen() int {
	// Return the number of possible entries we can create from one single
	// line of cron.
	l := 1
	if len(c.Mon) > 1 {
		l += len(c.Mon) - 1
	}
	if len(c.Dom) > 1 {
		l += len(c.Dom) - 1
	}
	if len(c.Hrs) > 1 {
		l += len(c.Hrs) - 1
	}
	if len(c.Mins) > 1 {
		l += len(c.Mins) - 1
	}
	if len(c.Dow) > 1 {
		l += len(c.Dow) - 1
	}
	return l
}

func (c *Cron) anyR() bool {
	// If any repeat field is true, return true, else false.
	//
	if !c.MinsR && !c.HrsR && !c.DomR && !c.MonR && !c.DowR {
		return false
	}
	return true
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
	// Work backwards from command, and calculate increments (big-endian)
	// processing.

	t := getTime()
	ts := []time.Time{}

	for i := 0; i < limit*c.getCronLen(); i++ {
		ts = append(ts, t)
	}

	// Hours and minutes are going to be consistent across entries. Set those
	// here.
	if len(c.Hrs) > 0 {
		idx := 0
		for i := 0; i < len(ts); i += len(c.Hrs) {
			for _, val := range c.Hrs {
				ts[idx] = setHours(val, ts[idx])
				idx++
			}
		}
		// TODO: Add a day so return is on this hour, on said say.
	}
	if len(c.Mins) > 0 {
		idx := 0
		for i := 0; i < len(ts); i += len(c.Mins) {
			for _, val := range c.Mins {
				ts[idx] = setMins(val, ts[idx])
				idx++
			}
		}
		// TODO: Add an hour so return is on the hour.
	}

	// The interplay between the big-endian values is more difficult. For
	// example, setting a DoW and DoM is like setting a conditional.
	//
	// E.g. Run on 1 January, if 1 January is also a Monday.
	//

	var exceptionDow = false

	if len(c.Dow) > 0 {
		idx := 0
		for i := 0; i < len(ts); i += len(c.Dow) {
			for _, val := range c.Dow {
				ts[idx] = setDow(val, ts[idx])
				// Set the next value relative to this one, so we're not going
				// back in time
				if idx+1 < len(ts) {
					ts[idx+1] = ts[idx]
				}
				idx++
			}
		}
		exceptionDow = true
	}

	if len(c.Mon) > 0 {
		for _, _ = range c.Mon {
			for _, _ = range ts {
				// Set month.
			}
		}
	}

	if len(c.Dom) > 0 {
		for _, _ = range c.Dom {
			for _, _ = range ts {
				// Set Day of month.
				//
				if exceptionDow == true {
					// if exceptionDOW then day of month must also be equal,
					// else add one year until we find that day.
				}
			}
		}
	}

	for _, v := range ts {
		fmt.Printf("Generated specific date: %s %s\n", v.Format(timeFormat), c.Command)
	}

	return fmt.Sprintf("---")
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
