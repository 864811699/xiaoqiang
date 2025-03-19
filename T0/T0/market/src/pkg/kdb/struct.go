package kdb

import (
	kdb "github.com/nsxy/kdbgo"
	"time"
)

var entrustCols = []string{ "time", "sym","szWindCode", "nActionDay", "nTime", "nStatus", "nPreClose", "nOpen", "nHigh", "nLow", "nMatch", "nAskPrice1", "nAskPrice2", "nAskPrice3", "nAskPrice4", "nAskPrice5", "nAskPrice6", "nAskPrice7", "nAskPrice8", "nAskPrice9", "nAskPrice10", "nAskVol1", "nAskVol2", "nAskVol3", "nAskVol4", "nAskVol5", "nAskVol6", "nAskVol7", "nAskVol8", "nAskVol9", "nAskVol10", "nBidPrice1", "nBidPrice2", "nBidPrice3", "nBidPrice4", "nBidPrice5", "nBidPrice6", "nBidPrice7", "nBidPrice8", "nBidPrice9", "nBidPrice10", "nBidVol1", "nBidVol2", "nBidVol3", "nBidVol4", "nBidVol5", "nBidVol6", "nBidVol7", "nBidVol8", "nBidVol9", "nBidVol10", "nNumTrades", "iVolume", "iTurnover", "nTotalBidVol", "nTotalAskVol", "nWeightedAvgBidPrice", "nWeightedAvgAskPrice", "nIOPV", "nYieldToMaturity", "nHighLimited", "nLowLimited", "nSyl1", "nSyl2", "nSD2"}

type Market struct {
	Sym                  string
	Time                 time.Time
	ZWindCode            string
	ActionDay            int32
	NTime                 int32
	Status               int32
	PreClose             int32
	Open                 int32
	High                 int32
	Low                  int32
	Match                int32
	AskPrice1            int32
	AskPrice2            int32
	NAskPrice3           int32
	NAskPrice4           int32
	NAskPrice5           int32
	NAskPrice6           int32
	NAskPrice7           int32
	NAskPrice8           int32
	NAskPrice9           int32
	NAskPrice10          int32
	NAskVol1             int32
	NAskVol2             int32
	NAskVol3             int32
	NAskVol4             int32
	NAskVol5             int32
	NAskVol6             int32
	NAskVol7             int32
	NAskVol8             int32
	NAskVol9             int32
	NAskVol10            int32
	NBidPrice1           int32
	NBidPrice2           int32
	NBidPrice3           int32
	NBidPrice4           int32
	NBidPrice5           int32
	NBidPrice6           int32
	NBidPrice7           int32
	NBidPrice8           int32
	NBidPrice9           int32
	NBidPrice10          int32
	NBidVol1             int32
	NBidVol2             int32
	NBidVol3             int32
	NBidVol4             int32
	NBidVol5             int32
	NBidVol6             int32
	NBidVol7             int32
	NBidVol8             int32
	NBidVol9             int32
	NBidVol10            int32
	NNumTrades           int32
	IVolume              int32
	ITurnover            int32
	NTotalBidVol         int32
	NTotalAskVol         int32
	NWeightedAvgBidPrice int32
	NWeightedAvgAskPrice int32
	NIOPV                int32
	NYieldToMaturity     int32
	NHighLimited         int32
	NLowLimited          int32
	NSyl1                int32
	NSyl2                int32
	NSD2                 int32
}


func parseMarketToKdb(market Market) *kdb.K {
	return kdb.NewTable(entrustCols, []*kdb.K{
		kdb.TimeNV([]time.Time{market.Time}),
		kdb.SymbolV([]string{market.Sym}),
		kdb.SymbolV([]string{market.ZWindCode}),
		kdb.IntV([]int32{market.ActionDay}),
		kdb.IntV([]int32{market.NTime}),
		kdb.IntV([]int32{market.Status}),
		kdb.IntV([]int32{market.PreClose}),
		kdb.IntV([]int32{market.Open}),
		kdb.IntV([]int32{market.High}),
		kdb.IntV([]int32{market.Low}),
		kdb.IntV([]int32{market.Match}),
		kdb.IntV([]int32{market.AskPrice1}),
		kdb.IntV([]int32{market.AskPrice2}),
		kdb.IntV([]int32{market.NAskPrice3}),
		kdb.IntV([]int32{market.NAskPrice4}),
		kdb.IntV([]int32{market.NAskPrice5}),
		kdb.IntV([]int32{market.NAskPrice6}),
		kdb.IntV([]int32{market.NAskPrice7}),
		kdb.IntV([]int32{market.NAskPrice8}),
		kdb.IntV([]int32{market.NAskPrice9}),
		kdb.IntV([]int32{market.NAskPrice10}),
		kdb.IntV([]int32{market.NAskVol1}),
		kdb.IntV([]int32{market.NAskVol2}),
		kdb.IntV([]int32{market.NAskVol3}),
		kdb.IntV([]int32{market.NAskVol4}),
		kdb.IntV([]int32{market.NAskVol5}),
		kdb.IntV([]int32{market.NAskVol6}),
		kdb.IntV([]int32{market.NAskVol7}),
		kdb.IntV([]int32{market.NAskVol8}),
		kdb.IntV([]int32{market.NAskVol9}),
		kdb.IntV([]int32{market.NAskVol10}),
		kdb.IntV([]int32{market.NBidPrice1}),
		kdb.IntV([]int32{market.NBidPrice2}),
		kdb.IntV([]int32{market.NBidPrice3}),
		kdb.IntV([]int32{market.NBidPrice4}),
		kdb.IntV([]int32{market.NBidPrice5}),
		kdb.IntV([]int32{market.NBidPrice6}),
		kdb.IntV([]int32{market.NBidPrice7}),
		kdb.IntV([]int32{market.NBidPrice8}),
		kdb.IntV([]int32{market.NBidPrice9}),
		kdb.IntV([]int32{market.NBidPrice10}),
		kdb.IntV([]int32{market.NBidVol1}),
		kdb.IntV([]int32{market.NBidVol2}),
		kdb.IntV([]int32{market.NBidVol3}),
		kdb.IntV([]int32{market.NBidVol4}),
		kdb.IntV([]int32{market.NBidVol5}),
		kdb.IntV([]int32{market.NBidVol6}),
		kdb.IntV([]int32{market.NBidVol7}),
		kdb.IntV([]int32{market.NBidVol8}),
		kdb.IntV([]int32{market.NBidVol9}),
		kdb.IntV([]int32{market.NBidVol10}),
		kdb.IntV([]int32{market.NNumTrades}),
		kdb.IntV([]int32{market.IVolume}),
		kdb.IntV([]int32{market.ITurnover}),
		kdb.IntV([]int32{market.NTotalBidVol}),
		kdb.IntV([]int32{market.NTotalAskVol}),
		kdb.IntV([]int32{market.NWeightedAvgBidPrice}),
		kdb.IntV([]int32{market.NWeightedAvgAskPrice}),
		kdb.IntV([]int32{market.NIOPV}),
		kdb.IntV([]int32{market.NYieldToMaturity}),
		kdb.IntV([]int32{market.NHighLimited}),
		kdb.IntV([]int32{market.NLowLimited}),
		kdb.IntV([]int32{market.NSyl1}),
		kdb.IntV([]int32{market.NSyl2}),
		kdb.IntV([]int32{market.NSD2}),
	})
}
