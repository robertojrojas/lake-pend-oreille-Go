package statistics

import (
	"models"
	"fmt"
	"net/http"
	"encoding/json"
	"strings"
)

const (
	MEAN   = "mean"
	MEDIAN = "median"
	DATE_FIELD = "date"
)

var INVALID_URL = map[string]string{
	"Error": "Invalid URL",
}

type ReportData struct {
	Mean float64
	Median float64
}

type ReportOutput struct {
	Date string
	ReportData map[string]ReportData
}

func GenerateReportDisplay(date string) {

	reportOutput, err := GenerateReport(date)

	if err != nil {
		panic(err)
	}
	reportValues := reportOutput.ReportData
	fmt.Printf("================================================================ \n")
	fmt.Printf("            STATISTICAL ANALYSIS For %s  \n", reportOutput.Date     )
	fmt.Printf("================================================================ \n")
	fmt.Printf("        Air Temperature    Barometric Pressuare   Wind Speed     \n")
    fmt.Printf("  MEAN     %.2f                %.2f                %.2f          \n",
		              reportValues[models.DATASOURCE_TYPES[0]].Mean,
		              reportValues[models.DATASOURCE_TYPES[1]].Mean,
		              reportValues[models.DATASOURCE_TYPES[2]].Mean                   )
	fmt.Printf(" MEDIAN    %.2f                %.2f                %.2f           \n",
					  reportValues[models.DATASOURCE_TYPES[0]].Median,
					  reportValues[models.DATASOURCE_TYPES[1]].Median,
					  reportValues[models.DATASOURCE_TYPES[2]].Median   )
	fmt.Printf("================================================================ \n")


}

func GenerateReport(date string) (ReportOutput, error) {

	reportValues := map[string]ReportData{}

	for _, dataSourceType := range models.DATASOURCE_TYPES {

		var lakeDatas models.DataRecs

		//fmt.Printf("Checking records for %s %s \n", date, dataSourceType)
		recordsExist, err := models.CheckDBRecordsFor(date, dataSourceType)

		if err != nil {
			fmt.Printf("CheckDBRecordsFor - Problems checking Records for %s %s \n", date, dataSourceType)
			return ReportOutput{}, err
		}

		if !recordsExist {
			models.FetchData(date, dataSourceType)
		}

		lakeDatas, err = models.GetDBRecordsFor(date, dataSourceType)
		if err != nil {
			fmt.Printf("GetDBRecordsFor - Problems checking Records for %s %s \n", date, dataSourceType)
			return ReportOutput{}, err
		}

		meanValue := lakeDatas.Mean()
		medianValue := lakeDatas.Median()

		reportValues[dataSourceType] = ReportData{
			Mean   : meanValue,
			Median : medianValue,
		}
	}

	reportOutput := ReportOutput{
		Date: date,
		ReportData:reportValues,
	}

	return reportOutput, nil

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

func generateJson(w http.ResponseWriter, data interface{}) {
	jsonEnc := json.NewEncoder(w)
	jsonEnc.Encode(data)
}



func CreateServer() {
	http.HandleFunc("/reports/", reportHandler)
	http.ListenAndServe(":8888", nil)
}


