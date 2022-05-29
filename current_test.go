package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestCurrentWeather(t *testing.T) {
	// read original
	data, err := ioutil.ReadFile("current_sample.json")
	if err != nil {
		t.Fatal(err)
	}

	var resp Response
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	if err = decoder.Decode(&resp); err != nil {
		t.Error(err)
	}
}
