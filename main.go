package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1000,
	WriteBufferSize: 1000,
	CheckOrigin: func(r *http.Request) bool {
		return true
	}, // allows from any origin
}

var clients = make(map[string]UserClient)
var broadcast = make(chan BroadcastMessage)

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/count", handleCount)
	http.HandleFunc("/ws", handleConnections)

	go receiveBroadcast()

	fmt.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Error starting server: " + err.Error())
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page!\n")
}

func handleCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("ACCESS-CONTROL-ALLOW-ORIGIN", "http://localhost:5173")
	fmt.Printf("%d\n", len(clients))
	fmt.Fprintf(w, "%d", len(clients))
}

func closeClient(conn *websocket.Conn) {
	conn.Close()
	msg := BroadcastMessage{
		Clients: clients,
	}
	broadcast <- msg
	fmt.Print(clients)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer closeClient(conn)

	// init connect
	var msg BroadcastMessage
	err = conn.ReadJSON(&msg)
	username := msg.Username
	client := UserClient{Conn: conn, Username: username}
	clients[username] = client

	if err != nil {
		fmt.Println(err)
		delete(clients, client.Username)
		return
	}
	broadcastMsg(&msg)

	// continuous
	for {
		err = conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			delete(clients, client.Username)
			return
		}
		broadcastMsg(&msg)
	}
}

func broadcastMsg(message *BroadcastMessage) {
	message.Clients = clients
	broadcast <- *message
}

func receiveBroadcast() {
	for {
		msg := <-broadcast

		for username, client := range clients {
			err := client.Conn.WriteJSON(msg)
			if err != nil {
				fmt.Println(err)
				client.Conn.Close()
				delete(clients, username)
			}
		}
	}
}

type BroadcastMessage struct {
	Username string                `json:"username"`
	Message  string                `json:"message"`
	Clients  map[string]UserClient `json:"clients"`
}

type DirectMessage struct {
	FromUser string `json:"fromUser"`
	ToUser   string `json:"toUser"`
	Message  string `json:"message"`
}

type UserClient struct {
	Conn     *websocket.Conn `json:"conn"`
	Username string          `json:"username"`
}

type JsonMessage struct {
	Type   string          `json:"type"`
	RawMsg json.RawMessage `json:"rawMsg"`
}
