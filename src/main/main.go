package main

import (
	"fmt"
	"models"
	"log"
	"os"
)

func main() {

	fmt.Printf("Hello From Lake Pend Oreille!")

	lakeDatas, err := models.GetLakeDataRecordsFor("2015_01_01")
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

}
