package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
)

func main() {

	if len(os.Args) < 2 {
		getWeather("auto:ip")
	} else {
		var wg sync.WaitGroup
		for _, v := range os.Args[1:] {
			query := getQuery(v)
			wg.Add(1)
			go func() {
				getWeather(query)
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func getQuery(query string) string {
	if val, ok := cities[query]; ok {
		return val
	}
	return query
}

func getWeather(query string) {

	resp, err := http.Get(createAPIQuery(query))
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

	printWeather(weather)

}

func createAPIQuery(query string) string {
	endPoint, err := url.Parse(APIurl)
	if err != nil {
		log.Fatal(err)
	}

	values := url.Values{}
	values.Add("key", APIkey)
	values.Add("q", query)
	values.Add("aqi", "no")

	endPoint.RawQuery = values.Encode()

	return endPoint.String()
}

func printWeather(weather Response) {
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
