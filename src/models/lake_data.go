package models

import (
	"fmt"
	"db"
	"log"
	"strings"
)

func init() {
	fmt.Printf("Initializing Lake Data Model...\n")
}

type LakeData struct {
	DataSource     string
	DateRecorded   string
	TimeRecorded   string
	RecordedValue  string
}

type DataRecs []LakeData

func (records DataRecs) Len() int {
	return len(records)
}

func formatDate(dateStr string, timeStr string) (string) {

	dateStr = strings.Replace(dateStr, "_", "-", len(dateStr))

	return dateStr + " " + timeStr

}

func (records DataRecs) GetData(recordIdx int) []interface{} {

	selectedRecord := records[recordIdx]
	values := make([]interface{}, 3)

	values[0] = selectedRecord.DataSource
	values[1] = formatDate(selectedRecord.DateRecorded, selectedRecord.TimeRecorded)
	values[2] = selectedRecord.RecordedValue

	fmt.Printf("returning values len(%d)\n", len(values))

	return values

}

func GetDBRecordsFor(dateForData string) (DataRecs, error) {

	dbConnection, err := db.GetDBConnection()

	if err != nil {
		log.Fatal("Houston we have a problem!")
		return []LakeData{}, err
	}

	defer dbConnection.Close()

	return []LakeData{
		LakeData{"airtemp","2015_01_01","22:38:52","23.24"},
		LakeData{"airtemp","2015_01_01","22:48:02","23.30"},
		LakeData{"airtemp","2015_01_01","22:48:52","23.31"},
		LakeData{"airtemp","2015_01_01","22:53:02","23.32"},
	}, nil

}

func StoreRecords(inputRecs DataRecs) (error) {


	for _, inputRec := range inputRecs {
		fmt.Printf("Inserting %s, %s, %s\n",
			       inputRec.DateRecorded,
			       inputRec.TimeRecorded,
			       inputRec.RecordedValue)

	}

	err := db.InsertData(inputRecs)

	return err
}

