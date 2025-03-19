package com.tailai.intradata.client.util;

import java.nio.charset.Charset;
/**
 * 常量类
 * @author Administrator
 *
 */
public class Constants {
    public static final String CHARSET_NAME = "UTF-8";
    public static Charset charset = Charset.forName(CHARSET_NAME);

    public static final String MSGPROPERTIES_DATATYPE = "DataType";
    public static final String MSGPROPERTIES_MARKETID = "MarketId";
    public static final String MSGPROPERTIES_CLIENTID = "ClientID";

    public static final String THEME_REQUEST = "QUOTATION.REQUEST.QUEUE";
    public static final String THEME_RESPONSE = "QUOTATION.RESPONSE.QUEUE";
    public static final String THEME_SUBSCRIBE_STOCK_SNAPSHOT = "QUOTATION.STOCK.SNAPSHOT";
    public static final String THEME_SUBSCRIBE_STOCK_TRANSACTION = "QUOTATION.STOCK.TRANSACTION";
    public static final String THEME_SUBSCRIBE_FUTURE_SNAPSHOT = "QUOTATION.FUTURE.SNAPSHOT";
    public static final String THEME_SUBSCRIBE_INDICATOR_TICK = "QUOTATION.INDICATOR.TICK";
    public static final String THEME_SUBSCRIBE_INDICATOR_MINUTE = "QUOTATION.INDICATOR.MINUTE";
    public static final String THEME_SUBSCRIBE_SIGNAL_MINUTE = "QUOTATION.SIGNAL.MINUTE";
    public static final String THEME_SUBSCRIBE_SIGNAL_TICK = "QUOTATION.SIGNAL.TICK";

    public static final String QUOTATION_CONNECTED = "CONNECTED";
    public static final String QUOTATION_DISCONNECTED = "DISCONNECTED";

    public static final int StkType_UnKnown = -1; //未知
    public static final int StkType_Index_Exchange = 0x01; //交易所指数
    public static final int StkType_Stock_AS = 0x10; //A股
    public static final int StkType_Stock_SME = 0x11;//中小板
    public static final int StkType_Stock_GEM = 0x12;//创业板
    public static final int StkType_Stock_KCB = 0x19;//科创板
    public static final int StkType_Fund_ETF = 0x23;//交易型开放式指数基金
    public static final int StkType_Future_Index = 0x70;//指数期货
    public static final int StkType_Future_SP = 0x71;//商品期货
    public static final int StkType_Option_Call_Index = 0x92;///指数认购期权
    public static final int StkType_Option_Put_Index = 0x93;///指数认沽期权
}
