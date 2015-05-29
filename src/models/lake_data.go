package models

import (
	"fmt"
	"db"
	"log"
)

func init() {
	fmt.Printf("Initializing Lake Data Model...")
}

type LakeData struct {
	DateRecorded   string
	TimeRecorded   string
	RecordedValue  string
}

func GetDBRecordsFor(dateForData string) ([]LakeData, error) {

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

func StoreRecords(inputRecs []LakeData) (error) {

	dbConnection, err := db.GetDBConnection()

	if err != nil {
		log.Fatal("Houston we have a problem!")
		return err
	}

	defer dbConnection.Close()

	for _, inputRec := range inputRecs {

		fmt.Printf("Inserting %s, %s, %s",
			       inputRec.DateRecorded,
			       inputRec.TimeRecorded,
			       inputRec.RecordedValue)

	}

	return nil
}

