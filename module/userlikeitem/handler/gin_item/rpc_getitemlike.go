package ginitem

import (
	"net/http"
	"social-todo-list/common"
	usecase "social-todo-list/module/userlikeitem/use_case"

	"github.com/gin-gonic/gin"
)

type likeRPC struct {
	store usecase.UserLikeItemStore
}

func NewLikeRPC(store usecase.UserLikeItemStore) *likeRPC {
	return &likeRPC{
		store: store,
	}
}

func (h *likeRPC) GetLikeItem(ctx *gin.Context) {

	type RequestData struct {
		Ids []int `json:"ids"`
	}

	var data RequestData

	if err := ctx.ShouldBind(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrInternal(err))
		return
	}

	mapRs, err := h.store.GetLikeItem(ctx, data.Ids)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrInternal(err))
		return
	}

	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(mapRs))
}
