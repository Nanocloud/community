package main

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/streamrail/concurrent-map"
	"math/rand"
)

var ch *amqp.Channel

type rpcHandler func(map[string]interface{}) (int, []byte, error)

var kHandler rpcHandler

var kAmqpChannel *amqp.Channel
var kRPCContexts = cmap.New()
var kAmqpQueueName string

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func rpcRequest(module string, action string, args interface{}) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	m["action"] = action
	m["args"] = args

	body, err := json.Marshal(m)
	if err != nil {
		log.Error("Unable to serialize the RPC request")
		return nil, err
	}

	corrId := randomString(32)
	resChan := make(chan *amqp.Delivery)

	_, exists := kRPCContexts.Get(corrId)
	if exists {
		log.Error("Correlation id already exists")
		return nil, errors.New("Correlation id already exists")
	}

	kRPCContexts.Set(corrId, resChan)

	err = kAmqpChannel.Publish(
		"",     // exchange
		module, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrId,
			ReplyTo:       kAmqpQueueName,
			Body:          body,
		})

	if err != nil {
		return nil, err
	}

	delivery := <-resChan

	if delivery.ContentType != "application/json" {
		err = errors.New(fmt.Sprintf("Invalid RPC request reponse type: %s", delivery.ContentType))
		log.Error(err)
		return nil, err
	}

	response := make(map[string]interface{})
	err = json.Unmarshal(delivery.Body, &response)
	if err != nil {
		err = errors.New(fmt.Sprintf("Unable to parse RPC request reponse: %s", err))
		log.Error(err)
		return nil, err
	}

	return response, nil
}

func init() {
	initConf()
	conn, err := amqp.Dial(conf.QueueURI)

	if err != nil {
		log.Fatalf("%s: %s", "Unable to connect to RabbitMQ", err)
	}

	kAmqpChannel, err = conn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "Unable to open RabbitMQ channel", err)
	}

	q, err := kAmqpChannel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
	}

	kAmqpQueueName = q.Name

	msgs, err := kAmqpChannel.Consume(
		kAmqpQueueName, // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to register a consumer", err)
	}

	go func() {
		for d := range msgs {
			corrId := d.CorrelationId
			context, exists := kRPCContexts.Get(corrId)
			if exists {
				kRPCContexts.Remove(corrId)
				context.(chan *amqp.Delivery) <- &d
			} else {
				log.Errorf("Cannot handle RPC request (correlation id = %s)", corrId)
			}
		}
	}()
}

func rpcListen(dbURI string, module string, handler rpcHandler) error {
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
		module, // name
		false,  // durable
		false,  // delete when usused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
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
