package kdbtdx

import (
	"encoding/json"
	"fmt"
	kdb "github.com/864811699/kdbgo"
	logger "github.com/alecthomas/log4go"
	"reflect"
	"testing"
	"time"
)

func TestKdb(t *testing.T) {

	s := newKdb("localhost", 5001, "test:pass", "./", []string{"testMS2", "testMS022"}, 10)
	//s.connect()
	//v, err := s.querySql()
	//for ix := range v {
	//	fmt.Printf("query res: %#v\n", v[ix])
	//}
	//fmt.Println("err", err)
	//	_ = s.KDBConn.Close()
	s.start()
	defer s.stop()
}


func TestRef(t *testing.T) {

	o := Request{
		Time: time.Now(),
		ParentId:  "q01",
		Trader:    "tr",
		Fund:      "ds",
		Strategy:  "dsx",
		Broker:    "cs",
		Stockcode: "601818",
		Stocktype: 0,
		Side:      3,
		Askprice:  0.25,
		Ordertype: 1,
		Orderqty:  1,
		Algorithm: "DMA",
		Params:    "",
	}

	objType := reflect.TypeOf(o)
	objVal := reflect.ValueOf(o)

	var data = make(map[string]string)
	for i := 0; i < objType.NumField(); i++ {
		if objVal.Field(i).CanInterface() {
			data[objType.Field(i).Name] = fmt.Sprintf("%v", objVal.Field(i).Interface())
		}
	}
	fmt.Println(data)

}

func TestUpd2(t *testing.T)  {
	k,err:=kdb.DialKDBTimeout("115.238.186.119",40000,"test:pass",time.Second*3)
	if err!=nil{
		t.Fatal("connect fail")
	}
	td:=new(tradeKdb)
	td.KDBConn=k
	sub("HeartBeat", []string{},*k)
	msg:=`{"ret":0,"query":[{"djje":0,"sjye":0,"orizzc":12484983.39,"kyzj":1018219.59,"fzje":5712544.95,"jzc":6772438.44,"zjye":1792271.03,"gpsz":10692712.36,"rzky":4649151,"rqfz":5712544.95,"rzfz":0,"zzc":12484983.39}]}`
	asset:=QueryAssetRsp{}

	err=json.Unmarshal([]byte(msg),&asset)
	if err != nil {
		t.Fatal(err)
	}
	asset.Responses[0].Account="aa"
	asset.Responses[0].Available=1.1
	td.send2Tab(PubAndUpdFunc, AssetTab, assetToKdb(asset.Responses[0]))

	//order:=Order{
	//	Request:Request{
	//		basic:     basic{Sym:"a",Qid:"a"},
	//		Time:      time.Time{},
	//
	//	},
	//
	//}
	//td.send2Tab(PubAndUpdFunc, AssetTab, assetToKdb(&order))
	//td.send2Tab(PubAndUpdFunc, ResponseTab, assetToKdb(&order))


	//err=k.AsyncCall(PubAndUpdFunc,kdb.Symbol(AssetTab),k_a)
	//if err != nil {
	//	t.Fatal(err)
	//}

}

func sub(tab string, sym []string,tb kdb.KDBConn)  {


		logger.Debug("Subscribing Kdb, table: %v, sym: %v", tab, sym)
		var symK *kdb.K
		if symNum := len(sym); symNum == 0 {
			symK = kdb.Symbol("")
		} else {
			symK = kdb.SymbolV(sym)
		}

		err := tb.AsyncCall(".u.sub", kdb.Symbol(tab), symK)
		if err != nil {
			_ = logger.Error("Failed to subscribe table: %v, sym: %v", tab, sym)
		}
}
type QueryAssetRsp struct {
	Ret       int         `json:"ret"`
	Responses []Asset `json:"query"`
}