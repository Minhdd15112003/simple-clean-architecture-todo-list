package ginitem

import (
	"net/http"
	"social-todo-list/common"
	"social-todo-list/module/userlikeitem/handler"
	"social-todo-list/module/userlikeitem/model"

	"github.com/gin-gonic/gin"
)

type GinLikeHandler struct {
	service *handler.LikeService
}

func NewGinLikeHandler(service *handler.LikeService) *GinLikeHandler {
	return &GinLikeHandler{
		service: service,
	}
}

func (h *GinLikeHandler) LikeItem(ctx *gin.Context) {

	id, err := common.FromBase58(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrInternal(err))
		return
	}

	requester, err := common.GetRequesterFromContext(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	if err := h.service.LikeItem(
		ctx.Request.Context(),
		&model.Like{
			UserID: requester.GetUserID(),
			ItemID: int(id),
		}); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
}

func (h *GinLikeHandler) UnLikeItem(ctx *gin.Context) {

	id, err := common.FromBase58(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrInternal(err))
		return
	}

	requester, err := common.GetRequesterFromContext(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	if err := h.service.UnLikeItem(
		ctx.Request.Context(),
		requester.GetUserID(),
		int(id),
	); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
}

func (h *GinLikeHandler) ListUserLikeItem(ctx *gin.Context) {

	id, err := common.FromBase58(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrInternal(err))
		return
	}

	// requester, err := common.GetRequesterFromContext(ctx.Request.Context())
	// if err != nil {
	// 	ctx.JSON(http.StatusBadRequest, err)
	// 	return
	// }
	var paging common.Paging
	if err := ctx.ShouldBind(&paging); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	paging.Process()
	data, err := h.service.ListUsers(ctx, int(id), &paging)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	for i := range data {
		data[i].Mask()
	}

	ctx.JSON(
		http.StatusOK,
		common.NewSuccessResponse(data, paging, nil),
	)
}
