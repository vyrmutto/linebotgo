package auth

import "errors"

var (
	ErrTokenExpired = errors.New("auth token expired")
	ErrInvalidToken = errors.New("auth token is invalid")
	ErrRateLimited  = errors.New("rate limited by LINE")
	ErrLoginBlocked = errors.New("login blocked by LINE")
	ErrLoginFailed  = errors.New("login failed: invalid credentials")
)
