package ginuser

import (
	"net/http"
	"social-todo-list/common"
	"social-todo-list/module/user/handler"
	"social-todo-list/module/user/model"

	"github.com/gin-gonic/gin"
)

type GinUserHandler struct {
	service *handler.IUserService
}

func NewGinUserHandler(service *handler.IUserService) *GinUserHandler {
	return &GinUserHandler{
		service: service,
	}
}

func (h *GinUserHandler) Register(ctx *gin.Context) {
	var data model.UserCreation

	if err := ctx.ShouldBind(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	if err := h.service.Register(ctx.Request.Context(), &data); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(data.Id))
}

func (h GinUserHandler) Login(ctx *gin.Context) {
	var data model.UserLogin
	if err := ctx.ShouldBind(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	accessToken, err := h.service.Login(ctx.Request.Context(), &data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(accessToken))
}

func (h GinUserHandler) Profile(ctx *gin.Context) {
	// user := ctx.Value(common.CurrentUser)
	user, err := h.service.Profile(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(user))
}
