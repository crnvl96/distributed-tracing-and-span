package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"go.opentelemetry.io/otel"
)

type CurrentWeather struct {
	TemperatureC float64 `json:"temp_c"`
	TemperatureF float64 `json:"temp_f"`
}

type WeatherData struct {
	Current CurrentWeather `json:"current"`
}

type Weather struct {
	Temp_C float64 `json:"temp_C"`
	Temp_F float64 `json:"temp_F"`
	Temp_K float64 `json:"temp_K"`
	City   string  `json:"city"`
}

func GetWeather(city string, ctx context.Context) (*Weather, error) {
	ctx, span := otel.GetTracerProvider().Tracer("weather_integration").Start(ctx, "weather_integration")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.weatherapi.com/v1/current.json?key="+os.Getenv("WEATHER_API_KEY")+"&q="+url.QueryEscape(city)+"&aqi=no", nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var weather WeatherData
	err = json.NewDecoder(res.Body).Decode(&weather)
	if err != nil {
		return nil, err
	}

	return &Weather{
		Temp_C: weather.Current.TemperatureC,
		Temp_F: weather.Current.TemperatureF,
		Temp_K: weather.Current.TemperatureC + 273,
		City:   city,
	}, nil
}
