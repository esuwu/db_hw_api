package delivery

import (
	"github.com/valyala/fasthttp"
	"log"
	models "main/models"
	"net/http"
)



func (handlers *Handlers) CreateThread(ctx *fasthttp.RequestCtx) {
	threadDetails := models.Thread{}
	threadDetails.UnmarshalJSON(ctx.PostBody())

	slug := ctx.UserValue("slug")

	threadExisting, err := handlers.usecases.CreateThread(&slug, &threadDetails)

	var response []byte

	switch err {
	case nil:
		ctx.SetStatusCode(http.StatusCreated)
		response, _ = threadExisting.MarshalJSON()

	case models.UserNotFound, models.ForumNotFound:
		ctx.SetStatusCode(http.StatusNotFound)
		response = models.ErrorMessage

	case models.ThreadAlreadyExists:
		ctx.SetStatusCode(http.StatusConflict)
		response, _ = threadExisting.MarshalJSON()
	}

	ctx.SetContentType("application/json")
	ctx.Write(response)
}

func (handlers *Handlers) GetThreadDetails(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	slugOrID := ctx.UserValue("slug_or_id")

	threadDetails, err := handlers.usecases.GetThread(slugOrID)
	if err != nil {
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.Write(models.ErrorMessage)
		return
	}

	var resp []byte
	resp, err = threadDetails.MarshalJSON()
	if err != nil {
		log.Fatalln(err)
	}

	ctx.Write(resp)
}

func (handlers *Handlers) UpdateThreadDetails(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	slugOrID := ctx.UserValue("slug_or_id").(string)

	threadUpd := models.ThreadUpdate{}
	threadUpd.UnmarshalJSON(ctx.PostBody())

	thread, statusCode := handlers.usecases.UpdateThreadDetails(&slugOrID, &threadUpd)
	ctx.SetStatusCode(statusCode)

	switch statusCode {
	case http.StatusOK:
		resp, _ := thread.MarshalJSON()
		ctx.Write(resp)
	case http.StatusNotFound:
		ctx.Write(models.ErrorMessage)
	}
}

func (handlers *Handlers) GetThreadPosts(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	slugOrID := ctx.UserValue("slug_or_id").(string)
	limit := ctx.QueryArgs().Peek("limit")
	since := ctx.QueryArgs().Peek("since")
	sort := ctx.QueryArgs().Peek("sort")
	desc := ctx.QueryArgs().Peek("desc")

	postArr, statusCode := handlers.usecases.GetThreadPosts(&slugOrID, limit, since, sort, desc)

	ctx.SetStatusCode(statusCode)

	switch statusCode {
	case http.StatusOK:
		if len(*postArr) != 0 {
			response, _ := postArr.MarshalJSON()
			ctx.Write(response)
		} else {
			ctx.Write([]byte("[]"))
		}
	case http.StatusNotFound:
		ctx.Write(models.ErrorMessage)
	}
}

func (handlers *Handlers)  VoteThread(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	var vote models.Vote
	vote.UnmarshalJSON(ctx.PostBody())

	slugOrID := ctx.UserValue("slug_or_id")

	thread, err := handlers.usecases.PutVote(slugOrID, &vote)
	if err != nil {
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.Write(models.ErrorMessage)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
	response, _ := thread.MarshalJSON()
	ctx.Write(response)
}