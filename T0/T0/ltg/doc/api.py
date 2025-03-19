import time
import os
import sys
import threading
from pyapi import PyApi, DataCallBack
from datetime import datetime, date, timedelta

pwd = os.path.dirname(os.path.realpath(__file__))
sys.path.insert(0, pwd + '/../')
import queue

from kdb.util import logger


def genreq():
    global reqcnt
    reqcnt += 1
    return int(round(time.time()) + reqcnt)


QUERY_ORDER = "query Order"
ACK_ORDER = "ack Order"
RSP_ORDER = "rsp Order"


class MyCallBack(DataCallBack):
    def __init__(self):
        super(MyCallBack, self).__init__()

        # 创建api实例，可以写在回调之外，这里放在回调class内只是方便回调函数调用api
        self.api = PyApi.PT_QuantApi(self, True, "Td_Real", "MD_Real")

        # 初始化API必须在创建之后立刻执行
        PyApi.Init()

        self.cv = threading.Condition()
        self.isLogin = False
        self.queue = queue.Queue()

    def OnConnect(self, type):
        # if type == 1:   #0为实时服务器 1为历史服务器 2为缓存服务器  推荐在0.1.2三种行情服务器全部连接成功后做行情请求能保证全时段行情稳定回调        10为交易服务器，连接成功之后才可做交易业务请求
        logger.info("OnConnnect{} success !!:".format(type))
        if type == 10:
            logger.info("OnConnnect trade{} success !!:".format(type))
        with self.cv:
            self.isLogin = True
            self.cv.notify()

    def OnDisconnect(self, type):
        # 断线回调，目前底层会自动处理断线，此接口已废弃
        print("OnDisconnect:", type)

    def OnRtnLoginWarn(self, pInfo):
        logger.error("Login fail !!:{}".format(pInfo))
        # print("OnRtnLoginWarn", pInfo)
        # 登陆失败 建议人工介入

    def OnRtnUserInfo(self, pInfo):
        # 登陆/重连成功
        # 此接口只有在登陆时被回调，如果盘中被回调表示网络异常，底层进行了重连重登操作
        # 程序会自动处理网络异常，并在网络恢复后自动重登，重登需要重新订阅行情。请在此接口被回调时尝试重新订阅行情。
        logger.warn("OnRtnUserInfo success !!:{}".format(pInfo))
        # print("OnRtnUserInfo", pInfo)

    def OnRspQryOrder(self, pRsp, nErr, isEnd):
        # 查询委托回调
        data = {}
        data["type"] = QUERY_ORDER
        data["pRsp"] = pRsp
        data["isEnd"] = isEnd
        data["nErr"] = nErr
        self.queue.put(data)
        logger.info("OnRspQryOrder:{} nErr:{} ,isEnd:{}".format(pRsp, nErr, isEnd))
        # print("OnRspQryOrder", pRsp, " nErr", nErr, " isEnd", isEnd)

    def OnRspQryMatch(self, pRsp, nErr, isEnd):
        # 查询成交回调
        pass
        #print("OnRspQryMatch", pRsp, " nErr", nErr, " isEnd", isEnd)

    def OnRspQryPosition(self, pRsp, nErr, isEnd):
        # 查询持仓回调
        pass
        #print("OnRspQryPositions推送底仓:", pRsp, " nErr", nErr, " isEnd", isEnd)

    def OnRspOrderInsert(self, pRsp, nErr):
        # 下单返回
        # 'nReqID': 2, 'szContractCode': '', 'nTradeType': 1, 'nOrderId': 0, 'nOrderPrice': 0, 'nOrderVol': 0, 'nStatus': 0, 'errDesc': '超过持仓数据异常!', 'szOrderOrigId': ''} ,nErr:10
        data = {}
        data["type"] = ACK_ORDER
        data["pRsp"] = pRsp
        data["nErr"] = nErr
        self.queue.put(data)
        logger.info("OnRspOrderInsert:{} ,nErr:{}".format(pRsp, nErr))

    def OnRspOrderDelete(self, pRsp, nErr):
        # 撤单返回
        print("OnRspOrderDelete", pRsp, " nErr", nErr)

    def OnRtnOrderStatusChangeNotice(self, pNotice):
        # 委托变化推送
        data = {}
        data["type"] = RSP_ORDER
        data["pRsp"] = pNotice
        data["nErr"] = 0
        self.queue.put(data)
        logger.info("OnRtnOrderStatusChangeNotice:{}".format(pNotice))

    def OnRtnOrderMatchNotice(self, pNotice):
        # 成交变化推送
        print("OnRtnOrderMatchNotice", pNotice)

    def OnRspSubQuote(self, pData):
        # 行情订阅返回
        pass
        #print("OnRspSubQuote", pData)

    def OnRtnTransaction(self, pTransaction):
        # 行情逐笔委托推送
        pass
        #print("OnRtnTransaction ", pTransaction)

    def OnRtnMarket(self, pMarket):
        # 行情市场快照推送
        pass
        #print("OnRtnMarket", pMarket)

    def OnRspQryAccountMaxEntrustCount(self, pFund, nErr):
        print("OnRspQryAccountMaxEntrustCount", pFund, " nErr", nErr)
