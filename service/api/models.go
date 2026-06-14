package api

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

type LoginRequest struct {
	Name string
}

type LoginResponse struct {
	Identifier string
	Token      string
}

type User struct {
	Userid    string
	Name      string
	UserPhoto string
}

type Group struct {
	GroupId          string
	ChatId           string
	GroupName        string
	GroupDescription string
	PhotoUrl         string
	CreatedAt        time.Time
	Members          []string
}

type Chat struct {
	ChatId      string
	CreatedAt   time.Time
	ModifiedAt  time.Time
	LastMessage *Message
}

type ChatMember struct {
	ChatId   string
	UserId   string
	JoinedAt time.Time
	Role     Role
}

type Message struct {
	MessageId   string
	ChatId      string
	SenderId    string
	Text        string
	State       MessageState
	ContentType ContentType
	PhotoUrl    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Comment struct {
	CommentId string
	MessageId string
	SenderId  string
	Emoji     string
	CreatedAt time.Time
}
