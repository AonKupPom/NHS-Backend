package controllers

import (
	"github.com/gin-gonic/gin"
)

func UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	fileName := ctx.PostForm("image")

	if err != nil {
		return
	}
	ctx.SaveUploadedFile(file, "./uploads/"+fileName)
}
