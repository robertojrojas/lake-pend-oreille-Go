package main

import (
	"os"
	"statistics"
	"runtime"
)

func main() {

	runtime.GOMAXPROCS(3)

	var dateParam string

	if len(os.Args) > 1 {
		dateParam = os.Args[1]
		statistics.GenerateReportDisplay(dateParam)
	} else {
		dateParam = "2015-01-01"
		statistics.CreateServer()
	}


}
