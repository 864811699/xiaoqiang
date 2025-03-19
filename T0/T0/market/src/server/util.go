package server

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/go-stomp/stomp/v3"
	"github.com/golang/protobuf/proto"
	"github.com/jordan-wright/email"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"io"
	"log"
	"market/pkg/kdb"
	pb "market/pkg/pb"
	"net/smtp"
	"time"
)

func MarketActiveToKdb(amqMsg stomp.Message) ([]*kdb.Market, error) {
	msg :=amqMsg.Body
	var err error
	if compressType := amqMsg.Header.Get("CompressType"); compressType == "1" {
		msg, err = ungzip(amqMsg.Body)
		if err != nil {
			return nil, errors.Wrapf(err, "active market msg ungzip fail,gzipmsg::%s", string(amqMsg.Body))
		}
	}

	sq, err := decodeActiveMq(msg)
	if err != nil {
		return nil, errors.Wrapf(err, "active market decode fail")
	}
	if sq != nil {
		return prasePbToKdbMarket(sq), nil
	}

	return nil, fmt.Errorf("this msg is not market")
}

func decodeActiveMq(msg []byte) (*pb.StockQuotationList, error) {
	market := &pb.Message{}
	err := proto.Unmarshal(msg, proto.Message(market))
	if err != nil {
		return nil, err
	}
	if market.MsgType == pb.MsgType_Notify_Quotation_Stock {
		return market.Notify.StockQuotationList, nil
	}
	return nil, nil
}

func ungzip(bts []byte) ([]byte, error) {
	reader := bytes.NewReader(bts)
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create GZIP reader")
	}
	defer gzipReader.Close()
	var outBuffer bytes.Buffer
	_, err = io.Copy(&outBuffer, gzipReader)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to decompress data")
	}
	return outBuffer.Bytes(), nil
}

//40503 434691000 11h15m3.434691s
func prasePbToKdbMarket(m *pb.StockQuotationList) []*kdb.Market {
	//fmt.Printf("prasePbToKdbMarket  %#v\n",m)
	markets := make([]*kdb.Market, len(m.StockQuotations), len(m.StockQuotations))
	for k, _ := range m.StockQuotations {
		m := m.StockQuotations[k]
		markets[k] = &kdb.Market{
			Sym:                  m.Instrument_ID,
			Time:                 time.Now(),
			ZWindCode:            fmt.Sprintf("%s.%s", m.Instrument_ID, m.Exchange_ID),
			ActionDay:            0,
			NTime:                m.Time,
			Status:               0,
			PreClose:             int32(m.PreClosePrice),
			Open:                 int32(m.OpenPrice),
			High:                 int32(m.HighPrice),
			Low:                  int32(m.LowPrice),
			Match:                int32(m.LastPrice),
			AskPrice1:            int32(m.AskPrice[0]),
			AskPrice2:            int32(m.AskPrice[1]),
			NAskPrice3:           int32(m.AskPrice[2]),
			NAskPrice4:           int32(m.AskPrice[3]),
			NAskPrice5:           int32(m.AskPrice[4]),
			NAskPrice6:           int32(m.AskPrice[5]),
			NAskPrice7:           int32(m.AskPrice[6]),
			NAskPrice8:           int32(m.AskPrice[7]),
			NAskPrice9:           int32(m.AskPrice[8]),
			NAskPrice10:          int32(m.AskPrice[9]),
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
			NBidPrice1:           int32(m.BidPrice[0]),
			NBidPrice2:           int32(m.BidPrice[1]),
			NBidPrice3:           int32(m.BidPrice[2]),
			NBidPrice4:           int32(m.BidPrice[3]),
			NBidPrice5:           int32(m.BidPrice[4]),
			NBidPrice6:           int32(m.BidPrice[5]),
			NBidPrice7:           int32(m.BidPrice[6]),
			NBidPrice8:           int32(m.BidPrice[7]),
			NBidPrice9:           int32(m.BidPrice[8]),
			NBidPrice10:          int32(m.BidPrice[9]),
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
			ITurnover:            0, //todo
			NTotalBidVol:         int32(m.TotalBidVol),
			NTotalAskVol:         int32(m.TotalAskVol),
			NWeightedAvgBidPrice: int32(m.AvgBidPrice),
			NWeightedAvgAskPrice: int32(m.AvgAskPrice),
			NIOPV:                int32(m.IOPV),
			NYieldToMaturity:     0,
			NHighLimited:         int32(m.HighLimited),
			NLowLimited:          int32(m.LowLimited),
			NSyl1:                0,
			NSyl2:                0,
			NSD2:                 0,
		}
	}
	return markets
}

func floatToInt10000(f float64) int32 {
	f1 := decimal.NewFromFloat(f)
	i := decimal.NewFromInt(10000)
	v := f1.Mul(i)
	return int32(v.IntPart())
}

func sendEmail() {
	e := email.NewEmail()
	//设置发送方的邮箱
	e.From = "smj <864811699@qq.com>"
	// 设置接收方的邮箱
	e.To = []string{"shumingjun@minghsiim.com"}
	//设置主题
	e.Subject = "这是主题"
	//设置文件发送的内容
	e.Text = []byte("www.topgoer.com是个不错的go语言中文文档")
	//设置服务器相关的配置
	err := e.Send("smtp.qq.com:25", smtp.PlainAuth("", "864811699@qq.com", "gaymswketnhobecf", "smtp.qq.com"))
	if err != nil {
		log.Fatal(err)
	}
}
