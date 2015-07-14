package statistics

import (
	"models"
	"fmt"
	"net/http"
	"encoding/json"
	"strings"
	"text/template"
	"os"
	"io"
	"sync"
)

const (
	MEAN   = "mean"
	MEDIAN = "median"
	DATE_FIELD = "date"
)


const REPORT_TEMPLATE = `
================================================================
	STATISTICAL ANALYSIS For {{.ReportDate}}
================================================================
	  Air Temperature   Barometric Pressuare    Wind Speed
  MEAN		 {{.AirTemp.Mean}}		{{.Barometric.Mean}}		     {{.WinSpeed.Mean}}
 MEDIAN		 {{.AirTemp.Median}}		{{.Barometric.Median}}		     {{.WinSpeed.Median}}
================================================================
	`

var INVALID_URL = map[string]string{
	"Error": "Invalid URL",
}

type ReportData struct {
	Mean   float64
	Median float64
}


type ReportOutput struct {
	ReportDate  string
	AirTemp     ReportData
	Barometric  ReportData
	WinSpeed    ReportData
}

func GenerateReportDisplay(date string) {

	reportOutput, err := GenerateReport(date)

	if err != nil {
		panic(err)
	}

	t := template.Must(template.New("letter").Parse(REPORT_TEMPLATE))

	err = t.Execute(os.Stdout, reportOutput)
	if err != nil {
		panic(err)
	}


}

func GenerateReport(date string) (ReportOutput, error) {

	reportOutput := ReportOutput{
		ReportDate: date,
	}

	// Create an unbuffered channel to receive match results to display.
	results := make(chan *models.DataRecs)

	// Setup a wait group so we can process all the statistics data.
	var waitGroup sync.WaitGroup

	// Set the number of goroutines we need to wait for while
	// they process the individual statistic data.
	waitGroup.Add(len(models.DATASOURCE_TYPES))

	for _, dataSourceType := range models.DATASOURCE_TYPES {

		// Launch the goroutine to perform the statistics data retrieval.
		go func (theDate string, theDataSourceType string) {
			models.GetStatisticDataFor(theDate, theDataSourceType, results)
			fmt.Printf("Done getting data for %s\n", theDataSourceType)
			waitGroup.Done()
		}(date, dataSourceType)


	}

	// Launch a goroutine to monitor when all the work is done.
	go func() {
		// Wait for everything to be processed.
		waitGroup.Wait()

		// Close the channel to signal to the Display
		// function that we can exit the program.
		close(results)
	}()

	collectReportData(&reportOutput, results)

	return reportOutput, nil

}

func collectReportData(reportOutput *ReportOutput, results chan *models.DataRecs) {
	// Now collect all the data
	for lakeDataRecord := range results {

		meanValue   := lakeDataRecord.Mean()
		medianValue := lakeDataRecord.Median()

		switch lakeDataRecord.GetDataSource() {
		case models.AirTemp:
			reportOutput.AirTemp    = ReportData{Mean:meanValue, Median:medianValue}
		case models.BarometricPress:
			reportOutput.Barometric = ReportData{Mean:meanValue, Median:medianValue}
		case models.Wind_Speed:
			reportOutput.WinSpeed   = ReportData{Mean:meanValue, Median:medianValue}
		}
	}

}


func reportHandler(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path
	pathParts := strings.Split(path, "/")

	if len(pathParts) != 3 {
		generateJson(w, INVALID_URL)
		return
	}

	date := pathParts[2]

	p, err := GenerateReport(date)
	if err != nil {
		fmt.Printf("Error parsing the %s\n", err)
		http.Error(w, "File not found", http.StatusInternalServerError)
	} else {
		generateJson(w, p)
	}
}

/*
   ##########################
         Sample JSON
   ##########################
   {
      "ReportDate":"2014-01-01",
      "AirTemp":{
         "Mean":36.18,
         "Median":36.3
      },
      "Barometric":{
          "Mean":28.26,
          "Median":28.3
      },
      "WinSpeed":{
          "Mean":3.89,
          "Median":2.8
      }
    }
   ##########################

*/
func generateJson(w io.Writer, data interface{}) {
	jsonEnc := json.NewEncoder(w)
	jsonEnc.Encode(data)
}


func CreateServer() {
	http.HandleFunc("/reports/", reportHandler)
	http.ListenAndServe(":8888", nil)
}


