package main
import (
	"encoding/json"
	"fmt"
	"time"
	"tradeApi/trade"
)

func cmd(req string) string{

	//queryInfo:=`{"cmd":"actinfo"}`
	t1:=time.Now()
	fmt.Println("开始执行时间:",t1.Format("15:04:05.000"))

	rsp,err:=trade.CallRequestData([]byte(req))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	t2:=time.Now()
	fmt.Printf("调用结束::%s\t耗时::%v\n",t2.Format("15:04:05.000"),t2.Sub(t1))
	utf8_str,err:=trade.GbkToUtf8(rsp)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	fmt.Println(utf8_str)
	return utf8_str
}


func main() {
	loginMsg:=`{"cmd":"auth", "ps":{"user":"002762004", "pwd":"1"}}`
	queryInfo:=`{"cmd":"actinfo"}`
	//queryOrders:=`{"cmd":"wt_query_weituo","id":8571}`
	//queryDeal:=`{"cmd":"wt_query_chengjiao","id":8571}`
	cmd(loginMsg)
	msg:=cmd(queryInfo)
	//cmd(queryOrders)
	//cmd(queryDeal)
	accountInfos := trade.AccountInfoRsp{}
	json.Unmarshal([]byte(msg), &accountInfos)
	for _,v:=range accountInfos.AccountInfos {
		fmt.Printf("%#v",*v)
	}

}