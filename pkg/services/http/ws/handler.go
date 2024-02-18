package ws

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	httpserverhelper "github.com/qweSunset/WS-RMQ/pkg/helpers/httpServerHelper"
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

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		httpserverhelper.ReturnErr(w, err, err.Error())
	}
	defer conn.Close()

	for {
		mtype, message, err := conn.ReadMessage()
		if err != nil || mtype == websocket.CloseMessage {
			fmt.Println("error read ws: "+err.Error(), "-- mTypr: ", uint(mtype))
			break
		}

		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Println(err)
			break
		}

		wsMess := wsMessage{}

		_ = json.Unmarshal(message, &wsMess)

		fmt.Println(wsMess)
	}
}
