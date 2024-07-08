package main

import (
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

const RABBITMQ_URL = "amqp://guest:guest@localhost:5672/"

func main() {
	conn, err := amqp.Dial(RABBITMQ_URL)
	failOnError(err, "Failed to connect RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to register a consumer")

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
		wg.Done()
	}()
	wg.Wait()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
