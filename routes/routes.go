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
var customerService services.CustomerService
var accountService services.AccountService
var moneyTransferService services.MoneyTransferService
var jwtService services.JwtService
var passwordService services.PasswordService

func setupServices() {
	jwtService = services.JwtService{}
	userService = *services.NewUserService(*repositories.NewUserRepository(), jwtService)
	bankService = *services.NewBankService(*repositories.NewBankRepository())
	customerService = *services.NewCustomerService(*repositories.NewCustomerRepository())
	accountService = *services.NewAccountService(*repositories.NewAccountRepository())
	moneyTransferService = *services.NewMoneyTransferService(*repositories.NewMoneyTransferRepository(), accountService, userService)
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
	setupCustomerRoutes(router)
	setupAccountRoutes(router)
	setupMoneyTransferRoutes(router)
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
		GET("", middleware.Auth(), controller.FindCurrentUser).
		GET(":username", middleware.Auth(), controller.FindByUsername).
		PUT(":id", middleware.Auth(), controller.Update).
		POST("", controller.Create).
		DELETE(":id", middleware.Auth(), controller.Delete)
}

/**
 * Sets up the bank routes at `/banks`.
 */
func setupBankRoutes(router *gin.Engine) {
	controller := controllers.NewBankController(bankService)
	router.Group("/banks").
		GET("", middleware.Auth(), controller.FindBanksByUserID).
		GET(":id", middleware.Auth(), controller.FindByID).
		GET(":id/customers", middleware.Auth(), controller.FindCustomers).
		PUT(":id", middleware.Auth(), controller.Update).
		POST("", middleware.Auth(), controller.Create).
		DELETE(":id", middleware.Auth(), controller.Delete)
}

/**
 * Sets up the customer routes at `/customers`.
 */
func setupCustomerRoutes(router *gin.Engine) {
	controller := controllers.NewCustomerController(customerService, bankService, accountService)
	router.Group("/customers").
		GET(":id", middleware.Auth(), controller.FindByID).
		GET(":id/accounts", middleware.Auth(), controller.FindAllAccounts).
		PUT(":id", middleware.Auth(), controller.Update).
		POST("", middleware.Auth(), controller.Create).
		DELETE(":id", middleware.Auth(), controller.Delete)
}

/**
 * Sets up the accounts routes at `/accounts`.
 */
func setupAccountRoutes(router *gin.Engine) {
	controller := controllers.NewAccountController(accountService, moneyTransferService)
	router.Group("/accounts").
		GET(":id", middleware.Auth(), controller.FindByID).
		GET(":id/money-transfers", middleware.Auth(), controller.FindMoneyTransfers)
}

func setupMoneyTransferRoutes(router *gin.Engine) {
	controller := controllers.NewMoneyTransferController(moneyTransferService, accountService)
	router.Group("/money-transfers").
		POST("", middleware.Audit(), controller.Create).
		PUT(":id/approve", middleware.Auth(), controller.Approve).
		PUT(":id/decline", middleware.Auth(), controller.Decline)
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
