package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func addSourceCode(serverBaseUrl string) uint64 {

	code, err := ioutil.ReadFile("codesToSend/testCode.go")
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

	fmt.Println("response Body:", body)

	return true

}

func main() {

	serverIpAddrPtr := flag.String("sa", "localhost", "a string")
	flag.Parse()

	serverBaseUrl := "http://" + *serverIpAddrPtr + ":8080"

	receivedId := addSourceCode(serverBaseUrl)

	_ = receivedId

	if checkSourceCode(serverBaseUrl, receivedId) {
		fmt.Println("SUCCESS")
	}

}
