package com.tailai.intradata.client.netcomm.amq;
import com.googlecode.protobuf.format.JsonFormat;
import com.tailai.intradata.client.netcomm.MyRunnable;
import com.tailai.intradata.client.netcomm.RPCClientMessageHandler;
import com.tailai.intradata.client.netcomm.proto.MarketData;
import com.tailai.intradata.client.util.MQUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.jms.Message;
import javax.jms.MessageListener;

/**
 * @author yilj
 */
public class ActivemqMessageListener implements MessageListener {
    private static Logger logger = LoggerFactory.getLogger(ActivemqMessageListener.class);
    private final RPCClientMessageHandler handler;
    private final MessageDecoderTask decoder;

    public ActivemqMessageListener(RPCClientMessageHandler handler) {
        this.handler = handler;
        decoder = new MessageDecoderTask("ActivemqMessageDecoderTask", handler);
        startThread(decoder);
    }

    void startThread(MyRunnable<Message> r) {
        Thread t = new Thread(r);
        t.setDaemon(true);
        t.start();
    }

    @Override
    public void onMessage(Message message) {
        try {
            if (!decoder.offer(message)) {
                logger.error("ignored for slow handling, queue size is {}. {}", decoder.getCount(), message);
            }
        } catch (InterruptedException e) {
        }
    }

    public void close() {
        decoder.close();
    }

    private class MessageDecoderTask extends MyRunnable<Message> {
        private final RPCClientMessageHandler handler;
        private JsonFormat jsonFormat = new JsonFormat();

        public MessageDecoderTask(String name, RPCClientMessageHandler handler) {
            super(name);
            this.handler = handler;
        }

        @Override
        protected void doIt(Message message) {
            try {
                byte[] byteArray = MQUtils.getBytes(message);
                if (byteArray == null) {
                    return;
                }

                MarketData.Message protoMsg = MarketData.Message.parseFrom(byteArray);
                if (logger.isTraceEnabled()) {
                    logger.trace("receive message : {}", jsonFormat.printToString(protoMsg));
                }

                handler.onProtoMessage(protoMsg);
            } catch (Throwable e) {
                logger.error("DecoderTask.doIt error", e);
            }
        }
    }
}



