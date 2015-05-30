package statistics

import (
	"models"
	"fmt"
)

func GenerateReportDisplay(date string) {

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
		fmt.Printf("The MEAN for  %s %s is %f\n", dataSourceType, date, meanValue)

		medianValue := lakeDatas.Median()
		fmt.Printf("The MEDIAN for %s %s is %f\n", dataSourceType, date, medianValue)

	}


}
