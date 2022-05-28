package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

const APIkey = "8ac4791f20af452e978111630222805"
const AksaiLatLong = "51.17,53.01"

type Response struct {
	Location `json:"location"`
	Current  `json:"current"`
}

type Location struct {
	Name    string  `json:"name"`
	Country string  `json:"country"`
	Lat     float32 `json:"lat"`
	Long    float32 `json:"lon"`
	TzID    string  `json:"tz_id"`
}

type Current struct {
	Temperature float32 `json:"temp_c"`
	Wind        float32 `json:"wind_kph"`
	WindDir     string  `json:"wind_dir"`
	Condition   `json:"condition"`
}

type Condition struct {
	Text string `json:"text"`
}

func main() {

	var query string
	if len(os.Args) < 2 {
		query = "auto:ip"
	} else {
		query = os.Args[1]
		if query == "Aksai" {
			query = AksaiLatLong
		}
	}

	endPoint, err := url.Parse("http://api.weatherapi.com/v1/current.json")
	if err != nil {
		log.Fatal(err)
	}

	values := url.Values{}
	values.Add("key", APIkey)
	values.Add("q", query)
	values.Add("aqi", "no")

	endPoint.RawQuery = values.Encode()
	resp, err := http.Get(endPoint.String())
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	weather := Response{}
	if err = json.Unmarshal(body, &weather); err != nil {
		log.Fatal(err)
	}

	var location string
	if weather.Country == weather.Name {
		location = weather.Country
	} else {
		location = fmt.Sprintf("%v/%v", weather.Country, weather.Name)
	}

	if location == "" {
		log.Fatal("unknown location")
	}

	fmt.Printf("temp: %v °C, %s, wind: %v kph %q (%s)\n",
		weather.Temperature, weather.Text, weather.Wind, weather.WindDir,
		location)
}
