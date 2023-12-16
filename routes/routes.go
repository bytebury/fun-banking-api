package routes

import (
	"golfer/controllers"
	"golfer/mailers"
	"golfer/middleware"
	"golfer/repositories"
	"golfer/services"

	"github.com/gin-gonic/gin"
)

var userService services.UserService
var bankService services.BankService
var jwtService services.JwtService
var passwordService services.PasswordService

func setupServices() {
	jwtService = services.JwtService{}
	userService = *services.NewUserService(*repositories.NewUserRepository(), jwtService)
	bankService = *services.NewBankService(*repositories.NewBankRepository())
	passwordService = *services.NewPasswordService(userService, jwtService, *mailers.NewPasswordResetMailer())
}

/**
 * Sets up all of the routes for the application.
 */
func SetupRoutes(router *gin.Engine) {
	setupServices()

	setupHealthCheckRoutes(router)
	setupAuthRoutes(router)
	setupUserRoutes(router)
	setupPasswordRoutes(router)
	setupBankRoutes(router)
}

/**
 * Setups the health check route found at `/health`.
 */
func setupHealthCheckRoutes(router *gin.Engine) {
	router.GET("/health", controllers.GetHealthCheckController)
}

/**
 * Setups the users routes at `/users`.
 */
func setupUserRoutes(router *gin.Engine) {
	controller := controllers.NewUserController(userService)
	router.Group("/users").
		GET(":id", controller.FindByID).
		PUT(":id", middleware.AuthMiddleware(), controller.Update).
		POST("", controller.Create).
		DELETE(":id", middleware.AuthMiddleware(), controller.Delete)
}

/**
 * Sets up the bank routes at `/banks`.
 */
func setupBankRoutes(router *gin.Engine) {
	controller := controllers.NewBankController(bankService)
	router.Group("/banks").
		GET(":id", controller.FindByID).
		GET("", middleware.AuthMiddleware(), controller.Search).
		PUT(":id", middleware.AuthMiddleware(), controller.Update).
		POST("", middleware.AuthMiddleware(), controller.Create).
		DELETE(":id", middleware.AuthMiddleware(), controller.Delete)
}

/**
 * Setups the authentication routes like login
 */
func setupAuthRoutes(router *gin.Engine) {
	controller := controllers.NewSessionsController(userService)
	router.Group("").
		POST("login", controller.Login)
}

func setupPasswordRoutes(router *gin.Engine) {
	controller := controllers.NewPasswordController(passwordService)
	router.Group("/passwords").
		POST("forgot", controller.SendForgotPasswordEmail).
		POST("reset", controller.ResetPassword)
}
