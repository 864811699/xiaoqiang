package trade

import (
	"encoding/json"
	logger "github.com/alecthomas/log4go"
	"io"
	"strconv"
	"time"

	"net"
)

type RspEntrust struct {
	Sym       string
	Qid       string
	Trader    string
	Time      time.Time
	Entrustno int32
	OrderID   string
	Stockcode string
	Side      int32
	Askprice  float64
	Askvol    int32
	Bidprice  float64
	Bidvol    int32
	Withdraw  int32
	Resptime  time.Time
	Status    int32
	Note      string
}

type MsgType struct {
	MsgType string `json:"msg_type"` //O委托回报  A资金  P持仓	T状态  CA信用资产
}

type XunTouRsp struct {
	Account_type   int     `json:"account_type"` //2普通  3信用
	Account_id     string  `json:"account_id"`
	Stock_code     string  `json:"stock_code"`
	Order_id       int     `json:"order_id"`
	Order_sysid    string  `json:"order_sysid"`
	Order_time     int     `json:"order_time"`
	Order_type     int     `json:"order_type"`
	Order_volume   int     `json:"order_volume"`
	Price_type     int     `json:"price_type"`
	Price          float64 `json:"price"`
	Traded_volume  int     `json:"traded_volume"`
	Traded_price   float64 `json:"traded_price"`
	Order_status   int     `json:"order_status"`
	Status_msg     string  `json:"status_msg"`
	KdbAcountName  string  `json:"strategy_name"` //KDB表名
	KdbLocalID_str string  `json:"order_remark"`  //委托备注
}

type XunTouReq struct {
	ReqType         ReqType          `json:"reqType"` //O C S  下单/撤单/查询
	StockCode       string           `json:"stock_code"`
	Side            XunTou_Side      `json:"side"`
	AskQty          int              `json:"askQty"`
	AskPrice        float64          `json:"askPrice"`
	PriceMode       XunTou_OrderType `json:"priceMode"`     //PRICE_MODE_LIMIT = 11
	KdbAcountName   string           `json:"KdbAcountName"` //KDB表名
	LocalID         string           `json:"KdbLocalID"`    //委托备注
	XunTouEntrustNo int              `json:"entrust_no"`    //讯投委托号
}

type Asset struct {
	AccountID   string  `json:"account_id"`
	AccountType int     `json:"account_type"`
	Cash        float64 `json:"cash"`
	FrozenCash  float64 `json:"frozen_cash"`
	MarketValue float64 `json:"market_value"`
	TotalAsset  float64 `json:"total_asset"`
}

type CreditAsset struct {
	AccountID           string  `json:"account_id"`
	AccountType         int     `json:"account_type"`
	Balance             float64 `json:"m_dBalance"`
	Available           float64 `json:"m_dAvailable"`
	MarketValue         float64 `json:"m_dMarketValue"`
	StockValue          float64 `json:"m_dStockValue"`
	TotalDebt           float64 `json:"m_dTotalDebt"`
	EnableBailBalance   float64 `json:"m_dEnableBailBalance"`
	AssureAsset         float64 `json:"m_dAssureAsset"`
	PerAssurescaleValue float64 `json:"m_dPerAssurescaleValue"`
}

type Position struct {
	AccountID    string  `json:"account_id"`
	AccountType  int     `json:"account_type"`
	Stock_code   string  `json:"stock_code"`
	Volume       int     `json:"volume"`
	CanUseVolume int     `json:"can_use_volume"`
	OpenPrice    float64 `json:"open_price"`
	MarketValue  float64 `json:"market_value"`
}
type ReqType string

const (
	NewOrder    ReqType = "O"
	CancelOrder ReqType = "C"
	SearchOrder ReqType = "S"
)

type XunTou_OrderType int

const (
	PRICE_MODE_LIMIT              XunTou_OrderType = 11 //指定价/限价
	MARKET_SH_CONVERT_5_CANCEL    XunTou_OrderType = 42 //上海最优五档即时成交剩余撤销
	LATEST_PRICE                  XunTou_OrderType = 5  //最新价
	MARKET_SH_CONVERT_5_LIMIT     XunTou_OrderType = 43 //上海最优五档即时成交剩余转限价
	MARKET_PEER_PRICE_FIRST       XunTou_OrderType = 44 //深圳对手方最优价格
	MARKET_MINE_PRICE_FIRST       XunTou_OrderType = 45 //深圳本方最优价格
	MARKET_SZ_INSTBUSI_RESTCANCEL XunTou_OrderType = 46 //深圳即时成交剩余撤销
	MARKET_SZ_CONVERT_5_CANCEL    XunTou_OrderType = 47 //深圳最优五档即时成交剩余撤销
	MARKET_SZ_FULL_OR_CANCEL      XunTou_OrderType = 48 //深圳全额成交或撤销

)

type XunTou_Side int

const (
	STOCK_BUY                        XunTou_Side = 23
	STOCK_SELL                       XunTou_Side = 24
	CREDIT_BUY                       XunTou_Side = 23 //担保品买入
	CREDIT_SELL                      XunTou_Side = 24 //担保品卖出
	CREDIT_FIN_BUY                   XunTou_Side = 27 //融资买入
	CREDIT_SLO_SELL                  XunTou_Side = 28 //融券卖出
	CREDIT_BUY_SECU_REPAY            XunTou_Side = 29 //买券还券
	CREDIT_DIRECT_SECU_REPAY         XunTou_Side = 30 //直接还券
	CREDIT_SELL_SECU_REPAY           XunTou_Side = 31 //卖券还款
	CREDIT_DIRECT_CASH_REPAY         XunTou_Side = 32 //直接还款
	CREDIT_FIN_BUY_SPECIAL           XunTou_Side = 40 //专项融资买入
	CREDIT_SLO_SELL_SPECIAL          XunTou_Side = 41 //专项融券卖出
	CREDIT_BUY_SECU_REPAY_SPECIAL    XunTou_Side = 42 //专项买券还券
	CREDIT_DIRECT_SECU_REPAY_SPECIAL XunTou_Side = 43 //专项直接还券
	CREDIT_SELL_SECU_REPAY_SPECIAL   XunTou_Side = 44 //专项卖券还款
	CREDIT_DIRECT_CASH_REPAY_SPECIAL XunTou_Side = 45 //专项直接还款

)




//账号状态
type AccountStatus int

const (
	ACCOUNT_STATUS_INVALID AccountStatus = -1
	//正常
	ACCOUNT_STATUS_OK AccountStatus = 0
	//连接中
	ACCOUNT_STATUS_WAITING_LOGIN AccountStatus = 1
	//登陆中
	ACCOUNT_STATUSING AccountStatus = 2
	//失败
	ACCOUNT_STATUS_FAIL AccountStatus = 3
	//初始化中
	ACCOUNT_STATUS_INITING AccountStatus = 4
	//数据刷新校正中
	ACCOUNT_STATUS_CORRECTING AccountStatus = 5
	//收盘后
	ACCOUNT_STATUS_CLOSED AccountStatus = 6
	//穿透副链接断开
	ACCOUNT_STATUS_ASSIS_FAIL AccountStatus = 7
	//系统停用（总线使用-密码错误超限）
	ACCOUNT_STATUS_DISABLEBYSYS AccountStatus = 8
	//用户停用（总线使用）
	ACCOUNT_STATUS_DISABLEBYUSER AccountStatus = 9
)

// TradeClient implements the quickfix.Application interface
type WsApi struct {
	conn            net.Conn
	rspChan         chan *XunTouRsp
	orderChan       chan *XunTouReq
	AssetChan       chan *Asset
	CreditAssetChan chan *CreditAsset
	PositionChan    chan *Position
	Addr            string //"127.0.0.1:12000"
	rcvMsgTemp      []byte
	unpackChan      chan []byte
}

func NewWSApi(addr string) *WsApi {
	wsApi := new(WsApi)
	wsApi.rspChan = make(chan *XunTouRsp, 100000)
	wsApi.orderChan = make(chan *XunTouReq, 100000)
	wsApi.AssetChan = make(chan *Asset, 100000)
	wsApi.PositionChan = make(chan *Position, 100000)
	wsApi.CreditAssetChan = make(chan *CreditAsset, 100000)
	wsApi.rcvMsgTemp = make([]byte, 0)
	wsApi.unpackChan = make(chan []byte, 10000)
	wsApi.Addr = addr
	return wsApi
}

func (this *WsApi) wsConnect() error {
	conn, err := net.Dial("tcp", this.Addr)
	if err != nil {
		logger.Error("NewContext err ::%v", err)
		return err
	}

	this.conn = conn
	return nil
}

func (this *WsApi) wsStart() error {
	this.Unmashal()
	this.send()
	this.recv()
	return nil
}

func (this *WsApi) send() {
	go func() {
		tk := time.NewTicker(time.Second)
		for {
			select {
			case orderMsg := <-this.orderChan:
				logger.Info("ws receive new order :: %#v", orderMsg)
				msg, err := json.Marshal(orderMsg)
				if err != nil {
					logger.Error("order :: %#v, marshal fail :: %v", orderMsg, err)

				} else {

					_, err = this.conn.Write(PacketMsg(msg))
					if err != nil {
						logger.Crashf("send msg :: %s , fail :: %v", string(msg), err)

					} else {
						logger.Info("order send success ,loaclID :: %v,XunTouID ::%v", orderMsg.LocalID, orderMsg.XunTouEntrustNo)
					}
				}
				if err != nil {
					rsp := XunTouRsp{}
					rsp.KdbAcountName = orderMsg.KdbAcountName
					rsp.KdbLocalID_str = orderMsg.LocalID
					rsp.Order_status = 57
					this.rspChan <- &rsp
				}
			case <-tk.C:
			}
		}
	}()
}

func (this *WsApi) recv() {

	go func() {
		for {
			buf := make([]byte, 1024)
			n, err := this.conn.Read(buf)
			if err != nil && err != io.EOF {
				logger.Crashf("receive msg fail :: %v", err)
			}
			//logger.Debug("ws receive response :: %#v", string(buf[:n]))

			this.rcvMsgTemp = UnpackMsg(append(this.rcvMsgTemp, buf[:n]...), this.unpackChan)
		}
	}()
}

func (this *WsApi) Unmashal() {
	go func() {
		tk := time.NewTicker(time.Second)
		for {
			select {
			case unpackMsg := <-this.unpackChan:
				//rsp := &XunTouRsp{}
				t := &MsgType{} //O委托回报  A资金  P持仓	T状态
				err := json.Unmarshal(unpackMsg, t)
				if err != nil {
					logger.Error("msg : %s,unmashal err: %v", unpackMsg, err)
					continue
				}
				switch t.MsgType {
				case "O", "T":
					rsp := &XunTouRsp{}
					err := json.Unmarshal(unpackMsg, rsp)
					if err != nil {
						logger.Error("msg : %s,unmashal err: %v", unpackMsg, err)
						continue
					}
					this.rspChan <- rsp
				case "A":
					a := &Asset{}
					err := json.Unmarshal(unpackMsg, a)
					if err != nil {
						logger.Error("msg : %s,unmashal err: %v", unpackMsg, err)
						continue
					}
					this.AssetChan <- a
				case "CA":
					a := &CreditAsset{}
					err := json.Unmarshal(unpackMsg, a)
					if err != nil {
						logger.Error("msg : %s,unmashal err: %v", unpackMsg, err)
						continue
					}
					this.CreditAssetChan <- a
				case "P":
					p := &Position{}
					err := json.Unmarshal(unpackMsg, p)
					if err != nil {
						logger.Error("msg : %s,unmashal err: %v", unpackMsg, err)
						continue
					}
					this.PositionChan <- p

				}
				//logger.Info("response order :: %#v", rsp)


			case <-tk.C:

			}
		}

	}()
}
func accountIsError(i int) (bl bool, msg string) {

	switch AccountStatus(i) {
	case ACCOUNT_STATUS_INVALID:
		bl = true
		msg = "无效"
		//无效
	case ACCOUNT_STATUS_OK:
		bl = false
		msg = "正常"
		//正常
	case ACCOUNT_STATUS_WAITING_LOGIN:
		bl = true
		msg = "连接中"
		//连接中
	case ACCOUNT_STATUSING:
		bl = true
		msg = "登陆中"
		//登陆中
	case ACCOUNT_STATUS_FAIL:
		bl = true
		msg = "失败"
		//失败
	case ACCOUNT_STATUS_INITING:
		bl = true
		msg = "初始化中"
		//初始化中
	case ACCOUNT_STATUS_CORRECTING:
		bl = true
		msg = "数据刷新校正中"
		//数据刷新校正中
	case ACCOUNT_STATUS_CLOSED:
		bl = true
		msg = "收盘后"
		//收盘后
	case ACCOUNT_STATUS_ASSIS_FAIL:
		bl = true
		msg = "穿透副链接断开"
		//穿透副链接断开
	case ACCOUNT_STATUS_DISABLEBYSYS:
		bl = true
		msg = "系统停用"
		//系统停用（总线使用-密码错误超限）
	case ACCOUNT_STATUS_DISABLEBYUSER:
		bl = true
		msg = "用户停用"
		//用户停用（总线使用）
	default:
		bl = false
		msg = "Unkunow Account Status ::" + strconv.Itoa(i)
	}
	return
}
