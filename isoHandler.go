package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/mofax/iso8583"
	"github.com/rivo/uniseg"
)

// Handle all ISO Client request

// Process ISO message in body request
func sendIso(writer http.ResponseWriter, request *http.Request) {

	var response Response
	var iso Iso8583

	// Read body request
	reqBody, _ := ioutil.ReadAll(request.Body)
	req := string(reqBody)
	log.Printf("ISO Message: %v\n", req)

	// Produce event
	err := doProducer(broker, topic1, req)

	if err != nil {
		errDesc := fmt.Sprintf("Failed sent to Kafka\nError: %v", err)
		log.Println(err)
		response.ResponseCode, response.ResponseDescription = 500, errDesc
		jsonFormatter(writer, response, 500)
	} else {
		// Read response
		msg, err := consumeResponse(broker, group, []string{topic2})
		if err != nil {
			errDesc := fmt.Sprintf("Failed to get response from Kafka\nError: %v", err)
			log.Println(err)
			response.ResponseCode, response.ResponseDescription = 500, errDesc
			jsonFormatter(writer, response, 500)
		} else {

			// Return empty response
			if msg == "" {
				errDesc := "Got empty response"
				log.Println(errDesc)
				response.ResponseCode, response.ResponseDescription = 500, errDesc
				jsonFormatter(writer, response, 500)
			} else {

				// Parse response string to ISO8583 data
				header := msg[0:4]
				data := msg[4:]

				isoStruct := iso8583.NewISOStruct("spec1987.yml", false)

				isoParsed, err := isoStruct.Parse(data)
				if err != nil {
					log.Printf("Error parsing iso message\nError: %v", err)
				}

				iso.Header, _ = strconv.Atoi(header)
				iso.MTI = isoParsed.Mti.String()
				iso.Hex, _ = iso8583.BitMapArrayToHex(isoParsed.Bitmap)

				iso.Message, err = isoParsed.ToString()
				if err != nil {
					log.Printf("Iso Parsed failed convert to string.\nError: %v", err)
				}

				//event := header + iso.Message

				iso.ResponseStatus.ResponseCode, iso.ResponseStatus.ResponseDescription = 200, "Success"
				jsonFormatter(writer, iso, 200)

			}

		}
	}

}

// Get response from mock server in ISO Format
func responseIso(message string) {

	var response Iso8583
	data := message[4:]

	isoStruct := iso8583.NewISOStruct("spec1987.yml", false)

	msg, err := isoStruct.Parse(data)
	if err != nil {
		log.Println(err)
	}

	// Convert ISO message to JSON format
	jsonIso := convertIsoToJson(msg)

	// Send JSON data to mock server
	serverResp := responseJson(jsonIso)

	// Convert response from JSON data to ISO8583 format
	isoParsed := convertJsonToIso(serverResp)

	// Change MTI response
	isoParsed.AddMTI("0210")

	isoMessage, _ := isoParsed.ToString()
	isoMTI := isoParsed.Mti.String()
	isoHex, _ := iso8583.BitMapArrayToHex(isoParsed.Bitmap)
	isoHeader := fmt.Sprintf("%04d", uniseg.GraphemeClusterCount(isoMessage))

	response.Header, _ = strconv.Atoi(isoHeader)
	response.MTI = isoMTI
	response.Hex = isoHex
	response.Message = isoMessage

	event := isoHeader + isoMessage
	log.Printf("\n\nResponse: \n\tHeader: %v\n\tMTI: %v\n\tHex: %v\n\tIso Message: %v\n\tFull Message: %v\n\n",
		response.Header,
		response.MTI,
		response.Hex,
		response.Message,
		event)

	// create file from response
	filename := "Response_to_" + isoParsed.Elements.GetElements()[3] + "@" + fmt.Sprintf(time.Now().Format("2006-01-02 15:04:05"))
	file := CreateFile("storage/response/"+filename, event)
	log.Println("File created: ", file)

	// Produce event
	err = doProducer(broker, topic2, event)
	if err != nil {
		log.Printf("Error producing message %v\n", message)
		log.Println(err)
	}

}
