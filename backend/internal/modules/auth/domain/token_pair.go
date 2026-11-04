package domain

// TokenPair represents the minted access and refresh token pair (ADR-014).
type TokenPair struct {
	AccessToken      string
	RefreshToken     string
	TokenType        string
	ExpiresInSeconds int
}
