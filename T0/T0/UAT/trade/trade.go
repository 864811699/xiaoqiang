package trade

import (
	"encoding/json"
	kdbtdx "github.com/864811699/T0kdb"
	logger "github.com/alecthomas/log4go"
	"io/ioutil"
	"strconv"
	"sync"
	"time"
)

type Api struct {

	sync.RWMutex
	*TradeCfg
	ConfigFile string
	TradeFile  string
	respChan chan *kdbtdx.Order
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
}



func NewApi(configFile, tradeCfg string) *Api {
	return &Api{
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
	go this.assetPosition()
	logger.Info("API run success!")

}


func (this *Api) errOrdReturn(o kdbtdx.Order, err error) {
	o.Status = 6
	o.Note = err.Error()
	logger.Error("Order Error : %v", err)
	this.respChan <- &o
}

func (this *Api) Trade(o kdbtdx.Order) {
	o.Status=0
	kdbtdx.Store(&o)
	time.Sleep(time.Millisecond*100)
	o2:=o
	if o2.Orderqty<=1000{//挂单
		o2.Status=1
	}else if o2.Orderqty>1000&&o2.Orderqty<=2000{//全成
		o2.CumQty=o.Orderqty
		o2.AvgPx=o.Askprice
		o2.Status=4
	}else if o2.Orderqty>2000&&o2.Orderqty<=4000{//部成
		o2.CumQty=o2.Orderqty/2
		o2.AvgPx=o2.Askprice
		o2.Status=2
	}else if  o.Orderqty>4000{//废单
		o2.Status=6
		o2.Note="模拟,废单"
	}

	this.respChan<-&o2
	logger.Info("trade order success :: entrustNo [%d]  ", o2.EntrustNo)
}

func (this *Api) Cancel(c kdbtdx.CancelReq) {

	orderInfo, ok := kdbtdx.GetOrder(c.Entrustno)

	if ok {
		logger.Info("cancel entrustNo: %v", c.Entrustno)
		if orderInfo.Status<4{
			orderInfo.Status=5
			orderInfo.Withdraw=orderInfo.Orderqty-orderInfo.CumQty
			this.respChan<-&orderInfo
			return
		}
	} else {
		logger.Info("cant found order :%v", c.Entrustno)
	}
}

func (this *Api)assetPosition()  {
	tk:=time.NewTicker(5*time.Second)

	for range tk.C {
		t:=time.Now().Format("150405")
		f,_:=strconv.ParseFloat(t,64)
		asset:=kdbtdx.Asset{
			Account:             "123456",
			Available:           f,
			StockValue:          f,
			AssureAsset:         f,
			TotalDebit:          f,
			Balance:             f,
			EnableBailBalance:   f,
			PerAssurescaleValue: f,
		}
		position:=kdbtdx.Position{
			Account:   "123456",
			Stockcode: "000001",
			StockName: "xxxx",
			Cost:      f,
			Gpye:      100,
			Fdyk:      100,
			Ykbl:      100,
			Sj:        100,
			Gpsz:      100,
			Kyye:      100,
		}
		kdbtdx.Tb.PushPosition(position)
		kdbtdx.Tb.PushAsset(asset)
	}
	
}


func (this *Api) GetUpdatedInfo() chan *kdbtdx.Order {

	return this.respChan
}

//func (this *Api) update() {
//	logger.Info("########### UPDATE  START ########")
//	now := time.Now().Format("20060102")
//	for rspStr := range this.RspStrCh {
//		logger.Info("read resopnse msg :%s", rspStr)
//
//		//T0交易时效性不高,增加睡眠解决异步返回导致数据被覆盖
//		time.Sleep(2 * time.Millisecond)
//
//		fields := strings.Split(rspStr, ",")
//		if len(fields) < 16 {
//			continue
//		}
//		if fields[this.Date] != now {
//			continue
//		}
//
//		packID := fields[this.EntrustNo]
//		tag, entrustNo := unpack(packID)
//		if tag != this.TagID {
//			//logger.Warn("fileds error : %v", fields)
//			continue
//		}
//
//		//this.cancelOrderMap_lock.RLock()
//		//cancelFunc, ok := this.cancelOrderMap[entrustNo]
//		//this.cancelOrderMap_lock.RUnlock()
//		//if ok {
//		//	cancelFunc()
//		//	delete(this.cancelOrderMap, entrustNo)
//		//}
//
//		order, ok := kdbtdx.GetOrder(entrustNo)
//		if ok {
//			if order.Status > 3 {
//				continue
//			}
//			logger.Info("read resopnse msg :%s", rspStr)
//			order.OrderId = fields[this.CancelID]
//			order.Status = StatusMapInt32[fields[this.Status]]
//			if order.OrderId == "" &&order.Status<4{ //广发需注释此行
//				continue
//			}
//			order.SecurityId = fmt.Sprintf("%s.%s", order.Stockcode, get_exchange(order.Stockcode))
//
//			cumQty, _ := strconv.Atoi(fields[this.CumQty])
//			order.CumQty = int32(cumQty)
//
//			avgPx, _ := strconv.ParseFloat(fields[this.AvgPx], 10)
//			order.AvgPx = avgPx
//
//			if order.Status == 5 {
//				withdraw, _ := strconv.Atoi(fields[this.CancelQty])
//				order.Withdraw = int32(withdraw)
//			}
//			msg := getMsg(fields[this.Note])
//			order.Note = getNote(order.Status, msg)
//			this.respChan <- &order
//		}
//
//	}
//
//}
//func (this *Api) errorOrder() {
//	logger.Info("########### ERROR ORDER  START ########")
//	//now := time.Now().Format("20060102")
//	for rspStr := range this.ErrStrCh {
//
//		fields := strings.Split(rspStr, ",")
//		if len(fields) < 2 {
//			continue
//		}
//
//		packID := fields[0]
//		tag, entrustNo := unpack(packID)
//		if tag != this.TagID {
//			//logger.Warn("fileds error : %v", fields)
//			continue
//		}
//
//		order, ok := kdbtdx.GetOrder(entrustNo)
//		if ok {
//			if order.Status > 3 {
//				continue
//			}
//			logger.Info("read error msg :%s", rspStr)
//
//			order.Status = 6
//			order.Note = getNote(order.Status, fields[1])
//			this.respChan <- &order
//		}
//
//	}
//
//}

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
}
