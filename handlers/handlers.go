package handlers

import (
	"encoding/json"
	"github.com/adigunhammedolalekan/rest-unit-testing-sample/repository"
	"github.com/go-chi/render"
	"net/http"
)

type response struct {
	Success bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Handler struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Handler {
	return &Handler{repo: repo}
}

func (handler *Handler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var body struct{
		Title string `json:"title"`
		Body string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.Respond(w, r, &response{Success: false, Message: "malformed body"})
		return
	}
	token := r.Header.Get("x-auth-token")
	user, err := handler.repo.GetUser(token)
	if err != nil {
		render.Status(r, http.StatusForbidden)
		render.Respond(w, r, &response{Success: false, Message: "forbidden"})
		return
	}
	p, err := handler.repo.CreatePost(user.ID.String(), body.Title, body.Body)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.Respond(w, r, &response{Success: false, Message: "server error"})
		return
	}
	render.Status(r, http.StatusOK)
	render.Respond(w, r, &response{Success: true, Message: "post.created", Data: p})
}

func (handler *Handler) GetPosts(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("x-auth-token")
	user, err := handler.repo.GetUser(token)
	if err != nil {
		render.Status(r, http.StatusForbidden)
		render.Respond(w, r, &response{Success: false, Message: "forbidden"})
		return
	}
	posts, err := handler.repo.GetPosts(user.ID.String())
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.Respond(w, r, &response{Success: false, Message: "server error"})
		return
	}
	render.Status(r, http.StatusOK)
	render.Respond(w, r, &response{Success: true, Message: "post.retrieved", Data: posts})
}
