package mq

import (
	"fmt"
	gproto "github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
	"go.uber.org/zap"

)

type rabbitMqManage struct {
	conn                  *amqp.Connection
	pubParentOrderchannel *amqp.Channel
	hasSubParentOrder     bool
}

type HandleParentOrder interface {
	ReceiveParentOrder(order *proto.ParentOrder)
}

var RabbitMqManage = &rabbitMqManage{
	hasSubParentOrder: false,
}

func (x *rabbitMqManage) GetConn() *amqp.Connection {
	return x.conn
}

func (x *rabbitMqManage) Close() {
	if x.pubParentOrderchannel != nil {
		if err := x.pubParentOrderchannel.Close(); err != nil {
		}
	}
	if x.conn != nil {
		if err := x.conn.Close(); err != nil {
		}
	}
	x.pubParentOrderchannel = nil
	x.conn = nil
}

func (x *rabbitMqManage) Init(rabbitUrl string) error {
	if x.conn == nil {
		newConn, err := amqp.Dial(rabbitUrl)
		if err != nil {
			return err
		}
		x.conn = newConn
	}
	return nil

}

func (x *rabbitMqManage) InitAlgoExchange() error {
	ch, err := x.GetConn().Channel()
	if err != nil {
		return err
	}
	err = ch.ExchangeDeclare("algo", "topic", false, false, false, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func (x *rabbitMqManage) PubParentOrder(order *proto.ParentOrder) error {
	routeKey := fmt.Sprintf("algo.%s", order.Algo)
	body, err := gproto.Marshal(order)
	if err != nil {
		return err
	}
	if x.pubParentOrderchannel == nil {
		x.pubParentOrderchannel, err = x.GetConn().Channel()
		if err != nil {
			return err
		}
	}
	err = x.pubParentOrderchannel.Publish("algo", routeKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        body,
	})
	if err != nil {
		return err
	}
	return nil
}

// SubAllParentOrder 订阅所有母单回报
func (x *rabbitMqManage) SubAllParentOrder(handler HandleParentOrder, queueName string) error {
	if x.hasSubParentOrder {
		return nil
	}

	ch, err := x.GetConn().Channel()
	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare("algo", "topic", false, false, false, false, nil)
	if err != nil {
		return err
	}

	var q amqp.Queue
	if queueName == "" {
		q, err = ch.QueueDeclare("", false, true, true, false, nil)
		if err != nil {
			return err
		}
	} else {
		q, err = ch.QueueDeclare(queueName, false, false, false, false, nil)
		if err != nil {
			return err
		}
	}

	err = ch.QueueBind(q.Name, "algo.*", "algo", false, nil)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			order := &proto.ParentOrder{}
			if err := gproto.Unmarshal(d.Body, order); err != nil {
				zap.L().Warn("unmarsh body to order", zap.Error(err))
				continue
			}
			handler.ReceiveParentOrder(order)
		}
		if err := ch.Close(); err != nil {
		}
	}()

	return nil
}

// SubParentOrderByAlgoName 订阅制定算法回报
func (x *rabbitMqManage) SubParentOrderByAlgoName(handler HandleParentOrder, algoNames []string, queueName string) error {
	ch, err := x.GetConn().Channel()
	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare("algo", "topic", false, false, false, false, nil)
	if err != nil {
		return err
	}

	var q amqp.Queue
	if queueName == "" {
		q, err = ch.QueueDeclare("", false, true, true, false, nil)
		if err != nil {
			return err
		}
	} else {
		q, err = ch.QueueDeclare(queueName, false, false, true, false, nil)
		if err != nil {
			return err
		}
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
			order := &proto.ParentOrder{}
			if err := gproto.Unmarshal(d.Body, order); err != nil {
				zap.L().Warn("unmarsh body to order", zap.Error(err))
				continue
			}
			handler.ReceiveParentOrder(order)
		}
		if err := ch.Close(); err != nil {
		}
	}()

	return nil
}
