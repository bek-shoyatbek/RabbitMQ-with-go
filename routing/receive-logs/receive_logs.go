package main

import (
	"log"
	"os"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

const RABBIT_MQ_URL = "amqp://guest:guest@localhost:5672/"

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial(RABBIT_MQ_URL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_direct",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare exchange")

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	failOnError(err, "Failed to declare queue")

	if len(os.Args) < 2 {
		log.Printf("Usage: %s [info] [warning] [error]", os.Args[0])
		os.Exit(0)
	}

	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s", q.Name, "logs_direct", s)
		err = ch.QueueBind(
			q.Name,
			s,
			"logs_direct",
			false,
			nil,
		)
		failOnError(err, "Failed to bind a queue")
	}

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
			log.Printf(" [x] %s", d.Body)
		}
		wg.Done()
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	wg.Wait()
}
