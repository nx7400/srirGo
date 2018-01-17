package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func main() {

	code, err := ioutil.ReadFile("codesToSend/testCode.go")
	if err != nil {
		panic(err)
	}

	addSourceCodeUrl := "http://localhost:8080/add_source_code"
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

	recivedId, _ := binary.Uvarint(body)

	fmt.Println("response Body: recived Id:", strconv.FormatUint(recivedId, 10))
}
