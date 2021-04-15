package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a App

const tableCreationQuery = `DROP TABLE IF EXISTS product;
	CREATE TABLE product (
		id SERIAL,
		name TEXT NOT NULL,
		price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
		createdOn TIMESTAMP NOT NULL DEFAULT NOW(),
		CONSTRAINT pk_product_id PRIMARY KEY (id),
		CONSTRAINT uk_product_name UNIQUE(name)
	)`

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM product")
	a.DB.Exec("ALTER SEQUENCE product_id_seq RESTART WITH 1")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code: %d. Got %d\n", expected, actual)
	}
}

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize(
		os.Getenv("TEST_DB_USERNAME"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"),
	)

	ensureTableExists()
	code := m.Run()
	fmt.Println(code)
	clearTable()
	os.Exit(code)
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/products", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}
