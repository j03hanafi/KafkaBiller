package main

import "github.com/gorilla/mux"

func server() *mux.Router {
	router := mux.NewRouter()

	// Endpoints

	// JSON Client

	// ISO Client
	router.HandleFunc("/payment/biller/iso", sendIso).Methods("POST")

	return router
}
