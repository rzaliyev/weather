package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var debug = flag.Bool("debug", false, "show orgingal respons in JSON format")

func main() {

	flag.Parse()

	if len(flag.Args()) < 1 {
		getWeather("auto:ip")
	} else {
		var wg sync.WaitGroup
		for _, v := range flag.Args() {
			wg.Add(1)
			go func(query string) {
				getWeather(query)
				wg.Done()
			}(v)
		}
		wg.Wait()
	}
}

func getWeather(query string) {

	query = convertQuery(query)

	resp, err := http.Get(createAPIQuery(query))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	weather := GetCurrentWeather(body)

	if *debug {
		var bs []byte
		if bs, err = json.Marshal(weather); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(string(bs))
		}

	} else {
		printWeather(*weather)
	}

}

func GetCurrentWeather(b []byte) (resp *Response) {
	resp = &Response{}
	if err := json.Unmarshal(b, resp); err != nil {
		log.Fatal(err)
	}
	return
}

func convertQuery(query string) string {
	if val, ok := cities[strings.ToLower(query)]; ok {
		return val
	}
	return query
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
