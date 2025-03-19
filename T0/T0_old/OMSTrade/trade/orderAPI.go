package trade

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type OrderRequest struct {
	OriginId string `json:"originId"`
	//%.2f
	InputPrice  float64 `json:"inputPrice"`
	InputVol    int     `json:"inputVol"`
	OrderAction int     `json:"orderAction"` //4=BUY_TO_OPEN;5=SELL_TO_CLOSE<br/>8=SELL_TO_OPEN;1=BUY_TO_CLOSE
	Remark      string  `json:"remark"`
	Broker      string  `json:"broker"`
	Symbol      string  `json:"symbol"`
	AccountId   string  `json:"accountId"`
	User        string  `json:"user"`
	Group       string  `json:"group"`
	Channel     int     `json:"channel"` //1-conn 2-qfii
	Portfolio   string  `json:"portfolio"`
	StrategyId  string  `json:"strategyId"`
}

type CancelRequest struct {
	OriginIds []string `json:"originIds"`
}

//"errorId":1,"errorMsg"
type Response struct {
	ErrorId  int    `json:"errorId"`
	ErrorMsg string `json:"errorMsg"`
}

type QueryRsp struct {
	ErrorId  int    `json:"ErrorId"`
	ErrorMsg string `json:"ErrorMsg"`
	Data
}

type Data struct {
	ParentOrders []ParentOrder `json:"Data"`
}

type QueryRequest struct {
	Portfolios []string `json:"Portfolios"`
}

func newRequest(urlStr string, bts []byte) ([]byte, error) {
	rsp, err := http.Post(urlStr, "application/json", bytes.NewReader(bts))
	if err != nil {
		return nil, fmt.Errorf("postfail,%v", err)
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("io read fail,%v", err)
	}
	return b, nil

}

func NewOrder(urlStr string, request OrderRequest) (*Response, error) {
	bts, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("marshal fail,%v", err)
	}
	resp, err := newRequest(urlStr, bts)
	if err != nil {
		return nil, fmt.Errorf("new order fail,%v", err)
	}

	r := &Response{}
	fmt.Printf("%s\n", resp)
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return nil, fmt.Errorf("unmarshal fail,%v", err)
	}
	return r, err
}

func NewCancel(urlStr string, request CancelRequest) ([]*Response, error) {
	bts, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("marshal fail,%v", err)
	}
	bts, err = newRequest(urlStr, bts)
	if err != nil {
		return nil, fmt.Errorf("cancel fail,%v", err)
	}
	fmt.Printf("%s\n", bts)
	r := []*Response{}
	err = json.Unmarshal(bts, &r)
	if err != nil {
		return nil, fmt.Errorf("unmarshal fail,%v", err)
	}
	return r, err
}

func SearchOrder(urlStr string, request QueryRequest) (*QueryRsp, error) {
	bts, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf(" marshal request fail,%v", err)
	}

	b, err := newRequest(urlStr, bts)
	if err != nil {
		return nil, fmt.Errorf("newRequest url[%s],bts[%s] fail,%v",urlStr,bts, err)
	}
	r := QueryRsp{}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return nil, fmt.Errorf("unmarshal rsp[%s] fail,%v",b, err)
	}

	return &r, err
}

var StatusMapInt32 = map[OrderStatus_Enum]int32{
	OrderStatus_UNKNOW:         1,
	OrderStatus_CREATE:         1,
	OrderStatus_LOCAL:          1,
	OrderStatus_STARTED:        1,
	OrderStatus_TRIGGERED:      1,
	OrderStatus_PLACED:         1,
	OrderStatus_PARTIAL:        2,
	OrderStatus_FILLED:         4,
	OrderStatus_CANCELED:       5,
	OrderStatus_ERROR:          6,
	OrderStatus_EXPIRED:        6,
	OrderStatus_PENDING_CANCEL: 3,
}
