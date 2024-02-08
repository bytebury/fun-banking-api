package api

import (
	"funbanking/internal/interfaces/api/handlers"

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
	r.setupHealthRoutes()
	r.setupUserRoutes()
	r.setupBankRoutes()
	r.setupCustomerRoutes()
	r.setupAccountRoutes()
	r.setupTransactionRoutes()
}

func (r runner) setupHealthRoutes() {
	handler := handlers.NewHealthHandler()
	r.router.Group("/health").GET("/", handler.GetHealthCheck)
}

func (r runner) setupUserRoutes() {
	handler := handlers.NewUserHandler()
	r.router.GET("/current-user", handler.GetCurrentUser)
	r.router.Group("users").
		GET("", handler.GetCurrentUser).
		GET(":id", handler.FindByID).
		PUT("", handler.Create).
		PATCH(":id", handler.Update)
}

func (r runner) setupBankRoutes() {
	handler := handlers.NewBankHandler()
	r.router.Group("/banks").
		GET("", handler.FindAllMyBanks).
		GET(":id", handler.FindByID).
		GET(":id/customers", handler.FindAllCustomers).
		POST("", handler.FindByUsernameAndSlug).
		PUT("", handler.Create).
		PATCH(":id", handler.Update).
		DELETE(":id", handler.Delete)
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

func (r runner) setupAccountRoutes() {
	handler := handlers.NewAccountHandler()
	r.router.Group("/accounts").
		GET(":id", handler.FindByID).
		GET(":id/transactions", handler.FindTransactions).
		PATCH(":id", handler.Update)
}

func (r runner) setupTransactionRoutes() {
	handler := handlers.NewTransactionHandler()
	r.router.Group("/transactions").
		GET(":id", handler.FindByID).
		PATCH(":id/approve", handler.Approve).
		PATCH(":id/decline", handler.Decline).
		PUT("", handler.Create)
}
