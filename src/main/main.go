package main

import (
	"fmt"
//	_ "models"
//	_ "log"
//	_ "os"
	"web"
)

func main() {

	//fmt.Printf("Hello From Lake Pend Oreille!\n")

//	lakeDatas, err := models.GetDBRecordsFor("airtemp", "2015-01-01")
//	if err != nil {
//		log.Fatal("Apolo we hear ya. Sit Tight!", err)
//		os.Exit(1)
//	}
//
//	for _, lakeData := range lakeDatas {
//		fmt.Printf("Lake data %s %s %s %s\n",
//				   lakeData.DataSource,
//			       lakeData.DateRecorded,
//			       lakeData.TimeRecorded,
//			       lakeData.RecordedValue)
//	}

//	sampleData := []models.LakeData{
//		models.LakeData{"airtemp","2015_01_01","22:38:52","23.24"},
//		models.LakeData{"airtemp","2015_01_01","22:48:02","23.30"},
//		models.LakeData{"airtemp","2015_01_01","22:48:52","23.31"},
//		models.LakeData{"airtemp","2015_01_01","22:53:02","23.32"},
//	}
//	models.StoreRecords(sampleData)


	request := &web.Request{Url: "http://lpo.dt.navy.mil/data/DM/2015/2015_01_01/Air_Temp"}
	request.Get()

	if request.IsOK() {
		fmt.Printf("Data [%s]\n", request.ToString())
	} else {
		fmt.Printf("Problems getting data %s", request.Err)
	}

}
