# Stock-watcher
背景：工作中一直使用 PHP 脚本语言做开发想选择一门静态语言作为一种补充，经过考虑选择了 go；
通过项目实践学习一门语言是最高效的。最近刚入坑股票，又没时间去时刻关注，所以用 go 实现了这个简易的股票价格通知程序。

## Todo
* 股票异常波动（5分中 ± 2%）提醒
* 使用 Mysql
* 使用 Web 添加关注股票

## Usage

1. 复制配置文件 `cp config.json.example config.json`

2. 修改 config.json 文件
```
{
  "stocks": [
    {
      "code": "600536",
      "highest": 14.44,
      "lowest": 14.39
    },
    {
      "code": "000016",
      "highest": 6.77,
      "lowest": 6.75
     }
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

stocks：指定监控股票代码
|———— code 股票代码
|———— highest 最高目标价，如果为 0 则无忽略
|———— lowest 最低目标价，如果为 0 则无忽略
|____ remark 备注，自己看


time：交易时间
|———— start 开盘时间
|____ end 收盘时间

notification：通知
|———— url 通知地址，我使用的 bearychat 会向这个地址 Post content-type:application/json 数据{"text": "股票数据", "channel": "指定股票接收频道"}
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