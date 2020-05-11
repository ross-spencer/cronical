package main

import (
	"fmt"
)

var app = "cronical"
var ver = "0.0.1"

func version() string {
	return fmt.Sprintf("%s-%s", app, ver)
}

var icalHdr = `BEGIN:VCALENDAR
VERSION:2.0
PRODID:cronical-0.0.1`

var icalFtr = "END:VCALENDAR"

func icalHeader() string {
	return icalHdr
}

func icalFooter() string {
	return icalFtr
}
