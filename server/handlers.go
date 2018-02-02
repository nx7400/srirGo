package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	diff "github.com/hattya/go.diff"
)

// Source code response.
type SourceCodeResponse struct {
	Status string // TODO refactor to enum
	Output string
}

var sourceCodesMap = make(map[uint64]string)
var lastSourceCodeId = uint64(0)

// AddSourceCode adds source code to the database. Processes source code passed within r HTTP request
// adds it to database and if no error occurs assigns id and sends it back to the client
// within w HTTP response. HTTP StatusOK is set if source code has been successfuly added.
func AddSourceCode(w http.ResponseWriter, r *http.Request) {

	fmt.Println()

	code, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	var id = uint64(time.Now().UnixNano())
	var pathToSourceCode = "receivedSourceCodes/" + strconv.FormatUint(id, 10) + ".go"

	err = ioutil.WriteFile("receivedSourceCodes/"+strconv.FormatUint(id, 10)+".go", code, 0644)
	if err != nil {
		panic(err)
	} else {
		sourceCodesMap[id] = pathToSourceCode
		fmt.Println("Added source code with id: " + strconv.FormatUint(id, 10))
	}

	idBuf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(idBuf, id)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(idBuf)

}

// CheckSourceCode checks source code denoted by an id sent within r HTTP request. Makes an attempt to compile
// program. If source code has been not found, then "NOT_FOUND" message is sent inside a body of
// HTTP response. If source code failed to compile, then "FAILED" message is sent. If source
// code was compiled succesfuly "SUCCESS" message is being send. In case of error
// HTTP status is set to StatusInternalServerError. Otherwise StatusOK is being sent.
func CheckSourceCode(w http.ResponseWriter, r *http.Request) {

	fmt.Println()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	var response SourceCodeResponse
	sourceCodeId, _ := binary.Uvarint(body)
	pathToSourceCode, ok := sourceCodesMap[sourceCodeId]

	if !ok {
		response = SourceCodeResponse{"NOT_FOUND", "Source Code Id: " + strconv.FormatUint(sourceCodeId, 10) + " not found"}
	} else {

		app := "go"
		cmd := exec.Command(app, "build", pathToSourceCode)

		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		err = cmd.Run()
		if err != nil {
			fmt.Println("Build result: " + fmt.Sprint(err) + ": " + stderr.String())
			response = SourceCodeResponse{"FAILED", fmt.Sprint(err) + ": " + stderr.String()}
		} else {
			fmt.Println("Build result: " + out.String() + "OK")
			response = SourceCodeResponse{"SUCCESS", " "}
		}

	}

	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

// RunSourceCode runs source code denoted by an id sent within r HTTP request. Standard output is
// catched and send back to the client within HTTP response. If program has failed
// to execute HTTP status is set to StatusBadRequest. Otherwise it is set to StatusOK.
func RunSourceCode(w http.ResponseWriter, r *http.Request) {

	fmt.Println()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	sourceCodeId, _ := binary.Uvarint(body)
	pathToSourceCode := sourceCodesMap[sourceCodeId]

	app := "go"
	out, err := exec.Command(app, "run", pathToSourceCode).Output()

	if err != nil {
		fmt.Println("Run failed: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		fmt.Printf("Program result: %s", out)
		w.WriteHeader(http.StatusOK)
		w.Write(out)
	}
}

// CompareSourceCode compares given source codes.
// In case of error HTTP status is set to StatusInternalServerError. Otherwise it is
// being set to StatusOK.
func CompareSourceCode(w http.ResponseWriter, r *http.Request) {

	fmt.Println()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	var response SourceCodeResponse
	sourceCodeId, _ := binary.Uvarint(body)
	pathToSourceCode, ok := sourceCodesMap[sourceCodeId]
	if !ok {
		response = SourceCodeResponse{"NOT_FOUND", "Source Code Id: " + strconv.FormatUint(sourceCodeId, 10) + " not found"}
	} else {

		code1, err := ioutil.ReadFile(pathToSourceCode)
		if err != nil {
			panic(err)
		}

		var code2 []byte

		if lastSourceCodeId == 0 {
			code2 = []byte("A")
		} else {
			pathToComparedSourceCode, ok := sourceCodesMap[lastSourceCodeId]

			if !ok {
				fmt.Println("NOT FOUND", strconv.FormatUint(lastSourceCodeId, 10))
			} else {
				fmt.Println("FOUND", strconv.FormatUint(lastSourceCodeId, 10))
			}

			code2, err = ioutil.ReadFile(pathToComparedSourceCode)
			if err != nil {
				panic(err)
			}
		}

		df := diff.Bytes(code1, code2)

		var diff = int(0)
		if len(df) != 0 {
			diff = df[0].Del + df[0].Ins
		}

		if diff == 0 {
			response = SourceCodeResponse{"REPORT", "Source codes are the same"}
		} else {
			response = SourceCodeResponse{"REPORT", "Difference between codes: " + strconv.Itoa(diff) + " bytes"}
		}

		fmt.Println("Compare result: " + response.Output)

	}

	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)

	lastSourceCodeId = sourceCodeId

}
