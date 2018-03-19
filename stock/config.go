package stock

import "time"

type Time struct {
	Start string
	End string
}

type Notification struct{
	Url string
	Channel string
}

type Config struct{
	Stocks []string `json:"stocks"`
	Time Time `json:"time"`
	Notify Notification `json:"notification"`
	Interval time.Duration `json:"interval"`
	StockQueryHost string `json:"stock_host"`
}
