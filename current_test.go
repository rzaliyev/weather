package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestCurrentWeather(t *testing.T) {

	t.Run("test current weather", func(t *testing.T) {
		data, err := ioutil.ReadFile("current_sample.json")
		if err != nil {
			t.Fatal(err)
		}

		assertJSONUnknownFields(t, data)
	})

	t.Run("test forecast weather", func(t *testing.T) {
		data, err := ioutil.ReadFile("forecast_sample.json")
		if err != nil {
			t.Fatal(err)
		}

		assertJSONUnknownFields(t, data)
	})
}

func assertJSONUnknownFields(t testing.TB, data []byte) {
	t.Helper()

	var resp Response
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&resp); err != nil {
		t.Error(err)
	}
}
