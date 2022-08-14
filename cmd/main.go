package main

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type ChatMessage struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}
type broadcastMessage struct {
	ChatMessage ChatMessage
	Ws          *websocket.Conn
}

var clients = make(map[*websocket.Conn]bool)
var broadcaster = make(chan broadcastMessage)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// ensure connection close when function returns
	defer ws.Close()
	clients[ws] = true

	for {
		var msg ChatMessage
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(clients, ws)
			break
		}
		broadcastMessage := broadcastMessage{ChatMessage: msg, Ws: ws}
		// send new message to the channel
		broadcaster <- broadcastMessage
	}
}

// If a message is sent while a client is closing, ignore the error
func unsafeError(err error) bool {
	return !websocket.IsCloseError(err, websocket.CloseGoingAway) && err != io.EOF
}

func handleMessages() {
	for {
		msg := <-broadcaster

		messageClients(msg)
	}
}

func messageClients(msg broadcastMessage) {
	// send to every client currently connected except the one who sent the message
	for client := range clients {
		if client == msg.Ws {
			continue
		}
		messageClient(client, msg.ChatMessage)
	}
}

func messageClient(client *websocket.Conn, msg ChatMessage) {
	err := client.WriteJSON(msg)
	if err != nil && unsafeError(err) {
		log.Printf("error: %v", err)
		client.Close()
		delete(clients, client)
	}
}

func main() {

	port := "8080"

	http.HandleFunc("/websocket", handleConnections)
	go handleMessages()

	log.Print("Server starting at localhost:" + port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
