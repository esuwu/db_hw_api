package delivery

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	models "main/models"
	"net/http"
)


func (handlers *Handlers) CreateUser(ctx *fasthttp.RequestCtx) {
	var user models.User
	user.UnmarshalJSON(ctx.PostBody())

	name := ctx.UserValue("nickname")

	Users, err := handlers.usecases.CreateUser(&user, name)

	var resp []byte

	switch err {
	case nil:
		ctx.SetStatusCode(http.StatusCreated)
		user.Nickname = ctx.UserValue("nickname").(string)
		resp, _ = user.MarshalJSON()

	case models.UserAlreadyExists:
		ctx.SetStatusCode(http.StatusConflict)
		resp, _ = Users.MarshalJSON()
	}

	ctx.SetContentType("application/json")
	ctx.Write(resp)
}

func (handlers *Handlers) GetUser(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	name := ctx.UserValue("nickname")

	var response []byte

	userFromDb, err := handlers.usecases.GetUserProfile(name)

	switch err {
	case nil:
		response, _ = userFromDb.MarshalJSON()
	case models.UserNotFound:

		ctx.SetStatusCode(http.StatusNotFound)
		response, _ = json.Marshal(err)
		ctx.SetContentType("application/json")
	}
	ctx.Write(response)
}

func (handlers *Handlers) UpdateUser(ctx *fasthttp.RequestCtx) {
	user := models.UserUpd{}
	user.UnmarshalJSON(ctx.PostBody())

	nickname := ctx.UserValue("nickname")

	userUpdated, error := handlers.usecases.UpdateUserProfile(&user, &nickname)
	var result []byte
	switch error {
	case nil:
		result, _ = userUpdated.MarshalJSON()

	case models.ConflictOnUsers:
		ctx.SetStatusCode(http.StatusConflict)
		result, _ = json.Marshal(error)
	case models.UserNotFound:
		ctx.SetStatusCode(http.StatusNotFound)
		result, _ = json.Marshal(error)
	}

	ctx.SetContentType("application/json")
	ctx.Write(result)
}
