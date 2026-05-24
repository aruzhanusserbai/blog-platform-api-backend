package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func performRequest(handler gin.HandlerFunc, method, route, target, body string, userID uint) *httptest.ResponseRecorder {
	router := gin.New()
	router.Handle(method, route, func(c *gin.Context) {
		if userID != 0 {
			c.Set("user_id", userID)
		}
		handler(c)
	})

	req := httptest.NewRequest(method, target, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	return recorder
}

func assertStatus(t *testing.T, recorder *httptest.ResponseRecorder, expected int) {
	t.Helper()

	if recorder.Code != expected {
		t.Fatalf("expected status %d, got %d: %s", expected, recorder.Code, recorder.Body.String())
	}
}

func TestRegisterRejectsMissingRequiredFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := performRequest(Register, http.MethodPost, "/register", "/register", `{"username":"alice"}`, 0)

	assertStatus(t, recorder, http.StatusBadRequest)
}

func TestLoginRejectsMissingPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := performRequest(LogIn, http.MethodPost, "/login", "/login", `{"username":"alice"}`, 0)

	assertStatus(t, recorder, http.StatusBadRequest)
}

func TestCreateUserRejectsMalformedJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := performRequest(CreateUser, http.MethodPost, "/users", "/users", `{"username":`, 0)

	assertStatus(t, recorder, http.StatusBadRequest)
}

func TestCreatePostRejectsMalformedJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := performRequest(CreatePost, http.MethodPost, "/posts", "/posts", `{"title":`, 1)

	assertStatus(t, recorder, http.StatusBadRequest)
}

func TestGetPostRejectsInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := performRequest(GetPost, http.MethodGet, "/posts/:id", "/posts/not-a-number", "", 0)

	assertStatus(t, recorder, http.StatusBadRequest)
}

func TestAddCommentToPostRejectsInvalidPostID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := performRequest(AddCommentToPost, http.MethodPost, "/posts/:id/comments", "/posts/not-a-number/comments", `{"content":"Nice"}`, 1)

	assertStatus(t, recorder, http.StatusBadRequest)
}

func TestAddCommentToPostRejectsMalformedJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := performRequest(AddCommentToPost, http.MethodPost, "/posts/:id/comments", "/posts/1/comments", `{"content":`, 1)

	assertStatus(t, recorder, http.StatusBadRequest)
}

func TestGetCommentsRejectsInvalidPostID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := performRequest(GetComments, http.MethodGet, "/posts/:id/comments", "/posts/not-a-number/comments", "", 0)

	assertStatus(t, recorder, http.StatusBadRequest)
}

func TestCreateTagsRejectsMalformedJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := performRequest(CreateTags, http.MethodPost, "/tags", "/tags", `{"name":`, 0)

	assertStatus(t, recorder, http.StatusBadRequest)
}

func TestAddTagsToPostRejectsInvalidPostID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := performRequest(AddTagsToPost, http.MethodPost, "/posts/:id/tags", "/posts/not-a-number/tags", `{"tag_ids":[1]}`, 0)

	assertStatus(t, recorder, http.StatusBadRequest)
}
