package repository

import (
	"bytes"
	"context"
	"fmt"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/jackc/pgx"
	"log"
	"main/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)


type forumUser struct {
	userNickname *string
	userEmail    *string
	userAbout    *string
	userFullname *string
}
type forumUserArr []forumUser

func (fu forumUserArr) Len() int {
	return len(fu)
}
func (fu forumUserArr) Swap(i, j int) {
	fu[i], fu[j] = fu[j], fu[i]
}
func (fu forumUserArr) Less(i, j int) bool {
	return *(fu[i].userNickname) < *(fu[j].userNickname)
}


func StringsCompare(a, b interface{}) int {
	return strings.Compare(a.(string), b.(string))
}

//type scanner interface {
//	Scan(dst ...interface{}) error
//}

//func scanPosts(scanner Scanner, postDst *models.Post) error {
//	err := scanner.Scan(
//		&postDst.Created,
//	)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//
//func getPosts(batchResults pgx.BatchResults, batchLen int) (models.Posts, error) {
//	insertedPosts := make(models.Post, batchLen)
//
//	for i := 0; i < batchLen; i++ {
//		if err := scanPosts(batchResults.QueryRow(), &insertedPosts[i]); err != nil {
//			return nil, err
//		}
//	}
//
//	return insertedPosts, nil
//}

func (store *DBStore) CreatePosts(timer time.Time, slugOrID interface{}, postsArr *models.PostArr) (*models.PostArr, error) {
	tx := TxBegin(store)
	defer tx.Rollback()

	batch := tx.BeginBatch()
	created := time.Unix(0, 0)

	var err error
	var forumID, threadID int
	var forumSlug string

	//Claiming thread ID
	threadID, err = strconv.Atoi(slugOrID.(string))
	if err != nil {
		if err = tx.QueryRow("SELECT id, forum_slug::TEXT FROM thread WHERE slug=$1", slugOrID).Scan(&threadID, &forumSlug); err != nil {
			log.Println(err)
			return nil, models.ThreadNotFound
		}
	} else {
		if err = tx.QueryRow("SELECT id, forum_slug::TEXT FROM thread WHERE id=$1", threadID).Scan(&threadID, &forumSlug); err != nil {
			log.Println(err)
			return nil, models.ThreadNotFound
		}
	}

	if len(*postsArr) == 0 {
		return nil, nil
	}

	if err = tx.QueryRow("getForumIDBySlug", &forumSlug).Scan(&forumID); err != nil {
		log.Fatalln(err)
	}

	ids := make([]int64, 0, len(*postsArr))
	if err = tx.QueryRow("SELECT array_agg(nextval('post_id_seq')::BIGINT) FROM generate_series(1,$1)", len(*postsArr)).Scan(&ids); err != nil {
		log.Fatalln(err)
	}

	var postsWaitingParents []int
	userNicknameSet := treeset.NewWith(StringsCompare)

	for i, post := range *postsArr {
		userNicknameSet.Add(strings.ToLower(post.User_nick))

		if post.Parent != 0 {
			postsWaitingParents = append(postsWaitingParents, i)
			batch.Queue("selectParentAndParents", []interface{}{int(post.Parent)}, nil, nil)
		}
	}

	userNicknameOrderedSet := userNicknameSet.Values()

	for _, userNickname := range userNicknameOrderedSet {
		batch.Queue("getUserProfileQuery", []interface{}{userNickname}, nil, nil)
	}

	var parentThreadID int64
	if err = batch.Send(context.Background(), nil);
		err != nil {
		log.Fatalln(err)
	}

	for _, postIdx := range postsWaitingParents {
		if err = batch.QueryRowResults().
			Scan(&parentThreadID, &(*postsArr)[postIdx].Parents);
			err != nil {
			return nil, models.PostsConflict
		}
		if parentThreadID != 0 && parentThreadID != int64(threadID) {
			return nil, models.PostsConflict
		}
	}

	userRealNicknameMap := make(map[string]string)
	var userModelsOrderedSet models.UsersArr

	for _, userNickname := range userNicknameOrderedSet {
		user := models.User{}
		if err = batch.QueryRowResults().
			Scan(&user.Nickname, &user.Email, &user.About, &user.Fullname);
			err != nil {
			return nil, models.UserNotFound
		}
		userModelsOrderedSet = append(userModelsOrderedSet, &user)
		userRealNicknameMap[userNickname.(string)] = user.Nickname
	}


	for index, post := range *postsArr {
		post.Id = int(ids[index])
		post.Thread_id = threadID
		post.Forum_slug = forumSlug
		post.Created = created
		fmt.Print(post.Created)
		post.User_nick = userRealNicknameMap[strings.ToLower(post.User_nick)]
		post.Parents = append(post.Parents, int32(ids[index]))

		batch.Queue("insertIntoPost", []interface{}{post.Id, post.User_nick, post.Message, post.Created, post.Forum_slug, post.Thread_id, post.Parent, post.Parents, post.Parents[0]}, nil, nil)
	}

	//for _, post := range *postsArr {
	//	batch.QueryRowResults().Scan(&post.Created)
	//}

	for _, user := range userModelsOrderedSet {
		batch.Queue("insertIntoForumUsers", []interface{}{forumID, user.Nickname, user.Email, user.About, user.Fullname}, nil, nil)
	}

	if err = batch.Send(context.Background(), nil);
		err != nil {
		log.Fatalln(err)
	}

	//for _, post := range *postsArr {
	//	if err := batch.QueryRowResults().Scan(&post.Created); err != nil {
	//		log.Print(err)
	//	}
	//}
	for range *postsArr {
		if _, err := batch.ExecResults(); err != nil {
			log.Fatalln(err)
		}
	}

	for range userModelsOrderedSet {
		if _, err := batch.ExecResults(); err != nil {
			log.Fatalln(err)
		}
	}

	_, err = tx.Exec(`UPDATE forum SET posts=posts+$2 WHERE slug=$1`, forumSlug, len(*postsArr))
	if err != nil {
		log.Fatalln(err)
	}

	if err = tx.Commit(); err != nil {
		log.Fatalln(err)
	}

	tx.Commit()
	return postsArr, nil
}

func (store *DBStore) PutVote(slugOrID interface{}, vote *models.Vote) (*models.Thread, error) {
	tx, err := store.DB.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	_, err = strconv.Atoi(slugOrID.(string))

	thread := models.Thread{}

	if err != nil {
		err = tx.QueryRow("putVoteByThrSLUG", vote.Nickname, slugOrID, vote.Voice).Scan(&thread.Id, &thread.Slug, &thread.Title, &thread.Message, &thread.Forum_slug, &thread.User_nick, &thread.Created, &thread.Votes_count)
	} else {
		err = tx.QueryRow("putVoteByThrID", vote.Nickname, slugOrID, vote.Voice).Scan(&thread.Id, &thread.Slug, &thread.Title, &thread.Message, &thread.Forum_slug, &thread.User_nick, &thread.Created, &thread.Votes_count)
	}

	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return &thread, nil
}

func (store *DBStore) GetThread(slugOrID interface{}) (*models.Thread, error) {
	thread := models.Thread{}

	_, err := strconv.Atoi(slugOrID.(string))

	if err != nil {
		err = store.DB.QueryRow("getThreadBySlug", slugOrID).
			Scan(&thread.Id, &thread.Slug, &thread.Title, &thread.Message, &thread.Forum_slug, &thread.User_nick, &thread.Created, &thread.Votes_count)
		return &thread, err
	}

	err = store.DB.QueryRow("getThreadById", slugOrID).Scan(&thread.Id, &thread.Slug, &thread.Title, &thread.Message, &thread.Forum_slug, &thread.User_nick, &thread.Created, &thread.Votes_count)
	return &thread, err
}

const threadUpdateQuery = `UPDATE thread
SET message = coalesce($1, message),
	title = coalesce($2,title)
WHERE id = $3
RETURNING  id,
	slug::TEXT,
	title,
	message,
	forum_slug::TEXT,
	user_nick::TEXT,
	created,
	votes_count `

func (store *DBStore) UpdateThreadDetails(slugOrID *string, thrUpdate *models.ThreadUpdate) (*models.Thread, int) {
	tx, err := store.DB.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	var ID int
	var fs string
	if ID, err = strconv.Atoi(*slugOrID); err != nil {
		if err = tx.QueryRow("SELECT id, forum_slug::TEXT FROM thread WHERE slug=$1", slugOrID).Scan(&ID, &fs);
			err != nil {
			return nil, http.StatusNotFound
		}
	}

	var thread models.Thread

	if err = tx.QueryRow(threadUpdateQuery, thrUpdate.Message, thrUpdate.Title, ID).
		Scan(&thread.Id, &thread.Slug, &thread.Title, &thread.Message, &thread.Forum_slug,
		&thread.User_nick, &thread.Created, &thread.Votes_count);
		err != nil {
		return nil, http.StatusNotFound
	}
	return &thread, http.StatusOK
}

func getThreadPostsFlat(store *DBStore, ID int, limit []byte, since []byte, desc []byte) (*models.PostArr, int) {
	var err error
	var rows *pgx.Rows

	if since != nil {
		if limit != nil {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = store.DB.Query("getPostsFlatSinceLimitDesc", ID, limit, since)
			} else {
				rows, err = store.DB.Query("getPostsFlatSinceLimit", ID, limit, since)
			}
		} else {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = store.DB.Query("getPostsFlatSinceLimitDesc", ID, nil, since)
			} else {
				rows, err = store.DB.Query("getPostsFlatSinceLimit", ID, nil, since)
			}
		}
	} else {
		if limit != nil {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = store.DB.Query("getPostsFlatLimitDesc", ID, limit)
			} else {
				rows, err = store.DB.Query("getPostsFlatLimit", ID, limit)
			}
		} else {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = store.DB.Query("getPostsFlatLimitDesc", ID, nil)
			} else {
				rows, err = store.DB.Query("getPostsFlatLimit", ID, nil)
			}
		}
	}

	if err != nil {
		log.Fatalln(err)
	}

	var posts models.PostArr

	for rows.Next() {
		post := models.Post{}

		if err = rows.Scan(&post.Id, &post.User_nick, &post.Message,
			&post.Created, &post.Forum_slug, &post.Thread_id,
			&post.Is_edited, &post.Parent);
			err != nil {
			log.Fatalln(err)
		}
		posts = append(posts, &post)
	}
	rows.Close()

	return &posts, http.StatusOK
}

func getThreadPostsTree(store *DBStore, ID int, limit []byte, since []byte, desc []byte) (*models.PostArr, int) {
	var err error
	var rows *pgx.Rows

	if since != nil {
		if limit != nil {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = store.DB.Query("getPostsTreeSinceLimitDesc", ID, limit, since)
			} else {
				rows, err = store.DB.Query("getPostsTreeSinceLimit", ID, limit, since)
			}
		} else {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = store.DB.Query("getPostsTreeSinceLimitDesc", ID, nil, since)
			} else {
				rows, err = store.DB.Query("getPostsTreeSinceLimit", ID, nil, since)
			}
		}
	} else {
		if limit != nil {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = store.DB.Query("getPostsTreeLimitDesc", ID, limit)
			} else {
				rows, err = store.DB.Query("getPostsTreeLimit", ID, limit)
			}
		} else {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = store.DB.Query("getPostsTreeLimitDesc", ID, nil)
			} else {
				rows, err = store.DB.Query("getPostsTreeLimit", ID, nil)
			}
		}
	}

	if err != nil {
		log.Fatalln(err)
	}

	var posts models.PostArr

	for rows.Next() {
		post := models.Post{}

		if err = rows.Scan(&post.Id, &post.User_nick, &post.Message,
			&post.Created, &post.Forum_slug, &post.Thread_id,
			&post.Is_edited, &post.Parent);
			err != nil {
			log.Fatalln(err)
		}
		posts = append(posts, &post)
	}
	rows.Close()

	return &posts, http.StatusOK
}
/*
post.ParentTreeSort: `
			WITH roots AS (
				SELECT DISTINCT path[1]
				FROM posts
				WHERE thread_id = $1
				ORDER BY path[1] DESC
				LIMIT $2
			)
			SELECT id,
				   thread_id,
				   author_nickname,
				   forum_slug,
				   is_edited,
				   message,
				   parent,
				   created
			FROM posts
			WHERE thread_id = $1
			  AND path[1] IN (SELECT * FROM roots)
			ORDER BY path[1] DESC, path[2:]`,
*/

func getThreadPostsParentTree(store *DBStore, ID int, limit []byte, since []byte, desc []byte) (*models.PostArr, int) {
	var err error
	var rows *pgx.Rows

	if since != nil {
		if limit != nil {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = store.DB.Query("getPostsParentTreeSinceLimitDesc", ID, limit, since)
			} else {
				rows, err = store.DB.Query("getPostsParentTreeSinceLimit", ID, limit, since)
			}
		} else {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = store.DB.Query("getPostsParentTreeSinceLimitDesc", ID, nil, since)
			} else {
				rows, err = store.DB.Query("getPostsParentTreeSinceLimit", ID, nil, since)
			}
		}
	} else {
		if limit != nil {
			if bytes.Equal(desc, []byte("true")) {
	/*men*/			rows, err = store.DB.Query("getPostsParentTreeLimitDesc", ID, limit)
			} else {
				rows, err = store.DB.Query("getPostsParentTreeLimit", ID, limit)
			}
		} else {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = store.DB.Query("getPostsParentTreeLimitDesc", ID, nil)
			} else {
				rows, err = store.DB.Query("getPostsParentTreeLimit", ID, nil)
			}
		}
	}

	if err != nil {
		log.Fatalln(err)
	}

	var posts models.PostArr

	for rows.Next() {
		post := models.Post{}

		if err = rows.Scan(&post.Id, &post.User_nick, &post.Message,
			&post.Created, &post.Forum_slug, &post.Thread_id,
			&post.Is_edited, &post.Parent);
			err != nil {
			log.Fatalln(err)
		}
		posts = append(posts, &post)
	}
	rows.Close()

	return &posts, http.StatusOK
}

func (store *DBStore) GetThreadPosts(slugOrID *string, limit []byte, since []byte, sort []byte, desc []byte) (*models.PostArr, int) {
	var ID int
	var err error

	if _, err = strconv.Atoi(*slugOrID); err != nil {
		if err = store.DB.QueryRow("checkThreadIdBySlug", slugOrID).Scan(&ID); err != nil {
			return nil, http.StatusNotFound
		}
	} else {
		if err = store.DB.QueryRow("checkThreadIdById", slugOrID).Scan(&ID); err != nil {
			return nil, http.StatusNotFound
		}
	}

	switch true {
	case bytes.Equal([]byte("tree"), sort):
		postsTree, status := getThreadPostsTree(store, ID, limit, since, desc)
		return postsTree, status
	case bytes.Equal([]byte("parent_tree"), sort):
		postsParentTree, status := getThreadPostsParentTree(store, ID, limit, since, desc)
		return postsParentTree, status
	default:
		PostsFlat, status := getThreadPostsFlat(store, ID, limit, since, desc)
		return PostsFlat, status
	}
}
