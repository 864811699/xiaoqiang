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
	sync.RWMutex
	*TradeCfg
	respChan chan *kdbtdx.Order
	*FileOrder
	tradeCfgFile string

	//cancelOrderMap      map[int32]context.CancelFunc
	//cancelOrderMap_lock sync.RWMutex
}

type TradeCfg struct {
	Host        string   `json:"kdb_host"`
	Port        int      `json:"kdb_port"`
	Auth        string   `json:"auth"`
	SubAccounts []string `json:"accounts"`
	DB          string   `json:"db"`

	DataCfg  string   `json:"dataCfgPath"`
	Broker   string   `json:"broker"`
	LoadSrc  string   `json:"loadSrc"`
	MaxNo    int32    `json:"maxNo"`
	WTFile string `json:"wt_file"`
	EntrustFile string `json:"entrust_file"`
	Infos map[string]Info `json:"infos"`
	//DataInfo DataInfo `json:"data_info"`
}

func NewApi(configFile, tradeCfg string) *Api {
	return &Api{
		RWMutex:      sync.RWMutex{},
		TradeCfg:     nil,
		respChan:     make(chan *kdbtdx.Order, 10000),
		tradeCfgFile: tradeCfg,
		//cancelOrderMap: make(map[int32]context.CancelFunc),
	}
}

func (this *Api) LoadCfg() *kdbtdx.Cfg {
	readByte, err := ioutil.ReadFile(this.tradeCfgFile)
	if err != nil {
		logger.Crashf("read trade cfg :%v fail ,err: %v", this.tradeCfgFile, err)
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

	this.FileOrder = NewFileOrderApi(this.WTFile, this.EntrustFile)

	go this.update()
	go this.OrderApiStart()
	go this.returnAssetPosition()
	//this.updateAfterReboot()
	logger.Info("API run success!")

}

func (this *Api) Trade(order kdbtdx.Order) {

	if order.ReqType!=-1{
		return
	}

	stockcode := order.Stockcode
	marketID := get_exchange_id(stockcode)

	side := this.Infos[order.Sym].Buy
	if order.Side == 1 {
		side = this.Infos[order.Sym].Sell
	}
	qty := Abs(int(order.Orderqty))

	price := order.Askprice
	id := pack(order.Sym, order.EntrustNo)

	line := fmt.Sprintf("1|%s||%s|%s|%s|%s|%4f|%d\r\n", id, this.Infos[order.Sym].FundAccount, marketID, stockcode, side, price, qty)
	this.orderChan <- line
	kdbtdx.Store(&order)
	logger.Info("trade success,order[%#v]", order)
}

func (this *Api) Cancel(c kdbtdx.CancelReq) {
	order, ok := kdbtdx.GetOrder(c.Entrustno)
	if ok && order.OrderId != "" {
		logger.Info("Cancel Order :: %v----->%v", order.Sym, order.EntrustNo)
		line := fmt.Sprintf("2||%s||||||\r\n", order.OrderId)
		this.orderChan <- line
	} else {
		logger.Warn("no entrustNo: ", order.EntrustNo)
	}

}

func (this *Api) GetUpdatedInfo() chan *kdbtdx.Order {
	return this.respChan
}

func (this *Api) update() {

	tick := time.NewTicker(500 * time.Millisecond)

	go func() {
		logger.Info("::::::Update Return::::::")

		for {
			select {
			case rspStr := <-this.updateChan:

				fields := strings.Split(rspStr, "|")
				if len(fields) < 12 {
					continue
				}

				logger.Info("read from updateChan :: %s", rspStr)
				entrustNo := fields[0]
				local_id := unpack(entrustNo)


				req, ok := kdbtdx.GetOrder(local_id)

				if ok {
					//if req.Status>3 {
					//	continue
					//}
					logger.Info("Check update :: %v", rspStr)

					req.OrderId = fields[1]

					bidqty, err := strconv.Atoi(fields[10])
					if err != nil {
						logger.Info("bidQty :: %v , parse err ::%v", fields[10], err)
						continue
					}

					bidPrice, err := strconv.ParseFloat(fields[9], 64)
					if err != nil {
						logger.Info("avgprice :: %v , parse err ::%v", fields[9], err)
						continue
					}
					req.AvgPx=bidPrice

					status := int32(0)
					statuStr := fields[8]
					switch statuStr {
					case "2":
						status = 1
					case "7":
						status = 2
					case "3", "4":
						status = 3
					case "8":
						status = 4
					case "5", "6":
						status = 5
						req.Withdraw = req.Orderqty- int32(bidqty)
					case "9":
						status = 6
					}
					req.CumQty=int32(bidqty)
					req.Status=status
					req.Note=fields[len(fields)-1]
					this.respChan<-&req
					logger.Info("Update entrust :: %#v", req)
				}
			case <-tick.C:

			}
		}

	}()

}

func (this *Api) Stop() {

}

func (this *Api) returnAssetPosition() {
	for  {
		select {
		case asset:=<-this.assetChan:
			kdbtdx.Tb.PushAsset(asset)
		case position:=<-this.positionChan:
			kdbtdx.Tb.PushPosition(position)
		}
	}
}