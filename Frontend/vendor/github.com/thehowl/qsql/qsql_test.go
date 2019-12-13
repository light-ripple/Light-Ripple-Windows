package qsql_test

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/thehowl/qsql"
)

/*
database:

CREATE DATABASE qsql;
USE qsql;
CREATE TABLE qsql_test(
	id INT(11) NOT NULL AUTO_INCREMENT,
	potato VARCHAR(128) NOT NULL,
	apple VARCHAR(128) NOT NULL,
	PRIMARY KEY(id)
);
INSERT INTO qsql_test(potato, apple) VALUES ("test", "test");
*/

func TestQuery(t *testing.T) {
	db, err := qsql.Open("mysql", "root@/qsql")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	data, err := db.Query("SELECT * FROM qsql_test")
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range data {
		if row["potato"] != "test" || row["apple"] != "test" {
			t.Fatal("Expected row to have potato=test and apple=test, got", row, "instead")
		}
	}
}

func TestQueryRow(t *testing.T) {
	db, err := qsql.Open("mysql", "root@/qsql")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	row, err := db.QueryRow("SELECT * FROM qsql_test")
	if err != nil {
		t.Fatal(err)
	}
	if row["potato"] != "test" || row["apple"] != "test" {
		t.Fatal("Expected row to have potato=test and apple=test, got", row, "instead")
	}
}

func Example() {
	db, err := qsql.Open("mysql", "root@/")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	row, err := db.QueryRow("SELECT 5 AS test, 1 AS test_bool, 13.37 AS test_float")
	if err != nil {
		panic(err)
	}
	fmt.Printf(
		"test: %d | test_bool: %v | test_float: %.3f\n",
		row["test"].Int(), row["test_bool"].Bool(), row["test_float"].Float64(),
	)
	// Output: test: 5 | test_bool: true | test_float: 13.370
}
