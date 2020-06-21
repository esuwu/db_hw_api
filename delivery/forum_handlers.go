package delivery

import (
	"bytes"
	"github.com/jackc/pgx"
	"github.com/valyala/fasthttp"
	"log"
	models "main/models"
	"net/http"
	"strings"
)



func (handlers *Handlers) CreateForum(ctx *fasthttp.RequestCtx) {
	newForum := models.Forum{}
	newForum.UnmarshalJSON(ctx.PostBody())

	forum, err := handlers.usecases.CreateForum(&newForum)

	var response []byte

	switch err {
	case nil:
		ctx.SetStatusCode(http.StatusCreated)
		response, _ = forum.MarshalJSON()

	case models.ForumAlreadyExists:
		ctx.SetStatusCode(	http.StatusConflict)
		response, _ = forum.MarshalJSON()

	case models.UserNotFound:
		ctx.SetStatusCode(http.StatusNotFound)
		response = models.ErrorMessage
	}
	ctx.SetContentType("application/json")
	ctx.Write(response)
}

func (handlers *Handlers) GetForumDetails(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug")
	var response []byte


	forum, err := handlers.usecases.GetForumDetails(slug)

		switch err {
		case nil:
			response, _ = forum.MarshalJSON()
		default:
			ctx.SetStatusCode(404)
			response = models.ErrorMessage
		}

	ctx.SetContentType("application/json")
	ctx.Write(response)
}

func (handlers *Handlers) GetForumThreads(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	slug := ctx.UserValue("slug")
	limit := ctx.QueryArgs().Peek("limit")
	since := ctx.QueryArgs().Peek("since")
	desc := ctx.QueryArgs().Peek("desc")

	threadArr, error := handlers.usecases.GetForumThreads(&slug, limit, since, desc)

	var response []byte
	switch error {
	case nil:
		if len(*threadArr) == 0 {
			ctx.Write([]byte("[]"))
			return
		}
		response, _ = threadArr.MarshalJSON()
	case models.ForumNotFound:
		ctx.SetStatusCode(http.StatusNotFound)
		response = models.ErrorMessage
	}

	ctx.Write(response)
}

func (handlers *Handlers) GetForumUsers(ctx *fasthttp.RequestCtx) {

	slug := ctx.UserValue("slug")
	limit := ctx.QueryArgs().Peek("limit")
	since := ctx.QueryArgs().Peek("since")
	desc := ctx.QueryArgs().Peek("desc")

	users, err := handlers.usecases.GetForumUsers(&slug, limit, since, desc)

	var response []byte

	switch err {
	case nil:
		ctx.SetStatusCode(http.StatusOK)
		if len(*users) != 0 {
			response, _ = users.MarshalJSON()
		} else {
			response = []byte("[]")
		}
	case models.ForumNotFound:
		ctx.SetStatusCode(http.StatusNotFound)
		response = models.ErrorMessage
	}

	ctx.SetContentType("application/json")
	ctx.Write(response)
}