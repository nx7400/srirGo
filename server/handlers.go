package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func AddProgram(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Add Program!!!")

	code, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	fmt.Print(string(code))

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

}

func CheckProgram(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Check Program!!!")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func RunProgram(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Run Program!!!")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func CompareProgram(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Compare Program!!!")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
