package trade

import (
	"bytes"
	"encoding/json"
	"fmt"
	logger "github.com/alecthomas/log4go"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type UrlConfig struct {
	LoginUrl      string `json:"login_url"`
	HeartUrl      string `json:"heart_url"`
	NewOrderUrl   string `json:"new_order_url"`
	NewCancelUrl  string `json:"new_cancel_url"`
	QueryOrderUrl string `json:"query_order_url"`
	SubUrl        string `json:"sub_url"`
}

type RestfulApi struct {
	url   *UrlConfig
	Token string
	RspCh chan QueryOrderRsp
	IsOk  bool
}

func NewRestfulApi(url *UrlConfig) *RestfulApi {
	return &RestfulApi{
		url:   url,
		RspCh: make(chan QueryOrderRsp,100000),
	}
}

func (api *RestfulApi) Login(userName, password string) error {
	loginReq := LoginReqParams{
		UserName: userName,
		Password: password,
	}
	bts, _ := json.Marshal(loginReq)
	rsp, err := http.Post(api.url.LoginUrl, "application/json", bytes.NewReader(bts))
	if err != nil {
		return errors.Wrapf(err, "login[%s] request fail", userName)
	}
	defer rsp.Body.Close()
	bts, _ = ioutil.ReadAll(rsp.Body)
	logger.Info("login resp:%s",bts)
	loginRsp := RespBase{}
	if err := json.Unmarshal(bts, &loginRsp); err != nil {
		return errors.Wrapf(err, "login[%s] response unmarshal fail,msg:%s", userName, bts)
	}
	if loginRsp.Code == 0 {
		api.Token = loginRsp.Data.Token
		api.IsOk = true
		api.url.SubUrl = fmt.Sprintf("%s/%s", api.url.SubUrl, loginRsp.Data.Token)
		return nil
	}
	return fmt.Errorf("login fail ,msg:%s", loginRsp.Message)
}

func (api *RestfulApi) setHeader(r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Token", api.Token)
}

func (api *RestfulApi) HeartBeat() error {
	r, _ := http.NewRequest("GET", api.url.HeartUrl, strings.NewReader(""))
	api.setHeader(r)
	rsp, err := (&http.Client{}).Do(r)
	if err != nil {
		api.IsOk = false
		return errors.Wrapf(err, "heartBeat do fail")
	}
	defer rsp.Body.Close()
	bts, _ := ioutil.ReadAll(rsp.Body)
	heartBeatRsp := RespBase{}
	if err := json.Unmarshal(bts, &heartBeatRsp); err != nil {
		api.IsOk = false
		return errors.Wrapf(err, "heartBeat response[%s] unmarshal fail", bts)
	}
	if heartBeatRsp.Code != 0 {
		api.IsOk = false
		return fmt.Errorf("heartBeat check fail,msg:%s", heartBeatRsp.Message)
	}
	return nil
}

func (api *RestfulApi) NewOrder(o OrderReq) (string, error) {
	bts, _ := json.Marshal(o)
	logger.Info("send new order request,%s",bts)
	r, _ := http.NewRequest("POST", api.url.NewOrderUrl, bytes.NewReader(bts))
	api.setHeader(r)
	rsp, err := (&http.Client{}).Do(r)
	if err != nil {
		return "", errors.Wrapf(err, "[%s] send order[%s]  fail,", o.AcctId, o.ClOrdId)
	}

	defer rsp.Body.Close()
	orderRsp := RespBase{}
	bts, _ = ioutil.ReadAll(rsp.Body)
	logger.Info("new order resp:%s",bts)
	//{"data":{"orderId":"1000005"},"code":0,"message":"成功"}
	if err := json.Unmarshal(bts, &orderRsp); err != nil {
		return "", errors.Wrapf(err, "account[%s] new order[%s] response[%s] unmarshal fail", o.AcctId, o.ClOrdId, bts)
	}
	if orderRsp.Code != 0 {
		return "", fmt.Errorf("account[%s] new order[%s] is rejected by broker,msg:%s", o.AcctId, o.ClOrdId, orderRsp.Message)
	}
	return orderRsp.Data.OrderId, nil
}

func (api *RestfulApi) NewCancel(c CancelReq) error {
	bts, _ := json.Marshal(c)
	logger.Info("send cancel request,%s",bts)
	r, _ := http.NewRequest("POST", api.url.NewCancelUrl, bytes.NewReader(bts))
	api.setHeader(r)
	rsp, err := (&http.Client{}).Do(r)
	if err != nil {
		return errors.Wrapf(err, "[%s] send cancel[%s]  fail,", c.AcctId, c.ClOrdId)
	}
	defer rsp.Body.Close()
	cancelRsp := RespBase{}
	bts, _ = ioutil.ReadAll(rsp.Body)
	logger.Info("new cancel resp:%s",bts)
	if err := json.Unmarshal(bts, &cancelRsp); err != nil {
		return errors.Wrapf(err, "account[%s] new order[%s] response[%s] unmarshal fail", c.AcctId, c.ClOrdId, bts)
	}
	if cancelRsp.Code != 0 {
		return fmt.Errorf("account[%s] new order[%s] is rejected by broker,msg:%s", c.AcctId, c.ClOrdId, cancelRsp.Message)
	}
	return nil
}

func (api *RestfulApi) NewQueryOrder(account string) ([]QueryOrderRsp, error) {
	url:=fmt.Sprintf("%s?acctId=%s",api.url.QueryOrderUrl,account)
	r, _ := http.NewRequest("GET", url, strings.NewReader(""))
	api.setHeader(r)
	rsp, err := (&http.Client{}).Do(r)
	if err != nil {
		api.IsOk = false
		return nil, errors.Wrapf(err, "query [%s] orders  fail", account)
	}
	defer rsp.Body.Close()
	bts, _ := ioutil.ReadAll(rsp.Body)
	logger.Info("query order resp:%s",bts)
	queryRsp := QueryRsp{}
	if err := json.Unmarshal(bts, &queryRsp); err != nil {
		api.IsOk = false
		return nil, errors.Wrapf(err, "query [%s] orders response[%s] unmarshal fail", account, bts)
	}
	if queryRsp.Code != 0 {
		api.IsOk = false
		return nil, fmt.Errorf("query [%s] orders fail,msg:%s", account, queryRsp.Message)
	}
	return queryRsp.Data, nil
}

func (api *RestfulApi) SubOrder() error {
	c, _, err := websocket.DefaultDialer.Dial(api.url.SubUrl, nil)
	if err != nil {
		return errors.Wrapf(err, "subOrder[%s] fail", api.url.SubUrl)
	}
	defer c.Close()

	for {
		_, msg, err := c.ReadMessage()
		if err == io.EOF {
			logger.Warn("接收订阅结束")
			return nil
		}
		if err != nil {
			return errors.Wrapf(err, "接收数据异常")
		}
		logger.Info("sub order resp:%s",msg)
		orderRsp := OrderRspStream{}
		if err := json.Unmarshal(msg, &orderRsp); err != nil {
			logger.Warn("sub order[%s] unmarshal fail,err:%v",msg,err)
			continue
		}

		if orderRsp.DataType=="Order"{
			api.RspCh<-orderRsp.Data
		}

	}

}
