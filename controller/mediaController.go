package controller

import (
	"Product_Mangement_Api/models"
	"Product_Mangement_Api/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding" // Thêm import này
)

type MediaController struct {
	mediaService service.MediaService
}

func NewMediaController(mediaService service.MediaService) *MediaController {
	return &MediaController{mediaService: mediaService}
}

func (ctrl *MediaController) Upload(c *gin.Context) {
	var req models.CreateMediaRequest

	if err := c.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		SendError(c, "Invalid request", err.Error(), http.StatusBadRequest)
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		SendError(c, "File not found", err.Error(), http.StatusBadRequest)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		SendError(c, "Cannot open file", err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 3. Gọi service (Service sẽ tự detect IMAGE/VIDEO)
	res, err := ctrl.mediaService.Upload(
		c.Request.Context(),
		file,
		fileHeader.Filename,
		req,
	)

	if err != nil {
		SendError(c, "Upload failed", err.Error(), http.StatusBadRequest)
		return
	}

	SendResponse(c, true, "Upload successfully", res, http.StatusOK)
}
