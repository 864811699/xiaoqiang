package activemq

import (
	"context"
	logger "github.com/alecthomas/log4go"
	"github.com/go-stomp/stomp/v3"
	"github.com/pkg/errors"
	"time"
)

type Active struct {
	conn   *stomp.Conn
	sub    *stomp.Subscription
	addr   string
	topic  string
	ch     chan []byte
	ctx    context.Context
	cancel context.CancelFunc
}

func NewActive(addr, topic string) *Active {
	ctx, cancel := context.WithCancel(context.Background())
	return &Active{
		conn:   new(stomp.Conn),
		addr:   addr,
		topic:  topic,
		ch:     make(chan []byte, 100000),
		sub:    new(stomp.Subscription),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (a *Active) connect() error {
	conn, err := stomp.Dial("tcp", a.addr,stomp.ConnOpt.HeartBeat(10*time.Second,86400*time.Second),stomp.ConnOpt.RcvReceiptTimeout(10800*time.Second))
	if err != nil {
		return errors.Wrapf(err, "Dial[%s] fail", a.addr)
	}
	a.conn = conn
	return nil
}

func (a *Active) subMarket() error {
	sub, err := a.conn.Subscribe(a.topic, stomp.AckAuto)
	if err != nil {
		return err
	}
	a.sub = sub
	return nil
}

func (a *Active) reciveMarket() {

	for {
		select {
		case msg, ok := <-a.sub.C:
			if !ok {
				logger.Warn("activeMQ receive market fail,ok:%v",ok)
				close(a.ch)
				return
			}
			if msg.Err!=nil{
				logger.Warn("activeMQ receive market fail,err:%v",msg.Err)
				continue
			}
			a.ch <- msg.Body
		case <-a.ctx.Done():
			close(a.ch)
			return
		}
	}

}

func (a *Active) Stop() {
	a.sub.Unsubscribe()
	a.cancel()
	close(a.ch)
}

func (a *Active) GetMsgChan() <-chan []byte {
	return a.ch
}

func (a *Active) Start() error{
	if err:=a.connect();err != nil {
		return errors.Wrapf(err,"active start fail")
	}
	if err:=a.subMarket();err!=nil{
		return errors.Wrapf(err,"active sub fail")
	}
	go a.reciveMarket()

	return nil
}
