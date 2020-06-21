package models

import "time"
//easyjson:json
type Thread struct {
	Author   string     `json:"author"`
	Created  *time.Time `json:"created"`
	Forum    string     `json:"forum"`
	ID       int64      `json:"id"`
	Message  string     `json:"message"`
	Slug     string     `json:"slug"`
	Title    string     `json:"title"`
	Votes    int32      `json:"votes"`
	AuthorID int64      `json:"-"`
	ForumID  int64      `json:"-"`
}
//easyjson:json
type Threads []*Thread
//easyjson:json
type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
	ThreadID int64  `json:"-"`
	AuthorID int64  `json:"-"`
}
//easyjson:json
type ThreadParams struct {
	Limit int
	Since time.Time
	Desc  bool
}
