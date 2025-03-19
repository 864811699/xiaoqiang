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
	//sync.RWMutex
	*TradeCfg
	respChan           chan *kdbtdx.Order
	tradeCfgFile       string
	unfinishOrder      map[int]int
	unfinishorder_lock sync.RWMutex
	responseChan       chan string
	orderMap           map[string]*kdbtdx.Order //通过合同编号映射实际委托
	orderMap_lock      sync.RWMutex
	isTraded_m         map[int]bool
	isTraded_m_lock    sync.RWMutex
	accountToSym       map[int]string
}

type TradeCfg struct {
	Host         string                  `json:"kdb_host"`
	Port         int                     `json:"kdb_port"`
	Auth         string                  `json:"auth"`
	SubAccounts  []string                `json:"accounts"`
	DB           string                  `json:"db"`
	Broker       string                  `json:"broker"`
	MaxNo        int32                   `json:"maxNo"`
	AccountInfos map[string]*AccountInfo `json:"account_infos"`
	TraderUser   string                  `json:"trader_user"`
	TraderPwd    string                  `json:"trader_pwd"`
	QuerryTk     time.Duration           `json:"querry_tk"`
}

type AccountInfo struct {
	Name string `json:"name"`
	Id   int
	Buy  SideType `json:"buy"`
	Sell SideType `json:"sell"`
}

func NewApi(configFile, tradeCfg string) *Api {
	return &Api{
		TradeCfg:      nil,
		respChan:      make(chan *kdbtdx.Order, 10000),
		tradeCfgFile:  tradeCfg,
		unfinishOrder: make(map[int]int, 10000),
		responseChan:  make(chan string, 1000),
		orderMap:      make(map[string]*kdbtdx.Order, 10000),
		isTraded_m:    make(map[int]bool, 100),
		accountToSym:  make(map[int]string, 100),
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
	for _, info := range tradeCfg.AccountInfos {
		if info.Sell != MAIQUANHUANKUAN_SELL && info.Sell != DANBAOPING_SELL && info.Buy != DANBAOPING_BUY && info.Buy != RONGZI_BUY && info.Sell != RONGQUAN_SELL {
			logger.Crashf("account info error,融资买入[%s],担保品买入[%s],卖券还款[%s],担保品卖出[%s],融券卖出[%s]", RONGZI_BUY, DANBAOPING_BUY, MAIQUANHUANKUAN_SELL, DANBAOPING_SELL, RONGQUAN_SELL)
		}
	}
	logger.Info("load cfg success!!")

	return &kdbtdx.Cfg{tradeCfg.DB, tradeCfg.Host, tradeCfg.Port, tradeCfg.Auth, tradeCfg.SubAccounts, tradeCfg.MaxNo}
}

func (this *Api) RunApi() {
	this.login()
	this.afterReboot()
	this.update()
	go this.crontabQueryAsset()
	go this.crontabQueryPosition()
	go this.crontabQueryResponse()
	logger.Info("API run success!")

}

func (this *Api) login() {

	loginReq := LoginReq{
		Request: Request{Cmd: string(LOGIN)},
		Ps: Ps{
			User: this.TraderUser,
			Pwd:  this.TraderPwd,
		},
	}
	bts, _ := json.Marshal(loginReq)
	name, err := CallRequestData(bts)
	if err != nil {
		logger.Warn("trader[%s] login fail,%v ", name, err)
		logger.Warn("fail msg:%s", bts)
		return
	}
	name, _ = GbkToUtf8(name)
	logger.Info("trader[%s] login success", name)

	r := Request{Cmd: string(QUERRY_ACCOUNT_INFO)}
	bts, _ = json.Marshal(r)
	if msg, err := CallRequestData(bts); err != nil {
		logger.Warn("account[%s] query accountInfo fail,", name, err)
		return
	} else {
		msg, _ = GbkToUtf8(msg)
		accountInfos := AccountInfoRsp{}
		json.Unmarshal([]byte(msg), &accountInfos)
		if accountInfos.SignalRsp.Ret != 0 {
			logger.Warn("account query fail,ret:%v", accountInfos.SignalRsp.Ret)
			return
		}
		logger.Info("account query success,%#v", accountInfos)
		for _, accountInfo := range accountInfos.AccountInfos {
			for sym, info := range this.AccountInfos {
				if info.Name == accountInfo.Name {
					this.AccountInfos[sym].Id = accountInfo.ID
					this.unfinishOrder[accountInfo.ID] = 0
					if syms, ok := this.accountToSym[accountInfo.ID]; ok {
						syms = fmt.Sprintf("%s|%s", syms, sym)
					} else {
						this.accountToSym[accountInfo.ID] = sym
					}
				}
			}
		}

	}

}

func (this *Api) Trade(order kdbtdx.Order) {
	if order.ReqType != -1 && order.ReqType != 1 {
		return
	}
	var wtFlag SideType
	side := BUY
	wtFlag = this.AccountInfos[order.Sym].Buy
	if order.Side == 1 {
		side = SELL
		wtFlag = this.AccountInfos[order.Sym].Sell
	}

	o := OrderReq{
		Request:   Request{Cmd: string(side)},
		Stockcode: order.Stockcode,
		Price:     order.Askprice,
		ID:        this.AccountInfos[order.Sym].Id,
		Amount:    int(order.Orderqty),
		WtFlag:    wtFlag,
	}
	bts, _ := json.Marshal(o)
	orderRsp := OrderRsp{}
	flage := "fail"
	msg, err := CallRequestData(bts)
	if err != nil {
		order.Status = 6
	} else {
		msg, _ = GbkToUtf8(msg)
		_ = json.Unmarshal([]byte(msg), &orderRsp)
		if orderRsp.Ret != 0 {
			order.Status = 6
			order.Note = orderRsp.ErrMsg
		} else {
			order.Status = 1
			order.OrderId = orderRsp.OrderId
			flage = "success"
		}
	}
	this.unfinishorder_lock.Lock()
	this.unfinishOrder[this.AccountInfos[order.Sym].Id]++
	this.unfinishorder_lock.Unlock()
	this.orderMap_lock.Lock()
	this.orderMap[orderRsp.OrderId] = &order
	this.orderMap_lock.Unlock()
	kdbtdx.Store(&order)
	this.isTraded_m_lock.RLock()
	this.isTraded_m[this.AccountInfos[order.Sym].Id] = true
	this.isTraded_m_lock.RUnlock()
	logger.Info("trade %s,order[%#v]", flage, order)
}

func (this *Api) Cancel(c kdbtdx.CancelReq) {
	order, ok := kdbtdx.GetOrder(c.Entrustno)
	if ok && order.OrderId != "" {
		if order.Status > 3 {
			logger.Info("Cancel Order fail:: %v----->%v---->%v", order.Sym, order.EntrustNo, order.Status)
			return
		}
		logger.Info("Cancel Order :: %v----->%v", order.Sym, order.EntrustNo)
		r := CancelReq{
			Request:   Request{Cmd: string(CANCELOrder)},
			ID:        this.AccountInfos[c.Sym].Id,
			OrderId:   order.OrderId,
			Stockcode: order.Stockcode,
		}

		bts, _ := json.Marshal(r)
		CallRequestData(bts)
	} else {
		logger.Warn("no entrustNo: ", order.EntrustNo)
	}

}

func (this *Api) GetUpdatedInfo() chan *kdbtdx.Order {
	return this.respChan
}

func (this *Api) crontabQueryResponse() {

	querys := make(map[int][]byte, 0)
	for _, accountInfo := range this.AccountInfos {
		r := QueryReq{
			Request: Request{Cmd: string(QUERY_ORDERS)},
			ID:      accountInfo.Id,
		}
		bts, _ := json.Marshal(r)
		querys[accountInfo.Id] = bts
	}
	for id, bts := range querys {
		go func(id int, bts []byte) {
			tk := time.NewTicker(time.Second * this.QuerryTk)
			tk_query := time.NewTicker(5 * 60 * time.Second)
			for {
				select {
				case _ = <-tk.C:
					this.unfinishorder_lock.RLock()
					unfinishNum, ok := this.unfinishOrder[id]
					this.unfinishorder_lock.RUnlock()
					if ok && unfinishNum > 0 {
						msg, err := CallRequestData(bts)
						if err != nil {
							logger.Warn("querry orders fail,", err)
							continue
						}
						msg, _ = GbkToUtf8(msg)
						this.responseChan <- msg
					}

				case _ = <-tk_query.C:
					this.unfinishorder_lock.RLock()
					unfinishNum, ok := this.unfinishOrder[id]
					this.unfinishorder_lock.RUnlock()
					this.isTraded_m_lock.RLock()
					isTrade, ok2 := this.isTraded_m[id]
					this.isTraded_m_lock.RUnlock()

					if ok && unfinishNum <= 0 && ok2 && isTrade {
						msg, err := CallRequestData(bts)
						if err != nil {
							logger.Warn("querry orders fail,", err)
							continue
						}
						msg, _ = GbkToUtf8(msg)
						this.responseChan <- msg
					}
				}
			}
		}(id, bts)
	}

}

func (this *Api) update() {

	tick := time.NewTicker(500 * time.Millisecond)

	go func() {
		logger.Info("::::::Update Return::::::")

		for {
			select {
			case rspStr := <-this.responseChan:
				logger.Warn("query api msg:%v",rspStr)
				r := QueryRsp{}

				if err := json.Unmarshal([]byte(rspStr), &r); err != nil {
					logger.Warn("query orders unmarshal fail,msg:%s,err:%v", rspStr, err)
					continue
				}
				if r.Ret != 0 {
					logger.Warn("querry orders fail,ret:%v", r.Ret)
					continue
				}

				//{
				//		"op_name": "卖出",
				//		"wtlx": "",
				//		"wtrq": "20240306",
				//		"note": "全部成交",
				//		"opid": 2100,
				//		"uuid": 9589,
				//		"htbh": "X1456535",
				//		"op": "S",
				//		"reqtm": 93,
				//		"cjjg": 39.05,
				//		"tact": "zjlqxi",
				//		"wtjg": 39.05,
				//		"wtsj": "09:49:56",
				//		"status": 1,
				//		"cdsl": 0,
				//		"wtsl": 2500,
				//		"jys": "SHA",
				//		"uname": "002812100",
				//		"cjsl": 2500,
				//		"cjje": 97625,
				//		"zqdm": "688159",
				//		"zqmc": "有方科技",
				//		"unick": "zh2023",
				//		"qsname": "中金证券(信用)",
				//		"qsid": 17
				//	}
				for _, response := range r.Responses {
					this.orderMap_lock.RLock()
					order, ok := this.orderMap[response.OrderId]
					this.orderMap_lock.RUnlock()
					if ok && order.Status == 6 {
						logger.Info("receive error order :: %#v", response)
						continue
					}
					if ok && (order.Status < 4 || order.Status == 6) {
						logger.Info("receive order update :: %#v", response)

						order := *order
						order.AvgPx = response.BidPrice
						order.CumQty = int32(response.BidQty)

						order.Status = OrderStatusMap[response.Status]
						if order.Status == 5 {
							order.Withdraw = order.Orderqty - order.CumQty
						}
						order.Note = response.Note
						this.orderMap[response.OrderId] = &order
						this.respChan <- &order
						logger.Info("Update entrust :: %#v", order)

						if order.Status > 3 {
							this.unfinishorder_lock.Lock()
							this.unfinishOrder[this.AccountInfos[order.Sym].Id]--
							this.unfinishorder_lock.Unlock()
						}

					}
				}

			case <-tick.C:

			}
		}

	}()

}

func (this *Api) afterReboot() {
	//entrustNoMap := kdbtdx.GetUnFinalizedOrderNo()
	//for entrustno, _ := range entrustNoMap {
	//	if o, ok := kdbtdx.GetOrder(entrustno); ok {
	//		this.orderMap[o.OrderId] = &o
	//		this.unfinishOrder[this.AccountInfos[o.Sym].Id]++
	//		this.isTraded_m_lock.RLock()
	//		this.isTraded_m[this.AccountInfos[o.Sym].Id] = true
	//		this.isTraded_m_lock.RUnlock()
	//	}
	//}
	//
	orders := kdbtdx.GetAllOrder()
	for entrustNo, _ := range orders {
		order:=orders[entrustNo]
		this.orderMap[order.OrderId] = &order
		this.isTraded_m[this.AccountInfos[orders[entrustNo].Sym].Id] = true
		if orders[entrustNo].Status < 4 {
			this.unfinishOrder[this.AccountInfos[orders[entrustNo].Sym].Id]++
		}
	}

}

func (this *Api) crontabQueryAsset() {

	querys := make(map[int][]byte, 0)
	for _, accountInfo := range this.AccountInfos {
		r := QueryReq{
			Request: Request{Cmd: string(QUERY_FUND)},
			ID:      accountInfo.Id,
		}
		bts, _ := json.Marshal(r)
		querys[accountInfo.Id] = bts
	}
	for id, bts := range querys {
		go func(id int, bts []byte) {
			tk := time.NewTicker(time.Second * 5)
			for range tk.C {
				this.isTraded_m_lock.RLock()
				_, ok := this.isTraded_m[id]
				this.isTraded_m_lock.RUnlock()
				if !ok {
					continue
				}
				msg, err := CallRequestData(bts)
				if err != nil {
					logger.Warn("querry asset fail,", err)
					continue
				}
				msg, _ = GbkToUtf8(msg)
				logger.Info("query asset:%s", msg)
				asset := QueryAssetRsp{}
				err = json.Unmarshal([]byte(msg), &asset)
				if err != nil {
					logger.Warn("position unmarshal fail,err:%v", err)
					continue
				}
				if len(asset.Responses) == 1 {
					sym := ""
					if syms, ok := this.accountToSym[id]; ok {
						sym = syms
					}
					asset.Responses[0].Sym = sym
					asset.Responses[0].Account = strconv.Itoa(id)
					logger.Info("asset:%#v", asset.Responses[0])
					kdbtdx.Tb.PushAsset(asset.Responses[0])
				}

			}

		}(id, bts)
	}
}

func (this *Api) crontabQueryPosition() {

	querys := make(map[int][]byte, 0)
	for _, accountInfo := range this.AccountInfos {
		r := QueryReq{
			Request: Request{Cmd: string(QUERY_HOLDING)},
			ID:      accountInfo.Id,
		}
		bts, _ := json.Marshal(r)
		querys[accountInfo.Id] = bts
	}

	for id, bts := range querys {
		go func(id int, bts []byte) {
			tk := time.NewTicker(time.Second * 3)
			for range tk.C {

				this.isTraded_m_lock.RLock()
				_, ok := this.isTraded_m[id]
				this.isTraded_m_lock.RUnlock()
				if !ok {
					continue
				}

				msg, err := CallRequestData(bts)
				if err != nil {
					logger.Warn("querry position fail,", err)
					continue
				}
				msg, _ = GbkToUtf8(msg)
				logger.Info("query position:%s", msg)
				//position:=kdbtdx.Position{}
				rsp := QueryPositionRsp{}
				err = json.Unmarshal([]byte(msg), &rsp)
				if err != nil {
					logger.Warn("position unmarshal fail,err:%v", err)
					continue
				}
				for k, _ := range rsp.Responses {
					position := *rsp.Responses[k]
					sym := ""
					if syms, ok := this.accountToSym[id]; ok {
						sym = syms
					}
					position.Sym = sym
					position.Account = strconv.Itoa(id)
					kdbtdx.Tb.PushPosition(position)
				}
				//position.Account=strconv.Itoa(id)
				//kdbtdx.Tb.PushPosition(position)
			}
		}(id, bts)
	}
}

func (this *Api) Stop() {

}
