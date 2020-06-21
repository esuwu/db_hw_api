package delivery

import (
	"fmt"
	"github.com/valyala/fasthttp"
	models "main/models"
	"net/http"
)


// POST /user/{nickname}/create
func (handlers *Handlers) CreateForum(ctx *fasthttp.RequestCtx) {
	var newForum models.Forum

	err := newForum.UnmarshalJSON(ctx.PostBody())

	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	forum, e := handlers.usecases.PutForum(&newForum)
	var body []byte
	if e != nil {
		if e.Code == http.StatusConflict {
			body, _ = forum.MarshalJSON()
			ctx.SetStatusCode(e.Code)
			ctx.SetContentType("application/json")
			ctx.Write(body)
			return
		}
		body, _ = e.MarshalJSON()
		ctx.SetStatusCode(e.Code)
		ctx.SetContentType("application/json")
		ctx.Write(body)
		return
	}

	body, _ = forum.MarshalJSON()
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusCreated)
	ctx.Write(body)
}
// GET /forum/{slug}/details
func (handlers *Handlers) GetForum(ctx *fasthttp.RequestCtx) {

	slug := ctx.UserValue("slug")
	slugStr := fmt.Sprintf("%v", slug)
	forum, err := handlers.usecases.GetForumBySlug(slugStr)
	if err != nil {
		body, _ := err.MarshalJSON()
		ctx.SetStatusCode(err.Code)
		ctx.SetContentType("application/json")
		ctx.Write(body)
		return
	}

	body, _ := forum.MarshalJSON()
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetContentType("application/json")
	ctx.Write(body)
}
