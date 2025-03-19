package trade

import (
	"fmt"
	logger "github.com/alecthomas/log4go"
	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
)


type rabbitMqManage struct {
	conn *amqp.Connection
}

type HandleParentOrder interface {
	ReceiveParentOrder (order *ParentOrder)
}

func NewRabbitMq() *rabbitMqManage {
	return &rabbitMqManage{}
}


func (r *rabbitMqManage) Init (rabbitUrl string) error {
	if r.conn==nil{
		newConn,err:=amqp.Dial(rabbitUrl)
		if err!=nil{
			return err
		}
		r.conn=newConn
	}
	return nil
}

func (r *rabbitMqManage) GetConn() *amqp.Connection{
	return r.conn
}

func (r *rabbitMqManage)Sub(handler HandleParentOrder,queueName string,algoNames []string) error {
	ch,err:=r.GetConn().Channel()
	if err!=nil{
		return err
	}

	err=ch.ExchangeDeclare("algo","topic",false,false,false,false,nil)
	if err != nil {
		return err
	}
	var q amqp.Queue
	q, err = ch.QueueDeclare(queueName, false, false, true, false, nil)
	if err != nil {
		return err
	}
	for _, algoName := range algoNames {
		err = ch.QueueBind(q.Name, fmt.Sprintf("algo.%s", algoName), "algo", false, nil)
		if err != nil {
			return err
		}
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			order := &ParentOrder{}
			if err := proto.Unmarshal(d.Body, order); err != nil {
				logger.Warn("unmarsh body to order  %v", err)
				continue
			}
			handler.ReceiveParentOrder(order)
		}
		if err := ch.Close(); err != nil {
		}
	}()

	return nil

}

