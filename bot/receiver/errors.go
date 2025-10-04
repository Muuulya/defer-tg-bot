package receiver

import "errors"

var (
	ErrorUserIDNotFound   = errors.New("user id not found in update")
	ErrorUserNameNotFound = errors.New("user name not found in update")
)
