package delivery

import (
	"fmt"
	"github.com/valyala/fasthttp"
	models "main/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (handlers *Handlers) CreatePost(ctx *fasthttp.RequestCtx) {
	var posts models.Posts
	var tempID int
	id := -1

	err := posts.UnmarshalJSON(ctx.PostBody())

	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetContentType("application/json")
		return
	}

	slugOrId := ctx.UserValue("slug_or_id")
	slugOrIdstr := fmt.Sprintf("%v", slugOrId)

	tempID, err = strconv.Atoi(slugOrIdstr)
	if err == nil {
		id = tempID
	}

	postsAdded := make(models.Posts, len(posts))
	var e *models.Error
	created := time.Now()

	for i, _ := range posts {
		posts[i].Created = created
		if id != -1 {
			posts[i].Thread = int64(id)
			postsAdded[i], e = handlers.usecases.PutPost(posts[i])
		} else {
			postsAdded[i], e = handlers.usecases.PutPostWithSlug(posts[i], slugOrIdstr)
		}
		if e != nil {
			body, _ := e.MarshalJSON()
			ctx.SetStatusCode(e.Code)
			ctx.SetContentType("application/json")
			ctx.Write(body)
			return
		}
	}

	if len(posts) == 0 {
		//var thread models.Thread
		if id != -1 {
			_, e = handlers.usecases.GetThreadByID(int64(id))
		} else {
			_, e = handlers.usecases.GetThreadBySlug(slugOrIdstr)
		}
		if e != nil {
			body, _ := e.MarshalJSON()
			ctx.SetStatusCode(e.Code)
			ctx.SetContentType("application/json")
			ctx.Write(body)
			return
		}
	}
	body, _ := postsAdded.MarshalJSON()
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusCreated)
	ctx.Write(body)
}
// POST /post/{id}/details
func (handlers *Handlers) UpdatePost(ctx *fasthttp.RequestCtx) {
	var setPost, post models.Post
	var e *models.Error

	err := setPost.UnmarshalJSON(ctx.PostBody())

	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}


	idInterface := ctx.UserValue("id")
	idStr := fmt.Sprintf("%v", idInterface)
	id, err := strconv.Atoi(idStr)

	setPost.ID = int64(id)

	post, e = handlers.usecases.ChangePost(&setPost)
	if e != nil {
		body, _ := e.MarshalJSON()
		ctx.SetStatusCode(e.Code)
		ctx.SetContentType("application/json")
		ctx.Write(body)
		return
	}

	body, _ := post.MarshalJSON()
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	ctx.Write(body)
}
// GET /post/{id}/details
func (handlers *Handlers) GetPostFull(ctx *fasthttp.RequestCtx) {


	idInterface := ctx.UserValue("id")
	idStr := fmt.Sprintf("%v", idInterface)

	id, _ := strconv.Atoi(idStr)

	key := ctx.QueryArgs().Peek("related")
	fields := strings.Split(string(key), ",")

	postFull, err := handlers.usecases.GetPostFull(int64(id), fields)
	if err != nil {
		body, _ := err.MarshalJSON()
		ctx.SetStatusCode(err.Code)
		ctx.SetContentType("application/json")
		ctx.Write(body)
		return
	}

	body, _ := postFull.MarshalJSON()
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	ctx.Write(body)
}

func (handlers *Handlers) GetPosts(ctx *fasthttp.RequestCtx) {
	var posts models.Posts
	var e *models.Error


	slug_or_idInterface := ctx.UserValue("slug_or_id")
	slug_or_id := fmt.Sprintf("%v", slug_or_idInterface)

	getLimit := ctx.QueryArgs().Peek("limit")
	getSince := ctx.QueryArgs().Peek("since")
	getSort := ctx.QueryArgs().Peek("sort")
	getDesc := ctx.QueryArgs().Peek("desc")

	var params models.PostParams
	var err error

	params.Limit, err = strconv.Atoi(string(getLimit))

	if err != nil {
		params.Limit = -1
	}
	params.Since, err = strconv.Atoi(string(getSince))
	if err != nil {
		params.Since = -1
	}
	fmt.Println("SINCE: ", params.Since)
	params.Desc = string(getDesc) == "true"

	switch string(getSort) {
	case "flat":
		params.Sort = models.Flat
	case "tree":
		params.Sort = models.Tree
	case "parent_tree":
		params.Sort = models.ParentTree
	}

	if id, err := strconv.Atoi(slug_or_id); err == nil {
		posts, e = handlers.usecases.GetPostsByThreadID(int64(id), params)
	} else {
		posts, e = handlers.usecases.GetPostsByThreadSlug(slug_or_id, params)
	}
	if e != nil {
		body, _ := e.MarshalJSON()
		ctx.SetStatusCode(e.Code)
		ctx.SetContentType("application/json")
		ctx.Write(body)
		return
	}


	body, _ := posts.MarshalJSON()
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	ctx.Write(body)
}

func (handlers *Handlers) Vote(ctx *fasthttp.RequestCtx) {
	var vote models.Vote

	err := vote.UnmarshalJSON(ctx.PostBody())

	if err != nil {

	}


	slug_or_idInterface := ctx.UserValue("slug_or_id")
	slug_or_id := fmt.Sprintf("%v", slug_or_idInterface)


	var thread models.Thread
	var e *models.Error

	if id, err := strconv.Atoi(slug_or_id); err == nil {
		vote.ThreadID = int64(id)
		thread, e = handlers.usecases.PutVote(&vote)
	} else {
		thread, e = handlers.usecases.PutVoteWithSlug(&vote, slug_or_id)
	}
	if e != nil {
		body, _ := e.MarshalJSON()
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetContentType("application/json")
		ctx.Write(body)
		return
	}

	body, _ := thread.MarshalJSON()
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	ctx.Write(body)
}
