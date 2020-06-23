package repository

import (
	"log"

	"github.com/jackc/pgx"
)

func TxBegin(store *DBStore) *pgx.Tx {
	tx, err := store.DB.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	return tx
}






