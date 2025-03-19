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
	logger.Info("还券指令:29买券还券/30直接还券/31卖券还款/32直接还款")
	logger.Info("account key : account-accountType")
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

	this.PBApi = NewPBAPI(this.AccountInfos, this.OrderDir, this.TempDir, this.ResonseDir)
	//
	go this.update()
	go this.errorOrder()
	go this.updateAssetPosition()
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
			order.Status = StatusMapInt32[fields[this.Status]]
			if order.OrderId == "" && order.Status < 4 { //广发需注释此行
				continue
			}
			order.SecurityId = fmt.Sprintf("%s.%s", order.Stockcode, get_exchange(order.Stockcode))

			cumQty, _ := strconv.Atoi(fields[this.CumQty])
			order.CumQty = int32(cumQty)

			avgPx, _ := strconv.ParseFloat(fields[this.AvgPx], 10)
			order.AvgPx = avgPx

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

func (this *Api) updateAssetPosition() {
	logger.Info("########### read asset and position ########")
	for {
		select {
		case rspStr := <-this.AssetCh:
			logger.Info("asset msg:: %s ", rspStr)
			fields := strings.Split(rspStr, ",")
			if len(fields) < 13 {
				logger.Error("asset read fail,%s", rspStr)
				continue
			}

			if _, err := strconv.ParseFloat(fields[0],64); err != nil {
				continue
			}

			//总资产	可用金额	总市值	股票总市值	债券总市值	总负债	可用保证金	个人维持担	净资产	融资合约金	融资合约费	融资合约利	融券合约金	融券市值	融券合约费	融券合约利	其他费用	融资授信额	融资可用额	融资已用额	融券授信额	融券可用额	融券已用额	剩余融券卖出资金	账号状态	账号key	盈亏	交易日

			fileName := fields[len(fields)-1]
			account, accountType := getAccountKey(fileName)
			key := createAccountKey(account, accountType)
			syms := ""
			if accountInfo, ok := this.AccountInfos[key]; ok {
				syms = accountInfo.Syms
			}
			asset := kdbtdx.Asset{
				Account: key,
				Sym:     syms,
			}
			if available, ok := strconv.ParseFloat(fields[1], 64); ok == nil {
				asset.Available = available
			}
			if stockValue, ok := strconv.ParseFloat(fields[3], 64); ok == nil {
				asset.StockValue = stockValue
			}
			if assureAsset, ok := strconv.ParseFloat(fields[8], 64); ok == nil {
				asset.AssureAsset = assureAsset
			}
			if TotalDebit, ok := strconv.ParseFloat(fields[5], 64); ok == nil {
				asset.TotalDebit = TotalDebit
			}
			if Balance, ok := strconv.ParseFloat(fields[0], 64); ok == nil {
				asset.Balance = Balance
			}
			if EnableBailBalance, ok := strconv.ParseFloat(fields[6], 64); ok == nil {
				asset.EnableBailBalance = EnableBailBalance
			}
			if PerAssurescaleValue, ok := strconv.ParseFloat(fields[7], 64); ok == nil {
				asset.PerAssurescaleValue = PerAssurescaleValue
			}

			kdbtdx.Tb.PushAsset(asset)


		case rspStr := <-this.PositionCh:
			//证券代码	市场名称	持仓数量	可用数量	在途股份	卖出冻结	成本价	最新价	昨夜拥股	盈亏	市值
			fields := strings.Split(rspStr, ",")
			if len(fields) < 11 {
				logger.Error("position read fail,%s", rspStr)
				continue
			}

			if _, err := strconv.Atoi(fields[0]); err != nil {
				continue
			}

			fileName := fields[len(fields)-1]
			account, accountType := getAccountKey(fileName)
			key := createAccountKey(account, accountType)
			syms := ""
			if accountInfo, ok := this.AccountInfos[key]; ok {
				syms = accountInfo.Syms
			}
			position := kdbtdx.Position{
				Account:   key,
				Stockcode: fields[0],
				StockName: "",
				Ykbl:      0,
				Sym:       syms,
			}

			if Cost, ok := strconv.ParseFloat(fields[6], 64); ok == nil {
				position.Cost = Cost
			}
			if Gpye, ok := strconv.Atoi(fields[2]); ok == nil {
				position.Gpye = int32(Gpye)
			}
			if Fdyk, ok := strconv.ParseFloat(fields[9], 64); ok == nil {
				position.Fdyk = Fdyk
			}
			if Sj, ok := strconv.ParseFloat(fields[7], 64); ok == nil {
				position.Sj = Sj
			}
			if Gpsz, ok := strconv.ParseFloat(fields[10], 64); ok == nil {
				position.Gpsz = Gpsz
			}
			if Kyye, ok := strconv.Atoi(fields[3]); ok == nil {
				position.Kyye = int32(Kyye)
			}
			kdbtdx.Tb.PushPosition(position)
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
	close(this.orderChan)

}
