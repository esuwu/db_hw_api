package delivery

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	models "main/models"
	"net/http"
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
		response, _ = json.Marshal(err)

	}
	ctx.SetContentType("application/json")
	ctx.Write(response)
}

func (handlers *Handlers) GetForum(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug")
	var response []byte


	forum, err := handlers.usecases.GetForumDetails(slug)

		switch err {
		case nil:
			response, _ = forum.MarshalJSON()
		default:
			ctx.SetStatusCode(http.StatusNotFound)
			response, _ = json.Marshal(err)
		}

	ctx.SetContentType("application/json")
	ctx.Write(response)
}

func (handlers *Handlers) GetThreads(ctx *fasthttp.RequestCtx) {
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
		response, _ = json.Marshal(error)
		ctx.SetContentType("application/json")
	}

	ctx.Write(response)
}

func (handlers *Handlers) GetUsers(ctx *fasthttp.RequestCtx) {

	slugKey := ctx.UserValue("slug")
	limitKey := ctx.QueryArgs().Peek("limit")
	sinceKey := ctx.QueryArgs().Peek("since")
	descKey := ctx.QueryArgs().Peek("desc")

	users, err := handlers.usecases.GetForumUsers(&slugKey, limitKey, sinceKey, descKey)

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
		response, _ = json.Marshal(err)
		ctx.SetContentType("application/json")
	}

	ctx.SetContentType("application/json")
	ctx.Write(response)
}