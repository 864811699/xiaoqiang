# T0 交易接口文档

测试环境: 139.224.56.241

生产环境: 47.103.115.44

## 1: 下单
- [post] : `http://ip:6026/manual/insertOrder`
### 1.1 request 


| 字段       | 描述       | 数据类型   |备注|
|----------|----------|--------|------|
| originId | 订单号,UUID | string | [必填，唯一] |
| inputPrice | 下单价格 | %.2f   | [必填] >0|
| inputVol | 下单数量 | int    | [必填] >0|
| orderAction | 订单方向 | int    | [必填]<br/>4=BUY_TO_OPEN;5=SELL_TO_CLOSE<br/>8=SELL_TO_OPEN;1=BUY_TO_CLOSE |
| broker | 券商| string | [必填] |
| remark | 备用字段| string | [必填] |
| symbol | 股票代码| string | [必填] |
| accountId | 账户name| string | [必填] |
| channel | 渠道| int | [必填]  1-conn 2-qfii|
| trader | 交易员 | string | [可选]  |


```json
{
      "originId": "effe58d2-5b0a-4a84-a65e-d5bf4a79f9f3",
      "inputPrice": 10.50,
      "inputVol": 100,  
      "orderAction": 4, 
      "broker": "anxin",
      "remark": "测试",
      "symbol": "601818", 
      "accountId": "msrc8yh8888"
}
```




### 1.2 response 

```json
[{
  "errorId": -1,  //0成功
  "errorMsg": "xxxxxxx",
  }]

```


## 2: 撤单
- [post] : `http://ip:6026/manual/cancelOrder`

### 2.1 request

| 字段       | 描述       | 数据类型   |备注|
|----------|----------|--------|------|
| originId | 订单号,UUID | string | [必填，唯一] |

```json
{
  "originIds":["effe58d2-5b0a-4a84-a65e-d5bf4a79f9f3"]
}
```


### 2.2 response

```json
[{
    "errorId":-6,   //0成功
    "errorData": "xxxxxx"
}]
```

## 3: 母单查询
- [post] : `http://ip:6026/order/getParentOrders`

### 2.1 request

```
{
	"Portfolios": ["intraday_strategy_2"]
}
```

### 3.2 response


| 字段       | 描述       |
|----------|---------|
|originId|订单号|
|orderStatus|订单状态|
|insertTime|创建时间|
|filledQuantity|成交数量|
|inputPrice|下单价格|
|inputVol|下单数量|
|orderType|订单类型|
|orderAction|订单方向|
|broker|券商|
|latestUpdateTime|更新时间|
|symbol|合约|
|accoutId|账户|
|strategyId|策略|
|avgPrice|成交均价|
|portfolio|组合|
|algo|算法名|
|group||
|channel|通道|
|runingStatus|运行状态|
|params|参数列表|
|signalId|信号id|
|tradingBatchId|批次号|

```json
{
    "ErrorId": 0,
    "ErrorMsg": "",
    "Data": [
            {
                "originId": "effe58d2-5b0a-4a84-a65e-d5bf4a79f2f9",
                "orderStatus": 11,
                "insertTime": 1665376178714297520,
                "filledQuantity": 0,
                "inputPrice": 2.84,
                "inputVol": 100,
                "orderType": 1,
                "orderAction": 4,
                "broker": "anxin",
                "latestUpdateTime": 1665377895755783440,
                "errorID": 0,
                "errorMsg": "",
                "hedgeFlag": 0,
                "remark": "BrokerAlgoAndDMA",
                "symbol": "601818",
                "accoutId": "msrc8yh8888",
                "strategyId": "manual T0",
                "user": "",
                "avgPrice": 0,
                "portfolio": "portfolio_manualT0",
                "algo": "DMA",
                "group": "",
                "trader": "",
                "channel": 3,
                "runingStatus": 3,
                "params": "",
                "pm": "",
                "signalId": "",
                "tradingBatchId": ""
            },
            {
                "originId": "dffe58d2-5b0a-4a84-a65e-d5bf4a79f2f9",
                "orderStatus": 11,
                "insertTime": 1665377928594541117,
                "filledQuantity": 0,
                "inputPrice": 2.84,
                "inputVol": 100,
                "orderType": 1,
                "orderAction": 4,
                "broker": "anxin",
                "latestUpdateTime": 1665377945146262870,
                "errorID": 0,
                "errorMsg": "",
                "hedgeFlag": 0,
                "remark": "BrokerAlgoAndDMA",
                "symbol": "601818",
                "accoutId": "msrc8yh8888",
                "strategyId": "manual T0",
                "user": "",
                "avgPrice": 0,
                "portfolio": "portfolio_manualT0",
                "algo": "DMA",
                "group": "",
                "trader": "",
                "channel": 3,
                "runingStatus": 3,
                "params": "",
                "pm": "",
                "signalId": "",
                "tradingBatchId": ""
            }
      ]
}
```

## 4: 母单订阅
rabbitMq 订阅方式参考代码样例
