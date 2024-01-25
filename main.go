package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Command map[string]interface{}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		for {

			msgType, data, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				continue
			} else if msgType != websocket.TextMessage {
				fmt.Println("unexpected message type")
				continue
			}
			command := Command{}
			if json.Unmarshal(data, &command) != nil {
				return
			}
			fmt.Println(command)

			if conn.WriteJSON(command) != nil {
				return
			}
		}
	})

	http.ListenAndServe(":8079", nil)
}
