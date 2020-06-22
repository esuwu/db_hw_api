package repository

import (
	"github.com/jackc/pgx"
	models "main/models"
)

type Repo interface {
	CreateForum(forum *models.Forum) (*models.Forum, error)
	GetForumDetails(slug interface{}) (*models.Forum, error)
	CreateThread(forumSlug interface{}, threadDetails *models.Thread) (*models.Thread, error)
	GetForumThreads(slug interface{}, limit []byte, since []byte, desc []byte) (*models.ThreadArr, error)
	GetForumUsers(slug interface{}, limit []byte, since []byte, desc []byte) (*models.UsersArr, error)
	GetPostDetails(id *string, related []byte) (*models.PostDetails, int)
	UpdatePostDetails(id *string, postUpd *models.PostUpdate) (*models.Post, int)
	GetStatus() *models.Status
	Clear()
	CreatePosts(slugOrID interface{}, postsArr *models.PostArr) (*models.PostArr, error)
    GetThread(slugOrID interface{}) (*models.Thread, error)
	PutVote(slugOrID interface{}, vote *models.Vote) (*models.Thread, error)
	UpdateThreadDetails(slugOrID *string, thrUpdate *models.ThreadUpdate) (*models.Thread, int)
	GetThreadPosts(slugOrID *string, limit []byte, since []byte, sort []byte, desc []byte) (*models.PostArr, int)
	CreateUser(user *models.User, nickname interface{}) (*models.UsersArr, error)
	UpdateUserProfile(newData *models.UserUpd, nickname interface{}) (*models.User, error)
	GetUserProfile(nickname interface{}) (*models.User, error)
}

type DBStore struct {
	DB *pgx.ConnPool
}

func NewDBStore(db *pgx.ConnPool) Repo {
	return &DBStore{
		db,
	}
}


