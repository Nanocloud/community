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

func rpcRequest(module string, action string, args interface{}) (interface{}, error) {
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
	conn, err := amqp.Dial(conf.QueueUri)

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
