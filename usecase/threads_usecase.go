package useCase

import "main/models"

func (u *useCase) CreateThread(slug interface{}, threadDetails *models.Thread) (*models.Thread, error){
	thread, err := u.repository.CreateThread(slug,threadDetails)
	return thread, err
}

func (u *useCase) GetThread(slugOrID interface{}) (*models.Thread, error){
	thread, err := u.repository.GetThread(slugOrID)
	return thread, err
}
func (u *useCase) UpdateThreadDetails(slugOrID *string, threadUpd *models.ThreadUpdate) (*models.Thread, int){
	thread, status := u.repository.UpdateThreadDetails(slugOrID, threadUpd)
	return thread, status
}
func (u *useCase)  GetThreadPosts(slugOrID *string, limit []byte, since []byte, sort []byte, desc []byte) (*models.PostArr, int){
	posts, status := u.repository.GetThreadPosts(slugOrID, limit, since, sort, desc)
	return posts, status
}