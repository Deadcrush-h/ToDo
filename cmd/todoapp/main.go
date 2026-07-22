package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/Deadcrush-h/ToDo/internal/core/logger"
	core_postgres_pool "github.com/Deadcrush-h/ToDo/internal/core/repository/postgres/conn"
	core_http_middleware "github.com/Deadcrush-h/ToDo/internal/core/transport/http/middleware"
	core_http_server "github.com/Deadcrush-h/ToDo/internal/core/transport/http/server"
	user_postgres_repository "github.com/Deadcrush-h/ToDo/internal/features/users/repository/postgres"
	users_service "github.com/Deadcrush-h/ToDo/internal/features/users/service"
	users_transport_http "github.com/Deadcrush-h/ToDo/internal/features/users/transport/http"

	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	appLogger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer appLogger.Close()

	appLogger.Debug("initializing postgres connection pool")

	pool, err := core_postgres_pool.NewConnectionPool(
		ctx,
		core_postgres_pool.NewConfigMust(),
	)
	if err != nil {
		appLogger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	appLogger.Debug("initializing feature", zap.String("feature", "users"))

	usersRepository := user_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTTPHandler(usersService)

	appLogger.Debug("initializing HTTP server")
	appLogger.Debug("Starting ToDo application!")

	userRoutes := usersTransportHTTP.Routes()

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		appLogger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(appLogger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(userRoutes...)
	httpServer.RegisterAPIRouters(*apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		appLogger.Error("HTTP server run error", zap.Error(err))
	}
}
