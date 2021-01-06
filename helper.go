package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// Response formatter
func jsonFormatter(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Create file for request/response
func CreateFile(fileName string, content string) string {

	if !strings.Contains(fileName, ".txt") {
		fileName += ".txt"
	}

	file, err := os.Create(fileName)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	defer file.Close()

	_, err = file.WriteString(content)

	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}

	return fileName

}

// Check existing file
func CheckExist(namaFile string) bool {

	info, err := os.Stat(namaFile)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Read existing file
func ReadFile(fileName string) string {

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Panicf("failed reading data from file: %s", err)
	}

	return string(data)

}
