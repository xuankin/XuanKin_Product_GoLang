package controller

import (
	"Product_Mangement_Api/models"
	"Product_Mangement_Api/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type ProductController struct {
	service service.ProductService
}

func NewProductController(service service.ProductService) *ProductController {
	return &ProductController{
		service: service,
	}
}
func (ctrl *ProductController) Create(c *gin.Context) {
	var req models.CreateProductRequest
	if err := c.ShouldBind(&req); err != nil {
		SendError(c, "Invalid data", err.Error(), http.StatusBadRequest)
		return
	}
	res, err := ctrl.service.Create(c.Request.Context(), req)
	if err != nil {
		SendError(c, "Create product error", err.Error(), http.StatusInternalServerError)
		return
	}
	SendResponse(c, true, "Created product successfully", res, http.StatusOK)
}
func (ctrl *ProductController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		SendError(c, "Invalid ID format", err.Error(), http.StatusBadRequest)
		return
	}
	var req models.UpdateProductRequest
	if err := c.ShouldBind(&req); err != nil {
		SendError(c, "Invalid data", err.Error(), http.StatusBadRequest)
		return
	}
	res, err := ctrl.service.UpdateById(c.Request.Context(), id, req)
	if err != nil {
		SendError(c, "Update product error", err.Error(), http.StatusInternalServerError)
		return
	}
	SendResponse(c, true, "Update product successfully", res, http.StatusOK)
}
func (ctrl *ProductController) SearchProduct(c *gin.Context) {
	var params models.FilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		SendError(c, "Invalid query params", err.Error(), http.StatusBadRequest)
		return
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	res, err := ctrl.service.Search(c.Request.Context(), params)
	if err != nil {
		SendError(c, "Search error", err.Error(), http.StatusInternalServerError)
		return
	}
	SendResponse(c, true, "Search successfully", res, http.StatusOK)
}
func (ctrl *ProductController) GetProductById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		SendError(c, "Invalid ID format", err.Error(), http.StatusBadRequest)
		return
	}
	res, err := ctrl.service.GetById(c.Request.Context(), id)
	if err != nil {
		SendError(c, "Get product error", err.Error(), http.StatusInternalServerError)
		return
	}
	SendResponse(c, true, "Get product successfully", res, http.StatusOK)
}
func (ctrl *ProductController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		SendError(c, "Invalid ID format", err.Error(), http.StatusBadRequest)
		return
	}
	if err := ctrl.service.DeleteById(c.Request.Context(), id); err != nil {
		SendError(c, "Delete product error", err.Error(), http.StatusInternalServerError)
		return
	}
	SendResponse(c, true, "Delete product successfully", nil, http.StatusOK)
}
func (ctrl *ProductController) GetProductList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	res, err := ctrl.service.List(c.Request.Context(), page, limit)
	if err != nil {
		SendError(c, "Get product list error", err.Error(), http.StatusInternalServerError)
		return
	}
	SendResponse(c, true, "Get product list successfully", res, http.StatusOK)
}
