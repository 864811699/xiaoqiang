package com.tailai.intradata.client.netcomm.proto;

import com.tailai.intradata.client.struct.*;

/**
 * @author wangcf
 */
public class MarketDataMessageRecoder {
    public static MarketData.Message record(TLMdReqFutureInstrumentField reqFutureInstrumentField, int requestId) {
        MarketData.QryFutureInstrumentsRequest.Builder
                qryBuilder =
                MarketData.QryFutureInstrumentsRequest.newBuilder();
        qryBuilder.setExchangeID(reqFutureInstrumentField.getExchangeID());

        MarketData.Request.Builder
                requestBuilder =
                MarketData.Request.newBuilder();
        requestBuilder.setRequestId(requestId);
        requestBuilder.setQryFutureInstruments(qryBuilder);

        MarketData.Message.Builder
                messageBuilder =
                MarketData.Message.newBuilder();
        messageBuilder.setMsgType(MarketData.MsgType.Qry_Future_Instruments_Request);
        messageBuilder.setSequence(requestId);
        messageBuilder.setRequest(requestBuilder);
        return messageBuilder.build();
    }

    public static MarketData.Message record(TLMdReqIndexQuotationField reqIndexQuotationField, int requestId) {
        MarketData.QryIndexQuotationRequest.Builder
                indexQuotationQryBuilder =
                MarketData.QryIndexQuotationRequest.newBuilder();
        indexQuotationQryBuilder.setExchangeID(reqIndexQuotationField.getExchangeID());

        MarketData.Request.Builder
                requestBuilder =
                MarketData.Request.newBuilder();
        requestBuilder.setRequestId(requestId);
        requestBuilder.setQryIndexQuotation(indexQuotationQryBuilder);

        MarketData.Message.Builder
                messageBuilder =
                MarketData.Message.newBuilder();
        messageBuilder.setMsgType(MarketData.MsgType.Qry_Index_Quotation_Request);
        messageBuilder.setSequence(requestId);
        messageBuilder.setRequest(requestBuilder);
        return messageBuilder.build();
    }

    public static MarketData.Message record(TLMdReqMinuteIndicatorField reqMinuteIndicatorField,
                                            int requestId) {
        MarketData.QryMinuteIndicatorRequest.Builder builder = MarketData.QryMinuteIndicatorRequest.newBuilder();
        builder.setExchangeID(reqMinuteIndicatorField.getExchangeID());

        MarketData.Request.Builder
                requestBuilder =
                MarketData.Request.newBuilder();
        requestBuilder.setRequestId(requestId);
        requestBuilder.setQryMinuteIndicator(builder);

        MarketData.Message.Builder
                messageBuilder =
                MarketData.Message.newBuilder();
        messageBuilder.setMsgType(MarketData.MsgType.Qry_Minute_Indicator_Request);
        messageBuilder.setSequence(requestId);
        messageBuilder.setRequest(requestBuilder);
        return messageBuilder.build();
    }

    public static MarketData.Message record(TLMdReqFutureQuotationField reqFutureQuotationField,
                                            int requestId) {
        MarketData.QryFutureQuotationRequest.Builder
                quotationQryBuilder =
                MarketData.QryFutureQuotationRequest.newBuilder();
        quotationQryBuilder.setExchangeID(reqFutureQuotationField.getExchangeID());

        MarketData.Request.Builder
                requestBuilder =
                MarketData.Request.newBuilder();
        requestBuilder.setRequestId(requestId);
        requestBuilder.setQryFutureQuotation(quotationQryBuilder);

        MarketData.Message.Builder
                messageBuilder =
                MarketData.Message.newBuilder();
        messageBuilder.setMsgType(MarketData.MsgType.Qry_Future_Quotation_Request);
        messageBuilder.setSequence(requestId);
        messageBuilder.setRequest(requestBuilder);
        return messageBuilder.build();

    }

    public static MarketData.Message record(TLMdReqSecuInstrumentField reqSecuInstrumentField, int requestId) {
        MarketData.QrySecuInstrumentsRequest.Builder
                qryBuilder =
                MarketData.QrySecuInstrumentsRequest.newBuilder();
        qryBuilder.setExchangeID(reqSecuInstrumentField.getExchangeID());

        MarketData.Request.Builder
                requestBuilder =
                MarketData.Request.newBuilder();
        requestBuilder.setRequestId(requestId);
        requestBuilder.setQrySecuInstruments(qryBuilder);

        MarketData.Message.Builder
                messageBuilder =
                MarketData.Message.newBuilder();
        messageBuilder.setMsgType(MarketData.MsgType.Qry_Secu_Instruments_Request);
        messageBuilder.setSequence(requestId);
        messageBuilder.setRequest(requestBuilder);
        return messageBuilder.build();
    }

    public static MarketData.Message record(TLMdReqStockQuotationField reqStockQuotationField, int requestId) {
        MarketData.QryStockQuotationRequest.Builder
                stockQuotationQryBuilder =
                MarketData.QryStockQuotationRequest.newBuilder();
        stockQuotationQryBuilder.setExchangeID(reqStockQuotationField.getExchangeID());

        MarketData.Request.Builder
                requestBuilder =
                MarketData.Request.newBuilder();
        requestBuilder.setRequestId(requestId);
        requestBuilder.setQryStockQuotation(stockQuotationQryBuilder);

        MarketData.Message.Builder
                messageBuilder =
                MarketData.Message.newBuilder();
        messageBuilder.setMsgType(MarketData.MsgType.Qry_Stock_Quotation_Request);
        messageBuilder.setSequence(requestId);
        messageBuilder.setRequest(requestBuilder);
        return messageBuilder.build();
    }

    public static MarketData.Message record(TLMdReqTickIndicatorField reqSecuSimpleIndicatorField,
                                            int requestId) {
        MarketData.QryTickIndicatorRequest.Builder
                simpleIndicatorQryBuilder =
                MarketData.QryTickIndicatorRequest.newBuilder();
        simpleIndicatorQryBuilder.setExchangeID(reqSecuSimpleIndicatorField.getExchangeID());

        MarketData.Request.Builder
                requestBuilder =
                MarketData.Request.newBuilder();
        requestBuilder.setRequestId(requestId);
        requestBuilder.setQryTickIndicator(simpleIndicatorQryBuilder);

        MarketData.Message.Builder
                messageBuilder =
                MarketData.Message.newBuilder();
        messageBuilder.setMsgType(MarketData.MsgType.Qry_Tick_Indicator_Request);
        messageBuilder.setSequence(requestId);
        messageBuilder.setRequest(requestBuilder);
        return messageBuilder.build();
    }
}
