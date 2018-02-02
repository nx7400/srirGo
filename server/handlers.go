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

type SourceCodeResponse struct {
	Status string // TODO refactor to enum
	Output string
}

var sourceCodesMap = make(map[uint64]string)

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

}
