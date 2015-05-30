package statistics

import (
	"models"
	"fmt"
)

func GenerateReportDisplay(date string) {



	reportValues := map[string][]float64{}

	for _, dataSourceType := range models.DATASOURCE_TYPES {

		var lakeDatas models.DataRecs

		//fmt.Printf("Checking records for %s %s \n", date, dataSourceType)
		recordsExist, err := models.CheckDBRecordsFor(date, dataSourceType)

		if err != nil {
			fmt.Printf("CheckDBRecordsFor - Problems checking Records for %s %s \n", date, dataSourceType)
			continue
		}

		if !recordsExist {
			models.FetchData(date, dataSourceType)
		}

		lakeDatas, err = models.GetDBRecordsFor(date, dataSourceType)
		if err != nil {
			fmt.Printf("GetDBRecordsFor - Problems checking Records for %s %s \n", date, dataSourceType)
			continue
		}

		meanValue := lakeDatas.Mean()
		//fmt.Printf("The MEAN for  %s %s is %f\n", dataSourceType, date, meanValue)

		medianValue := lakeDatas.Median()
		//fmt.Printf("The MEDIAN for %s %s is %f\n", dataSourceType, date, medianValue)

		reportValues[dataSourceType] = []float64{
			meanValue,
			medianValue,
		}
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
