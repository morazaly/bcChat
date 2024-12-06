package service

import (
	"bccChat/internal/models"
	"bccChat/internal/repository"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type ChatServer struct {
	repo      *repository.MessageRepository
	clients   map[*websocket.Conn]bool
	broadcast chan models.Message
	mutex     sync.Mutex
}

func NewChatServer(repo *repository.MessageRepository) *ChatServer {
	return &ChatServer{
		repo:      repo,
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan models.Message),
	}
}

func (s *ChatServer) HandleMessages() {
	for {
		msg := <-s.broadcast

		//mm := new(models.Message(Roomid:))
		/*
			if err := s.repo.SaveMessage(msg); err != nil {
				log.Printf("Error saving message: %v", err)
			}
		*/
		s.mutex.Lock()
		for client := range s.clients {
			if err := client.WriteJSON(msg); err != nil {
				log.Printf("Error writing JSON: %v", err)
				client.Close()
				delete(s.clients, client)
			}
		}
		s.mutex.Unlock()
	}
}

func (s *ChatServer) AddClient(conn *websocket.Conn) {
	s.mutex.Lock()
	s.clients[conn] = true
	fmt.Println(conn.RemoteAddr())
	s.mutex.Unlock()
}

func (s *ChatServer) RemoveClient(conn *websocket.Conn) {
	s.mutex.Lock()
	delete(s.clients, conn)
	s.mutex.Unlock()
	conn.Close()
}

func (s *ChatServer) BroadcastMessage(msg models.Message) {
	s.broadcast <- msg
}

func (s *ChatServer) StoreMessage(msg models.SendMessage, Username string) {
	s.mutex.Lock()
	s.repo.SaveMessage(msg, Username)
	s.mutex.Unlock()
}

func (s *ChatServer) GetHistory(roomID string) []models.Message {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	mess, _ := s.repo.GetHistory(roomID)
	return mess
}
