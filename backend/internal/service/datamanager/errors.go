package datamanager

import "errors"

var (
	ErrObjectNotFound = errors.New("object not found")
	ErrConflict       = errors.New("conflict: duplicate key value violates unique constraint")
	ErrDuplication    = errors.New("duplication")
)
