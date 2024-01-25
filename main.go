package main

import (
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
			command := Command{}
			if conn.ReadJSON(&command) != nil {
				fmt.Println("failed to unmarshal")
				continue
			}
			fmt.Println(command)

			if conn.WriteJSON(command) != nil {
				fmt.Println("failed to write back")
				return
			}
		}
	})

	http.ListenAndServe(":8079", nil)
}
