//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	url := "http://localhost:8080/write"     // Replace with your destination URL
	jsonStr := []byte(`{"content":"value"}`) // Replace with your request body JSON payload
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	// Handle response
	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))

}
