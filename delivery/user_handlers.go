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

	err := user.UnmarshalJSON(ctx.PostBody())

	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	nicknameInterface := ctx.UserValue("nickname")
	nickname := fmt.Sprintf("%v", nicknameInterface)


	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	user.Nickname = nickname

	users, e := handlers.usecases.PutUser(&user)
	if e != nil {
		body, _ := e.MarshalJSON()
		ctx.SetStatusCode(e.Code)
		ctx.SetContentType("application/json")
		ctx.Write(body)

		return
	}
	if users != nil {
		body, _ := users.MarshalJSON()
		ctx.SetStatusCode(http.StatusConflict)
		ctx.SetContentType("application/json")
		ctx.Write(body)
		return
	}
	body, _ := user.MarshalJSON()
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusCreated)
	ctx.Write(body)
}

func (handlers *Handlers) GetUser(ctx *fasthttp.RequestCtx) {
	nicknameInterface := ctx.UserValue("nickname")
	nickname := fmt.Sprintf("%v", nicknameInterface)

	user, err := handlers.usecases.GetUserByNickname(nickname)
	if err != nil {
		body, _ := json.Marshal(err)
		ctx.SetStatusCode(err.Code)
		ctx.SetContentType("application/json")
		ctx.Write(body)
		return
	}

	body, _ := user.MarshalJSON()
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	ctx.Write(body)
}

func (handlers *Handlers) UpdateUser(ctx *fasthttp.RequestCtx) {
	var userUpd models.UpdateUserFields
	var e *models.Error

	err := userUpd.UnmarshalJSON(ctx.PostBody())

	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	nicknameInterface := ctx.UserValue("nickname")
	nickname := fmt.Sprintf("%v", nicknameInterface)


	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	user, e := handlers.usecases.ChangeUser(&userUpd, nickname)
	if e != nil {
		body, _ := e.MarshalJSON()
		ctx.SetStatusCode(e.Code)
		ctx.SetContentType("application/json")
		ctx.Write(body)
		return
	}

	body, _ := user.MarshalJSON()
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	ctx.Write(body)
}

func (handlers *Handlers) GetUsers(ctx *fasthttp.RequestCtx) {

	slugInterface := ctx.UserValue("slug")
	slug := fmt.Sprintf("%v", slugInterface)

	getLimit := ctx.QueryArgs().Peek("limit")
	getSince := ctx.QueryArgs().Peek("since")
	getDesc := ctx.QueryArgs().Peek("desc")

	var params models.UserParams
	var err error

	params.Limit, err = strconv.Atoi(string(getLimit))
	if err != nil {
		params.Limit = -1
	}
	params.Since = string(getSince)
	fmt.Println("SINCE: ", params.Since)
	params.Desc = string(getDesc) == "true"

	users, e := handlers.usecases.GetUsersByForum(slug, params)
	if e != nil {
		body, _ := e.MarshalJSON()
		ctx.SetStatusCode(e.Code)
		ctx.SetContentType("application/json")
		ctx.Write(body)
		return
	}

	body, _ := users.MarshalJSON()
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	ctx.Write(body)
}
