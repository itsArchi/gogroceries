# GoGroceries 🛒

A RESTful API backend for an online groceries marketplace built with Go (Golang), Gin framework, and PostgreSQL(supabase). This application provides a complete e-commerce solution with user authentication, store management, product catalog, and transaction processing.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Technology Stack](#technology-stack)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
- [Project Structure](#project-structure)
- [Database Schema](#database-schema)

## Features

### User Management

- User registration and authentication
- JWT-based authorization
- User profile management
- Address management for delivery

### Store (Toko) Management

- Store creation and management
- Store profile updates
- Browse all stores
- View store details

### Product Management

- Create, read, update, and delete products
- Product categorization
- Product images support
- Product search and filtering
- Product inventory tracking with logs

### Category Management

- Product category management
- Admin-only category operations (create, update, delete)
- Public category browsing

### Transaction Management

- Shopping cart and checkout
- Transaction history
- Transaction details with line items
- Order tracking

### Security

- JWT authentication
- Password hashing
- Role-based access control (Admin/User)
- Protected routes with middleware

## Architecture

This project follows **Clean Architecture** principles with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────┐
│                    Delivery Layer                        │
│              (HTTP Handlers & Middleware)                │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                    Usecase Layer                         │
│                  (Business Logic)                        │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                   Repository Layer                       │
│              (Data Access & Persistence)                 │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                    Domain Layer                          │
│            (Entities & Business Rules)                   │
└─────────────────────────────────────────────────────────┘
```

## Technology Stack

- **Language**: Go 1.25.3
- **Web Framework**: Gin v1.11.0
- **ORM**: GORM v1.31.0
- **Database**: PostgreSQL (via pgx driver v5.6.0)
- **Authentication**: JWT (golang-jwt/jwt v5.3.0)
- **Password Hashing**: bcrypt
- **Environment Variables**: godotenv v1.5.1
- **Validation**: go-playground/validator v10.27.0
- **UUID Generation**: google/uuid v1.6.0

## Prerequisites

Before running this application, make sure you have:

- **Go** 1.25.3 or higher installed
- **PostgreSQL** database server running
- **Git** for cloning the repository

## Installation

1. **Clone the repository**

```bash
git clone https://github.com/itsArchi/gogroceries
cd gogroceries
```

2. **Install dependencies**

```bash
go mod download
```

3. **Set up the database**

Create a PostgreSQL database:

```sql
CREATE DATABASE gogroceries;
```

## Configuration

1. **Create environment file**

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

2. **Configure environment variables**

Edit the `.env` file with your settings:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=gogroceries
JWT_SECRET=your_jwt_secret_key_here
SERVER_PORT=8080
```

### Environment Variables Description

| Variable      | Description                         | Default       |
| ------------- | ----------------------------------- | ------------- |
| `DB_HOST`     | PostgreSQL host address             | `localhost`   |
| `DB_PORT`     | PostgreSQL port                     | `5432`        |
| `DB_USER`     | Database username                   | `postgres`    |
| `DB_PASSWORD` | Database password                   | -             |
| `DB_NAME`     | Database name                       | `gogroceries` |
| `JWT_SECRET`  | Secret key for JWT token generation | `secret`      |
| `SERVER_PORT` | Port for the API server             | `8080`        |

## Running the Application

1. **Run the application**

```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:8080` (or your configured port).

2. **Database migrations**

The application automatically runs database migrations on startup, creating all necessary tables:

- users
- tokos (stores)
- alamats (addresses)
- categories
- produks (products)
- foto_produks (product photos)
- trxs (transactions)
- detail_trxs (transaction details)
- log_produks (product logs)

## API Documentation

### Base URL

```
http://localhost:8080/api/v1
```

### Authentication Endpoints

#### Register

```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "nama": "Uzumaki Udin",
  "kata_sandi": "udinGantenk123",
  "no_telp": "081234567890",
  "email": "uzumaki@udin.com",
  "tanggal_Lahir": "1990-01-01",
  "pekerjaan": "Software Engineer",
  "id_provinsi": "1",
  "id_kota": "1",
  "jenis_kelamin": "L",
  "tentang": "About me"
}
```

#### Login

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "no_telp": "081234567890",
  "kata_sandi": "udinGantenk123"
}
```

### User Endpoints (Protected)

All user endpoints require JWT authentication via `Authorization: Bearer <token>` header.

#### Get My Profile

```http
GET /api/v1/user
Authorization: Bearer <token>
```

#### Update Profile

```http
PUT /api/v1/user
Authorization: Bearer <token>
Content-Type: application/json

{
  "nama": "Uzumaki Udin Gantenk",
  "email": "uzumakigantenk@udin.com"
}
```

### Address Endpoints (Protected)

#### Get All User Addresses

```http
GET /api/v1/user/alamat
Authorization: Bearer <token>
```

#### Create Address

```http
POST /api/v1/user/alamat
Authorization: Bearer <token>
Content-Type: application/json

{
  "judul_alamat": "Home",
  "nama_penerima": "Uzumaki Udin Gantenk",
  "no_telp": "081234567890",
  "detail_alamat": "Jl. jalan doang jadian kaga No. 123"
}
```

#### Get Address by ID

```http
GET /api/v1/user/alamat/:id
Authorization: Bearer <token>
```

#### Update Address

```http
PUT /api/v1/user/alamat/:id
Authorization: Bearer <token>
```

#### Delete Address

```http
DELETE /api/v1/user/alamat/:id
Authorization: Bearer <token>
```

### Store (Toko) Endpoints

#### Get All Stores

```http
GET /api/v1/toko
```

#### Get Store by ID

```http
GET /api/v1/toko/:id_toko
```

#### Get My Store (Protected)

```http
GET /api/v1/toko/my
Authorization: Bearer <token>
```

#### Update Store (Protected)

```http
PUT /api/v1/toko/:id_toko
Authorization: Bearer <token>
```

### Product Endpoints

#### Get All Products

```http
GET /api/v1/product
```

#### Get Product by ID

```http
GET /api/v1/product/:id
```

#### Create Product (Protected)

```http
POST /api/v1/product
Authorization: Bearer <token>
Content-Type: application/json

{
  "nama_produk": "Fresh Tomatoes",
  "slug": "fresh-tomatoes",
  "harga_reseller": "5000",
  "harga_konsumen": "7000",
  "stok": 100,
  "deskripsi": "Fresh organic tomatoes",
  "id_category": 1
}
```

#### Update Product (Protected)

```http
PUT /api/v1/product/:id
Authorization: Bearer <token>
```

#### Delete Product (Protected)

```http
DELETE /api/v1/product/:id
Authorization: Bearer <token>
```

### Category Endpoints

#### Get All Categories

```http
GET /api/v1/category
```

#### Get Category by ID

```http
GET /api/v1/category/:id
```

#### Create Category (Admin Only)

```http
POST /api/v1/category
Authorization: Bearer <token>
Content-Type: application/json

{
  "nama_category": "Vegetables"
}
```

#### Update Category (Admin Only)

```http
PUT /api/v1/category/:id
Authorization: Bearer <token>
```

#### Delete Category (Admin Only)

```http
DELETE /api/v1/category/:id
Authorization: Bearer <token>
```

### Transaction Endpoints (Protected)

#### Create Transaction

```http
POST /api/v1/trx
Authorization: Bearer <token>
Content-Type: application/json

{
  "id_alamat": 1,
  "detail_trx": [
    {
      "id_produk": 1,
      "kuantitas": 2
    }
  ]
}
```

#### Get All User Transactions

```http
GET /api/v1/trx
Authorization: Bearer <token>
```

#### Get Transaction by ID

```http
GET /api/v1/trx/:id
Authorization: Bearer <token>
```

##  Project Structure

```
gogroceries/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── config/
│   └── config.go                # Configuration management
├── delivery/
│   ├── http/
│   │   ├── router.go            # Route definitions
│   │   ├── handler_auth.go      # Authentication handlers
│   │   ├── user_handler.go      # User handlers
│   │   ├── toko_handler.go      # Store handlers
│   │   ├── produk_handler.go    # Product handlers
│   │   ├── category_handler.go  # Category handlers
│   │   ├── trx_handler.go       # Transaction handlers
│   │   └── alamat_handler.go    # Address handlers
│   └── middleware/
│       └── auth.go              # Authentication middleware
├── domain/
│   ├── user.go                  # User domain models
│   ├── auth.go                  # Auth domain models
│   ├── toko.go                  # Store domain models
│   ├── product.go               # Product domain models
│   ├── category.go              # Category domain models
│   ├── trx.go                   # Transaction domain models
│   ├── alamat.go                # Address domain models
│   ├── detail_transaction.go    # Transaction detail models
│   └── response.go              # Response models
├── internal/
│   └── helper/
│       ├── jwt.go               # JWT utilities
│       ├── hash.go              # Password hashing
│       └── helper.go            # General helpers
├── repository/
│   └── postgres/
│       ├── db.go                # Database connection
│       ├── user_repository.go   # User repository
│       ├── toko_repository.go   # Store repository
│       ├── produk_repository.go # Product repository
│       ├── category_repository.go # Category repository
│       ├── trx_repository.go    # Transaction repository
│       └── alamat_repository.go # Address repository
├── usecase/
│   ├── auth_usecase.go          # Authentication business logic
│   ├── user_usecase.go          # User business logic
│   ├── toko_usecase.go          # Store business logic
│   ├── produk_usecase.go        # Product business logic
│   ├── category_usecase.go      # Category business logic
│   ├── trx_usecase.go           # Transaction business logic
│   └── alamat_usecase.go        # Address business logic
├── .env                         # Environment variables (not in git)
├── .env.example                 # Environment variables template
├── go.mod                       # Go module dependencies
├── go.sum                       # Go module checksums
└── README.md                    # This file
```

##  Database Schema

### Main Entities

1. **Users** - User accounts and profiles
2. **Tokos** - Store/shop information
3. **Alamats** - User delivery addresses
4. **Categories** - Product categories
5. **Produks** - Product catalog
6. **FotoProduk** - Product images
7. **Trxs** - Transactions/orders
8. **DetailTrxs** - Transaction line items
9. **LogProduks** - Product inventory logs

### Key Relationships

- One User can have one Toko (store)
- One User can have multiple Alamats (addresses)
- One Toko can have multiple Produks (products)
- One Category can have multiple Produks
- One Produk can have multiple FotoProduk (images)
- One Trx belongs to one User and one Alamat
- One Trx can have multiple DetailTrxs
- Products have LogProduks for inventory tracking

## Security Best Practices

- Passwords are hashed using bcrypt before storage
- JWT tokens are used for stateless authentication
- Sensitive routes are protected with authentication middleware
- Admin-only routes have additional role-based authorization
- Environment variables are used for sensitive configuration
- SQL injection protection via GORM ORM

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Authors

itsArchi

## Support

For support, please open an issue in the GitHub repository or contact the development team.

---

