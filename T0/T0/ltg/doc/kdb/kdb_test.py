import os
import sys
import numpy as np
from qpython import qconnection as qp
from qpython.qcollection import qlist
from qpython.qcollection import qtable
from qpython.qtype import QException, QINT, QFLOAT, QSYMBOL, QDATETIME,QSTRING,QDOUBLE
from datetime import datetime, timedelta, timezone
import qpython.qtype
import uuid

q = qp.QConnection('58.216.145.86', 40000, "test", "pass")
q.open()

# subscribe
# q.sendAsync('.u.sub', np.string_('request'), np.string_('test1'))
# msg = q.receive(data_only=True)

# msg = [b'test1', b'cddeceef', b'test', np.datetime64(np.datetime64(datetime.now()) , 'ms'), 0, b'601818', 3., 200, 0., 0, 0, 0, b'', -1, b'']
# data = [QTable([msg[0]], qtype=QSYMBOL),  # sym
#         qlist([msg[1]], qtype=QSYMBOL),  # qid
#         qlist([msg[2]], qtype=QSYMBOL),  # accountname
#         qlist([msg[3]], qtype=qpython.qtype.QDATETIME),  # time
#         qlist([msg[4]], qtype=QINT_LIST),  # entrustno
#         qlist([msg[5]], qtype=QSYMBOL_LIST),  # stockcode
#         qlist([msg[6]], qtype=QFLOAT_LIST),  # askprice
#         qlist([msg[7]], qtype=QINT_LIST),  # askvol
#         qlist([msg[8]], qtype=QFLOAT_LIST),  # bidprice
#         qlist([msg[9]], qtype=QINT_LIST),  # bidvol
#         qlist([msg[10]], qtype=QINT_LIST),  # withdraw
#         qlist([msg[11]], qtype=QINT_LIST),  # status
#         qlist([msg[12]], qtype=QSYMBOL_LIST),  # note
#         qlist([msg[13]], qtype=QINT_LIST),  # reqtype
#         qlist([msg[14]], qtype=QSYMBOL_LIST), ]  # params
# print(data)

msg=['test1', 'cd17ee98-c089-10a3-8992-d437a566f081', 'test', '2024-02-29T13:39:00.786', 2, '601818', 3., 200, 0., 0, 0, 1, '', -1, '']
data = qtable(
            ['sym', 'qid', 'accountname', 'time', 'entrustno', 'stockcode', 'askprice', 'askvol', 'bidprice', 'bidvol',
             'withdraw', 'status', 'note', 'reqtype', 'params'],
            [[msg[0]], [msg[1]], [msg[2]], [np.datetime64(np.datetime64(datetime.now()), 'ms')],
             [msg[4]], [msg[5]], [msg[6]], [msg[7]], [msg[8]], [msg[9]], [msg[10]], [msg[11]], [msg[12]], [msg[13]],
             [msg[14]]],
            sym=QSYMBOL, qid=QSYMBOL, accountname=QSYMBOL,
            time=QDATETIME, entrustno=QINT, stockcode=QSYMBOL,
            askprice=QDOUBLE, askvol=QINT, bidprice=QDOUBLE, bidvol=QINT,
            withdraw=QINT, status=QINT, note=QSYMBOL, reqtype=QINT,
            params=QSYMBOL,
        )
#cf59ee6a-e1c6-ea04-fac4-4ce0bf20f1c1
#1111111111111111111111111111111111111
data = qtable(
            ['sym', 'qid', 'accountname','time','entrustno', 'stockcode', 'askprice'],
            [[msg[0]], ["asdfaf111"], [msg[2]] ,[np.datetime64(np.datetime64(datetime.now()), 'ms')],[msg[4]], [msg[5]], [msg[6]]],
            sym=QSYMBOL, qid=QSYMBOL, accountname=QSYMBOL, time=qpython.qtype.QDATETIME,entrustno=qpython.qtype.QINT, stockcode=qpython.qtype.QSYMBOL,askprice=qpython.qtype.QDOUBLE
        )

q.sendAsync('updTest', np.string_("response"), data,single_char_strings = True)

data = qtable(
    ['sym', 'time', 'no'],
    [['test2'], [np.datetime64(np.datetime64(datetime.now()), 'ms')], [1]],
    sym=qpython.qtype.QSYMBOL, time=qpython.qtype.QDATETIME, no=qpython.qtype.QINT
)
q.sendSync('updTest', np.string_("testPush"), data)
q.sendSync('updTest', np.string_("testPush2"), data)

dtype=[('sym', 'S5'), ('qid', 'S36'), ('accountname', 'S4'), ('time', '<M8[ms]'), ('entrustno', '<i4'), ('stockcode', 'S6'), ('askprice', '<f4'), ('askvol', '<i4'), ('bidprice', '<f4'), ('bidvol', '<i4'), ('withdraw', '<i4'), ('status', '<i4'), ('note', 'S1'), ('reqtype', '<i4'), ('params', 'S1')])
dtype=[('sym', 'S5'), ('qid', 'S36'), ('accountname', 'S4'), ('time', '<f8'), ('entrustno', '<i4'), ('stockcode', 'S6'), ('askprice', '<f8'), ('askvol', '<i4'), ('bidprice', '<f8'), ('bidvol', '<i4'), ('withdraw', '<i4'), ('status', '<i4'), ('note', 'S1'), ('reqtype', '<i4'), ('params', 'S1')])]
