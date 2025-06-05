package db

import (
	"database/sql"
	"time"
	_ "github.com/mattn/go-sqlite3"
	"github.com/google/uuid"
)

type Token struct {
	ID          string
	Value       string
	Type        string
	Description string
	Active      bool
	CreatedAt   int64
}

type StatEvent struct {
	ID         int64
	Event      string
	Package    string
	Version    string
	OccurredAt int64
}

type DB struct {
	sql *sql.DB
}

func Open(path string) (*DB, error) {
	sqlite, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return &DB{sql: sqlite}, nil
}

// Token CRUD

func (db *DB) ListTokens() ([]Token, error) {
	rows, err := db.sql.Query(`SELECT id, value, type, description, active, created_at FROM tokens`)
	if err != nil { return nil, err }
	defer rows.Close()
	var tokens []Token
	for rows.Next() {
		var t Token
		var active int
		if err := rows.Scan(&t.ID, &t.Value, &t.Type, &t.Description, &active, &t.CreatedAt); err != nil {
			return nil, err
		}
		t.Active = active == 1
		tokens = append(tokens, t)
	}
	return tokens, nil
}

func (db *DB) CreateToken(ttype, desc string) (Token, error) {
	token := Token{
		ID:          uuid.NewString(),
		Value:       uuid.NewString(),
		Type:        ttype,
		Description: desc,
		Active:      true,
		CreatedAt:   time.Now().Unix(),
	}
	_, err := db.sql.Exec(`INSERT INTO tokens (id, value, type, description, active, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		token.ID, token.Value, token.Type, token.Description, 1, token.CreatedAt)
	return token, err
}

func (db *DB) ToggleToken(id string) error {
	_, err := db.sql.Exec(`UPDATE tokens SET active = NOT active WHERE id = ?`, id)
	return err
}

func (db *DB) DeleteToken(id string) error {
	_, err := db.sql.Exec(`DELETE FROM tokens WHERE id = ?`, id)
	return err
}

func (db *DB) GetTokenValue(tokenValue string) (*Token, error) {
	row := db.sql.QueryRow(`SELECT id, value, type, description, active, created_at FROM tokens WHERE value = ?`, tokenValue)
	var t Token
	var active int
	err := row.Scan(&t.ID, &t.Value, &t.Type, &t.Description, &active, &t.CreatedAt)
	if err != nil { return nil, err }
	t.Active = active == 1
	return &t, nil
}

// Stats

func (db *DB) AddStat(event, pkg, version string) error {
	_, err := db.sql.Exec(`INSERT INTO stats (event, package, version, occurred_at) VALUES (?, ?, ?, ?)`,
		event, pkg, version, time.Now().Unix())
	return err
}

func (db *DB) RecentStats(n int) ([]StatEvent, error) {
	rows, err := db.sql.Query(`SELECT id, event, package, version, occurred_at FROM stats ORDER BY occurred_at DESC LIMIT ?`, n)
	if err != nil { return nil, err }
	defer rows.Close()
	var events []StatEvent
	for rows.Next() {
		var e StatEvent
		if err := rows.Scan(&e.ID, &e.Event, &e.Package, &e.Version, &e.OccurredAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func (db *DB) CountByEvent(event string) (map[string]int, error) {
	rows, err := db.sql.Query(`SELECT package, COUNT(*) FROM stats WHERE event = ? GROUP BY package`, event)
	if err != nil { return nil, err }
	defer rows.Close()
	m := map[string]int{}
	for rows.Next() {
		var pkg string
		var count int
		if err := rows.Scan(&pkg, &count); err != nil {
			return nil, err
		}
		m[pkg] = count
	}
	return m, nil
}