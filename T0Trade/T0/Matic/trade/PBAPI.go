package trade

import (
	"bufio"
	"fmt"
	logger "github.com/alecthomas/log4go"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os"
	"time"
)

type XunTouReq struct {
	IsCancel    bool
	KdbAccount  string
	Side        string
	PriceType   string
	Price       float64
	SecurityID  string
	Qty         int32
	PackID      string
	XunTouID    string
	Exchange  string

}



type XunTou_Side int

const (
	STOCK_BUY                XunTou_Side = 23
	STOCK_SELL               XunTou_Side = 24
	CREDIT_BUY               XunTou_Side = 23 //担保品买入
	CREDIT_SELL              XunTou_Side = 24 //担保品卖出
	CREDIT_FIN_BUY           XunTou_Side = 27 //融资买入
	CREDIT_SLO_SELL          XunTou_Side = 28 //融券卖出
	CREDIT_BUY_SECU_REPAY    XunTou_Side = 29 //买券还券
	CREDIT_DIRECT_SECU_REPAY XunTou_Side = 30 //直接还券

	CREDIT_SELL_SECU_REPAY           XunTou_Side = 31 //卖券还款
	CREDIT_DIRECT_CASH_REPAY         XunTou_Side = 32 //直接还款
	CREDIT_FIN_BUY_SPECIAL           XunTou_Side = 40 //专项融资买入
	CREDIT_SLO_SELL_SPECIAL          XunTou_Side = 41 //专项融券卖出
	CREDIT_BUY_SECU_REPAY_SPECIAL    XunTou_Side = 42 //专项买券还券
	CREDIT_DIRECT_SECU_REPAY_SPECIAL XunTou_Side = 43 //专项直接还券
	CREDIT_SELL_SECU_REPAY_SPECIAL   XunTou_Side = 44 //专项卖券还款
	CREDIT_DIRECT_CASH_REPAY_SPECIAL XunTou_Side = 45 //专项直接还款

)

// TradeClient implements the quickfix.Application interface
type PBApi struct {
	accounts    []AccountInfo
	responseFile map[string]string
	errorFile   map[string]string
	orderChan   chan *XunTouReq
	RspStrCh      chan string
	ErrStrCh	 chan string
	accountChan map[string]chan string
	ordFileMap  map[string]string
	tempFileMap map[string]string
}

type AccountInfo struct {
	Account     string `json:"account"`
	AccountType string `json:"account_type"`
	Buy string `json:"buy"`
	Sell string `json:"sell"`
}

func NewPBAPI(accounts map[string]AccountInfo, orderDir, tempDir, resonseDir string) *PBApi {

	pbApi := &PBApi{
		orderChan:   make(chan *XunTouReq, 10000),
		RspStrCh:      make(chan string, 10000),
		ErrStrCh:      make(chan string, 10000),
		accountChan: make(map[string]chan string),
		ordFileMap:  make(map[string]string, 100),
		tempFileMap: make(map[string]string, 100),
		responseFile:make(map[string]string,100),
		errorFile: make(map[string]string,100),
	}

	os.MkdirAll(tempDir, 0644)
	os.MkdirAll(orderDir, 0644)
	os.MkdirAll(resonseDir, 0644)

	day:=time.Now().Format("20060102")
	for kdbAccount, account := range accounts {
		//signal.%s_%s.%d.%s.txt   signal.资金账号_账号类型.下单序号.YYYYMMDD.txt
		orderFile := fmt.Sprintf("%s/signal.%s_%s", orderDir, account.Account, account.AccountType)
		tempFile := fmt.Sprintf("%s/signal.%s_%s", tempDir, account.Account, account.AccountType)

		pbApi.ordFileMap[kdbAccount] = orderFile
		pbApi.tempFileMap[kdbAccount] = tempFile
		pbApi.responseFile[kdbAccount] = fmt.Sprintf("%s/order.%s_%s.%s.csv",resonseDir,account.Account,account.AccountType,day)
		pbApi.errorFile[kdbAccount] = fmt.Sprintf("%s/order_error.%s_%s.%s.csv",resonseDir,account.Account,account.AccountType,day)
		pbApi.accountChan[kdbAccount] = make(chan string, 10000)

	}

	logger.Info("Accounts init success!!!")
	logger.Info("accountType : 2普通, 3信用")
	logger.Info("side: 23买\t	24卖\t	23担保品买入\t	24担保品卖出\t	27融资买入\t	28 融券卖出\t	29 买券还券\t	30 直接还券\t")

	return pbApi
}

func (this *PBApi) start(isFull bool) {
	//if isFull {
	//	this.readFullCsv()
	//} else {
	for k, _ := range this.responseFile {
		responseFile := this.responseFile[k]
		go readCsv(responseFile,this.RspStrCh)
	}
	for k, _ := range this.errorFile {
		responseFile := this.errorFile[k]
		go readCsv(responseFile,this.ErrStrCh)
	}
	//}

	this.sendOrder()
	this.marshalOrder()

}



func (this *PBApi) marshalOrder() {

	go func() {
		for order := range this.orderChan {
			msg := ""
			if order.IsCancel {
				msg = fmt.Sprintf("cancel_sys,%s,%s\n", order.XunTouID,order.Exchange)
			} else {
				//1,        23,    3,       3,   SH601818, 300, strategy
				//entrustno side priceType price stockcode qty stratege_name
				msg = fmt.Sprintf("%s,%s,%s,%.2f,%s,%d,%s",
					order.PackID, order.Side, order.PriceType, order.Price, order.SecurityID, order.Qty, order.KdbAccount)

			}
			this.accountChan[order.KdbAccount] <- msg
		}
	}()

}

func (this *PBApi) sendOrder() {

	bufSize := 10
	nowDay := time.Now().Format("20060102")
	for account, accountCh := range this.accountChan {
		account, accountCh := account, accountCh
		order_tick := time.NewTicker(100 * time.Millisecond)
		order_cache := make([]string, 0, bufSize)
		go func(string, chan string) {
			logger.Debug("account:%s ,chan : %v ", account, accountCh)
			tempFile := this.tempFileMap[account]
			ordFile := this.ordFileMap[account]

			for {
				select {
				case msg := <-accountCh:

					logger.Info("send order,chan [%s] ,order msg :%v ", account, msg)

					order_cache = append(order_cache, msg)
					if len(order_cache) >= bufSize {
						//signal.%s_%s.%d.%s.txt   signal.资金账号_账号类型.下单序号.YYYYMMDD.txt
						//orderFile := fmt.Sprintf("./order/signal.%s_%s", account.Account, account.AccountType)
						//tempFile := fmt.Sprintf("./temp/signal.%s_%s", account.Account, account.AccountType)
						now := time.Now()

						currentTempFile := fmt.Sprintf("%s.%d.%s.txt", tempFile, now.UnixNano(), nowDay)
						currentOrdFile := fmt.Sprintf("%s.%d.%s.txt", ordFile, now.UnixNano(), nowDay)
						f, err := os.OpenFile(currentTempFile, os.O_APPEND|os.O_CREATE, 0777)
						if err != nil {
							logger.Warn("create temp file[%s] error: %v\n", currentTempFile, err)
							continue
						}

						w := bufio.NewWriter(f)
						for _, ordStr := range order_cache {
							fmt.Fprintln(w, ordStr)
						}
						w.Flush()
						f.Close()

						_, err = os.Stat(currentOrdFile)
						if os.IsNotExist(err) {
							//移动temp文件到扫单目录
							os.Rename(currentTempFile, currentOrdFile)
						}

						order_cache = make([]string, 0, bufSize)

					}
				case <-order_tick.C:
					if len(order_cache) > 0 {
						now := time.Now()
						currentTempFile := fmt.Sprintf("%s.%d.%s.txt", tempFile, now.UnixNano(), nowDay)
						currentOrdFile := fmt.Sprintf("%s.%d.%s.txt", ordFile, now.UnixNano(), nowDay)

						f, err := os.OpenFile(currentTempFile, os.O_APPEND|os.O_CREATE, 0777)
						if err != nil {
							logger.Warn("create temp file[%s] error: %v\n", currentTempFile, err)
							continue
						}

						w := bufio.NewWriter(f)
						for _, ordStr := range order_cache {
							fmt.Fprintln(w, ordStr)
						}
						w.Flush()
						f.Close()

						_, err = os.Stat(currentOrdFile)
						if os.IsNotExist(err) {
							//移动temp文件到扫单目录
							os.Rename(currentTempFile, currentOrdFile)
						}

						order_cache = make([]string, 0, bufSize)

					}
				}
			}
		}(account, accountCh)
	}

}

func  readCsv(responseFile string,ch chan string) {
			tk:=time.NewTicker(100*time.Millisecond)
			isExist(responseFile)

			f, err := os.Open(responseFile)
			if err != nil {
				logger.Crashf("open file [%s] fail,err:%v", responseFile, err)
			}
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := GBKToUTF8(scanner.Text())
				ch <- line
			}

			logger.Info("notify %s beggin !!!",responseFile)

			for  {
				select {
				case <-tk.C:
					scanner := bufio.NewScanner(f)
					for scanner.Scan() {
						line := GBKToUTF8(scanner.Text())
						ch <- line
					}
				}

			}




}

//func (this *PBApi) readFullCsv() {
//	for k, _ := range this.resonseFile {
//		responseFile := this.resonseFile[k]
//		go func(responseFile string) {
//			tk := time.NewTicker(time.Second)
//			logger.Info("Scan Full CSV Beggin !!!")
//			for range tk.C {
//				f, err := os.Open(responseFile)
//				if err != nil {
//					logger.Crashf("open file [%s] fail,err:%v", responseFile, err)
//					continue
//				}
//				scanner := bufio.NewScanner(f)
//				for scanner.Scan() {
//					line := GBKToUTF8(scanner.Text())
//					this.RspStr <- line
//				}
//				f.Close()
//			}
//
//		}(responseFile)
//	}
//
//}

func GBKToUTF8(s string) string {
	utf8s, _ := simplifiedchinese.GB18030.NewDecoder().String(s)
	return utf8s
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

func isExist(f string) {
	_, err := os.Stat(f)
	if os.IsNotExist(err) {
		f1, err := os.Create(f)
		if err == nil {
			f1.Close()
		}
	}
}
