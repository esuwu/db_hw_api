package repository

import (
	"log"

	"github.com/jackc/pgx"
	"time"
)

func TxBegin(store *DBStore) *pgx.Tx {
	tx, err := store.DB.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	return tx
}

func Vaccuum(store *DBStore){
	time.Sleep(20 * time.Second)
	store.DB.Exec("CLUSTER forum_users USING forum_users_forum_id_nickname_index")
	store.DB.Exec("CLUSTER users USING users_nickname_index")
	store.DB.Exec("CLUSTER post USING parent_tree_3_1")
	store.DB.Exec("CLUSTER thread USING thread_forum_id_created_index")
	store.DB.Exec("CLUSTER forum USING forum_slug_id_index")
	store.DB.Exec("VACUUM ANALYZE")
}




