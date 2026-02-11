package controller

import (
	"Product_Mangement_Api/models"
	"Product_Mangement_Api/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MediaController struct {
	mediaService service.MediaService
}

func NewMediaController(s service.MediaService) *MediaController {
	return &MediaController{mediaService: s}
}
func (ctrl *MediaController) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		SendError(c, "File not found", err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	var req models.CreateMediaRequest
	if err := c.ShouldBind(&req); err != nil {
		SendError(c, "Invalid request", err.Error(), http.StatusBadRequest)
		return
	}
	res, err := ctrl.mediaService.Upload(c.Request.Context(), file, header.Filename, req)
	if err != nil {
		SendError(c, "Upload failed", err.Error(), http.StatusBadRequest)
		return
	}
	SendResponse(c, true, "Upload successfully", res, http.StatusOK)
}
