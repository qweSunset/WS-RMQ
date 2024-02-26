package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	httpserverhelper "github.com/qweSunset/WS-RMQ/pkg/helpers/httpServerHelper"
	amqp "github.com/rabbitmq/amqp091-go"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type wsMessage struct {
	Message string `json:"message"`
}

func WebSocketListener(w http.ResponseWriter, r *http.Request) {

	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		httpserverhelper.ReturnErr(w, err, err.Error())
	}
	defer wsConn.Close()

	mqConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		httpserverhelper.ReturnErr(w, err, err.Error())
	}
	defer mqConn.Close()

	mqChan, err := mqConn.Channel()
	if err != nil {
		httpserverhelper.ReturnErr(w, err, err.Error())
	}
	defer mqChan.Close()

	q1, err := mqChan.QueueDeclare(
		"message", // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		httpserverhelper.ReturnErr(w, err, err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for {
		mtype, message, err := wsConn.ReadMessage()
		if err != nil || mtype == websocket.CloseMessage {
			fmt.Println("error read ws: "+err.Error(), "-- mTypr: ", uint(mtype))
			break
		}

		err = mqChan.PublishWithContext(ctx,
			"",      // exchange
			q1.Name, // routing key
			false,   // mandatory
			false,   // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(message),
			})
		if err != nil {
			fmt.Println(err)
			break
		}

		err = wsConn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Println(err)
			break
		}

		wsMess := wsMessage{}

		_ = json.Unmarshal(message, &wsMess)

		fmt.Println(wsMess)
	}
}
