/*
 * 使用流程:
 * 1. visual studio 新建 C++ 空项目
 * 2. 将本CPP文件添加到项目中
 * 3. 解决方案平台选择 x86 (暂不支持x64)
 * 4. 编译完成后, 将 MetaTrade.dll/MetaTrade.lic/MetaTrade.ini 放到 exe 输出目录中
 * 5. 运行调试
*/
#include <Windows.h>

#include <iostream>
#include <stdexcept>
#include <string>

// API初始化, 返回授权成功的交易账号数量
// 返回值 < 1 时, 无需调用Deinit接口, 也不能调用其它接口, 否则会出错!
typedef int (*InitFn)();
// API反初始化
typedef void (*DeinitFn)();
// 登录交易账户
typedef int (*LogonFn)(const char* Ip, short Port, const char* Version,
                       short Yybid, const char* Account,
                       const char* TradeAccount, const char* JyPassword,
                       const char* TxPassword, char* ErrorInfo);
// 登出交易账户
typedef void (*LogoffFn)(int ClientId);
// 查询各类交易数据
typedef void (*QueryDataFn)(int ClientId, int Category, char* Result,
                            char* ErrorInfo);
// 单账户批量查询各类交易数据
typedef void (*QueryDatasFn)(int ClientId, int Category[], int Count,
                             char* Result[], char* ErrorInfo[]);
// 多账户批量查询各类交易数据
typedef void (*QueryMultiAccountsDatasFn)(int ClientId[], int Category[],
                                          int Count, char* Result[],
                                          char* ErrorInfo[]);
// 查询各类历史数据
typedef void (*QueryHistoryDataFn)(int ClientId, int Category,
                                   const char* StartDate, const char* EndDate,
                                   char* Result, char* ErrorInfo);
// 委托下单
typedef void (*SendOrderFn)(int ClientId, int Category, int EntrustType,
                            const char* Gddm, const char* Zqdm, float Price,
                            int Quantity, char* Result, char* ErrorInfo);
// 单账户批量下单
typedef void (*SendOrdersFn)(int ClientId, int Category[], int EntrustType[],
                             const char* Gddm[], const char* Zqdm[],
                             float Price[], int Quantity[], int Count,
                             char* Result[], char* ErrorInfo[]);
// 多账户批量下单
typedef void (*SendMultiAccountsOrdersFn)(int ClientId[], int Category[],
                                          int EntrustType[], const char* Gddm[],
                                          const char* Zqdm[], float Price[],
                                          int Quantity[], int Count,
                                          char* Result[], char* ErrorInfo[]);
// 委托撤单
typedef void (*CancelOrderFn)(int ClientId, const char* ExchangeId,
                              const char* EntrustId, char* Result,
                              char* ErrorInfo);
// 单账户批量撤单
typedef void (*CancelOrdersFn)(int ClientId, const char* ExchangeId[],
                               const char* EntrustId[], int Count,
                               char* Result[], char* ErrorInfo[]);
// 多账户批量撤单
typedef void (*CancelMultiAccountsOrdersFn)(int ClientId[],
                                            const char* ExchangeId[],
                                            const char* EntrustId[], int Count,
                                            char* Result[], char* ErrorInfo[]);
// 获取五档报价
typedef void (*GetQuoteFn)(int ClientId, const char* Zqdm, char* Result,
                           char* ErrorInfo);
// 单账户批量获取五档报价
typedef void (*GetQuotesFn)(int ClientId, const char* Zqdm[], int Count,
                            char* Result[], char* ErrorInfo[]);
// 多账户批量获取五档报价
typedef void (*GetMultiAccountsQuotesFn)(int ClientId[], const char* Zqdm[],
                                         int Count, char* Result[],
                                         char* ErrorInfo[]);
// 融资融券账户直接还款
typedef void (*RepayFn)(int ClientId, const char* Amount, char* Result,
                        char* ErrorInfo);
// 查询API授权到期日期
typedef int (*GetExpireDateFn)(int ClientId);

class MetaTradeAPI {
 public:
  MetaTradeAPI() {
    // 加载DLL
    m_hInstance = LoadLibraryA("MetaTrade.dll");
    if (!m_hInstance) {
      throw std::runtime_error("加载DLL失败!");
    }
    // 加载函数
    m_initFn = reinterpret_cast<InitFn>(GetProcAddress(m_hInstance, "Init"));
    m_deinitFn =
        reinterpret_cast<DeinitFn>(GetProcAddress(m_hInstance, "Deinit"));
    m_logonFn = reinterpret_cast<LogonFn>(GetProcAddress(m_hInstance, "Logon"));
    m_logoffFn =
        reinterpret_cast<LogoffFn>(GetProcAddress(m_hInstance, "Logoff"));
    m_queryDataFn =
        reinterpret_cast<QueryDataFn>(GetProcAddress(m_hInstance, "QueryData"));
    m_queryDatasFn = reinterpret_cast<QueryDatasFn>(
        GetProcAddress(m_hInstance, "QueryDatas"));
    m_queryMultiAccountsDatasFn = reinterpret_cast<QueryMultiAccountsDatasFn>(
        GetProcAddress(m_hInstance, "QueryMultiAccountsDatas"));
    m_queryHistoryDataFn = reinterpret_cast<QueryHistoryDataFn>(
        GetProcAddress(m_hInstance, "QueryHistoryData"));
    m_sendOrderFn =
        reinterpret_cast<SendOrderFn>(GetProcAddress(m_hInstance, "SendOrder"));
    m_sendOrdersFn = reinterpret_cast<SendOrdersFn>(
        GetProcAddress(m_hInstance, "SendOrders"));
    m_sendMultiAccountsOrdersFn = reinterpret_cast<SendMultiAccountsOrdersFn>(
        GetProcAddress(m_hInstance, "SendMultiAccountsOrders"));
    m_cancelOrderFn = reinterpret_cast<CancelOrderFn>(
        GetProcAddress(m_hInstance, "CancelOrder"));
    m_cancelOrdersFn = reinterpret_cast<CancelOrdersFn>(
        GetProcAddress(m_hInstance, "CancelOrders"));
    m_cancelMultiAccountsOrdersFn =
        reinterpret_cast<CancelMultiAccountsOrdersFn>(
            GetProcAddress(m_hInstance, "CancelMultiAccountsOrders"));
    m_getQuoteFn =
        reinterpret_cast<GetQuoteFn>(GetProcAddress(m_hInstance, "GetQuote"));
    m_getQuotesFn =
        reinterpret_cast<GetQuotesFn>(GetProcAddress(m_hInstance, "GetQuotes"));
    m_getMultiAccountsQuotesFn = reinterpret_cast<GetMultiAccountsQuotesFn>(
        GetProcAddress(m_hInstance, "GetMultiAccountsQuotes"));
    m_repayFn = reinterpret_cast<RepayFn>(GetProcAddress(m_hInstance, "Repay"));
    m_getExpireDateFn = reinterpret_cast<GetExpireDateFn>(
        GetProcAddress(m_hInstance, "GetExpireDate"));
    // API初始化
    if (m_initFn() < 1) {
      // 无需 Deinit()
      FreeLibrary(m_hInstance);
      throw std::runtime_error("没有可用的授权账户!");
    }
  }

  ~MetaTradeAPI() {
    // API清理
    m_deinitFn();
    // 卸载DLL
    FreeLibrary(m_hInstance);
  }

  // 登录交易账户
  int Logon(const char* Ip, short Port, const char* Version, short Yybid,
            const char* Account, const char* TradeAccount,
            const char* JyPassword, const char* TxPassword,
            char* ErrorInfo) const {
    return m_logonFn(Ip, Port, Version, Yybid, Account, TradeAccount,
                     JyPassword, TxPassword, ErrorInfo);
  }
  // 登出交易账户
  void Logoff(int ClientId) const { m_logoffFn(ClientId); }
  // 查询各类交易数据
  void QueryData(int ClientId, int Category, char* Result,
                 char* ErrorInfo) const {
    m_queryDataFn(ClientId, Category, Result, ErrorInfo);
  }
  // 单账户批量查询各类交易数据
  void QueryDatas(int ClientId, int Category[], int Count, char* Result[],
                  char* ErrorInfo[]) const {
    m_queryDatasFn(ClientId, Category, Count, Result, ErrorInfo);
  }
  // 多账户批量查询各类交易数据
  void QueryMultiAccountsDatas(int ClientId[], int Category[], int Count,
                               char* Result[], char* ErrorInfo[]) const {
    m_queryMultiAccountsDatasFn(ClientId, Category, Count, Result, ErrorInfo);
  }
  // 查询各类历史数据
  void QueryHistoryData(int ClientId, int Category, const char* StartDate,
                        const char* EndDate, char* Result,
                        char* ErrorInfo) const {
    m_queryHistoryDataFn(ClientId, Category, StartDate, EndDate, Result,
                         ErrorInfo);
  }
  // 委托下单
  void SendOrder(int ClientId, int Category, int EntrustType, const char* Gddm,
                 const char* Zqdm, float Price, int Quantity, char* Result,
                 char* ErrorInfo) const {
    m_sendOrderFn(ClientId, Category, EntrustType, Gddm, Zqdm, Price, Quantity,
                  Result, ErrorInfo);
  }
  // 单账户批量下单
  void SendOrders(int ClientId, int Category[], int EntrustType[],
                  const char* Gddm[], const char* Zqdm[], float Price[],
                  int Quantity[], int Count, char* Result[],
                  char* ErrorInfo[]) const {
    m_sendOrdersFn(ClientId, Category, EntrustType, Gddm, Zqdm, Price, Quantity,
                   Count, Result, ErrorInfo);
  }
  // 多账户批量下单
  void SendMultiAccountsOrders(int ClientId[], int Category[],
                               int EntrustType[], const char* Gddm[],
                               const char* Zqdm[], float Price[],
                               int Quantity[], int Count, char* Result[],
                               char* ErrorInfo[]) const {
    m_sendMultiAccountsOrdersFn(ClientId, Category, EntrustType, Gddm, Zqdm,
                                Price, Quantity, Count, Result, ErrorInfo);
  }
  // 委托撤单
  void CancelOrder(int ClientId, const char* ExchangeId, const char* EntrustId,
                   char* Result, char* ErrorInfo) const {
    m_cancelOrderFn(ClientId, ExchangeId, EntrustId, Result, ErrorInfo);
  }
  // 单账户批量撤单
  void CancelOrders(int ClientId, const char* ExchangeId[],
                    const char* EntrustId[], int Count, char* Result[],
                    char* ErrorInfo[]) const {
    m_cancelOrdersFn(ClientId, ExchangeId, EntrustId, Count, Result, ErrorInfo);
  }
  // 多账户批量撤单
  void CancelMultiAccountsOrders(int ClientId[], const char* ExchangeId[],
                                 const char* EntrustId[], int Count,
                                 char* Result[], char* ErrorInfo[]) const {
    m_cancelMultiAccountsOrdersFn(ClientId, ExchangeId, EntrustId, Count,
                                  Result, ErrorInfo);
  }
  // 获取五档报价
  void GetQuote(int ClientId, const char* Zqdm, char* Result,
                char* ErrorInfo) const {
    m_getQuoteFn(ClientId, Zqdm, Result, ErrorInfo);
  }
  // 单账户批量获取五档报价
  void GetQuotes(int ClientId, const char* Zqdm[], int Count, char* Result[],
                 char* ErrorInfo[]) const {
    m_getQuotesFn(ClientId, Zqdm, Count, Result, ErrorInfo);
  }
  // 多账户批量获取五档报价
  void GetMultiAccountsQuotes(int ClientId[], const char* Zqdm[], int Count,
                              char* Result[], char* ErrorInfo[]) const {
    m_getMultiAccountsQuotesFn(ClientId, Zqdm, Count, Result, ErrorInfo);
  }
  // 融资融券账户直接还款
  void Repay(int ClientId, const char* Amount, char* Result,
             char* ErrorInfo) const {
    m_repayFn(ClientId, Amount, Result, ErrorInfo);
  }
  // 查询API授权到期日期
  int GetExpireDate(int ClientId) const { return m_getExpireDateFn(ClientId); }

 private:
  HINSTANCE m_hInstance;
  InitFn m_initFn;
  DeinitFn m_deinitFn;
  LogonFn m_logonFn;
  LogoffFn m_logoffFn;
  QueryDataFn m_queryDataFn;
  QueryDatasFn m_queryDatasFn;
  QueryMultiAccountsDatasFn m_queryMultiAccountsDatasFn;
  QueryHistoryDataFn m_queryHistoryDataFn;
  SendOrderFn m_sendOrderFn;
  SendOrdersFn m_sendOrdersFn;
  SendMultiAccountsOrdersFn m_sendMultiAccountsOrdersFn;
  CancelOrderFn m_cancelOrderFn;
  CancelOrdersFn m_cancelOrdersFn;
  CancelMultiAccountsOrdersFn m_cancelMultiAccountsOrdersFn;
  GetQuoteFn m_getQuoteFn;
  GetQuotesFn m_getQuotesFn;
  GetMultiAccountsQuotesFn m_getMultiAccountsQuotesFn;
  RepayFn m_repayFn;
  GetExpireDateFn m_getExpireDateFn;
};

// 获取 std::string 缓冲区
#define BUF(stdstr) const_cast<char*>(stdstr.c_str())

int main() {
  using namespace std;

  MetaTradeAPI api;  // 唯一一个API对象

  string Result(1024 * 1024, NULL);  // 预留1024*1024字节空间
  string ErrorInfo(256, NULL);       // 预留256字节空间

  cout << "===== 登录第一个账号 =====" << endl;
  // 登录账号，如果是两融账号需要额外添加 .C 结尾
  // 交易账号, 一般与登录账号相同, 但无需额外添加 .C 结尾
  // 通讯密码, 一般为空字符串或与交易密码相同
  const int ClientId1 = api.Logon("交易服务器IP", 7788, "", 0, "123456", "123456",
                                  "password", "", BUF(ErrorInfo));
  if (ClientId1 <= 0) {
    cout << "第一个账号登录失败!" << endl;
    cout << ErrorInfo.c_str() << endl;
    return -1;
  } else {
    cout << "第一个账号登录成功!" << endl;
  }

  // 查询资金, 其它信息可通过修改Category参数查询
  api.QueryData(ClientId1, 0 /*Category*/, BUF(Result), BUF(ErrorInfo));
  if (NULL != ErrorInfo[0]) {  // 注意这里不能通过ErrorInfo.empty()判断
    cout << "===== 查询资金出错 =====" << endl;
    cout << ErrorInfo.c_str() << endl;
  } else {
    cout << "===== 查询资金成功 =====" << endl;
    // Result格式为表格数据, 每一行通过换行符\n分割，每一列通过制表符\t分割
    cout << Result.c_str() << endl;
  }

  // 获取五档报价
  api.GetQuote(ClientId1, "600006", BUF(Result), BUF(ErrorInfo));
  if (NULL != ErrorInfo[0]) {
    cout << "===== 获取五档报价出错 =====" << endl;
    cout << ErrorInfo.c_str() << endl;
  } else {
    cout << "===== 获取五档报价成功 =====" << endl;
    cout << Result.c_str() << endl;
  }

  // 单笔委托下单
  // 上海股东代码, 下单深圳股票时需使用深圳股东代码
  api.SendOrder(ClientId1, 0 /*卖入*/, 0 /*限价委托*/, "A123456", "600006",
                1.23f, 100, BUF(Result), BUF(ErrorInfo));
  if (NULL != ErrorInfo[0]) {  // 注意这里不能通过ErrorInfo.empty()判断
    cout << "===== 单笔委托下单出错 =====" << endl;
    cout << ErrorInfo.c_str() << endl;
  } else {
    cout << "===== 单笔委托下单成功 =====" << endl;
    // Result包含委托编号, 可用于撤单
    cout << Result.c_str() << endl;
  }

  cout << "===== 单账户批量下单 =====" << endl;
  const int Count = 20;
  int Categorys[Count] = {0};
  int EntrustTypes[Count] = {0};
  const char* Gddms[Count] = {0};
  const char* Zqdms[Count] = {0};
  float Prices[Count] = {0};
  int Quantities[Count] = {0};
  char* Results[Count] = {0};
  char* ErrorInfos[Count] = {0};
  for (int i = 0; i < Count; i++) {
    Categorys[i] = 0;         // 第i项委托的委托类别
    EntrustTypes[i] = 0;      // 第i项委托的报价方式
    Gddms[i] = "A123456";  // 第i项委托的股东代码
    Zqdms[i] = "600006";      // 第i项委托的证券代码
    Prices[i] = 1.23f;        // 第i项委托的委托价格
    Quantities[i] = 100;      // 第i项委托的委托数量
    Results[i] = (char*)malloc(1024 * 1024);
    ErrorInfos[i] = (char*)malloc(256);
  }
  api.SendOrders(ClientId1, Categorys, EntrustTypes, Gddms, Zqdms, Prices,
                 Quantities, Count, Results, ErrorInfos);
  for (int i = 0; i < Count; i++) {
    if (NULL != ErrorInfos[i][0]) {
      cout << "第" << i + 1 << "项委托失败:" << endl;
      cout << ErrorInfos[i] << endl;
    } else {
      cout << "第" << i + 1 << "项委托结果:" << endl;
      cout << Results[i] << endl;
    }
    free(Results[i]);
    Results[i] = NULL;
    free(ErrorInfos[i]);
    ErrorInfos[i] = NULL;
  }

  cout << "===== 登录第二个账号 =====" << endl;
  const int ClientId2 = api.Logon("2.4.6.8", 8899, "", 0, "987654", "987654",
                                  "password", "", BUF(ErrorInfo));
  if (ClientId2 <= 0) {
    cout << "第二个账号登录失败!" << endl;
    cout << ErrorInfo.c_str() << endl;
    api.Logoff(ClientId1);
    return -1;
  } else {
    cout << "第二个账号登录成功!" << endl;
  }

  cout << "===== 多账户批量下单 =====" << endl;
  int ClientIds[Count] = {0};
  for (int i = 0; i < Count; i++) {
    // 第i项委托的下单账户Id
    ClientIds[i] = (i % 2) ? ClientId2 : ClientId1;
    Categorys[i] = 0;     // 第i项委托的委托类别
    EntrustTypes[i] = 0;  // 第i项委托的报价方式
    // 第i项委托的股东代码
    Gddms[i] = (i % 2) ? "A987654" : "A123456";
    Zqdms[i] = "600006";  // 第i项委托的证券代码
    Prices[i] = 1.23f;    // 第i项委托的委托价格
    Quantities[i] = 100;  // 第i项委托的委托数量
    Results[i] = (char*)malloc(1024 * 1024);
    ErrorInfos[i] = (char*)malloc(256);
  }
  api.SendMultiAccountsOrders(ClientIds, Categorys, EntrustTypes, Gddms, Zqdms,
                              Prices, Quantities, Count, Results, ErrorInfos);
  for (int i = 0; i < Count; i++) {
    if (NULL != ErrorInfos[i][0]) {
      cout << "第" << i + 1 << "项委托失败:" << endl;
      cout << ErrorInfos[i] << endl;
    } else {
      cout << "第" << i + 1 << "项委托结果:" << endl;
      cout << Results[i] << endl;
    }
    free(Results[i]);
    Results[i] = NULL;
    free(ErrorInfos[i]);
    ErrorInfos[i] = NULL;
  }

  // 登出账户
  api.Logoff(ClientId1);
  api.Logoff(ClientId2);

  return 0;
}
