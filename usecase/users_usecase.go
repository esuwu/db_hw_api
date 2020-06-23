package useCase

import "main/models"

func (u *useCase) PutVote(slugOrID interface{}, vote *models.Vote) (*models.Thread, error){
	thread, err := u.repository.PutVote(slugOrID, vote)
	return thread, err
}
func (u *useCase) CreateUser(user *models.User, nickname interface{}) (*models.UsersArr, error){
	users, err := u.repository.CreateUser(user, nickname)
	return users, err
}
func (u *useCase) GetUserProfile(nickname interface{}) (*models.User, error){
	user, err := u.repository.GetUserProfile(nickname)
	return user, err
}
func (u *useCase) UpdateUserProfile(userUpd *models.UserUpd, nickname interface{}) (*models.User, error){
	user, err := u.repository.UpdateUserProfile(userUpd, nickname)
	return user, err
}