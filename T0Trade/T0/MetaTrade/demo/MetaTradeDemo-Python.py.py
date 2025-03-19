# -*- coding: utf-8 -*-

"""

此DLL接口为32位,用python解释器调用此接口的时候必须要用python32位解释器

1.用32位系统的直接安装32位python解释器
2.用64位系统的要设置兼容32位python解释器

"""

from ctypes import *


# 填入相对路径加载DLL
Dll = windll.LoadLibrary(".\\MetaTrade.dll")

# DLL初始化,只调用一次
# 初始化返回值<=0,为初始化失败,同时不能够掉用其他功能否则会出错
Dll.Init()

# 为返回结果和返回错误分配空间,结果信息空间为1024*1024字节,错误信息空间为256字节
# 要使用批量功能应创建多个Result和ErrorInfo空间
Result = create_string_buffer(1024 * 1024)
ErrorInfo = create_string_buffer(256)
Result1 = create_string_buffer(1024 * 1024)
ErrorInfo1 = create_string_buffer(256)

# 登入交易账户
# 参数是str类型的必须转换成bytes类型,int类型不用转换,某些参数可以为空
Ip = b'111.222.333.444'
Port = 1234
Version = b''
YyBid = 0
Account = b'123456789.C'
TradeAccount = b'123456789'
JyPassword = b'123456'
TxPassword = b''

# 按位置传参,注意C的类型
# 要使用批量功能应要创建多个ClientId
ClientId = Dll.Logon(c_char_p(Ip), c_int(Port), c_char_p(Version),
                     c_int(YyBid), c_char_p(Account), c_char_p(TradeAccount),
                     c_char_p(JyPassword), c_char_p(TxPassword), ErrorInfo)

# 登入成功ClientId返回结果为1,如果返回结果为-1表示登入失败,可以调ErrorInfo查看错误信息
print(ErrorInfo1.value.decode('gbk'))

# 查询各类交易数据
# ErrorInfo为空则返回结果成功请调Result查看结果,反之返回结果失败请调ErrorInfo查看错误信息
Dll.QueryData(ClientId, c_int(Category), Result, ErrorInfo)
print(Result.value.decode('gbk'))

# 单账户批量查询各类交易数据
# 要批量查询多类数据Result和ErrorInfo必须建立相同数量,可用for循环建立,Count参数必须等于查询数据列表的长度,反之内存报错
CategoryList = [Category1, Category2]
CategoryArray = (c_int * len(CategoryList))(*CategoryList)
ResultList = [cast(Result1, c_char_p), cast(Result2, c_char_p)]
ResultArray = (c_char_p * len(ResultList))(*ResultList)
ErrorInfoList = [cast(ErrorInfo1, c_char_p), cast(ErrorInfo2, c_char_p)]
ErrorInfoArray = (c_char_p * len(ErrorInfoList))(*ErrorInfoList)
Count = len(CategoryList)

Dll.QueryDatas(ClientId, CategoryArray, c_int(Count), ResultArray, ErrorInfoArray)

# 多账号批量查询各类交易数据
# 这里的Category参数只能对应账户,不能一个账户查询多种类型数据
ClientIdList = [ClientId1, ClientId2]
ClientIdArray = (c_int * len(ClientIdList))(*ClientIdList)
CategoryList = [Category1, Category2]
CategoryArray = (c_int * len(CategoryList))(*CategoryList)
ResultList = [cast(Result1, c_char_p), cast(Result2, c_char_p)]
ResultArray = (c_char_p * len(ResultList))(*ResultList)
ErrorInfoList = [cast(ErrorInfo1, c_char_p), cast(ErrorInfo2, c_char_p)]
ErrorInfoArray = (c_char_p * len(ErrorInfoList))(*ErrorInfoList)
Count = len(CategoryList)

Dll.QueryMultiAccountsDatas(ClientIdArray, CategoryArray, c_int(Count), ResultArray, ErrorInfoArray)

# 查询各类历史数据
Dll.QueryHistoryData(ClientId, c_int(Category), c_char_p(StartDate), c_char_p(EndDate), Result, ErrorInfo)

# 委托下单
Dll.SendOrder(ClientId, c_int(Category), c_int(EntrustType), c_char_p(Gddm), c_char_p(Zqdm), c_float(Price),
              c_int(Quantity), Result, ErrorInfo)

# 单账户批量下单
# Count参数和单账户批量查询意义一样
CategoryList = [Category1, Category2]
CategoryArray = (c_int * len(CategoryList))(*CategoryList)
EntrustTypeList = [EntrustType1, EntrustType2]
EntrustTypeArray = (c_int * len(EntrustTypeList))(*EntrustTypeList)
GddmList = [b'Gddm1', b'Gddm2']
GddmArray = (c_char_p * len(GddmList))(*GddmList)
ZqdmList = [b'Zqdm1', b'Zqdm2']
ZqdmArray = (c_char_p * len(ZqdmList))(*ZqdmList)
PriceList = [Price1, Price2]
PriceArray = (c_float * len(PriceList))(*PriceList)
QuantityList = [Quantity1, Quantity2]
QuantityArray = (c_int * len(QuantityList))(*QuantityList)
Count = len(CategoryList)
ResultList = [cast(Result1, c_char_p), cast(Result2, c_char_p)]
ResultArray = (c_char_p * len(ResultList))(*ResultList)
ErrorInfoList = [cast(ErrorInfo1, c_char_p), cast(ErrorInfo2, c_char_p)]
ErrorInfoArray = (c_char_p * len(ErrorInfoList))(*ErrorInfoList)

Dll.SendOrders(ClientId, CategoryArray, EntrustTypeArray, GddmArray, ZqdmArray, PriceArray, QuantityArray, c_int(Count),
               ResultArray, ErrorInfoArray)

# 多账户批量下单
# Count参数和多账号批量查询意义一样
ClientIdList = [ClientId1, ClientId2]
ClientIdArray = (c_int * len(ClientIdList))(*ClientIdList)
CategoryList = [Category1, Category2]
CategoryArray = (c_int * len(CategoryList))(*CategoryList)
EntrustTypeList = [EntrustType1, EntrustType2]
EntrustTypeArray = (c_int * len(EntrustTypeList))(*EntrustTypeList)
GddmList = [b'Gddm1', b'Gddm2']
GddmArray = (c_char_p * len(GddmList))(*GddmList)
ZqdmList = [b'Zqdm1', b'Zqdm2']
ZqdmArray = (c_char_p * len(ZqdmList))(*ZqdmList)
PriceList = [Price1, Price2]
PriceArray = (c_float * len(PriceList))(*PriceList)
QuantityList = [Quantity1, Quantity2]
QuantityArray = (c_int * len(QuantityList))(*QuantityList)
Count = len(CategoryList)
ResultList = [cast(Result1, c_char_p), cast(Result2, c_char_p)]
ResultArray = (c_char_p * len(ResultList))(*ResultList)
ErrorInfoList = [cast(ErrorInfo1, c_char_p), cast(ErrorInfo2, c_char_p)]
ErrorInfoArray = (c_char_p * len(ErrorInfoList))(*ErrorInfoList)

Dll.SendMultiAccountsOrders(ClientIdArray, CategoryArray, EntrustTypeArray, GddmArray, ZqdmArray, PriceArray,
                            QuantityArray, c_int(Count), ResultArray, ErrorInfoArray)

# 委托撤单
Dll.CancelOrder(ClientId, c_char_p(ExchangeId), c_char_p(EntrustId), Result, ErrorInfo)

# 单账户批量撤单
ExchangeIdList = [b'ExchangeId1', b'ExchangeId2']
ExchangeIdArray = (c_char_p * len(ExchangeIdList))(*ExchangeIdList)
EntrustIdList = [b'EntrustId1', b'EntrustId2']
EntrustIdArray = (c_char_p * len(EntrustIdList))(*EntrustIdList)
Count = len(ExchangeIdList)
ResultList = [cast(Result1, c_char_p), cast(Result2, c_char_p)]
ResultArray = (c_char_p * len(ResultList))(*ResultList)
ErrorInfoList = [cast(ErrorInfo1, c_char_p), cast(ErrorInfo2, c_char_p)]
ErrorInfoArray = (c_char_p * len(ErrorInfoList))(*ErrorInfoList)

Dll.CancelOrders(ClientId, ExchangeIdArray, EntrustIdArray, c_int(Count), ResultArray, ErrorInfoArray)

# 多账户批量撤单
ClientIdList = [ClientId1, ClientId2]
ClientIdArray = (c_int * len(ClientIdList))(*ClientIdList)
ExchangeIdList = [b'ExchangeId1', b'ExchangeId2']
ExchangeIdArray = (c_char_p * len(ExchangeIdList))(*ExchangeIdList)
EntrustIdList = [b'EntrustId1', b'EntrustId2']
EntrustIdArray = (c_char_p * len(EntrustIdList))(*EntrustIdList)
Count = len(ExchangeIdList)
ResultList = [cast(Result1, c_char_p), cast(Result2, c_char_p)]
ResultArray = (c_char_p * len(ResultList))(*ResultList)
ErrorInfoList = [cast(ErrorInfo1, c_char_p), cast(ErrorInfo2, c_char_p)]
ErrorInfoArray = (c_char_p * len(ErrorInfoList))(*ErrorInfoList)

Dll.CancelMultiAccountsOrders(ClientIdArray, ExchangeIdArray, EntrustIdArray, c_int(Count), ResultArray, ErrorInfoArray)

# 获取五档报价
Dll.GetQuote(ClientId, c_char_p(Zqdm), Result, ErrorInfo)

# 单账户批量获取五档报价
ZqdmList = [b'Zqdm1', b'Zqdm2']
ZqdmArray = (c_char_p * len(ZqdmList))(*ZqdmList)
Count = len(ZqdmList)
ResultList = [cast(Result1, c_char_p), cast(Result2, c_char_p)]
ResultArray = (c_char_p * len(ResultList))(*ResultList)
ErrorInfoList = [cast(ErrorInfo1, c_char_p), cast(ErrorInfo2, c_char_p)]
ErrorInfoArray = (c_char_p * len(ErrorInfoList))(*ErrorInfoList)

Dll.GetQuote(ClientId, ZqdmArray, c_int(Count), ResultArray, ErrorInfoArray)

# 多账户批量获取五档报价
ClientIdList = [ClientId1, ClientId2]
ClientIdArray = (c_int * len(ClientIdList))(*ClientIdList)
ZqdmList = [b'Zqdm1', b'Zqdm2']
ZqdmArray = (c_char_p * len(ZqdmList))(*ZqdmList)
Count = len(ZqdmList)
ResultList = [cast(Result1, c_char_p), cast(Result2, c_char_p)]
ResultArray = (c_char_p * len(ResultList))(*ResultList)
ErrorInfoList = [cast(ErrorInfo1, c_char_p), cast(ErrorInfo2, c_char_p)]
ErrorInfoArray = (c_char_p * len(ErrorInfoList))(*ErrorInfoList)

Dll.GetMultiAccountsQuotes(ClientIdArray, ZqdmArray, c_int(Count), ResultArray, ErrorInfoArray)

# 融资融券账户直接还款
Dll.Repay(ClientId, c_char_p(Amount), Result, ErrorInfo)

# 查询API授权到期日期
Dll.GetExpireDate(ClientId)

# 登出交易账号
# 连接多少交易账号,相应登出多少个
Dll.Logoff(ClientId)

# DLL反初始化
# Init失败后无需调用反初始化
Dll.Deinit()
