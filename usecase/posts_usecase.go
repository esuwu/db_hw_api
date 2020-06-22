package useCase

import "main/models"

func (u *useCase) GetPostDetails(id *string, related []byte) (*models.PostDetails, int){
	postDetails, status := u.repository.GetPostDetails(id, related)
	return postDetails, status
}
func (u *useCase) UpdatePostDetails(id *string, postUpd *models.PostUpdate) (*models.Post, int){
	post, status := u.repository.UpdatePostDetails(id, postUpd)
	return post, status
}

func (u *useCase)  CreatePosts(slugOrID interface{}, postsArr *models.PostArr) (*models.PostArr, error){
	posts, err := u.repository.CreatePosts(slugOrID, postsArr)
	return posts, err
}