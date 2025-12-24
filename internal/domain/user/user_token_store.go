package user

type RefreshTokenStore interface {
	// Store refresh token with its expiry
	Set(userID string, refreshToken string) error

	// Check if refresh token exists and is valid (not revoked)
	Exists(userID string, refreshToken string) (bool, error)

	// Delete refresh token (on logout or rotation)
	Delete(userID string, refreshToken string) error
}
