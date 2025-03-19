package trade

import (
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"strconv"
	"strings"
)

func pack(acct string, id int32) string {

	return fmt.Sprintf("%s-%d", acct, id)
}

func unpack(entrustNo string) (string, int32) {
	if strings.Contains(entrustNo, "-") {
		fields := strings.Split(entrustNo, "-")
		localId, _ := strconv.Atoi(fields[1])
		return fields[0], int32(localId)
	}
	return "", -1
}

var StatusMapInt32 = map[string]int32{
	"48": 1,
	"49": 1,
	"50": 1,
	"51": 3,
	"52": 3,
	"53": 5,
	"54": 5,
	"55": 2,
	"56": 4,
	"57": 6,

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