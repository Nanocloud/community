/*
 * Nanocloud Community, a comprehensive platform to turn any application
 * into a cloud solution.
 *
 * Copyright (C) 2015 Nanocloud Software
 *
 * This file is part of Nanocloud community.
 *
 * Nanocloud community is free software; you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * Nanocloud community is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
