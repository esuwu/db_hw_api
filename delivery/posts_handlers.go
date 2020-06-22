package delivery

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	models "main/models"
	"net/http"
)
func (handlers *Handlers) CreatePost(ctx *fasthttp.RequestCtx) {
	slugOrID := ctx.UserValue("slug_or_id")

	postsArr := models.PostArr{}
	postsArr.UnmarshalJSON(ctx.PostBody())

	newPosts, err := handlers.usecases.CreatePosts(slugOrID, &postsArr)

	var response []byte

	switch err {
	case nil:
		ctx.SetStatusCode(http.StatusCreated)
		if newPosts != nil {
			response, _ = newPosts.MarshalJSON()
		} else {
			response = []byte("[]")
		}

	case models.ThreadNotFound, models.UserNotFound:
		ctx.SetStatusCode(http.StatusNotFound)
		response, _ = json.Marshal(err)

	case models.PostsConflict:
		ctx.SetStatusCode(http.StatusConflict)
		response, _ = json.Marshal(err)
	}

	ctx.SetContentType("application/json")
	ctx.Write(response)
}

func  (handlers *Handlers)  GetPostFull(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	id := ctx.UserValue("id").(string)
	related := ctx.QueryArgs().Peek("related")

	postDetails, statusCode := handlers.usecases.GetPostDetails(&id, related)
	ctx.SetStatusCode(statusCode)

	switch statusCode {
	case http.StatusOK:
		resp, _ := postDetails.MarshalJSON()
		ctx.Write(resp)
	case http.StatusNotFound:

		ctx.Write(models.ErrorMessage)
	}
}

func (handlers *Handlers) UpdatePost(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	id := ctx.UserValue("id").(string)
	postUpd := models.PostUpdate{}

	postUpd.UnmarshalJSON(ctx.PostBody())

	post, statusCode := handlers.usecases.UpdatePostDetails(&id, &postUpd)
	ctx.SetStatusCode(statusCode)

	switch statusCode {
	case http.StatusOK:
		resp, _ := post.MarshalJSON()
		ctx.Write(resp)
	case http.StatusNotFound:
		ctx.Write(models.ErrorMessage)
	}
}
