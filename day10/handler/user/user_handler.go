package userhandler

import (
	usermodel "assignment8/models/user"
	"strconv"

	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"
	"gofr.dev/pkg/gofr/http/response"
)

type UserHandler struct {
	userService UserServiceInterface
}

func NewUserHandler(service UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: service,
	}
}

func (h *UserHandler) CreateUserHandler(ctx *gofr.Context) (any, error) {
	var user usermodel.User

	err := ctx.Bind(&user)
	if err != nil {
		return nil, err
	}

	err = h.userService.CreateUser(ctx, &user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h *UserHandler) GetUserHandler(ctx *gofr.Context) (any, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"Give Correct Input"}}
	}

	user, err := h.userService.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return response.Raw{Data: user}, nil
}
