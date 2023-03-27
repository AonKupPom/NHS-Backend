package controllers

import (
	"log"
	"os"

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

func UploadMultipleFiles(ctx *gin.Context) {
	form, _ := ctx.MultipartForm()
	files := form.File["files"]

	for _, file := range files {
		ctx.SaveUploadedFile(file, "./uploads/"+file.Filename)
	}
}

func UploadAndRemoveMultipleFiles(ctx *gin.Context) {
	form, _ := ctx.MultipartForm()
	files := form.File["files"]
	filesDelete := ctx.Request.Form["filesDelete"]

	for _, file := range files {
		ctx.SaveUploadedFile(file, "./uploads/"+file.Filename)
	}

	for _, fileDelete := range filesDelete {
		err_remove := os.Remove("./uploads/" + string(fileDelete))
		if err_remove != nil {
			log.Fatal(err_remove)
			return
		}
	}
}

func UploadAndRemoveFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	fileName := ctx.PostForm("image")
	fileDelete := ctx.PostForm("fileDelete")

	if err != nil {
		return
	}
	ctx.SaveUploadedFile(file, "./uploads/"+fileName)
	err_remove := os.Remove("./uploads/" + fileDelete)
	if err_remove != nil {
		log.Fatal(err_remove)
		return
	}
}

func RemoveFile(ctx *gin.Context) {
	fileDelete := ctx.Param("fileDelete")

	err := os.Remove("./uploads/" + fileDelete)
	if err != nil {
		log.Fatal(err)
		return
	}
}
