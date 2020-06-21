package models

import (
	"regexp"
)
//easyjson:json
type User struct {
	About    string `json:"about,omitempty"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Nickname string `json:"nickname,omitempty"`
}
//easyjson:json
type UserUpd struct {
	About    *string `json:"about,omitempty"`
	Email    *string `json:"email"`
	Fullname *string `json:"fullname"`
	Nickname *string `json:"nickname,omitempty"`
}



//easyjson:json
type UsersArr []*User

var (
	nicknameRegexp *regexp.Regexp
	emailRegexp    *regexp.Regexp
)

