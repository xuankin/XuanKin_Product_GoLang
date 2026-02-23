
# 🚀 Product Management System API

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat\&logo=go)
![Database](https://img.shields.io/badge/PostgreSQL-16-316192?style=flat\&logo=postgresql)
![Cache](https://img.shields.io/badge/Redis-7-DC382D?style=flat\&logo=redis)
![Search](https://img.shields.io/badge/Elasticsearch-7.10-005571?style=flat\&logo=elasticsearch)
![License](https://img.shields.io/badge/License-MIT-green?style=flat)

---

## 📖 Overview

The **Product Management System API** is a high-performance, enterprise-grade backend solution designed for modern e-commerce platforms.

Built with **Golang (Gin Framework)**, this system manages:

* Complex product lifecycles
* Multi-warehouse inventory
* High-performance search
* Deep 3-tier product architecture

### 🔥 Key Integrations

* 🔎 **Elasticsearch** – Full-text search
* ⚡ **Redis** – Multi-layer caching
* 🐘 **PostgreSQL** – Primary relational database
* 🧱 Clean Architecture

---

# 🏗 System Architecture

The project follows **Clean Architecture principles**, separating:

* Transport Layer
* Business Logic
* Data Access Layer

graph TD
    subgraph Client_Side
        Client["Web / Mobile Client"]
        Postman["Postman / Tester"]
    end

    subgraph API_Gateway
        Gin["Gin Web Server"]
        Middleware["Middleware: CORS, Logger, Recovery"]
    end

    subgraph Business_Logic
        Controller["Controllers"]
        Service["Services"]
        Repo["Repositories"]
    end

    subgraph Infrastructure
        PG[(PostgreSQL)]
        Redis[(Redis)]
        ES[(Elasticsearch)]
    end

    Client -->|HTTP/REST| Gin
    Postman -->|HTTP/REST| Gin
    Gin --> Middleware --> Controller
    Controller --> Service
    Service --> Repo

    Repo -->|GORM| PG
    Repo -->|Cache| Redis
    Repo -->|Search| ES

    Service -.->|Async Sync| ES

# ⚡ Key Features

## 📦 Advanced Catalog Management

* Hierarchical **Categories**
* Dynamic Global Attributes
* Product-specific attribute binding via `PRODUCT_ATTRIBUTE`
* Multi-language support using `JSONB` (EN / VI)

---

## 🧱 Deep 3-Tier Product Architecture

```
Product → Variant → Option
```

### Product

* Base information
* Brand
* Category
* Slug
* Name (JSONB i18n)

### Variant

* Group-level configuration
* Example: `"iPhone 15 Pro VN/A"`

### Option

* Actual sellable unit
* SKU
* Price
* Weight

Example:

```
Color: Black
Storage: 128GB
```

---

## 📊 Smart Inventory Management

* Option-level tracking
* Multi-warehouse support
* Stock movement logs
* Reservation logic
* Full audit trail

---

## 🔎 High-Performance Search

* Elasticsearch full-text search
* Typo-tolerant matching
* Real-time sync with PostgreSQL
* Optimized indexing

---

## 🚀 Performance Optimization

* Redis caching
* Optimized DB indexes
* Cascade deletion support

---

# 📂 Database Schema (ERD)

```mermaid
erDiagram
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
```

---

# 🛠 Tech Stack

| Component | Technology    |
| --------- | ------------- |
| Language  | Go 1.22+      |
| Framework | Gin Gonic     |
| Database  | PostgreSQL    |
| ORM       | GORM          |
| Cache     | Redis         |
| Search    | Elasticsearch |
| Config    | Viper         |
| UUID      | Google UUID   |

---

# 🚀 Getting Started

## ✅ Prerequisites

* Go 1.22+
* Docker & Docker Compose (Recommended)
* PostgreSQL
* Redis
* Elasticsearch

---

## 📦 Clone Repository

```bash
git clone https://github.com/your-username/product-management-api.git
cd product-management-api
```

---

## ⚙ Environment Configuration

Create `.env` file:

```env
# Database
DB_SOURCE=postgres://user:password@localhost:5432/product_db?sslmode=disable

# Server
SERVER_ADDRESS=:8080
BASE_URL=http://localhost:8080

# Redis
REDIS_ADDRESS=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Elasticsearch
ELASTICSEARCH_ADDRESS=http://localhost:9200
```

---

## 📥 Install Dependencies

```bash
go mod tidy
```

---

## ▶ Run Application

```bash
go run main.go
```

---

## 🌱 Seed Initial Data

```bash
go run cmd/seed/main.go
```

Seeds:

* Categories
* Brands
* Warehouses
* 3-tier product example (iPhone 15 Pro)

---

# 🔌 API Documentation (v1)

---

## 📦 Products

```
GET    /api/v1/products
GET    /api/v1/products/search
GET    /api/v1/products/:id
POST   /api/v1/products
PUT    /api/v1/products/:id
DELETE /api/v1/products/:id
POST   /api/v1/products/:id/attributes
```

---

## 🧱 Variants & Options

```
POST   /api/v1/variants
POST   /api/v1/variants/:id/options
PUT    /api/v1/options/:optionId
DELETE /api/v1/variants/:id
```

---

## 📊 Inventory

```
GET  /api/v1/inventory/option/:id
POST /api/v1/inventory/adjust
```

---

## 🏷 Master Data

```
/api/v1/categories
/api/v1/brands
/api/v1/attributes
/api/v1/warehouses
```

---

## 🖼 Media

```
POST /api/v1/media/upload
```

