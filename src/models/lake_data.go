package models

import (
	"fmt"
	"db"
	_ "log"
	"strings"
	"web"
	"strconv"
	"sort"
)

const (
	// http://lpo.dt.navy.mil/data/DM/2015/2015_01_01/Air_Temp
	ROOT_URL = "http://lpo.dt.navy.mil/data/DM/%s/%s/%s"
)


const (
	AirTemp         = "Air_Temp"
	BarometricPress = "Barometric_Press"
	Wind_Speed      = "Wind_Speed"
)

var DATASOURCE_TYPES = []string {
	AirTemp,
	BarometricPress,
	Wind_Speed,
}


func init() {
	//fmt.Printf("Initializing Lake Data Model...\n")
}

type LakeData struct {
	DataSource     string
	DateRecorded   string
	TimeRecorded   string
	RecordedValue  string
}

type DataRecs []LakeData


/*

    #####  db.Insertable Interface Impl ######

*/

func (records DataRecs) Len() int {
	return len(records)
}

func (records DataRecs) GetData(recordIdx int) []interface{} {

	selectedRecord := records[recordIdx]
	values := make([]interface{}, 3)

	values[0] = selectedRecord.DataSource
	values[1] = formatDate(selectedRecord.DateRecorded, selectedRecord.TimeRecorded)
	values[2] = selectedRecord.RecordedValue

	return values

}

/*

    #####  statistics.Mean Interface Impl ######

*/
func (records DataRecs) Mean() (float64) {

	total := 0.0
	for _, record := range records {
		//fmt.Printf("converting to float: %s", record.RecordedValue)
		recordedValue := fromStringToFloat(record.RecordedValue)
		total += recordedValue
	}

	if total == 0.0 {
		return total
	}

	meanValue := total / float64(len(records))
	meanValue = fromStringToFloat(fmt.Sprintf("%.2f", meanValue))
	return meanValue
}

/*

    #####  statistics.Mediam Interface Impl ######

*/
func (records DataRecs) Median() (float64) {


	medianValue := 0.0

	// sort the numbers from lowest to highest
	sort.Sort(records)

	numberOfItems := len(records)

	// For an even number of values, calculate the average of the two central numbers
	if numberOfItems % 2 == 0 {

		firstIdx := numberOfItems / 2
		secondIdx := firstIdx + 1

		//fmt.Printf(" numberOfItems %d  idx1st %d 2ndidx %d\n", numberOfItems, firstIdx, secondIdx)
		recordedValueFirst  := fromStringToFloat(records[firstIdx].RecordedValue)
		recordedValueSecond := fromStringToFloat(records[secondIdx].RecordedValue)

		if recordedValueSecond == 0.0 { // To avoid Division by Zero just return 0.0
			return 0.0
		}
		//fmt.Printf("Odd number Median %f %f\n", recordedValueFirst, recordedValueSecond)
		medianValue = (recordedValueFirst + recordedValueSecond) / float64(2)
		//fmt.Printf("Meian value is %f \n", medianValue)


	} else { // For an odd number of values, just take the middle number

		idx := (numberOfItems / 2) + 1
		//fmt.Printf(" numberOfItems %d  idx %d", numberOfItems, idx)
		medianValue = fromStringToFloat(records[idx].RecordedValue)
		//fmt.Printf("Meian value is %f \n", medianValue)
	}

	medianValue = fromStringToFloat(fmt.Sprintf("%.2f", medianValue))

	return medianValue
}

/*
	 #####  sort.Interface Impl ######
*/

func (records DataRecs) Swap(i, j int) {
	records[i], records[j] = records[j], records[i]
}

func (records DataRecs) Less(i, j int) bool {

	leftValue := fromStringToFloat(records[i].RecordedValue)
	rightValue := fromStringToFloat(records[j].RecordedValue)

	return leftValue < rightValue
}


func GetDBRecordsFor(dateForData string, recordType string) (DataRecs, error) {

    returnedRecords, err := db.Query("", recordType, dateForData)

	if err != nil {
		return nil, err
	}

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

func CheckDBRecordsFor(date string, recordType string) (bool, error) {

	recordCount, err := db.CountQuery(recordType, date)
	return recordCount > 0, err

}


func StoreRecords(inputRecs DataRecs) (error) {

	err := db.InsertData(inputRecs)

	return err
}

func ParseData(dataSourceType string, rawData string) (DataRecs) {

	parsedData := []LakeData{}

	for _, line  := range strings.Split(rawData, "\n") {
		if line == "" {
			continue
		}

		columns := strings.Split(line, " ")
		currentLakeData := LakeData{
			        DataSource:dataSourceType,
			        DateRecorded:columns[0],
			        TimeRecorded:columns[1],
	    }

		// Sometimes there is an extra 'space' after the second column
		// which causes the columns to be more then 3
		var lastColumns string
		for idx := 2; idx < len(columns); idx++ {
			lastColumns += columns[idx];
		}
		currentLakeData.RecordedValue = lastColumns

		parsedData = append(parsedData, currentLakeData)
	}

	return parsedData

}

func FetchData(date string, dataSourceType string) {

	dateParam := strings.Replace(date, "-", "_", len(date))
	yearPart  := strings.Split(dateParam, "_")[0]

	request := &web.Request{}
	url := fmt.Sprintf(ROOT_URL, yearPart, dateParam, dataSourceType)
	request.Url = url
	request.Get()

	if !request.IsOK() {
		fmt.Printf("Unable to get data for datasourceType: %s , date: %s , error: %s",
			dataSourceType,
			date,
			request.Err)
	}

	lakeDatas := ParseData(dataSourceType,request.ToString())
	StoreRecords(lakeDatas)

}

/*
	#### PRIVATE METHODS
*/

func formatDate(dateStr string, timeStr string) (string) {

	dateStr = strings.Replace(dateStr, "_", "-", len(dateStr))

	return dateStr + " " + timeStr

}

func fromStringToFloat(floatStr string) (float64) {

	if floatStr == "" {
		return 0.0
	}
	convertedValue, err := strconv.ParseFloat(floatStr, 64)
	if err != nil {
		convertedValue = 0.0
	}

	return convertedValue
}


