import os
import sys
import threading

import numpy as np
from qpython import qconnection as qp
from qpython.qcollection import qlist, qtable
from qpython.qtype import QException, QINT, QFLOAT, QSYMBOL, QDATETIME, QDOUBLE
import queue as thQueue
# from multiprocessing import Queue
from datetime import datetime
from numpy.lib.recfunctions import append_fields

pwd = os.path.dirname(os.path.realpath(__file__))
sys.path.insert(0, pwd + '/../')

from .util import logger


class K:

    def __init__(self, host, port, username, password, reqOrderQueue, rspOrderQueue):
        self.__queue = thQueue.Queue()
        self.reqOrderQueue = reqOrderQueue
        self.rspOrderQueue = rspOrderQueue
        self.host = host
        self.port = port
        self.user = username
        self.pwd = password
        logger.info("init K ,host:{},port:{}".format(host, port))

    def getQ(self):
        return self.__q

    def run(self):
        self.__runing = True
        logger.info("kdb set mode is Runing!!!")

    def is_runing(self):
        return self.__runing

    def stop(self):
        self.__runing = False
        logger.info("kdb set mode is Stop!!!")

    def connect(self):
        # 连接到KDB服务器
        self.__q = qp.QConnection(self.host, self.port, self.user, self.pwd)
        logger.info("kdb connect success!!!")

    def open(self):
        # 打开链接
        self.__q.open()
        logger.info("kdb open success!!!")

    def subscribe(self, table, sym):
        # 订阅A表更新
        self.__q.sendAsync('.u.sub', np.string_(table), np.string_(sym))
        logger.info("kdb subscribe table:{} sym:{} success!!!".format(table, sym))

    def receive(self):
        """ msg:
        [b'upd', b'request', rec.array([(b'test1', b'cddeceef-9ee9-3847-9172-3e3d7ab39b26', b'test', 8820.627192, 0, b'601818', 3., 200, 0., 0, 0, 0, b'', -1, b''),
           (b'test1', b'97911f28-7bfa-7efa-6c10-4fb48ceeb4f1', b'test', 8820.627192, 0, b'601818', 3., 200, 0., 0, 0, 0, b'', -1, b'')],
          dtype=[('sym', 'S5'), ('qid', 'S36'), ('accountname', 'S4'), ('time', '<f8'), ('entrustno', '<i4'), ('stockcode', 'S6'), ('askprice', '<f8'), ('askvol', '<i4'), ('bidprice', '<f8'), ('bidvol', '<i4'), ('withdraw', '<i4'), ('status', '<i4'), ('note', 'S1'), ('reqtype', '<i4'), ('params', 'S1')])]
        """

        def receiveFromKdb():
            while self.__runing:
                # 接收数据
                msg = self.__q.receive(data_only=True)
                # print(msg)
                # 将数据放入队列
                self.__queue.put(msg)
                # logger.info("kdb received message: {}".format(msg))

        t = threading.Thread(target=receiveFromKdb)
        t.start()

    def publisher(self, table, msg):
        """       'sym',   'qid',   'accountname'   'time','entrustno','stockcode', 'askprice','askvol','bidprice', 'bidvol','withdraw', 'status','note', 'reqtype','params'
        msg 格式: [b'test1', b'cddeceef-', b'test', 8820.627192, 0,      b'601818', 3.,           200,     0.,        0,        0,          0,      b'',     -1,       b'']
        """
        data = qtable(
            ['sym', 'qid', 'accountname', 'time', 'entrustno', 'stockcode', 'askprice', 'askvol', 'bidprice', 'bidvol',
             'withdraw', 'status', 'note', 'reqtype', 'params'],
            [[msg[0].encode('utf-8')], [msg[1].encode('utf-8')], [msg[2].encode('utf-8')],
             [np.datetime64(np.datetime64(datetime.now()), 'ms')],
             [msg[4]], [msg[5].encode('utf-8')], [msg[6]], [msg[7]], [msg[8]], [msg[9]], [msg[10]], [msg[11]],
             [msg[12].encode('utf-8')], [msg[13]],
             [msg[14].encode('utf-8')]],
            sym=QSYMBOL, qid=QSYMBOL, accountname=QSYMBOL,
            time=QDATETIME, entrustno=QINT, stockcode=QSYMBOL,
            askprice=QDOUBLE, askvol=QINT, bidprice=QDOUBLE, bidvol=QINT,
            withdraw=QINT, status=QINT, note=QSYMBOL, reqtype=QINT,
            params=QSYMBOL,
        )
        self.__q.sendAsync('upd0', np.string_(table), data)
        logger.info("push order to kdb :{}".format(msg))

    def close(self):
        self.__q.close()
        logger.info("kdb close success!!!")

    def getQueue(self):
        return self.__queue

    def query(self, table, syms):
        symListStr = "("
        for sym in syms:
            symListStr = symListStr + "`" + sym
        symListStr += ")"
        sql = 'select from `{} where sym in {},status<4'.format(table, symListStr)
        logger.info("qry kdb sql={}".format(sql))
        data = self.__q.sendSync(sql)

    def producer(self, requestTab):
        def getReqOrder():
            logger.info("kdb start load order from db {}!!!!!".format(requestTab))
            try:
                while self.is_runing():
                    rsp = self.__queue.get()
                    #  [b'upd', b'request', rec.array([])
                    if rsp[1] == requestTab.encode():
                        try:
                            spaceStr = np.array([b''])
                            orders = rsp[2]
                            fields = ['orderid', 'note', 'params']
                            for field in fields:
                                if field not in orders.dtype.names:
                                    orders = append_fields(orders, field, spaceStr, usemask=False)
                            # orders = append_fields(orders, 'orderid', spaceStr, usemask=False)
                            # orders = append_fields(orders, 'note', spaceStr, usemask=False)
                            # orders = append_fields(orders, 'params', spaceStr, usemask=False)
                            orders = orders[
                                ['sym', 'qid', 'accountname', 'time', 'entrustno', 'stockcode', 'askprice', 'askvol',
                                 'bidprice', 'bidvol',
                                 'withdraw', 'status', 'note', 'reqtype', 'params', 'orderid']]
                            for row in orders:
                                # orderid = np.array([""])
                                # row = append_fields(row, 'orderid', orderid, usemask=False)
                                data = list(row.tolist())
                                if data[11] == 3:
                                    logger.info("kdb load cancel:{}".format(data))
                                else:
                                    logger.info("kdb load order:{}".format(data))
                                self.reqOrderQueue.put(data)
                        except Exception as e:
                            logger.error("decode kdb msg fail,error={} ,msg ={}".format(e, rsp))
            except Exception as e:
                logger.error("未知错误 kdb producer: {}".format(e))

        threading.Thread(target=getReqOrder).start()

    def consumer(self, responseTab):
        def pushRspOrder():
            logger.info("kdb start push order to db {}!!!!!".format(responseTab))
            try:
                while self.is_runing():
                    rsp = self.rspOrderQueue.get()
                    logger.info("kdb push order:{}".format(rsp))
                    self.publisher(responseTab, rsp)
            except Exception as e:
                logger.error("未知错误 kdb consumer: {}".format(e))
        threading.Thread(target=pushRspOrder).start()


def StartKdb(cfg, reqOrderQueue, rspOrderQueue):
    k = K(cfg["kdbIP"], cfg["kdbPort"], cfg["user"], cfg["pwd"], reqOrderQueue, rspOrderQueue)
    k.connect()
    k.open()
    k.run()
    k.receive()
    k.subscribe("HeartBeat", [])
    k.subscribe(cfg["requestTab"], cfg["accounts"])
    k.producer(cfg["requestTab"])
    k.consumer(cfg["responseTab"])
