package models

import "time"
//easyjson:json
type Post struct {
	Id         int        `json:"id"`
	User_nick  string     `json:"author"`
	Message    string     `json:"message"`
	Created    time.Time `json:"created,omitempty"`
	Forum_slug string     `json:"forum"`
	Thread_id  int        `json:"thread"`
	Is_edited  bool       `json:"isEdited"`
	Parent     int32      `json:"parent,omitempty"`
	Parents    []int32
}
//easyjson:json
type PostDetails struct {
	AuthorDetails *User   `json:"author,omitempty"`
	ForumDetails  *Forum  `json:"forum,omitempty"`
	PostDetails   *Post   `json:"post,omitempty"`
	ThreadDetails *Thread `json:"thread,omitempty"`
}

//easyjson:json
type PostUpdate struct {
	Message *string `json:"message"`
}

//easyjson:json
type PostArr []*Post

//easyjson:json
type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
}

//easyjson:json
type VoteDB struct{
	ID int
	Nickname string
	Thread_id int
	Voice int
}

const (
	Flat = iota
	Tree
	ParentTree
)
