package api

import (
	"net/http"

	"github.com/Roixys/e-fast-store-api/exception"
	"github.com/Roixys/e-fast-store-api/model"
	"github.com/gin-gonic/gin"
)

func (server *Server) getListProduct(ctx *gin.Context) {
	categoryName := ctx.Query("category")

	if categoryName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Category parameter is required",
		})
		return
	}

	var products []model.Product

	result := server.DB.Joins("JOIN categories ON categories.id = products.category_id").
		Where("categories.name = ?", categoryName).
		Find(&products)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, exception.ServerErrorResponse(result.Error))
		return
	}

	// result := server.DB.Joins("categories").Where("categories.name = ?", categoryName).Find(&products)
	// if result.Error != nil {
	// 	ctx.JSON(http.StatusInternalServerError, exception.ServerErrorResponse(result.Error))
	// 	return
	// }

	ctx.JSON(http.StatusOK, products)

}
