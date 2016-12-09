package postgres

import (
	"bytes"
	"unicode"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// User SQL queries
const (
	createUserSQL = `INSERT INTO 
        users (username, email, password, actual_name, dob) 
		VALUES ($1, $2, $3, $4, $5) 
        RETURNING id`

	selectUserSQL = `SELECT users.id, users.created_at, users.updated_at,
		users.username, users.email, users.password, users.actual_name, dob`

	getUserByIDSQL = selectUserSQL + ", " +
		`(SELECT COUNT(*) FROM followers WHERE followers.follower_id=users.id) AS num_following,
		 (SELECT COUNT(*) FROM followers WHERE followers.followee_id=users.id) AS num_followers 
		 FROM users WHERE users.id=$1`

	// Don't need follower/following count since these are only used for auth
	getUserByUsernameSQL = selectUserSQL + " FROM users WHERE users.username=$1"
	getUserByEmailSQL    = selectUserSQL + " FROM users WHERE users.email=$1"

	// List queries never return follower/following count
	getFollowersByIDSQL = selectUserSQL +
		` FROM users INNER JOIN followers 
		  ON followers.follower_id=users.id WHERE followers.followee_id=$1`

	getFollowingByIDSQL = selectUserSQL +
		` FROM users INNER JOIN followers 
		  ON followers.followee_id=users.id WHERE followers.follower_id=$1`

	updateUserSQL = `UPDATE users SET 
		username=$1, email=$2, actual_name=$3, dob=$4, updated_at=now() 
		WHERE id=$5`

	deleteUserSQL = `DELETE FROM users WHERE id=$1`

	followUserSQL = `INSERT INTO 
		followers (follower_id, followee_id) 
		VALUES ($1, $2)`

	unFollowUserSQL = `DELETE FROM followers 
		WHERE follower_id=$1 AND followee_id=$2`
)

// Post SQL queries
const (
	createPostSQL = `INSERT INTO 
		posts (author_id, contents) VALUES ($1, $2) RETURNING id`

	selectPostSQL = `SELECT posts.id, posts.created_at, 
		posts.author_id, posts.contents,
		(SELECT COUNT(*) FROM post_keks WHERE post_keks.post_id=posts.id) AS keks, 
		(SELECT COUNT(*) FROM post_nos WHERE post_nos.post_id=posts.id) AS nos`

	getPostByIDSQL      = selectPostSQL + " FROM posts WHERE posts.id=$1"
	getPostsByUserIDSQL = selectPostSQL +
		` FROM posts
		  INNER JOIN users ON posts.author_id=users.id
		  WHERE users.id=$1`
	getFeedByUserIDSQL = selectPostSQL +
		` FROM posts 
		  INNER JOIN users ON posts.author_id=users.id
		  INNER JOIN followers ON followers.followee_id=users.id
		  WHERE followers.follower_id=$1 OR users.id=$1`

	updatePostSQL = `UPDATE posts SET contents=$1, updated_at=now() WHERE id=$2`

	deletePostSQL = `DELETE FROM posts WHERE id=$1`
)

// InitDB creates a postgres database instance using the given connection
// information
func InitDB(user, pass, host, port, database string) (*sqlx.DB, error) {
	db, err := sqlx.Open(
		"postgres",
		"postgres://"+user+":"+pass+"@"+host+":"+port+"/"+database)
	if err != nil {
		return nil, err
	}
	db.MapperFunc(mapperFunc)
	return db, nil
}

func mapperFunc(name string) string {
	if len(name) == 0 {
		return ""
	}

	var buf bytes.Buffer
	var prev rune
	for i, r := range name {
		if i == 0 {
			buf.WriteRune(unicode.ToLower(r))
		} else if unicode.IsUpper(r) && !unicode.IsUpper(prev) {
			buf.WriteRune('_')
			buf.WriteRune(unicode.ToLower(r))
		} else {
			buf.WriteRune(unicode.ToLower(r))
		}
		prev = r
	}
	return buf.String()
}
