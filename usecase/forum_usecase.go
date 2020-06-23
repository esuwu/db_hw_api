package useCase

import (
	models "main/models"
)

func (u *useCase) CreateForum(newForum *models.Forum) (*models.Forum, error){
	forum, err := u.repository.CreateForum(newForum)
	return forum, err
}

func (u *useCase) GetForumDetails(slug interface{}) (*models.Forum, error){
	forum, err := u.repository.GetForumDetails(slug)
	return forum, err
}
func (u *useCase) GetForumThreads(slug interface{}, limit []byte, since []byte, desc []byte) (*models.ThreadArr, error){
	thread, err := u.repository.GetForumThreads(slug, limit, since, desc)
	return thread, err
}
func (u *useCase) GetForumUsers(slug interface{}, limit []byte, since []byte, desc []byte) (*models.UsersArr, error){
	users, err := u.repository.GetForumUsers(slug, limit, since, desc)
	return users, err
}



