package server

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"market/pkg/kdb"
	pb "market/pkg/pb"
	"time"
)

func MarketActiveToKdb(msg []byte) (*kdb.Market,error) {
	sq,err:=decodeActiveMq(msg)
	if err!=nil{
		return nil,errors.Wrapf(err,"active market decode fail")
	}
	if sq!=nil{
		return prasePbToKdbMarket(sq),nil
	}
	return nil,fmt.Errorf("this msg is not market")
}


func decodeActiveMq(msg []byte) (*pb.StockQuotation,error) {
	market:=&pb.Message{}
	err:=proto.Unmarshal(msg,proto.Message(market))
	if err != nil {
		return nil,err
	}
	if market.MsgType==pb.MsgType_Notify_Quotation_Stock{
		return market.Notify.StockQuotation,nil
	}
	return nil,nil
}

//40503 434691000 11h15m3.434691s
func prasePbToKdbMarket(m *pb.StockQuotation) *kdb.Market {
	//fmt.Printf("prasePbToKdbMarket  %#v\n",m)
	return &kdb.Market{
		Sym:                  m.Instrument_ID,
		Time:                 time.Now(),
		ZWindCode:            fmt.Sprintf("%s.%s",m.Instrument_ID,m.Exchange_ID),
		ActionDay:            0,
		NTime:                m.Time,
		Status:               0,
		PreClose:             floatToInt10000(m.PreClosePrice),
		Open:                 floatToInt10000(m.OpenPrice),
		High:                 floatToInt10000(m.HighPrice),
		Low:                  floatToInt10000(m.LowPrice),
		Match:                floatToInt10000(m.LastPrice),
		AskPrice1:            floatToInt10000(m.AskPrice[0]),
		AskPrice2:            floatToInt10000(m.AskPrice[1]),
		NAskPrice3:           floatToInt10000(m.AskPrice[2]),
		NAskPrice4:           floatToInt10000(m.AskPrice[3]),
		NAskPrice5:           floatToInt10000(m.AskPrice[4]),
		NAskPrice6:           floatToInt10000(m.AskPrice[5]),
		NAskPrice7:           floatToInt10000(m.AskPrice[6]),
		NAskPrice8:           floatToInt10000(m.AskPrice[7]),
		NAskPrice9:           floatToInt10000(m.AskPrice[8]),
		NAskPrice10:          floatToInt10000(m.AskPrice[9]),
		NAskVol1:             int32(m.AskSize[0]),
		NAskVol2:             int32(m.AskSize[1]),
		NAskVol3:             int32(m.AskSize[2]),
		NAskVol4:             int32(m.AskSize[3]),
		NAskVol5:             int32(m.AskSize[4]),
		NAskVol6:             int32(m.AskSize[5]),
		NAskVol7:             int32(m.AskSize[6]),
		NAskVol8:             int32(m.AskSize[7]),
		NAskVol9:             int32(m.AskSize[8]),
		NAskVol10:            int32(m.AskSize[9]),
		NBidPrice1:           floatToInt10000(m.BidPrice[0]),
		NBidPrice2:           floatToInt10000(m.BidPrice[1]),
		NBidPrice3:           floatToInt10000(m.BidPrice[2]),
		NBidPrice4:           floatToInt10000(m.BidPrice[3]),
		NBidPrice5:           floatToInt10000(m.BidPrice[4]),
		NBidPrice6:           floatToInt10000(m.BidPrice[5]),
		NBidPrice7:           floatToInt10000(m.BidPrice[6]),
		NBidPrice8:           floatToInt10000(m.BidPrice[7]),
		NBidPrice9:           floatToInt10000(m.BidPrice[8]),
		NBidPrice10:          floatToInt10000(m.BidPrice[9]),
		NBidVol1:             int32(m.BidSize[0]),
		NBidVol2:             int32(m.BidSize[1]),
		NBidVol3:             int32(m.BidSize[2]),
		NBidVol4:             int32(m.BidSize[3]),
		NBidVol5:             int32(m.BidSize[4]),
		NBidVol6:             int32(m.BidSize[5]),
		NBidVol7:             int32(m.BidSize[6]),
		NBidVol8:             int32(m.BidSize[7]),
		NBidVol9:             int32(m.BidSize[8]),
		NBidVol10:            int32(m.BidSize[9]),
		NNumTrades:           int32(m.NumTrades),
		IVolume:              int32(m.TotalVolume),
		ITurnover:            0,//todo
		NTotalBidVol:         int32(m.TotalBidVol),
		NTotalAskVol:         int32(m.TotalAskVol),
		NWeightedAvgBidPrice: int32(m.AvgBidPrice),
		NWeightedAvgAskPrice: int32(m.AvgAskPrice),
		NIOPV:                int32(m.IOPV),
		NYieldToMaturity:     0,
		NHighLimited:         floatToInt10000(m.HighLimited),
		NLowLimited:          floatToInt10000(m.LowLimited),
		NSyl1:                0,
		NSyl2:                0,
		NSD2:                 0,
	}
}

func floatToInt10000(f float64)  int32{
	f1:=decimal.NewFromFloat(f)
	i:=decimal.NewFromInt(10000)
	v:=f1.Mul(i)
	return int32(v.IntPart())
}