package api

// TokenResponse is the model for an Access Token response
type TokenResponse struct {
	AccessToken string `json:"accessToken"`
}

// IDResponse is the model for a response containing
// the ID of a created entity
type IDResponse struct {
	ID int64 `json:"id"`
}
