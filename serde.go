package main

import (
	"encoding/json"
	"fmt"
)

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
