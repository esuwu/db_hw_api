package useCase

import (
	models "main/models"
	repository "main/repository"
)

type UseCase interface {
	CreateForum(newForum *models.Forum) (*models.Forum, error)
	CreateThread(slug interface{}, threadDetails *models.Thread) (*models.Thread, error)
	GetForumDetails(slug interface{}) (*models.Forum, error)
	GetForumThreads(slug interface{}, limit []byte, since []byte, desc []byte) (*models.ThreadArr, error)
	GetForumUsers(slug interface{}, limit []byte, since []byte, desc []byte) (*models.UsersArr, error)
	GetPostDetails(id *string, related []byte) (*models.PostDetails, int)
	UpdatePostDetails(id *string, postUpd *models.PostUpdate) (*models.Post, int)
	GetStatus() *models.Status
	Clear()
	CreatePosts(slugOrID interface{}, postsArr *models.PostArr) (*models.PostArr, error)
	GetThread(slugOrID interface{}) (*models.Thread, error)
	UpdateThreadDetails(slugOrID *string, threadUpd *models.ThreadUpdate) (*models.Thread, int)
	GetThreadPosts(slugOrID *string, limit []byte, since []byte, sort []byte, desc []byte) (*models.PostArr, int)
	PutVote(slugOrID interface{}, vote *models.Vote) (*models.Thread, error)
	CreateUser(user *models.User, nickname interface{}) (*models.UsersArr, error)
	GetUserProfile(nickname interface{}) (*models.User, error)
	UpdateUserProfile(userUpd *models.UserUpd, nickname interface{}) (*models.User, error)
}

type useCase struct {
	repository repository.Repo
}

func NewUseCase(repo repository.Repo) UseCase {
	return &useCase{
		repository: repo,
	}
}


func (u *useCase) GetStatus() *models.Status {
	return u.repository.GetStatus()
}
func (u *useCase) Clear() {
	u.repository.Clear()
	return
}
