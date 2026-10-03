package main

import (

	// "io"
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}
func main() {
	hub:=NewHub()
	go hub.Run()
	http.HandleFunc("/", handler)
	http.HandleFunc("/ws",wsHandler(hub))
	http.ListenAndServe(":8080", nil)
}
func handler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func wsHandler(hub *Hub) http.HandlerFunc{
     return func(w http.ResponseWriter,r *http.Request){
	 conn, err := upgrader.Upgrade(w, r, nil)
	 if err!=nil{
	    return
	 }
	 username:=r.URL.Query().Get("username")
	 room:=r.URL.Query().Get("room")
	 client,err:=hub.Register(username,conn,room)
	 if(err!=nil){
		response:=OutgoingMessage{
			Message: err.Error(),
			Type: "error",
		}
        data,_:=json.Marshal(response)
		conn.WriteMessage(websocket.TextMessage,data)
		conn.Close()
		return
	 }
	 hub.broadcast<-Message{
		sender: client,
		message:[]byte(username+ " joined "+room),
		typeof: "system",
	 }
	 go client.writePump()
	 defer conn.Close()
     defer hub.Unregister(username)
     defer close(client.done)
	 for {
	 _,msg,err:=conn.ReadMessage()
	 if err!=nil{
	    return
	 }
	 hub.broadcast<-Message{
		sender: client,
		message: msg,
	 }
	}
  }
}
