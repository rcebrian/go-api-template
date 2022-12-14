package bootstrap

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/rcebrian/users-service/configs"
	users "github.com/rcebrian/users-service/internal"
	"github.com/rcebrian/users-service/internal/platform/server/handler/health"
	server "github.com/rcebrian/users-service/internal/platform/server/openapi"
	"github.com/rcebrian/users-service/internal/users/creating"
	"github.com/rcebrian/users-service/internal/users/finding"

	"github.com/mvrilo/go-redoc"
)

// RunInternalServer starts a server for healthcheck status
func RunInternalServer(sqlClient *sql.DB) error {
	addr := fmt.Sprintf(":%d", configs.ServiceConfig.HttpInternalPort)
	internal := http.NewServeMux()
	internal.HandleFunc("/health", health.GetHealth(sqlClient).Handler)

	doc := redoc.Redoc{
		Title:       "API Docs",
		Description: "API documentation",
		SpecFile:    "./api/openapi-spec/openapi.yaml",
		SpecPath:    "/openapi.yaml",
		DocsPath:    "/docs",
	}

	internal.HandleFunc(doc.DocsPath, doc.Handler())
	internal.HandleFunc(doc.SpecPath, doc.Handler())

	return http.ListenAndServe(addr, internal)
}

// NewServer create a new configured server
func NewServer(userRepo users.UserRepository) *http.Server {
	addr := fmt.Sprintf(":%d", configs.HttpServerConfig.Port)

	// users
	UsersApiController := usersApiController(userRepo)

	router := server.NewRouter(UsersApiController)

	return &http.Server{
		Addr: addr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * time.Duration(configs.HttpServerConfig.WriteTimeout),
		ReadTimeout:  time.Second * time.Duration(configs.HttpServerConfig.ReadTimeout),
		IdleTimeout:  time.Second * time.Duration(configs.HttpServerConfig.IdleTimeout),
		Handler:      router,
	}
}

// usersApiController configure users controller with dependency injection
func usersApiController(userRepo users.UserRepository) server.Router {
	createService := creating.NewCreatingService(userRepo)
	findAllService := finding.NewFindAllUsersUseCase(userRepo)
	findByIdService := finding.NewFindUserByIdUseCase(userRepo)

	UsersApiService := server.NewUsersApiService(createService, findAllService, findByIdService)

	return server.NewUsersApiController(UsersApiService)
}
