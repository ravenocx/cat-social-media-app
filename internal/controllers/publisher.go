package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ravenocx/cat-socialx/internal/models"
)

func PublishToRabbitMQ(cm *models.CatMatch) error {
	urlString := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USERNAME"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
	)

	conn, err := amqp.Dial(urlString)

	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		"cat_matches", // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	if err != nil {
		return err
	}

	qu, err := ch.QueueDeclare(
		"log", // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	if err != nil {
		return err
	}


	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	body, err := json.Marshal(cm)

	if err != nil {
		return err
	}

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType: "application/json",
			Body:        body,
	})

	if err != nil {
		return err
	}

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		qu.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType: "application/json",
			Body:        body,
	})
	if err != nil {
		return err
	}

	log.Printf(" [x] Sent to rabbitmq :  %s\n", body)
	return nil
}