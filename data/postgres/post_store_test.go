package postgres

import (
	"testing"
	"time"

	"github.com/boxtown/gotag"
	"github.com/boxtown/meirl/data"
	"github.com/boxtown/meirl/data/datatest"
	"github.com/jmoiron/sqlx"
)

func TestCreatePost(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanupPostStoreTest(t, db)

			userID, err := populateUsersTable(t, db, datatest.ExampleUser())
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			store := NewPostStore(db)
			id, err := store.Create(datatest.ExamplePost(userID))
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			if id < 1 {
				t.Error("Invalid user ID generated")
				t.Fail()
			}
			return nil
		})
	})
}

func TestCreatePostWithBadAuthorID(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanupPostStoreTest(t, db)

			store := NewPostStore(db)
			_, err := store.Create(datatest.ExamplePost(0))
			if err == nil {
				t.Error("Create post with bad author ID should not have succeeded")
				t.Fail()
			}
			return nil
		})
	})
}

func TestGetPost(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanupPostStoreTest(t, db)

			userID, err := populateUsersTable(t, db, datatest.ExampleUser())
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			store := NewPostStore(db)
			check := datatest.ExamplePost(userID)
			id, err := store.Create(check)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			err = populateKeksTable(t, db, 3, userID, id)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			err = populateNosTable(t, db, 2, userID, id)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			post, err := store.Get(id)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			check.Keks, check.Nos = 3, 2
			if !datatest.PostsEqual(post, check) {
				t.Error("Post returned by Get has incorrect fields")
				t.Fail()
			}
			return nil
		})
	})
}

func TestGetNonExistentPost(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanupPostStoreTest(t, db)

			store := NewPostStore(db)
			_, err := store.Get(1)
			if err == nil {
				t.Error("Get non-existent post should not have succeeded")
				t.FailNow()
			}
			if err != data.ErrNoEnt {
				t.Error(err.Error())
				t.Fail()
			}
			return nil
		})
	})
}

func TestUpdatePost(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanupPostStoreTest(t, db)

			userID, err := populateUsersTable(t, db, datatest.ExampleUser())
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			store := NewPostStore(db)
			check := datatest.ExamplePost(userID)
			id, err := store.Create(check)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			check.Contents = []byte("updated")
			err = store.Update(id, check.Contents)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			post, err := store.Get(id)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			if !datatest.PostsEqual(post, check) {
				t.Error("Update failed")
				t.Fail()
			}
			return nil
		})
	})
}

func TestUpdatePostIsIdempotent(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanupPostStoreTest(t, db)

			store := NewPostStore(db)
			err := store.Update(3, []byte("updated"))
			if err != nil {
				t.Error(err.Error())
				t.Fail()
			}
			return nil
		})
	})
}

func TestDeletePost(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanupPostStoreTest(t, db)

			userID, err := populateUsersTable(t, db, datatest.ExampleUser())
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			store := NewPostStore(db)
			id, err := store.Create(datatest.ExamplePost(userID))
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			err = store.Delete(id)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			_, err = store.Get(id)
			if err != data.ErrNoEnt {
				t.Error("Post was not properly deleted")
				t.Fail()
			}
			return nil
		})
	})
}

func TestDeletePostIsIdempotent(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanupPostStoreTest(t, db)

			store := NewPostStore(db)
			err := store.Delete(1)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			return nil
		})
	})
}

func TestGetUserPosts(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanupPostStoreTest(t, db)

			userID, err := populateUsersTable(t, db, datatest.ExampleUser())
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			store := NewPostStore(db)
			postID, err := store.Create(datatest.ExamplePost(userID))
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			posts, err := store.UserPosts(
				userID,
				data.ListOptions{Marker: time.Now(), Desc: true},
				data.PostSortByDate,
			)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			if len(posts) != 1 {
				t.Error("Incorrect number of user posts retrieved")
				t.FailNow()
			}
			if posts[0].ID != postID {
				t.Error("Incorrect user posts retrieved")
				t.Fail()
			}
			return nil
		})
	})
}

func TestGetPostFeed(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanupPostStoreTest(t, db)

			followerID, err := populateUsersTable(t, db, datatest.ExampleUser())
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			followee := datatest.ExampleUser()
			followee.Username = "test2"
			followee.Email = "test2@test.com"
			followeeID, err := populateUsersTable(t, db, followee)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			err = populateFollowersTable(t, db, followerID, followeeID)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}

			store := NewPostStore(db)
			postID, err := store.Create(datatest.ExamplePost(followeeID))
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			feed, err := store.Feed(
				followerID,
				data.ListOptions{Marker: time.Now(), Desc: true},
				data.PostSortByDate,
			)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			if len(feed) != 1 {
				t.Error("Incorrect number of feed posts retrieved")
				t.FailNow()
			}
			if feed[0].ID != postID {
				t.Error("Inccorrect feed posts retrieved")
				t.Fail()
			}
			return nil
		})
	})
}

func populateUsersTable(t gotag.T, db *sqlx.DB, user *data.User) (int64, error) {
	store := NewUserStore(db)
	return store.Create(user)
}

func populateFollowersTable(t gotag.T, db *sqlx.DB, followerID, followeeID int64) error {
	store := NewUserStore(db)
	return store.Follow(followerID, followeeID)
}

func populateKeksTable(t gotag.T, db *sqlx.DB, keks int, authorID, postID int64) error {
	for i := 0; i < keks; i++ {
		_, err := db.Exec(
			"INSERT INTO post_keks (author_id, post_id) VALUES ($1, $2)",
			authorID, postID)
		if err != nil {
			return err
		}
	}
	return nil
}

func populateNosTable(t gotag.T, db *sqlx.DB, nos int, authorID, postID int64) error {
	for i := 0; i < nos; i++ {
		_, err := db.Exec(
			"INSERT INTO post_nos (author_id, post_id) VALUES ($1, $2)",
			authorID, postID)
		if err != nil {
			return err
		}
	}
	return nil
}

func cleanupPostStoreTest(t gotag.T, db *sqlx.DB) {
	_, err := db.Exec("DELETE FROM post_keks")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = db.Exec("DELETE FROM post_nos")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = db.Exec("DELETE FROM posts")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = db.Exec("DELETE FROM followers")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = db.Exec("DELETE FROM users")
	if err != nil {
		t.Fatal(err.Error())
	}
}
