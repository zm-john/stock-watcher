# Stock-watcher
背景：工作中一直使用 PHP 脚本语言做开发想选择一门静态语言作为一种补充，经过考虑选择了 go；
通过项目实践学习一门语言是最高效的。最近刚入坑股票，又没时间去时刻关注，所以用 go 实现了这个简易的股票价格通知程序。

## Todo
* 设定指定股票最低目标价和最高目标价进行通知，比如股票 10元，设定 9.5 和 10.5 那么会在小于等于 9.5 或大于等于 10.5 元左右进行通知
* 股票异常波动（5分中 ± 2%）提醒
* 后期更多功能....

## Usage

1. 配置 config.json 文件
```
{
  "stocks": [
    "sh600000", "sz000001"
  ],
  "time": {
    "start": "09:30",
    "end": "15:00"
  },
  "notification": {
    "url": "https://hook.bearychat.com/=bwCEN/incoming/88b7af379613e62f292f487e4c08d42e",
    "channel": "stock"
  },
  "interval": 30
}

stocks：指定股票代码，前缀 sh 表示上交所股票, 前缀 sz 表示深交所股票， 股票代码以 6 开头的上交所股票, 0、3 开头的为深交所股票

time：交易时间
|---- start 开盘时间
|____ end 收盘时间

notification：通知
|---- url 通知地址，我使用的 bearychat 会向这个地址 Post content-type:application/json 数据{"text": "股票数据", "channel": "指定股票接收频道"}
|____ channel 频道

interval：通知间隔时间，单位秒

```


2. 运行
```
1）克隆项目
git clone git@github.com:zm-john/stock-watcher.git // 最好放到 ~/go/src 目录下面，不然需要手动设置 GOPATH 为当前目录

2）进入项目目录
cd stock-watcher

3）安装依赖包
./install.sh or glide install

4）运行
go run main.go

```

## License
MIT