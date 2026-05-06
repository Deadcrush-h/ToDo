package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/Deadcrush-h/ToDo/internal/core/logger"
	core_http_middleware "github.com/Deadcrush-h/ToDo/internal/core/transport/http/middleware"
	core_http_server "github.com/Deadcrush-h/ToDo/internal/core/transport/http/server"
	users_transport_http "github.com/Deadcrush-h/ToDo/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cansel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cansel()

	fmt.Println("Hello, TODO app")
	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to unut application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("Starting ToDo application!")

	usersTransportHTTP := users_transport_http.NewUsersHTTTPHandler(nil)

	userRoutes := usersTransportHTTP.Routes()

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(userRoutes...)

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)

	httpServer.RegisterAPIRouters(*apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
