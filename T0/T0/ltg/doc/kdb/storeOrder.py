import os
import datetime
import json

class db:
    def __init__(self):
        self.maxId = 0
        self.filePath = "./data/{}.txt".format(datetime.datetime.now().strftime('%Y%m%d'))

    def openDB(self):
        os.makedirs('./data', exist_ok=True)
        self.f = open(self.filePath, 'a+')

    def closeDB(self):
        self.f.close()

    def __readDB(self):
        try:
            with open(self.filePath, 'r') as f:
                lines = f.readlines()
            return lines
        except FileNotFoundError:
            print("db file:{} not found ".format(self.filePath))

    def loadAllOrders(self):
        lines = self.__readDB()
        dict = {}
        if lines is None:
            return dict
        #   'sym',   'qid',   'accountname'   'time','entrustno','stockcode', 'askprice','askvol','bidprice', 'bidvol','withdraw', 'status','note', 'reqtype','params','orderid'
        #  [b'test1'^ b'cddeceef-'^ b'test'^ 8820.627192^ 0^      b'601818'^ 3.^           200^     0.^        0^        0^          0^      b''^     -1^       b'']
        #     0           1            2       3          4        5         6             7        8          9         10          11      12       13         14
        for line in lines:
            fields = json.loads(line)

            # fields = line.split("^")
            # if len(fields) != 16:
            #     print("存储数据异常:", line)
            #     continue
            # fields[4] = int(fields[4])
            # fields[6] = float(fields[6])
            # fields[7] = int(fields[7])
            # fields[8] = float(fields[8])
            # fields[9] = int(fields[9])
            # fields[10] = int(fields[10])
            # fields[11] = int(fields[11])
            # fields[13] = int(fields[13])
            if fields[4] > self.maxId:
                self.maxId = fields[4]
            if fields[1] in dict:
                oldOrder = dict[fields[1]]
                if fields[9] >= oldOrder[9] and fields[11] >= oldOrder[11]:
                    dict[fields[1]] = fields
                    dict[fields[4]] = fields
                    orderids = fields[15].split(",,")
                    for orderid in orderids:
                        dict[orderid] = fields
            else:
                dict[fields[1]] = fields
                dict[fields[4]] = fields
                orderids = fields[15].split(",,")
                for orderid in orderids:
                    dict[orderid] = fields
        return dict

    def getMaxId(self):
        return self.maxId

    def storeOrder(self, order):
        jsonData = json.dumps(order)
        self.f.write(jsonData + '\n')
        self.f.flush()
