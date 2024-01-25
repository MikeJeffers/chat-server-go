package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Command map[string]interface{}

func auth(conn *websocket.Conn, redis *redis.Client) *User {
	command := Command{}
	if conn.ReadJSON(&command) != nil {
		conn.Close()
		return nil
	} else if command["command"] != "AUTH" || command["token"] == nil {
		conn.Close()
		return nil
	}
	token := command["token"].(string)
	user := verifyToken(token, redis)
	if user == nil {
		conn.Close()
		return nil
	}
	conn.WriteJSON(user)
	return user
}

func main() {
	redis := redisClient()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("new connection")
		user := auth(conn, redis)
		if user == nil {
			fmt.Println("Failed auth")
			conn.Close()
			return
		}
		for {
			command := Command{}
			if conn.ReadJSON(&command) != nil {
				fmt.Println("failed to unmarshal")
				return
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
