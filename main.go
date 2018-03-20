package main

import (
	"io/ioutil"
	"encoding/json"
	"stock-watcher/stock"
	"strings"
)

func main() {
	config := loadConfig()
	aliasStockCode(&config)
	watcher := stock.Watcher{}
	watcher.Config(config)
	watcher.Watch()
}

func loadConfig() stock.Config {
	config := stock.Config{}
	data, err := ioutil.ReadFile("./config.json")

	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &config); err != nil {
		panic(err)
	}

	return config
}

func aliasStockCode(config *stock.Config) {
	var code string
	for index, _ := range config.Stocks {
		code = config.Stocks[index].Code
		if strings.HasPrefix(code, "6") {
			config.Stocks[index].Alias = strings.Join([]string{"sh", code}, "")
		} else {
			config.Stocks[index].Alias = strings.Join([]string{"sz", code}, "")
		}
	}
}