package controller

import (
	"Product_Mangement_Api/models"
	"Product_Mangement_Api/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type BrandController struct {
	service service.BrandService
}

func NewBrandController(service service.BrandService) *BrandController {
	return &BrandController{
		service: service,
	}
}
func (ctrl *BrandController) CreateBrand(c *gin.Context) {
	var req models.CreateBrandRequest
	if err := c.ShouldBind(&req); err != nil {
		SendError(c, "Invalid data", err.Error(), http.StatusBadRequest)
		return
	}
	res, err := ctrl.service.Create(c.Request.Context(), req)
	if err != nil {
		SendError(c, "Server error", err.Error(), http.StatusInternalServerError)
	}
	SendResponse(c, true, "Create successfully", res, http.StatusCreated)
}
func (ctrl *BrandController) ListBrand(c *gin.Context) {
	res, err := ctrl.service.List(c.Request.Context())
	if err != nil {
		SendError(c, "Server error", err.Error(), http.StatusInternalServerError)
		return
	}
	SendResponse(c, true, "List successfully", res, http.StatusOK)
}
func (ctrl *BrandController) GetBrandById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		SendError(c, "Invalid ID format", err.Error(), http.StatusBadRequest)
		return
	}
	res, err := ctrl.service.GetById(c.Request.Context(), id)
	if err != nil {
		SendError(c, "Server error", err.Error(), http.StatusBadRequest)
		return
	}
	SendResponse(c, true, "Get successfully", res, http.StatusOK)
}
func (ctrl *BrandController) DeleteBrand(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		SendError(c, "Invalid ID format", err.Error(), http.StatusBadRequest)
		return
	}
	if err := ctrl.service.DeleteById(c.Request.Context(), id); err != nil {
		SendError(c, "Server error", err.Error(), http.StatusBadRequest)
		return
	}
	SendResponse(c, true, "Delete successfully", nil, http.StatusOK)
}
func (ctrl *BrandController) UpdateBrand(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		SendError(c, "Invalid ID format", err.Error(), http.StatusBadRequest)
		return
	}
	var req models.UpdateBrandRequest
	if err := c.ShouldBind(&req); err != nil {
		SendError(c, "Invalid data", err.Error(), http.StatusInternalServerError)
	}
	updatedBrand, err := ctrl.service.UpdateById(c.Request.Context(), id, req)
	if err != nil {
		SendError(c, "Server error", err.Error(), http.StatusInternalServerError)
	}
	SendResponse(c, true, "Update successfully", updatedBrand, http.StatusOK)
}
