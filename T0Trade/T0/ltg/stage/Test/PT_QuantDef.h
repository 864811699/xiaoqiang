#ifndef _QuantDef_H_
#define _QuantDef_H_
#include <stdint.h>
#include <string.h>

namespace QuantPlus
{
	/// <summary>
	/// QuantApi
	/// </summary>
	typedef char PT_CodeType[64];

	// ���黷��
	enum PT_QuantMdAppEType
	{
		PT_QuantMdAppEType_Test = 0,            // ���Ի�����ע���û�����ֹͣʹ�ã�������
		PT_QuantMdAppEType_Real,                // ��������
	};

	// ���׻���
	enum PT_QuantTdAppEType
	{
		PT_QuantTdAppEType_Real,                // ��������
		PT_QuantTdAppEType_Test,                // ���Ի��� (ע���û�����ֹͣʹ�ã�����)    
		PT_QuantTdAppEType_Simulation           // ģ�⻷��
	};

	// ҵ�����������
	enum PT_Quant_APPServerType
	{
		PT_Quant_APPServerType_RealMdServer = 0,            // ʵʱ���������
		PT_Quant_APPServerType_HisrotyMdServer,             // ��ʷ���������
		PT_Quant_APPServerType_CacheMdServer,               // ʵʱ���������

		PT_Quant_APPServerType_TDServer = 10,               //  ���׷�����
	};

	// �û���ɫ
	enum PT_QuantUserType
	{
		PT_QuantUserType_Risk = 1,                       //  ���Ա
		PT_QuantUserType_Trade                           //  ����Ա
	};

	// �û���Ϣ
	class PT_QuantUserBase
	{
	public:
		int64_t                        nId;
		char               szUserName[128];                 // �û���
		char               szNickName[128];                 // �û�����
		int                       nGroupId;                 // ��ID

		PT_QuantUserType   nUserRole;                       // �û���ɫ

		double          nStampTax;                   //  ӡ��˰
		double          nTransferFees;               //  ������
		double          nCommissions;                //  Ӷ��

		char			szSecurityCode[32];					//��֤��

		PT_QuantUserBase()
		{
			nId = 0;
			nUserRole = PT_QuantUserType_Trade;
			memset(szUserName, 0, 128);
			memset(szNickName, 0, 128);
			nGroupId = 0;

			nStampTax = 0;
			nTransferFees = 0;
			nCommissions = 0;

			memset(szSecurityCode, 0, 32);
		}
		virtual ~PT_QuantUserBase()
		{

		}
	};

	// ȯ����
	class PT_QuantUserCodeControl
	{
	public:
		char               szWinCode[64];
		double             nCaptial;                   // �������ʽ�
		int                nLendingAmount;             // ������ȯ
		PT_QuantUserCodeControl()
		{
			memset(szWinCode, 0, 64);
			nCaptial = 0;
			nLendingAmount = 0;
		}
	};

	class PT_QuantUser : public PT_QuantUserBase
	{
	public:
		bool               ifStopTrade;                     // �Ƿ�ͣ��
		int               nStopTradePostion;                // ͣ��λ(�����ʽ���)
		double     nStopPercentTradePostion;                // ͣ��λ(�������)

		int                nSinglePositionHoldTime;         // ���ʳֲ�ʱ����ֵ
		int                nSinglePositionLoss;             // ���ʳֲֿ�����ֵ(�����ʽ���)
		double      nSinglePercentPositionLoss;             // ���ʳֲֿ�����ֵ(�������)

		int                      nCodeControlNum;       //  ����ȯ����
		PT_QuantUserCodeControl* pCodeControl;          // ������ȯ��Ϣ��ָ��ƫ��

		int                      nDisableCodeNum;
		PT_CodeType*             pDisableCode;

		/////add by long.wang 
		int64_t  	basePoint;							 	 // bpֵ
		int32_t		forceClosePostion;						 // ǿƽλ
		double 		nPercentForceClosePostion ; 			 // ǿƽλ����		
		/////
		
		PT_QuantUser()
		{
			ifStopTrade = false;
			nStopTradePostion = 1000000;
			nStopPercentTradePostion = 0.1;

			nSinglePositionHoldTime = 1000000;
			nSinglePositionLoss = 1000000;
			nSinglePercentPositionLoss = 0.1;

			nCodeControlNum = 0;
			pCodeControl = NULL;

			nDisableCodeNum = 0;
			pDisableCode = NULL;

			///// add by long.wang
			basePoint = 0;
			forceClosePostion = 0;
			nPercentForceClosePostion = 0.0;
			/////
		}
		virtual ~PT_QuantUser()
		{
			if(pCodeControl)
			{
				delete[] pCodeControl;
				pCodeControl = NULL;
			}

			if(pDisableCode)
			{
				delete[] pDisableCode;
				pDisableCode = NULL;
			}
		}
	};

	struct PT_BackTestReq
	{
		int64_t nSimAccountId;         // ģ���ʽ��˺�


		PT_BackTestReq()
		{
			nSimAccountId = 0;
		}
		virtual ~PT_BackTestReq()
		{

		}
	};

	/// <summary>
	/// ����ṹ��
	/// </summary>
#pragma pack(push)
#pragma pack(1)
	typedef int  MD_ReqID;
	typedef char MD_CodeType[32];
	typedef char MD_CodeName[32];
	typedef char MD_ISODateTimeType[21];    //���ں�ʱ������(��ʽ yyyy-MM-dd hh:mm:ss)
	typedef char MD_ShortText[128];
	typedef char MD_Text[1024];

	// ����������
	typedef int MD_SrvType;
#define MD_SrvType_none            0x0000   // δ֪����
#define MD_SrvType_history         0x0001   // ��ʷ�������������
#define MD_SrvType_cache           0x0002   // ʵʱ�������������
#define MD_SrvType_realtime        0x0004   // ʵʱ�������������

	// ��������
	typedef int MD_CycType;
#define MD_CycType_none            0x0000   // δ֪����
#define MD_CycType_second_10       0x0001   // 10��
#define MD_CycType_minute          0x0002   // ��
#define MD_CycType_minute_5        0x0004   // 5��
#define MD_CycType_minute_15       0x0008   // 15��
#define MD_CycType_minute_30       0x0010   // 30��
#define MD_CycType_hour            0x0020   // Сʱ
#define MD_CycType_day             0x0040   // ��

	// ��������
	typedef int MD_SubType;
#define MD_SubType_none            0x0000   // δ֪����
#define MD_SubType_market          0x0001   // ��������
#define MD_SubType_index           0x0002   // ָ������
#define MD_SubType_trans           0x0004   // ��ʳɽ�
#define MD_SubType_order           0x0008   // ���ί��
#define MD_SubType_order_queue     0x0010   // ί�ж���
#define MD_SubType_future          0x0020   // �ڻ�����
#define MD_SubType_future_option   0x0040   // ��Ȩ����
#define MD_SubType_kline           0x0080   // K������

	struct MD_DATA_CODE
	{
		char    szWindCode[32];         //Wind Code: AG1302.SHF
		char    szMarket[8];            //market code: SHF
		char    szCode[32];             //original code:ag1302
		char    szENName[32];
		char    szCNName[32];           //chinese name: ����1302
		int     nType;
	};

	struct MD_DATA_OPTION_CODE
	{
		MD_DATA_CODE basicCode;

		char szContractID[32];          // ��Ȩ��Լ����
		char szUnderlyingSecurityID[32];// ���֤ȯ����
		char chCallOrPut;               // �Ϲ��Ϲ�C1        �Ϲ������ֶ�Ϊ��C������Ϊ�Ϲ������ֶ�Ϊ��P��
		int  nExerciseDate;             // ��Ȩ��Ȩ�գ�YYYYMMDD

		//�����ֶ�
		char chUnderlyingType;          // ���֤ȯ����C3    0-A�� 1-ETF (EBS �C ETF�� ASH �C A ��)
		char chOptionType;              // ŷʽ��ʽC1        ��Ϊŷʽ��Ȩ�����ֶ�Ϊ��E������Ϊ��ʽ��Ȩ�����ֶ�Ϊ��A��

		char chPriceLimitType;          // �ǵ�����������C1 ��N����ʾ���ǵ�����������, ��R����ʾ���ǵ�����������
		int  nContractMultiplierUnit;   // ��Լ��λ,         ������Ȩ��Ϣ������ĺ�Լ��λ, һ��������
		int  nExercisePrice;            // ��Ȩ��Ȩ��,       ������Ȩ��Ϣ���������Ȩ��Ȩ�ۣ��Ҷ��룬��ȷ����
		int  nStartDate;                // ��Ȩ�׸�������,YYYYMMDD
		int  nEndDate;                  // ��Ȩ�������/��Ȩ�գ�YYYYMMDD
		int  nExpireDate;               // ��Ȩ�����գ�YYYYMMDD
	};

	// ��ע: ���ṹ����ֶ��еĵ�λΪ�ο���λ, ��������, ���������ϳ���������������̿ڵ�λΪ׼
	struct MD_DATA_MARKET
	{
		char            szWindCode[32];         //600001.SH
		char            szCode[32];             //ԭʼCode
		int             nActionDay;             //����(YYYYMMDD)
		int             nTime;                  //ʱ��(HHMMSSmmm)
		int             nStatus;                //״̬
		unsigned int    nPreClose;              //ǰ���̼�=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		unsigned int    nOpen;                  //���̼�=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		unsigned int    nHigh;                  //��߼�=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		unsigned int    nLow;                   //��ͼ�=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		unsigned int    nMatch;                 //���¼�=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		unsigned int    nAskPrice[10];          //������=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		unsigned int    nAskVol[10];            //������=ʵ�ʹ���(��λ: ��)
		unsigned int    nBidPrice[10];          //�����=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		unsigned int    nBidVol[10];            //������=ʵ�ʹ���(��λ: ��)
		unsigned int    nNumTrades;             //�ɽ�����=ʵ�ʱ���(��λ: ��)
		int64_t         iVolume;                //�ɽ�����=ʵ�ʹ���(��λ: ��)
		int64_t         iTurnover;              //�ɽ��ܽ��=ʵ�ʽ��(��λ: Ԫ)
		int64_t         nTotalBidVol;           //ί����������=ʵ�ʹ���(��λ: ��)
		int64_t         nTotalAskVol;           //ί����������=ʵ�ʹ���(��λ: ��)
		unsigned int    nWeightedAvgBidPrice;   //��Ȩƽ��ί��۸�=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		unsigned int    nWeightedAvgAskPrice;   //��Ȩƽ��ί���۸�=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		int             nIOPV;                  //IOPV��ֵ��ֵ
		int             nYieldToMaturity;       //����������
		unsigned int    nHighLimited;           //��ͣ��=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		unsigned int    nLowLimited;            //��ͣ��=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
	};

	// ��ע: ���ṹ����ֶ��еĵ�λΪ�ο���λ, ��������, ���������ϳ���������������̿ڵ�λΪ׼
	struct MD_DATA_INDEX
	{
		char        szWindCode[32];         //600001.SH
		char        szCode[32];             //ԭʼCode
		int         nActionDay;             //����(YYYYMMDD)
		int         nTime;                  //ʱ��(HHMMSSmmm)
		int         nOpenIndex;             //����ָ��=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		int         nHighIndex;             //���ָ��=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		int         nLowIndex;              //���ָ��=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		int         nLastIndex;             //����ָ��=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		int64_t     iTotalVolume;           //���������Ӧָ���Ľ�������=ʵ�ʹ���(��λ: ��)
		int64_t     iTurnover;              //���������Ӧָ���ĳɽ����=ʵ�ʽ��(��λ: Ԫ)
		int         nPreCloseIndex;         //ǰ��ָ��=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
	};

	// ��ע: ���ṹ��Ϊͨ���ڻ����ݽṹ��, �ڻ�Ʒ�ֲ�ͬ���׵�λ��ͬ, ��������, ���������ϳ���������������̿ڵ�λΪ׼
	struct MD_DATA_FUTURE
	{
		char            szWindCode[32];         //600001.SH
		char            szCode[32];             //ԭʼCode
		int             nActionDay;             //����(YYYYMMDD)
		int             nTime;                  //ʱ��(HHMMSSmmm)
		int             nStatus;                //״̬
		int64_t         iPreOpenInterest;       //��ֲ�
		unsigned int    nPreClose;              //�����̼�=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		unsigned int    nPreSettlePrice;        //�����=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		unsigned int    nOpen;                  //���̼�=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		unsigned int    nHigh;                  //��߼�=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		unsigned int    nLow;                   //��ͼ�=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		unsigned int    nMatch;                 //���¼�=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		int64_t         iVolume;                //�ɽ�����=ʵ������(��λ: ��)
		int64_t         iTurnover;              //�ɽ��ܽ��=ʵ�ʽ��(��λ: Ԫ)
		int64_t         iOpenInterest;          //�ֲ�����=ʵ������(��λ: ��)
		unsigned int    nClose;                 //������=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		unsigned int    nSettlePrice;           //�����=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		unsigned int    nHighLimited;           //��ͣ��=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		unsigned int    nLowLimited;            //��ͣ��=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		unsigned int    nAskPrice[5];           //������=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		unsigned int    nAskVol[5];             //������=ʵ������(��λ: ��)
		unsigned int    nBidPrice[5];           //�����=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		unsigned int    nBidVol[5];             //������=ʵ������(��λ: ��)
	};

	// ��ע: ���ṹ����ֶ��еĵ�λΪ�ο���λ, ��������, ���������ϳ���������������̿ڵ�λΪ׼
	struct MD_DATA_ORDER_QUEUE
	{
		char    szWindCode[32]; //600001.SH
		char    szCode[32];     //ԭʼCode
		int     nActionDay;     //����(YYYYMMDD)
		int     nTime;          //ʱ��(HHMMSSmmm)
		int     nSide;          //��������('B':Bid 'A':Ask)
		int     nPrice;         //ί�м۸�=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		int     nOrders;        //�ҵ���λ
		int     nABItems;       //ί�е���
		int     nABVolume[200]; //ί������=ʵ�ʹ���(��λ: ��)
	};

	// ��ע: ���ṹ����ֶ��еĵ�λΪ�ο���λ, ��������, ���������ϳ���������������̿ڵ�λΪ׼
	struct MD_DATA_TRANSACTION
	{
		char    szWindCode[32]; //600001.SH
		char    szCode[32];     //ԭʼCode
		int     nActionDay;     //�ɽ�����(YYYYMMDD)
		int     nTime;          //�ɽ�ʱ��(HHMMSSmmm)
		int     nIndex;         //�ɽ����
		int     nPrice;         //�ɽ��۸�=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		int     nVolume;        //�ɽ�����=ʵ�ʹ���(��λ: ��)
		int     nTurnover;      //�ɽ����=ʵ�ʽ��(��λ: Ԫ)
		int     nBSFlag;        //��������(��'B', ����'A', ������' ')
		char    chOrderKind;    //�ɽ����
		int     nAskOrder;      //������ί�����
		int     nBidOrder;      //����ί�����
	};

	// ��ע: ���ṹ����ֶ��еĵ�λΪ�ο���λ, ��������, ���������ϳ���������������̿ڵ�λΪ׼
	struct MD_DATA_ORDER
	{
		char    szWindCode[32]; //600001.SH
		char    szCode[32];     //ԭʼCode
		int     nActionDay;     //ί������(YYYYMMDD)
		int     nTime;          //ί��ʱ��(HHMMSSmmm)
		int     nOrder;         //ί�к�
		int     nPrice;         //ί�м۸�=ʵ�ʼ۸�(��λ: Ԫ/��)x10000
		int     nVolume;        //ί������=ʵ�ʹ���(��λ: ��)
		char    chOrderKind;    //ί�����
		char    chFunctionCode; //ί�д���('B','S','C')
	};

	// ��ע: ���ṹ����ֶ��еĵ�λΪ�ο���λ, ��������, ���������ϳ���������������̿ڵ�λΪ׼
	struct MD_DATA_KLINE                    //ģ����������
	{
		MD_CycType nType;                   //��������
		char    szWindCode[32];             //600001.SH
		char    szCode[32];                 //ԭʼCode
		MD_ISODateTimeType szDatetime;      //ʱ��
		int nDate;                          //����(YYYYMMDD)
		int nTime;                          //ʱ��(HHMMSSmmm)
		double  nOpen;                      //���̼�=ʵ�ʼ۸�(��λ: Ԫ/��)
		double  nHigh;                      //��߼�=ʵ�ʼ۸�(��λ: Ԫ/��)
		double  nLow;                       //��ͼ�=ʵ�ʼ۸�(��λ: Ԫ/��)
		double  nClose;                     //���ռ�=ʵ�ʼ۸�(��λ: Ԫ/��)
		double  nPreClose;                  //���ռ�=ʵ�ʼ۸�(��λ: Ԫ/��)
		double  nHighLimit;                 //��ͣ��=ʵ�ʼ۸�(��λ: Ԫ/��)
		double  nLowLimit;                  //��ͣ��=ʵ�ʼ۸�(��λ: Ԫ/��)
		int64_t iVolume;                    //�ɽ�����=ʵ�ʹ���(��λ: ��)
		int64_t nTurover;                   //�ɽ����=ʵ�ʽ��(��λ: Ԫ)
	};
#pragma pack(pop)
	/// <summary>
	/// ���׽ṹ��
	/// </summary>
	typedef char TD_CodeType[32];
	typedef char TD_CodeName[32];
	typedef char TD_ISODateTimeType[21];    //���ں�ʱ������(��ʽ yyyy-MM-dd hh:mm:ss)
	typedef char TD_OrderIdType[64];
	typedef char TD_AccountType[64];
	typedef char TD_PassType[64];
	typedef char TD_Text[128];
	typedef char TD_ClientIdType[128];

	enum TD_TradeType
	{
		TD_TradeType_None,
		TD_TradeType_Sell,      //����
		TD_TradeType_Buy        //����
	};

	enum TD_OffsetType
	{
		TD_OffsetType_None,
		TD_OffsetType_Open,                 //����
		TD_OffsetType_Close,                //ƽ��
	};

	enum TD_OrderStatusType
	{
		TD_OrderStatusType_fail = -10,          ///ָ��ʧ��
		TD_OrderStatusType_removed,             ///�����ɹ�
		TD_OrderStatusType_allDealed,           ///ȫ���ɽ�

		TD_OrderStatusType_unAccpet = 0,        ///δ����
		TD_OrderStatusType_accpeted,            ///�ѽ���δ����
		TD_OrderStatusType_queued,              ///�����Ŷ�  (������״̬)
		//      TD_OrderStatusType_toModify,            ///�����ĵ�
		//      TD_OrderStatusType_modifing,            ///�ѱ��ĵ�
		//      TD_OrderStatusType_modified,            ///�ĵ�����
		TD_OrderStatusType_toRemove,            ///��������
		TD_OrderStatusType_removing,            ///�ѱ�����
		TD_OrderStatusType_partRemoved,         ///���ֳ���
		TD_OrderStatusType_partDealed,          ///���ֳɽ�
	};

	struct TD_OrderDetail
	{
		///ί�б�ţ�broker ��������Ψһ��ţ�
		TD_OrderIdType  szOrderStreamId;
		///ȯ���ʽ��˻�Id
		int             nAccountId;                 //  ����ָ���ʽ��˺��µ������ֶ����µ���ʱ����Ҫ��д
		///�ʽ��˻�����
		TD_AccountType  szAccountNickName;

		///ί���걨��
		int             nOrderVol;                  //  ����ָ���ʽ��˺��µ������ֶ����µ���ʱ����Ҫ��д

		///�ɽ�����  * 10000
		int64_t             nDealedPrice;
		///�ɽ���
		int             nDealedVol;

		///��������
		int             nWithDrawnVol;

		///ί��ʱ��
		TD_ISODateTimeType  szOrderTime;
		/// ״̬
		TD_OrderStatusType  nStatus;
		//��ע
		TD_Text             szText;

		//������
		double              nFee;
		TD_OrderDetail()
		{
			memset(szOrderStreamId, 0, sizeof(TD_OrderIdType));
			nAccountId = 0;
			memset(szAccountNickName, 0, sizeof(TD_AccountType));

			nOrderVol = 0;

			nDealedPrice = 0;
			nDealedVol = 0;

			nWithDrawnVol = 0;
			nStatus = TD_OrderStatusType_unAccpet;
			memset(szOrderTime, 0, sizeof(TD_ISODateTimeType));
			memset(szText, 0, sizeof(TD_Text));

			nFee = 0;
		}
		virtual ~TD_OrderDetail()
		{

		}
	};

	struct TD_Base_Msg
	{
		///����ID���пͻ���APIά����ΨһID��
		int     nReqId;

		int64_t nStrategyId;

		///�û������ֶ�
		int     nUserInt;
		double  nUserDouble;
		TD_Text szUseStr;

		///�û��ʺ�Id
		int64_t     nUserId;
		///�ͻ��˱��
		TD_ClientIdType  szClientId;

		TD_Base_Msg()
		{
			nReqId = 0;
			nStrategyId = 0;

			nUserInt = 0;
			nUserDouble = 0;
			memset(szUseStr, 0, sizeof(TD_Text));

			nUserId = 0;
			memset(szClientId, 0, sizeof(TD_ClientIdType));
		}
		virtual ~TD_Base_Msg()
		{

		}
	};

	// �û�Ȩ��
	class TD_QuantUserAuthen : public TD_Base_Msg
	{
	public:
		int64_t                        nId;
		int                         nGroupId;

		bool                        ifStopTrade;                     // �Ƿ�ͣ��
		int                        nStopTradePostion;                // ͣ��λ(�����ʽ���)
		double              nStopPercentTradePostion;                // ͣ��λ(�������)

		int                         nSinglePositionHoldTime;         // ���ʳֲ�ʱ����ֵ
		int                         nSinglePositionLoss;             // ���ʳֲֿ�����ֵ(�����ʽ���)
		double               nSinglePercentPositionLoss;             // ���ʳֲֿ�����ֵ(�������)
		TD_QuantUserAuthen()
		{
			nId = 0;
			nGroupId = 0;

			ifStopTrade = false;
			nStopTradePostion = 1000000;
			nStopPercentTradePostion = 0.1;

			nSinglePositionHoldTime = 1000000;
			nSinglePositionLoss = 1000000;
			nSinglePercentPositionLoss = 0.1;
		}
		virtual ~TD_QuantUserAuthen()
		{

		}
	};

	class TD_QuantUserCodeInfo
	{
	public:
		char               szWinCode[64];
		double             nCaptial;                   // �������ʽ�
		int                nLendingAmount;             // ������ȯ
	};

	// �û���Ʊ��
	class TD_QuantUserCodePool : public TD_Base_Msg
	{
	public:
		int64_t                        nId;
		int                         nGroupId;

		int                      nCodeControlNum;       //  ����ȯ����
		TD_QuantUserCodeInfo* pCodeControl;          // ������ȯ��Ϣ��ָ��ƫ��

		TD_QuantUserCodePool()
		{
			nId = 0;
			nGroupId = 0;

			nCodeControlNum = 0;
			pCodeControl = NULL;

		}
		~TD_QuantUserCodePool()
		{
			if(pCodeControl)
			{
				delete[] pCodeControl;
				pCodeControl = NULL;
			}
		}
	};

	// �û���Ʊ��
	class TD_QuantUserDisablePublicCode : public TD_Base_Msg
	{
	public:
		int64_t                        nId;
		int                         nGroupId;

		int                      nDisablePublicCodeNum;          //  �����ù���ȯ����
		TD_CodeType*             pDisablePublicCode;             //  �����ù���ȯָ��

		TD_QuantUserDisablePublicCode()
		{
			nId = 0;
			nGroupId = 0;

			nDisablePublicCodeNum = 0;
			pDisablePublicCode = NULL;

		}
		~TD_QuantUserDisablePublicCode()
		{
			if(pDisablePublicCode)
			{
				delete[] pDisablePublicCode;
				pDisablePublicCode = NULL;
			}
		}
	};

	struct TD_Login
	{
		int64_t nId;
		TD_PassType szPass;
		TD_Text szMachineId;
		TD_Login()
		{
			nId = 0;
			memset(szPass, 0, sizeof(szPass));
			memset(szMachineId, 0, sizeof(szPass));
		}
		virtual ~TD_Login()
		{

		}
	};

	// �µ�
	struct TD_ReqOrderInsert : TD_Base_Msg
	{
		///֤ȯ��Լ����
		TD_CodeType     szContractCode;
		///֤ȯ��Լ����
		TD_CodeName     szContractName;
		///�������� ����
		TD_TradeType    nTradeType;
		//��ƽ������
		TD_OffsetType   nOffsetType;
		///ί�м�  *10000
		int64_t             nOrderPrice;
		///ί����
		int             nOrderVol;
		///  �����ֱ���
		int            nOrderNum;
		///  ������ϸ��ָ��ƫ�ƣ�
		TD_OrderDetail*    pOrderDetail;                // ����ָ���ʽ��˺��µ�����Ҫ���Ӧ���ʽ��˺�id�Լ��µ����ʽ��˺ŵ�ί������

		//�Ƿ��Ƿ��ǿƽ
		int                nCloseR;           // 0����ƽ��,1Ϊ��ظ�Ԥƽ��,2Ϊ��������ز��Դﵽǿƽλƽ��
		TD_ReqOrderInsert()
		{

			memset(szContractCode, 0, sizeof(TD_CodeType));
			memset(szContractName, 0, sizeof(TD_CodeName));
			nTradeType = TD_TradeType_None;
			nOffsetType = TD_OffsetType_None;

			nOrderPrice = 0;
			nOrderVol = 0;

			nOrderNum = 0;
			pOrderDetail = NULL;
			nCloseR = 0;
		}

		virtual ~TD_ReqOrderInsert()
		{

		}
	};

	struct TD_RspOrderInsert : TD_ReqOrderInsert
	{
		/// ���������û�UserId
		int64_t         nOrderOwnerId;
		///������ά��������Ψһ��
		int64_t         nOrderId;

		///�ύ�걨��
		int             nSubmitVol;
		///�ɽ�����  * 10000
		int64_t             nDealedPrice;
		///�ɽ�����
		int             nDealedVol;

		///��������
		int             nTotalWithDrawnVol;

		//�ϵ�����
		int             nInValid;
		/// ״̬
		TD_OrderStatusType  nStatus;
		/// �µ�ʱ��
		TD_ISODateTimeType  szInsertTime;

		//������
		double               nFee;
		TD_RspOrderInsert()
		{
			nOrderOwnerId = 0;
			nOrderId = 0;

			nSubmitVol = 0;
			nDealedPrice = 0;
			nDealedVol = 0;

			nTotalWithDrawnVol = 0;
			nInValid = 0;

			nStatus = TD_OrderStatusType_unAccpet;
			memset(szInsertTime, 0, sizeof(TD_ISODateTimeType));

			nFee = 0;
		}
		virtual ~TD_RspOrderInsert()
		{

		}
	};

	// ����
	struct TD_ReqOrderDelete : TD_Base_Msg
	{
		///ԭʼ����������ΨһId
		int64_t         nOrderId;
		///ί�б�ţ�broker ��������Ψһ��ţ�
		TD_OrderIdType  szOrderStreamId;
		TD_ReqOrderDelete()
		{
			nOrderId = 0;
			memset(szOrderStreamId, 0, sizeof(TD_OrderIdType));
		}
		virtual ~TD_ReqOrderDelete()
		{

		}
	};

	///����������Ӧ
	typedef TD_ReqOrderDelete TD_RspOrderDelete;
	//ί�в�ѯ����
	struct TD_ReqQryPriority : TD_Base_Msg
	{
		/// ��Ҫ��ѯ���û�Id
		int64_t               nQueryUserid;

		TD_ReqQryPriority()
		{
			nQueryUserid = 0;
		}
		virtual ~TD_ReqQryPriority()
		{

		}
	};

	struct TD_Priority
	{
		int          nPriority;
		int64_t      nAccountId;
		TD_Priority()
		{

		}
		virtual ~TD_Priority()
		{

		}
	};
	//ί�в�ѯ����
	struct TD_RspQryPriority : TD_Base_Msg
	{
		// �û�id
		int64_t               nQueryUserid;
		///֤ȯ��Լ����
		TD_CodeType           szContractCode;

		int                   nNum;
		TD_Priority*          pPriority;
		TD_RspQryPriority()
		{
			nQueryUserid = 0;
			memset(szContractCode, 0, sizeof(TD_Priority));
			nNum = 0;
			pPriority = NULL;
		}
		virtual ~TD_RspQryPriority()
		{
			if(pPriority != NULL)
			{
				delete[] pPriority;
				pPriority = NULL;
			}
		}
	};
	typedef TD_RspQryPriority TD_ReqUpdatePriority;
	typedef TD_RspQryPriority TD_RspUpdatePriority;
	//ί�в�ѯ����
	struct TD_ReqQryOrder : TD_Base_Msg
	{
		///֤ȯ��Լ����
		TD_CodeType             szContractCode;
		///��ʼλ��(����Ĭ�ϴ�ͷ��ʼ)
		int                     nIndex;
		///����������Ĭ�ϲ�ȫ����
		int                     nNum;

		TD_ReqQryOrder()
		{
			memset(szContractCode, 0, sizeof(TD_CodeType));
			nIndex = 0;
			nNum = 0;
		}
		virtual ~TD_ReqQryOrder()
		{

		}
	};

	//�ɽ���ѯ����
	struct TD_ReqQryMatch : TD_Base_Msg
	{
		///֤ȯ��Լ����
		TD_CodeType     szContractCode;
		///��ʼλ��(����Ĭ�ϴ�ͷ��ʼ)
		int             nIndex;
		///����������Ĭ�ϲ�ȫ����
		int             nNum;

		TD_ReqQryMatch()
		{
			memset(szContractCode, 0, sizeof(TD_CodeType));
			nIndex = 0;
			nNum = 0;
		}
		virtual ~TD_ReqQryMatch()
		{

		}
	};

	///�ֲֲ�ѯ����
	struct TD_ReqQryPosition : TD_Base_Msg
	{
		/// ֤ȯ��Լ����
		TD_CodeType     szContractCode;
		TD_ReqQryPosition()
		{
			memset(szContractCode, 0, sizeof(TD_CodeType));
		}
		virtual ~TD_ReqQryPosition()
		{

		}
	};

	///�ֲֲ�ѯ�ص�
	struct TD_RspQryPosition : TD_Base_Msg
	{
		///  ֤ȯ��Լ����
		TD_CodeType     szContractCode;
		///  �ֲ�����
		int     nPosition;
		///  �ֲ־��� * 10000
		double      nPrice;
		///  ��ӯ
		double  nProfit;
		///  �ѽ����ӯ��
		double nSelltleProfit;
		TD_RspQryPosition()
		{
			memset(szContractCode, 0, sizeof(TD_CodeType));
			nPosition = 0;
			nPrice = 0;
			nProfit = 0;
			nSelltleProfit = 0;
		}
		virtual ~TD_RspQryPosition()
		{

		}
	};

	///����ί������ѯ����
	struct TD_ReqQryMaxEntrustCount : TD_Base_Msg
	{
		/// ��Ʊ����
		TD_CodeType szContractCode;
		TD_ReqQryMaxEntrustCount()
		{
			memset(szContractCode, 0, sizeof(TD_CodeType));
		}
		virtual ~TD_ReqQryMaxEntrustCount()
		{

		}
	};

	///�ʽ��˺�����ί������ѯ����
	struct TD_ReqQryAccountMaxEntrustCount : TD_Base_Msg
	{
		/// ��Ʊ����
		TD_CodeType szContractCode;
		/// �ʽ��˺�
		int         nAccountId;

		TD_ReqQryAccountMaxEntrustCount()
		{
			memset(szContractCode, 0, sizeof(TD_CodeType));
			nAccountId = 0;
		}
		virtual ~TD_ReqQryAccountMaxEntrustCount()
		{

		}
	};

	///  ��ί������
	struct StockMaxEntrustCount
	{
		///  ��Ʊ����
		TD_CodeType szContractCode;
		///  �������ʽ���
		double nMaxBuyCaptial;
		///  ������
		int nMaxSellVol;
		StockMaxEntrustCount()
		{
			memset(szContractCode, 0, sizeof(TD_CodeType));
			nMaxBuyCaptial = 0;
			nMaxSellVol = 0;
		}
		virtual ~StockMaxEntrustCount()
		{

		}
	};

	///  ��ѯ����ί�������ص�
	struct TD_RspQryMaxEntrustCount : public TD_Base_Msg
	{
		StockMaxEntrustCount pStockMaxEntrustCount;
		TD_RspQryMaxEntrustCount()
		{

		}
		virtual ~TD_RspQryMaxEntrustCount()
		{

		}
	};
	///  �ʽ��˺���Ϣ��ʼֵ�ص�
	struct TD_RspQryInitAccountMaxEntrustCount : TD_Base_Msg
	{
		///  �ʽ��˺�
		int nAccountId;
		///  �ʽ��˻�����
		TD_AccountType  szAccountNickName;
		///  ���ù�Ʊ��
		int nNum;
		///  ָ��ƫ��
		StockMaxEntrustCount* pStockMaxEntrustCount;
		///  �ʽ�ͨ���Ƿ����
		bool bStatus;
		///  �ʽ��˺Ų�������ʽ�
		int nAvailableCaptial;

		TD_RspQryInitAccountMaxEntrustCount()
		{
			memset(szAccountNickName, 0, sizeof(TD_AccountType));
			nAccountId = 0;
			nNum = 0;
			pStockMaxEntrustCount = NULL;
			bStatus = 0;
			nAvailableCaptial = 0;
		}
		virtual ~TD_RspQryInitAccountMaxEntrustCount()
		{
			if (pStockMaxEntrustCount)
			{
				delete[] pStockMaxEntrustCount;
			}
		}
	};
	struct StockAccountMaxEntrustCount
	{
		///  ֤ȯ��Լ����
		TD_CodeType     szContractCode;
		///  �ֲ�����
		int     nTotalPosition;
		///  ������
		int     nAvailablePosition;
		///  �ʽ��˺Ų��浱�ճ�ʼ�ֲ�
		int     nLastPosition;

		StockAccountMaxEntrustCount()
		{
			memset(szContractCode, 0, sizeof(TD_CodeType));
			nTotalPosition = 0;
			nAvailablePosition = 0;
			nLastPosition = 0;
		}
		virtual ~StockAccountMaxEntrustCount()
		{

		}
	};

	///  �ʽ��˺�����ί�����ص�
	struct TD_RspQryAccountMaxEntrustCount : TD_Base_Msg
	{
		///  �ʽ��˺�
		int nAccountId;
		///  �ʽ��˻�����
		TD_AccountType  szAccountNickName;
		///  ���ù�Ʊ��
		int nNum;
		///  ָ��ƫ��
		StockAccountMaxEntrustCount* pStockMaxEntrustCount;
		///  �ʽ�ͨ���Ƿ����
		bool bStatus;
		///  �ʽ��˺Ų�������ʽ�
		int nAvailableCaptial;
		//   �ʽ��˺Ų��浱�ճ�ʼ�ʽ�
		int nTotalCaptial;

		TD_RspQryAccountMaxEntrustCount()
		{
			memset(szAccountNickName, 0, sizeof(TD_AccountType));
			nAccountId = 0;
			nNum = 0;
			pStockMaxEntrustCount = NULL;
			bStatus = 0;
			nAvailableCaptial = 0;
			nTotalCaptial = 0;
		}
		virtual ~TD_RspQryAccountMaxEntrustCount()
		{
			if(pStockMaxEntrustCount)
			{
				delete[] pStockMaxEntrustCount;
			}
		}
	};

	///����״̬�仯֪ͨ
	typedef TD_RspOrderInsert TD_RtnOrderStatusChangeNotice;

	///�ɽ��ص�֪ͨ
	struct TD_RtnOrderMatchNotice : TD_Base_Msg
	{
		/// ί�б��
		int64_t             nOrderId;
		///ί�б�ţ�broker ��������Ψһ��ţ�
		TD_OrderIdType      szOrderStreamId;
		///ϵͳ���
		int64_t             nMatchStreamId;
		///�ɽ��� * 10000
		double              nMatchPrice;
		///�ɽ���
		int                 nMatchVol;
		/// ��Ʊ����
		TD_CodeType         szContractCode;
		///֤ȯ��Լ����
		TD_CodeName     szContractName;
		/// �ɽ�ʱ��
		TD_ISODateTimeType  szMatchTime;
		///�������� ����
		TD_TradeType    nTradeType;
		///ȯ���ʽ��˻�Id
		int             nAccountId;
		///�ʽ��˻�����
		TD_AccountType  szAccountNickName;


		TD_RtnOrderMatchNotice()
		{
			nOrderId = 0;
			nMatchStreamId = 0;
			nMatchPrice = 0;
			nMatchVol = 0;
			nTradeType = TD_TradeType_None;

			memset(szContractCode, 0, sizeof(TD_CodeType));
			memset(szMatchTime, 0, sizeof(TD_ISODateTimeType));
			memset(szOrderStreamId, 0, sizeof(TD_OrderIdType));
			memset(szContractName, 0, sizeof(szContractName));
			nAccountId = 0;
			memset(szAccountNickName, 0, sizeof(TD_AccountType));
		}
		virtual ~TD_RtnOrderMatchNotice()
		{

		}
	};

	///ί�в�ѯ��Ӧ
	struct TD_RspQryOrder : public TD_RtnOrderStatusChangeNotice
	{
		int            nIndex;
		TD_RspQryOrder()
		{
			nIndex = 0;
		}
		~TD_RspQryOrder()
		{

		}
	};
	///�ɽ���ѯ��Ӧ
	struct TD_RspQryMatch : public TD_RtnOrderMatchNotice
	{
		int nIndex;
		TD_RspQryMatch()
		{
			nIndex = 0;
		}
		~TD_RspQryMatch()
		{

		}
	};
	///����״̬
	typedef TD_RtnOrderStatusChangeNotice TD_OrderStatus;
	///ӯ�����Ͳ���
	typedef TD_RspQryPosition TD_RtnProfit;

	class TD_SimulationReservation
	{
	public:
		char               szWinCode[64];
		int                nLendingAmount;             // �ײ�����
		double             nPrice;                     // �ײ־���
	public:
		TD_SimulationReservation()
		{
			memset(szWinCode, 0, 64);
			nLendingAmount = 0;
			nPrice = 0;
		}
		virtual ~TD_SimulationReservation()
		{

		}
	};
	// ģ���ʽ��˺���Ϣ
	class TD_SimulationAccount
	{
	public:
		int64_t               nSimAccountId;
		char                  szNickName[64];
		char                  szText[128];

		double                      nTotalAmount;
		int                         nReservationNum;   // �ײֹ�Ʊ����
		TD_SimulationReservation*   pReservationCode;  // �ײֹ�Ʊ����

	public:
		TD_SimulationAccount()
		{
			nSimAccountId = 0;
			memset(szNickName, 0, 64);
			memset(szText, 0, 128);

			nTotalAmount = 0;
			nReservationNum = 0;
			pReservationCode = NULL;
		}
		virtual ~TD_SimulationAccount()
		{
			if(pReservationCode != NULL)
			{
				delete[] pReservationCode;
				pReservationCode = NULL;
			}
		}
	};
}



#endif//_QuantDef_H_
