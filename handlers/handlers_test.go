package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/adigunhammedolalekan/rest-unit-testing-sample/mocks"
	"github.com/adigunhammedolalekan/rest-unit-testing-sample/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_CreatePostHandler(t *testing.T) {
	w := httptest.NewRecorder()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type body struct{
		Title string `json:"title"`
		Body string `json:"body"`
	}

	mockToken := uuid.New().String()
	mockPost := &body{Title: "Test Title", Body: "Test Body"}
	mockUser := &types.User{ID: uuid.New(), Name: "Tester"}
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(mockPost)
	assert.Nil(t, err)

	r := httptest.NewRequest("POST", "/", buf)
	r.Header.Add("x-auth-token", mockToken)

	repo := mocks.NewMockRepository(ctrl)
	repo.EXPECT().GetUser(mockToken).Return(mockUser, nil).Times(1)
	repo.EXPECT().CreatePost(mockUser.ID.String(), mockPost.Title, mockPost.Body).Return(&types.Post{}, nil).Times(1)

	handler := New(repo)
	handler.CreatePostHandler(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, strings.Contains(w.Body.String(), "post.created"))
}

func TestHandler_CreatePostHandler_BadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := httptest.NewRequest("POST", "/", bytes.NewBufferString("{malformed json data}"))
	r.Header.Add("x-auth-token", uuid.New().String())

	repo := mocks.NewMockRepository(ctrl)
	handler := New(repo)
	handler.CreatePostHandler(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.True(t, strings.Contains(w.Body.String(), "malformed body"))
}

func TestHandler_CreatePostHandler_BadAuth(t *testing.T) {
	w := httptest.NewRecorder()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type body struct{
		Title string `json:"title"`
		Body string `json:"body"`
	}

	mockToken := uuid.New().String()
	mockPost := &body{Title: "Test Title", Body: "Test Body"}
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(mockPost)
	assert.Nil(t, err)

	r := httptest.NewRequest("POST", "/", buf)
	r.Header.Add("x-auth-token", mockToken)

	repo := mocks.NewMockRepository(ctrl)
	repo.EXPECT().GetUser(mockToken).Return(nil, errors.New("user not found")).Times(1)
	handler := New(repo)
	handler.CreatePostHandler(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.True(t, strings.Contains(w.Body.String(), "forbidden"))
}
