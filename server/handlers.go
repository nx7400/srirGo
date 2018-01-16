package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var sourceCodesMap = make(map[int64]string)

func AddSourceCode(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Add Source Code!!!")

	code, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	var id = int64(time.Now().UnixNano())
	var pathToSourceCode = "receivedSourceCodes/" + strconv.FormatInt(id, 10) + ".go"

	err = ioutil.WriteFile("receivedSourceCodes/"+strconv.FormatInt(id, 10)+".go", code, 0644)
	if err != nil {
		panic(err)
	} else {
		sourceCodesMap[id] = pathToSourceCode //it's not necessary in current implementation, but may by handy later
	}

	//fmt.Print(string(code))

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

}

func CheckSourceCode(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Check Source Code!!!")

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
