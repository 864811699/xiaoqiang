package main

import (
	"bufio"
	"flag"
	"fmt"
	logger "github.com/alecthomas/log4go"
	"market/pkg/activemq"
	"market/pkg/conf"
	"market/pkg/kdb"
	"market/server"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"
)

const cfgPath="./etc/config.json"

func init() {

	logPtr := flag.String("log", "./log/cim.log", "log path to read from")
	flag.Parse()

	file, err := os.Open(*logPtr)

	var logName *string

	defer func() {
		file.Close()
		if logName!=nil {
			if err := os.Rename(*logPtr, fmt.Sprintf("%s%s.log", strings.Split(*logPtr, ".log")[0], *logName)); err != nil {
				fmt.Println("move trade log file error  ", err)
				return
			}
		}


	}()

	if err == nil {
		lines := bufio.NewScanner(file)
		logReg := regexp.MustCompile(`\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}`)
		lines.Scan()
		if text := lines.Text(); logReg.MatchString(text) {
			dt, err := time.Parse("2006/01/02 15:04:05", logReg.FindString(text))
			if err != nil {
				logger.Error("parse time error: %v", err)
				return
			}
			if tm := dt.Format("20060102"); tm < time.Now().Format("20060102") {
				logName = &tm
			}
		}
	}

	fmt.Println("log init success!!!")
}

func main() {
	logger.LoadConfiguration("./etc/log.xml")
	c:=conf.LoadCfg(cfgPath)

	a:=activemq.NewActive(c.ActiveCfg.Addr,c.ActiveCfg.Topic)
	k:=kdb.NewKdbHandle(c.KdbCfg.IP,c.KdbCfg.Auth,c.KdbCfg.Table,c.KdbCfg.Fun,c.KdbCfg.Port)

	s:=server.NewServer(k,a)
	go s.Start()
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	s.Stop()
	logger.Warn("程序结束")

}

