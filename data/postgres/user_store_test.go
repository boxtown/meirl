package postgres

import (
	"testing"

	"github.com/boxtown/gotag"
	"github.com/boxtown/meirl/data"
	"github.com/boxtown/meirl/data/datatest"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func TestCreateUser(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanUserStoreTest(t, db)

			store := NewUserStore(db)
			id, err := store.Create(datatest.ExampleUser())
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

func TestCreateUserWithDuplicateUserName(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanUserStoreTest(t, db)

			store := NewUserStore(db)
			user := datatest.ExampleUser()
			_, err := store.Create(user)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}

			user.Email = "test2@test.com"
			_, err = store.Create(user)
			if err == nil {
				t.Error("Create user with duplicate username should not have succeeded")
				t.Fail()
			}
			return nil
		})
	})
}

func TestCreateUserWithDuplicateEmail(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanUserStoreTest(t, db)

			store := NewUserStore(db)
			user := datatest.ExampleUser()
			_, err := store.Create(user)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}

			user.Username = "test2"
			_, err = store.Create(user)
			if err == nil {
				t.Error("Create user with duplicate email should not have succeeded")
				t.Fail()
			}
			return nil
		})
	})
}

func TestCreateUserWithMissingUserName(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanUserStoreTest(t, db)

			store := NewUserStore(db)
			user := datatest.ExampleUser()
			user.Username = ""
			_, err := store.Create(user)
			if err == nil {
				t.Error("Create user with missing username should not have succeeded")
				t.Fail()
			}
			return nil
		})
	})
}

func TestCreateUserWithMissingEmail(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanUserStoreTest(t, db)

			store := NewUserStore(db)
			user := datatest.ExampleUser()
			user.Email = ""
			_, err := store.Create(user)
			if err == nil {
				t.Error("Create user with missing email should not have succeeded")
				t.Fail()
			}
			return nil
		})
	})
}

func TestCreateUserWithMissingPassword(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanUserStoreTest(t, db)

			store := NewUserStore(db)
			user := datatest.ExampleUser()
			user.Password = ""
			_, err := store.Create(user)
			if err == nil {
				t.Error("Create user with missing password should not have succeeded")
				t.Fail()
			}
			return nil
		})
	})
}

func TestCreateUserWithMissingActualName(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanUserStoreTest(t, db)

			store := NewUserStore(db)
			user := datatest.ExampleUser()
			user.ActualName = ""
			_, err := store.Create(user)
			if err == nil {
				t.Error("Create user with missing actual name should not have succeeded")
				t.Fail()
			}
			return nil
		})
	})
}

func TestGetUserByID(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanUserStoreTest(t, db)

			store := NewUserStore(db)
			check := datatest.ExampleUser()
			id, err := store.Create(check)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			user, err := store.Get(id)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			if !datatest.UsersEqual(user, check) {
				t.Error("User returned by Get has incorrect fields")
				t.Fail()
			}
			return nil
		})
	})
}

func TestGetUserWithNonExistentID(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanUserStoreTest(t, db)

			store := NewUserStore(db)
			_, err := store.Get(1)
			if err == nil {
				t.Error("Get non-existent user should not have succeeded")
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

func TestGetUserByUsername(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanUserStoreTest(t, db)

			store := NewUserStore(db)
			check := datatest.ExampleUser()
			_, err := store.Create(check)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			user, err := store.GetByUsername(check.Username)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			if !datatest.UsersEqual(user, check) {
				t.Error("User returned by Get has incorrect fields")
				t.Fail()
			}
			return nil
		})
	})
}

func TestGetUserWithNonExistentUsername(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanUserStoreTest(t, db)

			store := NewUserStore(db)
			_, err := store.GetByUsername("")
			if err == nil {
				t.Error("Get non-existent user should not have succeeded")
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

func TestGetUserByEmail(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanUserStoreTest(t, db)

			store := NewUserStore(db)
			check := datatest.ExampleUser()
			_, err := store.Create(check)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			user, err := store.GetByEmail(check.Email)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			if !datatest.UsersEqual(user, check) {
				t.Error("User returned by Get has incorrect fields")
				t.Fail()
			}
			return nil
		})
	})
}

func TestGetUserByNonExistentEmail(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanUserStoreTest(t, db)

			store := NewUserStore(db)
			_, err := store.GetByEmail("")
			if err == nil {
				t.Error("Get non-existent user should not have succeeded")
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

func TestUpdateUser(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanUserStoreTest(t, db)

			store := NewUserStore(db)
			check := datatest.ExampleUser()
			id, err := store.Create(check)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			check.Username = "updated username"
			check.Email = "updated@updated.com"
			check.ActualName = "updated actual name"
			newTime := data.Time{}
			check.DOB = newTime
			err = store.Update(id, check)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			user, err := store.Get(id)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			if !datatest.UsersEqual(user, check) {
				t.Error("Update failed")
				t.Fail()
			}
			return nil
		})
	})
}

func TestUpdateUserDoesNotUpdatePassword(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanUserStoreTest(t, db)

			store := NewUserStore(db)
			check := datatest.ExampleUser()
			id, err := store.Create(check)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			check.Password = "updated"
			err = store.Update(id, check)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			user, err := store.Get(id)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			if user.Password != datatest.ExampleUser().Password {
				t.Error("Update should not change user password")
				t.Fail()
			}
			return nil
		})
	})
}

func TestUpdateUserIsIdempotent(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanUserStoreTest(t, db)

			store := NewUserStore(db)
			err := store.Update(3, datatest.ExampleUser())
			if err != nil {
				t.Error(err.Error())
				t.Fail()
			}
			return nil
		})
	})
}

func TestDeleteUser(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanUserStoreTest(t, db)

			store := NewUserStore(db)
			id, err := store.Create(datatest.ExampleUser())
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
				t.Error("User was not properly deleted")
				t.Fail()
			}
			return nil
		})
	})
}

func TestDeleteUserIsIdempotent(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanUserStoreTest(t, db)

			store := NewUserStore(db)
			err := store.Delete(1)
			if err != nil {
				t.Error(err.Error())
				t.Fail()
			}
			return nil
		})
	})
}

func TestFollowUser(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanUserStoreTest(t, db)

			store := NewUserStore(db)
			followerID, err := store.Create(datatest.ExampleUser())
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			followee := datatest.ExampleUser()
			followee.Username = "test2"
			followee.Email = "test2@email.com"
			followeeID, err := store.Create(followee)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			err = store.Follow(followerID, followeeID)
			if err != nil {
				t.Error(err.Error())
				t.Fail()
			}
			return nil
		})
	})
}

func TestUnFollowUser(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanUserStoreTest(t, db)

			store := NewUserStore(db)
			followerID, err := store.Create(datatest.ExampleUser())
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			followee := datatest.ExampleUser()
			followee.Username = "test2"
			followee.Email = "test2@email.com"
			followeeID, err := store.Create(followee)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			err = store.Follow(followerID, followeeID)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			err = store.UnFollow(followerID, followeeID)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			return nil
		})
	})
}

func TestGetFollowers(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanUserStoreTest(t, db)

			store := NewUserStore(db)
			followerID, err := store.Create(datatest.ExampleUser())
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			followee := datatest.ExampleUser()
			followee.Username = "test2"
			followee.Email = "test2@email.com"
			followeeID, err := store.Create(followee)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			err = store.Follow(followerID, followeeID)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			followers, err := store.Followers(
				followeeID,
				data.ListOptions{Marker: 0},
				data.UserSortByID,
			)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			if len(followers) != 1 {
				t.Error("Incorrect number of followers retrieved")
				t.FailNow()
			}
			if followers[0].ID != followerID {
				t.Error("Incorrect followers retrieved")
				t.Fail()
			}
			return nil
		})
	})
}

func TestGetFollowing(t *testing.T) {
	gotag.Test(gotag.Integration, t, func(t gotag.T) {
		whileConnectedToTestDb(testDbName, func(db *sqlx.DB) error {
			defer cleanUserStoreTest(t, db)

			store := NewUserStore(db)
			followerID, err := store.Create(datatest.ExampleUser())
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			followee := datatest.ExampleUser()
			followee.Username = "test2"
			followee.Email = "test2@email.com"
			followeeID, err := store.Create(followee)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			err = store.Follow(followerID, followeeID)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			following, err := store.Following(
				followerID,
				data.ListOptions{Marker: 0},
				data.UserSortByID,
			)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			if len(following) != 1 {
				t.Error("Incorrect number of followed users retrieved")
				t.FailNow()
			}
			if following[0].ID != followeeID {
				t.Error("Incorrect followed users retrieved")
				t.Fail()
			}
			return nil
		})
	})
}

func cleanUserStoreTest(t gotag.T, db *sqlx.DB) {
	_, err := db.Exec("DELETE FROM followers")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = db.Exec("DELETE FROM users")
	if err != nil {
		t.Fatal(err.Error())
	}
}
