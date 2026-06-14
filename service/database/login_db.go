package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// DoLogin cerca l'utente per nome. Se non c'è, lo crea.
func (db *appdbimpl) DoLogin(username string) (string, error) {
	var id string
	queryCheck := `SELECT identifier FROM users WHERE name = ?`

	err := db.c.QueryRow(queryCheck, username).Scan(&id)

	if err == nil {
		return id, nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("Errore ricerca utente: %w", err)
	}

	newUser, err := db.CreateUser(username)
	if err != nil {
		return "", fmt.Errorf("Errore creazione utente: %w", err)
	}

	return newUser.UserId, nil
}
