package server

import (
	"log"
	"time"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9)/10 
	maxMessageSize = 64*1024
)

type Client struct{
	hub *Hub
	conn *websocket.Conn

	send chan[]byte
}

func NewCLient (hub *Hub, conn *websocket.Conn) *Client{
	return &Client{
		hub: hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
}

func (c *Client) ReadPump (){

	defer func(){
		c.hub.Unregister()  <- c
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func (string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()

		if err != nil {
						if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[client] read error: %v", err)
			}
			return
		}

		_ = message
	}

}

func (c *Client) WritePump (){

	ticker := time.NewTicker(pingPeriod)
	defer func (){
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}

		case <- ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage,nil); err != nil {
				return
			}
		}
	}
}