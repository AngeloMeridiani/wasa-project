package database

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	// Login
	DoLogin(name string) (string, error)

	// User
	CreateUser(name string) (*User, error)
	SetUserName(userId string, name string) error
	GetUser(id string) (*User, error)
	GetUsersByName(name string) ([]User, error)
	SetMyPhoto(userId string, photoUrl string) error

	// Chats
	CreateChat(userId1 string, userId2 string) (*Chat, error)
	DeleteChat(chatId string, userId string) error
	GetMyChats(userId string) ([]Chat, error)
	GetCurrentChat(chatId string, userId string, limit int) (*Chat, []Message, error)

	// Messages
	SendMessage(senderId string, chatId string, text string) (*Message, error)
	SendImageMessage(senderId string, chatId string, photoUrl string) (*Message, error)
	UpdateMessage(messageId string, text string, senderId string) (*Message, error)
	DeleteMessage(messageId string, userId string) error
	UpdateMessageStatus(messageId string, newStatus MessageState) (time.Time, error)
	ForwardMessage(messageId string, targetChatId string, userId string) (*Message, error)

	// Comments
	CreateComment(senderId string, messageId string, emoji string) (*Comment, error)
	DeleteComment(commentId string, userId string) error

	// Groups
	CreateGroup(groupName string, groupDescription string, photoUrl string, creatorId string, members []string) (*Group, error)
	GetGroup(groupId string) (*Group, error)
	SetGroupName(groupId string, userId string, name string) error
	SetGroupPhoto(groupId string, userId string, photoUrl string) error
	AddToGroup(groupId string, requesterId string, userId string) error
	LeaveGroup(groupId string, userId string) error

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Tabella USERS
	sqlUsers := `
	CREATE TABLE IF NOT EXISTS users (
		identifier TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		photo_url TEXT 
	);`

	// Tabella CHATS
	sqlChats := `
	CREATE TABLE IF NOT EXISTS chats (
		chat_id TEXT PRIMARY KEY,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		modified_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		type TEXT NOT NULL DEFAULT 'dm',
		last_message_id TEXT,
		FOREIGN KEY(last_message_id) REFERENCES messages(message_id) ON DELETE SET NULL
	);`

	// Tabella GROUPS
	sqlGroups := `
		CREATE TABLE IF NOT EXISTS groups (
			group_id TEXT PRIMARY KEY,
			group_name TEXT NOT NULL,
			group_description TEXT,
			photo_url TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			chat_id TEXT NOT NULL,
			FOREIGN KEY(chat_id) REFERENCES chats(chat_id) ON DELETE CASCADE
		);`

	// Tabella MESSAGES
	sqlMessages := `
		CREATE TABLE IF NOT EXISTS messages (
			message_id TEXT PRIMARY KEY,
			chat_id TEXT NOT NULL,           
			sender_id TEXT NOT NULL, 
			text TEXT,
			state TEXT DEFAULT 'sent',
			content_type TEXT DEFAULT 'text',
			photo_url TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME,
			FOREIGN KEY(chat_id) REFERENCES chats(chat_id) ON DELETE CASCADE,
			FOREIGN KEY(sender_id) REFERENCES users(identifier)
		);`

	// Tabella CHAT_MEMBERS
	sqlMembers := `
		CREATE TABLE IF NOT EXISTS chat_members (
			chat_id TEXT,
			user_id TEXT,
			joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			role TEXT DEFAULT 'member',
			PRIMARY KEY (chat_id, user_id),
			FOREIGN KEY(chat_id) REFERENCES chats(chat_id) ON DELETE CASCADE,
			FOREIGN KEY(user_id) REFERENCES users(identifier) ON DELETE CASCADE
		);`

	// Tabella COMMENTS
	sqlComments := `
		CREATE TABLE IF NOT EXISTS comments (
			comment_id TEXT PRIMARY KEY,
			message_id TEXT NOT NULL,
			sender_id TEXT NOT NULL,
			emoji TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME,
			FOREIGN KEY(message_id) REFERENCES messages(message_id) ON DELETE CASCADE,
			FOREIGN KEY(sender_id) REFERENCES users(identifier)
		);`

	// Tabella DELETED_CHATS
	sqlDeletedChats := `
		CREATE TABLE IF NOT EXISTS deleted_chats (
			chat_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			deleted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY(chat_id, user_id),
			FOREIGN KEY(chat_id) REFERENCES chats(chat_id) ON DELETE CASCADE,
			FOREIGN KEY(user_id) REFERENCES users(identifier)
		);`

	queries := []string{sqlUsers, sqlChats, sqlGroups, sqlMessages, sqlMembers, sqlComments, sqlDeletedChats}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return nil, fmt.Errorf("error creating database schema: %w", err)
		}
	}

	return &appdbimpl{c: db}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

// ID univoco di 20 caratteri
func generateID() string {
	id := make([]byte, 10)
	_, _ = rand.Read(id)
	return hex.EncodeToString(id)
}
