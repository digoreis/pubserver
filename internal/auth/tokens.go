package auth

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeCI   TokenType = "ci"
	TokenTypeUser TokenType = "user"
)

type TokenInfo struct {
	ID          string    `json:"id"`
	Value       string    `json:"value"`
	Type        TokenType `json:"type"`
	Description string    `json:"description"`
	Active      bool      `json:"active"`
	CreatedAt   int64     `json:"created_at"`
}

type TokenStore struct {
	mu     sync.Mutex
	Tokens []TokenInfo
	Path   string
}

func NewTokenStore(path string) (*TokenStore, error) {
	ts := &TokenStore{Path: path}
	if err := ts.Load(); err != nil {
		return nil, err
	}
	return ts, nil
}

func (ts *TokenStore) Load() error {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	f, err := os.Open(ts.Path)
	if err != nil {
		if os.IsNotExist(err) {
			ts.Tokens = []TokenInfo{}
			return nil
		}
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&ts.Tokens)
}

func (ts *TokenStore) Save() error {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	f, err := os.Create(ts.Path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(ts.Tokens)
}

func (ts *TokenStore) List() []TokenInfo {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	return append([]TokenInfo{}, ts.Tokens...)
}

func (ts *TokenStore) Add(t TokenInfo) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.Tokens = append(ts.Tokens, t)
	return ts.Save()
}

func (ts *TokenStore) Toggle(id string) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	for i := range ts.Tokens {
		if ts.Tokens[i].ID == id {
			ts.Tokens[i].Active = !ts.Tokens[i].Active
			return ts.Save()
		}
	}
	return errors.New("token not found")
}

func (ts *TokenStore) Delete(id string) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	for i := range ts.Tokens {
		if ts.Tokens[i].ID == id {
			ts.Tokens = append(ts.Tokens[:i], ts.Tokens[i+1:]...)
			return ts.Save()
		}
	}
	return errors.New("token not found")
}

func (ts *TokenStore) CreateToken(t TokenType, desc string) (TokenInfo, error) {
	token := TokenInfo{
		ID:          uuid.NewString(),
		Value:       uuid.NewString(),
		Type:        t,
		Description: desc,
		Active:      true,
		CreatedAt:   time.Now().Unix(),
	}
	if err := ts.Add(token); err != nil {
		return TokenInfo{}, err
	}
	return token, nil
}