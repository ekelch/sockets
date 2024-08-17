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
var receiveChannel = make(chan JsonMessage)
var dmChannel = make(chan DirectMessage)

func main() {
	// db
	// dbName := "file:./local.db"
	// db, dbErr := sql.Open("libsql", dbName)
	// if dbErr != nil {
	// 	fmt.Fprintf(os.Stderr, "Failed to open db: %s", dbErr)
	// 	os.Exit(1)
	// }
	// defer db.Close()
	// db
	http.HandleFunc("/", homePage)
	http.HandleFunc("/count", handleCount)
	http.HandleFunc("/ws", handleConnections)

	go scanBroadcast()
	go scanDm()

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
	fmt.Fprintf(w, "%d", len(clients))
}

func closeClient(conn *websocket.Conn) {
	conn.Close()
	broadMsg := BroadcastMessage{
		Clients: clients,
	}
	broadSer, err := json.Marshal(broadMsg)
	if err != nil {
		fmt.Println(err)
	}

	jsonMsg := JsonMessage{
		Type:   "broadcast",
		RawMsg: broadSer,
	}
	receiveChannel <- jsonMsg
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer closeClient(conn)

	// init connect
	var jsonMsg JsonMessage
	err = conn.ReadJSON(&jsonMsg)
	if err != nil {
		fmt.Println(err)
	}
	initBroadcast, err := decodeBroadcast(jsonMsg.RawMsg)
	clients[initBroadcast.FromUser] = UserClient{Conn: conn, Username: initBroadcast.FromUser}
	broadcastMsg(initBroadcast)

	// continuous
	for {
		err = conn.ReadJSON(&jsonMsg)
		if err != nil { //close client failing on this read json
			fmt.Println(err)
			closeClient(conn)
			return
		}
		if jsonMsg.Type == "broadcast" {
			bm, _ := decodeBroadcast(jsonMsg.RawMsg)
			broadcastMsg(bm)
		} else if jsonMsg.Type == "direct" {
			dm, err := decodeDm(jsonMsg.RawMsg)
			if err != nil {
				fmt.Println("err decoding dm")
				return
			}
			dmChannel <- dm
		} else if jsonMsg.Type == "disconnect" {
			bm, _ := decodeBroadcast(jsonMsg.RawMsg)
			fmt.Println(bm)
			delete(clients, bm.FromUser)
			broadcastMsg(bm)
			return
		} else {
			fmt.Println("type not found, closing conn")
			closeClient(conn)
		}
	}
}

func decodeBroadcast(rawBroadcast json.RawMessage) (BroadcastMessage, error) {
	var broadcast BroadcastMessage
	err := json.Unmarshal(rawBroadcast, &broadcast)
	if err != nil {
		fmt.Println(err)
		delete(clients, broadcast.FromUser)
		return broadcast, err
	}
	return broadcast, nil
}

func decodeDm(rawDirect json.RawMessage) (DirectMessage, error) {
	var dm DirectMessage
	err := json.Unmarshal(rawDirect, &dm)
	if err != nil {
		fmt.Println(err)
		return dm, err
	}
	return dm, nil
}

func broadcastMsg(message BroadcastMessage) {
	message.Clients = clients
	broadSer, err := json.Marshal(message)
	if err != nil {
		fmt.Print(err)
	}
	jsonMsg := JsonMessage{
		Type:   "broadcast",
		RawMsg: broadSer,
	}
	receiveChannel <- jsonMsg
}

func scanBroadcast() {
	for {
		jsonMsg := <-receiveChannel
		for username, client := range clients {
			err := client.Conn.WriteJSON(jsonMsg)
			if err != nil {
				client.Conn.Close()
				delete(clients, username)
			}
		}
	}
}

func scanDm() {
	for {
		dm := <-dmChannel
		dmEncode, _ := json.Marshal(dm)
		jsonMsg := JsonMessage{Type: "direct", RawMsg: dmEncode}
		// toUser
		toClient := clients[dm.ToUser]
		err := toClient.Conn.WriteJSON(jsonMsg)
		if err != nil {
			toClient.Conn.Close()
			delete(clients, dm.ToUser)
		}
		// echo fromUser
		fromClient := clients[dm.FromUser]
		err = fromClient.Conn.WriteJSON(jsonMsg)
		if err != nil {
			fromClient.Conn.Close()
			delete(clients, dm.FromUser)
		}
	}
}

type JsonMessage struct {
	Type   string          `json:"type"`
	RawMsg json.RawMessage `json:"rawMsg"`
}

type BroadcastMessage struct {
	FromUser string                `json:"fromUser"`
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
