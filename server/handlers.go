package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
	"time"
)

var sourceCodesMap = make(map[uint64]string)

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

func CheckSourceCode(w http.ResponseWriter, r *http.Request) {

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
	cmd := exec.Command(app, "run", pathToSourceCode)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	} else {
		fmt.Println("Result: " + out.String() + "OK")
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func RunSourceCode(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Run Source Code!!!")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func CompareSourceCode(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Compare Source Code!!!")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
