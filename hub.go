package main

import (
	"encoding/json"
	"sync"
    "errors"
	"github.com/gorilla/websocket"
)
type Client struct{
	conn *websocket.Conn
	done chan struct{}
	room string
	send chan OutgoingMessage
}
type Message struct {
    sender  *Client
    message []byte
	typeof string
}
type OutgoingMessage struct {
    Message string `json:"message"`
    Type    string `json:"type"`
}

type Hub struct{
	clients map[string]*Client
	mu sync.Mutex
    broadcast chan Message
}

func (h *Hub) Register(username string, con *websocket.Conn,room string) (*Client,error){
	   client:=&Client{
		 conn:con,
		 send: make(chan OutgoingMessage),
		 done: make(chan struct{}),
		 room:room,
	   }
        h.mu.Lock()
		if _,exists:=h.clients[username];exists{
			h.mu.Unlock()
			return nil,errors.New("Username already taken")
		}
		h.clients[username]=client
		h.mu.Unlock()
		return client,nil
}

func NewHub() *Hub{
    ch:=make(chan Message)
	clients:=make(map[string]*Client)
	hub:=&Hub{
		clients: clients,
		broadcast: ch,
	}
	return hub
}

func (h *Hub)Unregister(username string){
       h.mu.Lock()
	   client:=h.clients[username]
	   delete(h.clients,username)
	    h.mu.Unlock()
	   h.broadcast<-Message{
		 sender: client,
		 message:[]byte(username+" left "+client.room),
		 typeof: "system",
	   }
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
			if(client.room==message.sender.room){
            client.send<-OutgoingMessage{
				Message: string(message.message),
				Type: message.typeof,
			}
			}
		}
	}
}

func (c *Client) writePump(){
	defer c.conn.Close()
	for{
		select{
		  case message:=<-c.send:
	      data,err:=json.Marshal(message)
		  if(err!=nil){
			return
		  }
          err=c.conn.WriteMessage(websocket.TextMessage,data)
		  if err != nil {
           return
          }
		  case <-c.done:return
		}
	}
}