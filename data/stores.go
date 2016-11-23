package data

// ErrNoEnt is returned by stores when a desired entity could
// not be found
var ErrNoEnt = Error{Message: "Entity not found"}

// Error is a custom error type for the data layer
type Error struct {
	Message string
	Cause   error
}

// NewError wraps the error inside an Error type.
func NewError(err error) *Error {
	return &Error{Message: err.Error(), Cause: err}
}

func (e Error) Error() string {
	return e.Message
}

// ListOptions is a collection of options
// for listing entities
type ListOptions struct {
	Marker interface{}
	Offset int
	Limit  int
	Desc   bool
}

// UserSortMethod is the method of sorting
// for user lists
type UserSortMethod int

const (
	// UserSortByID designates a user sort by user id
	UserSortByID UserSortMethod = iota

	// UserSortByUsername designates a user sort by username
	UserSortByUsername

	// UserSortByEmail designates a user sort by email
	UserSortByEmail

	// UserSortByActualName designates a user sort by actual name
	UserSortByActualName

	// UserSortByDOB designates a user sort by date of birth
	UserSortByDOB
)

// PostSortMethod is the method of sorting
// for post lists
type PostSortMethod int

const (
	// PostSortByDate designates a post sort by creation date
	PostSortByDate PostSortMethod = iota

	// PostSortByKeks designates a post sort by keks
	PostSortByKeks

	// PostSortByNos designates a post sort by nos
	PostSortByNos
)

// Stores is a collection of all data
// stores
type Stores struct {
	UserStore
	PostStore
}

// UserStore represents a common gateway for
// user data stores
type UserStore interface {
	Create(user *User) (int64, error)
	Get(id int64) (*User, error)
	GetByUsername(username string) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(id int64, user *User) error
	Delete(id int64) error
	Follow(followerID, followeeID int64) error
	UnFollow(followerID, followeeID int64) error
	Followers(id int64, options ListOptions, sort UserSortMethod) ([]User, error)
	Following(id int64, options ListOptions, sort UserSortMethod) ([]User, error)
}

// PostStore represents a common gateway for
// post data stores
type PostStore interface {
	Create(post *Post) (int64, error)
	Get(id int64) (*Post, error)
	Update(id int64, contents []byte) error
	Delete(id int64) error
	UserPosts(userID int64, options ListOptions, sort PostSortMethod) ([]Post, error)
	Feed(userID int64, options ListOptions, sort PostSortMethod) ([]Post, error)
}
