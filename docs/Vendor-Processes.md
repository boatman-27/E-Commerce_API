## Vendor Processes

### overview

This document explains the Go code implementing the **Vendor Product Management** functionality in an eCommerce backend system. The implementation is divided into two main packages:

- `services/VendorService.go`: Contains business logic.
- `controllers/VendorController.go`: Exposes HTTP routes using Gin.

---

## `VendorService`

### Fields:

- `DB`: A pointer to a `sqlx.DB` instance for database operations.

### Methods:

#### AddProduct

- **Purpose**: Adds a new product for a vendor.
- **Inputs**: `*models.Product`, `vendorId string`
- **Returns**: `*models.Product`, `error`
- **Key Operations**:
  - Inserts the product into the `products` table.
  - Associates the product with the vendor.
  - Returns the inserted product.

#### GetProductByID

- **Purpose**: Fetches a single product by its ID (scoped to vendor).
- **Inputs**: `productId string`, `vendorId string`
- **Returns**: `*models.Product`, `error`
- **Key Operations**:
  - Validates that the product belongs to the vendor.
  - Fetches and returns the product.

#### GetVendorProducts

- **Purpose**: Fetches all products of a vendor.
- **Inputs**: `vendorId string`
- **Returns**: `[]*models.Product`, `error`
- **Key Operations**:
  - Queries the database for all vendor-owned products.

#### DeleteProduct

- **Purpose**: Deletes a vendor’s product.
- **Inputs**: `productId string`, `vendorId string`
- **Returns**: `error`
- **Key Operations**:
  - Checks if the product exists.
  - Verifies vendor ownership.
  - Deletes the product if valid.

#### UpdateProduct

- **Purpose**: Updates product details.
- **Inputs**: `*models.Product`, `vendorId string`
- **Returns**: `*models.Product`, `error`
- **Key Operations**:
  - Verifies vendor ownership.
  - Dynamically builds an SQL `SET` clause for provided fields.
  - Updates product details and timestamps.
  - Returns the updated product.

---

### `VendorController`

### Fields:

- `vendorService`: A pointer to the `VendorService` struct containing the business logic.

### Methods:

#### `AddProduct`

- **Method**: `POST`
- **Path**: `/vendor/products`
- **Behavior**:
  - Binds incoming JSON to a `Product` struct.
  - Extracts `UserId` (vendorId) from context.
  - Calls `VendorService.AddProduct`.
  - Returns the new product in JSON.

#### `GetProductById`

- **Method**: `GET`
- **Path**: `/vendor/products/id?productId={id}`
- **Behavior**:
  - Accepts `productId` as a query parameter.
  - Validates vendor ownership.
  - Returns the product JSON.

#### `GetVendorProducts`

- **Method**: `GET`
- **Path**: `/vendor/products`
- **Behavior**:
  - Fetches all products belonging to the vendor.
  - Returns a JSON array of products.

#### `DeleteProduct`

- **Method**: `POST`
- **Path**: `/vendor/products/id?productId={id}`
- **Behavior**:
  - Accepts `productId` as a query parameter.
  - Verifies vendor ownership.
  - Deletes the product.
  - Returns a success message.

#### `UpdateProduct`

- **Method**: `PATCH`
- **Path**: `/vendor/products`
- **Behavior**:
  - Binds JSON to a `Product` struct.
  - Validates `ProductId`.
  - Verifies vendor ownership.
  - Updates product fields.
  - Returns the updated product.

---

## Data Models in Golang

```go
// === === === === ===
//
//	=== Products ===
//
// === === === === ===

type Product struct {
	ProductId   string    `json:"ProductId" db:"productid"`
	VendorId    string    `json:"vendorId" db:"vendorid"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	SKU         string    `json:"sku" db:"sku"`
	Price       *float64  `json:"price" db:"price"`
	Discount    *float64  `json:"discount" db:"discount"`
	Stock       *int      `json:"stock" db:"stock"`
	Brand       string    `json:"brand" db:"brand"`
	Category    string    `json:"category" db:"category"`
	IsActive    *bool     `json:"isActive" db:"is_active"`
	CreatedAt   time.Time `json:"createdAt" db:"createdat"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updatedat"`
}
```

## sQL Tables

```sQL
CREATE TABLE products (
	  productid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	  vendorid VARCHAR(36) NOT NULL REFERENCES users(userid) ON DELETE CASCADE
	  name TEXT NOT NULL,
	  description TEXT,
	  sku TEXT UNIQUE NOT NULL,
	  price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
	  discount DECIMAL(5, 2) DEFAULT 0.00 CHECK (discount >= 0 AND discount <= 100),
	  stock INT NOT NULL DEFAULT 0 CHECK (stock >= 0),
	  brand TEXT NOT NULL,
	  category TEXT NOT NULL,
	  is_active BOOLEAN DEFAULT TRUE,
	  createdat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	  updatedat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## Dummy Data

```sQL
INSERT INTO products (
  vendorid,
  name,
  description,
  sku,
  price,
  discount,
  stock,
  brand,
  category,
  is_active
)
VALUES
-- 1
('b01ecc65-7139-4643-be15-c6fb5c14e297', 'Wireless Mouse', 'Ergonomic wireless mouse with 3 DPI settings', 'WM-001', 19.99, 5.00, 100, 'Logitech', 'Electronics', TRUE),
-- 2
('b01ecc65-7139-4643-be15-c6fb5c14e297', 'Mechanical Keyboard', 'Backlit mechanical keyboard with blue switches', 'MK-101', 49.99, 10.00, 50, 'Keychron', 'Electronics', TRUE),
-- 3
('b01ecc65-7139-4643-be15-c6fb5c14e297', 'Cotton T-Shirt', '100% organic cotton, unisex fit', 'TS-009', 14.99, 0.00, 200, 'H&M', 'Apparel', TRUE),
-- 4
('b01ecc65-7139-4643-be15-c6fb5c14e297', 'Ceramic Mug', '12oz ceramic mug with heat-resistant coating', 'MG-555', 9.99, 15.00, 120, 'Mugify', 'Home & Kitchen', TRUE),
-- 5
('b01ecc65-7139-4643-be15-c6fb5c14e297', 'Notebook', 'A5 size dotted notebook, 200 pages', 'NB-204', 6.49, 0.00, 300, 'Moleskine', 'Stationery', TRUE);
```

---

## Example JSON

### Add Product Request

```json
{
  "name": "Gaming Laptop",
  "description": "High-end laptop for gaming and productivity",
  "sku": "LAP12345",
  "price": 1500.0,
  "discount": 100.0,
  "stock": 10,
  "brand": "AlienTech",
  "category": "Electronics",
  "isActive": true
}
```

### Add Product Response

```json
{
  "product": {
    "productId": "123e4567-e89b-12d3-a456-426614174000",
    "vendorId": "dc22872c-f003-40df-b61a-743c97945b33",
    "name": "Gaming Laptop",
    "description": "High-end laptop for gaming and productivity",
    "sku": "LAP12345",
    "price": 1500.0,
    "discount": 100.0,
    "stock": 10,
    "brand": "AlienTech",
    "category": "Electronics",
    "isActive": true
  }
}
```

### Get Product Response

```json
{
  "product": {
    "ProductId": "123e4567-e89b-12d3-a456-426614174000",
    "vendorId": "dc22872c-f003-40df-b61a-743c97945b33",
    "name": "Gaming Laptop",
    "description": "High-end laptop for gaming and productivity",
    "sku": "LAP12345",
    "price": 1500.0,
    "discount": 100.0,
    "stock": 10,
    "brand": "AlienTech",
    "category": "Electronics",
    "isActive": true
  }
}
```

### Get Vendor Products Response

```json
{
  "products": [
    {
      "ProductId": "123e4567-e89b-12d3-a456-426614174000",
      "name": "Gaming Laptop",
      "price": 1500.0,
      "stock": 10,
      "brand": "AlienTech",
      "category": "Electronics"
    }
  ]
}
```

### Delete Product Response

```json
{
  "message": "Product deleted successfully"
}
```

### Update Product Request

```json
{
  "ProductId": "123e4567-e89b-12d3-a456-426614174000",
  "price": 1400.0,
  "stock": 8,
  "isActive": false
}
```

### Update Product Response

```json
{
  "Updated Product": {
    "ProductId": "123e4567-e89b-12d3-a456-426614174000",
    "name": "Gaming Laptop",
    "price": 1400.0,
    "stock": 8,
    "isActive": false
  }
}
```
