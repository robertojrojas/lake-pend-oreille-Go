package db

import (
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

const (

	LPO_DB_NAME       = "./lpo.db"
	LPO_TABLE_NAME    = "LAKE_DATA"
    FIND_TABLE_QUERY  = "SELECT name FROM sqlite_master WHERE type = 'table' and name = '%s'"
	CREATE_TABLE_STMT = `
		CREATE TABLE %s ( id INTEGER PRIMARY KEY, type TEXT, stamp TEXT NOT NULL, value INTEGER );
	    CREATE UNIQUE INDEX typestamp ON %s(type, stamp);
	`
	INSERT_DATA_STMT  = "INSERT INTO %s (type, stamp, value) VALUES (?, ?, ?)"
	COUNT_DATA_QUERY  = "SELECT COUNT(*) FROM %s WHERE type = ? AND stamp LIKE ? || '%'"
	GET_DATA_QUERY    = "SELECT type, stamp, value FROM %s WHERE type = ? AND stamp LIKE ? || '%s'"
	TYPE_COL          = "type"
	STAMP_COL         = "stamp"
	VALUE_COL         = "col"

)

type Insertable interface {
	 Len() int
	 GetData(recordIdx int) []interface{}
}

func init() {
	fmt.Printf("Initializing db...\n")
	CreateTableIfNeeded(LPO_TABLE_NAME)

}

func GetDBConnection() (*sql.DB, error) {

	//TODO: Need to change the database location from an environment var.
	db, err := sql.Open("sqlite3", LPO_DB_NAME)
	return db, err
}


/*
  Returns true if the table was create, otherwise false.
  If there are any errors, false and the error are returned
*/
func CreateTableIfNeeded(tableName string) (bool, error) {

	db, err := GetDBConnection()
	if err != nil {
		log.Fatal(err)
		return false, err
	}

	findTableQuery := fmt.Sprintf(FIND_TABLE_QUERY, tableName)
	fmt.Printf("Find table with query | %s |\n", findTableQuery)
	rows, err := db.Query(findTableQuery)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
		return false, err
	}

	if !rows.Next() {
		fmt.Printf("About to create Table %s\n", tableName)
		sqlStmt := fmt.Sprintf(CREATE_TABLE_STMT, tableName, tableName)
		fmt.Print("creating table with %s", sqlStmt)
		_, err = db.Exec(sqlStmt)
		if err != nil {
			log.Printf("%q: %s\n", err, sqlStmt)
			return false, err
		}
		return true, nil
	} else {
		fmt.Printf("Table %s was created already!\n", LPO_DB_NAME)
	}

	return false, nil

}

func InsertData(dataToInsert Insertable) (error) {

	db, err := GetDBConnection()
	if err != nil {
		log.Fatal(err)
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmtStr := fmt.Sprintf(INSERT_DATA_STMT, LPO_TABLE_NAME)
	fmt.Printf("stmtStr %s\n", stmtStr)
	stmt, err := tx.Prepare(stmtStr)
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	for i := 0; i < dataToInsert.Len(); i++ {
		_, err = stmt.Exec(dataToInsert.GetData(i)...)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	tx.Commit()

	return nil
}

func Query(queryStr string, args ...interface{}) ([]map[string]string, error) {

	db, err := GetDBConnection()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	//TODO: Need a better way to handle this
	findQuery := fmt.Sprintf(GET_DATA_QUERY, LPO_TABLE_NAME, "%")
	fmt.Printf("args %T %v\n", args, args)

	rows, err := db.Query(findQuery, args...)
	if err != nil {
		log.Fatal("Query Failed!", err)
		return nil, err
	}
	defer rows.Close()

	retValues := []map[string]string{}

	for rows.Next() {
		var dataSource,
		    stamp,
		    recordedValue string
		rows.Scan(&dataSource, &stamp, &recordedValue)

		currentRow := map[string]string {
			TYPE_COL  : dataSource,
			STAMP_COL : stamp,
			VALUE_COL :recordedValue,

		}
		retValues = append(retValues, currentRow)

	}

	return retValues, nil

}


