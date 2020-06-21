package delivery

import (
	"github.com/valyala/fasthttp"
	useCase2 "main/usecase"
)

type Handlers struct {
	usecases useCase2.UseCase
}

func NewHandlers(ucases useCase2.UseCase) *Handlers {
	return &Handlers{
		usecases: ucases,
	}
}


func (handlers *Handlers) GetStatus(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	status := handlers.usecases.GetStatus()
	resp, _ := status.MarshalJSON()
	ctx.Write(resp)
}

func (handlers *Handlers) Clear(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	handlers.usecases.Clear()
}
