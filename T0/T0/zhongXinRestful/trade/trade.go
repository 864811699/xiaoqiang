package trade

import (
	"encoding/json"
	"fmt"
	kdbtdx "github.com/864811699/T0kdb"
	logger "github.com/alecthomas/log4go"
	"io/ioutil"
	"time"
)

type Api struct {
	*TradeCfg
	respChan     chan *kdbtdx.Order
	tradeCfgPath string
	restfulApi   *RestfulApi
}

type TradeCfg struct {
	Host         string                 `json:"kdb_host"`
	Port         int                    `json:"kdb_port"`
	Auth         string                 `json:"auth"`
	SubAccounts  []string               `json:"accounts"`
	DB           string                 `json:"db"`
	LoadSrc      string                 `json:"loadSrc"`
	MaxNo        int32                  `json:"maxNo"`
	UrlCfg       *UrlConfig             `json:"urls"`
	UserInfo     *UserInfo              `json:"userInfo"`
	AccountInfos map[string]AccountInfo `json:"accountInfos"`
}

type UserInfo struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type AccountInfo struct {
	FundAccount string `json:"fundAccount"`
	CashType    string `json:"cashType"`    //1普通 2专项
	EntrustType string `json:"entrustType"` //0普通 6融资 7融券
}

func NewApi(configFile, tradeCfg string) *Api {
	return &Api{
		TradeCfg:     nil,
		respChan:     make(chan *kdbtdx.Order, 10000),
		tradeCfgPath: tradeCfg,
	}
}

func (this *Api) LoadCfg() *kdbtdx.Cfg {
	readByte, err := ioutil.ReadFile(this.tradeCfgPath)
	if err != nil {
		logger.Crashf("read trade cfg :%v fail ,err: %v", this.tradeCfgPath, err)
	}
	tradeCfg := new(TradeCfg)

	err = json.Unmarshal(readByte, tradeCfg)
	if err != nil {
		logger.Crashf("tradecfg unmarshal fail,err:%v", err)
	}
	this.TradeCfg = tradeCfg

	logger.Info("load cfg success!!")

	return &kdbtdx.Cfg{tradeCfg.DB, tradeCfg.Host, tradeCfg.Port, tradeCfg.Auth, tradeCfg.SubAccounts, tradeCfg.MaxNo}
}

func (this *Api) RunApi() {

	this.restfulApi = NewRestfulApi(this.UrlCfg)
	if err := this.restfulApi.Login(this.UserInfo.UserName, this.UserInfo.Password); err != nil {
		panic(err)
	}
	logger.Info("用户登入成功")
	go this.restfulRun()
	go this.update()

	this.updateAfterReboot()
	logger.Info("API run success!")

}

func (this *Api) Trade(o kdbtdx.Order) {
	if o.ReqType!=-1{
		return
	}
	if this.restfulApi.IsOk {
		packId := pack(o.Sym, o.EntrustNo)
		side := "B"
		if o.Side != 0 {
			side = "S"
		}
		order := OrderReq{
			AcctId:         this.AccountInfos[o.Sym].FundAccount,
			ClOrdId:        packId,
			ExchId:         get_exchange_id(o.Stockcode),
			StkId:          o.Stockcode,
			OrderQty:       int(o.Orderqty),
			OrderPrice:     o.Askprice,
			OrderSide:      side,
			EntrustType:    this.AccountInfos[o.Sym].EntrustType,
			CashType:       this.AccountInfos[o.Sym].CashType,
			InvestStrategy: o.Sym,
			InvestMemo:     packId,
		}

		orderId, err := this.restfulApi.NewOrder(order)
		if err != nil {
			o.Status = 6
			o.Note = err.Error()
		} else {
			o.OrderId = orderId
			o.Status = 1
		}

	} else {
		o.Status = 6
		o.Note = "账户未登入"
	}

	kdbtdx.Store(&o)
	logger.Info("trade success,entrustNo[%v]-->orderID[%v]",o.EntrustNo,o.OrderId )
}

func (this *Api) Cancel(c kdbtdx.CancelReq) {
	order, ok := kdbtdx.GetOrder(c.Entrustno)
	if ok && order.OrderId != "" {
		logger.Info("Cancel Order :: %v----->%v---->%v", order.Sym, order.EntrustNo,order.OrderId)

		cancelReq := CancelReq{
			AcctId:  this.AccountInfos[c.Sym].FundAccount,
			ClOrdId: fmt.Sprintf("C-%d", c.Entrustno),
			OrderId: order.OrderId,
		}
		if err := this.restfulApi.NewCancel(cancelReq); err != nil {
			logger.Warn("cancel order[%v] fail,err:%v", c.Entrustno, err)
		}
	} else {
		logger.Warn("no entrustNo: ", order.EntrustNo)
	}

}

func (this *Api) GetUpdatedInfo() chan *kdbtdx.Order {
	return this.respChan
}

func (this *Api) updateAfterReboot() {
	for _, accountInfo := range this.AccountInfos {
		if orderRsps, err := this.restfulApi.NewQueryOrder(accountInfo.FundAccount); err != nil {
			logger.Error("query order rsp fail,%v", err)
		} else {
			logger.Info("查询%s委托成功", accountInfo.FundAccount)
			for k, _ := range orderRsps {
				this.checkSym(orderRsps[k].InvestStrategy)
				packId := orderRsps[k].InvestMemo
				entrustNo := unpack(packId)

				req, ok := kdbtdx.GetOrder(entrustNo)
				if ok {
					req.OrderId = orderRsps[k].OrderId
					req.CumQty = int32(orderRsps[k].KnockQty)
					req.AvgPx = orderRsps[k].KnockPrice
					req.Status = StatusMapInt32[orderRsps[k].OrderStatus]
					req.Note = orderRsps[k].ErrorMsg
					this.respChan <- &req
					logger.Info("Update entrust :: %#v", req)
				}
			}
		}

	}
}

func (this *Api) update() {

	tick := time.NewTicker(500 * time.Millisecond)

	go func() {
		logger.Info("::::::Update Return::::::")

		for {
			select {
			case orderRsp := <-this.restfulApi.RspCh:
				logger.Info("read from updateChan :: %#v", orderRsp)
				this.checkSym(orderRsp.InvestStrategy)
				local_id := unpack(orderRsp.InvestMemo)

				req, ok := kdbtdx.GetOrder(local_id)

				if ok {
					if req.Status > 3 {
						continue
					}
					if req.Side == 1 {
						req.Orderqty = -req.Orderqty
					}
					if req.OrderId==""{
						req.OrderId = orderRsp.OrderId
					}
					req.CumQty = int32(orderRsp.KnockQty)
					req.AvgPx = orderRsp.KnockPrice
					req.Status = StatusMapInt32[orderRsp.OrderStatus]
					req.Note = orderRsp.ErrorMsg
					this.respChan <- &req
					logger.Info("Update entrust :: %#v", req)
				}
			case <-tick.C:

			}
		}

	}()

}

func (this *Api) restfulRun() {

	go func() {
		if err := this.restfulApi.SubOrder(); err != nil {
			panic(fmt.Sprintf("sub order fail,%v", err))
		}
	}()

	tk := time.NewTicker(20 * time.Second)
	n := 3
	for range tk.C {
		if err := this.restfulApi.HeartBeat(); err != nil {
			logger.Error("heartBeat fail : %v", err)
			n--
			if n == 0 {
				panic(fmt.Sprintf("heartBeat check fail n:%d", n))
			}
			continue
		}
		n = 3
	}
}

func (this *Api) Stop() {

}

func (this *Api) checkSym(sym string) bool {
	_,ok:=this.AccountInfos[sym]
	return ok
}