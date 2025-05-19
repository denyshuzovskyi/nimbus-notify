package error

import "errors"

var (
	LocationNotFound          = errors.New("no matching location found")
	SubscriptionAlreadyExists = errors.New("subscription already exists")
	InvalidToken              = errors.New("invalid token")
	TokenNotFound             = errors.New("token not found")
	UnexpectedState           = errors.New("unexpected state")
)
