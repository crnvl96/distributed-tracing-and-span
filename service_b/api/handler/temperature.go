package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/crnvl96/distributed-tracing-and-spam/service_b/internal"
	"github.com/crnvl96/distributed-tracing-and-spam/service_b/pkg/validation"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type Body struct {
	Zipcode string `json:"zipcode"`
}

func CalculateTemperature(w http.ResponseWriter, r *http.Request) {
	var body Body

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	marshalled, err := json.Marshal(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	fmt.Println(body)
	err = validation.ValidadeZipCode(ctx, w, marshalled)

	defer r.Body.Close()

	if err != nil && err.Error() == "invalid zipcode" {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	addr, err := internal.GetZipCode(body.Zipcode, ctx)

	if *addr == (internal.Address{}) {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	weather, err := internal.GetWeather(addr.City, ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&weather)
}
