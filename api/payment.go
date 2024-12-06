package api

import (
	"errors"
	"net/http"

	"github.com/Roixys/e-fast-store-api/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (server *Server) createPayment(ctx *gin.Context) {
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

	userID := user.ID.String()

	tx := server.DB.Begin()

	var carts []model.Cart
	if err := tx.Where("user_id = ?", userID).Preload("Product").Find(&carts).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart items"})
		return
	}

	if len(carts) == 0 {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	var totalAmount uint64
	var paymentOrders []model.PaymentOrder

	for _, cart := range carts {
		totalAmount += uint64(cart.Quantity) * cart.Product.Price

		paymentOrders = append(paymentOrders, model.PaymentOrder{
			Amount:    float64(cart.Product.Price),
			Quantity:  cart.Quantity,
			UserID:    cart.UserID,
			ProductID: cart.ProductID,
		})
	}

	paymentEvent := model.PaymentEvent{
		UserID:      uuid.MustParse(userID),
		TotalAmount: totalAmount,
	}
	if err := tx.Create(&paymentEvent).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment event"})
		return
	}

	if err := tx.Create(&paymentOrders).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment orders"})
		return
	}

	if err := tx.Where("user_id = ?", userID).Delete(&model.Cart{}).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":       "Payment successfully created",
		"total_amount":  totalAmount,
		"payment_event": paymentEvent,
	})
}
