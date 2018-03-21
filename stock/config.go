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

type Stock struct {
    Code string `json:"code"`
    Highest float64 `json:"highest"`
    Lowest float64 `json:"lowest"`
    Alias string
}


type Config struct{
    Stocks []Stock `json:"stocks"`
    Time Time `json:"time"`
    Notify Notification `json:"notification"`
    Interval time.Duration `json:"interval"`
    StockQueryHost string `json:"stock_host"`
}
