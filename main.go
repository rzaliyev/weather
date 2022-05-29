package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

var debug = flag.Bool("debug", false, "show orgingal respons in JSON format")
var days = flag.Int("days", 0, "set forecast days")
var forecast = flag.Bool("f", false, "query forecast for 3 days")

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

	var APIurl = APIcurrent
	values := url.Values{}
	values.Add("key", APIkey)
	values.Add("q", query)
	values.Add("aqi", "no")

	if *forecast {
		*days = 3
	}

	if *days > 0 {
		APIurl = APIforecast
		values.Add("days", strconv.Itoa(*days))
	}

	endPoint, err := url.Parse(APIurl)
	if err != nil {
		log.Fatal(err)
	}

	endPoint.RawQuery = values.Encode()

	return endPoint.String()
}

func printWeather(weather Response) {

	if *days > 0 {
		printForecastWeather(weather)
	} else {
		printCurrentWeather(weather)
	}
}

func printCurrentWeather(weather Response) {
	fmt.Printf("Current weather for %s:\n", getLocation(weather))
	fmt.Printf("temp: %v °C (feels %v °C), wind: %v kph %q, UV: %v, %s\n",
		weather.Temperature, weather.FeelsLike, weather.Wind, weather.WindDir, weather.UV, weather.Text)
}

func getLocation(weather Response) (location string) {

	if val, ok := geos[fmt.Sprintf("%v, %v", weather.Location.Lat, weather.Location.Long)]; ok {
		weather.Name = val
	}

	if weather.Country == weather.Name {
		location = fmt.Sprintf("%v/%v/Lat:%v,Lon:%v", weather.Country, weather.Region, weather.Location.Lat, weather.Location.Long)
	} else {
		location = fmt.Sprintf("%v/%v/%v", weather.Country, weather.Region, weather.Name)
	}

	if location == "" {
		log.Fatal("unknown location")
	}
	return
}

func printForecastWeather(weather Response) {
	fmt.Printf("Forecast for %s:\n", getLocation(weather))
	for _, day := range weather.ForecastDays {
		fmt.Printf("[%v] min: %v °C, max: %v °C, rain: %d%% wind: %v kph, UV: %v, %s\n",
			day.Date.Format("2006-01-02"),
			day.Day.MinTemp, day.Day.MaxTemp, day.Day.ChanceOfRain, day.Day.MaxWind, day.Day.UV, day.Day.Condition.Text)
	}
}
