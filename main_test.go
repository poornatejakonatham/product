package main

import (
	"bytes"
	"encoding/json"
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

func TestNonExistentProduct(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/product/10", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Product not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Product not found. Got '%s'", m["error"])
	}
}

func TestCreateProduct(t *testing.T) {
	clearTable()

	var jsonStr = []byte(`{"name": "test product", "price": 11.22, "created_on": "2021-04-15T19:00:00Z"}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test product" {
		t.Errorf("Expected product name to be 'test product'. Got '%v'", m["name"])
	}

	if m["price"] != 11.22 {
		t.Errorf("Expected product price to be '11.22'. Got '%v'", m["price"])
	}

	if m["created_on"] != "2021-04-15T19:00:00Z" {
		t.Errorf("Expected product created date to be '2021-04-15T19:00:00Z'. Got '%v'", m["created_on"])
	}

	// ID is compared to 1.0 because JSON unmarshling converts numbers to floats,
	// when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected product id to be '1'. Got '%v'", m["id"])
	}
}
