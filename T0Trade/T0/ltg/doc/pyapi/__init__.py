from QuantBaseApi import QuantCallBack,PT_QuantApi_Python36
import traceback

PyApi = PT_QuantApi_Python36

def excsafe(func):
    def wrapper(*args, **kwargs):
        try:
            func(*args, **kwargs)
        except Exception as ex:
            traceback.print_exc()
            raise ex
    
    return wrapper

class DataCallBack(QuantCallBack):

    @excsafe
    def RtnConnect(self):
        self.OnConnect(self.rspdict['nSrvType'])

    @excsafe
    def RtnDisconnect(self):
        self.OnDisconnect(0)

    @excsafe
    def RtnUserInfo(self):
        self.OnRtnUserInfo(self.rspdict)

    @excsafe
    def RtnLoginWarn(self):
        self.OnRtnLoginWarn(self.rspdict)

    @excsafe
    def RtnRspQryOrder(self):
        if self.errorid == 0 and len(self.rspdict['errDesc'].split(':')) == 2:
            extractedval = self.rspdict['errDesc'].split(':')
            self.rspdict['szOrderOrigId'] = extractedval[0];
            self.rspdict['szMemo'] = extractedval[1];
        self.OnRspQryOrder(self.rspdict, self.errorid, self.isend)

    @excsafe
    def RtnRspQryMatch(self):
        if self.errorid == 0:
            self.rspdict['szOrderOrigId'] = self.rspdict['errDesc']
            
        self.OnRspQryMatch(self.rspdict, self.errorid, self.isend)
        
    @excsafe
    def RtnRspQryPosition(self):
        if self.errorid == 0 and len(self.rspdict['errDesc'].split(':')) == 2:
            extractedval = self.rspdict['errDesc'].split(':')
            self.rspdict['szProductCode'] = extractedval[0];
            self.rspdict['nPrePrice'] = extractedval[1];
            
        self.OnRspQryPosition(self.rspdict, self.errorid, self.isend)

    @excsafe
    def RtnRspOrderInsert(self):
        self.rspdict['szOrderOrigId'] = ''
        if self.errorid == 0:
            self.rspdict['szOrderOrigId'] = self.rspdict['errDesc']
        self.OnRspOrderInsert(self.rspdict, self.errorid)

    @excsafe
    def RtnRspOrderDelete(self):
        self.OnRspOrderDelete(self.rspdict, self.errorid)

    @excsafe
    def RtnOrderStatusChangeNotice(self):
        if self.errorid == 0 and len(self.rspdict['errDesc'].split(':')) == 2:
            extractedval = self.rspdict['errDesc'].split(':')
            self.rspdict['szOrderOrigId'] = extractedval[0];
            self.rspdict['szMemo'] = extractedval[1];
        self.OnRtnOrderStatusChangeNotice(self.rspdict)

    @excsafe
    def RtnOrderMatchNotice(self):
        if self.errorid == 0:
            self.rspdict['szOrderOrigId'] = self.rspdict['errDesc']
        self.OnRtnOrderMatchNotice(self.rspdict)

    @excsafe
    def RtnRspSubQuote(self):
        self.OnRspSubQuote(self.rspdict)

    @excsafe
    def RtnMarket(self):
        self.OnRtnMarket(self.rspdict)

    @excsafe
    def RtnTransaction(self):
        self.OnRtnTransaction(self.rspdict)

    @excsafe		
    def RtnQryAccountMaxEntrustCount(self):
        self.OnRspQryAccountMaxEntrustCount(self.rspdict, self.errorid)
		
    