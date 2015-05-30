package main

import (
	"os"
	"statistics"
)

func main() {

	var dateParam string

	if len(os.Args) > 1 {
		dateParam = os.Args[1]
	} else {
		dateParam = "2015-01-01"
	}
	statistics.GenerateReportDisplay(dateParam)

}
