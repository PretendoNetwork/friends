package database

import (
	"errors"
)

var (
	// ErrPIDNotFound is returned if a given PID is not found in the database
	ErrPIDNotFound = errors.New("PID not found")

	// ErrFriendRequestNotFound is returned if a given friend request is not found in the database
	ErrFriendRequestNotFound = errors.New("Friend request not found")

	// ErrFriendshipNotFound is returned if a given friendship is not found in the database
	ErrFriendshipNotFound = errors.New("Friendship not found")

	// ErrBlockListNotFound is returned if a given PID does not have a blacklist
	ErrBlacklistNotFound = errors.New("Blacklist not found")

	// ErrEmptyList is returned if a given PID returned an empty list on an operation
	ErrEmptyList = errors.New("List is empty")
)
