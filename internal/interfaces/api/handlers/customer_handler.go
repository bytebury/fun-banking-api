package handlers

import (
	"funbanking/internal/domain/repository"
	"funbanking/internal/domain/service"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	customerService service.CustomerService
}

func NewCustomerHandler() CustomerHandler {
	return CustomerHandler{
		customerService: service.NewCustomerService(
			repository.NewCustomerRepository(),
		),
	}
}

func (h CustomerHandler) FindByID(c *gin.Context) {
	h.customerService.FindByID(c)
}

func (h CustomerHandler) FindAccounts(c *gin.Context) {
	h.customerService.FindAccounts(c)
}

func (h CustomerHandler) Create(c *gin.Context) {
	h.customerService.Create(c)
}

func (h CustomerHandler) Update(c *gin.Context) {
	h.customerService.Update(c)
}

func (h CustomerHandler) Delete(c *gin.Context) {
	h.customerService.Delete(c)
}
