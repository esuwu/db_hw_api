package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	models "main/models"
	"net/http"
	"strconv"
)


func (handlers *Handlers) CreateUser(ctx *fasthttp.RequestCtx) {
	var user models.User
	user.UnmarshalJSON(ctx.PostBody())

	nickname := ctx.UserValue("nickname")

	existingUsers, err := database.CreateUser(&user, nickname)

	var resp []byte

	switch err {
	case nil:
		ctx.SetStatusCode(201)
		user.Nickname = ctx.UserValue("nickname").(string)
		resp, _ = user.MarshalJSON()

	case models.UserAlreadyExists:
		ctx.SetStatusCode(409)
		resp, _ = existingUsers.MarshalJSON()
	}

	ctx.SetContentType("application/json")
	ctx.Write(resp)
}