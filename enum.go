package main

const (
	APIkey       = "8ac4791f20af452e978111630222805"
	APIcurrent   = "http://api.weatherapi.com/v1/current.json"
	APIforecast  = "http://api.weatherapi.com/v1/forecast.json"
	AksaiLatLon  = "51.17, 53.02"
	UralskLatLon = "51.22, 51.38"
)

var cities = map[string]string{
	"aksai":  AksaiLatLon,
	"uralsk": UralskLatLon,
}

var geos = map[string]string{
	AksaiLatLon:  "Aksai",
	UralskLatLon: "Uralsk",
}
