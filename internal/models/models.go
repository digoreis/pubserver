package models

// Estruturas principais do domínio
type Package struct {
    Name     string
    Versions []string
}

type User struct {
    Username string
    Token    string
}