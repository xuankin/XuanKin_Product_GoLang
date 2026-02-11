package controller

import (
	"Product_Mangement_Api/models"
	"Product_Mangement_Api/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type WarehouseController struct {
	service service.WarehouseService
}

func NewWarehouseController(service service.WarehouseService) *WarehouseController {
	return &WarehouseController{
		service: service,
	}
}
func (ctrl *WarehouseController) Create(c *gin.Context) {
	var req models.CreateWarehouseRequest
	if err := c.ShouldBind(&req); err != nil {
		SendError(c, "Invalid data", err.Error(), http.StatusBadRequest)
		return
	}
	res, err := ctrl.service.Create(c.Request.Context(), req)
	if err != nil {
		SendError(c, "Error when create warehouse", err.Error(), http.StatusInternalServerError)
		return
	}
	SendResponse(c, true, "Create successfully", res, http.StatusOK)
}
func (ctrl *WarehouseController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		SendError(c, "Invalid ID format", err.Error(), http.StatusBadRequest)
		return
	}
	var req models.UpdateWarehouseRequest
	if err := c.ShouldBind(&req); err != nil {
		SendError(c, "Invalid data", err.Error(), http.StatusBadRequest)
		return
	}
	res, err := ctrl.service.Update(c.Request.Context(), req, id)
	if err != nil {
		SendError(c, "Error when update warehouse", err.Error(), http.StatusInternalServerError)
		return
	}
	SendResponse(c, true, "Update successfully", res, http.StatusOK)
}
func (ctrl *WarehouseController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		SendError(c, "Invalid ID format", err.Error(), http.StatusBadRequest)
		return
	}
	if err := ctrl.service.Delete(c.Request.Context(), id); err != nil {
		SendError(c, "Error when delete warehouse", err.Error(), http.StatusInternalServerError)
		return
	}
	SendResponse(c, true, "Deleted successfully", nil, http.StatusOK)
}
func (ctrl *WarehouseController) ListAll(c *gin.Context) {
	Warehouses, err := ctrl.service.List(c.Request.Context())
	if err != nil {
		SendError(c, "Error when list warehouses", err.Error(), http.StatusInternalServerError)
		return
	}
	SendResponse(c, true, "List warehouses", Warehouses, http.StatusOK)
}
func (ctrl *WarehouseController) GetWarehouseById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		SendError(c, "Invalid ID format", err.Error(), http.StatusBadRequest)
		return
	}
	warehouse, err := ctrl.service.GetById(c.Request.Context(), id)
	if err != nil {
		SendError(c, "Error when get warehouse by id", err.Error(), http.StatusInternalServerError)
		return
	}
	SendResponse(c, true, "Get Warehouse successfully", warehouse, http.StatusOK)
}
