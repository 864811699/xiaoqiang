# encoding:utf-8
import threading
import time
from multiprocessing import Process, Queue
import os
import sys
import json
from multiprocessing import Queue
import copy

pwd = os.path.dirname(os.path.realpath(__file__))
sys.path.insert(0, pwd + '/../')

from api import MyCallBack
from api import QUERY_ORDER, ACK_ORDER, RSP_ORDER
from kdb import storeOrder, kdb
from kdb.util import logger


def readConfig(path):
    with open(path, 'r', encoding="UTF-8") as json_file:
        cfg = json.load(json_file)
    return cfg


status_dict = {
    -10: 6,  # 指令失败
    -9: 5,  # 撤单
    -8: 4,  # 全部成交
    0: 0,  # 订单初始化,还未上报
    1: 0,  # 已收到委托,但未受理
    2: 1,  # 正在排队处理中
    3: 3,  # 接收到撤单指令,但还未发出
    4: 3,  # 撤单指令已发出
    5: 5,  # 部分撤单
    6: 2,  # 部分成交
}


def getExchange(code):
    exchange = code + ".SZ"
    if code[0] == '6' or code[0] == '5':
        exchange = code + ".SH"
    return exchange


class trade:
    def __init__(self):
        # 接收kdb委托请求队列
        self.queue = Queue()
        self.rspOrderQueue = Queue()
        self.db = storeOrder.db()
        self.runing = True

    def initTradeApi(self, user, pwd):
        # 创建回调函数
        self.mspi = MyCallBack()

        # 创建/初始化API实例 （这里把创建API放到了回调函数的构造函数中，方便回调接口调用API）
        self.mapi = self.mspi.api

        # 登陆
        retLog = self.mapi.Login(user, pwd)
        logger.info("QuantPlus Login :ret:{}  user:{}  pwd:{}".format(retLog, user, pwd))  # 打印登录返回码

        # 必须在登陆成功后才能请求业务接口，demo使用condition variable来保证，仅供参考。
        with self.mspi.cv:
            while not self.mspi.isLogin:
                self.mspi.cv.wait()

    def storeOrder(self, order):
        # orderID
        # 检查当前order的orderID是否为空
        # 空,则不处理
        logger.info("store order[15]:[{}]:{}".format(order[15], order))
        if order[15] != "":
            # 不空则分2种
            old_order = self.order_dict.get(order[1])
            if old_order is not None:
                # 原始委托的orderID不存在则直接拼入
                if old_order[15] == "":
                    self.order_dict[order[15]] = order
                else:
                    # 原始委托的orderID存在
                    # 两orderID相同则不处理,不同则校验是否包含,不包含则拼接,包含不处理
                    if order[15] != old_order[15] and (order[15] not in old_order[15]):
                        self.order_dict[order[15]] = order
                        order[15] = old_order[15] + ",," + order[15]
                        logger.warn("发现 szOrderOrigId 新[{}] ,旧[{}] 不一致,建立多次映射,order[{}]".format(order[15],
                                                                                                             old_order[
                                                                                                                 15],
                                                                                                             order))

        # qid
        self.order_dict[order[1]] = order
        # entrustno
        self.order_dict[order[4]] = order
        self.db.storeOrder(order)
        self.rspOrderQueue.put(order)

    def newOrder(self, order):
        ret = self.mapi.OrderInsert(order)
        return ret

    def cancelOrder(self, req):
        # 使用服务器真实委托单号撤单
        # 下单回报/委托/成交接口在nErr为0(没有异常)的情况下，szOrderOrigId字段传送服务器真实委托单号
        # req = {"nReqId": genreq(), "nUserInt": 1, "nUserDouble": 1, "szUserStr": "", "szOrderOrigId": pRsp["szOrderOrigId"], "szOrderStreamId": "", "szProductCode": ""}
        self.mapi.OrderDelete(req)

    def qryOrder(self, accounts):
        # TODO
        if len(accounts) == 0:
            req = {"nReqId": self.maxId, "nUserInt": 1, "nUserDouble": 1, "szUserStr": "", "szContractCode": ""}
            ret = self.mapi.QryOrder(req)
            if ret != 0:
                logger.warn("send qry order fail,account is null", )
                return -1
            return ret
        for account in accounts:
            self.maxId += 1
            req = {"nReqId": self.maxId, "nUserInt": 1, "nUserDouble": 1, "szUserStr": account, "szContractCode": ""}
            ret = self.mapi.QryOrder(req)
            if ret != 0:
                logger.warn("send qry order fail,account:{}".format(account))

                return -1
        return ret

    def dealResponse(self):
        def dealRsp():
            try:
                while self.runing:
                    rsp = self.mspi.queue.get()
                    logger.info("receive from api,msg:{}".format(rsp))
                    try:
                        if rsp["nErr"] == 0 or (rsp["nErr"] != 0 and rsp["type"] == ACK_ORDER):
                            rspOrder = rsp["pRsp"]
                            orderid = rspOrder["szOrderOrigId"]
                            if rsp["type"] == QUERY_ORDER:
                                # data = {}
                                # data["type"] = QUERY_ORDER
                                # data["pRsp"] = pRsp
                                #       pRsp:{'nReqID': 10, 'szContractCode': '', 'nTradeType': 1, 'nOrderId': 18, 'nStatus': 0, 'nOrderPrice': 0.0, 'nOrderVol': 0, 'nDealedVol': 0, 'nDealedPrice': 0.0, 'szInsertTime': '', 'errDesc': 'no data'}
                                # data["isEnd"] = isEnd
                                # data["nErr"] = nErr

                                kdbOrder = self.order_dict.get(orderid, -1)
                                if kdbOrder == -1:
                                    logger.warn("qry order not found:{}".format(rspOrder))
                                    continue
                                kdbOrder = copy.copy(kdbOrder)
                            elif rsp["type"] == ACK_ORDER:
                                #   0      1        2           3           4           5           6           7           8       9           10      11          12      13          14
                                # ['sym', 'qid', 'accountname', 'time', 'entrustno', 'stockcode', 'askprice', 'askvol', 'bidprice','bidvol','withdraw', 'status', 'note', 'reqtype', 'params']
                                # 'nReqID': 2, 'szContractCode': '', 'nTradeType': 1, 'nOrderId': 0, 'nOrderPrice': 0, 'nOrderVol': 0, 'nStatus': 0, 'errDesc': '超过持仓数据异常!', 'szOrderOrigId': ''} ,nErr:10
                                # receive from api,msg:{'type': 'ack Order', 'pRsp': {'nReqID': 40, 'szContractCode': '600111.SH', 'nTradeType': 1, 'nOrderId': 27, 'nOrderPrice': 192800, 'nOrderVol': 10000, 'nStatus': 0, 'errDesc': '36108', 'szOrderOrigId': '36108'}, 'nErr': 0}
                                entrustno = rspOrder["nReqID"]
                                kdbOrder = self.order_dict.get(entrustno, -1)
                                kdbOrder = copy.copy(kdbOrder)
                                if kdbOrder == -1:
                                    logger.warn("ack order not found:{}".format(rspOrder))
                                    continue
                                if rsp["nErr"] != 0:
                                    kdbOrder[11] = 6
                                    kdbOrder[12] = rspOrder["errDesc"]
                                    self.storeOrder(kdbOrder)
                                    continue

                                kdbOrder[15] = orderid
                                kdbOrder[11] = 1
                                self.storeOrder(kdbOrder)
                                continue
                            elif rsp["type"] == RSP_ORDER:
                                kdbOrder = self.order_dict.get(orderid, -1)
                                kdbOrder = copy.copy(kdbOrder)
                                if kdbOrder == -1:
                                    logger.warn(
                                        "api response orderid[{}] not found:{}".format(orderid, rspOrder))

                                    continue
                            else:
                                logger.warn("receive from api error:{}".format(rsp))

                            status = rspOrder["nStatus"]
                            status = status_dict.get(status, -1)
                            if status == -1:
                                logger.warn("received from api,msg status incorrect:{}".format(rspOrder))
                                continue
                            bidvol = rspOrder["nDealedVol"]
                            bidPx10000 = rspOrder["nDealedPrice"]

                            note = rspOrder["errDesc"]
                            if bidvol >= kdbOrder[9] and status >> kdbOrder[11]:
                                kdbOrder[8] = bidPx10000 / 10000.0
                                kdbOrder[9] = bidvol
                                if kdbOrder[7] < 0:
                                    kdbOrder[9] = -kdbOrder[9]
                                # 撤单数量
                                kdbOrder[11] = status
                                if status == 5:
                                    kdbOrder[10] = abs(kdbOrder[7]) - abs(bidvol)
                                if status == 6:
                                    kdbOrder[12] = note
                                self.storeOrder(kdbOrder)
                            if status == 0 and "ITG_" in note and rspOrder['szMemo'] != "":
                                kdbOrder[11] = 6
                                kdbOrder[12] = note
                                self.storeOrder(kdbOrder)
                    except Exception as e:
                        logger.error("dealResponse error: {} ,rsp[{}]".format(e, rsp))
            except Exception as e:
                logger.error("未知错误 dealRsp: {}".format(e))
        self.t1 = threading.Thread(target=dealRsp)
        self.t1.start()

    def reboot(self, accounts):
        self.order_dict = self.db.loadAllOrders()
        self.db.openDB()
        self.maxId = self.db.getMaxId()
        ret = self.qryOrder(accounts)
        if ret != 0:
            logger.warn("reboot qry fail!!!!!")
            sys.exit()


def main():
    logger.warn("!!!!!!!!!!!!main run")
    cfg = readConfig("./etc/config.json")
    t = trade()
    p = Process(target=kdb.StartKdb, args=(cfg["kdb"], t.queue, t.rspOrderQueue))
    p.start()
    logger.warn("!!!!!!!!!!!!initTradeApi run")
    t.initTradeApi(cfg['api']["user"], cfg['api']["pwd"])
    logger.warn("!!!!!!!!!!!!dealResponse run")
    t.dealResponse()
    t.reboot(cfg['api']['fundAccounts'])

    # 接收委托,并下指令
    try:
        while True:
            order = t.queue.get()
            t.maxId += 1
            order[4] = t.maxId
            order[1] = order[1].decode("utf-8")
            order[0] = order[0].decode("utf-8")
            order[2] = order[2].decode("utf-8")
            order[5] = order[5].decode("utf-8")
            order[12] = order[12].decode("utf-8")
            order[14] = order[14].decode("utf-8")
            order[15] = order[15].decode("utf-8")
            # 解析
            #   'sym',   'qid',   'accountname'   'time','entrustno','stockcode', 'askprice','askvol','bidprice', 'bidvol','withdraw', 'status','note', 'reqtype','params'
            #  [b'test1'^ b'cddeceef-'^ b'test'^ 8820.627192^ 0^      b'601818'^ 3.^           200^     0.^        0^        0^          0^      b''^     -1^       b'']
            #     0           1            2       3          4        5         6             7        8          9         10          11      12       13         14
            if order[11] == 0:
                # 委托   status=0 委托,
                # reqtype  -1 交易 |  1 还券
                # 添加 orderid
                t.storeOrder(order)
                if order[13] == -1:
                    # 下单 szProductCode从查询持仓中获取，如果填写空字符或者req的dict里不存在这个key则从conf.json中的Trade.ProductCode字段获取
                    # nReqId： 外部编号，唯一
                    # nUserInt：固定1
                    # nUserDouble：固定1
                    # szUserStr：交易帐号
                    # szContractCode：股票编号，例 600000
                    # szProductCode：产品编号，可不变，使用conf.json中配置
                    # nOrderPrice：委托价格，需*10000
                    # nOrderVol：委托数量
                    # nTradeType：交易类型  2：买 1：卖

                    side = 2
                    if order[7] < 0:
                        side = 1
                    qty = abs(order[7])
                    req = {"nReqId": order[4], "nUserInt": 1, "nUserDouble": 1,
                           "szUserStr": cfg["accountInfos"][order[0]]["szUserStr"], "szContractCode": getExchange(order[5]),
                           "szProductCode": cfg["accountInfos"][order[0]]["szProductCode"], "nOrderPrice": order[6] * 10000,
                           "nOrderVol": qty, "nTradeType": side}
                    t.newOrder(req)
                    logger.info("server exec order:{}".format(order))
                else:
                    continue
            elif order[11] == 3:
                # 撤单 status =3 撤单
                # 撤单操作  szProductCode和下单时的对应，如果填写空字符或者req的dict里不存在这个key则从conf.json中的Trade.ProductCode字段获取
                orderFull = t.order_dict.get(order[1])
                if orderFull is not None:
                    orderids = orderFull[15].split(",,")
                    for orderid in orderids:
                        req = {"nReqId": t.maxId, "nUserInt": 1, "nUserDouble": 1,
                               "szUserStr": cfg["accountInfos"][order[0]]["szUserStr"], "szOrderOrigId": orderid,
                               "szOrderStreamId": "", "szProductCode": cfg["accountInfos"][order[0]]["szProductCode"]}
                        t.cancelOrder(req)
                    logger.info("server cancel order:{}".format(order))
            else:
                continue
    except Exception as e:
        logger.error("未知错误main: {}".format(e))

if __name__ == '__main__':
    main()
