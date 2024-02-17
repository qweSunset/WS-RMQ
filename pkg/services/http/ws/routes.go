package ws

import httpserverhelper "github.com/qweSunset/WS-RMQ/pkg/helpers/httpServerHelper"

var Routes = httpserverhelper.Routes{
	httpserverhelper.Route{
		Name:        "WebSocketListener",
		Method:      "GET",
		Pattern:     "/ws",
		HandlerFunc: WebSocketListener,
	},
}
