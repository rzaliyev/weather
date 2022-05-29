package main

import (
	"encoding/json"
	"log"

	"github.com/rzaliyev/config"
)

type Config struct {
	APIkey      string `json:"key"`
	CurrentAPI  string `json:"current_api,omitempty"`
	ForecastAPI string `json:"forecast_api,omitempty"`
	Cities      []City `json:"cities,omitempty"`
}

type City struct {
	Name string  `json:"name"`
	Lat  float32 `json:"lat"`
	Lon  float32 `json:"lon"`
}

func GetConfig() *Config {
	bs, err := config.Get(".config/weather", "config")
	if err != nil {
		log.Fatal(err)
	}

	cfg := &Config{}
	if err = json.Unmarshal(bs, cfg); err != nil {
		log.Fatal(err)
	}

	return cfg
}
