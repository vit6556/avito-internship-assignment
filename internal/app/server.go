package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/vit6556/avito-internship-assignment/internal/config"
	"github.com/vit6556/avito-internship-assignment/internal/database/postgres"
	"github.com/vit6556/avito-internship-assignment/internal/delivery/http/handler"
	httpmiddleware "github.com/vit6556/avito-internship-assignment/internal/delivery/http/middleware"
	"github.com/vit6556/avito-internship-assignment/internal/service/auth"
	"github.com/vit6556/avito-internship-assignment/internal/service/employee"
	"github.com/vit6556/avito-internship-assignment/internal/service/merch"
	"github.com/vit6556/avito-internship-assignment/internal/service/transaction"
)

func InitServer(cfg *config.ServerConfig, dbPool *pgxpool.Pool) *echo.Echo {
	employeeRepo := postgres.NewEmployeeRepository(dbPool)
	merchRepo := postgres.NewMerchRepository(dbPool)
	transactionRepo := postgres.NewTransaction(dbPool)

	authService := authservice.NewAuthService(employeeRepo, cfg.Secret, cfg.TokenTTL, cfg.User.DefaultBalance)
	employeeService := employeeservice.NewEmployeeService(employeeRepo, merchRepo, transactionRepo)
	transactionService := transactionservice.NewTransactionService(employeeRepo, transactionRepo)
	merchService := merchservice.NewMerchService(employeeRepo, merchRepo)

	jwtMiddleware := httpmiddleware.JWTMiddleware(authService)

	authHandler := httphandler.NewAuthHandler(authService, cfg.TokenTTL, cfg.HTTPServer.Secure)
	employeeHandler := httphandler.NewEmployeeHandler(employeeService)
	transactionHandler := httphandler.NewTransactionHandler(transactionService)
	merchHandler := httphandler.NewMerchHandler(merchService)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/api/auth", authHandler.GetToken)
	e.GET("/api/info", employeeHandler.GetEmployeeInfo, jwtMiddleware)
	e.POST("/api/sendCoin", transactionHandler.SendCoin, jwtMiddleware)
	e.GET("/api/buy/:item", merchHandler.BuyItem, jwtMiddleware)

	return e
}
