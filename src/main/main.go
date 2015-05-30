package main

import (
	"fmt"
	"models"
	"log"
	"os"
)

func main() {

	fmt.Printf("Hello From Lake Pend Oreille!")

	lakeDatas, err := models.GetDBRecordsFor("2015_01_01")
	if err != nil {
		log.Fatal("Apolo we hear ya. Sit Tight!", err)
		os.Exit(1)
	}

	for _, lakeData := range lakeDatas {
		fmt.Printf("Lake data %s %s %s\n",
			       lakeData.DateRecorded,
			       lakeData.TimeRecorded,
			       lakeData.RecordedValue)
	}

	sampleData := []models.LakeData{
		models.LakeData{"2015_01_01","22:38:52","23.24"},
		models.LakeData{"2015_01_01","22:48:02","23.30"},
		models.LakeData{"2015_01_01","22:48:52","23.31"},
		models.LakeData{"2015_01_01","22:53:02","23.32"},
	}
	models.StoreRecords(sampleData)

}
