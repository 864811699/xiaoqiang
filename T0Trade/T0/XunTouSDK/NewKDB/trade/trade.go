package trade

import (
	"encoding/json"
	"fmt"
	kdbtdx "github.com/864811699/T0kdb"
	logger "github.com/alecthomas/log4go"
	"io/ioutil"
	"strconv"
	"sync"
	"time"
)

type Api struct {
	*WsApi
	sync.RWMutex
	*TradeCfg
	ConfigFile   string
	TradeFile    string
	AccountToSym map[string]string
	respChan     chan *kdbtdx.Order
	offExchangeOrder map[string]*kdbtdx.Order
}

type TradeCfg struct {
	Host         string                  `json:"kdb_host"`
	Port         int                     `json:"kdb_port"`
	Auth         string                  `json:"auth"`
	DbPath       string                  `json:"dbPath"`
	Sym          []string                `json:"accounts"`
	Tag          string                  `json:"tag"`
	MaxId        int32                   `json:"maxNo"`
	WsAddr       string                  `json:"wsAddr"`
	Url          string                  `json:"url"`
	AccountInfos map[string]*AccountInfo `json:"account_infos"`
}

type AccountInfo struct {
	Buy    XunTou_Side `json:"buy"`
	Sell   XunTou_Side `json:"sell"`
	Return int         `json:"return"`
	Syms   string      `json:"syms"`
}

func NewApi(configFile, tradeCfg string) *Api {
	return &Api{
		WsApi:      nil,
		RWMutex:    sync.RWMutex{},
		TradeCfg:   nil,
		ConfigFile: configFile,
		TradeFile:  tradeCfg,
		respChan:   make(chan *kdbtdx.Order, 10000),
		offExchangeOrder:make(map[string]*kdbtdx.Order, 10000),
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

	logger.Info("load cfg success!!")
	logger.Info("return: 29 //买券还券,30 //直接还券,42 //专项买券还券,43 //专项直接还券")

	return &kdbtdx.Cfg{tradeCfg.DbPath, tradeCfg.Host, tradeCfg.Port, tradeCfg.Auth, tradeCfg.Sym, tradeCfg.MaxId}
}

func (this *Api) RunApi() {

	this.WsApi = NewWSApi(this.TradeCfg.WsAddr)
	logger.Info("init API success!!")
	go this.updateAssetPosition()
	go this.update()
	go this.updateAfterReboot()

	this.wsConnect()
	logger.Info("connect XunTou success!!")

	this.wsStart()
	logger.Info("API run success!")

}

func (this *Api) errOrdReturn(o kdbtdx.Order, err error) {
	o.Status = 6
	o.Note = err.Error()
	logger.Error("Order Error : %v", err)
	this.respChan <- &o
}

func (this *Api) Trade(o kdbtdx.Order) {
	// -1 交易 |  1 还券
	if o.ReqType != -1 && o.ReqType != 1 {
		return
	}
	reqType := NewOrder
	kdbAcountName := o.Sym
	stockCode := o.Stockcode
	//side := o.Side
	askQty := o.Orderqty
	askPrice := o.Askprice
	entrustno := o.EntrustNo
	orderType := o.Ordertype
	side := this.AccountInfos[o.Sym].Buy
	if o.Side == 1 {
		side = this.AccountInfos[o.Sym].Sell
	}
	if o.ReqType == 1 {
		side = XunTou_Side(this.AccountInfos[o.Sym].Return)
	}

	xunTou_orderType := XunTou_OrderType(0)
	switch orderType {
	case 0: //市价
		xunTou_orderType = LATEST_PRICE
	case 1: //限价
		xunTou_orderType = PRICE_MODE_LIMIT
	}

	o.SecurityId = get_SecurityId(stockCode)
	o.Status = 0
	kdbtdx.Store(&o)

	xunTouReq := &XunTouReq{
		ReqType:         reqType,
		StockCode:       o.SecurityId,
		Side:            side,
		AskQty:          int(askQty),
		AskPrice:        askPrice,
		PriceMode:       xunTou_orderType,
		KdbAcountName:   kdbAcountName,
		LocalID:         pack(this.Tag, entrustno),
		XunTouEntrustNo: 0,
	}
	this.orderChan <- xunTouReq

}

func (this *Api) Cancel(c kdbtdx.CancelReq) {

	orderInfo, ok := kdbtdx.GetOrder(c.Entrustno)

	if ok && orderInfo.Status < 3 && orderInfo.Status > 0 {
		logger.Info("cancel entrustNo: %v", c.Entrustno)
		orderID, err := strconv.Atoi(orderInfo.OrderId)
		if err != nil {
			logger.Warn("OrderID :%v fail,err :%v", orderInfo.OrderId, err)
			return
		}
		xunTouReq := &XunTouReq{
			ReqType:         CancelOrder,
			XunTouEntrustNo: orderID,
			KdbAcountName:   c.Sym,
			LocalID:         strconv.Itoa(int(c.Entrustno)),
		}
		this.orderChan <- xunTouReq
		return
	}
}

func (this *Api) GetUpdatedInfo() chan *kdbtdx.Order {

	return this.respChan
}

func (this *Api) update() {
	logger.Info("########### UPDATE  START ########")
	tk := time.NewTicker(time.Second)
	for {
		select {
		case rsp := <-this.rspChan:
			if rsp.KdbLocalID_str == "ACCOUNT STATUS" {
				if bl, msg := accountIsError(rsp.Order_status); bl {
					logger.Warn("ACCOUNT STATUS :: %v !!", msg)
					logger.Close()
					Send_weiXin("账户出现异常 :: "+msg, this.Url)
					panic(msg)
				} else {
					logger.Info("ACCOUNT STATUS :: %v !!", msg)
				}

			} else {

				tag, entrustno := unpack(rsp.KdbLocalID_str)
				if tag != this.Tag || entrustno == -1 {
					logger.Warn("非系统单 : %#v", *rsp)
					side:=getkdbSide(rsp.Order_type)
					status := int32(0)
					switch rsp.Order_status {
					case 48, 49:
						status = 0
					case 50:
						status = 1
					case 55:
						status = 2
					case 51, 52:
						status = 3
					case 56:
						status = 4
					case 53, 54:
						status = 5
					case 57:
						status = 6
					}
					withdraw := int32(0)
					if status == 5 {
						withdraw =  int32(rsp.Order_volume -rsp.Traded_volume)
					}
					sym:=""
					if accountInfo,ok:=this.AccountInfos[rsp.Account_id];ok{
						sym=accountInfo.Syms
					}

					o:=kdbtdx.Order{
						Request:    kdbtdx.Request{
							Basic:kdbtdx.Basic{
								Sym: sym,
								Qid: fmt.Sprintf("%s%s",sym,rsp.Order_sysid),
							},
							ReqType:   -1,
							Time:      time.Time{},
							ParentId:  "",
							Trader:    "",
							Fund:      "",
							Strategy:  "",
							Broker:    "",
							Stockcode: rsp.Stock_code,
							Stocktype: 0,
							Side:      side,
							Askprice:  rsp.Price,
							Ordertype: 0,
							Orderqty:  int32(rsp.Order_volume),
							Algorithm: "",
							Params:    "-1",
						},
						OrderId:    rsp.Order_sysid,
						SecurityId: rsp.Stock_code,
						EntrustNo:  -int32(rsp.Order_id),
						CumQty:     int32(rsp.Traded_volume),
						AvgPx:      rsp.Traded_price,
						Withdraw:   withdraw,
						Status:     status,
						Note:       getMsg(rsp.Status_msg),
					}
					this.respChan <-&o
					continue
				}

				order, ok := kdbtdx.GetOrder(entrustno)
				if ok && order.Status < 4 {
					logger.Info("receive msg :: %#v", *rsp)

					status := int32(0)
					switch rsp.Order_status {
					case 48, 49:
						status = 0
					case 50:
						status = 1
					case 55:
						status = 2
					case 51, 52:
						status = 3
					case 56:
						status = 4
					case 53, 54:
						status = 5
					case 57:
						status = 6
					}

					withdraw := int32(0)
					if status == 5 {
						withdraw = order.Orderqty - int32(rsp.Traded_volume)
					}
					order.OrderId = strconv.Itoa(rsp.Order_id)
					order.SecurityId = order.Stockcode
					order.AvgPx = rsp.Traded_price
					order.CumQty = int32(rsp.Traded_volume)
					order.Withdraw = withdraw
					order.Status = status
					order.Note = getMsg(rsp.Status_msg)

					this.respChan <- &order
				}
			}
		case <-tk.C:
		}
	}
}

func (this *Api) updateAssetPosition() {
	for {
		select {
		case a := <-this.AssetChan:
			syms:=""
			if accountInfo,ok:=this.AccountInfos[a.AccountID];ok{
				syms=accountInfo.Syms
			}
			asset := kdbtdx.Asset{
				Account:             a.AccountID,
				Available:           a.Cash,
				StockValue:          a.MarketValue,
				AssureAsset:         0,
				TotalDebit:          0,
				Balance:             a.TotalAsset,
				EnableBailBalance:   0,
				PerAssurescaleValue: 0,
				Sym:syms,
			}
			kdbtdx.Tb.PushAsset(asset)
		case p := <-this.PositionChan:
			syms:=""
			if accountInfo,ok:=this.AccountInfos[p.AccountID];ok{
				syms=accountInfo.Syms
			}
			position := kdbtdx.Position{
				Account:   p.AccountID,
				Stockcode: p.Stock_code[:6],
				StockName: "",
				Cost:      p.OpenPrice,
				Gpye:      int32(p.Volume),
				Fdyk:      0,
				Ykbl:      0,
				Sj:        0,
				Gpsz:      p.MarketValue,
				Kyye:      int32(p.CanUseVolume),
				Sym:syms,
			}
			kdbtdx.Tb.PushPosition(position)
		case a := <-this.CreditAssetChan:
			syms:=""
			if accountInfo,ok:=this.AccountInfos[a.AccountID];ok{
				syms=accountInfo.Syms
			}
			asset := kdbtdx.Asset{
				Account:             a.AccountID,
				Available:           a.Available,
				StockValue:          a.MarketValue,
				AssureAsset:         a.AssureAsset,
				TotalDebit:          a.TotalDebt,
				Balance:             a.Balance,
				EnableBailBalance:   a.EnableBailBalance,
				PerAssurescaleValue: a.PerAssurescaleValue,
				Sym:syms,
			}
			kdbtdx.Tb.PushAsset(asset)
		}
	}
}

func (this *Api) updateAfterReboot() {

	unfinalNos := kdbtdx.GetUnFinalizedOrderNo()
	logger.Info("unfinalNo : %v", unfinalNos)

	//if len(unfinalNos) < 1 {
	//	return
	//}

	for _, v := range this.Sym {
		xuntouSearch := &XunTouReq{
			ReqType:       SearchOrder,
			KdbAcountName: v,
		}
		this.orderChan <- xuntouSearch
	}

	go func() {
		tk := time.NewTicker(60 * time.Second)
		select {
		case <-tk.C:
			for entrustNo, _ := range unfinalNos {
				order, ok := kdbtdx.GetOrder(entrustNo)
				if ok {
					if order.Status == 0 && order.OrderId == "" {
						this.errOrdReturn(order, fmt.Errorf("Not Send Success !!!"))
					}
				}
			}
		}
	}()

}

func (this *Api) Stop() {
	close(this.respChan)
	close(this.orderChan)
	close(this.rspChan)
	close(this.unpackChan)
}

func get_SecurityId(stockCode string) string {
	sfx := ""
	switch stockCode[0] {
	case '0', '3','1':
		sfx = ".SZ"
	case '5', '6':
		sfx = ".SH"
	}
	return stockCode + sfx
}

func getkdbSide(i int)  (side int32){
	switch i {
	case 23,27,40:
		side=0
	default:
		side=1
	}
	return
}
