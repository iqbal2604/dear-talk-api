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

// GetMyProfile godoc
// @Summary      Lihat profil sendiri
// @Description  Mendapatkan data profil user yang sedang login
// @Tags         Users
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response{data=domain.User}
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /users/me [get]
func (h *UserHandler) GetMyProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	user, err := h.userUsecase.GetProfile(userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.OK(c, "profile fetched", user)
}

// UpdateMyProfile godoc
// @Summary      Update profil sendiri
// @Description  Mengupdate username atau avatar user yang sedang login
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      domain.UpdateProfileRequest  true  "Update Profile Request"
// @Success      200      {object}  response.Response{data=domain.User}
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Failure      409      {object}  response.Response
// @Router       /users/me [put]
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

// SearchUsers godoc
// @Summary      Cari user
// @Description  Mencari user berdasarkan username atau email
// @Tags         Users
// @Produce      json
// @Security     BearerAuth
// @Param        q    query     string  true  "Search query"
// @Success      200  {object}  response.Response{data=[]domain.User}
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /users/search [get]
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

// GetUserByID godoc
// @Summary      Lihat profil user lain
// @Description  Mendapatkan data profil user berdasarkan ID
// @Tags         Users
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  response.Response{data=domain.User}
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /users/{id} [get]
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
