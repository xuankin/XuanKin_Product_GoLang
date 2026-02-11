package controller

import (
	"Product_Mangement_Api/models"
	"Product_Mangement_Api/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type InventoryController struct {
	service service.InventoryService
}

func NewInventoryController(service service.InventoryService) *InventoryController {
	return &InventoryController{
		service: service,
	}
}
func (ctrl *InventoryController) AdjustStock(c *gin.Context) {
	var req models.UpdateInventoryRequest
	if err := c.ShouldBind(&req); err != nil {
		SendError(c, "Invalid data", err.Error(), http.StatusBadRequest)
		return
	}
	if err := ctrl.service.AdjustStock(c.Request.Context(), req); err != nil {
		SendError(c, "Adjust stock failed", err.Error(), http.StatusInternalServerError)
		return
	}
	SendResponse(c, true, "AdjustStock successfully", nil, http.StatusOK)
}
func (ctrl *InventoryController) GetStockByVariant(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		SendError(c, "Invalid ID format", err.Error(), http.StatusBadRequest)
		return
	}
	stock, err := ctrl.service.GetStockByVariant(c.Request.Context(), id)
	if err != nil {
		SendError(c, "Get stock failed", err.Error(), http.StatusInternalServerError)
		return
	}
	SendResponse(c, true, "Get stock by variant", stock, http.StatusOK)
}
