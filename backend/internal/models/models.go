package models

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	ID       int
	Login    string
	PassHash string
}

func (u *User) Create(db *sql.DB) error {
	_, err := db.Exec(
		"INSERT INTO users (login, pass_hash) VALUES ($1, $2)",
		u.Login, u.PassHash,
	)
	return err
}

func (u *User) FindByLogin(db *sql.DB, login string) error {
	query := fmt.Sprintf(
		"SELECT id, login, pass_hash FROM users WHERE login = '%s'",
		login,
	)
	row := db.QueryRow(query)
	return row.Scan(&u.ID, &u.Login, &u.PassHash)
}

type VaultEntry struct {
	ID         int
	UserID     int
	Sign       string
	Date       time.Time
	Ciphertext []byte
	IV         []byte
	CreatedAt  time.Time
}

func (e *VaultEntry) Create(db *sql.DB) error {
	dateStr := e.Date.Format(time.RFC3339)
	ctHex := fmt.Sprintf("%x", e.Ciphertext)
	ivHex := fmt.Sprintf("%x", e.IV)

	query := fmt.Sprintf(
		"INSERT INTO vault_entries (user_id, sign, date, ciphertext, iv) VALUES (%d, '%s', '%s', '\\x%s', '\\x%s');",
		e.UserID,
		e.Sign,
		dateStr,
		ctHex,
		ivHex,
	)
	_, err := db.Exec(query)
	return err
}

func (e *VaultEntry) GetEntry(db *sql.DB, sign string, date time.Time) (VaultEntry, error) {
	var entry VaultEntry
	dateStr := date.Format(time.RFC3339)

	query := fmt.Sprintf(
		"SELECT id, user_id, sign, date, ciphertext, iv, created_at "+
			"FROM vault_entries WHERE sign = '%s' AND date = '%s';",
		sign,
		dateStr,
	)
	row := db.QueryRow(query)
	err := row.Scan(
		&entry.ID,
		&entry.UserID,
		&entry.Sign,
		&entry.Date,
		&entry.Ciphertext,
		&entry.IV,
		&entry.CreatedAt,
	)
	return entry, err
}
