package middleware

import (
	"log"
	"social-todo-list/common"

	"github.com/gin-gonic/gin"
)

// Recover for Gin HTTP middleware
func Recover() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				ctx.Header("Content-Type", "application/json")

				if appErr, ok := err.(*common.AppError); ok {
					ctx.AbortWithStatusJSON(appErr.StatusCode, appErr)
					// panic(err)
					return
				}
				appErr := common.ErrInternal(err.(error))
				ctx.AbortWithStatusJSON(appErr.StatusCode, appErr)
				panic(err)
				// return
			}
		}()
		ctx.Next()
	}
}

func RecoverGoroutine() {
	if err := recover(); err != nil {
		log.Printf("[PANIC RECOVERED] %v", err)

		if appErr, ok := err.(*common.AppError); ok {
			log.Printf("[APP ERROR] Code: %d, Message: %s", appErr.StatusCode, appErr.Message)
			return
		}

		log.Printf("[INTERNAL ERROR] %v", err)
	}
}
