package main

import (
	"sync"

	"github.com/gorilla/websocket"
)
type Client struct{
	conn *websocket.Conn
	send chan []byte
	done chan struct{}
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
		 done: make(chan struct{}),
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
		clients:=make([]*Client,0,len(h.clients))
		for _,client:= range h.clients{
			clients=append(clients,client)
		}
		h.mu.Unlock()
		for _,client:=range clients{
            client.send<-message
		}
	}
}

func (c *Client) writePump(){
	for{
		select{
		  case message:=<-c.send:
	      err:=c.conn.WriteMessage(websocket.TextMessage,message)
		  if(err!=nil){
			return
		  }

		  case <-c.done:return
		}
	}
}