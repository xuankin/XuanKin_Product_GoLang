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

// Helper để chuyển đổi map sang JSON
func toJson(data interface{}) datatypes.JSON {
	bytes, _ := json.Marshal(data)
	return datatypes.JSON(bytes)
}

// Hàm hỗ trợ tạo tên đa ngôn ngữ
func multiLang(vi, en string) map[string]interface{} {
	return map[string]interface{}{
		"vi": vi,
		"en": en,
	}
}

func main() {
	// 1. Load Config (Giả sử chạy từ root folder, file app.env ở root)
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Không thể load config:", err)
	}

	// 2. Kết nối DB
	db := config.ConnectDB(&cfg)
	log.Println("Bắt đầu khởi tạo dữ liệu mẫu (Seed Data)...")

	// 3. Xóa dữ liệu cũ (Tùy chọn: cẩn thận khi dùng trên prod)
	// db.Exec("TRUNCATE TABLE products, categories, brands, attributes, warehouses, inventories, product_variants, media CASCADE")

	// 4. Tạo Warehouses (Kho hàng)
	hcmWarehouse := entity.Warehouse{
		Name:    toJson(multiLang("Kho Hồ Chí Minh", "Ho Chi Minh Warehouse")),
		Address: "123 Đường Lê Lợi, Q1, TP.HCM",
		Phone:   "0909123456",
		Status:  models.StatusActive,
	}
	hnWarehouse := entity.Warehouse{
		Name:    toJson(multiLang("Kho Hà Nội", "Ha Noi Warehouse")),
		Address: "456 Đường Kim Mã, Ba Đình, Hà Nội",
		Phone:   "0909654321",
		Status:  models.StatusActive,
	}
	db.Create(&hcmWarehouse)
	db.Create(&hnWarehouse)
	log.Println("-> Đã tạo Warehouses")

	// 5. Tạo Brands (Thương hiệu)
	apple := entity.Brand{Name: toJson(multiLang("Apple", "Apple")), Logo: "https://upload.wikimedia.org/wikipedia/commons/f/fa/Apple_logo_black.svg"}
	samsung := entity.Brand{Name: toJson(multiLang("Samsung", "Samsung")), Logo: "https://upload.wikimedia.org/wikipedia/commons/2/24/Samsung_Logo.svg"}
	nike := entity.Brand{Name: toJson(multiLang("Nike", "Nike")), Logo: "https://upload.wikimedia.org/wikipedia/a/a6/Logo_NIKE.svg"}

	db.Create(&apple)
	db.Create(&samsung)
	db.Create(&nike)
	log.Println("-> Đã tạo Brands")

	// 6. Tạo Categories (Danh mục)
	// Danh mục cha
	electronics := entity.Category{Name: toJson(multiLang("Điện tử", "Electronics"))}
	fashion := entity.Category{Name: toJson(multiLang("Thời trang", "Fashion"))}
	db.Create(&electronics)
	db.Create(&fashion)

	// Danh mục con
	phones := entity.Category{Name: toJson(multiLang("Điện thoại", "Smartphones")), ParentId: &electronics.ID}
	laptops := entity.Category{Name: toJson(multiLang("Máy tính xách tay", "Laptops")), ParentId: &electronics.ID}
	shoes := entity.Category{Name: toJson(multiLang("Giày dép", "Shoes")), ParentId: &fashion.ID}

	db.Create(&phones)
	db.Create(&laptops)
	db.Create(&shoes)
	log.Println("-> Đã tạo Categories")

	// 7. Tạo Attributes (Thuộc tính)
	// Màu sắc
	colorAttr := entity.Attribute{Name: toJson(multiLang("Màu sắc", "Color"))}
	db.Create(&colorAttr)

	// Giá trị màu
	colorBlack := entity.AttributeValue{AttributeID: colorAttr.ID, Value: toJson(multiLang("Đen", "Black"))}
	colorWhite := entity.AttributeValue{AttributeID: colorAttr.ID, Value: toJson(multiLang("Trắng", "White"))}
	colorBlue := entity.AttributeValue{AttributeID: colorAttr.ID, Value: toJson(multiLang("Xanh", "Blue"))}
	db.Create(&colorBlack)
	db.Create(&colorWhite)
	db.Create(&colorBlue)

	storageAttr := entity.Attribute{Name: toJson(multiLang("Dung lượng", "Storage"))}
	db.Create(&storageAttr)

	storage128 := entity.AttributeValue{AttributeID: storageAttr.ID, Value: toJson(multiLang("128GB", "128GB"))}
	storage256 := entity.AttributeValue{AttributeID: storageAttr.ID, Value: toJson(multiLang("256GB", "256GB"))}
	db.Create(&storage128)
	db.Create(&storage256)
	log.Println("-> Đã tạo Attributes & Values")

	iPhone15 := entity.Product{
		Name:        toJson(multiLang("iPhone 15 Pro", "iPhone 15 Pro")),
		Slug:        "iphone-15-pro-" + uuid.New().String()[:8],
		Description: toJson(multiLang("Điện thoại Apple mới nhất", "The latest Apple phone")),
		CategoryID:  phones.ID,
		BrandID:     apple.ID,
		Status:      models.StatusActive,
	}

	v1 := entity.ProductVariant{
		SKU:       "IP15-BLK-128",
		Price:     25000000,
		SalePrice: 24000000,
		Weight:    0.2,
		Status:    models.StatusActive,
	}

	v2 := entity.ProductVariant{
		SKU:       "IP15-WHT-256",
		Price:     30000000,
		SalePrice: 29000000,
		Weight:    0.2,
		Status:    models.StatusActive,
	}

	iPhone15.Variants = []entity.ProductVariant{v1, v2}

	if err := db.Create(&iPhone15).Error; err != nil {
		log.Printf("Lỗi tạo product: %v", err)
	} else {

		var savedVariants []entity.ProductVariant
		db.Where("product_id = ?", iPhone15.ID).Find(&savedVariants)

		for _, v := range savedVariants {
			if v.SKU == "IP15-BLK-128" {

				db.Model(&v).Association("Attributes").Append(&colorBlack, &storage128)

				db.Create(&entity.Inventory{VariantID: v.ID, WarehouseID: hcmWarehouse.ID, Quantity: 50})
				db.Create(&entity.Inventory{VariantID: v.ID, WarehouseID: hnWarehouse.ID, Quantity: 30})
			} else if v.SKU == "IP15-WHT-256" {

				db.Model(&v).Association("Attributes").Append(&colorWhite, &storage256)

				db.Create(&entity.Inventory{VariantID: v.ID, WarehouseID: hcmWarehouse.ID, Quantity: 20})
			}
		}
	}

	// 9. Tạo Product thứ 2 - Samsung S24
	samsungS24 := entity.Product{
		Name:        toJson(multiLang("Samsung Galaxy S24", "Samsung Galaxy S24")),
		Slug:        "samsung-s24-" + uuid.New().String()[:8],
		Description: toJson(multiLang("Siêu phẩm AI", "AI Phone")),
		CategoryID:  phones.ID,
		BrandID:     samsung.ID,
		Status:      models.StatusActive,
	}
	db.Create(&samsungS24)

	vSamsung := entity.ProductVariant{
		ProductID: samsungS24.ID,
		SKU:       "SS-S24-BLU",
		Price:     20000000,
		Status:    models.StatusActive,
	}
	db.Create(&vSamsung)
	db.Model(&vSamsung).Association("Attributes").Append(&colorBlue, &storage256)
	db.Create(&entity.Inventory{VariantID: vSamsung.ID, WarehouseID: hnWarehouse.ID, Quantity: 100})

	log.Println("-> Đã tạo Products & Variants")
	log.Println("=== HOÀN TẤT SEED DATA ===")
}
