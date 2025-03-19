package trade

import (
	"bufio"
	"fmt"
	"github.com/LindsayBradford/go-dbf/godbf"
	logger "github.com/alecthomas/log4go"
	"os"
	"strconv"
	"strings"
	"time"
)

type ZhongXinReq struct {
	ReqType     Req_Type
	AccountType string
	Account     string
	SecurityId  string
	Side        ZhongXin_Side
	Price       float64
	Qty         int32
	OrderType   OrderType

	PackID string

	CancelID string
}

type ZhongXinRsp struct {
	PackID    string
	CancelID  string
	Status    int32
	ErrMsg    string
	DealQty   int
	AvgPx     float64
	CancelQty int
	Note      string
}

//s1	资金账户	暂时和acct一样，预留字段
//s2	币种	参见数据字典章节
//s3	当前余额
//s4	可用余额
//s5	当前保证金/已用保证金	期货账户是当前保证金，期权账户是已用保证金
//s6	冻结保证金/可用保证金	期货账户是冻结保证金，期权账户是可用保证金
type Asset struct {
	Account string

}

type Position struct {

}

type OrderType string

const (
	LIMIT  OrderType = "0"
	MARKET OrderType = "U"
)

var OrderType_map = map[int32]OrderType{
	0: MARKET,
	1: LIMIT,
}

type Req_Type string

const (
	REQ_ORDER  Req_Type = "O"
	REQ_CANCEL Req_Type = "C"
)

var ZhongXin_Side_Map = map[int32]ZhongXin_Side{
	0: STOCK_BUY,
	1: STOCK_SELL,
	4: CREDIT_FIN_BUY,
	5: CREDIT_SLO_SELL,
}

type ZhongXin_Side string

const (
	STOCK_BUY       ZhongXin_Side = "1"
	STOCK_SELL      ZhongXin_Side = "2"
	CREDIT_FIN_BUY  ZhongXin_Side = "A" //融资买入
	CREDIT_SLO_SELL ZhongXin_Side = "B" //融券卖出
)

// TradeClient implements the quickfix.Application interface
type PBApi struct {
	resonseFile string
	assetFile   string
	OrderChan   chan *ZhongXinReq
	DBFChan     chan *ZhongXinRsp
	orderDir    string
	tempDir     string
}

func NewPBAPI(orderDir, tempDir, resonseFile, assetFile string) *PBApi {
	return &PBApi{
		OrderChan:   make(chan *ZhongXinReq, 10000),
		DBFChan:     make(chan *ZhongXinRsp, 10000),
		orderDir:    orderDir,
		tempDir:     tempDir,
		resonseFile: resonseFile,
		assetFile:   assetFile,
	}
}

func (this *PBApi) start() {
	this.InintAccounts()
	go this.readCsv()
	go this.sendOrder()

}

func (this *PBApi) InintAccounts() {
	os.RemoveAll(this.tempDir)
	os.MkdirAll(this.tempDir, 777)
	os.MkdirAll(this.orderDir, 777)
	logger.Info("Accounts init success!!!")
}

func (this *PBApi) sendOrder() {

	for order := range this.OrderChan {

		order_file := fmt.Sprintf("%s/%s.csv", this.orderDir, order.PackID)
		temp_order_file := fmt.Sprintf("%s/%s.csv", this.tempDir, order.PackID)
		f, err := os.Create(temp_order_file)
		if err != nil {
			logger.Warn("create temp file[%s] error: %v\n", temp_order_file, err)
			continue
		}

		w := bufio.NewWriter(f)
		reqStr := ""
		if order.ReqType == REQ_ORDER {

			reqStr = fmt.Sprintf("%s,%s,%s,%s,%d,%s,%.4f,%s", order.ReqType, order.AccountType, order.Account, order.SecurityId, order.Qty, order.Side, order.Price, order.OrderType)

		} else {

			reqStr = fmt.Sprintf("%s,%s,%s,%s", order.ReqType, order.AccountType, order.Account, order.CancelID)
		}
		fmt.Fprintln(w, reqStr)
		w.Flush()
		f.Close()
		os.Rename(temp_order_file, order_file)
		os.Remove(temp_order_file)

	}

}

func (this *PBApi) readCsv() {
	defer func() {
		err := recover()
		if err != nil {
			logger.Info("出现panic异常 :: %v", err)
		}
		go this.readCsv()
	}()

	order_ticker := time.NewTicker(500 * time.Millisecond)
	rows := 0
	for {
		select {
		case <-order_ticker.C:

			file, err := godbf.NewFromFile(this.resonseFile, "GBK")
			if err != nil {
				logger.Error(err)
				return
			}

			//file.Rows()
			row := file.NumberOfRecords()

			for i := rows; i < row; i++ {
				logger.Info("scan msg : %v", file.GetRowAsSlice(i))
				//2020-07-31 判断日期
				//if tm, err := file.FieldValueByName(i, "ORD_TIME"); err == nil {
				//	if ! strings.Contains(tm,now){
				//		continue
				//	}
				//
				//}
				if localid, err := file.FieldValueByName(i, "CLIENT_ID"); localid != "" && err == nil {
					if strings.Contains(localid, "Cancel") || !strings.Contains(localid, day) {
						continue
					}

					//err_msg 有数据时 ,ORD_NO 必定为空, 撤废/废单(status=已报)
					// CLIENT_ID 为空标识,订单回报先于委托确认回报
					//n++
					entrust := new(ZhongXinRsp)
					entrust.PackID = localid

					if entrust_no, err := file.FieldValueByName(i, "ORD_NO"); err == nil {
						entrust.CancelID = entrust_no
					}

					if entrust_status, err := file.FieldValueByName(i, "ORD_STATUS"); err == nil {
						entrust.Status = StatusMapInt32[entrust_status]

					}

					if err_msg, err := file.FieldValueByName(i, "ERR_MSG"); err_msg != "" && err == nil {
						entrust.ErrMsg = err_msg
					}

					if filled_qty, err := file.FieldValueByName(i, "FILLED_QTY"); err == nil {
						if deal_amount, err := strconv.Atoi(filled_qty); err == nil {
							entrust.DealQty = deal_amount

						}
					}

					if avg_px, err := file.FieldValueByName(i, "AVG_PX"); err == nil {
						if deal_price, err := strconv.ParseFloat(avg_px, 64); err == nil {
							entrust.AvgPx = deal_price
						}
					}

					if cxl_qty, err := file.FieldValueByName(i, "CXL_QTY"); err == nil {
						if cancel_amount, err := strconv.Atoi(cxl_qty); err == nil {
							entrust.CancelQty = cancel_amount
						}
					}

					if entrust_errMsg, err := file.FieldValueByName(i, "ERR_MSG"); err == nil {
						entrust.Note = entrust_errMsg
					}
					if entrust.ErrMsg != "" && entrust.CancelID == "" {
						entrust.Status = 6
					}
					this.DBFChan <- entrust
				}

			}
			rows = row

		}

	}
}

func (this *PBApi)loadAsset()  {
	defer func() {
		err := recover()
		if err != nil {
			logger.Info("出现panic异常 :: %v", err)
		}
		go this.loadAsset()
	}()

	order_ticker := time.NewTicker(500 * time.Millisecond)
	rows := 0
	for {
		select {
		case <-order_ticker.C:

			file, err := godbf.NewFromFile(this.assetFile, "GBK")
			if err != nil {
				logger.Error(err)
				return
			}

			//file.Rows()
			row := file.NumberOfRecords()

			for i := rows; i < row; i++ {
				logger.Info("scan msg : %v", file.GetRowAsSlice(i))

				if localid, err := file.FieldValueByName(i, "CLIENT_ID"); localid != "" && err == nil {
					if strings.Contains(localid, "Cancel") || !strings.Contains(localid, day) {
						continue
					}


				}

			}
			rows = row

		}

	}
}

var StatusMapInt32 = map[string]int32{
	"0": 1,
	"1": 2,
	"2": 4,
	"3": 5,
	"4": 5,
	"5": 6,
	"6": 6,
}
