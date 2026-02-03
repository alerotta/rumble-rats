package server

import "log"

type Hub struct {
	clients map[*Client]struct{}
	register chan *Client
	unregister chan *Client 
	broadcast chan []byte
}

func NewHub () *Hub {
	return &Hub{
		clients: make(map[*Client]struct{}),
		register: make(chan *Client, 64),
		unregister: make(chan *Client, 64),
		broadcast: make(chan []byte, 256),
	}
}

func (h *Hub) Register() chan<- *Client   { return h.register }
func (h *Hub) Unregister() chan<- *Client { return h.unregister }
func (h *Hub) Broadcast() chan<- []byte   { return h.broadcast }

func (h *Hub) Run() {
	for {
		select {

		case c := <-h.register:
			h.clients[c] = struct{}{}
			
			log.Printf("[hub] client registered (%d total)", len(h.clients))

		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				
				close(c.send)
				log.Printf("[hub] client unregistered (%d total)", len(h.clients))
			}

		case msg := <-h.broadcast:
			
			
			for c := range h.clients {
				select {
				case c.send <- msg:
					
				default:
					
					delete(h.clients, c)
					close(c.send)
					log.Printf("[hub] dropped slow client (%d total)", len(h.clients))
				}
			}
		}
	}
}

func (h *Hub) BroadcastToAll(msg []byte) {
	h.broadcast <- msg
}