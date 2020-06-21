package models

import "time"
//easyjson:json
type Post struct {
	Author   string    `json:"author"`
	Created  time.Time `json:"created"`
	Forum    string    `json:"forum"`
	ID       int64     `json:"id"`
	IsEdited bool      `json:"isEdited"`
	Message  string    `json:"message"`
	Parent   int64     `json:"parent"`
	Thread   int64     `json:"thread"`
	ForumID  int64     `json:"-"`
	AuthorID int64     `json:"-"`
}
//easyjson:json
type PostFull struct {
	Author *User   `json:"author"`
	Forum  *Forum  `json:"forum"`
	Post   *Post   `json:"post"`
	Thread *Thread `json:"thread"`
}
//easyjson:json
type PostParams struct {
	Limit int
	Since int
	Desc  bool
	Sort  int
}

//easyjson:json
type Posts []*Post

const (
	Flat = iota
	Tree
	ParentTree
)
