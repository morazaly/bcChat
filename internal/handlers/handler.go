package handlers

import (
	"bccChat/internal/models"
	service "bccChat/internal/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	service *service.ChatServer
}

func NewHandler(service *service.ChatServer) *Handler {
	return &Handler{
		service: service}
}

func (h *Handler) StartHandler( /*ctx context.Context,*/ ch chan error) { //ctx надо передать дальше по функциям
	r := mux.NewRouter()

	r.HandleFunc("/chat/connect", HandleConnections(h.service)).Methods("GET")
	r.HandleFunc("/chat/message", SendMessage(h.service)).Methods("POST")
	r.HandleFunc("/chat/history", GetHistory(h.service)).Methods("GET")

	go h.service.HandleMessages()

	fmt.Println("Server started on :8080")
	ch <- http.ListenAndServe(":8080", r)
	/*if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}*/
}

func HandleConnections(server *service.ChatServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Error upgrading connection: %v", err)
			return
		}

		server.AddClient(conn)

		defer server.RemoveClient(conn)

		for {
			var msg models.SendMessage
			if err := conn.ReadJSON(&msg); err != nil {
				log.Printf("Error reading JSON: %v", err)
				break
			}
			server.StoreMessage(msg, conn.RemoteAddr().String())
			t := time.Now()
			server.BroadcastMessage(models.Message{RoomID: msg.RoomID, Username: conn.RemoteAddr().String(), Content: msg.Content, Timestamp: t.String()})
		}
	}
}
func SendMessage(server *service.ChatServer) http.HandlerFunc { //Оставил этот метод гет по условии. там надо сделать идентификацию
	return func(w http.ResponseWriter, r *http.Request) {

		/*conn, err := upgrader.Upgrade(w, r, nil) //coonect
		if err != nil {
			log.Printf("Error upgrading connection: %v", err)
			return
		}*/

		var msg models.SendMessage
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		server.StoreMessage(msg, "Api user") //conn.RemoteAddr().String())

		t := time.Now()
		server.BroadcastMessage(models.Message{RoomID: msg.RoomID, Username: "Api user" /*conn.RemoteAddr().String()*/, Content: msg.Content, Timestamp: t.String()})

		w.WriteHeader(http.StatusOK)
	}
}

func GetHistory(server *service.ChatServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID := r.URL.Query().Get("room_id")
		if roomID == "" {
			http.Error(w, "room_id is required", http.StatusBadRequest)
			return
		}

		history := server.GetHistory(roomID)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(history); err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		}
	}
}
