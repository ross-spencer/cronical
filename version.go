package main

import (
	"fmt"
)

var app = "cronical"
var ver = "0.0.1"

func version() string {
	return fmt.Sprintf("%s-%s", app, ver)
}
