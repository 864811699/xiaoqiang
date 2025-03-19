package trade

import (
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"unsafe"
)
/*

#include <windows.h>
#include <stdio.h>
//和DLL交互函数
//pszRequest 请求json数据 szResult应答json数据，nMaxLen szResult的最大长度
//返回值是应答的实际数据长度，如果<=0,那么是系统错误
//用户名：ywu  密码：1
//int RequestData(const char* pszData,char* szResult,uint32_t nMaxLen);
//参数说明：
//pszData – 请求字符串 json 格式
//szResult – 应答数据返回内容
//nMaxLen – szResult 的最大长度

typedef int(*RequestData)(const char* pszRequest, char szResult[1024000], unsigned int nMaxLen);//查询各类历史数据
RequestData lpRequestData = NULL;
int Cmd(const char* pszData, char szResult[1024000],unsigned int nMaxLen){
	HINSTANCE m_hInstance = LoadLibrary("wtrader.dll");
	lpRequestData = (RequestData)GetProcAddress(m_hInstance, "RequestData");
	//认证交易员
	int nLen = 0;
	nLen = lpRequestData(pszData, szResult, nMaxLen);
	return  nLen;
}


*/
import "C"


const MAX_RESPONSE_LEN=1024000

//pszRequest=`{"cmd":"auth", "ps":{"user":"002762004", "pwd":"1"}}`
func CallRequestData(pszRequest []byte ) (string, error) {
	bts := make([]byte, MAX_RESPONSE_LEN)
	rsp := (*C.char)(C.CBytes(bts))
	defer C.free(unsafe.Pointer(rsp))

	req := (*C.char)(C.CBytes(pszRequest))
	defer C.free(unsafe.Pointer(req))

	c_rsp_len := C.Cmd(req, rsp, MAX_RESPONSE_LEN)
	if c_rsp_len > 0 {
		s := C.GoStringN(rsp, c_rsp_len)
		return s, nil
	}
	return "", fmt.Errorf("Call RequestData fail :%d", c_rsp_len)

}

func GbkToUtf8(s string) (string, error) {
	return simplifiedchinese.GB18030.NewDecoder().String(s)
}

