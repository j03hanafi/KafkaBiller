package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// Handle all JSON Client request

// Send request to mock server in JSON format
// Get response from mock server in JSON format
func responseJson(jsonIso Transaction) PaymentResponse {
	var response PaymentResponse

	// Initiate request body
	requestBody, err := json.Marshal(jsonIso)
	if err != nil {
		log.Fatalf("Preparing body request failed. Error: %v\n", err)
	}

	// Client setup for custom http request
	client := &http.Client{}

	log.Printf("Request to https://tiruan.herokuapp.com/biller\n")

	// Request to mock server
	req, err := http.NewRequest("GET", "https://tiruan.herokuapp.com/biller", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatalf("Failed to sent request to https://tiruan.herokuapp.com/biller. Error: %v\n", err)
	}

	// Check response from mock server
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to get response from https://tiruan.herokuapp.com/biller. Error: %v\n", err)
	}

	defer resp.Body.Close()

	log.Printf("Response from https://tiruan.herokuapp.com/biller\n")

	// Read response from mock server
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	// Get transactionData from mock server response

	return response
}
