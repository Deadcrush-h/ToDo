package users_transport_http

import (
	"context"
	"net/http"

	"github.com/Deadcrush-h/ToDo/internal/core/domain"
	core_http_server "github.com/Deadcrush-h/ToDo/internal/core/transport/http/server"
)

type UsersHTTPHandler struct {
	userService UserService
}

type UserService interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)
}

func NewUsersHTTPHandler(
	userService UserService,
) *UsersHTTPHandler {
	if userService == nil {
		panic("userService is required for UsersHTTPHandler")
	}

	return &UsersHTTPHandler{
		userService: userService,
	}
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
	}
}
