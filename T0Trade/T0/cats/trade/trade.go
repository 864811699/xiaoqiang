package trade

import (
	"context"
	"encoding/json"
	"fmt"
	kdbtdx "github.com/864811699/T0kdb"
	logger "github.com/alecthomas/log4go"
	"io/ioutil"
	"sync"
)

type Api struct {
	*PBApi
	sync.RWMutex
	*TradeCfg
	ConfigFile string
	TradeFile  string

	respChan            chan *kdbtdx.Order
	cancelOrderMap      map[int32]context.CancelFunc
	cancelOrderMap_lock sync.RWMutex
}

type TradeCfg struct {
	Host         string                 `json:"kdb_host"`
	Port         int                    `json:"kdb_port"`
	Auth         string                 `json:"auth"`
	DbPath       string                 `json:"dbPath"`
	Sym          []string               `json:"accounts"`
	MaxId        int32                  `json:"maxNo"`
	RespFile     string                 `json:"respFile"`
	AssetFile    string                 `json:"assetFile"`
	OrdDir       string                 `json:"ordDir"`
	TempDir      string                 `json:"tempDir"`
	AccountInfos map[string]AccountInfo `json:"accountInfos"`
}

type AccountInfo struct {
	Account     string `json:"account"`
	AccountType string `json:"accountType"`
	Buy         string `json:"buy"`
	Sell        string `json:"sell"`
}

func NewApi(configFile, tradeCfg string) *Api {
	return &Api{
		PBApi:          nil,
		RWMutex:        sync.RWMutex{},
		TradeCfg:       nil,
		ConfigFile:     configFile,
		TradeFile:      tradeCfg,
		respChan:       make(chan *kdbtdx.Order, 10000),
		cancelOrderMap: make(map[int32]context.CancelFunc),
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

	return &kdbtdx.Cfg{tradeCfg.DbPath, tradeCfg.Host, tradeCfg.Port, tradeCfg.Auth, tradeCfg.Sym, tradeCfg.MaxId}
}

func (this *Api) RunApi() {

	this.PBApi = NewPBAPI(this.OrdDir, this.TempDir, this.RespFile,this.AssetFile)

	go this.update()
	this.PBApi.start()
	//this.updateAfterReboot()
	logger.Info("API run success!")

}
func (this *Api) errOrdReturn(o kdbtdx.Order, err error) {
	o.Status = 6
	o.Note = err.Error()
	logger.Error("Order Error : %v", err)
	this.respChan <- &o
}

func (this *Api) Trade(o kdbtdx.Order) {

	if o.ReqType != -1 {
		return
	}

	side := ""
	switch o.Side {
	case 0:
		side = this.AccountInfos[o.Sym].Buy
	case 1:
		side = this.AccountInfos[o.Sym].Sell
	}

	orderType, ok := OrderType_map[o.Ordertype]
	if !ok {
		this.errOrdReturn(o, fmt.Errorf("orderType :%v , orderType should be 1,LimitPrice", o.Ordertype))
		return
	}

	exchange := get_exchange(o.Stockcode)
	securityId := fmt.Sprintf("%s.%s", o.Stockcode, exchange)

	catsReq := &ZhongXinReq{
		ReqType:     REQ_ORDER,
		AccountType: this.AccountInfos[o.Sym].AccountType,
		Account:     this.AccountInfos[o.Sym].Account,
		SecurityId:  securityId,
		Side:        ZhongXin_Side(side),
		Price:       o.Askprice,
		Qty:         o.Orderqty,
		OrderType:   orderType,
		PackID:      pack(day, o.EntrustNo),
	}
	this.OrderChan <- catsReq

	logger.Info("trade order success :: entrustNo [%d] --> packID [%s] ", o.EntrustNo, catsReq.PackID)
	o.SecurityId = fmt.Sprintf("%s.%s", o.Stockcode, exchange)
	o.Status = 0
	this.respChan <- &o

}

func (this *Api) Cancel(c kdbtdx.CancelReq) {

	orderInfo, ok := kdbtdx.GetOrder(c.Entrustno)

	if ok {
		logger.Info("cancel entrustNo: %v", c.Entrustno)
		cancelReq := &ZhongXinReq{
			ReqType:     REQ_CANCEL,
			AccountType: this.AccountInfos[c.Sym].AccountType,
			Account:     this.AccountInfos[c.Sym].Account,
			PackID:      pack(fmt.Sprintf("%s-Cancel", day), c.Entrustno),
			CancelID:    orderInfo.OrderId,
		}
		this.OrderChan <- cancelReq
		return
	} else {
		logger.Info("cant found order :%v", c.Entrustno)
	}
}

func (this *Api) GetUpdatedInfo() chan *kdbtdx.Order {

	return this.respChan
}

func (this *Api) update() {
	logger.Info("########### UPDATE  START ########")

	for catsRsp := range this.DBFChan {
		logger.Info("read resopnse msg :%v", catsRsp)

		_, entrustNo := unpack(catsRsp.PackID)

		order, ok := kdbtdx.GetOrder(entrustNo)
		if ok {
			if order.Status > 3 {
				continue
			}

			if catsRsp.CancelID != "" {
				order.OrderId = catsRsp.CancelID
			}
			order.CumQty = int32(catsRsp.DealQty)
			order.AvgPx = catsRsp.AvgPx
			order.Withdraw = int32(catsRsp.CancelQty)
			order.Status = catsRsp.Status
			order.Note = getNote(order.Status, catsRsp.Note)
			this.respChan <- &order
		}

	}

}

func get_exchange(stockCode string) string {
	sfx := ""
	switch stockCode[0] {
	case '0', '3':
		sfx = "SZ"
	case '5', '6':
		sfx = "SH"
	}
	return sfx
}

func (this *Api) Stop() {
	close(this.respChan)
	close(this.OrderChan)
	close(this.DBFChan)

}

func getNote(status int32, tag58 string) string {
	note := ""

	switch status {
	case 0:
		note = "未报"
	case 1:
		note = "已报"
	case 2:
		note = "部成"
	case 3:
		note = "待撤"
	case 4:
		note = "已成"
	case 5:
		note = "已撤"
	case 6:
		note = "废单"
	default:
		note = "未知"

	}
	s := tag58

	return fmt.Sprintf("%s : %s", note, s)
}
