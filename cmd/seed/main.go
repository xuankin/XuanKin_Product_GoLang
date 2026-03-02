package main

import (
	"Product_Mangement_Api/config"
	"Product_Mangement_Api/entity"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Helper function tạo JSON đa ngôn ngữ
func multiLang(vi, en string) datatypes.JSON {
	data := map[string]string{
		"vi": vi,
		"en": en,
	}
	bytes, _ := json.Marshal(data)
	return datatypes.JSON(bytes)
}

func main() {
	// Load config từ thư mục gốc
	cfg, err := config.LoadConfig("../../")
	if err != nil {
		log.Fatal("Không thể load config: ", err)
	}

	// Kết nối database
	db := config.ConnectDB(&cfg)

	fmt.Println("--- Đang bắt đầu seed data đa ngôn ngữ ---")

	// 1. Seed Warehouses (Kho hàng)
	whHCM := entity.Warehouse{
		Name:    multiLang("Kho Hồ Chí Minh", "Ho Chi Minh Warehouse"),
		Address: "Quận 1, TP.HCM",
		Phone:   "0909123456",
		Status:  "ACTIVE",
	}
	whHN := entity.Warehouse{
		Name:    multiLang("Kho Hà Nội", "Ha Noi Warehouse"),
		Address: "Cầu Giấy, Hà Nội",
		Phone:   "0909654321",
		Status:  "ACTIVE",
	}
	db.Where("address = ?", whHCM.Address).FirstOrCreate(&whHCM)
	db.Where("address = ?", whHN.Address).FirstOrCreate(&whHN)

	// 2. Seed Brands (Thương hiệu)
	brandApple := entity.Brand{Name: multiLang("Apple", "Apple"), Logo: "apple_logo.png"}
	brandSamsung := entity.Brand{Name: multiLang("Samsung", "Samsung"), Logo: "samsung_logo.png"}
	brandDell := entity.Brand{Name: multiLang("Dell", "Dell"), Logo: "dell_logo.png"}

	db.Where("name @> ?", `{"vi": "Apple"}`).FirstOrCreate(&brandApple)
	db.Where("name @> ?", `{"vi": "Samsung"}`).FirstOrCreate(&brandSamsung)
	db.Where("name @> ?", `{"vi": "Dell"}`).FirstOrCreate(&brandDell)

	// 3. Seed Categories (Danh mục)
	catElectronics := entity.Category{Name: multiLang("Thiết bị điện tử", "Electronics")}
	db.Where("name @> ?", `{"vi": "Thiết bị điện tử"}`).FirstOrCreate(&catElectronics)

	catPhone := entity.Category{Name: multiLang("Điện thoại", "Smartphones"), ParentId: &catElectronics.ID}
	catLaptop := entity.Category{Name: multiLang("Máy tính xách tay", "Laptops"), ParentId: &catElectronics.ID}

	db.Where("name @> ?", `{"vi": "Điện thoại"}`).FirstOrCreate(&catPhone)
	db.Where("name @> ?", `{"vi": "Máy tính xách tay"}`).FirstOrCreate(&catLaptop)

	// 4. Seed Attributes (Thuộc tính dùng chung)
	attrColor := entity.Attribute{Name: multiLang("Màu sắc", "Color"), Type: "TEXT", IsFilterable: true, IsRequired: true}
	attrStorage := entity.Attribute{Name: multiLang("Dung lượng", "Storage"), Type: "TEXT", IsFilterable: true, IsRequired: true}
	attrRAM := entity.Attribute{Name: multiLang("RAM", "RAM"), Type: "TEXT", IsFilterable: true, IsRequired: true}

	db.Where("name @> ?", `{"vi": "Màu sắc"}`).FirstOrCreate(&attrColor)
	db.Where("name @> ?", `{"vi": "Dung lượng"}`).FirstOrCreate(&attrStorage)
	db.Where("name @> ?", `{"vi": "RAM"}`).FirstOrCreate(&attrRAM)

	// 5. Seed Product 1: iPhone 15 Pro (Apple - Điện thoại)
	p1 := entity.Product{
		Name:        multiLang("iPhone 15 Pro Max", "iPhone 15 Pro Max"),
		Description: multiLang("Khung Titan bền bỉ", "Strong Titanium frame"),
		Slug:        "iphone-15-pro-max-" + uuid.New().String()[:5],
		BrandID:     brandApple.ID,
		CategoryID:  catPhone.ID,
		Status:      "ACTIVE",
	}
	seedProductWithVariants(db, &p1, whHCM.ID, whHN.ID)

	// 6. Seed Product 2: Samsung S24 Ultra (Samsung - Điện thoại)
	p2 := entity.Product{
		Name:        multiLang("Samsung Galaxy S24 Ultra", "Samsung Galaxy S24 Ultra"),
		Description: multiLang("Công nghệ AI đỉnh cao", "Ultimate AI Technology"),
		Slug:        "samsung-s24-ultra-" + uuid.New().String()[:5],
		BrandID:     brandSamsung.ID,
		CategoryID:  catPhone.ID,
		Status:      "ACTIVE",
	}
	seedProductWithVariants(db, &p2, whHCM.ID, whHN.ID)

	// 7. Seed Product 3: Dell XPS 15 (Dell - Laptop)
	p3 := entity.Product{
		Name:        multiLang("Dell XPS 15 9530", "Dell XPS 15 9530"),
		Description: multiLang("Laptop đồ họa mỏng nhẹ", "Thin and light graphic laptop"),
		Slug:        "dell-xps-15-9530-" + uuid.New().String()[:5],
		BrandID:     brandDell.ID,
		CategoryID:  catLaptop.ID,
		Status:      "ACTIVE",
	}
	seedProductWithVariants(db, &p3, whHN.ID, whHCM.ID)

	fmt.Println("--- Seed data hoàn tất! Hãy kiểm tra trong Database ---")
}

// Hàm hỗ trợ tạo luồng: Product -> Variant -> Option -> Values -> Inventory
func seedProductWithVariants(db *gorm.DB, p *entity.Product, whPrimaryID uuid.UUID, whSecondaryID uuid.UUID) {
	// Lưu Product
	if err := db.Where("slug = ?", p.Slug).FirstOrCreate(p).Error; err != nil {
		fmt.Println("Lỗi tạo SP:", err)
		return
	}

	// Tạo Variant (Ví dụ: Bản Quốc tế)
	variant := entity.ProductVariant{
		ProductID: p.ID,
		Code:      "VAR-" + uuid.New().String()[:8],
		Name:      multiLang("Bản Chính Hãng", "Official Version"),
		Status:    "ACTIVE",
	}
	db.Where(entity.ProductVariant{ProductID: p.ID, Code: variant.Code}).FirstOrCreate(&variant)

	// Tạo Option 1 (Màu Đen / 256GB)
	opt1 := entity.VariantOption{
		VariantID: variant.ID,
		SKU:       "SKU-" + uuid.New().String()[:8],
		Price:     30000000,
		SalePrice: 29500000,
		Weight:    250,
		Status:    "ACTIVE",
	}
	db.Create(&opt1)

	// Lưu Option Values cho Option 1
	db.Create(&entity.VariantOptionValue{OptionID: opt1.ID, Name: "Color", Value: "Black", SortOrder: 1})
	db.Create(&entity.VariantOptionValue{OptionID: opt1.ID, Name: "Storage", Value: "256GB", SortOrder: 2})

	// Tồn kho cho Option 1
	db.Create(&entity.Inventory{OptionID: opt1.ID, WarehouseID: whPrimaryID, Quantity: 100})

	// Tạo Option 2 (Màu Trắng / 512GB)
	opt2 := entity.VariantOption{
		VariantID: variant.ID,
		SKU:       "SKU-" + uuid.New().String()[:8],
		Price:     35000000,
		SalePrice: 34000000,
		Weight:    250,
		Status:    "ACTIVE",
	}
	db.Create(&opt2)

	// Lưu Option Values cho Option 2
	db.Create(&entity.VariantOptionValue{OptionID: opt2.ID, Name: "Color", Value: "White", SortOrder: 1})
	db.Create(&entity.VariantOptionValue{OptionID: opt2.ID, Name: "Storage", Value: "512GB", SortOrder: 2})

	// Tồn kho cho Option 2 ở cả 2 kho
	db.Create(&entity.Inventory{OptionID: opt2.ID, WarehouseID: whPrimaryID, Quantity: 50})
	db.Create(&entity.Inventory{OptionID: opt2.ID, WarehouseID: whSecondaryID, Quantity: 30})
}
