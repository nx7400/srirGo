package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type SourceCodeResponse struct {
	Status string
	Output string
}

func addSourceCode(serverBaseUrl string, sourceCodePath string) uint64 {

	code, err := ioutil.ReadFile(sourceCodePath)
	if err != nil {
		panic(err)
	}

	addSourceCodeUrl := serverBaseUrl + "/add_source_code"
	req, err := http.NewRequest("POST", addSourceCodeUrl, bytes.NewBuffer(code))

	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)

	receivedId, _ := binary.Uvarint(body)

	fmt.Println("response Body: received Id:", strconv.FormatUint(receivedId, 10))

	return receivedId

}

func checkSourceCode(serverBaseUrl string, sourceCodeId uint64) bool {

	idBuf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(idBuf, sourceCodeId)

	checkSourceCodeUrl := serverBaseUrl + "/check_source_code"
	req, err := http.NewRequest("POST", checkSourceCodeUrl, bytes.NewBuffer(idBuf))

	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)

	var response SourceCodeResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	// TODO refactor to swich
	if response.Status == "SUCCESS" {
		fmt.Println("Code Check Passed")
		return true
	} else if response.Status == "FAILED" {
		fmt.Println("Code Check Failed. Output: " + response.Output)
		return false
	} else if response.Status == "NOT_FOUND" {
		fmt.Println("Code Check Failed. Output: " + response.Output)
		return false
	}
	return false
}

func runSourceCode(serverBaseUrl string, sourceCodeId uint64) string {

	idBuf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(idBuf, sourceCodeId)

	checkSourceCodeUrl := serverBaseUrl + "/run_source_code"
	req, err := http.NewRequest("POST", checkSourceCodeUrl, bytes.NewBuffer(idBuf))

	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("response Body:", body)

	return string(body[:])
}
func compareSourceCode(serverBaseUrl string, sourceCodeId uint64) bool {


	idBuf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(idBuf, sourceCodeId)

	checkSourceCodeUrl := serverBaseUrl + "/compare_source_code"
	req, err := http.NewRequest("POST", checkSourceCodeUrl, bytes.NewBuffer(idBuf))

	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	
var response SourceCodeResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	if response.Status == "REPORT" {
		fmt.Println("Report recived. Output: " + response.Output)
		return true
	} else {
		fmt.Println("Invalid response")
		return false
	}
}

func main() {

	serverIpAddrPtr := flag.String("sa", "localhost", "server address")
    sourceCodePath := flag.String("src", "codesToSend/testCode.go", "source code path")
	flag.Parse()

	serverBaseUrl := "http://" + *serverIpAddrPtr + ":8080"

    fmt.Println("Passing " + *sourceCodePath + " to process on " + serverBaseUrl)

	receivedId := addSourceCode(serverBaseUrl, *sourceCodePath)

	if checkSourceCode(serverBaseUrl, receivedId) {
		fmt.Println(runSourceCode(serverBaseUrl, receivedId))
	}

}
