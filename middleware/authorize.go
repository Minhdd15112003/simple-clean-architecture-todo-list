package middleware

import (
	"context"
	"fmt"
	"social-todo-list/common"
	"social-todo-list/components/tokenprovider"
	"social-todo-list/module/user/model"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthenStore interface {
	FindUser(ctx context.Context,
		conditions map[string]interface{},
		moreInfor ...string) (*model.User, error)
}

func ErrWrongAuthHeader(err error) *common.AppError {
	return common.NewCustomError(err,
		fmt.Sprintf("wrong authen header"),
		fmt.Sprintf("ErrWrongAuthHeader"))
}

func extractTokenFromHeaderString(s string) (string, error) {
	parts := strings.Split(s, " ")
	if parts[0] != "Bearer" || len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		return "", ErrWrongAuthHeader(nil)
	}
	return parts[1], nil
}

func RequiredAuth(authStore AuthenStore, tokenProvider tokenprovider.Provider) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token, err := extractTokenFromHeaderString(ctx.GetHeader("Authorization"))
		if err != nil {
			ctx.JSON(500, common.ErrInternal(err))
			return
		}

		payload, err := tokenProvider.Validate(token)
		if err != nil {
			ctx.JSON(500, common.ErrInternal(err))
			return
		}

		user, err := authStore.FindUser(ctx.Request.Context(), map[string]interface{}{
			"id": payload.UserId(),
		})
		if err != nil {
			ctx.JSON(400, common.ErrCannotGetEntity(model.EntityName, err))
			return
		}
		if user.Status == 0 {
			ctx.JSON(400, common.ErrNoPermission(err))
			return
		}

		// user.Mask(common.Db)
		ctx.Set(common.CurrentUser, user)
		ctx.Next()
	}
}
