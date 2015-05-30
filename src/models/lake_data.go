package models

import (
	"fmt"
	"db"
	_ "log"
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

	return values

}

func GetDBRecordsFor(dateForData string, recordType string) (DataRecs, error) {

    returnedRecords, _ := db.Query("", dateForData, recordType)

	lakeDataValues := make([]LakeData, len(returnedRecords))

	for idx, record := range returnedRecords {
		dateTime := record[db.STAMP_COL]
		dateTimeParts := strings.Split(dateTime, " ")
		lakeDataValues[idx] = LakeData {
			     record[db.TYPE_COL],
			     dateTimeParts[0],
			     dateTimeParts[1],
			     record[db.VALUE_COL],
		}
	}

	return lakeDataValues, nil

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

