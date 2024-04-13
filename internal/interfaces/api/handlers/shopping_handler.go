package handlers

import (
	"funbanking/internal/domain/shopping"
	"funbanking/internal/domain/users"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ShoppingHandler interface {
	// Shops
	FindShopByID(ctx *gin.Context)
	FindShopsByUser(ctx *gin.Context)
	SaveShop(ctx *gin.Context)
	DeleteShop(ctx *gin.Context)
	// Items
	FindItemByID(ctx *gin.Context)
	FindItemsByStoreID(ctx *gin.Context)
	SaveItem(ctx *gin.Context)
	DeleteItem(ctx *gin.Context)
	// Purchasing
	BuyItems(ctx *gin.Context)
}

type shoppingHandler struct {
	shopService shopping.ShopService
}

func NewShoppingHandler(shopService shopping.ShopService) ShoppingHandler {
	return shoppingHandler{shopService}
}

func (handler shoppingHandler) FindShopByID(ctx *gin.Context) {
	shopID := ctx.Param("id")
	shop, err := handler.shopService.FindByID(shopID)

	if err != nil {
		handler.handleError(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, shop)
}

func (handler shoppingHandler) FindShopsByUser(ctx *gin.Context) {
	shop, err := handler.shopService.FindAllByUser(ctx.MustGet("user").(users.User))

	if err != nil {
		handler.handleError(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, shop)
}

func (handler shoppingHandler) SaveShop(ctx *gin.Context) {
	var request shopping.Shop

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	shop, err := handler.shopService.Save(request, ctx.MustGet("user").(users.User))

	if err != nil {
		handler.handleError(err, ctx)
		return
	}

	ctx.JSON(http.StatusAccepted, shop)
}

func (handler shoppingHandler) DeleteShop(ctx *gin.Context) {
	shopID := ctx.Param("id")

	if err := handler.shopService.Delete(shopID, ctx.MustGet("user").(users.User)); err != nil {
		handler.handleError(err, ctx)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "Successfully deleted that shop"})
}

func (handler shoppingHandler) FindItemByID(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)
}

func (handler shoppingHandler) FindItemsByStoreID(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)
}

func (handler shoppingHandler) SaveItem(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)
}

func (handler shoppingHandler) DeleteItem(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)
}

func (handler shoppingHandler) BuyItems(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)
}

func (handler shoppingHandler) handleError(err error, ctx *gin.Context) {
	if strings.Contains(err.Error(), "forbidden") {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "You do not have permissions to do that"})
		return
	}

	if strings.Contains(err.Error(), "idx_user_name") {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Looks like you have a shop with that name already"})
		return
	}

	if strings.Contains(err.Error(), "fk_shops_user") {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "The user you provided is either missing or invalid"})
		return
	}

	if strings.Contains(err.Error(), "record not found") {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Unable to find that resource"})
		return
	}

	if strings.Contains(err.Error(), "insufficient funds") {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "You do not have enough money to do that"})
		return
	}

	if strings.Contains(err.Error(), "empty cart") {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "You need to add items to your cart first"})
		return
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
}
