package main

import (
	"time"
)

func getTime() time.Time {
	// Return time set back to midnight.
	t := time.Now()
	t = t.Add(-time.Minute * time.Duration(t.Minute()))
	t = t.Add(-time.Hour * time.Duration(t.Hour()))
	return t
}

func setHours(hour int, t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), hour, t.Minute(), 00, 00, time.UTC)
}

func setMins(minute int, t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), minute, 00, 00, time.UTC)
}

func setDow(day int, t time.Time) time.Time {
	for true {
		t = t.Add(time.Hour * 24)
		if t.Weekday() == time.Weekday(day) {
			break
		}
	}
	return t
}

func setMon(mon int, t time.Time) time.Time {
	return time.Date(t.Year(), time.Month(mon), t.Day(), t.Hour(), t.Minute(), 00, 00, time.UTC)
}

func setDom(dom int, t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), dom, t.Hour(), t.Minute(), 00, 00, time.UTC)
}
