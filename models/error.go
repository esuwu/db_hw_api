package models

import "errors"

//easyjson:json
type Error struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}
var ErrorMessage []byte

var UserNotFound = errors.New("UserNotFound")
var UserAlreadyExists = errors.New("UserAlreadyExists")
var ConflictOnUsers = errors.New("ConflictOnUsers")

var ForumNotFound = errors.New("ForumNotFound")
var ForumAlreadyExists = errors.New("ForumAlreadyExists")

var ThreadAlreadyExists = errors.New("ThreadAlreadyExists")
var ThreadNotFound = errors.New("ThreadNotFound")

var PostsConflict = errors.New("ConflictOnPosts")

//easyjson:json
type ErrorString struct {
	Message string `json:"message"`
}

