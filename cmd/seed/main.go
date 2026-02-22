package main

import (
	"Product_Mangement_Api/config"
	"Product_Mangement_Api/entity"
	"Product_Mangement_Api/models"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

func toJson(data interface{}) datatypes.JSON {
	bytes, _ := json.Marshal(data)
	return datatypes.JSON(bytes)
}

func multiLang(vi, en string) map[string]interface{} {
	return map[string]interface{}{"vi": vi, "en": en}
}

func main() {
	cfg, err := config.LoadConfig("../../")
	if err != nil {
		log.Fatal("Không thể load config:", err)
	}

	db := config.ConnectDB(&cfg)
	log.Println("Bắt đầu khởi tạo dữ liệu mẫu (Seed Data)...")

	hcmWarehouse := entity.Warehouse{Name: toJson(multiLang("Kho Hồ Chí Minh", "Ho Chi Minh Warehouse")), Address: "123 Lê Lợi, Q1", Phone: "0909123456", Status: models.StatusActive}
	hnWarehouse := entity.Warehouse{Name: toJson(multiLang("Kho Hà Nội", "Ha Noi Warehouse")), Address: "456 Kim Mã, Ba Đình", Phone: "0909654321", Status: models.StatusActive}
	db.Create(&hcmWarehouse)
	db.Create(&hnWarehouse)

	apple := entity.Brand{Name: toJson(multiLang("Apple", "Apple")), Logo: "apple.png"}
	db.Create(&apple)

	electronics := entity.Category{Name: toJson(multiLang("Điện tử", "Electronics"))}
	db.Create(&electronics)
	phones := entity.Category{Name: toJson(multiLang("Điện thoại", "Smartphones")), ParentId: &electronics.ID}
	db.Create(&phones)

	colorAttr := entity.Attribute{Name: toJson(multiLang("Màu sắc", "Color")), Type: "TEXT", IsFilterable: true}
	storageAttr := entity.Attribute{Name: toJson(multiLang("Dung lượng", "Storage")), Type: "TEXT", IsFilterable: true}
	db.Create(&colorAttr)
	db.Create(&storageAttr)

	iPhone15 := entity.Product{
		Name:        toJson(multiLang("iPhone 15 Pro", "iPhone 15 Pro")),
		Slug:        "iphone-15-pro-" + uuid.New().String()[:8],
		Description: toJson(multiLang("Điện thoại Apple mới nhất", "The latest Apple phone")),
		CategoryID:  phones.ID,
		BrandID:     apple.ID,
		Status:      models.StatusActive,

		ProductAttributes: []entity.ProductAttribute{
			{
				AttributeID: colorAttr.ID,
				Values: []entity.ProductAttributeValue{
					{Value: toJson(multiLang("Đen", "Black"))},
					{Value: toJson(multiLang("Trắng", "White"))},
				},
			},
			{
				AttributeID: storageAttr.ID,
				Values: []entity.ProductAttributeValue{
					{Value: toJson(multiLang("128GB", "128GB"))},
					{Value: toJson(multiLang("256GB", "256GB"))},
				},
			},
		},

		Variants: []entity.ProductVariant{
			{
				Code:   "IP15P-VN",
				Name:   toJson(multiLang("Bản VN/A", "VN/A Version")),
				Status: models.StatusActive,

				Options: []entity.VariantOption{
					{
						SKU:       "IP15P-BLK-128",
						Price:     25000000,
						SalePrice: 24000000,
						Weight:    0.2,
						Status:    models.StatusActive,

						Values: []entity.VariantOptionValue{
							{Name: "Color", Value: "Black", SortOrder: 1},
							{Name: "Storage", Value: "128GB", SortOrder: 2},
						},

						Inventories: []entity.Inventory{
							{WarehouseID: hcmWarehouse.ID, Quantity: 50},
							{WarehouseID: hnWarehouse.ID, Quantity: 30},
						},
					},
					{
						SKU:       "IP15P-WHT-256",
						Price:     30000000,
						SalePrice: 29000000,
						Weight:    0.2,
						Status:    models.StatusActive,
						Values: []entity.VariantOptionValue{
							{Name: "Color", Value: "White", SortOrder: 1},
							{Name: "Storage", Value: "256GB", SortOrder: 2},
						},
						Inventories: []entity.Inventory{
							{WarehouseID: hcmWarehouse.ID, Quantity: 20},
						},
					},
				},
			},
		},
	}

	if err := db.Create(&iPhone15).Error; err != nil {
		log.Printf("Lỗi tạo product: %v", err)
	} else {
		log.Println("-> Đã tạo Product thành công với cấu trúc mới!")
	}

	log.Println("=== HOÀN TẤT SEED DATA ===")
}
