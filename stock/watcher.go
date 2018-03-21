package stock

import (
    "net/http"
    "io/ioutil"
    "golang.org/x/text/transform"
    "golang.org/x/text/encoding/simplifiedchinese"
    "regexp"
    "strings"
    "github.com/uniplaces/carbon"
    "time"
    "strconv"
    "encoding/json"
    "sync"
    "errors"
    "fmt"
)

const HOST = "http://hq.sinajs.cn/"

var stockField []string = []string{
    "name",
    "openingPrice",     // 开盘价
    "closingPrice",     // 收盘价
    "currentPrice",     // 当前价
    "highestPrice",     // 最高价
    "lowestPrice",      // 最低价
    "biddingBuyPrice",  // 竞买价
    "biddingSalePrice", // 竞卖价
    "quantity",         // 成交量
    "turnover",         // 成交额
    "buy1Quantity",     // 买1申请数量
    "buy1Price",        // 买1报价
    "buy2Quantity",     // 买2申请数量
    "buy2Price",        // 买2报价
    "buy3Quantity",     // 买3申请数量
    "buy3Price",        // 买3报价
    "buy4Quantity",     // 买4申请数量
    "buy4Price",        // 买4报价
    "buy5Quantity",     // 买5申请数量
    "buy5Price",        // 买5报价
    "sale1Quantity",    // 卖1申请数量
    "sale1Price",       // 卖1报价
    "sale2Quantity",    // 卖2申请数量
    "sale2Price",       // 卖2报价
    "sale3Quantity",    // 卖3申请数量
    "sale3Price",       // 卖3报价
    "sale4Quantity",    // 卖4申请数量
    "sale4Price",       // 卖4报价
    "sale5Quantity",    // 卖5申请数量
    "sale5Price",       // 卖5报价
    "date",             // 日期
    "time",             // 时间
    "end",              //
}

type Watcher struct{
    config Config
}

type Message struct {
    Text string `json:"text"`
    Channel string `json:"channel"`
} 

func (w *Watcher) Config(conf Config) {
    w.config = conf
}

func (w *Watcher) Watch() {
    var wg sync.WaitGroup
    for {
        if w.isTradeTime() == false {
            continue
        }

        for _, stock := range w.config.Stocks {
            wg.Add(1)
            go func(st Stock) {
                defer wg.Done()
                rest, err := fetch(st.Alias)
                if err != nil {
                    w.notify(err.Error())
                }
                message := format(*w, rest)
                if len(message) > 0 {
                    w.notify(message)
                }
            }(stock)
        }


        wg.Wait()

        time.Sleep(w.config.Interval * time.Second)
    }
}




func (w *Watcher) isTradeTime() bool {
    if carbon.Now().IsWeekend() {
        // weekend
        return false
    }

    start := strings.Split(w.config.Time.Start, ":")
    end := strings.Split(w.config.Time.End, ":")

    if len(start) != 2 {
        panic("start time config err.")
    }
    if len(end) != 2 {
        panic("end time config err.")
    }

    startTime := carbon.Now()
    startHour, err := strconv.Atoi(start[0])
    if err != nil {
        startHour = 0
    }
    startTime.SetHour(startHour)

    startMinute, err := strconv.Atoi(start[1])
    if err != nil {
        startMinute = 0
    }
    startTime.SetMinute(startMinute)

    endTime := carbon.Now()
    endHour, err := strconv.Atoi(end[0])
    if err != nil {
        endHour = 0
    }
    endTime.SetHour(endHour)

    endMinute, err := strconv.Atoi(end[1])
    if err != nil {
        endMinute = 0
    }
    endTime.SetMinute(endMinute)

    now := carbon.Now()

    if now.Gte(startTime) && now.Lte(endTime) {
        return true
    }

    return false
}

func (w *Watcher) notify(message string) {
    msg := Message{message, w.config.Notify.Channel}
    buf, err := json.Marshal(msg)

    if err != nil {
        panic("josn encode error")
    }
    http.Post(w.config.Notify.Url, "application/json", strings.NewReader(string(buf)))
}

func fetch(code string) (string, error) {
    url := strings.Join([]string{HOST, "list=", code}, "")
    response, err := http.Get(url)

    if err != nil {
        // error
        return "", err
    }

    reader := transform.NewReader(response.Body, simplifiedchinese.GBK.NewDecoder())
    body, err := ioutil.ReadAll(reader)

    if err != nil {
        // error
        return "", err
    }
    return string(body), nil
}


func format(w Watcher, str string) string {
    piece, err := match(str)
    if err != nil {
        return err.Error()
    }

    data := make(map[string]string)
    stockInfo := string(piece[2])
    p := strings.Split(stockInfo, ",")

    for index, _ := range stockField {
        data[stockField[index]] = p[index]
    }

    code := string(piece[1])
    stock, err := w.findStock(code)

    if err != nil {
        return err.Error()
    }
    price, err := strconv.ParseFloat(data["currentPrice"], 64)
    if err != nil {
        return err.Error()
    }
    if stock.Highest > 0 && price >= stock.Highest  {
        // 达到最高目标价
        return fmt.Sprintf("%s 达到最高目标价，现价 %s", data["name"], data["currentPrice"])

    }
    if stock.Lowest > 0 && price <= stock.Lowest {
        // 达到最低目标价
        return fmt.Sprintf("%s 达到最低目标价，现价 %s", data["name"], data["currentPrice"])
    }

    return ""
}


func (w *Watcher) findStock(code string) (Stock, error) {

    for _, val := range w.config.Stocks {
        if val.Code == code {
            return val, nil
        }
    }

    return Stock{}, errors.New("code not find")
}

func match(str string) ([][]byte, error) {
    reg := regexp.MustCompile(`var\W\w+[sz|sh](\d{6})=\"(.+)\"`)
    rest := reg.FindAllSubmatch([]byte(str), -1)

    if len(rest) == 0 {
        return [][]byte{}, errors.New("匹配失败")
    }
    return rest[0], nil
}