package kdb

import (
	"context"
	logger "github.com/alecthomas/log4go"
	kdb "github.com/nsxy/kdbgo"
	"github.com/pkg/errors"
	"time"
)

type KdbHandle struct {
	ip       string
	port     int
	auth     string
	conn     *kdb.KDBConn
	tab      string
	fun      string
	marketCh chan *Market
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewKdbHandle(ip, auth, tab, fun string, port int) *KdbHandle {
	ctx, cancel := context.WithCancel(context.Background())
	return &KdbHandle{
		ip:       ip,
		port:     port,
		auth:     auth,
		tab:      tab,
		fun:      fun,
		ctx:      ctx,
		cancel:   cancel,
		marketCh: make(chan *Market, 100000),
	}
}

func (tb *KdbHandle) Connect() error {
	conn, err := kdb.DialKDBTimeout(tb.ip, tb.port, tb.auth, time.Second*10)
	if err != nil {
		return errors.Wrapf(err, "connect kdb[%s:%d] fail", tb.ip, tb.port)
	}
	tb.conn = conn
	return nil
}

func (tb *KdbHandle) subscribe(sym []string) error {
	var symK *kdb.K
	if symNum := len(sym); symNum == 0 {
		symK = kdb.Symbol("")
	} else {
		symK = kdb.SymbolV(sym)
	}

	err := tb.conn.AsyncCall(".u.sub", kdb.Symbol(tb.tab), symK)
	if err != nil {
		return errors.Wrapf(err, "sub kdb[%s:%d] fail", tb.ip, tb.port)
	}
	return nil
}

func (tb *KdbHandle) send2Tab(data *kdb.K) error {
	if err := tb.conn.AsyncCall(tb.fun, kdb.Symbol(tb.tab), data); err != nil {
		return errors.Wrapf(err, "send data: %v to kdb table: %v with Func: %v", data, tb.tab, tb.fun)
	}
	return nil
}

func (tb *KdbHandle) readFromKdb() error {
	for {
		res, _, err := tb.conn.ReadMessage()
		if err != nil {
			return errors.Wrapf(err, "read KDB message error:", )
		}
		switch res.Data.(type) {
		case string:
			logger.Info(res.Data)
			if msg := res.Data.(string); msg == "\"heartbeat\"" {
				continue
			}

		case []*kdb.K:
			dataArr := res.Data.([]*kdb.K)
			funcName := dataArr[0].Data.(string)
			if funcName != "upd" {
				_ = logger.Error("function name is not upd, func_name: %v", funcName)
				continue
			}
			tableName := dataArr[1].Data.(string)
			if tableName != tb.tab {
				continue
			}

			for k, _ := range dataArr {
				if dataArr[k].Type == kdb.XT {
					//logger.Info("recv:%v",dataArr[k])
				}

			}
		}

	}
}

func (tb *KdbHandle) GetMarketCh()  chan<- *Market{
	return tb.marketCh
}

func (tb *KdbHandle) dealMarket() {
	for {
		select {
		case market, ok := <-tb.marketCh:
			if !ok {
				return
			}
			marketK:=parseMarketToKdb(*market)
			if err := tb.send2Tab(marketK); err != nil {
				logger.Warn("send to kdb fail,err:%v", err)
				return
			}
		case _ = <-tb.ctx.Done():
			return
		}
	}
}

func (tb *KdbHandle) Start() error {
	if err := tb.Connect(); err != nil {
		return errors.Wrapf(err, "kdb start fail")
	}

	go tb.dealMarket()

	return nil
}

func (tb *KdbHandle) Stop() {
	tb.conn.Close()
	tb.cancel()

}
