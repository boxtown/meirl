package api

import "github.com/boxtown/meirl/data"

/* *************** *
 * Mock User Store *
 * *************** */

type mockUserStore struct {
	OnCreate        func(user *data.User) (int64, error)
	OnGet           func(id int64) (*data.User, error)
	OnGetByUsername func(username string) (*data.User, error)
	OnGetByEmail    func(email string) (*data.User, error)
	OnUpdate        func(id int64, user *data.User) error
	OnDelete        func(id int64) error
	OnFollow        func(followerID, followeeID int64) error
	OnUnFollow      func(followerID, followeeID int64) error

	OnFollowers func(
		id int64,
		options data.ListOptions,
		sort data.UserSortMethod) ([]data.User, error)

	OnFollowing func(
		id int64,
		options data.ListOptions,
		sort data.UserSortMethod) ([]data.User, error)
}

func (store mockUserStore) Create(user *data.User) (int64, error) {
	return store.OnCreate(user)
}

func (store mockUserStore) Get(id int64) (*data.User, error) {
	return store.OnGet(id)
}

func (store mockUserStore) GetByUsername(username string) (*data.User, error) {
	return store.OnGetByUsername(username)
}

func (store mockUserStore) GetByEmail(email string) (*data.User, error) {
	return store.OnGetByEmail(email)
}

func (store mockUserStore) Update(id int64, user *data.User) error {
	return store.OnUpdate(id, user)
}

func (store mockUserStore) Delete(id int64) error {
	return store.OnDelete(id)
}

func (store mockUserStore) Follow(followerID, followeeID int64) error {
	return store.OnFollow(followerID, followeeID)
}

func (store mockUserStore) UnFollow(followerID, followeeID int64) error {
	return store.OnUnFollow(followerID, followeeID)
}

func (store mockUserStore) Followers(
	id int64,
	options data.ListOptions,
	sort data.UserSortMethod) ([]data.User, error) {
	return store.OnFollowers(id, options, sort)
}

func (store mockUserStore) Following(
	id int64,
	options data.ListOptions,
	sort data.UserSortMethod) ([]data.User, error) {
	return store.OnFollowing(id, options, sort)
}

/* *************** *
 * Mock Post Store *
 * *************** */

type mockPostStore struct {
	OnCreate func(post *data.Post) (int64, error)
	OnGet    func(id int64) (*data.Post, error)
	OnUpdate func(id int64, contents []byte) error
	OnDelete func(id int64) error

	OnUserPosts func(
		userID int64,
		options data.ListOptions,
		sort data.PostSortMethod) ([]data.Post, error)

	OnFeed func(
		userID int64,
		options data.ListOptions,
		sort data.PostSortMethod) ([]data.Post, error)
}

func (store mockPostStore) Create(post *data.Post) (int64, error) {
	return store.OnCreate(post)
}

func (store mockPostStore) Get(id int64) (*data.Post, error) {
	return store.OnGet(id)
}

func (store mockPostStore) Update(id int64, contents []byte) error {
	return store.OnUpdate(id, contents)
}

func (store mockPostStore) Delete(id int64) error {
	return store.OnDelete(id)
}

func (store mockPostStore) UserPosts(
	userID int64,
	options data.ListOptions,
	sort data.PostSortMethod) ([]data.Post, error) {
	return store.OnUserPosts(userID, options, sort)
}

func (store mockPostStore) Feed(
	userID int64,
	options data.ListOptions,
	sort data.PostSortMethod) ([]data.Post, error) {
	return store.OnFeed(userID, options, sort)
}

/* ************** *
 * Mock Auth Impl *
 * ************** */

type mockAuth struct {
	OnSecurePassword      func(passowrd string) (string, error)
	OnCheckPassword       func(password, storedPassword string) bool
	OnGenerateAccessToken func(user *data.User, signingKey []byte) (string, error)
}

func (auth mockAuth) SecurePassword(password string) (string, error) {
	return auth.OnSecurePassword(password)
}

func (auth mockAuth) CheckPassword(password, storedPassword string) bool {
	return auth.OnCheckPassword(password, storedPassword)
}

func (auth mockAuth) GenerateAccessToken(user *data.User, signingKey []byte) (string, error) {
	return auth.OnGenerateAccessToken(user, signingKey)
}
