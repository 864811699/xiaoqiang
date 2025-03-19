package server

import (
	logger "github.com/alecthomas/log4go"
	"market/pkg/activemq"
	"market/pkg/kdb"
)

type Server struct {
	kdbHandle    *kdb.KdbHandle
	activeHandle *activemq.Active
	inCh         <-chan []byte
	outCh        chan<- *kdb.Market
}

func NewServer(k *kdb.KdbHandle, a *activemq.Active) *Server {
	return &Server{
		kdbHandle:    k,
		activeHandle: a,
		outCh:        k.GetMarketCh(),
		inCh:         a.GetMsgChan(),
	}
}

func (s *Server) Start() {

	if err := s.kdbHandle.Start(); err != nil {
		logger.Crash("kdb start fail,err:%v", err)
	}
	logger.Info("kdb success!!!")

	if err:=s.activeHandle.Start();err!=nil{
		logger.Crash("activeMQ start fail,err:%v", err)
	}
	logger.Info("activeMQ success!!!")


	for msg := range s.inCh {
		if market, err := MarketActiveToKdb(msg); err != nil {
			logger.Warn("data error,%v", err)
			continue
		} else {
			logger.Debug("data zhongxin,%v", market)
			s.outCh <- market
		}

	}
}

func (s *Server) Stop() {
	s.activeHandle.Stop()
	s.kdbHandle.Stop()
}
