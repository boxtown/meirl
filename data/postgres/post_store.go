package postgres

import (
	"database/sql"

	"github.com/boxtown/meirl/data"
	"github.com/jmoiron/sqlx"
)

// PostStore is a PostgreSQL specific implementation
// of data.PostStore
type PostStore struct {
	db *sqlx.DB
}

// NewPostStore returns a newly constructed PostStore
// with the given database reference
func NewPostStore(db *sqlx.DB) *PostStore {
	return &PostStore{db}
}

// Create creates a record for the given post in Postgres
func (store *PostStore) Create(post *data.Post) (int64, error) {
	var id int64
	err := store.db.Get(&id, createPostSQL, post.AuthorID, post.Contents)
	if err != nil {
		return 0, data.NewError(err)
	}
	return id, nil
}

// Get retrieves a post by id
func (store *PostStore) Get(id int64) (*data.Post, error) {
	var p data.Post
	err := store.db.Get(&p, getPostByIDSQL, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, data.ErrNoEnt
		}
		return nil, data.NewError(err)
	}
	return &p, nil
}

// Update updates a post by id
func (store *PostStore) Update(id int64, contents []byte) error {
	_, err := store.db.Exec(updatePostSQL, contents, id)
	if err != nil {
		return data.NewError(err)
	}
	return nil
}

// Delete deletes a post by id
func (store *PostStore) Delete(id int64) error {
	_, err := store.db.Exec(deletePostSQL, id)
	if err != nil {
		return data.NewError(err)
	}
	return nil
}

// UserPosts returns the posts for the user with
// the given id
func (store *PostStore) UserPosts(userID int64, options data.ListOptions, sort data.PostSortMethod) ([]data.Post, error) {
	if options.Limit <= 10 || options.Limit > 1000 {
		options.Limit = 10
	}
	paginator := createPostsPaginator(options, sort)
	query := paginator.seekingQuery(getPostsByUserIDSQL, 2, true)
	var posts []data.Post
	err := store.db.Select(&posts, query, userID, options.Marker)
	if err != nil {
		return nil, data.NewError(err)
	}
	return posts, nil
}

// Feed retrieves the post feed for the user with
// the given id
func (store *PostStore) Feed(userID int64, options data.ListOptions, sort data.PostSortMethod) ([]data.Post, error) {
	if options.Limit <= 10 || options.Limit > 1000 {
		options.Limit = 10
	}
	paginator := createPostsPaginator(options, sort)
	query := paginator.seekingQuery(getFeedByUserIDSQL, 2, true)
	var posts []data.Post
	err := store.db.Select(&posts, query, userID, options.Marker)
	if err != nil {
		return nil, data.NewError(err)
	}
	return posts, nil
}

func createPostsPaginator(options data.ListOptions, sort data.PostSortMethod) *paginator {
	paginator := paginator{
		limit: options.Limit,
		desc:  options.Desc,
	}
	switch sort {
	case data.PostSortByKeks:
		paginator.field = "keks"
	case data.PostSortByNos:
		paginator.field = "nos"
	default:
		paginator.field = "posts.created_at"
	}
	return &paginator
}
