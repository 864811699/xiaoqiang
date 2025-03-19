package main

import (
	"bufio"
	"flag"
	"fmt"
	kdbtdx "github.com/864811699/T0kdb"
	logger "github.com/alecthomas/log4go"
	"os"
	"regexp"
	"strings"
	"time"

	"uat/trade"
)

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
	api := trade.NewApi("","./etc/trade.json")
	kdbtdx.Run(api)
}
