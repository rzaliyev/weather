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

	cfg := GetConfig()

	flag.Parse()

	if len(flag.Args()) < 1 {
		getWeather(cfg, "auto:ip")
	} else {
		var wg sync.WaitGroup
		for _, v := range flag.Args() {
			wg.Add(1)
			go func(query string) {
				getWeather(cfg, query)
				wg.Done()
			}(v)
		}
		wg.Wait()
	}
}

func getWeather(cfg *Config, query string) {

	query = convertQuery(cfg, query)

	resp, err := http.Get(createAPIQuery(cfg, query))
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
		printWeather(cfg, *weather)
	}

}

func GetCurrentWeather(b []byte) (resp *Response) {
	resp = &Response{}
	if err := json.Unmarshal(b, resp); err != nil {
		log.Fatal(err)
	}
	return
}

// find special cases form the configuration
func convertQuery(cfg *Config, query string) string {

	for _, v := range cfg.Cities {
		if strings.EqualFold(query, v.Name) {
			return fmt.Sprintf("%v,%v", v.Lat, v.Lon)
		}
	}
	return query
}

func createAPIQuery(cfg *Config, query string) string {

	var APIurl = APIcurrent
	values := url.Values{}
	values.Add("key", cfg.APIkey)
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

func printWeather(cfg *Config, weather Response) {

	if *days > 0 {
		printForecastWeather(cfg, weather)
	} else {
		printCurrentWeather(cfg, weather)
	}
}

func printCurrentWeather(cfg *Config, weather Response) {
	fmt.Printf("Current weather for %s:\n", getLocation(cfg, weather))
	fmt.Printf("temp: %v °C (feels %v °C), wind: %v kph %q, UV: %v, %s\n",
		weather.Temperature, weather.FeelsLike, weather.Wind, weather.WindDir, weather.UV, weather.Text)
}

func getLocation(cfg *Config, weather Response) (location string) {

	for _, v := range cfg.Cities {
		if v.Lat == weather.Lat && v.Lon == weather.Long {
			weather.Name = v.Name
		}
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

func printForecastWeather(cfg *Config, weather Response) {
	fmt.Printf("Forecast for %s:\n", getLocation(cfg, weather))
	for _, day := range weather.ForecastDays {
		fmt.Printf("[%v] min: %v °C, max: %v °C, rain: %d%% wind: %v kph, UV: %v, %s\n",
			day.Date.Format("2006-01-02"),
			day.Day.MinTemp, day.Day.MaxTemp, day.Day.ChanceOfRain, day.Day.MaxWind, day.Day.UV, day.Day.Condition.Text)
	}
}
