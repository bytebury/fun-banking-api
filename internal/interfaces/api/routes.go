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
	r.setupMetricsRoutes()
	r.setupUserRoutes()
	r.setupBankRoutes()
	r.setupCustomerRoutes()
	r.setupEmployeeRoutes()
	r.setupAccountRoutes()
	r.setupTransactionRoutes()
	r.setupAnnouncementRoutes()
	r.setupSessionRoutes()
}

func (r runner) setupMetricsRoutes() {
	handler := handlers.NewMetricsHandler()
	r.router.Group("/metrics").GET("/", handler.GetApplicationInfo)
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
	r.router.Group("/banks").
		// TODO: Search for banks GET("/banks")
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
		GET(":id", handler.FindByID).
		GET(":id/accounts", handler.FindAccounts).
		PUT("", handler.Create).
		PATCH(":id", handler.Update).
		DELETE(":id", handler.Delete)
}

func (r runner) setupEmployeeRoutes() {
	handler := handlers.NewEmployeeHandler()
	r.router.Group("/employees", middleware.Auth()).
		GET("banks/:id", handler.FindAllByBankID).
		GET("users/:id", handler.FindAllByUserID).
		PUT("", handler.Create)
}

func (r runner) setupAccountRoutes() {
	handler := handlers.NewAccountHandler()
	r.router.Group("/accounts").
		GET(":id", handler.FindByID).                      // CUSTOMER
		GET(":id/transactions", handler.FindTransactions). // CUSTOMER
		PATCH(":id", handler.Update)                       // CUSTOMER
}

func (r runner) setupTransactionRoutes() {
	handler := handlers.NewTransactionHandler()
	r.router.Group("/transactions").
		GET(":id", handler.FindByID). // CUSTOMER
		PATCH(":id/approve", middleware.Auth(), handler.Approve).
		PATCH(":id/decline", middleware.Auth(), handler.Decline).
		PUT("", handler.Create) // CUSTOMER
}

func (r runner) setupAnnouncementRoutes() {
	handler := handlers.NewAnnouncementHandler()
	r.router.Group("/announcements").
		// SEARCH ROUTE GET("/")
		GET(":id", handler.FindByID).
		PUT("", middleware.Admin(), handler.Create).
		PATCH(":id", middleware.Admin(), handler.Update)
}
