package statistics

import (
	"models"
	"fmt"
	"net/http"
	"encoding/json"
)

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
		              reportValues[models.DATASOURCE_TYPES[0]][0],
		              reportValues[models.DATASOURCE_TYPES[1]][0],
		              reportValues[models.DATASOURCE_TYPES[2]][0]                   )
	fmt.Printf(" MEDIAN    %.2f                %.2f                %.2f           \n",
					  reportValues[models.DATASOURCE_TYPES[0]][1],
					  reportValues[models.DATASOURCE_TYPES[1]][1],
					  reportValues[models.DATASOURCE_TYPES[2]][1]   )
	fmt.Printf("================================================================ \n")


}

func GenerateReport(date string) (map[string][]float64, error) {

	reportValues := map[string][]float64{}

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

		reportValues[dataSourceType] = []float64{
			meanValue,
			medianValue,
		}
	}

	return reportValues, nil

}

func reportHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	date := r.Form["date"][0]

	p, err := GenerateReport(date)
	if err != nil {
		http.Error(w, "File not found", http.StatusInternalServerError)
	} else {
		jsonEnc := json.NewEncoder(w)
		jsonEnc.Encode(p)
	}
}


func CreateServer() {
	http.HandleFunc("/reports", reportHandler)
	http.ListenAndServe(":8888", nil)
}


