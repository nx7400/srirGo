package main

import (
	"bytes"
	"encoding/binary"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

func TestCheckSourceCode(t *testing.T) {

	var sourceCodeId = uint64(time.Now().UnixNano())

	idBuf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(idBuf, sourceCodeId)

	req, err := http.NewRequest("POST", "/add_source_code", bytes.NewBuffer(idBuf))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CheckSourceCode)

	handler.ServeHTTP(rr, req)

	//Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler AddSourceCode returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestRunSourceCodeHandler(t *testing.T) {
	const VALID_CODE_ID = 0
	const NONEXISTING_CODE_ID = 1

	sourceCodesMap[VALID_CODE_ID] = "tests/validCode.go"
	sourceCodesMap[NONEXISTING_CODE_ID] = "tests"

	rr := runSourceCodeRequest(t, NONEXISTING_CODE_ID)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler RunSourceCode returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	rr = runSourceCodeRequest(t, VALID_CODE_ID)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler RunSourceCode returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func runSourceCodeRequest(t *testing.T, sourceCodeId uint64) *httptest.ResponseRecorder {
	idBuf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(idBuf, sourceCodeId)

	req, err := http.NewRequest("POST", "/run_source_code", bytes.NewBuffer(idBuf))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RunSourceCode)

	handler.ServeHTTP(rr, req)
	return rr
}

func TestCompareSourceCode(t *testing.T) {

	var sourceCodeId = uint64(time.Now().UnixNano())

	idBuf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(idBuf, sourceCodeId)

	req, err := http.NewRequest("POST", "/add_source_code", bytes.NewBuffer(idBuf))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CompareSourceCode)

	handler.ServeHTTP(rr, req)

	//Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler AddSourceCode returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
