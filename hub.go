package main

import (
	"sync"

	"github.com/gorilla/websocket"
)
type Hub struct{
	clients map[string]*websocket.Conn
	mu sync.Mutex
	broadcast chan []byte
}

func (h *Hub) Register(username string, con *websocket.Conn){
        h.mu.Lock()
		h.clients[username]=con
		h.mu.Unlock()
}

func NewHub() *Hub{
    ch:=make(chan []byte)
	clients:=make(map[string]*websocket.Conn)
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
		for _,conn:= range h.clients{
			conn.WriteMessage(websocket.TextMessage,message)
		}
		h.mu.Unlock()
	}
}