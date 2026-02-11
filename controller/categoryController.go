package controller

import (
	"Product_Mangement_Api/models"
	"Product_Mangement_Api/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type CategoryController struct {
	service service.CategoryService
}

func NewCategoryController(service service.CategoryService) *CategoryController {
	return &CategoryController{
		service: service,
	}
}
func (ctrl *CategoryController) CreateCategory(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		SendError(c, "Invalid Data", err.Error(), http.StatusBadRequest)
		return
	}
	res, err := ctrl.service.Create(c.Request.Context(), req)
	if err != nil {
		SendError(c, err.Error(), err.Error(), http.StatusBadRequest)
		return
	}
	SendResponse(c, true, "Create successfully", res, http.StatusOK)
}
func (ctrl *CategoryController) List(c *gin.Context) {
	res, err := ctrl.service.ListAll(c.Request.Context())
	if err != nil {
		SendError(c, err.Error(), err.Error(), http.StatusInternalServerError)
		return
	}
	SendResponse(c, true, "List successfully", res, http.StatusOK)
}
func (ctrl *CategoryController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		SendError(c, "Invalid ID format", err.Error(), http.StatusBadRequest)
		return
	}
	if err := ctrl.service.Delete(c.Request.Context(), id); err != nil {
		SendError(c, err.Error(), err.Error(), http.StatusBadRequest)
		return
	}
	SendResponse(c, true, "Delete successfully", nil, http.StatusOK)
}
func (ctrl *CategoryController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		SendError(c, "Invalid ID format", err.Error(), http.StatusBadRequest)
		return
	}
	var req models.UpdateCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		SendError(c, "Invalid Data", err.Error(), http.StatusInternalServerError)
		return
	}
	updatedCat, err := ctrl.service.Update(c.Request.Context(), req, id)
	if err != nil {
		SendError(c, err.Error(), err.Error(), http.StatusInternalServerError)
		return
	}
	SendResponse(c, true, "Update successfully", updatedCat, http.StatusOK)
}
