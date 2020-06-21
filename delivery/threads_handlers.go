package delivery

import (
	"fmt"
	"github.com/valyala/fasthttp"
	models "main/models"
	"net/http"
	"strconv"
	"time"
)



func (handlers *Handlers) CreateThread(ctx *fasthttp.RequestCtx) {
	var newThread models.Thread

	err := newThread.UnmarshalJSON(ctx.PostBody())

	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	slugOrId := ctx.UserValue("slug")
	slug := fmt.Sprintf("%v", slugOrId)

	fmt.Println("TIME: ", newThread.Created)

	newThread.Forum = slug

	thread, e := handlers.usecases.PutThread(&newThread)
	if e != nil {
		if e.Code == http.StatusConflict {
			body, _ := thread.MarshalJSON()
			ctx.SetStatusCode(e.Code)
			ctx.SetContentType("application/json")
			ctx.Write(body)
			return
		}
		body, _ := e.MarshalJSON()
		ctx.SetStatusCode(e.Code)
		ctx.SetContentType("application/json")
		ctx.Write(body)
		return
	}

	body, _ := thread.MarshalJSON()
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusCreated)
	ctx.Write(body)
}

func (handlers *Handlers) GetThreads(ctx *fasthttp.RequestCtx) {

	slugInterface := ctx.UserValue("slug")
	slug := fmt.Sprintf("%v", slugInterface)

	getLimit := ctx.QueryArgs().Peek("limit")
	getSince := ctx.QueryArgs().Peek("since")
	getDesc := ctx.QueryArgs().Peek("desc")

	var params models.ThreadParams
	var err error

	params.Limit, err = strconv.Atoi(string(getLimit))
	if err != nil {
		params.Limit = -1
	}

	params.Since, err = time.Parse(time.RFC3339Nano, string(getSince))
	if err != nil {
		params.Since = time.Time{}
	}
	params.Desc = string(getDesc) == "true"

	threads, e := handlers.usecases.GetThreadsByForum(slug, params)
	if e != nil {

		body, _ := e.MarshalJSON()
		ctx.SetStatusCode(e.Code)
		ctx.SetContentType("application/json")
		ctx.Write(body)
		return
	}

	body, _ := threads.MarshalJSON()
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	ctx.Write(body)
}

func (handlers *Handlers) GetThread(ctx *fasthttp.RequestCtx) {
	var thread models.Thread
	var e *models.Error

	slugInterface := ctx.UserValue("slug_or_id")
	slug_or_id := fmt.Sprintf("%v", slugInterface)

	if id, err := strconv.Atoi(slug_or_id); err == nil {
		thread, e = handlers.usecases.GetThreadByID(int64(id))
	} else {
		thread, e = handlers.usecases.GetThreadBySlug(slug_or_id)
	}
	if e != nil {
		body, _ := e.MarshalJSON()
		ctx.SetStatusCode(e.Code)
		ctx.SetContentType("application/json")
		ctx.Write(body)
		return
	}

	fmt.Println(thread)

	body, _ := thread.MarshalJSON()
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	ctx.Write(body)
}

func (handlers *Handlers) UpdateThread(ctx *fasthttp.RequestCtx) {
	var thread models.Thread
	var e *models.Error

	err := thread.UnmarshalJSON(ctx.PostBody())

	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	slugInterface := ctx.UserValue("slug_or_id")
	slug_or_id := fmt.Sprintf("%v", slugInterface)



	if id, err := strconv.Atoi(slug_or_id); err == nil {
		thread.ID = int64(id)
		thread, e = handlers.usecases.UpdateThreadWithID(&thread)
	} else {
		thread.Slug = slug_or_id
		thread, e = handlers.usecases.UpdateThreadWithSlug(&thread)
	}
	if e != nil {
		body, _ := e.MarshalJSON()
		ctx.SetStatusCode(e.Code)
		ctx.SetContentType("application/json")
		ctx.Write(body)
		return
	}


	body, _ := thread.MarshalJSON()
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	ctx.Write(body)
}
