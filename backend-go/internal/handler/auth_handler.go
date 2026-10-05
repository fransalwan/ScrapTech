package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"scrapflow-backend/internal/usecase"
	"scrapflow-backend/pkg/response"
)

type AuthHandler struct {
	db          *gorm.DB
	authUsecase usecase.AuthUsecase
	jwtSecret   string
}

func NewAuthHandler(db *gorm.DB, aUsecase usecase.AuthUsecase, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		db:          db,
		authUsecase: aUsecase,
		jwtSecret:   jwtSecret,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req usecase.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid input validation", err.Error())
		return
	}

	res, err := h.authUsecase.RegisterCompanyAndUser(h.db, &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Company and user registered successfully", res)
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid login credentials", err.Error())
		return
	}

	res, err := h.authUsecase.Login(h.db, req.Email, req.Password, h.jwtSecret)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Login successful", res)
}
