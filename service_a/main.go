package main

import (
	"net/http"

	"github.com/crnvl96/distributed-tracing-and-spam/service_a/api/handler"
)

type RequestData struct {
	Zipcode string `json:"zipcode"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /zipcode", handler.ValidateZipCode)
	http.ListenAndServe(":8080", mux)
}
