package trade

import (
	"encoding/json"
	"testing"
)

func TestJson(t *testing.T) {
	o:=OrderReq{
		Request:   Request{Cmd:"111"},
		Stockcode: "",
		Price:     0,
		ID:        0,
		Amount:    0,
		WtFlag:    "",
	}

	bts,_:=json.Marshal(o)
	t.Logf("%s",bts)

}
