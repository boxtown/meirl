package datatest

import (
	"bytes"
	"time"

	"github.com/boxtown/meirl/data"
)

// UsersEqual returns true if the two users are equivalent
func UsersEqual(a *data.User, b *data.User) bool {
	return a.Username == b.Username &&
		a.Email == b.Email &&
		a.Password == b.Password &&
		a.ActualName == b.ActualName &&
		a.DOB.Equal(b.DOB)
}

// PostsEqual returns true if the two posts are equivalent
func PostsEqual(a *data.Post, b *data.Post) bool {
	return a.AuthorID == b.AuthorID &&
		bytes.Equal(a.Contents, b.Contents) &&
		a.Keks == b.Keks &&
		a.Nos == b.Nos
}

// ExampleUser generates an example user for testing
func ExampleUser() *data.User {
	return &data.User{
		Username:   "test",
		Email:      "test@test.com",
		Password:   "test",
		ActualName: "testy mctest",
		DOB:        data.Time{Time: time.Now()},
	}
}

// ExamplePost generates an example post for testing
func ExamplePost(authorID int64) *data.Post {
	return &data.Post{
		AuthorID: authorID,
		Contents: []byte("test"),
	}
}
