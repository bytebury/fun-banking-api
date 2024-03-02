package api

import (
	"funbanking/internal/interfaces/api/handlers"
	"funbanking/internal/interfaces/api/middleware"

	"github.com/gin-gonic/gin"
)

func Run() {
	r := runner{router: gin.Default()}
	r.setup()
	r.router.Run()
}

type runner struct {
	router *gin.Engine
}

func (r runner) setup() {
	// Setup Middleware
	r.router.Use(middleware.CorsMiddleware())

	// Setup Routes
	r.setupMetricsRoutes()
	r.setupUserRoutes()
	r.setupBankRoutes()
	r.setupCustomerRoutes()
	r.setupEmployeeRoutes()
	r.setupAccountRoutes()
	r.setupTransactionRoutes()
	r.setupAnnouncementRoutes()
	r.setupSessionRoutes()
	r.setupPasswordRoutes()
	r.setupNotificationRoutes()
}

func (r runner) setupNotificationRoutes() {
	handler := handlers.NewTransactionHandler()

	r.router.Group("/notifications", middleware.Auth()).GET("", handler.FindAllPendingTransactions)
}

func (r runner) setupMetricsRoutes() {
	handler := handlers.NewMetricsHandler()
	r.router.Group("/metrics").
		GET("", handler.GetApplicationInfo).
		GET("visitors", middleware.Admin(), handler.GetVisitorsInfo).
		GET("users", middleware.Admin(), handler.GetUsersInfo)
}

func (r runner) setupUserRoutes() {
	handler := handlers.NewUserHandler()
	r.router.GET("/current-user", middleware.Auth(), handler.GetCurrentUser)
	r.router.Group("users").
		GET("", middleware.Admin(), handler.Search).
		GET(":username", middleware.Audit(), handler.FindByUsername).
		PUT("", handler.Create).
		PATCH(":id", middleware.Auth(), handler.Update)
}

func (r runner) setupSessionRoutes() {
	userHandler := handlers.NewUserHandler()
	customerHandler := handlers.NewCustomerHandler()
	r.router.Group("sessions").
		POST("users", userHandler.Login).
		POST("customers", customerHandler.Login)
}

func (r runner) setupBankRoutes() {
	handler := handlers.NewBankHandler()
	r.router.GET("/my-banks", middleware.Auth(), handler.FindAllByUserID)
	r.router.Group("/banks").
		GET("", middleware.Admin(), handler.Search).
		GET(":id", middleware.Auth(), handler.FindByID).
		GET(":id/customers", middleware.Auth(), handler.FindAllCustomers).
		POST("", handler.FindByUsernameAndSlug).
		PUT("", middleware.Auth(), handler.Create).
		PATCH(":id", middleware.Auth(), handler.Update).
		DELETE(":id", middleware.Auth(), handler.Delete)
}

func (r runner) setupCustomerRoutes() {
	handler := handlers.NewCustomerHandler()
	r.router.Group("/customers").
		GET(":id", middleware.Auth(), handler.FindByID).
		GET(":id/accounts", middleware.Customer(), handler.FindAccounts).
		PUT("", middleware.Auth(), handler.Create).
		PATCH(":id", middleware.Auth(), handler.Update).
		DELETE(":id", middleware.Auth(), handler.Delete)
}

func (r runner) setupEmployeeRoutes() {
	handler := handlers.NewEmployeeHandler()
	r.router.Group("/employees", middleware.Auth()).
		GET("banks/:id", middleware.Auth(), handler.FindAllByBankID).
		GET("users/:id", middleware.Auth(), handler.FindAllByUserID).
		PUT("", middleware.Auth(), handler.Create)
}

func (r runner) setupAccountRoutes() {
	handler := handlers.NewAccountHandler()
	r.router.Group("/accounts").
		GET(":id", middleware.Customer(), handler.FindByID).
		GET(":id/transactions", middleware.Customer(), handler.FindTransactions).
		GET(":id/insights/transactions", middleware.Customer(), handler.MonthlyTransactionInsights).
		PATCH(":id", middleware.Customer(), handler.Update).
		PUT("", middleware.Auth(), handler.Create)
}

func (r runner) setupTransactionRoutes() {
	handler := handlers.NewTransactionHandler()
	r.router.Group("/transactions").
		PATCH(":id/approve", middleware.Auth(), handler.Approve).
		PATCH(":id/decline", middleware.Auth(), handler.Decline).
		PUT("", middleware.Customer(), handler.Create)
}

func (r runner) setupAnnouncementRoutes() {
	handler := handlers.NewAnnouncementHandler()
	r.router.Group("/announcements").
		GET("", handler.FindAll).
		GET(":id", handler.FindByID).
		PUT("", middleware.Admin(), handler.Create).
		PATCH(":id", middleware.Admin(), handler.Update)
}

func (r runner) setupPasswordRoutes() {
	handler := handlers.NewPasswordHandler()
	r.router.Group("passwords").
		POST("forgot", handler.ForgotPassword).
		POST("reset", middleware.PasswordReset(), handler.ResetPassword)
}
