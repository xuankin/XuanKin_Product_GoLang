package controller

import (
	"Product_Mangement_Api/models"
	"Product_Mangement_Api/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type VariantController struct {
	service service.VariantService
}

func NewVariantController(service service.VariantService) *VariantController {
	return &VariantController{
		service: service,
	}
}
func (ctrl *VariantController) Create(c *gin.Context) {
	var req models.CreateVariantRequest
	if err := c.ShouldBind(&req); err != nil {
		SendError(c, "Invalid data", err.Error(), http.StatusBadRequest)
		return
	}
	res, err := ctrl.service.AddVariant(c.Request.Context(), req)
	if err != nil {
		SendError(c, "Something went wrong when add variant", err.Error(), http.StatusBadRequest)
		return
	}
	SendResponse(c, true, "Add variant successfully", res, http.StatusOK)
}
func (ctrl *VariantController) UpdateVariant(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		SendError(c, "Invalid ID format", err.Error(), http.StatusBadRequest)
		return
	}
	var req models.UpdateVariantRequest
	if err := c.ShouldBind(&req); err != nil {
		SendError(c, "Invalid data", err.Error(), http.StatusBadRequest)
		return
	}
	variantUpdated, err := ctrl.service.UpdateVariant(c.Request.Context(), id, req)
	if err != nil {
		SendError(c, "Error when update variant", err.Error(), http.StatusInternalServerError)
	}
	SendResponse(c, true, "Update variant successfully", variantUpdated, http.StatusOK)
}
func (ctrl *VariantController) DeleteVariant(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		SendError(c, "Invalid ID format", err.Error(), http.StatusBadRequest)
		return
	}
	if err := ctrl.service.DeleteVariant(c.Request.Context(), id); err != nil {
		SendError(c, "Error when delete variant", err.Error(), http.StatusInternalServerError)
		return
	}
	SendResponse(c, true, "Delete variant", id, http.StatusOK)
}
func (ctrl *VariantController) GetVariantByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		SendError(c, "Invalid ID format", err.Error(), http.StatusBadRequest)
		return
	}
	variant, err := ctrl.service.GetVariantByID(c.Request.Context(), id)
	if err != nil {
		SendError(c, "Error when get variant", err.Error(), http.StatusInternalServerError)
		return
	}
	SendResponse(c, true, "Get variant", variant, http.StatusOK)
}
