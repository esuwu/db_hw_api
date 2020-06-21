package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	useCase2 "main/usecase"
	"net/http"
)

type Handlers struct {
	usecases useCase2.UseCase
}

func NewHandlers(ucases useCase2.UseCase) *Handlers {
	return &Handlers{
		usecases: ucases,
	}
}

func WriteResponse(w http.ResponseWriter, body []byte, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(body)
}

func (handlers *Handlers) GetStatus(ctx *fasthttp.RequestCtx) {
	status, _ := handlers.usecases.GetStatus()

	fmt.Println(status)

	body, _ := json.Marshal(status)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Write(body)
}

func (handlers *Handlers) Clear(ctx *fasthttp.RequestCtx) {
	handlers.usecases.RemoveAllData()
}
