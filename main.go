package main

import (
	//"database/sql"
	"fmt"
	delivery "main/delivery"
	models "main/models"
	prepStat "main/prepareStat"
	repository "main/repository"
	useCase "main/usecase"
	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
	"log"
	"github.com/valyala/fasthttp"
	"github.com/buaazp/fasthttprouter"
)

func InitPrepStatement(db *pgx.ConnPool){

	prepStat.PrepareForum(db)
	prepStat.PrepareForumUsers(db)
	prepStat.PreparePost(db)
	prepStat.PrepateThread(db)
	prepStat.PrepareUsers(db)
	prepStat.PrepareVotes(db)
}

func RouteInit(api *delivery.Handlers) *fasthttprouter.Router {

	r := fasthttprouter.New()
	r.POST("/api/forum/:slug", api.CreateForum)
	r.POST("/api/forum/:slug/create", api.CreateThread)
	r.GET("/api/forum/:slug/details", api.GetForum)
	r.GET("/api/forum/:slug/threads", api.GetThreads)
	r.GET("/api/forum/:slug/users", api.GetUsers)

	r.POST("/api/thread/:slug_or_id/create", api.CreatePost)
	r.GET("/api/thread/:slug_or_id/details", api.GetThread)
	r.POST("/api/thread/:slug_or_id/details", api.UpdateThread)
	r.GET("/api/thread/:slug_or_id/posts", api.GetPosts)
	r.POST("/api/thread/:slug_or_id/vote", api.Vote)

	r.POST("/api/user/:nickname/create", api.CreateUser)
	r.GET("/api/user/:nickname/profile", api.GetUser)
	r.POST("/api/user/:nickname/profile", api.UpdateUser)

	r.GET("/api/post/:id/details", api.GetPostFull)
	r.POST("/api/post/:id/details", api.UpdatePost)

	r.GET("/api/service/status", api.GetStatus)
	r.POST("/api/service/clear", api.Clear)


	return r
}

func main() {
	db, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			User:     "docker",
			Password: "docker",
			Port:     5432,
			Database: "docker",
		},
		MaxConnections: 50,
	})
	InitPrepStatement(db)
	usecases := useCase.NewUseCase(repository.NewDBStore(db))
	api := delivery.NewHandlers(usecases)

	_, err = db.Exec(models.InitScript)

	if err != nil {
		fmt.Println(err)
	}

	router := RouteInit(api)

	log.Println("http server started on 5000 port: ")
	//err = http.ListenAndServe(":5000", r)
	err = fasthttp.ListenAndServe("localhost:5000", router.Handler)
	log.Println(err)
	if err != nil {
		log.Println(err)
		return
	}
}
