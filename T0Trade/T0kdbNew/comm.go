package kdbtdx

import (
	"time"
)

const (
	EntrustTab = "request"
	//RequestTab    = "requestv3"
	//CancelTab     = "cancelTab"
	HeartBeat   = "HeartBeat"
	ResponseTab = "response"
	//ResponseTabV3 = "responsev3"
	FuncUpdate    = "upd0"
	PositionTab   = "position"
	AssetTab      = "asset"
	PubAndUpdFunc = "upd2"
)

var (
	entrustCols = []string{"sym", "qid", "accountname", "time", "entrustno", "stockcode", "askprice", "askvol",
		"bidprice", "bidvol", "withdraw", "status", "note", "reqtype", "params"}
	assetCoks = []string{"account", "available", "stockvalue", "assureasset", "totaldebit", "balance",
		"enablebailbalance", "perassurescalevalue","sym"}
	positionCoks = []string{"account", "zqdm", "zqmc", "cbj", "gpye", "fdyk", "ykbl", "sj", "gpsz", "kyye","sym"}
	Tb           = &tradeKdb{}
)

type Basic struct {
	Sym, Qid string
}

type entrust struct {
	Basic
	Accountname              string
	Time                     time.Time
	Entrustno                int32
	Stockcode                string
	Askprice                 float64
	Askvol                   int32
	Bidprice                 float64
	Bidvol, Withdraw, Status int32
	Note                     string
	Reqtype                  int32
	Params                   string
}

type Request struct {
	Basic
	ReqType                                             int32
	Time                                                time.Time
	ParentId, Trader, Fund, Strategy, Broker, Stockcode string
	Stocktype, Side                                     int32
	Askprice                                            float64
	Ordertype, Orderqty                                 int32
	Algorithm, Params                                   string
}

type CancelReq struct {
	Basic
	Entrustno int32
}

type Order struct {
	Request
	OrderId           string
	SecurityId        string
	EntrustNo, CumQty int32
	AvgPx             float64
	Withdraw, Status  int32
	Note              string
}

type Cfg struct {
	DbPath string   `json:"dbPath"`
	Host   string   `json:"host"`
	Port   int      `json:"port"`
	Auth   string   `json:"auth"`
	Sym    []string `json:"sym"`
	MaxId  int32    `json:"maxId"`
}

type Hold struct {
	Account    string `json:"account"`
	Ticker     string `json:"ticker"`
	Count      int    `json:"count"`
	FrozenQty  int    `json:"frozenQty"`
	InitialQty int    `json:"initialQty"`
}

type queryResult struct {
	Entrustno, Status, Cumqty int32
}

type Asset struct {
	Account             string  `json:"account"`
	Available           float64 `json:"kyzj"`
	StockValue          float64 `json:"-"`    //股票总市值
	AssureAsset         float64 `json:"jzc"`  //净资产
	TotalDebit          float64 `json:"fzje"` //总负债
	Balance             float64 `json:"zzc"`  //总资产
	EnableBailBalance   float64 `json:"-"`    //可用保证金
	PerAssurescaleValue float64 `json:"-"`    //维持担保比例
	Sym                 string  `json:"sym"`
}

//`account`zqdm`zqmc`cbj`gpye`fdyk`ykbl`sj`gpsz`kyye
type Position struct {
	Account   string  `json:"-"`
	Stockcode string  `json:"zqdm"`
	StockName string  `json:"zqmc"`
	Cost      float64 `json:"cbj"`  //成本价
	Gpye      int32   `json:"gpye"` //股票余额
	Fdyk      float64 `json:"fdyk"` //浮动盈亏
	Ykbl      float64 `json:"ykbl"` //盈亏比例
	Sj        float64 `json:"sj"`   //市价
	Gpsz      float64 `json:"gpsz"` //市值
	Kyye      int32   `json:"kyye"` //可用
	Sym       string  `json:"sym"`
}
