package main

import (
	"sync"

	"github.com/gorilla/websocket"
)
type Hub struct{
	clients map[string]*websocket.Conn
	mu sync.Mutex
}

func (h *Hub) Register(username string, con *websocket.Conn){
        h.mu.Lock()
		h.clients[username]=con
		h.mu.Unlock()
}

func NewHub() *Hub{

	clients:=make(map[string]*websocket.Conn)
	hub:=&Hub{
		clients: clients,
	}
	return hub
}

func (h *Hub)Unregister(username string){
       h.mu.Lock()
	   delete(h.clients,username)
	   h.mu.Unlock()
}