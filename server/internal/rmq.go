package internal

import (
	"fmt"

	"github.com/streadway/amqp"
)

// RabbitMQConnection handles connection details
type RabbitMQConnection struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	cQueues map[string]*amqp.Queue
	pQueues map[string]*amqp.Queue
}

// Init initializes the connection
func (rmq *RabbitMQConnection) Init() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq-server:5672/")
	FailOnError(err, "Error connecting to the broker")

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")

	rmq.conn = conn
	rmq.ch = ch

	rmq.cQueues = make(map[string]*amqp.Queue)
	rmq.pQueues = make(map[string]*amqp.Queue)
}

// Kill terminates channel and connection
func (rmq *RabbitMQConnection) Kill() {
	rmq.ch.Close()
	rmq.conn.Close()
}

// CreateCQueue creates a consumer queue
func (rmq *RabbitMQConnection) CreateCQueue(queueName string) {
	q, err := rmq.ch.QueueDeclare(queueName, false, false, false, false, nil)
	FailOnError(err, "Error creating the queue")
	rmq.cQueues[queueName] = &q
}

// CreatePQueue creates a producer queue
func (rmq *RabbitMQConnection) CreatePQueue(queueName string) {
	q, err := rmq.ch.QueueDeclare(queueName, false, false, false, false, nil)
	FailOnError(err, "Error creating the queue")
	rmq.pQueues[queueName] = &q
}

// Consume begins consuming from a consumer queue
func (rmq *RabbitMQConnection) Consume(queueName string) <-chan amqp.Delivery {
	msgs, err := rmq.ch.Consume(queueName, "", false, false, false, false, nil)
	FailOnError(err, "Failed to register as a consumer")

	return msgs
}

// Publish publishes to a producer queue
func (rmq *RabbitMQConnection) Publish(queueName string, message string) {
	err := rmq.ch.Publish("", queueName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	})
	fmt.Println("Published result!")
	FailOnError(err, "Failed to publish a message")
}
