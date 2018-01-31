package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddSourceCodeHandler(t *testing.T) {

	testSourceCode := []byte("Test")

	req, err := http.NewRequest("POST", "/add_source_code", bytes.NewBuffer(testSourceCode))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddSourceCode)

	handler.ServeHTTP(rr, req)

	//Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler AddSourceCode returned wrong status code: got %v want %v", status, http.StatusOK)
	}

}
