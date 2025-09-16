// @title MCA Bank Auth API
// @version 1.0
// @description This is the authentication service for MCA Bank
// @BasePath /

// @host localhost:5001
// @schemes http
package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/rasteiro11/MCABankAuth/docs"
	"github.com/rasteiro11/MCABankAuth/entities"
	pbCustomer "github.com/rasteiro11/MCABankAuth/gen/proto/go"
	"github.com/rasteiro11/MCABankAuth/pkg/security"
	"github.com/rasteiro11/MCABankAuth/pkg/token"
	"github.com/rasteiro11/MCABankAuth/pkg/validator"
	authGrpc "github.com/rasteiro11/MCABankAuth/src/auth/delivery/grpc"
	authHttp "github.com/rasteiro11/MCABankAuth/src/auth/delivery/http"
	authService "github.com/rasteiro11/MCABankAuth/src/auth/service"
	"github.com/rasteiro11/MCABankAuth/src/user/domain"
	usersRepo "github.com/rasteiro11/MCABankAuth/src/user/repository"
	usersService "github.com/rasteiro11/MCABankAuth/src/user/service"
	"github.com/rasteiro11/PogCore/pkg/config"
	"github.com/rasteiro11/PogCore/pkg/database"
	"github.com/rasteiro11/PogCore/pkg/logger"
	"github.com/rasteiro11/PogCore/pkg/server"
	"github.com/rasteiro11/PogCore/pkg/transport/grpcserver"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func main() {
	ctx := context.Background()
	database, err := database.NewDatabase(database.GetMysqlEngineBuilder)
	if err != nil {
		logger.Of(ctx).Fatalf("[main] database.NewDatabase() retunrned error: %+v\n", err)
	}

	if err := database.Migrate(entities.GetEntities()...); err != nil {
		logger.Of(ctx).Fatalf("[main] database.Migrate() returned error: %+v\n", err)
	}

	server := server.NewServer()
	server.AddHandler("/swagger/*", "", http.MethodGet, fiberSwagger.WrapHandler)
	server.Use("/*", cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	db := database.Conn()

	usersRepo := usersRepo.NewRepository(db)

	hasher := security.NewPasswordHasher()
	emailValidator := validator.NewEmailValidator()
	jwtService := token.NewJWTService[domain.Claims](config.Instance().RequiredString("JWT_SECRET"), 2*time.Hour)

	usersService := usersService.NewUserService(usersRepo)

	authService := authService.NewAuthService(usersService, hasher, emailValidator, jwtService)

	authGrpcService := authGrpc.NewService(authGrpc.WithAuthService(authService))

	go func() {
		server := grpcserver.NewServer(grpcserver.WithReflectionEnabled())

		server.Register(pbCustomer.AuthService_ServiceDesc, authGrpcService)

		if err := server.Run(); err != nil {
			logger.Of(ctx).Fatalf("[main] server.Run() returned error: %+v", err)
		}
	}()

	authHttp.NewHandler(server, authHttp.WithAuthService(authService))

	server.PrintRouter()

	logger.Of(ctx).Debug("Testing deploy")

	if err := server.Start(os.Getenv("SERVER_PORT")); err != nil {
		logger.Of(ctx).Fatalf("[main] server.NewServer() returned error: %+v\n", err)
	}
}
