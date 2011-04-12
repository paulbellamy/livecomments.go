package main

import (
	"json"
	"github.com/madari/go-socket.io"
	"github.com/hoisie/redis.go"
	"http"
	"log"
)

var client redis.Client

var sio *socketio.SocketIO

func socketIOConnectHandler(c *socketio.Conn) {
}

func socketIODisconnectHandler(c *socketio.Conn) {
}

func socketIOMessageHandler(c *socketio.Conn, msg socketio.Message) {
	var j map[string]string
	err := json.Unmarshal([]uint8(msg.Data()), &j)

	if err != nil {
		log.Println("Error Parsing Message: ", err)
		return
	}

	switch j["event"] {
	case "initial":
		log.Println("got initial: ", j["data"])
		r, _ := json.Marshal(PaginateFor(j["data"], 0, 10))
		c.Send("{\"event\":\"initial\", \"data\":" + string(r) + "}")
	case "comment":
		log.Println("got comment: ", j["data"])
		if comment, err := Create([]uint8(j["data"])); err == nil {
			log.Println("Stored Comment: ", comment.ToJson())
			sio.Broadcast("{\"event\":\"comment\", \"data\":" + comment.ToJson() + "}")
		} else {
			log.Println("Error Storing Comment: ", err)
		}
	}
}

func main() {
	client.Addr = "localhost:6379"

	// create the socket.io server and mux it to /socket.io/
	config := socketio.DefaultConfig
	config.Origins = []string{"*"}
	sio = socketio.NewSocketIO(&config)

	sio.OnConnect(socketIOConnectHandler)
	sio.OnDisconnect(socketIODisconnectHandler)
	sio.OnMessage(socketIOMessageHandler)

	mux := sio.ServeMux()
	mux.Handle("/", http.FileServer("static/", "/"))
	if err := http.ListenAndServe(":80", mux); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
