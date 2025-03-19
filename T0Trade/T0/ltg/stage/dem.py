#encoding:utf-8

#from QuantBaseApi import QuantCallBack,PT_QuantApi_Python36
from pyapi import PyApi, DataCallBack
from datetime import datetime, date, timedelta
import time
import threading

reqcnt = 0

def genreq():
    global reqcnt
    reqcnt += 1
    return int(round(time.time()) + reqcnt)

class MyCallBack(DataCallBack):
    def __init__(self):
        super(MyCallBack, self).__init__()
        
        #创建api实例，可以写在回调之外，这里放在回调class内只是方便回调函数调用api
        self.api = PyApi.PT_QuantApi(self, True, "Td_Real", "MD_Real")
        
        #初始化API必须在创建之后立刻执行
        PyApi.Init()
		
        self.cv = threading.Condition()
        self.isLogin = False
    
    def OnConnect(self, type):
        #if type == 1:   #0为实时服务器 1为历史服务器 2为缓存服务器  推荐在0.1.2三种行情服务器全部连接成功后做行情请求能保证全时段行情稳定回调        10为交易服务器，连接成功之后才可做交易业务请求
        print("OnConnnect:", type)
        with self.cv:
            self.isLogin = True
            self.cv.notify()
        
    def OnDisconnect(self, type):
        #断线回调，目前底层会自动处理断线，此接口已废弃
        print("OnDisconnect:", type)
        
        
    def OnRtnLoginWarn(self, pInfo):
        print ("OnRtnLoginWarn", pInfo)
        #登陆失败 建议人工介入
        
        
    def OnRtnUserInfo(self, pInfo):
        #登陆/重连成功
        #此接口只有在登陆时被回调，如果盘中被回调表示网络异常，底层进行了重连重登操作
        #程序会自动处理网络异常，并在网络恢复后自动重登，重登需要重新订阅行情。请在此接口被回调时尝试重新订阅行情。
        print ("OnRtnUserInfo", pInfo)
        
        

    def OnRspQryOrder(self, pRsp, nErr, isEnd):
        #查询委托回调
        # !!!!!!!!!!!!!!!!!OnRspQryOrder {'nReqID': 10, 'szContractCode': '', 'nTradeType': 1, 'nOrderId': 0, 'nStatus': 0, 'nOrderPrice': 0.0, 'nOrderVol': 0, 'nDealedVol': 0, 'nDealedPrice': 0.0, 'szInsertTime': '', 'errDesc': 'no data'}  nErr -999  isEnd True
        print ("!!!!!!!!!!!!!!!!!OnRspQryOrder", pRsp, " nErr", nErr, " isEnd", isEnd)
        

    def OnRspQryMatch(self, pRsp, nErr, isEnd):
        #查询成交回调
        print ("OnRspQryMatch", pRsp, " nErr", nErr, " isEnd", isEnd)
        

    def OnRspQryPosition(self, pRsp, nErr, isEnd):
        #查询持仓回调
        print ("OnRspQryPositions推送底仓:", pRsp, " nErr", nErr, " isEnd", isEnd)

    def OnRspOrderInsert(self, pRsp, nErr):
        #下单返回
        #!!!!!!!!!下单返回OnRspOrderInsert {'nReqID': 104429, 'szContractCode': '1061.SZ', 'nTradeType': 2, 'nOrderId': 0, 'nOrderPrice': 41500, 'nOrderVol': 100, 'nStatus': 0, 'errDesc': '超过持仓数据异常!', 'szOrderOrigId': ''} 10
        print("!!!!!!!!!下单返回OnRspOrderInsert", pRsp, nErr)
        
        #撤单操作  szProductCode和下单时的对应，如果填写空字符或者req的dict里不存在这个key则从conf.json中的Trade.ProductCode字段获取      
        #req = {"nReqId": genreq(), "nUserInt": 1, "nUserDouble": 1, "szUserStr": "", "nOrderId": pRsp["nOrderId"], "szOrderStreamId": "", "szProductCode": ""}
        
        #使用服务器真实委托单号撤单
        #下单回报/委托/成交接口在nErr为0(没有异常)的情况下，szOrderOrigId字段传送服务器真实委托单号
        #req = {"nReqId": genreq(), "nUserInt": 1, "nUserDouble": 1, "szUserStr": "", "szOrderOrigId": pRsp["szOrderOrigId"], "szOrderStreamId": "", "szProductCode": ""}
        #self.api.OrderDelete(req)
        

    def OnRspOrderDelete(self, pRsp, nErr):
        #撤单返回
        print("!!!!!!!!!撤单返回OnRspOrderDelete", pRsp, " nErr", nErr)
        

    def OnRtnOrderStatusChangeNotice(self, pNotice):
        #委托变化推送
        print ("!!!!!!!!!!!!!!委托变化推送OnRtnOrderStatusChangeNotice",  pNotice)
        

    def OnRtnOrderMatchNotice(self, pNotice):
        pass
        #成交变化推送
        print ("!!!!!!!!!!!!!!OnRtnOrderMatchNotice", pNotice)
        

    def OnRspSubQuote(self, pData):
        pass
        #行情订阅返回
        print ("OnRspSubQuote", pData)
        

    def OnRtnTransaction(self, pTransaction):
        pass
        #行情逐笔委托推送
        print ("OnRtnTransaction ", pTransaction)
        

    def OnRtnMarket(self, pMarket):
        pass
        #行情市场快照推送
        print("OnRtnMarket", pMarket)
		
    def OnRspQryAccountMaxEntrustCount(self, pFund, nErr):
        pass
        # print("OnRspQryAccountMaxEntrustCount", pFund, " nErr", nErr)
        



def main():
    #创建回调函数
    mspi = MyCallBack()
    
    #创建/初始化API实例 （这里把创建API放到了回调函数的构造函数中，方便回调接口调用API）
    mapi = mspi.api
    
    #登陆
    retLog = mapi.Login("10016","666666")
    print("QuantPlus Login :",retLog)  #打印登录返回码

    #必须在登陆成功后才能请求业务接口，demo使用condition variable来保证，仅供参考。
    with mspi.cv:
        while not mspi.isLogin:
            mspi.cv.wait()
    print("!!!!!!!!!!!!登入成功")
    #nReqId请保证当日唯一，请自行生成，此处仅作示范
    
    #查询委托
    req = {"nReqId":10,"nUserInt":1,"nUserDouble":1,"szUserStr":"","szContractCode":""}
    ret = mapi.QryOrder(req)
    #print("查询委托请求结果:",ret)

    #查询持仓
    #req = {"nReqId": genreq(), "nUserInt": 1, "nUserDouble": 1, "szUserStr": "", "szContractCode": ""}
    #ret = mapi.QryPosition(req)
	
	#查询成交
    #req = {"nReqId": genreq(), "nUserInt": 1, "nUserDouble": 1, "szUserStr": "", "szContractCode": ""}
    #ret = mapi.QryMatch(req)
    
    #订阅行情
    #第二个字段market代表订阅市场快照（通过OnRtnMarket推送），transaction代表订阅逐笔成交（通过OnRtnTransaction推送）订阅结果通过OnRspSubQuote获取
    #ret = mapi.ReqSubQuote(genreq(), ["market","transaction"], [""], ["000403.SZ", "600030.SH", "000997.SZ", "600029.SH"], "", "")

    #退订行情 退订结果通过OnRspSubQuote获取
    #ret = mapi.ReqSubQuote(genreq(), ["unsub"], [""], ["000403.SZ", "600030.SH"], "", "")


    #查询资金接口  仅供投顾类型用户使用
    #req = {"nReqId": genreq()}
    #ret = mapi.QryAccountMaxEntrustCount(req)
    print("!!!!!!!!!!!!准备下单")
    def tfunc():

            print("!!!!!!!!!!!!开始下单")
            hour = "{:0>2d}".format(time.localtime(time.time()).tm_hour)
            minute = "{:0>2d}".format(time.localtime(time.time()).tm_min)
            sec = "{:0>2d}".format(time.localtime(time.time()).tm_sec)
            orderid = int(str(hour) + str(minute) + str(sec))
        
            #下单 szProductCode从查询持仓中获取，如果填写空字符或者req的dict里不存在这个key则从conf.json中的Trade.ProductCode字段获取
            stock = "000617.SZ"
            # 这边价格需要 * 10000...
            req = {"nReqId": 654322, "nUserInt": 1, "nUserDouble": 1, "szUserStr": "", "szContractCode": stock, "szProductCode": "200002","nOrderPrice": 6.3 * 10000, "nOrderVol": 100, "nTradeType": 1}
            ret = mapi.OrderInsert(req)
            print("!!!!!! 下单结果:",ret)
            time.sleep(1000)

            #撤单范例 见MyCallBack的OnRspOrderInsert方法

    
    # thread1 = threading.Thread(target=tfunc)
    # thread1.start()

    tfunc()
    
    #主程序挂在api.Run()函数上
    while True:
        #mapi.Run()
        pass
		
    print("exit")
        
if __name__ == '__main__':
    main()
