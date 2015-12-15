package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"
)

var ch *amqp.Channel

type rpcHandler func(map[string]interface{}) (int, []byte, error)

var kHandler rpcHandler

func rpcListen(dbURI string, handler rpcHandler) error {
	log.Debug("[RPC] listenning")

	conn, err := amqp.Dial(dbURI)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	kHandler = handler

	ch, err = conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		"users", // name
		false,   // durable
		false,   // delete when usused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	if err != nil {
		return err
	}

	log.Debugf("[Users] [RPC] AMQP Queue created: %s\n", q.Name)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	for d := range msgs {
		go handleReq(d)
	}
	return nil
}

func replyError(d *amqp.Delivery) {
	body := make(map[string]interface{})
	body["error"] = "Internal Server Error"

	m := make(map[string]interface{})
	m["body"] = body
	m["http_status_code"] = 500

	res, err := json.Marshal(m)
	if err != nil {
		log.Error(err)
		return
	}

	err = ch.Publish(
		"",        // exchange
		d.ReplyTo, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: d.CorrelationId,
			Body:          res,
		})
	if err != nil {
		log.Error(err)
		return
	}
}

func handleReq(d amqp.Delivery) {
	if d.ContentType != "application/json" {
		log.Errorf("Invalid RPC Request Content-Type: %s", d.ContentType)
		replyError(&d)
		return
	}

	req := make(map[string]interface{})
	err := json.Unmarshal(d.Body, &req)
	if err != nil {
		log.Errorf("Invalid RPC Request Body: %s", err)
		replyError(&d)
		return
	}

	statusCode, body, err := kHandler(req)
	log.Debugf("Status code = %d\n", statusCode)

	if err != nil {
		log.Error(err)
		return
	}

	err = ch.Publish(
		"",        // exchange
		d.ReplyTo, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: d.CorrelationId,
			Body:          body,
		})

	if err != nil {
		log.Error("Failed to publish a message: " + err.Error())
		return
	}

	d.Ack(false)
}
