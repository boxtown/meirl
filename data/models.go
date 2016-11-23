package data

import (
	"fmt"
	"strconv"
	"time"
)

// Time is a wrapper for time.Time that marshals to json
// as seconds since epoch
type Time struct {
	time.Time
}

// MarshalJSON marshals the Time instance to seconds since the epoch
func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", t.Unix())), nil
}

// UnmarshalJSON unmarshals a Time instance from seconds since the
// epoch
func (t *Time) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	unix, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return &errCouldNotUnmarshalTime{data: data}
	}
	t.Time = time.Unix(unix, 0)
	return nil
}

// Scan scans a time.Time into the Time instance
func (t *Time) Scan(src interface{}) error {
	if srcTime, ok := src.(time.Time); ok {
		t.Time = srcTime
		return nil
	}
	return &errCouldNotScanTime{val: src}
}

// Equal returns true if the time instance is equal
// to another time instance. Compares using Unix() since
// we lose some precision from Postgres
func (t *Time) Equal(other Time) bool {
	return t.Time.Unix() == other.Time.Unix()
}

// AutoIncr is a data model for objects created
// with an auto-incrementable integer key in the datastore
type AutoIncr struct {
	ID        int64 `json:"id"`
	CreatedAt Time  `json:"created_at"`
}

// Mutable indicates a type is mutable, giving it a
// `UpdatedAt` field
type Mutable struct {
	AutoIncr
	UpdatedAt Time `json:"updated_at"`
}

// User is the data model for a MeIRL user
type User struct {
	Mutable
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`

	ActualName string `json:"actual_name"`
	DOB        Time   `json:"dob"`

	NumFollowing int `json:"num_following"`
	NumFollowers int `json:"num_followers"`
}

// Post is the data model for a MeIRL post
type Post struct {
	AutoIncr
	AuthorID int64
	Contents []byte
	Keks     int
	Nos      int
}

type errCouldNotUnmarshalTime struct {
	data []byte
}

func (e *errCouldNotUnmarshalTime) Error() string {
	return fmt.Sprintf("Could not marshal %s into data.Time struct", e.data)
}

type errCouldNotScanTime struct {
	val interface{}
}

func (e *errCouldNotScanTime) Error() string {
	return fmt.Sprintf("Could not scan %s into data.Time struct", e.val)
}
