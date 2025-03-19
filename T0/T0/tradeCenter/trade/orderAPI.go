package trade

import (
	"bufio"
	"fmt"
	kdbtdx "github.com/864811699/T0kdb"
	logger "github.com/alecthomas/log4go"
	"github.com/fsnotify/fsnotify"
	"os"
	"strconv"
	"strings"
	"time"
)

//type OrdData struct {
//	// Broker   string  `json:"Broker" binding:"required"`
//	AlgoName  string            `json:"AlgoName" binding:"required" `
//	Params    map[string]string `json:"Params" binding:"required" `
//	Side      string            `json:"Side" binding:"required" `
//	OrderType string            `json:"OrderType" binding:"required"`
//	CancelID string
//	Entrust
//}

//type Entrust struct {
//	Sym         string `json:"user" binding:"required"`
//	Qid         string `json:"qid" binding:"-" `
//	Accountname string `json:"Account" binding:"required" `
//	Time        time.Time
//	Entrustno   int32
//	Stockcode   string  `json:"Ticker" binding:"required" `
//	Askprice    float64 `json:"Price" binding:"-" `
//	Askvol      int32   `json:"Qty" binding:"required"`
//	Bidprice    float64
//	Bidvol      int32
//	Withdraw    int32
//	Status      int32
//}

//type RejectedOrder struct {
//	Account   string  `json:"account"`
//	User      string  `json:"tradeUser"`
//	Date      string  `json:"date"`
//	Time      string  `json:"time"`
//	LocalID   int32   `json:"LocalID"`
//	Ticker    string  `json:"ticker"`
//	Qty       int32   `json:"qty"`
//	Price     float64 `json:"price"`
//	Side      string  `json:"side"`
//	IsQFII    bool    `json:"isQfii"`
//	AlgoName  string  `json:"AlgoName"`
//	ErrorText string  `json:"ErrorText"`
//}
//
//type rejMsg struct {
//	text string
//	OrdData
//}

type DataInfo struct {
	WTFile      string          `json:"wt_file"`
	EntrustFile string          `json:"entrust_file"`
	Infos       map[string]Info `json:"infos"`
}

type Info struct {
	Buy         string `json:"buy"`
	Sell        string `json:"sell"`
	FundAccount string `json:"fundAccount"`
}

var StatusMapInt32 = map[string]int32{
	"0": 1,
	"1": 2,
	"2": 4,
	"4": 5,
	"6": 3,
	"8": 6,
}

const (
	MARKET_SHANGHAI_code = "XSHG"
	MARKET_SHENZHEN_code = "XSHE"
)

const (
	LIMIT_Price  = "2"
	MARKET_Price = "1"
)

// TradeClient implements the FileOrder
//TODO
type FileOrder struct {
	entrust_path  string //委托回报文件
	write_path    string //下单文件
	asset_path    string
	position_path string

	orderChan    chan string
	updateChan   chan string
	assetChan    chan kdbtdx.Asset
	positionChan chan kdbtdx.Position
}

func NewFileOrderApi(write_path, entrust_path string) *FileOrder {
	t := time.Now().Format("20060102")
	orderApi := new(FileOrder)

	orderApi.write_path = fmt.Sprintf(write_path, t)
	orderApi.entrust_path = fmt.Sprintf(entrust_path, t)

	orderApi.orderChan = make(chan string, 10000)
	orderApi.updateChan = make(chan string, 10000)
	orderApi.assetChan = make(chan kdbtdx.Asset, 10000)
	orderApi.positionChan = make(chan kdbtdx.Position, 10000)

	return orderApi
}

func (this *FileOrder) OrderApiStart() {
	go this.fileNotify()
	go this.write_order()
	go this.readAssetFile()
	go this.readPositionFile()
}

func (this *FileOrder) write_order() error {
	timeout := time.NewTicker(100 * time.Millisecond)

	go func() {
		logger.Info("********写单文件模块启动*******")
		f, err := os.OpenFile(this.write_path, os.O_APPEND, 222)
		if err != nil {
			logger.Info("open write file :%v  fail ,err:%v", this.write_path, err)
			panic(err)
		}
		order_cache := make([]string, 0, 10)
		for {
			select {
			case order_str := <-this.orderChan:

				order_cache = append(order_cache, order_str)
				if len(order_cache) >= 10 {

					w := bufio.NewWriter(f)

					for _, reqStr := range order_cache {
						fmt.Fprintln(w, reqStr)
						w.Flush()
					}
					w.Flush()

					order_cache = make([]string, 0, 10)
				}
			case <-timeout.C:

				if len(order_cache) > 0 {

					w := bufio.NewWriter(f)
					for _, reqStr := range order_cache {
						fmt.Fprintln(w, reqStr)
						w.Flush()
					}
					w.Flush()

					order_cache = make([]string, 0, 10)
				}

			}
		}
	}()

	return nil
}

func (this *FileOrder) fileNotify() {
	logger.Info("beggin Notify file :: %v", this.entrust_path)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic("监控创建失败,")

	}
	err = watcher.Add(this.entrust_path)
	if err != nil {
		logger.Info("file : %v", this.entrust_path)
		panic("文件监控添加失败")
	}

	f, _ := os.Open(this.entrust_path)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		this.updateChan <- line
	}

	timeout := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case ev := <-watcher.Events:
			if ev.Op&fsnotify.Write == fsnotify.Write {
				scanner := bufio.NewScanner(f)
				for scanner.Scan() {

					line := scanner.Text()
					logger.Debug("read msg :: %v", line)

					this.updateChan <- line
				}
			}
		case <-timeout.C:
			continue
		}
	}
}

func tickNotify(file string, msgChan chan string) {
	logger.Info("beggin Notify file :: %v", file)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic("监控创建失败,")
	}
	err = watcher.Add(file)
	if err != nil {
		fmt.Println("file : ", file)
		panic("文件监控添加失败")
	}

	f, _ := os.Open(file)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		msgChan <- line
	}

	tk := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-tk.C:
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				if len(line) < 10 {
					continue
				}
				logger.Debug("read msg :: %v", line)

				msgChan <- line
			}


		}
	}
}

func (this *FileOrder) readAssetFile() {
	logger.Info("beggin Notify file :: %v", this.asset_path)
	tk := time.NewTicker(5 * time.Second)
	for range tk.C {
		f, _ := os.Open(this.asset_path)
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimSuffix(line, "_wjd")
			//交易账号|可用资金
			fields := strings.Split(line, "|")
			if len(fields) < 2 {
				logger.Warn("asset line error [%v]", line)
				continue
			}
			asset := kdbtdx.Asset{
				Account:   fields[0],
				Available: 0,
			}

			f, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				logger.Error("asset parse fail,text:%v", line)
			} else {
				asset.Available = f
			}

			this.assetChan <- asset

		}
	}

}

func (this *FileOrder) readPositionFile() {
	logger.Info("beggin Notify file :: %v", this.position_path)
	tk := time.NewTicker(5 * time.Second)
	for range tk.C {
		f, _ := os.Open(this.position_path)
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimSuffix(line, "_wjd")
			//交易账号|代码|市场|模式|余额|可用
			fields := strings.Split(line, "|")
			if len(fields) < 6 {
				logger.Warn("position line error [%v]", line)
				continue
			}
			position := kdbtdx.Position{
				Account:   fields[0],
				Stockcode: fields[1],
			}

			gpye, err := strconv.Atoi(fields[4])
			if err != nil {
				logger.Error("position parse fail,text:%v", line)
			} else {
				position.Gpye = int32(gpye)
			}

			kyye, err := strconv.Atoi(fields[5])
			if err != nil {
				logger.Error("position parse fail,text:%v", line)
			} else {
				position.Kyye = int32(kyye)
			}

			this.positionChan<-position

		}
	}
}
