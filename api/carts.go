package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Roixys/e-fast-store-api/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (server *Server) getCartProduct(ctx *gin.Context) {
	username := ctx.MustGet(authorizationPayloadKey).(string)
	var user model.User
	if err := server.DB.First(&user, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	var productCarts []model.Product

	result := server.DB.Joins("JOIN carts ON carts.product_id = products.id").
		Where("carts.user_id = ?", user.ID.String()).
		Find(&productCarts)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"products": productCarts})
}

func (server *Server) createCartProduct(ctx *gin.Context) {
	type AddToCartRequest struct {
		ProductID string `json:"productId" binding:"required"`
		Quantity  int    `json:"quantity" binding:"required,min=1"`
	}

	var req AddToCartRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := ctx.MustGet(authorizationPayloadKey).(string)

	var product model.Product
	if err := server.DB.First(&product, "id = ?", req.ProductID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	var user model.User
	if err := server.DB.First(&user, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	cart := model.Cart{
		UserID:    uuid.MustParse(user.ID.String()),
		ProductID: uuid.MustParse(req.ProductID),
		Quantity:  req.Quantity,
	}

	if err := server.DB.Create(&cart).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Product already exists in the cart"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var fullCart model.Cart
	if err := server.DB.Preload("User").Preload("Product").First(&fullCart, cart.ID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving full cart data"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Product added to cart",
		"cart":    fullCart,
	})
}

func (server *Server) deleteCartProduct(ctx *gin.Context) {
	username := ctx.MustGet(authorizationPayloadKey).(string)
	var user model.User
	if err := server.DB.First(&user, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	productID := ctx.Param("product_id")
	if productID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Product ID is required"})
		return
	}

	var cart model.Cart
	if err := server.DB.Where("user_id = ? AND product_id = ?", user.ID.String(), productID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found in the cart"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if err := server.DB.Delete(&cart).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product removed from cart"})
}
