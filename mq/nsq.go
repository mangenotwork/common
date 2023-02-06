package mq

import (
	"fmt"

	goNsq "github.com/nsqio/go-nsq"
)

type Producer struct {
	P *goNsq.Producer
}

func NewProducer(addr string) (producer *Producer, err error) {
	if addr == "" {
		err = fmt.Errorf("[NSQ] init failed：need nsq server addr")
		return
	}
	config := goNsq.NewConfig()
	p, err := goNsq.NewProducer(addr, config)
	if err != nil {
		return
	}
	p.SetLogger(nil, 0)
	producer = &Producer{
		P: p,
	}
	return
}

func (m *Producer) Publish(topic string, data []byte) (err error) {
	if m.P == nil {
		err = fmt.Errorf("[NSQ] init failed:%v", err)
	}
	err = m.P.Publish(topic, data)
	defer m.P.Stop()
	if err != nil {
		return fmt.Errorf("[NSQ] publish error:%v", err)
	}
	return
}

// MessageHandler MessageHandler
type MessageHandler struct {
	msgChan     chan *goNsq.Message
	stop        bool
	nsqServer   string
	Channel     string
	maxInFlight int
}

// NewMessageHandler return new MessageHandler
func NewMessageHandler(nsqServer string, channel string) (mh *MessageHandler, err error) {
	if nsqServer == "" {
		err = fmt.Errorf("[NSQ] need nsq server")
		return
	}
	mh = &MessageHandler{
		msgChan:   make(chan *goNsq.Message, 1024),
		stop:      false,
		nsqServer: nsqServer,
		Channel:   channel,
	}
	return
}

// SetMaxInFlight set nsq consumer MaxInFlight
func (m *MessageHandler) SetMaxInFlight(val int) {
	m.maxInFlight = val
}

// Registry register nsq topic
func (m *MessageHandler) Registry(topic string, ch chan []byte) {
	config := goNsq.NewConfig()
	if m.maxInFlight > 0 {
		config.MaxInFlight = m.maxInFlight
	}
	consumer, err := goNsq.NewConsumer(topic, m.Channel, config)
	if err != nil {
		panic(err)
	}
	consumer.SetLogger(nil, 0)
	consumer.AddHandler(goNsq.HandlerFunc(m.handlerMessage))
	err = consumer.ConnectToNSQLookupd(m.nsqServer)
	if err != nil {
		panic(err)
	}
	m.process(ch)
}

func (m *MessageHandler) process(ch chan<- []byte) {
	m.stop = false
	for {
		select {
		case message := <-m.msgChan:
			ch <- message.Body
			if m.stop {
				close(m.msgChan)
				return
			}
		}
	}
}

// handlerMessage handlerMessage
func (m *MessageHandler) handlerMessage(message *goNsq.Message) error {
	if !m.stop {
		m.msgChan <- message
	}
	message.Finish()
	return nil
}
