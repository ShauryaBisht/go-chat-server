package main

import (
	"sync"

	"github.com/gorilla/websocket"
)
type Client struct{
	conn *websocket.Conn
	send chan []byte
}
type Hub struct{
	clients map[string]*Client
	mu sync.Mutex
	broadcast chan []byte
}

func (h *Hub) Register(username string, con *websocket.Conn) *Client{
	   client:=&Client{
		 conn:con,
		 send: make(chan []byte),
	   }
        h.mu.Lock()
		h.clients[username]=client
		h.mu.Unlock()
		return client
}

func NewHub() *Hub{
    ch:=make(chan []byte)
	clients:=make(map[string]*Client)
	hub:=&Hub{
		clients: clients,
		broadcast: ch,
	}
	return hub
}

func (h *Hub)Unregister(username string){
       h.mu.Lock()
	   delete(h.clients,username)
	   h.mu.Unlock()
}

func (h* Hub) Run(){
	for{
		message:=<-h.broadcast
		h.mu.Lock()
		for _,client:= range h.clients{
			client.send<-message
		}
		h.mu.Unlock()
	}
}

func (c *Client) writePump(){
	for{
	message:=<-c.send
	c.conn.WriteMessage(websocket.TextMessage,message)
	}
}