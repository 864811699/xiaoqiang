package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type ActiveInfo struct {
	Addr  string `json:"addr"`
	Topic string `json:"topic"`
}

type KdbInfo struct {
	IP    string `json:"ip"`
	Port  int    `json:"port"`
	Auth  string `json:"auth"`
	Table string `json:"table"`
	Fun   string `json:"fun"`
}

type Conf struct {
	ActiveCfg *ActiveInfo `json:"active_cfg"`
	KdbCfg *KdbInfo `json:"kdb_cfg"`
}


func LoadCfg(path string) *Conf {
	f,err:=os.Open(path)
	if err != nil {
		panic(err)
	}

	bts,err:=ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	c:=Conf{}
	err=json.Unmarshal(bts,&c)
	if err != nil {
		panic(err)
	}
	return &c
}