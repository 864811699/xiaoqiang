package trade

import (
	"encoding/json"
	"fmt"
	kdbtdx "github.com/864811699/T0kdb"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
	"time"

	logger "github.com/alecthomas/log4go"
)

type Api struct {
	*PBApi
	sync.RWMutex
	*TradeCfg
	ConfigFile string
	TradeFile  string

	respChan chan *kdbtdx.Order
	//cancelOrderMap      map[int32]context.CancelFunc
	//cancelOrderMap_lock sync.RWMutex
}

type TradeCfg struct {
	Host   string   `json:"kdb_host"`
	Port   int      `json:"kdb_port"`
	Auth   string   `json:"auth"`
	DbPath string   `json:"dbPath"`
	Sym    []string `json:"accounts"`
	MaxId  int32    `json:"maxNo"`
	TagID  string   `json:"tagID"`

	Timeout      time.Duration          `json:"timeout"`
	IsFull       bool                   `json:"isFull"`
	OrderDir     string                 `json:"order_dir"`
	TempDir      string                 `json:"temp_dir"`
	ResonseDir   string                 `json:"resonse_dir"`
	AccountInfos map[string]AccountInfo `json:"account_infos"`
	Fields
}

type Fields struct {
	Date      int `json:"date"`
	CancelID  int `json:"cancelID"`
	EntrustNo int `json:"entrustNo"`
	CumQty    int `json:"dealQty"`
	AvgPx     int `json:"price"`
	CancelQty int `json:"cancelQty"`
	Status    int `json:"status"`
	Note      int `json:"note"`
}

func NewApi(configFile, tradeCfg string) *Api {
	return &Api{
		PBApi:      nil,
		RWMutex:    sync.RWMutex{},
		TradeCfg:   nil,
		ConfigFile: configFile,
		TradeFile:  tradeCfg,
		respChan:   make(chan *kdbtdx.Order, 10000),
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

	logger.Info("load cfg success!!")

	return &kdbtdx.Cfg{tradeCfg.DbPath, tradeCfg.Host, tradeCfg.Port, tradeCfg.Auth, tradeCfg.Sym, tradeCfg.MaxId}
}

func (this *Api) RunApi() {

	this.PBApi = NewPBAPI(this.AccountInfos, this.OrderDir, this.TempDir, this.ResonseDir)

	go this.update()
	go this.errorOrder()
	this.PBApi.start(this.IsFull)
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

	side := this.getSide(o.Sym, o.Side)
	if side == "" {
		this.errOrdReturn(o, fmt.Errorf("side :%v , side should be 0/1", o.Side))
		return
	}

	exchange := get_exchange(o.Stockcode)

	o.SecurityId = fmt.Sprintf("%s%s", exchange, o.Stockcode)
	o.Status = 0
	kdbtdx.Store(&o)

	xuntouReq := &XunTouReq{
		IsCancel:   false,
		KdbAccount: o.Sym,
		Side:       side,
		PriceType:  "3",
		Price:      o.Askprice,
		SecurityID: o.SecurityId,
		Qty:        o.Orderqty,
		PackID:     pack(this.TagID, o.EntrustNo),
		XunTouID:   "",
	}
	this.orderChan <- xuntouReq

	logger.Info("trade order success :: entrustNo [%d] --> packID [%s] ", o.EntrustNo, xuntouReq.PackID)

	//this.respChan <- &o

	//go this.cancelOrder(o)

}

func (this *Api) Cancel(c kdbtdx.CancelReq) {

	orderInfo, ok := kdbtdx.GetOrder(c.Entrustno)

	if ok {
		logger.Info("cancel entrustNo: %v", c.Entrustno)
		if orderInfo.OrderId == "" {
			logger.Info("cant found cancelID :%v", c.Entrustno)
			return
		}
		exchange := get_exchange(orderInfo.Stockcode)
		xunTouReq := &XunTouReq{
			IsCancel:   true,
			XunTouID:   orderInfo.OrderId,
			KdbAccount: orderInfo.Sym,
			Exchange:   exchange,
		}
		this.orderChan <- xunTouReq
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
	now := time.Now().Format("20060102")
	for rspStr := range this.RspStrCh {
		logger.Info("read resopnse msg :%s", rspStr)

		//T0交易时效性不高,增加睡眠解决异步返回导致数据被覆盖
		time.Sleep(2 * time.Millisecond)

		fields := strings.Split(rspStr, ",")
		if len(fields) < 16 {
			continue
		}
		if fields[this.Date] != now {
			continue
		}

		packID := fields[this.EntrustNo]
		tag, entrustNo := unpack(packID)
		if tag != this.TagID {
			//logger.Warn("fileds error : %v", fields)
			continue
		}

		//this.cancelOrderMap_lock.RLock()
		//cancelFunc, ok := this.cancelOrderMap[entrustNo]
		//this.cancelOrderMap_lock.RUnlock()
		//if ok {
		//	cancelFunc()
		//	delete(this.cancelOrderMap, entrustNo)
		//}

		order, ok := kdbtdx.GetOrder(entrustNo)
		if ok {
			if order.Status > 3 {
				continue
			}
			logger.Info("read resopnse msg :%s", rspStr)
			order.OrderId = fields[this.CancelID]
			if order.OrderId == "" {
				continue
			}
			order.SecurityId = fmt.Sprintf("%s.%s", order.Stockcode, get_exchange(order.Stockcode))

			cumQty, _ := strconv.Atoi(fields[this.CumQty])
			order.CumQty = int32(cumQty)

			avgPx, _ := strconv.ParseFloat(fields[this.AvgPx], 10)
			order.AvgPx = avgPx
			order.Status = StatusMapInt32[fields[this.Status]]
			if order.Status == 5 {
				withdraw, _ := strconv.Atoi(fields[this.CancelQty])
				order.Withdraw = int32(withdraw)
			}
			msg := getMsg(fields[this.Note])
			order.Note = getNote(order.Status, msg)
			this.respChan <- &order
		}

	}

}
func (this *Api) errorOrder() {
	logger.Info("########### ERROR ORDER  START ########")
	//now := time.Now().Format("20060102")
	for rspStr := range this.ErrStrCh {

		fields := strings.Split(rspStr, ",")
		if len(fields) < 2 {
			continue
		}

		packID := fields[0]
		tag, entrustNo := unpack(packID)
		if tag != this.TagID {
			//logger.Warn("fileds error : %v", fields)
			continue
		}

		order, ok := kdbtdx.GetOrder(entrustNo)
		if ok {
			if order.Status > 3 {
				continue
			}
			logger.Info("read error msg :%s", rspStr)

			order.Status = 6
			order.Note = getNote(order.Status, fields[1])
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

//func (this *Api) cancelOrder(o kdbtdx.Order) {
//	tk := time.NewTicker(this.Timeout * time.Second)
//
//	ctx, cancelFunc := context.WithCancel(context.Background())
//	this.cancelOrderMap_lock.Lock()
//	this.cancelOrderMap[o.EntrustNo] = cancelFunc
//	this.cancelOrderMap_lock.Unlock()
//
//	select {
//	case <-ctx.Done():
//		return
//	case <-tk.C:
//
//		o.Status = 6
//		o.Note = "超时废单"
//		this.respChan <- &o
//		logger.Info("订单未回报,超时废单,account:%s :%d", o.Sym, o.EntrustNo)
//
//	}
//}

func (this *Api) getSide(sym string, n int32) string {
	side := ""
	switch n {
	case 0:
		side = this.AccountInfos[sym].Buy
	case 1:
		side = this.AccountInfos[sym].Sell
	}
	return side
}
func (this *Api) Stop() {
	close(this.respChan)
	close(this.orderChan)

}
