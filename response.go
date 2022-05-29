package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type WeatherTime time.Time

func (wt *WeatherTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02 15:04", s)
	if err != nil {
		return err
	}
	*wt = WeatherTime(t)
	return nil
}

func (wt WeatherTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(wt).Format("2006-01-02 15:04"))
}

func (wt WeatherTime) Format(s string) string {
	t := time.Time(wt)
	return t.Format(s)
}

type WeatherDay time.Time

func (wd *WeatherDay) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*wd = WeatherDay(t)
	return nil
}

func (wd WeatherDay) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(wd).Format("2006-01-02"))
}

func (wd WeatherDay) Format(s string) string {
	t := time.Time(wd)
	return t.Format(s)
}

type WeatherBool bool

func (wb *WeatherBool) UnmarshalJSON(b []byte) error {
	bs := strings.Trim(string(b), "\"")
	i, err := strconv.Atoi(bs)
	if err != nil {
		return err
	}

	switch i {
	case 0:
		*wb = false
	case 1:
		*wb = true
	default:
		return fmt.Errorf("cannot unmarshal int into bool")
	}

	return nil
}

func (wb WeatherBool) MarshalJSON() ([]byte, error) {
	var result int
	switch wb {
	case true:
		result = 1
	case false:
		result = 0
	}
	return json.Marshal(result)
}

type Response struct {
	Location `json:"location"`
	Current  `json:"current"`
	Forecast `json:"forecast"`
}

type Location struct {
	Name           string      `json:"name"`
	Region         string      `json:"region"`
	Country        string      `json:"country"`
	Lat            float32     `json:"lat"`
	Long           float32     `json:"lon"`
	TzID           string      `json:"tz_id"`
	LocalTimeEpoch int         `json:"localtime_epoch"`
	LocalTime      WeatherTime `json:"localtime"`
}

type Current struct {
	LastUpdatedEpoch int         `json:"last_updated_epoch"`
	LastUpdated      WeatherTime `json:"last_updated"`
	Temperature      float32     `json:"temp_c"`
	IsDay            WeatherBool `json:"is_day"`
	Condition        `json:"condition"`
	Wind             float32 `json:"wind_kph"`
	WindDegree       int     `json:"wind_degree"`
	WindDir          string  `json:"wind_dir"`
	Pressure         float32 `json:"pressure_mb"`
	Precip           float32 `json:"precip_mm"`
	Humidity         int     `json:"humidity"`
	Cloud            int     `json:"cloud"`
	FeelsLike        float32 `json:"feelslike_c"`
	Vis              float32 `json:"vis_km"`
	UV               float32 `json:"uv"`
	Gust             float32 `json:"gust_kph"`
}

type Condition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
	Code int    `json:"code"`
}

type Forecast struct {
	ForecastDays []ForecastDay `json:"forecastday"`
}

type ForecastDay struct {
	Date  WeatherDay `json:"date"`
	Day   `json:"day"`
	Astro `json:"astro"`
}

type Day struct {
	MaxTemp      float32 `json:"maxtemp_c"`
	MinTemp      float32 `json:"mintemp_c"`
	MaxWind      float32 `json:"maxwind_kph"`
	ChanceOfRain int     `json:"daily_chance_of_rain"`
	Condition    `json:"condition"`
	UV           float32 `json:"uv"`
}

type Astro struct {
	Sunrise   string `json:"sunrise"`
	Sunset    string `json:"sunset"`
	MoonPhase string `json:"moon_phase"`
}
