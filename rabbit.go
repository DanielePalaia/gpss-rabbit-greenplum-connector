package main

import (
	"log"

	"github.com/streadway/amqp"
)

type rabbitClient struct {
	connString  string
	queueString string
	channel     *amqp.Channel
	queue       amqp.Queue
	size        int
	buffer      []string
	gpssclient  *gpssClient
}

func makeRabbitClient(connString string, queueString string, size int, gpssclient *gpssClient) *rabbitClient {
	client := new(rabbitClient)
	client.connString = connString
	client.queueString = queueString
	client.size = size
	client.buffer = make([]string, size)
	client.gpssclient = gpssclient

	return client
}

func (client *rabbitClient) failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func (client *rabbitClient) connect() (*amqp.Channel, amqp.Queue) {
	conn, err := amqp.Dial(client.connString)
	client.failOnError(err, "Failed to connect to RabbitMQ")
	//defer conn.Close()

	ch, err := conn.Channel()
	client.failOnError(err, "Failed to open a channel")
	//defer ch.Close()
	client.channel = ch

	q, err := ch.QueueDeclare(
		client.queueString, // name
		false,              // durable
		false,              // delete when usused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)

	client.failOnError(err, "Failed to declare a queue")
	client.queue = q

	return ch, q

}

func (client *rabbitClient) consume() {
	msgs, err := client.channel.Consume(
		client.queue.Name, // queue
		"",                // consumer
		true,              // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	client.failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		count := 0
		for d := range msgs {
			//log.Printf("Received a message: %s", d.Body)

			if count >= client.size {
				log.Printf("im writing")
				client.gpssclient.ConnectToGreenplumDatabase()
				client.gpssclient.WriteToGreenplum(client.buffer)
				client.gpssclient.CloseRequest()
				client.gpssclient.DisconnectToGreenplumDatabase()
				count = 0
			}
			client.buffer[count] = string(d.Body)
			count++
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
