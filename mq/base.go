package mq

import (
	"fmt"
	"sync"

	"github.com/mangenotwork/common/conf"
	"github.com/mangenotwork/common/log"

	"github.com/Shopify/sarama"
	"github.com/nsqio/go-nsq"
	"github.com/streadway/amqp"
)

var (
	nsqServer   = conf.Conf.Nsq.Producer // nsqServer
	kafkaServer = conf.Conf.Kafka.Addr   // kafkaServer
)

// MQer 消息队列接口
type MQer interface {
	Producer(topic string, data []byte)
	Consumer(topic, channel string, ch chan []byte, f func(b []byte))
}

// NewMQ 实例化消息队列对象
func NewMQ() MQer {
	switch conf.Conf.Mq { // mq 设置的类型
	case "nsq":
		return new(MQNsqService)
	case "rabbit":
		return new(MQRabbitService)
	case "kafka":
		return new(MQKafkaService)
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

// MQKafkaService Kafka消息队列
type MQKafkaService struct {
}

// Producer 生产者
func (m *MQNsqService) Producer(topic string, data []byte) {
	nsqConf := &nsq.Config{}
	client, err := nsq.NewProducer(nsqServer, nsqConf)
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

// Producer 生产者
func (m *MQKafkaService) Producer(topic string, data []byte) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follower都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner //写到随机分区中，我们默认设置32个分区
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回
	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.ByteEncoder(data)
	// 连接kafka
	client, err := sarama.NewSyncProducer(kafkaServer, config)
	if err != nil {
		log.Error("Producer closed, err:", err)
		return
	}
	defer client.Close()
	// 发送消息
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		log.Error("send msg failed, err:", err)
		return
	}
	log.InfoF("pid:%v offset:%v\n", pid, offset)
}

// Consumer 消费者
func (m *MQKafkaService) Consumer(topic, serverId string, ch chan []byte, f func(b []byte)) {
	var wg sync.WaitGroup
	consumer, err := sarama.NewConsumer(kafkaServer, nil)
	if err != nil {
		log.ErrorF("Failed to start consumer: %s", err)
		return
	}
	partitionList, err := consumer.Partitions("task-status-data") // 通过topic获取到所有的分区
	if err != nil {
		log.Error("Failed to get the list of partition: ", err)
		return
	}
	log.Info(partitionList)
	for partition := range partitionList { // 遍历所有的分区
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest) // 针对每个分区创建一个分区消费者
		if err != nil {
			log.ErrorF("Failed to start consumer for partition %d: %s\n", partition, err)
		}
		wg.Add(1)
		go func(sarama.PartitionConsumer) { // 为每个分区开一个go协程取值
			for msg := range pc.Messages() { // 阻塞直到有值发送过来，然后再继续等待
				log.DebugF("Partition:%d, Offset:%d, key:%s, value:%s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				f(msg.Value)
			}
			defer pc.AsyncClose()
			wg.Done()
		}(pc)
	}
	wg.Wait()
	consumer.Close()
}
