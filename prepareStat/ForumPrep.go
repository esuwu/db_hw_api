package prepareStat

import (
	"log"

	"github.com/jackc/pgx"
)


const selectForumQuery = `SELECT
	slug::TEXT,
	title,
	posts,
	threads,
	moderator::TEXT
FROM forum
WHERE slug=$1`

const getForumIDBySlug = `SELECT id
FROM forum
WHERE slug=$1`

func PrepareForum(tx *pgx.ConnPool) {
	if _, err := tx.Prepare("selectForumQuery", selectForumQuery); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("getForumIDBySlug", getForumIDBySlug); err != nil {
		log.Fatalln(err)
	}
}
