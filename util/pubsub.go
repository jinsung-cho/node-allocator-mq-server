package util

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Publish(b []byte) error {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		return envErr
	}
	id := os.Getenv("MQ_ID")
	passwd := os.Getenv("MQ_PASSWD")
	ip := os.Getenv("MQ_IP")
	port := os.Getenv("MQ_PORT")
	queue := os.Getenv("MQ_RESOURCE_QUE")
	conn, dialErr := amqp.Dial("amqp://" + id + ":" + passwd + "@" + ip + ":" + port)
	if dialErr != nil {
		return dialErr
	}

	defer conn.Close()

	ch, connErr := conn.Channel()
	if connErr != nil {
		return connErr
	}
	defer ch.Close()

	q, declareErr := ch.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if declareErr != nil {
		return declareErr
	}

	pubErr := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		})
	if pubErr != nil {
		return pubErr
	}

	return nil
}
func Subscribe(byteCh chan<- []byte, hash string) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	ip := os.Getenv("MQ_IP")
	port := os.Getenv("MQ_PORT")
	id := os.Getenv("MQ_ID")
	passwd := os.Getenv("MQ_PASSWD")

	conn, err := amqp.Dial("amqp://" + id + ":" + passwd + "@" + ip + ":" + port)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		hash,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	for msg := range msgs {
		fmt.Println("in for loop")
		body := msg.Body
		byteCh <- body

		return
	}
}
