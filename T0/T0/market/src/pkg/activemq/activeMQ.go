package activemq

import (
	"context"
	logger "github.com/alecthomas/log4go"
	"github.com/go-stomp/stomp/v3"
	"github.com/pkg/errors"
	"sync"
	"time"
)

type Active struct {
	conn   *stomp.Conn
	sub    []*stomp.Subscription
	addr   string
	topic  []string
	ch     chan  stomp.Message
	ctx    context.Context
	cancel context.CancelFunc
}

func NewActive(addr string, topic []string) *Active {
	ctx, cancel := context.WithCancel(context.Background())
	return &Active{
		conn:   new(stomp.Conn),
		addr:   addr,
		topic:  topic,
		ch:     make(chan stomp.Message, 100000),
		sub:    make([]*stomp.Subscription,0,10),
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
	for _,topic:=range a.topic {
		sub, err := a.conn.Subscribe(topic, stomp.AckAuto)
		if err != nil {
			return err
		}
		a.sub=append(a.sub,sub)
	}
	return nil
}

func (a *Active) reciveMarket() {
	var wg sync.WaitGroup
	wg.Add(2)
	for k,_:=range a.sub {
		sub:=a.sub[k]
		go func(sub *stomp.Subscription,wg *sync.WaitGroup ) {
			wg.Done()
			for {
				select {
				case msg, ok := <-sub.C:
					if !ok {
						logger.Warn("activeMQ receive market fail,ok:%v",ok)
						close(a.ch)
						return
					}
					if msg.Err!=nil{
						logger.Warn("activeMQ receive market fail,err:%v",msg.Err)
						continue
					}
					//CompressType:=msg.Header.Get("CompressType")
					msg_new:=*msg
					a.ch <- msg_new
				case <-a.ctx.Done():
					close(a.ch)
					return
				}
			}
		}(sub,&wg)
	}
	wg.Wait()

}

func (a *Active) Stop() {
	for k,_:=range a.sub {
		a.sub[k].Unsubscribe()
	}

	a.cancel()
	close(a.ch)
}

func (a *Active) GetMsgChan() <-chan stomp.Message {
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
