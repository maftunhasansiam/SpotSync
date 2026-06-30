package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/maftunhasansiam/SpotSync/config"
	"github.com/maftunhasansiam/SpotSync/handler"
	"github.com/maftunhasansiam/SpotSync/middleware"
	"github.com/maftunhasansiam/SpotSync/repository"
	"github.com/maftunhasansiam/SpotSync/service"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment variables")
	}

	config.ConnectDatabase()
	db := config.DB

	userRepo := repository.NewUserRepository(db)
	zoneRepo := repository.NewZoneRepository(db)
	reservationRepo := repository.NewReservationRepository(db)

	authService := service.NewAuthService(userRepo)
	zoneService := service.NewZoneService(zoneRepo)
	reservationService := service.NewReservationService(reservationRepo, zoneRepo)

	validate := validator.New()

	authHandler := handler.NewAuthHandler(authService, validate)
	zoneHandler := handler.NewZoneHandler(zoneService, validate)
	reservationHandler := handler.NewReservationHandler(reservationService, validate)

	e := echo.New()
	e.HideBanner = true

	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"message": "SpotSync API is running",
		})
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok", "service": "SpotSync API"})
	})

	api := e.Group("/api/v1")

	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	zones := api.Group("/zones")
	zones.GET("", zoneHandler.GetAllZones)
	zones.GET("/:id", zoneHandler.GetZone)
	zones.POST("", zoneHandler.CreateZone, middleware.JWTMiddleware, middleware.AdminOnly)
	zones.PUT("/:id", zoneHandler.UpdateZone, middleware.JWTMiddleware, middleware.AdminOnly)
	zones.DELETE("/:id", zoneHandler.DeleteZone, middleware.JWTMiddleware, middleware.AdminOnly)

	reservations := api.Group("/reservations", middleware.JWTMiddleware)
	reservations.GET("/my-reservations", reservationHandler.GetMyReservations)
	reservations.POST("", reservationHandler.CreateReservation)
	reservations.DELETE("/:id", reservationHandler.CancelReservation)
	reservations.GET("", reservationHandler.GetAllReservations, middleware.AdminOnly)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("SpotSync API running on port %s\n", port)
	if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatal("Failed to start server:", err)
	}
}
