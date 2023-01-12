package mq

import (
	"fmt"
	"github.com/mangenotwork/common/log"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	//连接
	conn *amqp.Connection
	//管道
	channel *amqp.Channel
	//队列名称
	QueueName string
	//交换机
	Exchange string
	//key Simple模式 几乎用不到
	Key string
	//连接信息
	MqUrl string
}

// NewRabbitMQ 创建RabbitMQ结构体实例
func NewRabbitMQ(queueName, exchange, key, amqpUrl string) (*RabbitMQ, error) {
	mq := &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, MqUrl: amqpUrl}
	var err error
	//创建rabbitMq连接
	mq.conn, err = amqp.Dial(mq.MqUrl)
	if err != nil {
		return nil, err
	}
	mq.failOnErr(err, "创建连接错误！")
	mq.channel, err = mq.conn.Channel()
	mq.failOnErr(err, "获取channel失败")
	return mq, err
}

// Destroy 断开channel和connection
func (r *RabbitMQ) Destroy() {
	r.channel.Close()
	r.conn.Close()
}

// failOnErr 错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.ErrorF(fmt.Sprintf("%s:%s", message, err))
	}
}

// NewRabbitMQSimple 简单模式step：1。创建简单模式下RabbitMQ实例
func NewRabbitMQSimple(queueName string) (*RabbitMQ, error) {
	return NewRabbitMQ(queueName, "", "", "")
}

// PublishSimple 简单模式Step:2、简单模式下生产代码
func (r *RabbitMQ) PublishSimple(message []byte) (err error) {
	//1、申请队列，如果队列存在就跳过，不存在创建
	//优点：保证队列存在，消息能发送到队列中
	_, err = r.channel.QueueDeclare(
		//队列名称
		r.QueueName,
		//是否持久化
		false,
		//是否为自动删除 当最后一个消费者断开连接之后，是否把消息从队列中删除
		false,
		//是否具有排他性 true表示自己可见 其他用户不能访问
		false,
		//是否阻塞 true表示要等待服务器的响应
		false,
		//额外数学系
		nil,
	)

	r.failOnErr(err, "failed to declare a queue")

	//2.发送消息到队列中
	err = r.channel.Publish(
		//默认的Exchange交换机是default,类型是direct直接类型
		r.Exchange,
		//要赋值的队列名称
		r.QueueName,
		//如果为true，根据exchange类型和rout key规则，如果无法找到符合条件的队列那么会把发送的消息返回给发送者
		false,
		//如果为true,当exchange发送消息到队列后发现队列上没有绑定消费者，则会把消息还给发送者
		false,
		//消息
		amqp.Publishing{
			//类型
			ContentType: "text/plain",
			//消息
			Body: []byte(message),
		})
	r.failOnErr(err, "publish 消息失败")
	return
}

// RegistryConsumeSimple 简单模式注册消费者
func (r *RabbitMQ) RegistryConsumeSimple() (msg <-chan amqp.Delivery) {
	//1、申请队列，如果队列存在就跳过，不存在创建
	//优点：保证队列存在，消息能发送到队列中
	_, err := r.channel.QueueDeclare(
		//队列名称
		r.QueueName,
		//是否持久化
		false,
		//是否为自动删除 当最后一个消费者断开连接之后，是否把消息从队列中删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		//额外参数
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	//接收消息
	msg, err = r.channel.Consume(
		r.QueueName,
		//用来区分多个消费者
		"",
		//是否自动应答
		true,
		//是否具有排他性
		false,
		//如果设置为true,表示不能同一个connection中发送的消息传递给这个connection中的消费者
		false,
		//队列是否阻塞
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	return
}

// NewRabbitMQPubSub 订阅模式创建 rabbitMq实例  (目前用的fanout模式)
func NewRabbitMQPubSub(exchangeName string) (*RabbitMQ, error) {
	mq, err := NewRabbitMQ("", exchangeName, "", "")
	if mq == nil || err != nil {
		return nil, err
	}

	//获取connection
	mq.conn, err = amqp.Dial(mq.MqUrl)
	mq.failOnErr(err, "failed to connect mq!")
	if mq.conn == nil || err != nil {
		return nil, err
	}

	//获取channel
	mq.channel, err = mq.conn.Channel()
	mq.failOnErr(err, "failed to open a channel!")
	return mq, err
}

// PublishPub 订阅模式生成
func (r *RabbitMQ) PublishPub(message []byte) (err error) {
	//尝试创建交换机，不存在创建
	err = r.channel.ExchangeDeclare(
		//交换机名称
		r.Exchange,
		//交换机类型 广播类型
		amqp.ExchangeFanout,
		//是否持久化
		true,
		//是否自动删除
		false,
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		//是否阻塞 true表示要等待服务器的响应  false 无等待
		false,
		//参数
		nil,
	)
	r.failOnErr(err, "failed to declare an exchange"+"nge")

	//2 发送消息
	err = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			//类型
			ContentType: "text/plain",
			//消息
			Body: message,
		})
	return
}

// RegistryReceiveSub 订阅模式消费端代码
func (r *RabbitMQ) RegistryReceiveSub() (msg <-chan amqp.Delivery) {
	//尝试创建交换机，不存在创建
	err := r.channel.ExchangeDeclare(
		//交换机名称
		r.Exchange,
		//交换机类型 广播类型
		amqp.ExchangeFanout,
		//是否持久化
		true,
		//是否字段删除
		false,
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		//是否阻塞 true表示要等待服务器的响应
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare an exchange")
	//2试探性创建队列，创建队列
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")
	//绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name,
		//在pub/sub模式下，这里的key要为空
		"",
		r.Exchange,
		false,
		nil,
	)
	//消费消息
	msg, err = r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	return
}

// NewRabbitMQTopic 话题模式 创建RabbitMQ实例
func NewRabbitMQTopic(exchange string, routingKey string) (*RabbitMQ, error) {
	mq, _ := NewRabbitMQ("", exchange, routingKey, "")
	var err error
	mq.conn, err = amqp.Dial(mq.MqUrl)
	mq.failOnErr(err, "failed   to connect rabbitMq!")
	mq.channel, err = mq.conn.Channel()
	mq.failOnErr(err, "failed to open a channel")
	return mq, err
}

// PublishTopic 话题模式发送信息
func (r *RabbitMQ) PublishTopic(message []byte) (err error) {
	//尝试创建交换机，不存在创建
	err = r.channel.ExchangeDeclare(
		//交换机名称
		r.Exchange,
		//交换机类型 话题模式
		amqp.ExchangeTopic,
		//是否持久化
		true,
		//是否字段删除
		false,
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		//是否阻塞 true表示要等待服务器的响应
		false,
		nil,
	)
	r.failOnErr(err, "topic failed to declare an exchange")
	//2发送信息
	err = r.channel.Publish(
		r.Exchange,
		//要设置
		r.Key,
		false,
		false,
		amqp.Publishing{
			//类型
			ContentType: "text/plain",
			//消息
			Body: message,
		})
	return
}

// RegistryReceiveTopic 话题模式接收信息
//
//	要注意key
//	其中* 用于匹配一个单词，#用于匹配多个单词（可以是零个）
//	匹配 xx.* 表示匹配xx.hello,但是xx.hello.one需要用xx.#才能匹配到
func (r *RabbitMQ) RegistryReceiveTopic() (msg <-chan amqp.Delivery) {
	//尝试创建交换机，不存在创建
	err := r.channel.ExchangeDeclare(
		//交换机名称
		r.Exchange,
		//交换机类型 话题模式
		amqp.ExchangeTopic,
		//是否持久化
		true,
		//是否字段删除
		false,
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		//是否阻塞 true表示要等待服务器的响应
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare an exchange")
	//2试探性创建队列，创建队列
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")
	//绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name,
		//在pub/sub模式下，这里的key要为空
		r.Key,
		r.Exchange,
		false,
		nil,
	)
	//消费消息
	msg, err = r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	return
}

// NewRabbitMQRouting 路由模式 创建RabbitMQ实例
func NewRabbitMQRouting(exchange string, routingKey string) (*RabbitMQ, error) {
	rabbitMQ, _ := NewRabbitMQ("", exchange, routingKey, "")
	var err error
	rabbitMQ.conn, err = amqp.Dial(rabbitMQ.MqUrl)
	rabbitMQ.failOnErr(err, "failed   to connect rabbitMq!")
	rabbitMQ.channel, err = rabbitMQ.conn.Channel()
	rabbitMQ.failOnErr(err, "failed to open a channel")
	return rabbitMQ, err
}

// PublishRouting 路由模式发送信息
func (r *RabbitMQ) PublishRouting(message []byte) (err error) {
	//尝试创建交换机，不存在创建
	err = r.channel.ExchangeDeclare(
		//交换机名称
		r.Exchange,
		//交换机类型 广播类型
		amqp.ExchangeDirect,
		//是否持久化
		true,
		//是否字段删除
		false,
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		//是否阻塞 true表示要等待服务器的响应
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare an exchange")
	//发送信息
	err = r.channel.Publish(
		r.Exchange,
		//要设置
		r.Key,
		false,
		false,
		amqp.Publishing{
			//类型
			ContentType: "text/plain",
			//消息
			Body: message,
		})
	return
}

// RegistryReceiveRouting 路由模式接收信息
func (r *RabbitMQ) RegistryReceiveRouting() (msg <-chan amqp.Delivery) {
	//尝试创建交换机，不存在创建
	err := r.channel.ExchangeDeclare(
		//交换机名称
		r.Exchange,
		//交换机类型 广播类型
		amqp.ExchangeDirect,
		//是否持久化
		true,
		//是否字段删除
		false,
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		//是否阻塞 true表示要等待服务器的响应
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare an exchange"+"nge")
	// 试探性创建队列，创建队列
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")
	//绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name,
		//在pub/sub模式下，这里的key要为空
		r.Key,
		r.Exchange,
		false,
		nil,
	)
	//消费消息
	msg, err = r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	return
}
