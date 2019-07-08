package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

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
	filePath    string
	fileBatch   *os.File
	mode        int
}

func makeRabbitClient(connString string, queueString string, size int, gpssclient *gpssClient, mode int, filePath string) *rabbitClient {
	client := new(rabbitClient)
	client.connString = connString
	client.queueString = queueString
	client.size = size
	client.buffer = make([]string, size)
	client.gpssclient = gpssclient
	client.mode = mode
	// if persistency is set then open the file
	if client.mode == 1 {
		client.filePath = filePath
		var err error
		client.fileBatch, err = os.OpenFile(client.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			client.failOnError(err, "error opening text file")
		}
	}

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
		true,               // durable
		false,              // delete when usused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)

	client.failOnError(err, "Failed to declare a queue")
	client.queue = q

	return ch, q

}

func (client *rabbitClient) consume(mode bool) {
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
			if client.mode == 1 {
				// persistence of batch activated
				client.writeToFile(string(d.Body))
			}

			client.buffer[count] = string(d.Body)
			count++

			if count == client.size {
				log.Printf("Batch reached: I'm sending request to write to gpss/gprc server")
				client.gpssclient.ConnectToGreenplumDatabase()
				client.gpssclient.WriteToGreenplum(client.buffer)
				if client.mode == 1 {
					// persistence of batch activated
					client.truncateFile()
				}
				client.gpssclient.DisconnectToGreenplumDatabase()
				count = 0
				// Just for testing purpose
				if mode == false {
					close(forever)
				}

			}
		}

	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	<-forever

}

func (client *rabbitClient) writeToFile(message string) {

	//log.Printf("writing to file")
	w := bufio.NewWriter(client.fileBatch)
	message = message + "\n"
	fmt.Fprintf(w, "%v", message)
	w.Flush()

}

func (client *rabbitClient) truncateFile() {
	//log.Printf("truncating the file")
	os.Truncate(client.filePath, 0)

}
