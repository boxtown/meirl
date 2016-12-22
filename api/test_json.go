package api

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/boxtown/meirl/data"
)

/*******************
 * To JSON helpers *
 *******************/

func userToJSON(user *data.User) (io.Reader, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(user)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}

func postToJSON(post *data.Post) (io.Reader, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(post)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}

/* ***************** *
 * From JSON Helpers *
 * ***************** */

func userFromJSON(r io.Reader) (*data.User, error) {
	var u data.User
	err := json.NewDecoder(r).Decode(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func tokenResponseFromJSON(r io.Reader) (*TokenResponse, error) {
	var t TokenResponse
	err := json.NewDecoder(r).Decode(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func postFromJSON(r io.Reader) (*data.Post, error) {
	var p data.Post
	err := json.NewDecoder(r).Decode(&p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func postsFromJSON(r io.Reader) ([]data.Post, error) {
	var p []data.Post
	err := json.NewDecoder(r).Decode(&p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
