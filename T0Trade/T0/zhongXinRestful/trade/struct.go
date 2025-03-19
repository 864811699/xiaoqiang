package trade

type RespBase struct {
	Code    int    `json:"code"` //0-成功  1-失败
	Message string `json:"message"`
	Data    Data   `json:"data"`

}

type LoginReqParams struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type Data struct {
	Token   string `json:"token"`
	OrderId string    `json:"orderId"`
}

type OrderReq struct {
	AcctId         string  `json:"acctId"`
	ClOrdId        string  `json:"clOrdId"`
	ExchId         string  `json:"exchId"`
	StkId          string  `json:"stkId"`
	OrderQty       int     `json:"orderQty"`
	OrderPrice     float64 `json:"orderPrice"`
	OrderSide      string  `json:"orderSide"`
	EntrustType    string  `json:"entrustType"`
	CashType       string  `json:"cashType"`
	InvestStrategy string  `json:"investStrategy"`
	InvestMemo     string  `json:"investMemo"`
}

type CancelReq struct {
	AcctId  string `json:"acctId"`
	ClOrdId string `json:"clOrdId"`
	OrderId string `json:"orderId"`
}


type QueryRsp struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data []QueryOrderRsp `json:"data"`
}

type QueryOrderRsp struct {
	ClOrdId        string  `json:"clOrdId"`
	BasketId       string  `json:"basketId"`
	OrderId        string  `json:"orderId"`
	AcctId         string  `json:"acctId"`
	ExchId         string  `json:"exchId"`
	StkId          string  `json:"stkId"`
	OrderSide      string  `json:"orderSide"`
	OrderQty       int     `json:"orderQty"`
	OrderPrice     float64 `json:"orderPrice"`
	OrderTime      string  `json:"orderTime"`
	ErrorMsg       string  `json:"errorMsg"`
	KnockQty       int     `json:"knockQty"`  //成交数量
	KnockPrice     float64 `json:"knockPrice"`  //成交均价
	KnockAmt       float64 `json:"knockAmt"`
	OrderStatus    string  `json:"orderStatus"`
	CanceledQty    int     `json:"canceledQty"`
	CashType       string  `json:"cashType"`
	EntrustType    string  `json:"entrustType"`
	UpdateTime     string  `json:"updateTime"`
	InvestStrategy string  `json:"investStrategy"`
	InvestMemo     string  `json:"investMemo"`
}

type OrderRspStream struct {
	DataType string `json:"dataType"` //Order/Trade
	Data QueryOrderRsp `json:"data"`
}
