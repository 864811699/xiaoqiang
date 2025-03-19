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
	//*PBApi
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
		//PBApi:      nil,
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
	logger.Info("还券指令:29买券还券/30直接还券/31卖券还款/32直接还款")
	//for sym, accountInfo := range this.AccountInfos {
	//
	//	key := createAccountKey(accountInfo.Account, accountInfo.AccountType)
	//	if accountInfo, ok := this.AccountInfos[key]; ok {
	//		syms := accountInfo.Syms + "|" + sym
	//		a := AccountInfo{Syms: syms}
	//		this.AccountInfos[key] = a
	//
	//	} else {
	//		this.AccountInfos[key] = AccountInfo{Syms: sym}
	//	}
	//}

	return &kdbtdx.Cfg{tradeCfg.DbPath, tradeCfg.Host, tradeCfg.Port, tradeCfg.Auth, tradeCfg.Sym, tradeCfg.MaxId}
}

func (this *Api) RunApi() {

	//this.PBApi = NewPBAPI(this.AccountInfos, this.OrderDir, this.TempDir, this.ResonseDir)
	//
	go this.update()
	go this.errorOrder()
	go this.updateAssetPosition()
	//this.PBApi.start(this.IsFull)
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
	if o.ReqType != -1 && o.ReqType != 1 {
		return
	}
	side := this.getSide(o.Sym, o.Side)
	if side == "" {
		this.errOrdReturn(o, fmt.Errorf("side :%v , side should be 0/1", o.Side))
		return
	}

	exchange := get_exchange(o.Stockcode)

	o.SecurityId = fmt.Sprintf("%s%s", exchange, o.Stockcode)
	o.Status = 0
	kdbtdx.Store(&o)
	if o.ReqType == 1 {
		side = this.AccountInfos[o.Sym].Return
		this.cover(o)
		return
	}
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
	//this.orderChan <- xuntouReq

	logger.Info("trade order success :: entrustNo [%d] --> packID [%s] ", o.EntrustNo, xuntouReq.PackID)

	go this.returnRsp(*xuntouReq)

}

func (this *Api) cover(o kdbtdx.Order) {
	time.Sleep(time.Second)

	o, ok := kdbtdx.GetOrder(o.EntrustNo)
	if ok {
		o.Status = 4
		o.AvgPx = o.Askprice + 0.1
		o.CumQty = o.Orderqty
		o.Note = "还券完成"
		this.respChan <- &o
	}
}

func (this *Api) returnRsp(req XunTouReq) {
	time.Sleep(time.Second)
	_, entrustno := unpack(req.PackID)
	o, ok := kdbtdx.GetOrder(entrustno)
	if ok {
		switch {
		case req.Qty < 300:
			o.Status = 1 //已报
		case req.Qty < 1000 && req.Qty >= 300:
			o.Status = 2 //部成
			o.CumQty = req.Qty / 2
			o.AvgPx = req.Price + 0.1
		case req.Qty >= 1000 && req.Qty < 2000:
			o.Status = 4 //全成
			o.CumQty = req.Qty
			o.AvgPx = req.Price + 0.1
		case req.Qty >= 2000:
			o.Status = 6 //废单
		}
		this.respChan <- &o
	}
}

func (this *Api) Cancel(c kdbtdx.CancelReq) {

	orderInfo, ok := kdbtdx.GetOrder(c.Entrustno)

	if ok {
		logger.Info("cancel entrustNo: %v", c.Entrustno)
		if orderInfo.OrderId == "" {
			logger.Info("cant found cancelID :%v", c.Entrustno)
			return
		}

		if orderInfo.Status < 4 {
			orderInfo.Status = 5
			orderInfo.Withdraw = orderInfo.Orderqty - orderInfo.CumQty
			this.respChan <- &orderInfo
		}
		//this.orderChan <- xunTouReq
		return
	} else {
		logger.Info("cant found order :%v", c.Entrustno)
	}
}

//func (this *Api) Trade(o kdbtdx.Order) {
//
//	side := this.getSide(o.Sym, o.Side)
//	if side == "" {
//		this.errOrdReturn(o, fmt.Errorf("side :%v , side should be 0/1", o.Side))
//		return
//	}
//
//	exchange := get_exchange(o.Stockcode)
//
//	o.SecurityId = fmt.Sprintf("%s%s", exchange, o.Stockcode)
//	o.Status = 0
//	kdbtdx.Store(&o)
//
//	xuntouReq := &XunTouReq{
//		IsCancel:   false,
//		KdbAccount: o.Sym,
//		Side:       side,
//		PriceType:  "3",
//		Price:      o.Askprice,
//		SecurityID: o.SecurityId,
//		Qty:        o.Orderqty,
//		PackID:     pack(this.TagID, o.EntrustNo),
//		XunTouID:   "",
//	}
//	//this.orderChan <- xuntouReq
//
//	logger.Info("trade order success :: entrustNo [%d] --> packID [%s] ", o.EntrustNo, xuntouReq.PackID)
//	order:=o
//	switch  {
//	case order.Orderqty>=10000:
//		order.Status=6
//		order.Note="数量过多"
//
//	case order.Orderqty>=5000&&order.Orderqty<10000:
//		order.CumQty=order.Orderqty
//		order.AvgPx=order.Askprice-0.1
//		order.Status=4
//
//	case order.Orderqty>1000&&order.Orderqty<5000:
//		order.CumQty=order.Orderqty/2
//		order.AvgPx=order.Askprice-0.1
//		order.Status=2
//
//	case order.Orderqty<=1000:
//		order.Status=1
//
//	}
//
//	this.respChan <- &order
//	//go this.cancelOrder(o)
//}

//func (this *Api) Cancel(c kdbtdx.CancelReq) {
//
//	orderInfo, ok := kdbtdx.GetOrder(c.Entrustno)
//
//	if ok {
//		logger.Info("cancel entrustNo: %v", c.Entrustno)
//		//if orderInfo.OrderId == "" {
//		//	logger.Info("cant found cancelID :%v", c.Entrustno)
//		//	return
//		//}
//		orderInfo.Status=5
//		this.respChan<-&orderInfo
//		//exchange := get_exchange(orderInfo.Stockcode)
//		//xunTouReq := &XunTouReq{
//		//	IsCancel:   true,
//		//	XunTouID:   orderInfo.OrderId,
//		//	KdbAccount: orderInfo.Sym,
//		//	Exchange:   exchange,
//		//}
//		//this.orderChan <- xunTouReq
//		return
//	} else {
//		logger.Info("cant found order :%v", c.Entrustno)
//	}
//}

func (this *Api) GetUpdatedInfo() chan *kdbtdx.Order {

	return this.respChan
}

func (this *Api) update() {
	logger.Info("########### UPDATE  START ########")
	//now := time.Now().Format("20060102")
	//for rspStr := range this.RspStrCh {
	//	logger.Info("read resopnse msg :%s", rspStr)
	//
	//	//T0交易时效性不高,增加睡眠解决异步返回导致数据被覆盖
	//	time.Sleep(2 * time.Millisecond)
	//
	//	fields := strings.Split(rspStr, ",")
	//	if len(fields) < 16 {
	//		continue
	//	}
	//	if fields[this.Date] != now {
	//		continue
	//	}
	//
	//	packID := fields[this.EntrustNo]
	//	tag, entrustNo := unpack(packID)
	//	if tag != this.TagID {
	//		//logger.Warn("fileds error : %v", fields)
	//		continue
	//	}
	//
	//	//this.cancelOrderMap_lock.RLock()
	//	//cancelFunc, ok := this.cancelOrderMap[entrustNo]
	//	//this.cancelOrderMap_lock.RUnlock()
	//	//if ok {
	//	//	cancelFunc()
	//	//	delete(this.cancelOrderMap, entrustNo)
	//	//}
	//
	//	order, ok := kdbtdx.GetOrder(entrustNo)
	//	if ok {
	//		if order.Status > 3 {
	//			continue
	//		}
	//		logger.Info("read resopnse msg :%s", rspStr)
	//		order.OrderId = fields[this.CancelID]
	//		order.Status = StatusMapInt32[fields[this.Status]]
	//		if order.OrderId == "" && order.Status < 4 { //广发需注释此行
	//			continue
	//		}
	//		order.SecurityId = fmt.Sprintf("%s.%s", order.Stockcode, get_exchange(order.Stockcode))
	//
	//		cumQty, _ := strconv.Atoi(fields[this.CumQty])
	//		order.CumQty = int32(cumQty)
	//
	//		avgPx, _ := strconv.ParseFloat(fields[this.AvgPx], 10)
	//		order.AvgPx = avgPx
	//
	//		if order.Status == 5 {
	//			withdraw, _ := strconv.Atoi(fields[this.CancelQty])
	//			order.Withdraw = int32(withdraw)
	//		}
	//		msg := getMsg(fields[this.Note])
	//		order.Note = getNote(order.Status, msg)
	//		this.respChan <- &order
	//	}
	//
	//}

}

func (this *Api) errorOrder() {
	logger.Info("########### ERROR ORDER  START ########")
	//now := time.Now().Format("20060102")
	//for rspStr := range this.ErrStrCh {
	//
	//	fields := strings.Split(rspStr, ",")
	//	if len(fields) < 2 {
	//		continue
	//	}
	//
	//	packID := fields[0]
	//	tag, entrustNo := unpack(packID)
	//	if tag != this.TagID {
	//		//logger.Warn("fileds error : %v", fields)
	//		continue
	//	}
	//
	//	order, ok := kdbtdx.GetOrder(entrustNo)
	//	if ok {
	//		if order.Status > 3 {
	//			continue
	//		}
	//		logger.Info("read error msg :%s", rspStr)
	//
	//		order.Status = 6
	//		order.Note = getNote(order.Status, fields[1])
	//		this.respChan <- &order
	//	}
	//
	//}

}

func (this *Api) updateAssetPosition() {
	logger.Info("########### read asset and position ########")

	for {
		now_str:=time.Now().Format("150405")
		now,_:=strconv.Atoi(now_str)
		now_f:=float64(now)
		for k, v := range this.AccountInfos {
			asset := kdbtdx.Asset{
				Account:             v.Account,
				Available:           now_f,
				StockValue:          now_f+0.1,
				AssureAsset:         now_f+0.2,
				TotalDebit:          now_f+0.3,
				Balance:             now_f+0.4,
				EnableBailBalance:   now_f+0.5,
				PerAssurescaleValue: now_f+0.6,
				Sym:                 k,
			}

			kdbtdx.Tb.PushAsset(asset)
			position := kdbtdx.Position{
				Account:   v.Account,
				Stockcode: now_str,
				StockName: now_str,
				Cost:      now_f,
				Gpye:      int32(now),
				Fdyk:      now_f+0.1,
				Ykbl:      now_f+0.2,
				Sj:        0, //需要读取你自己的行情价格
				Gpsz:      now_f+0.3,
				Kyye:      int32(now_f)+1,
				Sym:       k,
			}
			kdbtdx.Tb.PushPosition(position)
		}
		time.Sleep(5 * time.Second)
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

func getAccountKey(s string) (account, accountType string) {
	//position.32814516_11.20231114.csv
	fields := strings.Split(s, ".")
	accountKey := fields[1]
	fields = strings.Split(accountKey, "_")
	return fields[0], fields[1]

}

func createAccountKey(account, accountType string) string {
	return fmt.Sprintf("%s-%s", account, accountType)
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
	//close(this.orderChan)

}
