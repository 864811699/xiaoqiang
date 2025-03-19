package trade

import (
	"fmt"
	"github.com/shopspring/decimal"
	"golang.org/x/text/encoding/simplifiedchinese"
	"strconv"
	"strings"
)

type void struct{}

var inSet void

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
func Abs32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

func pack(acct string, id int32) string {
	return fmt.Sprintf("%s_%d", acct, id)
}

func unpack(entrustNo string) int32 {
	if strings.Contains(entrustNo, "_") {
		strSlice := strings.Split(entrustNo, "_")
		n := len(strSlice) - 1
		idx, err := strconv.Atoi(strSlice[n])
		if err == nil {
			return int32(idx)
		}
		return -1
	}
	return -1
}

func mul(num1 float64, num2 int) float64 {
	d1 := decimal.NewFromFloat(num1).Mul(decimal.NewFromFloat(float64(num2)))
	f1 := d1.Round(4).String()
	f2, _ := strconv.ParseFloat(f1, 32)
	return f2
}

func div(num1 float64, num2 int) float64 {
	if num2==0 {
		return 0
	}
	d1 := decimal.NewFromFloat(num1).Div(decimal.NewFromFloat(float64(num2)))
	f1 := d1.Round(4).String()
	f2, _ := strconv.ParseFloat(f1, 32)
	return f2
}

func sumPrice(lastPrice,bidPrice float64,bidQty int) float64 {
	curSumPrice:=mul(bidPrice,bidQty)
	return  curSumPrice+lastPrice
}

func getMsg(msg string) string{
	if isGBK([]byte(msg)) {
		msg,_=Gbk2Utf8(msg)
	}
	return msg
}

func isGBK(data []byte) bool {
	length := len(data)
	var i int = 0
	for i < length {
		//fmt.Printf("for %x\n", data[i])
		if data[i] <= 0xff {
			//编码小于等于127,只有一个字节的编码，兼容ASCII吗
			i++
			continue
		} else {
			//大于127的使用双字节编码
			if  data[i] >= 0x81 &&
				data[i] <= 0xfe &&
				data[i + 1] >= 0x40 &&
				data[i + 1] <= 0xfe &&
				data[i + 1] != 0xf7 {
				i += 2
				continue
			} else {
				return false
			}
		}
	}
	return true
}

func Gbk2Utf8(s string)(string,error)  {
	return simplifiedchinese.GB18030.NewDecoder().String(s)
}



