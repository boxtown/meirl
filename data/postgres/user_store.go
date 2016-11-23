package postgres

import (
	"database/sql"

	"github.com/boxtown/meirl/data"
	"github.com/jmoiron/sqlx"
)

// UserStore is a PostgreSQL specific implementation
// of data.UserStore
type UserStore struct {
	db *sqlx.DB
}

// NewUserStore returns a newly constructed UserStore
// with the given database reference
func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{db}
}

// Create creates a record for the given user in Postgres
func (store *UserStore) Create(user *data.User) (int64, error) {
	var id int64
	err := store.db.Get(&id, createUserSQL,
		user.Username, user.Email, user.Password,
		user.ActualName, user.DOB.Time)
	if err != nil {
		return 0, data.NewError(err)
	}
	return id, nil
}

// Get retrieves a user by id
func (store *UserStore) Get(id int64) (*data.User, error) {
	var u data.User
	err := store.db.Get(&u, getUserByIDSQL, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, data.ErrNoEnt
		}
		return nil, data.NewError(err)
	}
	return &u, nil
}

// GetByUsername retrieves a user by username
func (store *UserStore) GetByUsername(username string) (*data.User, error) {
	var u data.User
	err := store.db.Get(&u, getUserByUsernameSQL, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, data.ErrNoEnt
		}
		return nil, data.NewError(err)
	}
	return &u, nil
}

// GetByEmail retrieves a user by email
func (store *UserStore) GetByEmail(email string) (*data.User, error) {
	var u data.User
	err := store.db.Get(&u, getUserByEmailSQL, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, data.ErrNoEnt
		}
		return nil, data.NewError(err)
	}
	return &u, nil
}

// Update updates a user by id
func (store *UserStore) Update(id int64, user *data.User) error {
	_, err := store.db.Exec(updateUserSQL,
		user.Username, user.Email, user.ActualName,
		user.DOB.Time, id)
	if err != nil {
		return data.NewError(err)
	}
	return nil
}

// Delete deletes a given user by id
func (store *UserStore) Delete(id int64) error {
	_, err := store.db.Exec(deleteUserSQL, id)
	if err != nil {
		return data.NewError(err)
	}
	return nil
}

// Follow creates a follow relationship between
// the follower and followee
func (store *UserStore) Follow(followerID, followeeID int64) error {
	_, err := store.db.Exec(followUserSQL, followerID, followeeID)
	if err != nil {
		return data.NewError(err)
	}
	return nil
}

// UnFollow idempotently deletes a follow relationship
// between a follower and followee
func (store *UserStore) UnFollow(followerID, followeeID int64) error {
	_, err := store.db.Exec(unFollowUserSQL, followerID, followeeID)
	if err != nil {
		return data.NewError(err)
	}
	return nil
}

// Followers returns a slice of users that are the Followers
// of the user with the given id
func (store *UserStore) Followers(id int64, options data.ListOptions, sort data.UserSortMethod) ([]data.User, error) {
	if options.Limit <= 0 || options.Limit > 1000 {
		options.Limit = 10
	}
	paginator := createUserPaginator(options, sort)
	query := paginator.seekingQuery(getFollowersByIDSQL, 2, true)
	var followers []data.User
	err := store.db.Select(&followers, query, id, options.Marker)
	if err != nil {
		return nil, data.NewError(err)
	}
	return followers, nil
}

// Following returns a slice of users that the user with the given
// id is following
func (store *UserStore) Following(id int64, options data.ListOptions, sort data.UserSortMethod) ([]data.User, error) {
	if options.Limit <= 0 || options.Limit > 1000 {
		options.Limit = 10
	}
	paginator := createUserPaginator(options, sort)
	query := paginator.seekingQuery(getFollowingByIDSQL, 2, true)
	var following []data.User
	err := store.db.Select(&following, query, id, options.Marker)
	if err != nil {
		return nil, data.NewError(err)
	}
	return following, nil
}

func createUserPaginator(options data.ListOptions, sort data.UserSortMethod) *paginator {
	paginator := paginator{
		limit: options.Limit,
		desc:  options.Desc,
	}
	switch sort {
	case data.UserSortByUsername:
		paginator.field = "username"
	case data.UserSortByEmail:
		paginator.field = "email"
	case data.UserSortByActualName:
		paginator.field = "actual_name"
	case data.UserSortByDOB:
		paginator.field = "dob"
	default:
		paginator.field = "id"
	}
	return &paginator
}
