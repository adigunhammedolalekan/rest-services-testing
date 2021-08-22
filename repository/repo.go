package repository

import (
	"errors"
	"github.com/adigunhammedolalekan/rest-unit-testing-sample/types"
	"github.com/google/uuid"
)

type Database map[string] []*types.Post

func (d Database) addPost(p *types.Post) (*types.Post, error) {
	d[p.User] = append(d[p.User], p)
	return p, nil
}

func (d Database) posts(userId string) []*types.Post {
	return d[userId]
}

//go:generate mockgen -source=repo.go -destination=../mocks/repository_mock.go -package=mocks
type Repository interface {
	CreatePost(userId, title, body string) (*types.Post, error)
	GetPosts(userId string) ([]*types.Post, error)
	GetUser(userId string) (*types.User, error)
}

type repository struct {
	db Database
	userId uuid.UUID // single user system, keeping things simple :)
}

func New(db Database) Repository {
	return &repository{userId: uuid.New(), db: db}
}

func (r *repository) CreatePost(userId, title, body string) (*types.Post, error) {
	p := &types.Post{
		ID:      uuid.New(),
		User:    userId,
		Title:   title,
		Body:    body,
	}
	return r.db.addPost(p)
}

func (r *repository) GetPosts(userId string) ([]*types.Post, error) {
	return r.db.posts(userId), nil
}

func (r *repository) GetUser(token string) (*types.User, error) {
	if token == "" {
		return nil, errors.New("user not authorized")
	}
	return &types.User{
		ID:       r.userId,
		Name:     "Admin",
		Email:    "admin@testing.app",
	}, nil
}
