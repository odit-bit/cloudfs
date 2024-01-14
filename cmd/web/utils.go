package main

import (
	"log"
	"net/http"
)

// utils
func printErr(w http.ResponseWriter, err error, endpointScope string) {
	log.Printf("%v: %v \n", endpointScope, err)
	http.Error(w, "error", http.StatusInternalServerError)
}
