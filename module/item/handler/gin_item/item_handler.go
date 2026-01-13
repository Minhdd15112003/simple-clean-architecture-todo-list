package ginitem

import (
	"net/http"
	"social-todo-list/common"
	"social-todo-list/module/item/handler"
	"social-todo-list/module/item/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GinItemHandler struct {
	service *handler.ItemService
}

func NewGinItemHandler(service *handler.ItemService) *GinItemHandler {
	return &GinItemHandler{
		service: service,
	}
}

func (h *GinItemHandler) GetItems(ctx *gin.Context) {
	var queryString struct {
		common.Paging
		model.Fitter
	}

	if err := ctx.ShouldBind(&queryString); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	queryString.Process()

	itemData, err := h.service.GetItems(
		ctx.Request.Context(),
		&queryString.Fitter,
		&queryString.Paging,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}
	ctx.JSON(
		http.StatusOK,
		common.NewSuccessResponse(itemData, queryString.Fitter, queryString.Paging),
	)
}

func (h *GinItemHandler) GetItem(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	itemData, err := h.service.GetItemByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(itemData))
}

func (h *GinItemHandler) CreateItem(ctx *gin.Context) {
	var itemData model.TodoItemCreation
	if err := ctx.ShouldBind(&itemData); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	if err := h.service.CreateNewItem(ctx.Request.Context(), &itemData); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(itemData.Id))
}

func (h *GinItemHandler) UpdateItem(ctx *gin.Context) {
	var itemData model.TodoItemUpdation

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	if err := ctx.ShouldBind(&itemData); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	if err := h.service.UpdateItem(ctx.Request.Context(), id, &itemData); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
}

func (h *GinItemHandler) DeleteItem(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	if err := h.service.DeleteItem(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
}
