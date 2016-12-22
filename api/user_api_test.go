package api

import (
	"bytes"
	"strconv"
	"testing"

	"net/http"

	"net/http/httptest"

	"strings"

	"github.com/boxtown/meirl/api/apitest"
	"github.com/boxtown/meirl/data"
	"github.com/boxtown/meirl/data/datatest"
	jwt "github.com/dgrijalva/jwt-go"
)

func TestCreateUser(t *testing.T) {
	api := NewUserAPI(
		data.Stores{
			UserStore: mockUserStore{
				OnGetByUsername: func(username string) (*data.User, error) {
					return nil, data.ErrNoEnt
				},
				OnGetByEmail: func(email string) (*data.User, error) {
					return nil, data.ErrNoEnt
				},
				OnCreate: func(user *data.User) (int64, error) {
					return 1, nil
				},
			},
		},
		mockAuth{
			OnSecurePassword: func(password string) (string, error) {
				return password, nil
			},
		},
		false,
	)

	json, _ := userToJSON(datatest.ExampleUser())
	r, _ := http.NewRequest("", "", json)
	w := httptest.NewRecorder()
	api.CreateUser()(w, r)

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

func TestBadCreateUserJSON(t *testing.T) {
	api := NewUserAPI(data.Stores{}, nil, false)
	r, _ := http.NewRequest("", "", bytes.NewBufferString("{"))
	w := httptest.NewRecorder()
	api.CreateUser()(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, received %d", http.StatusBadRequest, w.Code)
		t.Fail()
	}
}

func TestBadCreateUserUsername(t *testing.T) {
	api := NewUserAPI(data.Stores{}, nil, false)
	user := datatest.ExampleUser()
	user.Username = ""
	json, _ := userToJSON(user)
	r, _ := http.NewRequest("", "", json)
	w := httptest.NewRecorder()
	api.CreateUser()(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, received %d", http.StatusBadRequest, w.Code)
		t.Fail()
	}

	user.Username = "test name"
	json, _ = userToJSON(user)
	r, _ = http.NewRequest("", "", json)
	w = httptest.NewRecorder()
	api.CreateUser()(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, received %d", http.StatusBadRequest, w.Code)
		t.Fail()
	}
}

func TestBadCreateUserEmail(t *testing.T) {
	api := NewUserAPI(data.Stores{}, nil, false)
	user := datatest.ExampleUser()
	user.Email = ""
	json, _ := userToJSON(user)
	r, _ := http.NewRequest("", "", json)
	w := httptest.NewRecorder()
	api.CreateUser()(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, received %d", http.StatusBadRequest, w.Code)
		t.Fail()
	}

	user.Email = "test email"
	json, _ = userToJSON(user)
	r, _ = http.NewRequest("", "", json)
	w = httptest.NewRecorder()
	api.CreateUser()(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, received %d", http.StatusBadRequest, w.Code)
		t.Fail()
	}
}

func TestTakenCreateUserUsernameOrEmail(t *testing.T) {
	api := NewUserAPI(
		data.Stores{
			UserStore: mockUserStore{
				OnGetByUsername: func(username string) (*data.User, error) {
					return nil, nil
				},
				OnGetByEmail: func(email string) (*data.User, error) {
					return nil, nil
				},
			},
		},
		nil,
		false,
	)
	json, _ := userToJSON(datatest.ExampleUser())
	r, _ := http.NewRequest("", "", json)
	w := httptest.NewRecorder()
	api.CreateUser()(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, received %d", http.StatusBadRequest, w.Code)
		t.Fail()
	}
}

func TestGetUser(t *testing.T) {
	stored := datatest.ExampleUser()
	api := NewUserAPI(
		data.Stores{
			UserStore: mockUserStore{
				OnGet: func(id int64) (*data.User, error) {
					return stored, nil
				},
			},
		},
		nil,
		false,
	)

	r, _ := http.NewRequest("", "", nil)
	r = apitest.RequestWithContextID(r, idContextKey, int64(1))
	w := httptest.NewRecorder()
	api.GetUser()(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, received %d", http.StatusOK, w.Code)
		t.Fail()
	}
	user, err := userFromJSON(w.Body)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
	}
	if !datatest.UsersEqual(user, stored) {
		t.Error("Retrieved user did not equal stored user")
		t.Fail()
	}
}

func TestGetMe(t *testing.T) {
	stored := datatest.ExampleUser()
	api := NewUserAPI(
		data.Stores{
			UserStore: mockUserStore{
				OnGet: func(id int64) (*data.User, error) {
					return stored, nil
				},
			},
		},
		nil,
		false,
	)
	r, _ := http.NewRequest("", "", nil)
	r = apitest.RequestWithClaims(r, claimsContextKey, jwt.MapClaims{
		"sub": int64(1),
	})
	w := httptest.NewRecorder()
	api.GetMe()(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, received %d", http.StatusOK, w.Code)
		t.Fail()
	}
	user, err := userFromJSON(w.Body)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
	}
	if !datatest.UsersEqual(user, stored) {
		t.Error("Retrieved user did not equal stored user")
		t.Fail()
	}
}

func TestGetFeed(t *testing.T) {
	api := NewUserAPI(
		data.Stores{
			UserStore: mockUserStore{
				OnGet: func(id int64) (*data.User, error) {
					return nil, nil
				},
			},
			PostStore: mockPostStore{
				OnFeed: func(
					userID int64,
					options data.ListOptions,
					sort data.PostSortMethod) ([]data.Post, error) {
					return []data.Post{
						*datatest.ExamplePost(1),
						*datatest.ExamplePost(2),
					}, nil
				},
			},
		},
		nil,
		false,
	)
	r, _ := http.NewRequest("", "", nil)
	r = apitest.RequestWithContextID(r, idContextKey, int64(1))
	w := httptest.NewRecorder()
	api.GetFeed()(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, received %d", http.StatusOK, w.Code)
		t.Fail()
	}
	posts, err := postsFromJSON(w.Body)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
	}
	if len(posts) != 2 {
		t.Error("Wrong number of posts returned")
		t.Fail()
	}
}

func TestFollowerUser(t *testing.T) {
	api := NewUserAPI(
		data.Stores{
			UserStore: mockUserStore{
				OnGet: func(id int64) (*data.User, error) {
					return nil, nil
				},
				OnFollow: func(followerID, followeeID int64) error {
					return nil
				},
			},
		},
		nil,
		false,
	)

	r, _ := http.NewRequest("", "", nil)
	r = apitest.RequestWithContextID(r, idContextKey, int64(1))
	r = apitest.RequestWithClaims(r, claimsContextKey, jwt.MapClaims{
		"sub": int64(1),
	})
	w := httptest.NewRecorder()
	api.FollowUser()(w, r)

	if w.Code != http.StatusAccepted {
		t.Errorf("Expected status code %d, received %d", http.StatusAccepted, w.Code)
		t.Fail()
	}
}

func TestDeleteUser(t *testing.T) {
	api := NewUserAPI(
		data.Stores{
			UserStore: mockUserStore{
				OnDelete: func(id int64) error {
					return nil
				},
			},
		},
		nil,
		false,
	)

	r, _ := http.NewRequest("", "", nil)
	r = apitest.RequestWithContextID(r, idContextKey, int64(1))
	w := httptest.NewRecorder()
	api.DeleteUser()(w, r)

	if w.Code != http.StatusAccepted {
		t.Errorf("Expected status code %d, received %d", http.StatusAccepted, w.Code)
		t.Fail()
	}
}

func TestLoginUser(t *testing.T) {
	api := NewUserAPI(
		data.Stores{
			UserStore: mockUserStore{
				OnGetByUsername: func(username string) (*data.User, error) {
					return datatest.ExampleUser(), nil
				},
			},
		},
		mockAuth{
			OnCheckPassword: func(password, storedPassword string) bool {
				return password == storedPassword
			},
			OnGenerateAccessToken: func(user *data.User, signingKey []byte) (string, error) {
				return "test-token", nil
			},
		},
		false,
	)

	json, _ := userToJSON(datatest.ExampleUser())
	r, _ := http.NewRequest("", "", json)
	w := httptest.NewRecorder()
	api.Login([]byte("test"))(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, received %d", http.StatusOK, w.Code)
		t.Fail()
	}
	tr, err := tokenResponseFromJSON(w.Body)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
	}
	if tr.AccessToken != "test-token" {
		t.Error("Wrong token returned")
		t.Fail()
	}
}
