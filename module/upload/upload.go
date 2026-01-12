package upload

import (
	"fmt"
	"net/http"
	"social-todo-list/common"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Upload(db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		fileHeader, err := ctx.FormFile("file")

		if err != nil {
			ctx.JSON(400, common.ErrInvalidRequest(err))
			return
		}
		dst := fmt.Sprintf("static/%d.%s", time.Now().UTC().Nanosecond(), strings.ReplaceAll(fileHeader.Filename, " ", "_"))
		if err := ctx.SaveUploadedFile(fileHeader, dst); err != nil {
			ctx.JSON(500, common.ErrInvalidRequest(err))
			return
		}

		img := common.Image{
			Id:        0,
			Url:       dst,
			Width:     100,
			Height:    100,
			CloudName: "local",
			Extension: strings.ReplaceAll(fileHeader.Filename, " ", "_"),
		}
		img.Fulfill("http://localhost:8000")
		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(img))
	}
}
