package trade

import (
	"encoding/json"
	"fmt"
	kdbtdx "github.com/864811699/T0kdb"
	logger "github.com/alecthomas/log4go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"

	"io/ioutil"
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
	WebAddr string `json:"webAddr"`
}

type AccountSide struct {
	//Sell   OrderAction_Enum `json:"sell"`
	//Buy    OrderAction_Enum `json:"buy"`
	Broker string `json:"broker"`
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
	//for _, accountInfo := range this.AccountSides {
	//	if  accountInfo.Sell != 8  && accountInfo.Buy != 1 {
	//		panic("买卖方向配置错误,8=SELL_TO_OPEN;1=BUY_TO_CLOSE")
	//	}
	//}
	logger.Info("load cfg success!!")

	return &kdbtdx.Cfg{tradeCfg.DB, tradeCfg.Host, tradeCfg.Port, tradeCfg.Auth, tradeCfg.SubAccounts, tradeCfg.MaxNo}
}

func (this *Api) RunApi() {

	if err := this.rbq.Init(this.RbqUrl); err != nil {
		logger.Crashf("rbq init fail,%v", err)
	}

	this.rbq.Sub(this, "", this.RbqSubAlgo)
	logger.Info("rbq sub success!")
	go this.update()

	this.updateAfterReboot()

	go func() {
		r:=gin.Default()
		r.GET("/getTrade",this.GetTrade)
		r.Run(this.WebAddr)
	}()

	//go Sub(this.SubUrl, this.ch,this.ctx)
	logger.Info("API run success!")
}

func (this *Api) GetUpdatedInfo() chan *kdbtdx.Order {

	return this.respChan
}

func (this *Api) Trade(o kdbtdx.Order) {
	defer kdbtdx.Store(&o)
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
	}

	if pass := checkFields(remark.Request.Channel, remark.Request.Side); !pass {
		errmsg := fmt.Sprintf("remark error,channel[%v],side[%v],channel should be conn/qfii,side should be 0/1/2/3", remark.Request.Channel, remark.Request.Side)
		logger.Warn(errmsg)
		o.Status = 6
		remark.Msg=errmsg
		bts,_:=json.Marshal(remark)
		o.Note = string(bts)
		return
	}
	side := getSide(remark.Request.Side)
	account,channelStr:=getAccountChannel(o.Sym)
	channel := getChannel(channelStr)

	o.Status = 0
	packId := pack(o.Sym, o.EntrustNo)
	req := OrderRequest{
		OriginId:    o.Qid,
		InputPrice:  o.Askprice,
		InputVol:    int(o.Orderqty),
		OrderAction: int(side),
		Remark:      packId,
		Broker:      this.AccountSides[o.Sym].Broker,
		Symbol:      o.Stockcode,
		AccountId:   account,
		Channel:     channel,
		Portfolio:this.SearchAlgo,
		StrategyId:this.StrategyId,
	}

	rsp, err := NewOrder(this.OrderUrl, req)
	if err != nil {
		o.Status = 6
		remark.Msg=err.Error()
		bts,_:=json.Marshal(remark)
		o.Note = string(bts)

		logger.Warn("order[%#v] send fail ,%v", req, err)
	} else if rsp.ErrorId != 0 {
		o.Status = 6
		remark.Msg=rsp.ErrorMsg
		bts,_:=json.Marshal(remark)
		o.Note = string(bts)

		logger.Warn("order[%#v] reject by OMS, %v", req, rsp.ErrorMsg)
	} else {
		logger.Info("trade order success :: entrustNo [%d] ", o.EntrustNo)
	}

}

func (this *Api) Cancel(c kdbtdx.CancelReq) {

	logger.Info("cancel entrustNo: %v", c.Entrustno)

	cancelRequest := CancelRequest{[]string{c.Qid}}
	rsp, err := NewCancel(this.CancelUrl, cancelRequest)
	logger.Warn("cancel send ,rsp:%#v,err:%v", rsp, err)
	return
}

func (this *Api) updateAfterReboot() {
	orders:=kdbtdx.GetOrders()
	this.orderMap_lock.Lock()
	for k, _ := range orders {
		order:=orders[k]
		this.orderMap[order.Qid]=order
	}
	this.orderMap_lock.Unlock()
	nos := kdbtdx.GetUnFinalizedOrderNo()
	if len(nos) < 1 {
		return
	}

	this.orderMap_lock.Lock()
	for unfinshNo, _ := range nos {
		if o, ok := kdbtdx.GetOrder(unfinshNo); ok {
			this.orderMap[o.Qid] = o
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
		logger.Debug("orderMap ok[%v] qid:%s,status:%v", ok, rspOrder.OriginId, order.Status)

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
			remark:=Remark{}
			if err:=json.Unmarshal([]byte(order.Note),&remark);err!=nil{
				logger.Warn("update unmarshal note fail,err:%v",err)
			}
			remark.Msg=rspOrder.ErrorMsg
			bts,_:=json.Marshal(remark)
			order.Note = string(bts)
			this.respChan <- &order

			this.orderMap_lock.RLock()
			this.orderMap[rspOrder.OriginId]=order
			this.orderMap_lock.RUnlock()
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
	case "inner":
		omsChannel=3
	default:
		omsChannel = -1
	}
	return
}

func checkFields(channel string, side int) (pass bool) {
	//if channel != "conn" && channel != "qfii" {
	//	return false
	//}

	if side != 0 && side != 1 && side != 2 && side != 3 {
		return false
	}
	return true
}

func getSide(side int) (omsSide OrderAction_Enum) {
	switch side {
	case 0:
		omsSide = OrderAction_BUY_TO_OPEN
	case 1:
		omsSide = OrderAction_SELL_TO_CLOSE
	case 2:
		omsSide = OrderAction_BUY_TO_CLOSE
	case 3:
		omsSide = OrderAction_SELL_TO_OPEN
	}
	return omsSide
}

func getAccountChannel(sym string)(account,channel string)  {
	fields:=strings.Split(sym,"-")
	return fields[0],fields[1]
}


type Request struct {
	Channel string `json:"channel"`
	Side    int    `json:"side"`
}

type Remark struct {
	Request *Request `json:"request"`
	Msg     string   `json:"msg"`
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

//func (this *Api) GetHoldings(c *gin.Context) {
//	r:=Response{}
//	client:=pb.NewOrderServiceClient(this.ClientConn)
//	res,err:=client.QueryHolding(context.Background(),&pb.QueryRequest{})
//	if err != nil {
//		Logger.Warnf("query holding fail,%v",err)
//		r.Status=false
//		return
//	}else {
//		r.Status=true
//		r.Holdings=res.Holdings
//	}
//	Logger.Infof("search hold %v",res.Holdings)
//	c.JSON(http.StatusOK,r)
//}