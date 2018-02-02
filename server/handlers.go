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
	"github.com/hattya/go.diff"
)

// Source code response.
type SourceCodeResponse struct {
	Status string // TODO refactor to enum
	Output string
}

var sourceCodesMap = make(map[uint64]string)

// AddSourceCode adds source code to the database. Processes source code passed within r HTTP request
// adds it to database and if no error occurs assigns id and sends it back to the client
// within w HTTP response. HTTP StatusOK is set if source code has been successfuly added.
func AddSourceCode(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Add Source Code!!!")

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
		sourceCodesMap[id] = pathToSourceCode //it's not necessary in current implementation, but may by handy later
	}

	//fmt.Print(string(code))
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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	var response SourceCodeResponse
	sourceCodeId, _ := binary.Uvarint(body)
	pathToSourceCode, ok := sourceCodesMap[sourceCodeId] // TODO check if requested Id exist

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
			fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			response = SourceCodeResponse{"FAILED", fmt.Sprint(err) + ": " + stderr.String()}
		} else {
			fmt.Println("Result: " + out.String() + "OK")
			response = SourceCodeResponse{"SUCCESS", " "}
		}

	}

	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//os.Stdout.Write(js)
	//fmt.Println("-v", response)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

// RunSourceCode runs source code denoted by an id sent within r HTTP request. Standard output is
// catched and send back to the client within HTTP response. If program has failed
// to execute HTTP status is set to StatusBadRequest. Otherwise it is set to StatusOK. 
func RunSourceCode(w http.ResponseWriter, r *http.Request) {
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
		fmt.Printf("Result: %s", out)
		w.WriteHeader(http.StatusOK)
		w.Write(out)
	}
}

// CompareSourceCode compares given source codes.
// In case of error HTTP status is set to StatusInternalServerError. Otherwise it is
// being set to StatusOK. 
func CompareSourceCode(w http.ResponseWriter, r *http.Request) {
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

		code2, err2 := ioutil.ReadFile("receivedSourceCodes/1517531559919090900.go")
		if err2 != nil {
			panic(err2)
		}

		df := diff.Bytes(code1, code2)

		diff := df[0].Del + df[0].Ins

		if diff == 0 {
			response = SourceCodeResponse{"REPORT", "Source codes are the same"}
		} else {
			response = SourceCodeResponse{"REPORT", "Difference between codes: " + strconv.Itoa(diff) + " bytes"}
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

