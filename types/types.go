package types

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID uuid.UUID
	Name string
	Email string
	Password string `json:"-"`
	Created time.Time
}

type Post struct {
	ID uuid.UUID `json:"id"`
	User string `json:"user"`
	Title string `json:"title"`
	Body string `json:"body"`
	Created time.Time
	Updated time.Time
}
