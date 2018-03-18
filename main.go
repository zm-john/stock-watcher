package main

import (
	"io/ioutil"
	"encoding/json"
	"stock-watcher/stock"
)

var config = stock.Config{}

func main() {
	loadConfig()
	watcher := stock.Watcher{}
	watcher.Config(config)
	watcher.Watch()
}

func loadConfig() {
	data, err := ioutil.ReadFile("./config.json")

	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &config); err != nil {
		panic(err)
	}
}