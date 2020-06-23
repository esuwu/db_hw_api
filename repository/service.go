package repository

import (
	"log"
	"main/models"
)


func (store *DBStore) GetStatus() *models.Status {
	tx, err := store.DB.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	status := models.Status{}
	tx.QueryRow("SELECT (SELECT count(*) FROM forum) as forum, (SELECT count(*) FROM post) as post, (SELECT count(*) FROM users) as user, (SELECT count(*) FROM thread) as thread").Scan(&status.Forum, &status.Post, &status.User, &status.Thread)

	return &status
}


func (store *DBStore) Clear() {
	tx, err := store.DB.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	tx.Exec("TRUNCATE forum_users, vote, post, thread, forum, users RESTART IDENTITY CASCADE")
}
