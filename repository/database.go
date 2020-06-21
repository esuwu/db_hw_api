package repository

import (
	"log"

	"github.com/jackc/pgx"
	"main/prepareStat"
	"time"
)

var db *pgx.ConnPool

func TxMustBegin() *pgx.Tx {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	return tx
}

func Vaccuum() {
	time.Sleep(20 * time.Second)
	//go func() {
	db.Exec("CLUSTER forum_users USING forum_users_forum_id_nickname_index")
	db.Exec("CLUSTER users USING users_nickname_index")
	db.Exec("CLUSTER post USING parent_tree_3_1")
	db.Exec("CLUSTER thread USING thread_forum_id_created_index")
	db.Exec("CLUSTER forum USING forum_slug_id_index")
	db.Exec("VACUUM ANALYZE")
	//}()
}

// InitDBSchema initializes tables, indexes, etc.
func InitDBSchema() {

	resources.PrepareForumQueries(db)
	resources.PrepareForumUsersQueries(db)
	resources.PreparePostQueries(db)
	resources.PrepateThreadQueries(db)
	resources.PrepareUsersQueries(db)
	resources.PrepareVotesQureies(db)
}


