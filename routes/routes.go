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
var transactionService services.TransactionService
var healthService services.HealthService
var jwtService services.JwtService
var passwordService services.PasswordService
var announcementService services.AnnouncementService
var employeeService services.EmployeeService

func setupServices() {
	jwtService = services.JwtService{}
	userService = *services.NewUserService(*repositories.NewUserRepository(), jwtService)
	bankService = *services.NewBankService(*repositories.NewBankRepository())
	customerService = *services.NewCustomerService(*repositories.NewCustomerRepository())
	accountService = *services.NewAccountService(*repositories.NewAccountRepository())
	passwordService = *services.NewPasswordService(userService, jwtService, *mailers.NewPasswordResetMailer())
	healthService = *services.NewHealthService(*repositories.NewHealthRepository())
	announcementService = *services.NewAnnouncementService(*repositories.NewAnnouncementRepository())
	employeeService = *services.NewEmployeeService(*repositories.NewEmployeeRepository(), userService, bankService)
	transactionService = *services.NewTransactionService(*repositories.NewTransactionRepository(), accountService, userService, employeeService)
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
	setupTransactionRoutes(router)
	setupAnnouncementRoutes(router)
	setupNotificationRoutes(router)
	setupEmployeeRoutes(router)

	bankController := controllers.NewBankController(bankService, employeeService)
	router.GET(":username/:slug", bankController.FindByUsernameAndSlug)
}

/**
 * Setups the health check route found at `/health`.
 */
func setupHealthCheckRoutes(router *gin.Engine) {
	userRepository := repositories.NewUserRepository()
	controller := controllers.NewHealthController(healthService)

	router.GET("/health", controller.GetHealthCheck).
		GET("/health/users", middleware.Admin(*userRepository), controller.GetUserInsights)
}

func setupNotificationRoutes(router *gin.Engine) {
	controller := controllers.NewTransactionController(transactionService, accountService, employeeService)
	router.GET("/notifications", middleware.Auth(), controller.Notifications)
}

func setupAnnouncementRoutes(router *gin.Engine) {
	userRepository := repositories.NewUserRepository()
	controller := controllers.NewAnnouncementController(announcementService)
	router.Group("/announcements").
		GET("", controller.Find).
		GET(":id", controller.FindByID).
		PUT(":id", middleware.Admin(*userRepository), controller.Update).
		POST("", middleware.Admin(*userRepository), controller.Create).
		DELETE(":id", middleware.Admin(*userRepository), controller.Delete)
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

func setupEmployeeRoutes(router *gin.Engine) {
	controller := controllers.NewEmployeeController(employeeService)
	router.Group("/employees").
		POST("", middleware.Auth(), controller.Create).
		GET("/bank/:id", middleware.Auth(), controller.FindByBank).
		GET("/user/:id", middleware.Auth(), controller.FindByUser).
		DELETE(":id", middleware.Auth(), controller.Delete)
}

/**
 * Sets up the bank routes at `/banks`.
 */
func setupBankRoutes(router *gin.Engine) {
	controller := controllers.NewBankController(bankService, employeeService)
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
	controller := controllers.NewCustomerController(customerService, bankService, accountService, employeeService)
	router.Group("/customers").
		GET(":id", middleware.Auth(), controller.FindByID).
		// TODO: This needs to be audit once we do customer tokens!
		GET(":id/accounts", middleware.Audit(), controller.FindAllAccounts).
		PUT(":id", middleware.Auth(), controller.Update).
		POST("", middleware.Auth(), controller.Create).
		POST("signin", controller.Login).
		DELETE(":id", middleware.Auth(), controller.Delete)
}

/**
 * Sets up the accounts routes at `/accounts`.
 */
func setupAccountRoutes(router *gin.Engine) {
	// TODO: Need to lock this down once we do tokens for customers
	//       E.g. it's Audit right now, we'll need it to be Auth
	controller := controllers.NewAccountController(accountService, transactionService)
	router.Group("/accounts").
		GET(":id", middleware.Audit(), controller.FindByID).
		GET(":id/transactions", middleware.Audit(), controller.FindTransactions).
		GET(":id/insights/transactions", middleware.Audit(), controller.GetTransactionHistoricalData)
}

func setupTransactionRoutes(router *gin.Engine) {
	controller := controllers.NewTransactionController(transactionService, accountService, employeeService)
	router.Group("/transactions").
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
