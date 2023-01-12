package mq

import (
	"fmt"
	"github.com/mangenotwork/common/conf"
	"github.com/mangenotwork/common/log"
	"github.com/nsqio/go-nsq"
	"github.com/streadway/amqp"
)

var (
	nsqServer = conf.Conf.Nsq.Producer // nsqServer
)

// MQInterface 消息队列接口
type MQInterface interface {
	Producer(topic string, data []byte)
	Consumer(topic, channel string, ch chan []byte, f func(b []byte))
}

// NewMQ 实例化消息队列对象
func NewMQ() MQInterface {
	switch conf.Conf.Mq { // mq 设置的类型
	case "nsq":
		return new(MQNsqService)
	case "rabbit":
		return new(MQRabbitService)
	default:
		return new(MQNsqService)
	}
}

// MQNsqService NSQ消息队列
type MQNsqService struct {
}

// MQRabbitService Rabbit消息队列
type MQRabbitService struct {
}

// Producer 生产者
func (m *MQNsqService) Producer(topic string, data []byte) {
	nsqConf := &nsq.Config{}
	client, err := nsq.NewProducer(conf.Conf.Nsq.Producer, nsqConf)
	if err != nil {
		log.Error("[nsq]无法连接到队列")
		return
	}
	log.DebugF(fmt.Sprintf("[生产消息] topic : %s -->  %s", topic, string(data)))
	err = client.Publish(topic, data)
	if err != nil {
		log.Error("[生产消息] 失败 ： " + err.Error())
	}
}

// Producer 生产者
func (m *MQRabbitService) Producer(topic string, data []byte) {
	mq, err := NewRabbitMQPubSub(topic)
	if err != nil {
		log.Error("[rabbit]无法连接到队列")
		return
	}
	//defer mq.Destroy()

	log.DebugF(fmt.Sprintf("[生产消息] topic : %s -->  %s", topic, string(data)))
	err = mq.PublishPub(data)
	if err != nil {
		log.Error("[生产消息] 失败 ： " + err.Error())
	}
}

// Consumer 消费者
func (m *MQNsqService) Consumer(topic, channel string, ch chan []byte, f func(b []byte)) {
	mh, err := NewMessageHandler(nsqServer, channel)
	if err != nil {
		log.Error(err)
		return
	}

	go func() {
		mh.SetMaxInFlight(1000)
		mh.Registry(topic, ch)
	}()

	go func() {
		for {
			select {
			case s := <-ch:
				f(s)
			}
		}
	}()

	log.DebugF("[NSQ] ServerID:%v => %v started", channel, topic)
}

// Consumer 消费者
func (m *MQRabbitService) Consumer(topic, serverId string, ch chan []byte, f func(b []byte)) {
	mh, err := NewRabbitMQPubSub(topic)
	if err != nil {
		log.Error("[rabbit]无法连接到队列")
		return
	}
	msg := mh.RegistryReceiveSub()

	go func(m <-chan amqp.Delivery) {
		for {
			select {
			case s := <-m:
				f(s.Body)
			}
		}
	}(msg)

	log.DebugF("[Rabbit] ServerID:%v => %v started", serverId, topic)
}
