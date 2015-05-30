package models

import (
	"fmt"
	"db"
	"log"
)

func init() {
	fmt.Printf("Initializing Lake Data Model...\n")
}

type LakeData struct {
	DateRecorded   string
	TimeRecorded   string
	RecordedValue  string
}

type DataRecs []LakeData

func (records DataRecs) Len() int {
	return len(records)
}

func (records DataRecs) GetData(recordIdx int) []interface{} {

	selectedRecord := records[recordIdx]
	values := make([]interface{}, 3)

	values = append(values,selectedRecord.DateRecorded )
	values = append(values,selectedRecord.TimeRecorded )
	values = append(values,selectedRecord.RecordedValue )
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
		LakeData{"2015_01_01","22:38:52","23.24"},
		LakeData{"2015_01_01","22:48:02","23.30"},
		LakeData{"2015_01_01","22:48:52","23.31"},
		LakeData{"2015_01_01","22:53:02","23.32"},
	}, nil

}

func StoreRecords(inputRecs DataRecs) (error) {


	for _, inputRec := range inputRecs {

		fmt.Printf("Inserting %s, %s, %s\n",
			       inputRec.DateRecorded,
			       inputRec.TimeRecorded,
			       inputRec.RecordedValue)

	}

	return nil
}

