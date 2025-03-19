package trade

import kdbtdx "github.com/864811699/T0kdb"

type RequestType string

const (
	LOGIN               RequestType = "auth"
	QUERRY_ACCOUNT_INFO RequestType = "actinfo"
	BUY                 RequestType = "wt_mairu"
	SELL                RequestType = "wt_maichu"
	CANCELOrder         RequestType = "wt_chedan"
	QUERY_FUND          RequestType = "wt_query_zijin"
	QUERY_HOLDING       RequestType = "wt_query_gupiao"
	QUERY_CANCEL        RequestType = "wt_query_chedan"
	QUERY_ORDERS        RequestType = "wt_query_weituo"
	QUERY_DEAL          RequestType = "wt_query_chengjiao"
)

type SideType string

const (
	AUTO SideType = "auto"
	//融资买入
	RONGZI_BUY SideType = "xy"
	//担保品买入
	DANBAOPING_BUY SideType = "zy"
	//卖券还款
	MAIQUANHUANKUAN_SELL SideType = "xy"
	//担保品卖出
	DANBAOPING_SELL SideType = "zy"
	//融券卖出
	RONGQUAN_SELL SideType = "rq"
)

type SignalRsp struct {
	Ret int `json:"ret"`
}

type Request struct {
	Cmd string `json:"cmd"`
}

type Ps struct {
	User string `json:"user"`
	Pwd  string `json:"pwd"`
}

type LoginReq struct {
	Request
	Ps Ps `json:"ps"`
}

type LoginRsp struct {
	SignalRsp
	Nickname string `json:"nickname"`
}

type ActInfo struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

//{"ret":0,"query":[{"account":"310210300347","qsname":"东方财富X(信用)","name":"testdc0215","id":8571}]}
type AccountInfoRsp struct {
	SignalRsp
	AccountInfos []*ActInfo `json:"query"`
}

type OrderReq struct {
	Request
	Stockcode string   `json:"stockcode"`
	Price     float64  `json:"price"`
	ID        int      `json:"id"`
	Amount    int      `json:"amount"`
	WtFlag    SideType `json:"wt_flag"`
}

type OrderRsp struct {
	SignalRsp
	ErrMsg  string `json:"err"`
	OrderId string `json:"htbh"`
}

type CancelReq struct {
	Request
	ID        int    `json:"id"`
	OrderId   string `json:"htbh"`
	Stockcode string `json:"stockcode"`
}

type QueryReq struct {
	Request
	ID int `json:"id"`
}

type QueryRsp struct {
	Ret       int         `json:"ret"`
	Responses []*Response `json:"query"`
}

type Response struct {
	OpName string `json:"op_name"`
	//Wtlx          string  `json:"wtlx"`
	//Date          string  `json:"wtrq"` //委托日期 20230227
	Note string `json:"note"` //备注信息
	//Uuid          int     `json:"uuid"`
	AskPrice float64 `json:"wtjg"` //委托价格
	Side     string  `json:"op"`   //B-buy   S-sell
	BidPrice float64 `json:"cjjg"` //成交价格
	//Tact          string  `json:"tact"`
	//ReqTime       string  `json:"wtsj"` //委托时间 "14:41:22"
	Status int `json:"statu` //0-未成交 1-全部成交 2-部分成交 3-部成部撤  4-全部撤单
	//Cdsl          int     `json:"cdsl"` //撤单数量
	OrderId string `json:"htbh"` //合同编号
	//Jys           string  `json:"jys":`
	//Gdzh          string  `json:"gdzh"` //股东账号
	BidQty        int     `json:"cjsl"` //成交数量
	AskTotalPrice float64 `json:"cjje"` //成交总金额
	Stockcode     string  `json:"zqdm"` //证券代码
	//StockName     string  `json:"zqmc"` //证券名称
	AskQty int `json:"wtsl"` //委托数量
	//Qsname        string  `json:"qsname"`
	//Qsid          int     `json:"qsid"`
}



var OrderStatusMap = map[int]int32{
	0: 1,
	1: 4,
	2: 2,
	3: 5,
	4: 5,
}

type QueryAssetRsp struct {
	Ret       int         `json:"ret"`
	Responses []kdbtdx.Asset `json:"query"`
}
type QueryPositionRsp struct {
	Ret       int         `json:"ret"`
	Responses []*kdbtdx.Position `json:"query"`
}