package activemq

import (
	"github.com/go-stomp/stomp/v3"
	"github.com/golang/protobuf/proto"
	pb "market/pkg/pb"
	"testing"
)

func TestConn(t *testing.T)  {

	conn,err:=stomp.Dial("tcp","115.238.186.119:61613")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("connect success")

	sub,err:=conn.Subscribe("/topic/QUOTATION.STOCK.SNAPSHOT",stomp.AckAuto)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("sub success")

	msg:=<-sub.C

	msgs:=&pb.Message{}
	err=proto.Unmarshal(msg.Body,proto.Message(msgs))
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%#v",msgs.Notify.StockQuotationList)


}

func TestListMap(t *testing.T) {
	l:=[]int{9,9,9}
	for i := range l {
		t.Log(i)
	}
	m:=map[int]int{9:9,8:8}
	for i := range m {
		t.Log(i)
	}
}