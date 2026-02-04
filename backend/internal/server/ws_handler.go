package server

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"

)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {return true},
}

func serveWS (hub *Hub, w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w,r,nil)
	if err != nil {
		log.Printf("[ws] upgrade error: %v", err)
		return
	}

	client := NewCLient(hub,conn)
	hub.Register() <- client

	go client.WritePump()
	go client.ReadPump()
}