
# Product Management System API

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
![Database](https://img.shields.io/badge/PostgreSQL-16-316192?style=flat&logo=postgresql)
![Cache](https://img.shields.io/badge/Redis-7-DC382D?style=flat&logo=redis)
![Search](https://img.shields.io/badge/Elasticsearch-7.10-005571?style=flat&logo=elasticsearch)
![License](https://img.shields.io/badge/License-MIT-green?style=flat)

## 📖 Overview

The **Product Management System API** is a high-performance, enterprise-grade backend solution designed for modern e-commerce platforms. Built with **Golang (Gin Framework)**, this system orchestrates complex product lifecycles, advanced inventory management across multiple warehouses, and blazing-fast search capabilities.

Key architectural highlights include **Elasticsearch** integration for full-text search, **Redis** for multi-layer caching, and a robust **PostgreSQL** schema handling dynamic product variants and internationalization (i18n).

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

```

---

## ⚡ Key Features

* **Advanced Catalog Management:**
* Support for hierarchical **Categories** and **Brands**.
* **Dynamic Attributes** (e.g., Color, Size, Material) configurable per product line.
* **Multi-language Support (i18n):** Native JSONB storage for English and Vietnamese content.


* **Product Variants & SKU System:**
* Flexible parent-child product relationships.
* Individual tracking of SKUs, prices, weights, and media for each variant.


* **Smart Inventory Management:**
* **Multi-Warehouse:** Track stock levels across different physical locations.
* **Stock Movements:** Audit trail for all inbound, outbound, and adjustment transactions.
* **Reservation Logic:** Support for reserved quantities during checkout.


* **High-Performance Search:**
* Integrated **Elasticsearch** for typo-tolerant, full-text search.
* Real-time synchronization between PostgreSQL and Elasticsearch via Goroutines.


* **Performance Optimization:**
* **Redis Caching** strategy for high-traffic endpoints (Product Details, Listings).
* Optimized Database indexing and query planning.



---

## 📂 Database Schema (ER Diagram)

The following Entity Relationship Diagram (ERD) illustrates the normalized database structure:

```mermaid
erDiagram
    PRODUCT ||--|{ PRODUCT_VARIANT : "has variants"
    PRODUCT }o--|| CATEGORY : "classified in"
    PRODUCT }o--|| BRAND : "manufactured by"
    PRODUCT ||--o{ MEDIA : "represented by"

    PRODUCT_VARIANT ||--o{ INVENTORY : "stocked in"
    PRODUCT_VARIANT }o--o{ VARIANT_ATTRIBUTE_VALUE : "defined by"
    
    INVENTORY }o--|| WAREHOUSE : "located at"
    INVENTORY ||--o{ STOCK_MOVEMENT : "audited by"

    ATTRIBUTE ||--|{ ATTRIBUTE_VALUE : "contains"
    VARIANT_ATTRIBUTE_VALUE }o--|| ATTRIBUTE_VALUE : "references"

    PRODUCT {
        uuid id PK
        jsonb name "i18n"
        string slug
        uuid category_id FK
        uuid brand_id FK
        string status
    }

    PRODUCT_VARIANT {
        uuid id PK
        uuid product_id FK
        string sku
        decimal price
        decimal sale_price
    }

    INVENTORY {
        uuid id PK
        uuid variant_id FK
        uuid warehouse_id FK
        int quantity
    }

```

---

## 🛠 Tech Stack

| Component | Technology | Description |
| --- | --- | --- |
| **Language** | [Go (Golang)](https://go.dev/) | Version 1.22+ |
| **Framework** | [Gin Gonic](https://github.com/gin-gonic/gin) | High-performance HTTP web framework |
| **Database** | [PostgreSQL](https://www.postgresql.org/) | Primary relational database |
| **ORM** | [GORM](https://gorm.io/) | The fantastic ORM library for Golang |
| **Caching** | [Redis](https://redis.io/) | In-memory data structure store |
| **Search Engine** | [Elasticsearch](https://www.elastic.co/) | Distributed search and analytics engine |
| **Config** | [Viper](https://github.com/spf13/viper) | Configuration management |
| **UUID** | [Google UUID](https://github.com/google/uuid) | Universally Unique Identifier generation |

---

## 📂 Project Structure

```text
.
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

```

---

## 🚀 Getting Started

### Prerequisites

* Go 1.22 or higher
* Docker & Docker Compose (Recommended)
* PostgreSQL, Redis, Elasticsearch instances running

### Installation

1. **Clone the repository**
```bash
git clone [https://github.com/your-username/product-management-api.git](https://github.com/your-username/product-management-api.git)
cd product-management-api

```


2. **Environment Configuration**
Create an `app.env` file in the root directory:
```properties
# Database Configuration
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
```


3. **Install Dependencies**
```bash
go mod tidy

```


4. **Run the Application**
```bash
go run main.go

```


5. **Seed Initial Data (Optional)**
Populate the database with sample categories, brands, warehouses, and products:
```bash
go run cmd/seed/main.go

## 🔌 API Documentation

### Products

* `GET /api/v1/products`: List products with pagination.
* `GET /api/v1/products/:id`: Get detailed product info.
* `POST /api/v1/products`: Create a new product with variants.
* `GET /api/v1/products/search`: Full-text search via Elasticsearch.

### Inventory

* `GET /api/v1/inventory/variant/:id`: Check stock for a specific SKU.
* `POST /api/v1/inventory/adjust`: Adjust stock (Inbound/Outbound/Correction).

### Categories & Brands

* `GET /api/v1/categories`: Retrieve category tree.
* `POST /api/v1/brands`: Create a new brand.

### Media

* `POST /api/v1/media/upload`: Upload images or videos for products.


## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request
