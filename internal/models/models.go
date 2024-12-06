package models

type Message struct {
	RoomID    string `json:"room_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

type SendMessage struct {
	RoomID  string `json:"room_id"`
	Content string `json:"content"`
}
