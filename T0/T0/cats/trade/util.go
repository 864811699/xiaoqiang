package trade

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var day = time.Now().Format("0102")

func pack(s string, id int32) string {
	return fmt.Sprintf("%s-%d", s, id)
}

func unpack(entrustNo string) (string, int32) {
	if strings.Contains(entrustNo, "-") {
		fields := strings.Split(entrustNo, "-")
		localId, _ := strconv.Atoi(fields[1])
		return fields[0], int32(localId)
	}
	return "", -1
}
