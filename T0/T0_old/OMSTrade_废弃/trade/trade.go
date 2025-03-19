package trade

import (
	"encoding/json"
	"fmt"
	kdbtdx "github.com/864811699/T0kdb"
	logger "github.com/alecthomas/log4go"
	"github.com/gin-gonic/gin"
	"net/http"

	"io/ioutil"
	"strings"
	"sync"
)

type Api struct {
	*TradeCfg
	respChan      chan *kdbtdx.Order
	TradeFile     string
	ch            chan *ParentOrder
	rbq           *rabbitMqManage
	orderMap      map[string]kdbtdx.Order
	orderMap_lock sync.RWMutex
	//cancelOrderMap      map[int32]context.CancelFunc
	//cancelOrderMap_lock sync.RWMutex
}

type TradeCfg struct {
	Host         string                  `json:"kdb_host"`
	Port         int                     `json:"kdb_port"`
	Auth         string                  `json:"auth"`
	DB           string                  `json:"db"`
	SubAccounts  []string                `json:"accounts"`
	Broker       string                  `json:"broker"`
	MaxNo        int32                   `json:"maxNo"`
	OrderUrl     string                  `json:"order_url"`
	CancelUrl    string                  `json:"cancel_url"`
	SearchUrl    string                  `json:"search_url"`
	SearchAlgo   string                  `json:"search_algo"`
	StrategyId string `json:"strategyId"`
	RbqUrl       string                  `json:"rbq_url"`
	RbqSubAlgo   []string                `json:"rbq_sub_algo"`
	AccountSides map[string]*AccountSide `json:"account_sides"`
	WebAddr	string `json:"webAddr"`
}

type AccountSide struct {
	Sell   OrderAction_Enum `json:"sell"`
	Buy    OrderAction_Enum `json:"buy"`
	Broker string           `json:"broker"`
}

func NewApi(configFile, tradeCfg string) *Api {
	return &Api{
		TradeFile: tradeCfg,

		TradeCfg: nil,
		respChan: make(chan *kdbtdx.Order, 10000),
		ch:       make(chan *ParentOrder, 10000),
		rbq:      NewRabbitMq(),
		orderMap: make(map[string]kdbtdx.Order, 10000),
		//cancelOrderMap: make(map[int32]context.CancelFunc),
	}
}

func (this *Api) LoadCfg() *kdbtdx.Cfg {
	readByte, err := ioutil.ReadFile(this.TradeFile)
	if err != nil {
		logger.Crashf("read trade cfg :%v fail ,err: %v", this.TradeFile, err)
	}
	tradeCfg := new(TradeCfg)

	err = json.Unmarshal(readByte, tradeCfg)
	if err != nil {
		logger.Crashf("tradecfg unmarshal fail,err:%v", err)
	}
	this.TradeCfg = tradeCfg
	//4=BUY_TO_OPEN;5=SELL_TO_CLOSE<br/>8=SELL_TO_OPEN;1=BUY_TO_CLOSE
	for _, accountInfo := range this.AccountSides {
		if  accountInfo.Sell != 8  && accountInfo.Buy != 1&& accountInfo.Buy != 4 &&accountInfo.Sell != 5{
			panic("买卖方向配置错误,8=SELL_TO_OPEN;1=BUY_TO_CLOSE")
		}
	}
	logger.Info("load cfg success!!")

	return &kdbtdx.Cfg{tradeCfg.DB, tradeCfg.Host, tradeCfg.Port, tradeCfg.Auth, tradeCfg.SubAccounts, tradeCfg.MaxNo}
}

func (this *Api) RunApi() {
	go func() {
		r:=gin.Default()
		r.GET("/getTrade",this.GetTrade)
		r.Run(this.WebAddr)
	}()
	if err := this.rbq.Init(this.RbqUrl); err != nil {
		logger.Crashf("rbq init fail,%v", err)
	}

	this.rbq.Sub(this, "", this.RbqSubAlgo)
	logger.Info("rbq sub success!")
	go this.update()

	this.updateAfterReboot()

	//go Sub(this.SubUrl, this.ch,this.ctx)
	logger.Info("API run success!")
}

func (this *Api) GetUpdatedInfo() chan *kdbtdx.Order {

	return this.respChan
}

func (this *Api) Trade(o kdbtdx.Order) {

	this.orderMap_lock.Lock()
	this.orderMap[o.Qid] = o
	this.orderMap_lock.Unlock()

	remark := Remark{}
	if err := json.Unmarshal([]byte(o.Note), &remark); err != nil {
		errmsg := fmt.Sprintf("json unmarshal fail,note:%s, err:%v", o.Note, err)
		logger.Warn(errmsg)
		o.Status = 6
		remark.Msg=errmsg
		bts,_:=json.Marshal(remark)
		o.Note = string(bts)
		return
	}else {
		//4=BUY_TO_OPEN;5=SELL_TO_CLOSE
		//8=SELL_TO_OPEN;1=BUY_TO_CLOSE
		if remark.Request.Side!=4&&remark.Request.Side!=5&&remark.Request.Side!=8&&remark.Request.Side!=1{
			errmsg := fmt.Sprintf("side[%d] error,side should be 1/4/5/8",remark.Request.Side)
			logger.Warn(errmsg)
			o.Status = 6
			remark.Msg=errmsg
			bts,_:=json.Marshal(remark)
			o.Note = string(bts)
			return
		}
	}

	//	OrderAction int     `json:"orderAction"` //bto4  stc5  sto8  btc1
	//side通过 remark.request.side识别
	kdbAccount, _, omsChannel := this.GetAccountSide(o.Sym, o.Side)
	if  omsChannel == -1 || kdbAccount == "" {
		err := fmt.Sprintf("trade fail ,webAccount[%s] or webSide[%d] error", o.Sym, o.Side)
		logger.Warn(err)
		o.Status = 6
		o.Note = err
	} else {
		o.Status = 0
		packId := pack(o.Sym, o.EntrustNo)
		req := OrderRequest{
			OriginId:    o.Qid,
			InputPrice:  o.Askprice,
			InputVol:    int(o.Orderqty),
			OrderAction: remark.Request.Side,
			Remark:      packId,
			Broker:      this.AccountSides[o.Sym].Broker,
			Symbol:      o.Stockcode,
			AccountId:   kdbAccount,
			Channel:     omsChannel,
			Group: o.Trader,
			Portfolio:this.SearchAlgo,
			StrategyId:this.StrategyId,

		}

		rsp, err := NewOrder(this.OrderUrl, req)
		if err != nil {
			o.Status = 6
			o.Note = err.Error()
			logger.Warn("order[%#v] send fail ,%v", req, err)
		} else if rsp.ErrorId != 0 {
			o.Status = 6
			o.Note = rsp.ErrorMsg
			logger.Warn("order[%#v] reject by OMS, %v", req, rsp.ErrorMsg)
		} else {
			logger.Info("trade order success :: entrustNo [%d] ,req:[%#v]", o.EntrustNo,req)
		}
	}

	kdbtdx.Store(&o)

}

func (this *Api) Cancel(c kdbtdx.CancelReq) {

	logger.Info("cancel entrustNo: %v", c.Entrustno)

	cancelRequest := CancelRequest{[]string{c.Qid}}
	rsp, err := NewCancel(this.CancelUrl, cancelRequest)
	logger.Warn("cancel send ,rsp:%#v,err:%v", rsp, err)
	return
}

func (this *Api) updateAfterReboot() {

	nos := kdbtdx.GetUnFinalizedOrderNo()
	if len(nos) < 1 {
		return
	}

	this.orderMap_lock.Lock()
	for unfinshNo, _ := range nos {
		if o,ok:=kdbtdx.GetOrder(unfinshNo);ok{
			this.orderMap[o.Qid]=o
		}
	}

	this.orderMap_lock.Unlock()


	q := QueryRequest{[]string{this.SearchAlgo}}

	rsp, err := SearchOrder(this.SearchUrl, q)
	if err != nil {
		panic(err)
	}
	if rsp.ErrorId != 0 {
		panic(rsp.ErrorMsg)
	}
	for _, rspOrder := range rsp.Data.ParentOrders {
		rspOrder := rspOrder
		this.ch <- &rspOrder
	}

}

func (this *Api) ReceiveParentOrder(order *ParentOrder) {
	this.ch <- order
}

func (this *Api) update() {
	logger.Info("########### UPDATE  START ########")

	for rspOrder := range this.ch {
		if rspOrder.Portfolio != this.SearchAlgo {
			continue
		}
		logger.Info("read resopnse msg :%#v", rspOrder)
		this.orderMap_lock.RLock()
		order, ok := this.orderMap[rspOrder.OriginId]
		this.orderMap_lock.RUnlock()
		logger.Debug("orderMap ok[%v] qid:%s,status:%v",ok,rspOrder.OriginId,order.Status)

		if ok {
			if order.Status > 3 {
				continue
			}

			//order.OrderId = order.Qid
			//if order.OrderId == "" {
			//	continue
			//}

			order.CumQty = rspOrder.FilledQuantity
			order.AvgPx = rspOrder.AvgPrice
			order.Status = StatusMapInt32[rspOrder.OrderStatus]
			if order.Status == 5 {
				order.Withdraw = order.Orderqty - order.CumQty
			}

			order.Note = rspOrder.ErrorMsg
			this.respChan <- &order
		}

	}

}

func (this *Api) Stop() {

	close(this.ch)
	close(this.respChan)
}

//webAccount  mshw_ms01_conn/mshw_ms01_qfii
//omsStocktype conn=1   qfii=2
func getChannel(channel string) (omsChannel int) {
	switch channel {
	case "conn":
		omsChannel = 1
	case "qfii":
		omsChannel = 2
	default:
		omsChannel = -1
	}
	return
}

//webSide     0/1
//omsSide  8=SELL_TO_OPEN;1=BUY_TO_CLOSE
func getOrderAction(webSide int32, accountSide *AccountSide) (omsSide OrderAction_Enum) {
	switch webSide {
	case 0:
		omsSide = accountSide.Buy
	case 1:
		omsSide = accountSide.Sell
	default:
		omsSide = -1
	}
	return
}

func (this *Api) GetAccountSide(webAccount string, webSide int32) (kdbAccount string, omsSide OrderAction_Enum, stockType int) {
	fields := strings.Split(webAccount, "_")
	if len(fields) < 2 || (webSide != 0 && webSide != 1) {
		return "", -1, -1
	}
	kdbAccount = strings.Join(fields[:len(fields)-1], "_")
	channel := fields[len(fields)-1]
	stockType = getChannel(channel)
	omsSide = getOrderAction(webSide, this.AccountSides[webAccount])
	return
}


func (this *Api) GetTrade(c *gin.Context)  {
	r:=WebResponse{}
	orders:=[]*kdbtdx.Order{}
	r.Status=true

	this.orderMap_lock.RLock()
	for _,order:=range this.orderMap{
		o:=order
		orders=append(orders,&o)
	}
	this.orderMap_lock.RUnlock()

	r.Tradings=orders

	logger.Info("search trade %v",)
	c.JSON(http.StatusOK,r)
}


type WebResponse struct {
	Status bool `json:"status"`
	Tradings []*kdbtdx.Order `json:"tradings"`
}

type Request struct {
	Channel string `json:"channel"`
	Side    int    `json:"side"`
}

type Remark struct {
	Request *Request `json:"request"`
	Msg     string   `json:"msg"`
}