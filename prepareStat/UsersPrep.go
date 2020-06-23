package prepareStat

import (
	"log"

	"github.com/jackc/pgx"
)


const selectUsrByNickOrEmailQuery = `SELECT nickname::TEXT, email::TEXT, about, fullname FROM users WHERE nickname=$1 OR email=$2`

const getUserProfileQuery = `SELECT nickname::TEXT, email::TEXT, about, fullname FROM users WHERE nickname = $1`

func PrepareUsers(tx *pgx.ConnPool) {
	if _, err := tx.Prepare("selectUsrByNickOrEmailQuery", selectUsrByNickOrEmailQuery); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("getUserProfileQuery", getUserProfileQuery); err != nil {
		log.Fatalln(err)
	}
}
