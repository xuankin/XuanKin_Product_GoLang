package router

import (
	"Product_Mangement_Api/config"
	"Product_Mangement_Api/controller"
	"Product_Mangement_Api/repository"
	"Product_Mangement_Api/service"
	"context"
	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func SetupRouter(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	r.Static("uploads", "./uploads")
	rdb := config.ConnectRedis(cfg)
	esClient := config.ConnectElasticsearch(cfg)
	esRepo := repository.NewElasticsearchRepository(esClient)
	esRepo.CreateIndexIfNotExists(context.Background())
	cacheRepo := repository.NewCacheRepository(rdb)
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	brandRepo := repository.NewBrandRepository(db)
	attributeRepo := repository.NewAttributeRepository(db)
	warehouseRepo := repository.NewWarehouseRepository(db)
	variantRepo := repository.NewVariantRepository(db)
	mediaRepo := repository.NewMediaRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)

	productService := service.NewProductService(productRepo, brandRepo, cacheRepo, categoryRepo, esRepo)
	categoryService := service.NewCategoryService(categoryRepo, cacheRepo)
	brandService := service.NewBrandService(brandRepo, cacheRepo)
	attributeService := service.NewAttributeService(attributeRepo)
	warehouseService := service.NewWarehouseService(warehouseRepo)
	variantService := service.NewVariantService(variantRepo, productRepo, esRepo, cacheRepo)
	mediaService := service.NewMediaService(mediaRepo, cfg.BaseUrl)
	inventoryService := service.NewInventoryService(inventoryRepo)

	productCtrl := controller.NewProductController(productService)
	categoryCtrl := controller.NewCategoryController(categoryService)
	brandCtrl := controller.NewBrandController(brandService)
	variantCtrl := controller.NewVariantController(variantService)
	attributeCtrl := controller.NewAttributeController(attributeService)
	warehouseCtrl := controller.NewWarehouseController(warehouseService)
	inventoryCtrl := controller.NewInventoryController(inventoryService)
	mediaCtrl := controller.NewMediaController(mediaService)

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api/v1")
	{
		api.POST("/media/upload", mediaCtrl.Upload)

		// --- Product Routes ---
		api.POST("/products", productCtrl.Create)
		api.GET("/products/search", productCtrl.SearchProduct)
		api.GET("/products", productCtrl.GetProductList)
		api.GET("/products/:id", productCtrl.GetProductById)
		api.PUT("/products/:id", productCtrl.Update)
		api.DELETE("/products/:id", productCtrl.Delete)

		// --- Variant Routes ---
		api.POST("/variants", variantCtrl.Create)
		api.GET("/variants/:id", variantCtrl.GetVariantByID)
		api.PUT("/variants/:id", variantCtrl.UpdateVariant)
		api.DELETE("/variants/:id", variantCtrl.DeleteVariant)

		// --- Category Routes ---
		api.POST("/categories", categoryCtrl.CreateCategory)
		api.GET("/categories", categoryCtrl.List)
		api.PUT("/categories/:id", categoryCtrl.Update)
		api.DELETE("/categories/:id", categoryCtrl.Delete)

		// --- Brand Routes ---
		api.POST("/brands", brandCtrl.CreateBrand)
		api.GET("/brands", brandCtrl.ListBrand)
		api.GET("/brands/:id", brandCtrl.GetBrandById)
		api.PUT("/brands/:id", brandCtrl.UpdateBrand)
		api.DELETE("/brands/:id", brandCtrl.DeleteBrand)

		// --- Attribute Routes ---
		api.POST("/attributes", attributeCtrl.Create)
		api.POST("/attributes/values", attributeCtrl.AddValue)
		api.GET("/attributes", attributeCtrl.ListAttributes)
		api.GET("/attributes/:id", attributeCtrl.GetByID)
		api.PUT("/attributes/:id", attributeCtrl.Update)
		api.DELETE("/attributes/:id", attributeCtrl.Delete)

		// --- Warehouse Routes ---
		api.POST("/warehouses", warehouseCtrl.Create)
		api.GET("/warehouses", warehouseCtrl.ListAll)
		api.GET("/warehouses/:id", warehouseCtrl.GetWarehouseById)
		api.PUT("/warehouses/:id", warehouseCtrl.Update)
		api.DELETE("/warehouses/:id", warehouseCtrl.Delete)

		// --- Inventory Routes ---
		api.POST("/inventory/adjust", inventoryCtrl.AdjustStock)
		api.GET("/inventory/variant/:id", inventoryCtrl.GetStockByVariant)

	}
}
