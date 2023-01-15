package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"strings"
	"unicode"
)

var (
	ddl    = `CREATE TABLE IF NOT EXISTS phone_numbers ("number" VARCHAR(50));`
	insert = `INSERT INTO phone_numbers (number) VALUES
                                           ('1234567890'),('123 456 7891'),('(123) 456 7892'),
                                           ('(123) 456-7893'),('123-456-7894'),('123-456-7890'),
                                           ('1234567892'),('(123)456-7892)');`
)

type phoneNumber struct {
	oldFmt string
	newFmt string
}

func main() {
	if err := run(); err != nil {
		log.Fatal("Error during run: ", err)
	}
}

func run() error {
	cxn, err := connect()
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	err = initDb(cxn)
	if err != nil {
		return fmt.Errorf("unable to initialize phone number table: %w", err)
	}
	records, err := loadRecords(cxn)
	if err != nil {
		return fmt.Errorf("unable to load records: %w", err)
	}
	return updateTable(cxn, normalizePhoneNumbers(records))
}

func connect() (*sql.DB, error) {
	cxn, err := sql.Open("pgx", "user=postgres password=mysecretpassword host=localhost port=5432")
	if err != nil {
		return nil, fmt.Errorf("couldn't open connection: %w", err)
	}
	if err := cxn.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}
	return cxn, nil
}

func initDb(cxn *sql.DB) error {
	_, err := cxn.Exec(ddl)
	if err != nil {
		return fmt.Errorf("unable to perform DDL for table: %w", err)
	}

	_, err = cxn.Exec(insert)
	if err != nil {
		return fmt.Errorf("unable to populate table with data: %w", err)
	}
	return nil
}

func loadRecords(cxn *sql.DB) ([]string, error) {
	rows, err := cxn.Query("SELECT * FROM phone_numbers")
	if err != nil {
		return nil, fmt.Errorf("loadRecords: %w", err)
	}
	var ret []string
	for rows.Next() {
		var num string
		if err := rows.Scan(&num); err == nil {
			ret = append(ret, num)
		}
	}
	return ret, rows.Close()
}

func normalizePhoneNumbers(numbers []string) []phoneNumber {
	ret := make([]phoneNumber, 0, len(numbers))
	for _, v := range numbers {
		norm := strings.Map(func(r rune) rune {
			if unicode.IsDigit(r) {
				return r
			}
			return -1
		}, v)
		ret = append(ret, phoneNumber{oldFmt: v, newFmt: norm})
	}
	return ret
}

func updateTable(cxn *sql.DB, nums []phoneNumber) error {
	for _, v := range nums {
		if err := processNumber(cxn, v); err != nil {
			return fmt.Errorf("updateTable: %w", err)
		}
	}
	return nil
}

func processNumber(cxn *sql.DB, num phoneNumber) error {
	row := cxn.QueryRow("SELECT COUNT(*) FROM phone_numbers WHERE number=$1", num.newFmt)
	cnt := 0
	err := row.Scan(&cnt)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("processNumber: %w", err)
	}
	if cnt == 0 {
		return updateRecord(cxn, num)
	} else {
		return deleteRecord(cxn, num)
	}
}

func updateRecord(cxn *sql.DB, num phoneNumber) error {
	_, err := cxn.Exec(`UPDATE phone_numbers SET number=$1 WHERE number=$2`, num.newFmt, num.oldFmt)
	if err != nil {
		return fmt.Errorf("updateRecord: %w", err)
	}
	return nil
}

func deleteRecord(cxn *sql.DB, num phoneNumber) error {
	_, err := cxn.Exec(`DELETE FROM phone_numbers WHERE number=$1`, num.oldFmt)
	if err != nil {
		return fmt.Errorf("deleteRecord: %w", err)
	}
	return nil
}
