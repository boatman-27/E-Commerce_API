# Product Browsing and Search

### Overview

This code defines a product browsing and searching using the Go programming language, Gin web framework, and PostgreSQL (via sqlx). It is split into two main files:

- `controllers/ProductController.go`: Handles HTTP requests/responses for user operations.
- `services/ProductService.go`: Contains business logic and database interactions.

---

## `Product Service`

### `GetProducts`

- **Purpose**: Retrieves a list of active products from the database, filtered and sorted based on query parameters.
- **Inputs**:
  - `filters map[string]string`: Contains optional filters such as:
    - `category`: Filters products by category (case-insensitive, partial match)
    - `brand`: Filters products by brand (case-insensitive, partial match)
    - `min_price`: Filters products with price >= min_price
    - `max_price`: Filters products with price <= max_price
    - `discount`: Filters products with discount >= value
    - `keyword`: Searches for the keyword in name or description (case-insensitive, partial match)
    - `sort_by`: Field to sort by (`price`, `discount`, `stock`, `createdat`)
    - `order`: Sorting order (`ASC` or `DESC`)
    - `limit`: Number of products to return
    - `offset`: Number of products to skip
- **Returns**: A slice of `*models.Product` and an error if any occurred.
- **Key Operations**:
  - Builds a dynamic SQL query with filters
  - Appends conditions and arguments securely
  - Adds sorting, ordering, and pagination
  - Executes the query and maps results to product model

---

## `Product Controller`

### `GetProducts`

- **Method**: `GET`
- **Path**: `/products`
- **Behavior**:
  - Parses query parameters from the HTTP request
  - Constructs a filter map and passes it to `ProductService.GetProducts`
  - On success: returns `200 OK` with JSON containing list of products
  - On error: returns `500 Internal Server Error` with error message

---

## sQL Table

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

CREATE INDEX idx_products_productid_stock_cover ON products(productid) INCLUDE (stock);
```

## Dummy data

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

## Example Requests

### Request

```HTTP
GET /products?category=electronics&min_price=100&max_price=500&sort_by=price&order=DESC&limit=10&offset=20
```

### Response

```json
{
  "products": [
    {
      "productId": "c41c78d2-31fd-4f6e-b35e-7719c111aa00",
      "name": "Wireless Headphones",
      "description": "Noise-cancelling over-ear headphones",
      "category": "Electronics",
      "brand": "SoundMax",
      "price": 299.99,
      "discount": 10.0,
      "stock": 35,
      "createdAt": "2025-08-01T12:30:00Z",
      "isActive": true
    }
  ]
}
```
