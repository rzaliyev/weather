package main

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
