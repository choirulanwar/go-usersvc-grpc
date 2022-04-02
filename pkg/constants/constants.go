package constants

import "github.com/pkg/errors"

var (
	ErrNotFound = errors.New("Not found")
	ErrInvalid  = errors.New("Invalid")
	ErrExists   = errors.New("Already exists")
)
