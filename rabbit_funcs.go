package main

import (
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func connect(connString string) *amqp.Channel {

	conn, err := amqp.Dial(connString)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	return ch

}

func send(ch *amqp.Channel, queueName string, count int) {
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	for i := 0; i < count; i++ {
		body := "777\r\n"
		body = body + "Rome\r\n"
		body = body + "2017-08-19 12:17:55\r\n"
		body = body + "my description\r\n"
		body = body + "{ \"cust_id\": 1313131, \"month\": 12, \"expenses\": 1313.13 }"

		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		failOnError(err, "Failed to publish a message")
	}
}
