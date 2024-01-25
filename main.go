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
		fmt.Println("failed to read command json")
		conn.Close()
		return nil
	} else if command["command"] != "AUTH" || command["token"] == nil {
		fmt.Println("command malformed")
		conn.Close()
		return nil
	}
	token := command["token"].(string)
	fmt.Println(token)
	user := verifyToken(token, redis)
	fmt.Println(user)
	if user == nil {
		conn.Close()
		return nil
	}
	conn.WriteJSON(user)
	return user
}

func reader(conn *websocket.Conn, clientChan chan string) {
	for {
		command := Command{}
		if conn.ReadJSON(&command) != nil {
			fmt.Println("failed to unmarshal")
			clientChan <- "hey"
			return
		}
		fmt.Println(command)
	}
}

type Message struct {
	Message string `json:"message"`
	From    string `json:"from"`
	At      string `json:"at"`
}

type Server struct {
	clients  map[*websocket.Conn]bool
	messages []Message
	users    map[int64]*User
}

func (s *Server) broadcast(command string, data map[string]interface{}) {
	payload := make(map[string]interface{})
	for k, v := range data {
		payload[k] = v
	}
	payload["command"] = command
	for client, _ := range s.clients {
		fmt.Println("sending data to client", payload)
		err := client.WriteJSON(payload)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (s *Server) addUser(user *User) {
	if s.users[user.Id] != nil {
		return
	}
	s.users[user.Id] = user
	s.broadcast("USER_JOIN", map[string]interface{}{"user": user})
}

func (s *Server) removeUser(user *User) {
	if s.users[user.Id] != nil {
		delete(s.users, user.Id)
	}
	s.broadcast("USER_LEAVE", map[string]interface{}{"user": user})
}

func main() {
	server := Server{make(map[*websocket.Conn]bool, 0), make([]Message, 0), make(map[int64]*User)}
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
		server.clients[conn] = true
		server.addUser(user)
		clientChan := make(chan string)
		go reader(conn, clientChan)
		//
		for {
			select {
			case <-clientChan:
				server.removeUser(user)
				delete(server.clients, conn)
				return
			}
		}
	})

	http.ListenAndServe(":8079", nil)
}
