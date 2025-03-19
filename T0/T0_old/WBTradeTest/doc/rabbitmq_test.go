package mq

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"oems/proto"
	"os"
	"os/signal"
	"testing"
	"time"
)

var rabbitMqUrl = "amqp://oms:oms@172.16.8.61:5672/"

func Test_rabbitMqManage_PubParentOrder(t *testing.T) {
	assert.Nil(t, RabbitMqManage.Init(rabbitMqUrl))
	order := proto.TestCreateParentOrder()
	for {
		time.Sleep(3 * time.Second)
		//order.FilledQuantity += 100
		order = proto.TestCreateParentOrder()
		assert.Nil(t, RabbitMqManage.PubParentOrder(order))
	}
}

type TestSubParentOrder struct {
}

func (x *TestSubParentOrder) ReceiveParentOrder(order *proto.ParentOrder) {
	fmt.Println(order.String())
}

func Test_rabbitMqManage_SubParentOrder(t *testing.T) {
	defer RabbitMqManage.Close()
	assert.Nil(t, RabbitMqManage.Init(rabbitMqUrl))
	testSubParentOrder := &TestSubParentOrder{}
	assert.Nil(t, RabbitMqManage.SubAllParentOrder(testSubParentOrder, ""))

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
}

func Test_rabbitMqManage_SubParentOrder2(t *testing.T) {
	defer RabbitMqManage.Close()
	assert.Nil(t, RabbitMqManage.Init(rabbitMqUrl))
	testSubParentOrder := &TestSubParentOrder{}
	assert.Nil(t, RabbitMqManage.SubAllParentOrder(testSubParentOrder, ""))

	order := proto.TestCreateParentOrder()
	for {
		time.Sleep(3 * time.Second)
		//order.FilledQuantity += 100
		order = proto.TestCreateParentOrder()
		assert.Nil(t, RabbitMqManage.PubParentOrder(order))
	}
}

func Test_rabbitMqManage_SubParentOrder1(t *testing.T) {
	defer RabbitMqManage.Close()
	assert.Nil(t, RabbitMqManage.Init(rabbitMqUrl))
	testSubParentOrder := &TestSubParentOrder{}
	assert.Nil(t, RabbitMqManage.SubAllParentOrder(testSubParentOrder, "db_parent_order"))

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
}

func Test_rabbitMqManage_SubParentOrderByAlgoName(t *testing.T) {
	assert.Nil(t, RabbitMqManage.Init(rabbitMqUrl))
	testSubParentOrder := &TestSubParentOrder{}
	assert.Nil(t, RabbitMqManage.SubParentOrderByAlgoName(testSubParentOrder, []string{"cornerstone_capstone_v2"}, ""))

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	RabbitMqManage.Close()
}

func Test_rabbitMqManage_InitAlgoExchange(t *testing.T) {
	assert.Nil(t, RabbitMqManage.Init(rabbitMqUrl))
	assert.Nil(t, RabbitMqManage.InitAlgoExchange())
}
