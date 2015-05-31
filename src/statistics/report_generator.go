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
)

var INVALID_URL = map[string]string{
	"Error": "Invalid URL",
}

func GenerateReportDisplay(date string) {

	reportValues, err := GenerateReport(date)

	if err != nil {
		panic(err)
	}

	fmt.Printf("================================================================ \n")
	fmt.Printf("            STATISTICAL ANALYSIS For %s  \n", date)
	fmt.Printf("================================================================ \n")
	fmt.Printf("        Air Temperature    Barometric Pressuare   Wind Speed     \n")
    fmt.Printf("  MEAN     %.2f                %.2f                %.2f          \n",
		              reportValues[models.DATASOURCE_TYPES[0]][MEAN],
		              reportValues[models.DATASOURCE_TYPES[1]][MEAN],
		              reportValues[models.DATASOURCE_TYPES[2]][MEAN]                   )
	fmt.Printf(" MEDIAN    %.2f                %.2f                %.2f           \n",
					  reportValues[models.DATASOURCE_TYPES[0]][MEDIAN],
					  reportValues[models.DATASOURCE_TYPES[1]][MEDIAN],
					  reportValues[models.DATASOURCE_TYPES[2]][MEDIAN]   )
	fmt.Printf("================================================================ \n")


}

func GenerateReport(date string) (map[string]map[string]float64, error) {

	reportValues := map[string]map[string]float64{}

	for _, dataSourceType := range models.DATASOURCE_TYPES {

		var lakeDatas models.DataRecs

		//fmt.Printf("Checking records for %s %s \n", date, dataSourceType)
		recordsExist, err := models.CheckDBRecordsFor(date, dataSourceType)

		if err != nil {
			fmt.Printf("CheckDBRecordsFor - Problems checking Records for %s %s \n", date, dataSourceType)
			return nil, err
		}

		if !recordsExist {
			models.FetchData(date, dataSourceType)
		}

		lakeDatas, err = models.GetDBRecordsFor(date, dataSourceType)
		if err != nil {
			fmt.Printf("GetDBRecordsFor - Problems checking Records for %s %s \n", date, dataSourceType)
			return nil, err
		}

		meanValue := lakeDatas.Mean()
		medianValue := lakeDatas.Median()

		reportValues[dataSourceType] = map[string]float64{
			MEAN   : meanValue,
			MEDIAN : medianValue,
		}
	}

	return reportValues, nil

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


