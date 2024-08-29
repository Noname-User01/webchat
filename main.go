package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"golang.org/x/net/websocket"
)

var logger *slog.Logger = slog.Default()

type Message struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

type JsonMessage struct {
	Type string `json:"type"`
	Obj  any    `json:"obj"`
}

type Server struct {
	Chats     []Message
	Broadcast map[*websocket.Conn]chan<- Message
}

func (s *Server) AppendMessage(msg Message) {
	s.Chats = append(s.Chats, msg)
	for _, ch := range s.Broadcast {
		ch <- msg
	}
}

// websoket
func (s *Server) WebsocketHandler(ws *websocket.Conn) {
	defer ws.Close()
	logger.Info(fmt.Sprintf("Websocket connected from %s", ws.RemoteAddr().String()))

	ch := make(chan Message, 10)
	s.Broadcast[ws] = ch
	go func() {
		for {
			msg := <-ch
			if msg.Name == "" {
				break
			}
			err := websocket.JSON.Send(ws, JsonMessage{
				Type: "append",
				Obj:  msg,
			})
			if err != nil {
				logger.Warn(fmt.Sprint(err))
			}
		}
	}()

	for {
		message := JsonMessage{}
		err := websocket.JSON.Receive(ws, &message)
		if err != nil {
			if err.Error() == "EOF" {
				logger.Info("Disconnect")
				break
			}
			logger.Warn(fmt.Sprint(err))
			continue

		}
		logger.Info(fmt.Sprint(message))
		switch message.Type {
		case "get":
			err := websocket.JSON.Send(ws, JsonMessage{
				Type: "messages",
				Obj:  s.Chats,
			})
			if err != nil {
				logger.Warn(fmt.Sprint(err))
			}
		case "send":
			mse, ok := message.Obj.(map[string]any)
			if !ok {
				logger.Warn("invalid message")
				break
			}
			msg := Message{
				Name: mse["name"].(string),
				Text: mse["text"].(string),
			}
			s.AppendMessage(msg)
		}

	}
	logger.Info("WebSocket disconnected")
	ch <- Message{}
	delete(s.Broadcast, ws)
}

func main() {
	server := &Server{
		Chats:     make([]Message, 0),
		Broadcast: make(map[*websocket.Conn]chan<- Message),
	}
	files := http.FileServer(http.Dir("public"))

	http.Handle("/", files)
	http.Handle("/ws", websocket.Handler(server.WebsocketHandler))

	logger.Info("Starting server...")
	logger.Info("Starting server...http://localhost:8080")

	//サーバー起動
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Error(fmt.Sprint(err))
	}
}
