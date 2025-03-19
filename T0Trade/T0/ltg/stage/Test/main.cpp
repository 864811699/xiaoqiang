#include "PT_QuantApi.h"
#include <string.h>
#include <iostream>
#include <windows.h>
#include <thread>
#include <mutex>
//#include <unistd.h>

using namespace std;
using namespace QuantPlus;

int orderid = 0;
mutex g_mtx;

#define LOGVAR(var) " "  #var ":" << var

// if pRspInfo field exists, you should always check it first.
class Callback : public PT_QuantSpi {
    void OnConnect(int nSrvType) {}
    ///@brief 通知断开
    ///@param nSrvType 业务服务器类型 参考QuantPlus::PT_Quant_APPServerType
    ///@return 无
    ///@note 在业务服务器断开时主动通知
    void OnDisconnect(int nSrvType) {}

    void onRtnUserInfo(const PT_QuantUserBase* pInfo) {

        cout << "Login OK " << endl;

    }

    void onRtnLoginWarn(int nLoginWarnType) {
        cout << "Login Error " << endl;

    }

    void onRspQryOrder(const TD_RspQryOrder *rsp, int error, bool isEnd) override {
        cout << "OnRspQryOrder " << LOGVAR(rsp->nReqId) << LOGVAR(rsp->nOrderId) << LOGVAR(error) << LOGVAR(isEnd) << endl;

    }

    void onRspQryMatch(const TD_RspQryMatch *rsp, int error, bool isEnd) {
        cout << "onRspQryMatch " << LOGVAR(rsp->nReqId) << LOGVAR(rsp->nOrderId) << LOGVAR(error) << LOGVAR(isEnd) << endl;

    };

	void onRspQryPosition(const TD_RspQryPosition *rsp, int error, bool isEnd) {
        cout << "onRspQryPosition " << LOGVAR(rsp->nReqId) << LOGVAR(error) << LOGVAR(isEnd) << endl;

    }

	void onRspOrderInsert(const TD_RspOrderInsert *rsp, int error) {
		unique_lock<mutex> lock(g_mtx);
        cout << "onRspOrderInsert " << LOGVAR(rsp->nReqId) << LOGVAR(error)  << endl;
        orderid = rsp->nOrderId;
    }

	void onRspOrderDelete(const TD_RspOrderDelete *rsp, int error) {
        cout << "onRspOrderDelete " << LOGVAR(rsp->nReqId) << LOGVAR(error)  << endl;
    }

    // if any order status change occured, OnRtnOrder would be called.
    void onRtnOrderStatusChangeNotice(const TD_RtnOrderStatusChangeNotice *notice) {
        cout << "onRtnOrderStatusChangeNotice " << LOGVAR(notice->nReqId) << LOGVAR(notice->nOrderId)  << endl;
    } 

    void onRtnOrderMatchNotice(const TD_RtnOrderMatchNotice *notice) {
        cout << "onRtnOrderMatchNotice " << LOGVAR(notice->nOrderId) << endl;

    }

    void OnRspSubQuote(MD_ReqID nReqID, MD_SubType nSubType, MD_CycType nCycType, 
                            const MD_CodeType *pSubWindCode, long nSubWindCodeNum, MD_ISODateTimeType szBeginTime, 
                            MD_ISODateTimeType szEndTime, int nErrNo, const char *szErrMsg) {
        cout << "OnRspSubQuote " << LOGVAR(nReqID) << LOGVAR(nErrNo)  << endl;
    }

	void OnRtnMarket(MD_ReqID nReqID, MD_DATA_MARKET *pMarket) {
        cout << "OnRtnMarket " << LOGVAR(pMarket->szCode)  << LOGVAR(pMarket->nMatch) << endl;        
    }

    void OnRtnTransaction(MD_ReqID nReqID, MD_DATA_TRANSACTION *pTrans) {
        cout << "OnRtnTransaction " << LOGVAR(pTrans->szCode) << LOGVAR(pTrans->nTime) << LOGVAR(pTrans->nPrice) << endl;        
    }


};

int main(int argc, char const *argv[])
{
    // Register Your Callback. 
    // Callback functions will all be called in a single thread, Please keep these functions simple.
    // large complex code will block the thread, slow down the whole ctpapi's performance.
    auto* cb = new Callback();

    // Get tradeapi interface. you SHOULD ONLY create this once.
    auto m_ptraderapi = PT_QuantApi::createApi(cb, true, PT_QuantTdAppEType::PT_QuantTdAppEType_Real, true, PT_QuantMdAppEType::PT_QuantMdAppEType_Real, true, true); 
    // 根据需求设置tdConnect mdConnect对应开启交易及行情功能

    // Initialize api
    // at this time, ctptrade will load the config file conf.json in the working dir, pls make sure you've configured that properly.
    /* conf.json template
    {
        "Log": {
            "LogFile": "./trade.log",   
            "LogLevel": "DEBUG"		
        },

        "Trade": {
            "ProductCode": "<product code>",
            "Addr": "<serverip:port>"
        }

        "Quote": {
            "ProductCode": "<product code>",
            "Addr": "<serverip:port>"
        }
    }
     */
    m_ptraderapi->Init();

    // Login in, result will be delivered by OnRspUserLogin.
    m_ptraderapi->Login((char*) "t10001", (char*) "888888");

    // now tradeapi starts functioning, any order changes would be informed through OnRtnOrder callback.
    // Besides, you can query server for data, queries and responses are linked by nRequestID.
    // It's important to guarantee nRequestID to be UNIQUE within a trading day. 
    // Queries with duplicate nRequestID will overlap each other causing undefined errors.

    Sleep(1000);    //此处做延时是为了等待异步connect连接成功，严格意义上需要等待onconnect函数返回判断对应的业务服务器连接成功之后才能做业务请求

    // pull all orders today. check OnRspQryOrder for results.
    TD_ReqQryOrder order;
    order.nReqId = 1000;
    m_ptraderapi->reqQryOrder(&order);

    // pull all transactions today. corresponding callback: OnRspQryTrade
    TD_ReqQryMatch trade;
    trade.nReqId = 1100;
    m_ptraderapi->reqQryMatch(&trade);

    // query for your quote. callback: OnRspQryInvestorPosition
    TD_ReqQryPosition pos;
    pos.nReqId = 1200;
    m_ptraderapi->reqQryPosition(&pos);

    // Cancel an order
    TD_ReqOrderDelete cancel;

    cancel.nOrderId = orderid; // OrderID,  OrderID = OrderSysID = OrderRef = unique order indentifier.

    //m_ptraderapi->reqOrderDelete(&cancel);


	int codeNum = 4;
	MD_CodeType *pSubWindCode = new MD_CodeType[codeNum];


	strncpy(pSubWindCode[0], "000725.SZ", sizeof(MD_CodeType));
	strncpy(pSubWindCode[1], "000997.SZ", sizeof(MD_CodeType));
	strncpy(pSubWindCode[2], "600000.SH", sizeof(MD_CodeType));
	strncpy(pSubWindCode[3], "603283.SH", sizeof(MD_CodeType));
	strncpy(pSubWindCode[4], "002971.SZ", sizeof(MD_CodeType));


	int reqId = 2000;
	//m_ptraderapi->ReqSubQuote(reqId, MD_SubType_market, MD_CycType_none, pSubWindCode, codeNum, "", "");

	Sleep(1);

	// Submit a new order, callback: OnRspOrderInsert
	std::thread t{ [m_ptraderapi]() {
		auto t = time(nullptr) / 10000;
		for (;;) {
			TD_ReqOrderInsert neworder;

			neworder.nReqId = t++;
			strcpy(neworder.szContractCode, "601933.SH"); // symbol
			neworder.nOrderPrice = 8.85 * 10000; // price
			neworder.nTradeType = TD_TradeType::TD_TradeType_Buy; // side
			neworder.nOrderVol = 200; // quantity

			//std::cout << " insert order " << std::endl;
			m_ptraderapi->reqOrderInsert(&neworder);

			Sleep(1000);

			neworder.nReqId = t++;
			strcpy(neworder.szContractCode, "000006.SZ"); // symbol
			neworder.nOrderPrice = 5.85 * 10000; // price
			neworder.nTradeType = TD_TradeType::TD_TradeType_Buy; // side
			neworder.nOrderVol = 100; // quantity

			//std::cout << " insert order " << std::endl;
			m_ptraderapi->reqOrderInsert(&neworder);

			Sleep(1000);

			neworder.nReqId = t++;
			strcpy(neworder.szContractCode, "000001.SZ"); // symbol
			neworder.nOrderPrice = 16.85 * 10000; // price
			neworder.nTradeType = TD_TradeType::TD_TradeType_Buy; // side
			neworder.nOrderVol = 100; // quantity

			//std::cout << " insert order " << std::endl;
			m_ptraderapi->reqOrderInsert(&neworder);

			Sleep(1000);
		}
	} };


    //std::cout << "done" << std::endl;
	
	t.join();

    //just prevent stop.
    int i;
    std::cin >> i;
    
    return 0;
}
