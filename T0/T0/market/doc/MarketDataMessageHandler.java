package com.tailai.intradata.client.netcomm;

import com.tailai.intradata.client.TLMdSpi;
import com.tailai.intradata.client.netcomm.proto.MarketData;
import com.tailai.intradata.client.struct.*;
import com.tailai.intradata.client.util.Constants;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * @author wangcf
 */
public class MarketDataMessageHandler extends MyRunnable<MarketData.Message>
        implements RPCClientMessageHandler {
    private static final Logger logger = LoggerFactory.getLogger(MarketDataMessageHandler.class);

    private TLMdSpi spi;

    public MarketDataMessageHandler(TLMdSpi spi) {
        super("MarketDataMessageHandler");
        enablePerf();
        this.spi = spi;
    }

    @Override
    public void onConnected(String mqAddress) {
        if (spi != null) {
            spi.onConnected();
        }
    }

    @Override
    public void onDisConnected(String mqAddress) {
        if (spi != null) {
            spi.onDisconnected(mqAddress);
        }
    }

    @Override
    public void onProtoMessage(MarketData.Message message) {
        try {
            super.offer(message);
        } catch (InterruptedException e) {
            logger.error(e.getMessage(), e);
        }
    }

    @Override
    protected void doIt(MarketData.Message msg) {
        try {
            this.clientMsgIn(msg);
        } catch (Exception ex) {
            logger.error(ex.getMessage(), ex);
        }
    }

    private void clientMsgIn(MarketData.Message msg) {
        switch (msg.getMsgType()) {
            case Notify_Quotation_Stock:
                handleStockQuotation(msg);
                break;
            case Notify_Quotation_Index:
                handleIndexQuotation(msg);
                break;
            case Notify_Quotation_Future:
                handleFutureQuotation(msg);
                break;
            case Notify_Indicator_Tick:
                handleTickIndicator(msg);
                break;
            case Notify_Indicator_Minute:
                handleMinuteIndicator(msg);
                break;
            case Notify_Signal_Minute:
                handleMinuteSignal(msg);
                break;
            case Notify_Stock_Order_Queue:
                handleOrderQueue(msg);
                break;
            case Notify_Transaction_Order:
                handleTransactionOrder(msg);
                break;
            case Notify_Transaction_Trade:
                handleTransactionTrade(msg);
                break;
            case Notify_Signal_Tick:
                handleTickSignal(msg);
                break;
            case Notify_Secu_Instrument:
                handleSecuInstrument(msg);
                break;
            case Notify_Future_Instrument:
                handleFutureInstrument(msg);
                break;
            case Qry_Secu_Instruments_Response:
                handleQrySecuInstrumentsResponse(msg);
                break;
            case Qry_Future_Instruments_Response:
                handleQryFutureInstrumentsResponse(msg);
                break;
            case Qry_Stock_Quotation_Response:
                handleQryStockQuotationResponse(msg);
                break;
            case Qry_Index_Quotation_Response:
                handleQryIndexQuotationResponse(msg);
                break;
            case Qry_Future_Quotation_Response:
                handleQryFutureQuotationResponse(msg);
                break;
            case Qry_Tick_Indicator_Response:
                handleQryTickIndicatorResponse(msg);
                break;
            case Qry_Minute_Indicator_Response:
                handleQryMinuteIndicatorResponse(msg);
                break;
            default:
                break;
        }
    }

    private void handleTickSignal(MarketData.Message msg) {
        if (spi == null) {
            return;
        }
        MarketData.TickSignal tickSignal = msg.getNotify().getTickSignal();
        TLMdTickSignalField target = new TLMdTickSignalField();
        target.setExchangeID(tickSignal.getExchangeID());
        target.setInstrumentID(tickSignal.getInstrumentID());
        target.setTime(tickSignal.getTime());
        target.setSignals(tickSignal.getSignalsMap());
        spi.onRtnTickSignal(target);
    }

    private void handleMinuteSignal(MarketData.Message msg) {
        if (spi == null) {
            return;
        }
        MarketData.MinuteSignal minuteSignal = msg.getNotify().getMinuteSignal();
        TLMdMinuteSignalField target = new TLMdMinuteSignalField();
        target.setExchangeID(minuteSignal.getExchangeID());
        target.setInstrumentID(minuteSignal.getInstrumentID());
        target.setMinuteTime(minuteSignal.getTime());
        target.setSignals(minuteSignal.getSignalsMap());
        spi.onRtnMinuteSignal(target);
    }

    private void handleFutureInstrument(MarketData.Message msg) {
        if (spi == null) {
            return;
        }

        MarketData.FutureInstrument instrument = msg.getNotify().getFutureInstrument();
        TLMdFutureInstrumentField target = new TLMdFutureInstrumentField();
        target.setExchangeID(instrument.getExchangeID());
        target.setInstrumentID(instrument.getInstrumentID());
        target.setInstrumentName(instrument.getInstrumentName().toString(Constants.charset));
        target.setHighLimited(instrument.getHighLimited());
        target.setLowLimited(instrument.getLowLimited());
        target.setPreClosePrice(instrument.getPreClosePrice());
        target.setOpenPrice(instrument.getOpenPrice());
        target.setSecurityType(instrument.getStkType());
        target.setPreSettlePrice(instrument.getPreSettlePrice());
        spi.onRtnFutureInstrument(target);
    }

    private void handleSecuInstrument(MarketData.Message msg) {
        if (spi == null) {
            return;
        }

        MarketData.SecuInstrument instrument = msg.getNotify().getSecuInstrument();
        TLMdSecuInstrumentField target = new TLMdSecuInstrumentField();
        target.setExchangeID(instrument.getExchangeID());
        target.setInstrumentID(instrument.getInstrumentID());
        target.setInstrumentName(instrument.getInstrumentName().toString(Constants.charset));
        target.setHighLimited(instrument.getHighLimited());
        target.setLowLimited(instrument.getLowLimited());
        target.setPreClosePrice(instrument.getPreClosePrice());
        target.setOpenPrice(instrument.getOpenPrice());
        target.setSecurityType(instrument.getStkType());
        spi.onRtnSecuInstrument(target);
    }

    private void handleTransactionOrder(MarketData.Message msg) {
        if (spi == null) {
            return;
        }
        MarketData.SecuTransactionOrder order = msg.getNotify().getTransactionOrder();
        TLMdSecuTransactionOrderField target = new TLMdSecuTransactionOrderField();
        target.setExchangeID(order.getExchangeID());
        target.setInstrumentID(order.getInstrumentID());
        target.setTime(order.getTime());
        target.setPrice(order.getPrice());
        target.setVolume(order.getVolume());
        target.setOrderKind(order.getOrderKind());
        target.setFunctionCode(order.getFunctionCode());
        target.setOrderNum(order.getOrderNum());
        target.setServerReceivedTime(order.getServerReceivedTime());
        target.setServerSentTime(order.getServerSentTime());
        spi.onRtnSecuTransactionOrder(target);
    }

    private void handleTransactionTrade(MarketData.Message msg) {
        if (spi == null) {
            return;
        }
        MarketData.SecuTransactionTrade trade = msg.getNotify().getTransactionTrade();
        TLMdSecuTransactionTradeField target = new TLMdSecuTransactionTradeField();
        target.setExchangeID(trade.getExchangeID());
        target.setInstrumentID(trade.getInstrumentID());
        target.setTime(trade.getTime());
        target.setBsFlag(trade.getBsFlag());
        target.setPrice(trade.getPrice());
        target.setVolume(trade.getVolume());
        target.setTurnover(trade.getTurnover());
        target.setOrderKind(trade.getOrderKind());
        target.setFunctionCode(trade.getFunctionCode());
        target.setAskOrder(trade.getAskOrder());
        target.setBidOrder(trade.getBidOrder());
        target.setIndex(trade.getIndex());
        target.setServerReceivedTime(trade.getServerReceivedTime());
        target.setServerSentTime(trade.getServerSentTime());
        spi.onRtnSecuTransactionTrade(target);
    }

    private void handleOrderQueue(MarketData.Message msg) {
        if (spi == null) {
            return;
        }
        MarketData.SecuOrderQueue orderQueue = msg.getNotify().getOrderQueue();
        TLMdSecuOrderQueueField target = new TLMdSecuOrderQueueField();
        target.setExchangeID(orderQueue.getExchangeID());
        target.setInstrumentID(orderQueue.getInstrumentID());
        target.setTime(orderQueue.getTime());
        target.setABItems(orderQueue.getABItems());
        target.setOrders(orderQueue.getOrders());
        target.setSide(orderQueue.getSide());
        target.setPrice(orderQueue.getPrice());
        int[] v = new int[orderQueue.getABVolumeCount()];
        for (int i = 0; i < v.length; i++) {
            v[i] = orderQueue.getABVolume(i);
        }
        target.setABVolume(v);
        target.setServerReceivedTime(orderQueue.getServerReceivedTime());
        target.setServerSentTime(orderQueue.getServerSentTime());
        spi.onRtnSecuOrderQueue(target);
    }

    private void handleMinuteIndicator(MarketData.Message msg) {
        if (spi == null) {
            return;
        }
        MarketData.MinuteIndicator minuteIndicator = msg.getNotify().getMinuteIndicator();
        TLMdMinuteIndicatorField target = new TLMdMinuteIndicatorField();
        target.setExchangeID(minuteIndicator.getExchangeID());
        target.setInstrumentID(minuteIndicator.getInstrumentID());
        target.setMinuteTime(minuteIndicator.getTime());
        target.setIndicators(minuteIndicator.getIndicatorsMap());
        spi.onRtnMinuteIndicator(target);
    }

    private void handleTickIndicator(MarketData.Message msg) {
        if (spi == null) {
            return;
        }
        MarketData.TickIndicator tickIndicator = msg.getNotify().getTickIndicator();
        TLMdTickIndicatorField target = new TLMdTickIndicatorField();
        target.setExchangeID(tickIndicator.getExchangeID());
        target.setInstrumentID(tickIndicator.getInstrumentID());
        target.setTime(tickIndicator.getTime());
        target.setAccuAskSize(tickIndicator.getAccuAskSize());
        target.setAccuBidSize(tickIndicator.getAccuBidSize());
        target.setAccuFillSize(tickIndicator.getAccuFillSize());
        target.setAccuFillSizePerSec(tickIndicator.getAccuFillSizePerSec());
        target.setAvgAskSize1(tickIndicator.getAvgAskSize1());
        target.setAvgAskSize2(tickIndicator.getAvgAskSize2());
        target.setAvgBidSize1(tickIndicator.getAvgBidSize1());
        target.setAvgBidSize2(tickIndicator.getAvgBidSize2());
        target.setAvgSpreed(tickIndicator.getAvgSpreed());
        target.setFillSize(tickIndicator.getFillSize());
        target.setFillSizePerSec(tickIndicator.getFillSizePerSec());
        target.setMinVolatility(tickIndicator.getMinVolatility());
        target.setVolumeInterval(tickIndicator.getVolumeInterval());
        target.setUpdownTick(tickIndicator.getUpdownTick());
        target.setNumTradesInterval(tickIndicator.getNumTradesInterval());
        target.setMoneyFlowSmall(tickIndicator.getMoneyFlowSmall());
        target.setMoneyFlowLarge(tickIndicator.getMoneyFlowLarge());
        spi.onRtnTickIndicator(target);
    }

    private void handleFutureQuotation(MarketData.Message msg) {
        if (spi == null) {
            return;
        }
        MarketData.FutureQuotation futureQuotation = msg.getNotify().getFutureQuotation();
        TLMdFutureQuotationField target = new TLMdFutureQuotationField();
        target.setExchangeID(futureQuotation.getExchangeID());
        target.setInstrumentID(futureQuotation.getInstrumentID());
        target.setHighPrice(futureQuotation.getHighPrice());
        target.setLowPrice(futureQuotation.getLowPrice());
        target.setTime(futureQuotation.getTime());
        target.setLastPrice(futureQuotation.getLastPrice());
        target.setTotalVolume(futureQuotation.getTotalVolume());
        target.setTotalAmount((long) futureQuotation.getTotalValue());
        for (int i = 0; i < futureQuotation.getAskPriceCount(); i++) {
            target.getAskPrice()[i] = futureQuotation.getAskPrice(i);
        }
        for (int i = 0; i < futureQuotation.getAskSizeCount(); i++) {
            target.getAskSize()[i] = (int) futureQuotation.getAskSize(i);
        }
        for (int i = 0; i < futureQuotation.getBidPriceCount(); i++) {
            target.getBidPrice()[i] = futureQuotation.getBidPrice(i);
        }
        for (int i = 0; i < futureQuotation.getBidSizeCount(); i++) {
            target.getBidSize()[i] = (int) futureQuotation.getBidSize(i);
        }
        target.setHighLimited(futureQuotation.getHighLimited());
        target.setLowLimited(futureQuotation.getLowLimited());
        target.setPreClosePrice(futureQuotation.getPreClosePrice());
        target.setServerReceivedTime(futureQuotation.getServerReceivedTime());
        target.setServerSentTime(futureQuotation.getServerSentTime());
        spi.onRtnFutureQuotation(target);
    }

    private void handleIndexQuotation(MarketData.Message msg) {
        if (spi == null) {
            return;
        }

        MarketData.IndexQuotation indexQuotation = msg.getNotify().getIndexQuotation();
        TLMdIndexQuotationField target = new TLMdIndexQuotationField();
        target.setExchangeID(indexQuotation.getExchangeID());
        target.setInstrumentID(indexQuotation.getInstrumentID());
        target.setTime(indexQuotation.getTime());
        target.setOpenPrice(indexQuotation.getOpenPrice());
        target.setHighPrice(indexQuotation.getHighPrice());
        target.setLowPrice(indexQuotation.getLowPrice());
        target.setLastPrice(indexQuotation.getLastPrice());
        target.setTotalVolume(indexQuotation.getTotalVolume());
        target.setTotalAmount((long) indexQuotation.getTotalValue());
        target.setPreClosePrice(indexQuotation.getPreClosePrice());
        target.setServerReceivedTime(indexQuotation.getServerReceivedTime());
        target.setServerSentTime(indexQuotation.getServerSentTime());
        spi.onRtnIndexQuotation(target);
    }

    private void handleStockQuotation(MarketData.Message msg) {
        MarketData.StockQuotation
                stockQuotation =
                msg.getNotify().getStockQuotation();

        TLMdStockQuotationField target = new TLMdStockQuotationField();
        target.setExchangeID(stockQuotation.getExchangeID());
        target.setInstrumentID(stockQuotation.getInstrumentID());
        target.setOpenPrice(stockQuotation.getOpenPrice());
        target.setHighPrice(stockQuotation.getHighPrice());
        target.setLowPrice(stockQuotation.getLowPrice());
        target.setTime(stockQuotation.getTime());
        target.setLastPrice(stockQuotation.getLastPrice());
        target.setTotalVolume(stockQuotation.getTotalVolume());
        target.setTotalAmount((long) stockQuotation.getTotalValue());
        target.setTotalNumTrades(stockQuotation.getNumTrades());
        for (int i = 0; i < stockQuotation.getAskPriceCount(); i++) {
            target.getAskPrice()[i] = stockQuotation.getAskPrice(i);
        }
        for (int i = 0; i < stockQuotation.getAskSizeCount(); i++) {
            target.getAskSize()[i] = stockQuotation.getAskSize(i);
        }
        for (int i = 0; i < stockQuotation.getBidPriceCount(); i++) {
            target.getBidPrice()[i] = stockQuotation.getBidPrice(i);
        }
        for (int i = 0; i < stockQuotation.getBidSizeCount(); i++) {
            target.getBidSize()[i] = stockQuotation.getBidSize(i);
        }
        target.setIOPV(stockQuotation.getIOPV());
        target.setHighLimited(stockQuotation.getHighLimited());
        target.setLowLimited(stockQuotation.getLowLimited());
        target.setPreClosePrice(stockQuotation.getPreClosePrice());
        target.setAfterVolume(stockQuotation.getAfterVolume());
        target.setAfterTurnover(stockQuotation.getAfterTurnover());
        target.setAfterPrice(stockQuotation.getAfterPrice());
        target.setAfterMatchItems(stockQuotation.getAfterMatchItems());
        target.setTotalAskVol(stockQuotation.getTotalAskVol());
        target.setTotalBidVol(stockQuotation.getTotalBidVol());
        target.setAvgAskPrice(stockQuotation.getAvgAskPrice());
        target.setAvgBidPrice(stockQuotation.getAvgBidPrice());
        target.setServerReceivedTime(stockQuotation.getServerReceivedTime());
        target.setServerSentTime(stockQuotation.getServerSentTime());
        spi.onRtnStockQuotation(target);
    }

    private void handleQryFutureQuotationResponse(MarketData.Message msg) {
        if (spi == null) {
            return;
        }

        MarketData.Response response = msg.getResponse();
        TLMdRspInfoField info = new TLMdRspInfoField();
        info.setSuccess(response.getSuccess());
        if (!response.getSuccess()) {
            info.setErrorMsg(response.getInfo());
            spi.onRspQryFutureQuotation(null, info, (int) response.getRequestId(), true);
            return;
        }

        MarketData.FutureQuotationList futureQuotationList = response.getFutureQuotationList();
        if (futureQuotationList != null) {
            int count = futureQuotationList.getFutureQuotationsCount();
            int row = 0;
            for (MarketData.FutureQuotation futureQuotation : futureQuotationList.getFutureQuotationsList()) {
                row++;
                TLMdFutureQuotationField target = new TLMdFutureQuotationField();
                target.setExchangeID(futureQuotation.getExchangeID());
                target.setInstrumentID(futureQuotation.getInstrumentID());
                target.setOpenPrice(futureQuotation.getOpenPrice());
                target.setHighPrice(futureQuotation.getHighPrice());
                target.setLowPrice(futureQuotation.getLowPrice());
                target.setTime(futureQuotation.getTime());
                target.setLastPrice(futureQuotation.getLastPrice());
                target.setTotalVolume(futureQuotation.getTotalVolume());
                target.setTotalAmount((long) futureQuotation.getTotalValue());
                target.setTotalNumTrades(futureQuotation.getNumTrades());
                for (int i = 0; i < futureQuotation.getAskPriceCount(); i++) {
                    target.getAskPrice()[i] = futureQuotation.getAskPrice(i);
                }
                for (int i = 0; i < futureQuotation.getAskSizeCount(); i++) {
                    target.getAskSize()[i] = (int) futureQuotation.getAskSize(i);
                }
                for (int i = 0; i < futureQuotation.getBidPriceCount(); i++) {
                    target.getBidPrice()[i] = futureQuotation.getBidPrice(i);
                }
                for (int i = 0; i < futureQuotation.getBidSizeCount(); i++) {
                    target.getBidSize()[i] = (int) futureQuotation.getBidSize(i);
                }
                target.setHighLimited(futureQuotation.getHighLimited());
                target.setLowLimited(futureQuotation.getLowLimited());
                target.setPreClosePrice(futureQuotation.getPreClosePrice());
                boolean isLast = row == count ? true : false;
                spi.onRspQryFutureQuotation(target, info, (int) response.getRequestId(), isLast);
            }
        }
    }

    private void handleQryIndexQuotationResponse(MarketData.Message msg) {
        if (spi == null) {
            return;
        }

        MarketData.Response response = msg.getResponse();
        TLMdRspInfoField info = new TLMdRspInfoField();
        info.setSuccess(response.getSuccess());
        if (!response.getSuccess()) {
            info.setErrorMsg(response.getInfo());
            spi.onRspQryIndexQuotation(null, info, (int) response.getRequestId(), true);
            return;
        }

        MarketData.IndexQuotationList indexQuotationList = response.getIndexQuotationList();
        if (indexQuotationList != null) {
            int count = indexQuotationList.getIndexQuotationsCount();
            int row = 0;
            if (count == 0) {
                spi.onRspQryIndexQuotation(null, info, (int) response.getRequestId(), true);
            }
            for (MarketData.IndexQuotation indexQuotation : indexQuotationList.getIndexQuotationsList()) {
                row++;
                TLMdIndexQuotationField target = new TLMdIndexQuotationField();
                target.setExchangeID(indexQuotation.getExchangeID());
                target.setInstrumentID(indexQuotation.getInstrumentID());
                target.setTime(indexQuotation.getTime());
                target.setOpenPrice(indexQuotation.getOpenPrice());
                target.setHighPrice(indexQuotation.getHighPrice());
                target.setLowPrice(indexQuotation.getLowPrice());
                target.setLastPrice(indexQuotation.getLastPrice());
                target.setTotalVolume(indexQuotation.getTotalVolume());
                target.setTotalAmount((long) indexQuotation.getTotalValue());
                target.setPreClosePrice(indexQuotation.getPreClosePrice());
                boolean isLast = row == count ? true : false;
                spi.onRspQryIndexQuotation(target, info, (int) response.getRequestId(), isLast);
            }
        }
    }

    private void handleQryStockQuotationResponse(MarketData.Message msg) {
        if (spi == null) {
            return;
        }

        MarketData.Response response = msg.getResponse();
        TLMdRspInfoField info = new TLMdRspInfoField();
        info.setSuccess(response.getSuccess());
        if (!response.getSuccess()) {
            info.setErrorMsg(response.getInfo());
            spi.onRspQryStockQuotation(null, info, (int) response.getRequestId(), true);
            return;
        }

        MarketData.StockQuotationList stockQuotationList = response.getStockQuotationList();
        if (stockQuotationList != null) {
            int count = stockQuotationList.getStockQuotationsCount();
            int row = 0;
            for (MarketData.StockQuotation stockQuotation : stockQuotationList.getStockQuotationsList()) {
                row++;
                TLMdStockQuotationField target = new TLMdStockQuotationField();
                target.setExchangeID(stockQuotation.getExchangeID());
                target.setInstrumentID(stockQuotation.getInstrumentID());
                target.setOpenPrice(stockQuotation.getOpenPrice());
                target.setHighPrice(stockQuotation.getHighPrice());
                target.setLowPrice(stockQuotation.getLowPrice());
                target.setTime(stockQuotation.getTime());
                target.setLastPrice(stockQuotation.getLastPrice());
                target.setTotalVolume(stockQuotation.getTotalVolume());
                target.setTotalAmount((long) stockQuotation.getTotalValue());
                target.setTotalNumTrades(stockQuotation.getNumTrades());
                for (int i = 0; i < stockQuotation.getAskPriceCount(); i++) {
                    target.getAskPrice()[i] = stockQuotation.getAskPrice(i);
                }
                for (int i = 0; i < stockQuotation.getAskSizeCount(); i++) {
                    target.getAskSize()[i] = stockQuotation.getAskSize(i);
                }
                for (int i = 0; i < stockQuotation.getBidPriceCount(); i++) {
                    target.getBidPrice()[i] = stockQuotation.getBidPrice(i);
                }
                for (int i = 0; i < stockQuotation.getBidSizeCount(); i++) {
                    target.getBidSize()[i] = stockQuotation.getBidSize(i);
                }
                target.setIOPV(stockQuotation.getIOPV());
                target.setHighLimited(stockQuotation.getHighLimited());
                target.setLowLimited(stockQuotation.getLowLimited());
                target.setPreClosePrice(stockQuotation.getPreClosePrice());
                target.setAfterVolume(stockQuotation.getAfterVolume());
                target.setAfterTurnover(stockQuotation.getAfterTurnover());
                target.setAfterPrice(stockQuotation.getAfterPrice());
                target.setAfterMatchItems(stockQuotation.getAfterMatchItems());
                target.setTotalAskVol(stockQuotation.getTotalAskVol());
                target.setTotalBidVol(stockQuotation.getTotalBidVol());
                target.setAvgAskPrice(stockQuotation.getAvgAskPrice());
                target.setAvgBidPrice(stockQuotation.getAvgBidPrice());
                boolean isLast = row == count ? true : false;
                spi.onRspQryStockQuotation(target, info, (int) response.getRequestId(), isLast);
            }
        }
    }

    private void handleQryTickIndicatorResponse(MarketData.Message msg) {
        if (spi == null) {
            return;
        }

        MarketData.Response response = msg.getResponse();
        TLMdRspInfoField info = new TLMdRspInfoField();
        info.setSuccess(response.getSuccess());
        if (!response.getSuccess()) {
            info.setErrorMsg(response.getInfo());
            spi.onRspQryTickIndicator(null, info, (int) response.getRequestId(), true);
            return;
        }

        MarketData.TickIndicatorList tickIndicatorList = response.getTickIndicatorList();
        if (tickIndicatorList != null) {
            int count = tickIndicatorList.getTickIndicatorsCount();
            int row = 0;
            for (MarketData.TickIndicator tickIndicator : tickIndicatorList.getTickIndicatorsList()) {
                row++;
                TLMdTickIndicatorField target = new TLMdTickIndicatorField();
                target.setExchangeID(tickIndicator.getExchangeID());
                target.setInstrumentID(tickIndicator.getInstrumentID());
                target.setTime(tickIndicator.getTime());
                target.setAccuAskSize(tickIndicator.getAccuAskSize());
                target.setAccuBidSize(tickIndicator.getAccuBidSize());
                target.setAccuFillSize(tickIndicator.getAccuFillSize());
                target.setAccuFillSizePerSec(tickIndicator.getAccuFillSizePerSec());
                target.setAvgAskSize1(tickIndicator.getAvgAskSize1());
                target.setAvgAskSize2(tickIndicator.getAvgAskSize2());
                target.setAvgBidSize1(tickIndicator.getAvgBidSize1());
                target.setAvgBidSize2(tickIndicator.getAvgBidSize2());
                target.setAvgSpreed(tickIndicator.getAvgSpreed());
                target.setFillSize(tickIndicator.getFillSize());
                target.setFillSizePerSec(tickIndicator.getFillSizePerSec());
                target.setMinVolatility(tickIndicator.getMinVolatility());
                target.setVolumeInterval(tickIndicator.getVolumeInterval());
                target.setUpdownTick(tickIndicator.getUpdownTick());
                target.setNumTradesInterval(tickIndicator.getNumTradesInterval());
                target.setMoneyFlowSmall(tickIndicator.getMoneyFlowSmall());
                target.setMoneyFlowLarge(tickIndicator.getMoneyFlowLarge());
                boolean isLast = row == count ? true : false;
                spi.onRspQryTickIndicator(target, info, (int) response.getRequestId(), isLast);
            }
        }
    }

    private void handleQryMinuteIndicatorResponse(MarketData.Message msg) {
        if (spi == null) {
            return;
        }

        MarketData.Response response = msg.getResponse();
        TLMdRspInfoField info = new TLMdRspInfoField();
        info.setSuccess(response.getSuccess());
        if (!response.getSuccess()) {
            info.setErrorMsg(response.getInfo());
            spi.onRspQryMinuteIndicator(null, info, (int) response.getRequestId(), true);
            return;
        }

        MarketData.MinuteIndicatorList tickIndicatorList = response.getMinuteIndicatorList();
        if (tickIndicatorList != null) {
            int count = tickIndicatorList.getMinuteIndicatorsCount();
            int row = 0;
            for (MarketData.MinuteIndicator minuteIndicator : tickIndicatorList.getMinuteIndicatorsList()) {
                row++;
                TLMdMinuteIndicatorField target = new TLMdMinuteIndicatorField();
                target.setExchangeID(minuteIndicator.getExchangeID());
                target.setInstrumentID(minuteIndicator.getInstrumentID());
                target.setMinuteTime(minuteIndicator.getTime());
                target.setIndicators(minuteIndicator.getIndicatorsMap());
                boolean isLast = row == count ? true : false;
                spi.onRspQryMinuteIndicator(target, info, (int) response.getRequestId(), isLast);
            }
        }
    }

    private void handleQrySecuInstrumentsResponse(MarketData.Message msg) {
        if (spi == null) {
            return;
        }

        MarketData.Response response = msg.getResponse();
        TLMdRspInfoField info = new TLMdRspInfoField();
        info.setSuccess(response.getSuccess());
        if (!response.getSuccess()) {
            info.setErrorMsg(response.getInfo());
            spi.onRspQrySecuInstruments(null, info, (int) response.getRequestId(), true);
            return;
        }

        MarketData.SecuInstrumentList instrumentList = response.getSecuInstrumentList();
        if (instrumentList != null) {
            int count = instrumentList.getSecuInstrumentsCount();
            int row = 0;
            for (MarketData.SecuInstrument instrument : instrumentList.getSecuInstrumentsList()) {
                row++;
                TLMdSecuInstrumentField target = new TLMdSecuInstrumentField();
                target.setExchangeID(instrument.getExchangeID());
                target.setInstrumentID(instrument.getInstrumentID());
                target.setInstrumentName(instrument.getInstrumentName().toString(Constants.charset));
                target.setHighLimited(instrument.getHighLimited());
                target.setLowLimited(instrument.getLowLimited());
                target.setPreClosePrice(instrument.getPreClosePrice());
                target.setOpenPrice(instrument.getOpenPrice());
                target.setSecurityType(instrument.getStkType());
                boolean isLast = row == count ? true : false;
                spi.onRspQrySecuInstruments(target, info, (int) response.getRequestId(), isLast);
            }
        }
    }

    private void handleQryFutureInstrumentsResponse(MarketData.Message msg) {
        if (spi == null) {
            return;
        }

        MarketData.Response response = msg.getResponse();
        TLMdRspInfoField info = new TLMdRspInfoField();
        info.setSuccess(response.getSuccess());
        if (!response.getSuccess()) {
            info.setErrorMsg(response.getInfo());
            spi.onRspQryFutureInstruments(null, info, (int) response.getRequestId(), true);
            return;
        }

        MarketData.FutureInstrumentList instrumentList = response.getFutureInstrumentList();
        if (instrumentList != null) {
            int count = instrumentList.getFutureInstrumentsCount();
            int row = 0;
            for (MarketData.FutureInstrument instrument : instrumentList.getFutureInstrumentsList()) {
                row++;
                TLMdFutureInstrumentField target = new TLMdFutureInstrumentField();
                target.setExchangeID(instrument.getExchangeID());
                target.setInstrumentID(instrument.getInstrumentID());
                target.setInstrumentName(instrument.getInstrumentName().toString(Constants.charset));
                target.setHighLimited(instrument.getHighLimited());
                target.setLowLimited(instrument.getLowLimited());
                target.setPreClosePrice(instrument.getPreClosePrice());
                target.setOpenPrice(instrument.getOpenPrice());
                target.setSecurityType(instrument.getStkType());
                target.setPreSettlePrice(instrument.getPreSettlePrice());
                boolean isLast = row == count ? true : false;
                spi.onRspQryFutureInstruments(target, info, (int) response.getRequestId(), isLast);
            }
        }
    }

}
