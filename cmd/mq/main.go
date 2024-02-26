package main

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type wsMessage struct {
	Message string `json:"message"`
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"message", // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
	}

	for d := range msgs {
		wsMess := wsMessage{}

		_ = json.Unmarshal(d.Body, &wsMess)

		fmt.Println(wsMess)
	}
}

// type wsMessage1 struct {
// 	MessageQueue string  `json:"messageQueue"` // имя канала (telegram, whatsApp, service)
// 	UserId       string  `json:"userId"`       // uuid пользователя
// 	ChatId       string  `json:"chatId"`       // id чата
// 	MessageType  string  `json:"MessageType"`  // тип данных сообщения (text/ogg/jpg)
// 	Message      []*byte `json:"data"`         // само сообщение
// }
