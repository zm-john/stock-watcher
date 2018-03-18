package stock

import (
	"net/http"
	"io/ioutil"
	"golang.org/x/text/transform"
	"golang.org/x/text/encoding/simplifiedchinese"
	"regexp"
	"fmt"
	"strings"
	"github.com/uniplaces/carbon"
	"time"
	"strconv"
	"encoding/json"
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
	for {
		if w.isTradeTime() == false {
			continue
		}

		query := w.stocks()
		rest, err := fetch(query)

		if err != nil {
			continue
		}

		message := format(rest)
		w.notify(message)

		time.Sleep(30 * time.Second)
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
	fmt.Println(string(buf))
	http.Post(w.config.Notify.Url, "application/json", strings.NewReader(string(buf)))
}

func (w *Watcher) stocks() string {
	return fmt.Sprintf("%slist=%s", HOST, strings.Join(w.config.Stocks, ","))
}

func fetch(url string) (string, error) {
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


func format(str string) string {
	piece := match(str)
	data := make(map[string]string)
	substr := make([]string, len(piece))
	for index, val := range piece {
		stock := string(val[1])
		p := strings.Split(stock, ",")

		for index, _ := range stockField {
			data[stockField[index]] = p[index]
		}


		substr[index] = fmt.Sprintf(
			`## %s
当前价：%s / 开盘价：%s / 最高价：%s / 最低价：%s
卖一 %s：%s / 卖二 %s：%s / 卖三 %s：%s / 卖四 %s：%s / 卖五 %s：%s
买一 %s：%s / 买二 %s：%s / 买三 %s：%s / 买四 %s：%s / 买五 %s：%s
`,
				data["name"], data["currentPrice"],
				data["openingPrice"],
				data["highestPrice"],
				data["lowestPrice"],
				data["sale1Price"],
				data["sale1Quantity"],
				data["sale2Price"],
				data["sale2Quantity"],
				data["sale3Price"],
				data["sale3Quantity"],
				data["sale4Price"],
				data["sale4Quantity"],
				data["sale5Price"],
				data["sale5Quantity"],
				data["buy1Price"],
				data["buy1Quantity"],
				data["buy2Price"],
				data["buy2Quantity"],
				data["buy3Price"],
				data["buy3Quantity"],
				data["buy4Price"],
				data["buy4Quantity"],
				data["buy5Price"],
				data["buy5Quantity"])
	}

	return strings.Join(substr, "\r\n\r\n")
}


func match(str string) [][][]byte {
	reg := regexp.MustCompile(`var\W\w+=\"(.+)\"`)
	rest := reg.FindAllSubmatch([]byte(str), -1)

	return rest
}