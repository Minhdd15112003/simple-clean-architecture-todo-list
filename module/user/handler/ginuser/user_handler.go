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
	}

	if err := h.service.Register(ctx.Request.Context(), &data); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}

	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(data.Id))
}
