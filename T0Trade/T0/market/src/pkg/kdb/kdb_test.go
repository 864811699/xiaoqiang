package kdb

import "testing"

func TestSubkdb(t *testing.T) {
	//ip:="139.196.52.182"
	ip:="115.238.186.119"
	port:=40000
	auth:="test:pass"
	fun:="upd0"
	tab:="Market"

	k:=NewKdbHandle(ip,auth,tab,fun,port)
	err:=k.Connect()
	if err != nil {
		t.Fatal(err)
	}



	err=k.subscribe([]string{"603000"})
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		err=k.readFromKdb()
		if err != nil {
			t.Fatal(err)
		}
	}()

	for msg:=range k.marketCh{
		t.Logf("%#v",*msg)

	}

}


//40503434691000 11h15m3.434691s