package database

import (
	"time"
)

type MessageState string
type ContentType string
type Role string

const (
	// Enum per MessageState
	StateSent     MessageState = "sent"
	StateReceived MessageState = "received"
	StateRead     MessageState = "read"

	// Enum per ContentType
	TypeText  ContentType = "text"
	TypeImage ContentType = "image"
	TypeAudio ContentType = "audio"

	// Enum per il ruolo degli utenti in una chat
	Member Role = "member"
	Admin  Role = "admin"
)

type User struct {
	UserId    string `json:"userId"`
	Name      string `json:"name"`
	UserPhoto string `json:"userPhoto"`
}

type Group struct {
	GroupId          string    `json:"groupId"`
	ChatId           string    `json:"chatId"`
	GroupName        string    `json:"groupName"`
	GroupDescription string    `json:"groupDescription"`
	PhotoUrl         string    `json:"groupPhoto"`
	CreatedAt        time.Time `json:"createdAt"`
	Members          []string  `json:"members"`
}

type Chat struct {
	ChatId      string    `json:"chatId"`
	GroupId     *string   `json:"groupId,omitempty"`
	Title       string    `json:"title,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	ModifiedAt  time.Time `json:"modifiedAt"`
	LastMessage *Message  `json:"lastMessage,omitempty"`
}

type ChatMember struct {
	ChatId   string    `json:"chatId"`
	UserId   string    `json:"userId"`
	JoinedAt time.Time `json:"joinedAt"`
	Role     Role      `json:"role"`
}

type Message struct {
	MessageId   string       `json:"messageId"`
	ChatId      string       `json:"chatId"`
	SenderId    string       `json:"senderId"`
	Text        string       `json:"text"`
	State       MessageState `json:"state"`
	ContentType ContentType  `json:"contentType"`
	PhotoUrl    string       `json:"photoUrl,omitempty"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

type Comment struct {
	CommentId string    `json:"commentId"`
	MessageId string    `json:"messageId"`
	SenderId  string    `json:"senderId"`
	Emoji     string    `json:"emoji"`
	CreatedAt time.Time `json:"createdAt"`
}

type LoginRequest struct {
	Name string `json:"name"`
}

type LoginResponse struct {
	Identifier string `json:"identifier"`
	Token      string `json:"token"`
}
