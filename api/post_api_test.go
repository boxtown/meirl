package api

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/boxtown/meirl/data"
	"github.com/boxtown/meirl/data/datatest"
	jwt "github.com/dgrijalva/jwt-go"
)

func TestCreatePost(t *testing.T) {
	api := NewPostAPI(
		data.Stores{
			PostStore: mockPostStore{
				OnCreate: func(post *data.Post) (int64, error) {
					return 1, nil
				},
			},
		},
		false,
	)

	json, _ := postToJSON(datatest.ExamplePost(1))
	r, _ := http.NewRequest("", "", json)
	r = r.WithContext(context.WithValue(r.Context(), claimsContextKey, jwt.MapClaims{
		"sub": int64(1),
	}))
	w := httptest.NewRecorder()
	api.CreatePost()(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, received %d", http.StatusCreated, w.Code)
		t.Fail()
	}
	location := w.HeaderMap.Get("Location")
	if location == "" {
		t.Errorf("Expected non-empty location header")
		t.Fail()
	}
	_, err := strconv.Atoi(location[strings.LastIndex(location, "/")+1:])
	if err != nil {
		t.Errorf("Expected generated id as last path element in location %s", location)
		t.Fail()
	}
}

func TestGetPost(t *testing.T) {
	stored := datatest.ExamplePost(1)
	api := NewPostAPI(
		data.Stores{
			PostStore: mockPostStore{
				OnGet: func(id int64) (*data.Post, error) {
					return stored, nil
				},
			},
		},
		false,
	)

	r, _ := http.NewRequest("", "", nil)
	r = r.WithContext(context.WithValue(r.Context(), idContextKey, int64(1)))
	w := httptest.NewRecorder()
	api.GetPost()(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, received %d", http.StatusOK, w.Code)
		t.Fail()
	}
	post, err := postFromJSON(w.Body)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
	}
	if !datatest.PostsEqual(post, stored) {
		t.Error("Retrieved post did not equal stored post")
		t.Fail()
	}
}
