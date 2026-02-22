Dưới đây là file README.md đã được căn chỉnh lại toàn bộ formating, sửa lỗi hiển thị và cập nhật chuẩn xác theo ERD mới nhất của bạn.Bạn chỉ cần bấm nút Copy code ở góc trên cùng bên phải của khung dưới đây và dán đè toàn bộ vào file README.md của dự án là xong nhé:Markdown# Product Management System API

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
![Database](https://img.shields.io/badge/PostgreSQL-16-316192?style=flat&logo=postgresql)
![Cache](https://img.shields.io/badge/Redis-7-DC382D?style=flat&logo=redis)
![Search](https://img.shields.io/badge/Elasticsearch-7.10-005571?style=flat&logo=elasticsearch)
![License](https://img.shields.io/badge/License-MIT-green?style=flat)

## 📖 Overview

The **Product Management System API** is a high-performance, enterprise-grade backend solution designed for modern e-commerce platforms. Built with **Golang (Gin Framework)**, this system orchestrates complex product lifecycles, advanced inventory management across multiple warehouses, and blazing-fast search capabilities.

Key architectural highlights include **Elasticsearch** integration for full-text search, **Redis** for multi-layer caching, and a robust **PostgreSQL** schema handling a deep, 3-tier product structure (Product -> Variant -> Option) with internationalization (i18n).

---

## 🏗 System Architecture

The project adheres to **Clean Architecture** principles, ensuring separation of concerns between the transport layer, business logic, and data access.

```mermaid
graph TD
    subgraph Client_Side
        Client[Web/Mobile Client]
        Postman[Postman/Tester]
    end

    subgraph API_Gateway
        Gin[Gin Web Server]
        Middleware[Middleware\n(CORS, Logger, Recovery)]
    end

    subgraph Business_Logic
        Controller[Controllers]
        Service[Services]
        Repo[Repositories]
    end

    subgraph Infrastructure
        PG[(PostgreSQL\nPrimary DB)]
        Redis[(Redis\nCache Store)]
        ES[(Elasticsearch\nSearch Engine)]
    end

    Client -->|HTTP/REST| Gin
    Postman -->|HTTP/REST| Gin
    Gin --> Middleware --> Controller
    Controller --> Service
    Service --> Repo

    Repo -->|ORM/GORM| PG
    Repo -->|Cache Hit/Miss| Redis
    Repo -->|Indexing/Search| ES

    %% Async Sync Process
    Service -.->|Async Routine| ES
⚡ Key FeaturesAdvanced Catalog Management:Support for hierarchical Categories and Brands.Dynamic Global Attributes (e.g., Material, Style) configurable per product line via PRODUCT_ATTRIBUTE.Multi-language Support (i18n): Native JSONB storage for English and Vietnamese content.Deep 3-Tier Product Architecture:Product: Base information (Name, Description, Brand, Category).Variant: Groupings (e.g., "iPhone 15 Pro VN/A").Option: Specific sellable items with individual SKUs, Prices, and Weights (e.g., "Color: Black, Storage: 128GB").Smart Inventory Management:Option-Level Tracking: Inventory is strictly tied to specific VariantOption IDs.Multi-Warehouse: Track stock levels across different physical locations.Stock Movements: Audit trail for all inbound, outbound, and adjustment transactions (STOCK_MOVEMENT).Reservation Logic: Support for reserved quantities during checkout.High-Performance Search:Integrated Elasticsearch for typo-tolerant, full-text search.Real-time synchronization between PostgreSQL and Elasticsearch.Performance Optimization:Redis Caching strategy for high-traffic endpoints (Product Details, Listings).Optimized Database indexing and cascade deletions.📂 Database Schema (ER Diagram)The following Entity Relationship Diagram (ERD) illustrates the highly normalized database structure used in this project:Đoạn mãerDiagram
    PRODUCT ||--|{ PRODUCT_VARIANT : has
    PRODUCT }o--|| CATEGORY : belongs_to
    PRODUCT }o--|| BRAND : belongs_to
    PRODUCT ||--o{ MEDIA : has
    PRODUCT ||--o{ PRODUCT_ATTRIBUTE : has

    PRODUCT_ATTRIBUTE }o--|| ATTRIBUTE : defines
    PRODUCT_ATTRIBUTE ||--o{ PRODUCT_ATTRIBUTE_VALUE : has

    PRODUCT_VARIANT ||--o{ VARIANT_OPTION : has
    VARIANT_OPTION ||--o{ VARIANT_OPTION_VALUE : has

    VARIANT_OPTION ||--o{ INVENTORY : stored_in
    INVENTORY }o--|| WAREHOUSE : belongs_to
    INVENTORY ||--o{ STOCK_MOVEMENT : tracked_by

    PRODUCT {
        uuid id PK
        jsonb name
        string slug
        uuid category_id FK
        uuid brand_id FK
    }

    PRODUCT_VARIANT {
        uuid id PK
        uuid product_id FK
        string code
    }

    VARIANT_OPTION {
        uuid id PK
        uuid variant_id FK
        string sku
        decimal price
    }
    
    INVENTORY {
        uuid id PK
        uuid option_id FK
        uuid warehouse_id FK
        integer quantity
    }
🛠 Tech StackComponentTechnologyDescriptionLanguageGo (Golang)Version 1.22+FrameworkGin GonicHigh-performance HTTP web frameworkDatabasePostgreSQLPrimary relational databaseORMGORMThe fantastic ORM library for GolangCachingRedisIn-memory data structure storeSearch EngineElasticsearchDistributed search and analytics engineConfigViperConfiguration managementUUIDGoogle UUIDUniversally Unique Identifier generation📂 Project StructurePlaintext.
├── cmd/
│   └── seed/            # Database seeding script
├── config/              # Configuration (DB, Redis, ES, Viper)
├── controller/          # HTTP Handlers (Gin)
├── entity/              # Database Models (GORM structs)
├── models/              # DTOs (Request/Response structs)
├── repository/          # Data Access Layer (DAL)
├── router/              # API Routes & Middleware setup
├── service/             # Business Logic Layer
├── app.env              # Environment variables
├── main.go              # Application entry point
└── go.mod               # Go module definition
🚀 Getting StartedPrerequisitesGo 1.22 or higherDocker & Docker Compose (Recommended)PostgreSQL, Redis, Elasticsearch instances runningInstallationClone the repositoryBashgit clone [https://github.com/your-username/product-management-api.git](https://github.com/your-username/product-management-api.git)
cd product-management-api
Environment ConfigurationCreate an app.env file in the root directory:Properties# Database Configuration
DB_SOURCE=postgres://user:password@localhost:5432/product_db?sslmode=disable

# Server Configuration
SERVER_ADDRESS=:8080
BASE_URL=http://localhost:8080

# Redis Configuration
REDIS_ADDRESS=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Elasticsearch Configuration
ELASTICSEARCH_ADDRESS=http://localhost:9200
Install DependenciesBashgo mod tidy
Run the ApplicationBashgo run main.go
Seed Initial Data (Optional)Populate the database with sample categories, brands, warehouses, and a complete 3-tier product structure (e.g., iPhone 15 Pro):Bashgo run cmd/seed/main.go
🔌 API Documentation (v1)Products & AttributesGET /api/v1/products: List products with pagination.GET /api/v1/products/search: Full-text search via Elasticsearch.GET /api/v1/products/:id: Get detailed product info.POST /api/v1/products: Create a new product.PUT /api/v1/products/:id: Update product base info.DELETE /api/v1/products/:id: Delete product (Cascade deletes variants, options, media).POST /api/v1/products/:id/attributes: Sync specific attributes to a product.Variants & OptionsPOST /api/v1/variants: Create a new variant group.POST /api/v1/variants/:id/options: Add a sellable configuration (Option) to a variant.PUT /api/v1/options/:optionId: Update Option details (Price, SKU, Weight).DELETE /api/v1/variants/:id: Delete a variant group.Inventory (Option-based)GET /api/v1/inventory/option/:id: Check stock for a specific Option across warehouses.POST /api/v1/inventory/adjust: Adjust stock (Inbound/Outbound/Correction) and auto-create Stock Movement logs.Master Data (Categories, Brands, Attributes, Warehouses)Standard CRUD operations available at:/api/v1/categories/api/v1/brands/api/v1/attributes/api/v1/warehousesMediaPOST /api/v1/media/upload: Upload images/videos and map them to Product, Variant, or Option (Supports automated Primary Image reset).
