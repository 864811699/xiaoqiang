package com.tailai.intradata.client.netcomm.amq;

import com.google.protobuf.GeneratedMessageV3;
import com.googlecode.protobuf.format.JsonFormat;
import com.tailai.intradata.client.TLDataType;
import com.tailai.intradata.client.TLMdException;
import com.tailai.intradata.client.netcomm.Connector;
import com.tailai.intradata.client.netcomm.RPCClientMessageHandler;
import com.tailai.intradata.client.netcomm.RPCException;
import com.tailai.intradata.client.netcomm.proto.MarketData;
import com.tailai.intradata.client.util.Constants;
import org.apache.activemq.ActiveMQConnectionFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.jms.*;
import java.lang.IllegalStateException;
import java.util.Arrays;
import java.util.UUID;
import java.util.concurrent.atomic.AtomicBoolean;

public class ActivemqConnector implements Connector {
    private static final Logger logger = LoggerFactory.getLogger(ActivemqConnector.class);
    private JsonFormat jsonFormat = new JsonFormat();
    private AtomicBoolean initialized = new AtomicBoolean(false);
    private TLMdApiAmqConfig config;
    private RPCClientMessageHandler messageHandler;
    private ActivemqMessageListener messageListener;
    private Connection connection;
    private Session managingSession;
    private MessageProducer requestProducer;
    private Session subscribeSession;
    private final String clientID;

    public ActivemqConnector(TLMdApiAmqConfig config, RPCClientMessageHandler messageHandler) {
        this.config = config;
        this.messageHandler = messageHandler;
        this.clientID = UUID.randomUUID().toString();
    }

    @Override
    public void init() {
        if (config == null) {
            throw new TLMdException("config is not allowed to be null.");
        }
        if (initialized.getAndSet(true)) {
            throw new IllegalStateException("Api is already initialized. ");
        }

        try {
            logger.info("ActiveMQ connection initializing. ");
            ActiveMQConnectionFactory connectionFactory = new ActiveMQConnectionFactory(config.getUrl());
            connectionFactory.setTransportListener(new ActivemqTransportListener(config.getUrl(), messageHandler));
            connectionFactory.setDispatchAsync(true);
            connectionFactory.setOptimizeAcknowledge(true);
            connectionFactory.setOptimizeAcknowledgeTimeOut(300);
            if (config.getUserName() != null && !config.getUserName().equals("")) {
                connection = connectionFactory.createConnection(config.getUserName(), config.getPassword());
            } else {
                connection = connectionFactory.createConnection();
            }
            connection.start();

            messageListener = new ActivemqMessageListener(messageHandler);

            initQuerySession();

            logger.info("ActiveMQ connection initialized successfully.");
        } catch (Exception e) {
            throw new TLMdException("ActiveMQ manage connection initialized error.", e);
        }
    }

    private void initQuerySession() throws JMSException {
        managingSession = connection.createSession(false, Session.AUTO_ACKNOWLEDGE);

        Queue requestQueue = managingSession.createQueue(Constants.THEME_REQUEST);
        requestProducer = managingSession.createProducer(requestQueue);
        requestProducer.setDeliveryMode(DeliveryMode.NON_PERSISTENT);

        String messageSelector = "ClientID='" + this.clientID + "'";
        Queue responseQueue = managingSession.createQueue(Constants.THEME_RESPONSE);
        MessageConsumer responseConsumer = managingSession.createConsumer(responseQueue, messageSelector);
        responseConsumer.setMessageListener(messageListener);
    }

    @Override
    public void release() {
        this.initialized.set(false);
        if (requestProducer != null) {
            try {
                requestProducer.close();
            } catch (JMSException e) {
                logger.error("requestProducer close error", e);
            }
        }
        if (managingSession != null) {
            try {
                managingSession.close();
            } catch (JMSException e) {
                logger.error("managingSession close error", e);
            }
        }
        if (subscribeSession != null) {
            try {
                subscribeSession.close();
            } catch (JMSException e) {
                logger.error("subscribeSession close error", e);
            }
        }
        if (connection != null) {
            try {
                connection.close();
            } catch (JMSException e) {
                logger.error("connection close error", e);
            }
        }
        if (messageListener != null) {
            messageListener.close();
        }
    }

    @Override
    public void send(GeneratedMessageV3 message) throws RPCException {
        if (logger.isTraceEnabled()) {
            logger.trace("send message: {}", jsonFormat.printToString(message));
        }

        try {
            if (message instanceof MarketData.Message) {
                MarketData.Message protoMsg = (MarketData.Message) message;
                BytesMessage bytesMessage = managingSession.createBytesMessage();
                bytesMessage.setObjectProperty(Constants.MSGPROPERTIES_CLIENTID, clientID);
                bytesMessage.writeBytes(protoMsg.toByteArray());
                requestProducer.send(bytesMessage);
            }
        } catch (Exception e) {
            logger.error("send message error", e);
            throw new RPCException(e);
        }
    }

    @Override
    public void subscribe() throws RPCException {
        try {
            subscribeSession = connection.createSession(false, Session.AUTO_ACKNOWLEDGE);
            StringBuilder selector = new StringBuilder();
            selector.append("MarketId IN (");
            int i = 0;
            for (String marketId : config.getMarkets()) {
                if (i > 0) {
                    selector.append(",");
                }
                selector.append("'");
                selector.append(marketId);
                selector.append("'");
                i++;
            }
            selector.append(")");
            String selectorStr = selector.toString();

            if (config.isSubQuotation()) {
                logger.debug("starting stock market data consumer, selector = {}", selectorStr);
                Topic stockTopic = subscribeSession.createTopic(Constants.THEME_SUBSCRIBE_STOCK_SNAPSHOT);
                MessageConsumer stockConsumer = subscribeSession.createConsumer(stockTopic, selectorStr);
                stockConsumer.setMessageListener(messageListener);
            }

            if (Arrays.asList(config.getMarkets()).contains(TLDataType.EXCHANGE_CFFEX)) {
                logger.debug("starting future market data consumer, selector = {}", selectorStr);
                Topic futureTopic = subscribeSession.createTopic(Constants.THEME_SUBSCRIBE_FUTURE_SNAPSHOT);
                MessageConsumer futureConsumer = subscribeSession.createConsumer(futureTopic, selectorStr);
                futureConsumer.setMessageListener(messageListener);
            }

            if (config.isSubSimpleIndicator()) {
                logger.debug("starting simple indicator consumer, selector = {}", selectorStr);
                Topic indicatorTopic = subscribeSession.createTopic(Constants.THEME_SUBSCRIBE_INDICATOR_TICK);
                MessageConsumer indicatorConsumer = subscribeSession.createConsumer(indicatorTopic, selectorStr);
                indicatorConsumer.setMessageListener(messageListener);
            }

            if (config.isSubMinuteIndicator()) {
                logger.debug("starting minute indicator consumer, selector = {}", selectorStr);
                Topic minuterTopic = subscribeSession.createTopic(Constants.THEME_SUBSCRIBE_INDICATOR_MINUTE);
                MessageConsumer minuteConsumer = subscribeSession.createConsumer(minuterTopic, selectorStr);
                minuteConsumer.setMessageListener(messageListener);
            }

            if (config.isSubMinuteSignal()) {
                logger.debug("starting minute signal consumer, selector = {}", selectorStr);
                Topic minuterTopic = subscribeSession.createTopic(Constants.THEME_SUBSCRIBE_SIGNAL_MINUTE);
                MessageConsumer minuteConsumer = subscribeSession.createConsumer(minuterTopic, selectorStr);
                minuteConsumer.setMessageListener(messageListener);
            }

            if (config.isSubTickSignal()) {
                logger.debug("starting tick signal consumer, selector = {}", selectorStr);
                Topic minuterTopic = subscribeSession.createTopic(Constants.THEME_SUBSCRIBE_SIGNAL_TICK);
                MessageConsumer minuteConsumer = subscribeSession.createConsumer(minuterTopic, selectorStr);
                minuteConsumer.setMessageListener(messageListener);
            }

            if (config.isSubOrderQueue() || config.isSubOrder() || config.isSubTransaction()) {
                selectorStr += " AND " + Constants.MSGPROPERTIES_DATATYPE + " IN (''";
                if (config.isSubOrderQueue()) {
                    selectorStr += ",'" + MarketData.MsgType.Notify_Stock_Order_Queue_VALUE + "'";
                }
                if (config.isSubOrder()) {
                    selectorStr += ",'" + MarketData.MsgType.Notify_Transaction_Order_VALUE + "'";
                }
                if (config.isSubTransaction()) {
                    selectorStr += ",'" + MarketData.MsgType.Notify_Transaction_Trade_VALUE + "'";
                }
                selectorStr += ")";

                logger.debug("starting stock transaction data consumer, selector = {}", selectorStr);
                Topic transactionTopic = subscribeSession.createTopic(Constants.THEME_SUBSCRIBE_STOCK_TRANSACTION);
                MessageConsumer minuteConsumer = subscribeSession.createConsumer(transactionTopic, selectorStr);
                minuteConsumer.setMessageListener(messageListener);
            }
        } catch (Exception e) {
            throw new RPCException(e);
        }
    }
}
