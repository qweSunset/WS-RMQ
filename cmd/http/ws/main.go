package main

import (
	"log"
	"net/http"

	httpserverhelper "github.com/qweSunset/WS-RMQ/pkg/helpers/httpServerHelper"
	"github.com/qweSunset/WS-RMQ/pkg/services/http/ws"
	"github.com/rs/cors"
)

func main() {

	router := httpserverhelper.NewRouter(ws.Routes)

	castomCors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "authorization", "Accept", "Content-Type", "Content-Length", "Accept-Encoding", "", "", ""},
		Debug:            true,
		AllowedMethods:   []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
	})

	mainHandler := castomCors.Handler(router)

	log.Println("Chat http listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", mainHandler))
}
