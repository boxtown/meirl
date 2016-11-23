package api

// FollowUserRequest is the model for a follow user POST
// API request
type FollowUserRequest struct {
	FollowerID int64 `json:"followerID"`
	FolloweeID int64 `json:"followeeID"`
}

// TokenResponse is the model for an Access Token response
type TokenResponse struct {
	AccessToken string `json:"accessToken"`
}
