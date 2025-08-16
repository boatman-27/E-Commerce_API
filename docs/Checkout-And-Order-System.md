# Checkout & Order System

### Overview

This code handles checkout operations, including order summary retrieval and purchase confirmation, using Go, Gin, and PostgreSQL via `sqlx`. It is split into two main files:

- `services/CheckoutService.go`: Contains business logic for generating order summaries, confirming purchases, and handling stock transactions.
- `controllers/CheckoutController.go`: Handles HTTP requests/responses for checkout-related operations.
- `services/OrderService.go`
- `controllers/OrderController.go`

---

## Checkout Operations

## `Checkout Service`

### `OrderSummary`

- **Purpose**: Retrieves the current cart summary and the user’s default shipping address.
- **Inputs**:
  - `cartId string`: ID of the cart to summarize.
  - `userId string`: ID of the user making the request.
- **Returns**: A pointer to `models.Summary` containing:
  - `ReturnedCart` – cart contents and totals
  - `ShippingAddress` – default shipping address for the user
- **Key Operations**:
  - Calls `CartService.ViewCart` to fetch the cart.
  - Calls `ShippingService.GetDefaultShippingAddress` for shipping info.
  - Combines results into a `Summary` object.

### `ConfirmPurchase`

- **Purpose**: Finalizes a purchase by creating an order, reducing stock, inserting order items, and sending confirmation emails.
- **Inputs**:
  - `cartId string`: ID of the cart to checkout.
  - `userId string`: ID of the purchasing user.
  - `userEmail string`: User’s email for order confirmation.
- **Returns**: A pointer to `models.Summary` for the completed order.
- **Key Operations**:
  - Begins a SQL transaction.
  - Retrieves the order summary using `OrderSummary`.
  - Validates that the cart is not empty.
  - Inserts a new order into `orders` and retrieves the new `orderId`.
  - Locks product rows for update using `SELECT ... FOR UPDATE` to prevent race conditions.
  - Checks stock availability for each product.
  - Updates stock for each product in the database.
  - Inserts each item into `order_items`.
  - Deletes the cart using `CartService.DeleteCart`.
  - Sends a confirmation email via `helpers.SendConfirmationEmail`.
  - Commits the transaction.

---

## `Checkout Controller`

### `OrderSummary`

- **Method**: `GET`
- **Path**: `/summary?cartId=<id>`
- **Behavior**:
  - Retrieves `cartId` from query parameters.
  - Retrieves `UserId` from context (set by authentication middleware).
  - Calls `CheckoutService.OrderSummary`.
  - On success: returns `200 OK` with JSON containing the summary.
  - On error: returns `400` or `500` with error message.

---

### `ConfirmPurchase`

- **Method**: `POST`
- **Path**: `/confirm?cartId=<id>`
- **Behavior**:
  - Retrieves `cartId` from query parameters.
  - Retrieves `UserId` and `Email` from context.
  - Calls `CheckoutService.ConfirmPurchase`.
  - On success: returns `200 OK` with JSON containing the order summary.
  - On error: returns `400` or `500` with error message.

---

## Order Operations

## `OrderService`

### Fields:

- `DB`: A pointer to a `sqlx.DB` instance for database operations.

### Methods:

#### checkOrderOwnership

- **Purpose**: Ensures the given order belongs to the user.
- **Inputs**: `orderId string`, `userId string`
- **Returns**: `error`
- **Key Operations**:
  - Fetches `userid` from the orders table.
  - Returns `ErrNotFound` if order doesn’t exist.
  - Returns `ErrUnauthorized` if user does not own the order.

#### TrackOrder

- **Purpose**: Returns the current status of a user’s order.
- **Inputs**: `orderId string`, `userId string`
- **Returns**: `status string`, `error`
- **Key Operations**:
  - Queries `orders` table for order matching `orderId` and `userId`.
  - Returns `ErrNotFound` if no match.

#### ViewOrders

- **Purpose**: Returns all orders of a user with their items.
- **Inputs**: `userId string`
- **Returns**: `[]*models.Order`, `error`
- **Key Operations**:
  - Joins `orders`, `order_items`, and `products` tables.
  - Groups items under each order.
  - Returns sorted list by `orderedAt`.

---

### `OrderController`

### Fields:

- `orderService`: A pointer to `OrderService`.

### Methods:

#### TrackOrder

- **Method**: `GET`
- **Path**: `/orders/track?orderId={id}`
- **Behavior**:
  - Accepts `orderId` as query parameter.
  - Retrieves `UserId` from context.
  - Calls `OrderService.TrackOrder`.
  - Returns order status or `404` if not found.

#### ViewOrders

- **Method**: `GET`
- **Path**: `/orders`
- **Behavior**:
  - Retrieves `UserId` from context.
  - Calls `OrderService.ViewOrders`.
  - Returns JSON array of orders with items.

---

## Data Models in Golang

```go
// === === === === ===
//
//	=== Orders ===
//
// === === === === ===

type OrderItem struct {
	ProductId   string `json:"productId" db:"productid"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Brand       string `json:"brand" db:"brand"`
	Category    string `json:"category" db:"category"`
	Quantity    int    `json:"quantity" db:"quantity"`
}

type Order struct {
	OrderId    string       `json:"orderId" db:"orderid"`
	TotalPrice float64      `json:"totalPrice" db:"total_price"`
	OrderedAt  time.Time    `json:"orderedAt" db:"orderedat"`
	Items      []*OrderItem `json:"items"`
}

```

---

## sQL Tables

```sQL
CREATE TABLE orders (
	orderid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	cartid UUID NOT NULL REFERENCES carts(cartid) ON DELETE CASCADE,
    userid VARCHAR(36) NOT NULL UNIQUE REFERENCES users(userid) ON DELETE CASCADE,
    status VARCHAR(10) NOT NULL DEFAULT 'processing' CHECK (status  IN ('processing', 'shipping', 'delivery', 'delivered'))
    total_price FLOAT NOT NULL
    orderedat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
);

CREATE INDEX idx_orders_userid ON orders(userid);
```

```sQL
CREATE TABLE order_items (
	orderitemid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	orderid UUID NOT NULL REFERENCES orders(orderid) ON DELETE CASCADE,
	productid UUID NOT NULL REFERENCES products(productid) ON DELETE CASCADE,
    quantity INT NOT NULL CHECK (quantity >= 1),
    addedat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(orderid, productid)
);

CREATE INDEX idx_order_items_orderid ON order_items(orderid);
CREATE INDEX idx_order_items_productid ON order_items(productid);
```

---

## Example Requests

### Order Summary

```HTTP
GET /checkout/summary?cartId=123e4567-e89b-12d3-a456-426614174000
```

### Response

```json
{
  "summary": {
    "returnedCart": {
      "items": [
        {
          "productId": "c41c78d2-31fd-4f6e-b35e-7719c111aa00",
          "name": "Wireless Mouse",
          "quantity": 2,
          "price": 19.99
        }
      ],
      "total": 39.98
    },
    "shippingAddress": {
      "addressLine1": "123 Main St",
      "city": "Muscat",
      "postalCode": "1001",
      "country": "Oman"
    }
  }
}
```

### Confirm Purchase

```HTTP
GET /checkout/confirm?cartId=123e4567-e89b-12d3-a456-426614174000
```

### Response

```json
{
  "summary": {
    "returnedCart": { ... },
    "shippingAddress": { ... }
  }
}
```

### Track Order Response

```json
{
  "status": "shipping"
}
```

### View Orders Response

```json
{
  "orders": [
    {
      "orderId": "order-uuid-1",
      "totalPrice": 129.99,
      "orderedAt": "2025-08-16T11:00:00Z",
      "items": [
        {
          "productId": "product-uuid-1",
          "name": "Wireless Mouse",
          "description": "Ergonomic wireless mouse with 3 DPI settings",
          "brand": "Logitech",
          "category": "Electronics",
          "quantity": 2
        },
        {
          "productId": "product-uuid-2",
          "name": "Mechanical Keyboard",
          "description": "Backlit mechanical keyboard with blue switches",
          "brand": "Keychron",
          "category": "Electronics",
          "quantity": 1
        }
      ]
    }
  ]
}
```
