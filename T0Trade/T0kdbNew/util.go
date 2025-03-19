package kdbtdx

import (
	"time"

	kdb "github.com/864811699/kdbgo"
)

func ord2resp(order *Order) *kdb.K {

	return kdb.NewTable(entrustCols, []*kdb.K{
		kdb.SymbolV([]string{order.Sym}),
		kdb.SymbolV([]string{order.Qid}),
		kdb.SymbolV([]string{order.Trader}),
		kdb.TimeZV([]time.Time{time.Now()}),
		kdb.IntV([]int32{order.EntrustNo}),
		kdb.SymbolV([]string{order.Stockcode}),
		kdb.FloatV([]float64{order.Askprice}),
		kdb.IntV([]int32{order.Orderqty}),
		kdb.FloatV([]float64{order.AvgPx}),
		kdb.IntV([]int32{order.CumQty}),
		kdb.IntV([]int32{order.Withdraw}),
		kdb.IntV([]int32{order.Status}),
		kdb.SymbolV([]string{order.Note}),
		kdb.IntV([]int32{order.ReqType}),
		kdb.SymbolV([]string{order.Params}),
	})
}
//"account","available","stockvalue","assureasset","totaldebit","balance","enablebailbalance","perassurescalevalue"}
//sfffffff
func assetToKdb(asset Asset)*kdb.K  {
	return kdb.NewTable(assetCoks, []*kdb.K{
		kdb.SymbolV([]string{asset.Account}),
		kdb.FloatV([]float64{asset.Available}),
		kdb.FloatV([]float64{asset.StockValue}),
		kdb.FloatV([]float64{asset.AssureAsset}),
		kdb.FloatV([]float64{asset.TotalDebit}),
		kdb.FloatV([]float64{asset.Balance}),
		kdb.FloatV([]float64{asset.EnableBailBalance}),
		kdb.FloatV([]float64{asset.PerAssurescaleValue}),
		kdb.SymbolV([]string{asset.Sym}),

	})
}
//"account","zqdm","zqmc","cbj","gpye","fdyk","ykbl","sj","gpsz","kyye"}
//sssfiffffi
func positionToKdb(position Position)*kdb.K  {
	return kdb.NewTable(positionCoks, []*kdb.K{
		kdb.SymbolV([]string{position.Account}),
		kdb.SymbolV([]string{position.Stockcode}),
		kdb.SymbolV([]string{position.StockName}),
		kdb.FloatV([]float64{position.Cost}),
		kdb.IntV([]int32{position.Gpye}),
		kdb.FloatV([]float64{position.Fdyk}),
		kdb.FloatV([]float64{position.Ykbl}),
		kdb.FloatV([]float64{position.Sj}),
		kdb.FloatV([]float64{position.Gpsz}),
		kdb.IntV([]int32{position.Kyye}),
		kdb.SymbolV([]string{position.Sym}),
	})
}

func ord2entrust(order *Order) *kdb.K {

	qty := order.Orderqty
	cumQty := order.CumQty
	if order.Side%2 == 1 {
		qty = -qty
		cumQty = -cumQty
	}

	var wd int32
	if order.Status == 5 {
		wd = order.Orderqty - order.CumQty
	}

	return kdb.NewTable(entrustCols, []*kdb.K{
		kdb.SymbolV([]string{order.Trader}),
		kdb.SymbolV([]string{order.Qid}),
		kdb.SymbolV([]string{order.Sym}),
		kdb.TimeZV([]time.Time{order.Time.Add(-8 * time.Hour)}),
		kdb.IntV([]int32{order.EntrustNo}),
		kdb.SymbolV([]string{order.Stockcode}),
		kdb.FloatV([]float64{order.Askprice}),
		kdb.IntV([]int32{qty}),
		kdb.FloatV([]float64{order.AvgPx}),
		kdb.IntV([]int32{cumQty}),
		kdb.IntV([]int32{wd}),
		kdb.IntV([]int32{order.Status}),
	})
}

func GetOrder(key int32) (Order, bool) {

	vi, ok := Tb.orderMap.Load(key)
	if !ok {
		return Order{}, false
	}
	v := vi.(*Order)
	return *v, true
}

func GetOrders() []Order {
	orders:=[]Order{}

	Tb.orderMap.Range(func(key, value interface{}) bool {
		order_ptr:=value.(*Order)
		orders=append(orders,*order_ptr)
		return true
	})

	return orders
}

func GetUnFinalizedOrderNo() map[int32]bool {

	res := make(map[int32]bool)
	Tb.orderMap.Range(func(_, value interface{}) bool {
		if o := value.(*Order); o.Status < 4 {
			res[o.EntrustNo] = true
		}
		return true
	})
	return res
}

func GetAllOrder() map[int32]Order {

	res := make(map[int32]Order)
	Tb.orderMap.Range(func(_, value interface{}) bool {
		if o,ok := value.(*Order); ok {
			res[o.EntrustNo] = *o
		}
		return true
	})
	return res
}

func Store(o *Order) {

	Tb.updateTab(o)
	Tb.oChan <- o
	Tb.orderMap.Store(o.EntrustNo, o)
}