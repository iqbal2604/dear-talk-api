package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iqbal2604/dear-talk-api.git/internal/domain"
	"github.com/iqbal2604/dear-talk-api.git/pkg/response"
)

type UserHandler struct {
	userUsecase domain.UserManagementUsecase
}

func NewUserHandler(userUsecase domain.UserManagementUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

// ─── Get My Profile ───────────────────────────────────────────────────────────

func (h *UserHandler) GetMyProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	user, err := h.userUsecase.GetProfile(userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.OK(c, "profile fetched", user)
}

// ─── Update My Profile ────────────────────────────────────────────────────────

func (h *UserHandler) UpdateMyProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req domain.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request", err.Error())
		return
	}

	user, err := h.userUsecase.UpdateProfile(userID, &req)
	if err != nil {
		response.Conflict(c, err.Error())
	}

	response.OK(c, "profile updated", user)
}

// ─── Search User ────────────────────────────────────────────────────────

func (h *UserHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		response.BadRequest(c, "search query is required", nil)
		return
	}

	users, err := h.userUsecase.SearchUsers(query)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, "user fetched", users)
}

// ─── Get User By ID ─────────────────────────────────────────────────────────

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id", nil)
		return
	}

	user, err := h.userUsecase.GetUserByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.OK(c, "user fetched", user)
}
