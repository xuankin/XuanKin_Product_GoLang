package controller

import (
	"Product_Mangement_Api/models"
	"Product_Mangement_Api/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type AttributeController struct {
	service service.AttributeService
}

func NewAttributeController(service service.AttributeService) *AttributeController {
	return &AttributeController{
		service: service,
	}
}
func (ctrl *AttributeController) Create(c *gin.Context) {
	var req models.CreateAttributeRequest
	if err := c.ShouldBind(&req); err != nil {
		SendError(c, "Invalid data", err.Error(), http.StatusBadRequest)
		return
	}
	res, err := ctrl.service.CreateAttribute(c.Request.Context(), req)
	if err != nil {
		SendError(c, "Server error", err.Error(), http.StatusInternalServerError)
		return
	}
	SendResponse(c, true, "Create successfully", res, http.StatusOK)
}
func (ctrl *AttributeController) AddValue(c *gin.Context) {
	var req models.CreateAttributeValueRequest
	if err := c.ShouldBind(&req); err != nil {
		SendError(c, "Invalid data", err.Error(), http.StatusBadRequest)
		return
	}
	res, err := ctrl.service.AddValue(c.Request.Context(), req)
	if err != nil {
		SendError(c, "Server error", err.Error(), http.StatusInternalServerError)
		return
	}
	SendResponse(c, true, "Add successfully", res, http.StatusOK)
}
func (ctrl *AttributeController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		SendError(c, "Invalid ID format", err.Error(), http.StatusBadRequest)
		return
	}
	if err := ctrl.service.DeleteAttribute(c.Request.Context(), id); err != nil {
		SendError(c, "Server error", err.Error(), http.StatusInternalServerError)
		return
	}
	SendResponse(c, true, "Delete successfully", gin.H{"id": id}, http.StatusOK)
}
func (ctrl *AttributeController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		SendError(c, "Invalid ID format", err.Error(), http.StatusBadRequest)
		return
	}
	var req models.UpdateAttributeRequest
	if err := c.ShouldBind(&req); err != nil {
		SendError(c, "Invalid data", err.Error(), http.StatusBadRequest)
		return
	}
	attrUpdated, err := ctrl.service.UpdatedAttribute(c.Request.Context(), id, req)
	if err != nil {
		SendError(c, "Server error", err.Error(), http.StatusInternalServerError)
		return
	}
	SendResponse(c, true, "Update successfully", attrUpdated, http.StatusOK)
}
func (ctrl *AttributeController) ListAttributes(c *gin.Context) {
	listAttr, err := ctrl.service.ListAttributes(c.Request.Context())
	if err != nil {
		SendError(c, "Server error", err.Error(), http.StatusInternalServerError)
		return
	}
	SendResponse(c, true, "List successfully", gin.H{"list": listAttr}, http.StatusOK)
}
func (ctrl *AttributeController) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		SendError(c, "Invalid ID format", err.Error(), http.StatusBadRequest)
		return
	}
	attr, err := ctrl.service.GetAttributesById(c.Request.Context(), id)
	if err != nil {
		SendError(c, "Server error", err.Error(), http.StatusInternalServerError)
		return
	}
	SendResponse(c, true, "Get successfully", gin.H{"attr": attr}, http.StatusOK)
}
