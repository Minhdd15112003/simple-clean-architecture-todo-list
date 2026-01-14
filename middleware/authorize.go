package middleware

import (
	"context"
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
		"wrong authen header",
		"ErrWrongAuthHeader")
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
			appErr := common.ErrInternal(err)
			ctx.AbortWithStatusJSON(appErr.StatusCode, appErr)
			return
		}

		payload, err := tokenProvider.Validate(token)
		if err != nil {
			appErr := common.ErrInternal(err)
			ctx.AbortWithStatusJSON(appErr.StatusCode, appErr)
			return
		}

		user, err := authStore.FindUser(ctx.Request.Context(), map[string]interface{}{
			"id": payload.UserId(),
		})
		if err != nil {
			appErr := common.ErrCannotGetEntity(model.EntityName, err)
			ctx.AbortWithStatusJSON(appErr.StatusCode, appErr)
			return
		}
		if user.Status == 0 {
			appErr := common.ErrNoPermission(err)
			ctx.AbortWithStatusJSON(appErr.StatusCode, appErr)
			return
		}

		// Store requester in request context (not gin context) for clean architecture
		ctx.Request = ctx.Request.WithContext(
			context.WithValue(ctx.Request.Context(), common.CurrentUser, user),
		)
		ctx.Next()
	}
}
