package main

import (
	"Product_Mangement_Api/config"
	"Product_Mangement_Api/entity"
	"Product_Mangement_Api/router"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// 1. Load Config
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Unable to load config:", err)
	}

	db := config.ConnectDB(&cfg)
	db.AutoMigrate(
		&entity.Product{},
		&entity.Brand{},
		&entity.Category{},
		&entity.Inventory{},
		&entity.Warehouse{},
		&entity.Media{},
		&entity.ProductVariant{},
		&entity.VariantOption{},
		&entity.VariantOptionValue{},
		&entity.StockMovement{},
		&entity.Attribute{},
		&entity.ProductAttribute{},
		&entity.ProductAttributeValue{},
	)
	r := gin.Default()

	router.SetupRouter(r, db, &cfg)

	serverAddress := cfg.ServerAddress
	if serverAddress == "" {
		serverAddress = ":8080"
	}

	log.Printf("Server đang chạy tại %s", serverAddress)
	if err := r.Run(serverAddress); err != nil {
		log.Fatal("Server startup error:", err)
	}
}
