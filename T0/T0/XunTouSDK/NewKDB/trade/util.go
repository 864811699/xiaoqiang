package trade

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	logger "github.com/alecthomas/log4go"
	"github.com/shopspring/decimal"
	"golang.org/x/text/encoding/simplifiedchinese"
	"net/http"
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

	return fmt.Sprintf("%s-%d", acct, id)
}

func unpack(entrustNo string)(acct string,idx int32) {

	if strings.Contains(entrustNo, "-") {
		fields:=strings.Split(entrustNo,"-")
		if len(fields)==2 {
			entrustno,err:=strconv.Atoi(fields[1])
			if err == nil {
				acct=fields[0]
				idx=int32(entrustno)
				return acct,idx
			}
		}

	}
	return "",-1

}

func get_exchange_id(stockcode string) string {
	var id string
	switch stockcode[0] {
	case '0', '3':
		id = "SZE"
	case '5', '6':
		id = "SSE"
	}
	return id
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

const (
	ConstHeader         = "***"
	ConstHeaderLength   =3
	ConstSaveDataLength = 4
)

//整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//封包
func PacketMsg(message []byte) []byte {
	return append(append([]byte(ConstHeader), IntToBytes(len(message))...), message...)
}




//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}


//解包
func UnpackMsg(buffer []byte, readerChannel chan []byte) []byte {
	length := len(buffer)

	var i int
	for i = 0; i < length; i = i + 1 {
		if length < i+ConstHeaderLength+ConstSaveDataLength {
			break
		}
		if string(buffer[i:i+ConstHeaderLength]) == ConstHeader {
			messageLength := BytesToInt(buffer[i+ConstHeaderLength : i+ConstHeaderLength+ConstSaveDataLength])
			if length < i+ConstHeaderLength+ConstSaveDataLength+messageLength {
				break
			}
			data := buffer[i+ConstHeaderLength+ConstSaveDataLength : i+ConstHeaderLength+ConstSaveDataLength+messageLength]
			readerChannel <- data

			i += ConstHeaderLength + ConstSaveDataLength + messageLength - 1
		}
	}

	if i == length {
		return make([]byte, 0)
	}
	return buffer[i:]


}

type MessageBody struct {
	MsgType string `json:"msgtype"`
	Text    Text   `json:"text"`
}
type Text struct {
	Content        string `json:"content"`
	Mentioned_list string `json:"mentioned_list"`
}

func Send_weiXin(s,url_str string) {

	text := Text{s, ""}
	messageBody := MessageBody{"text", text}
	by, _ := json.Marshal(&messageBody)
	reader := bytes.NewReader(by)
	_, err := http.Post(url_str, "application/json", reader)
	if err != nil {
		logger.Info("err::%v", err)
		return
	}

}

func getMsg(msg string) string{
	//if isGBK([]byte(msg)) {
	//	msg,_=Gbk2Utf8(msg)
	//}
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