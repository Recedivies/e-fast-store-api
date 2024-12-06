package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/Roixys/e-fast-store-api/exception"
	"github.com/Roixys/e-fast-store-api/model"
	"github.com/Roixys/e-fast-store-api/util"
	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(user model.User) userResponse {
	return userResponse{
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, exception.ErrorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, exception.ErrorResponse(err))
		return
	}

	user := model.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
	}

	if result := server.DB.Create(&user); result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
			ctx.JSON(http.StatusConflict, exception.ViolationUniqueConstraint("Username or Email already exists"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, exception.ErrorResponse(err))
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}
