package models

import "time"
//easyjson:json
type Thread struct {
	Id          *int      `json:"id"`
	Slug        *string   `json:"slug"`
	Title       string    `json:"title"`
	Message     string    `json:"message"`
	Forum_slug  string    `json:"forum"`
	User_nick   string    `json:"author"`
	Created     time.Time `json:"created,omitempty"`
	Votes_count *int      `json:"votes"`
}
//easyjson:json
type ThreadUpdate struct {
	Message *string `json:"message"`
	Title   *string `json:"title"`
}

//easyjson:json
type ThreadArr []*Thread