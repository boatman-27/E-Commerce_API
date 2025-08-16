# Cart Management System

### Overview

This code defines a shopping cart system using Go, Gin, and PostgreSQL (via sqlx). It supports creating carts, adding/editing/deleting cart items, viewing cart contents, and calculating total prices. It is split into two main files.

- `services/CartService.go`: Contains the business logic and database interactions.
- `controllers/CartController.go`: Handles HTTP requests/responses for cart operations.

---

## `CartService`

### Fields:

- `DB`: A pointer to a `sqlx.DB` instance for database operations.

### Methods:

#### `checkOwnership`

- **Purpose**: Checks if the cart belongs to a given user.
- **Inputs**: `cartId string`, `userId string`
- **Returns**: `error`
- **Key Operations**:
  - Fetches the cart's owner from the database.
  - Returns an unauthorized error if the user does not own the cart.

#### `getAvailableStock`

- **Purpose**: Returns the available stock for a product.
- **Inputs**: `productId string`
- **Returns**: `int`, `error`

#### `getExistingQty`

- **Purpose**: Returns the current quantity of a product in a cart.
- **Inputs**: `cartId string`, `productId string`
- **Returns**: `int`, `error`

#### `getTotal`

- **Purpose**: Calculates the total cost of a cart.
- **Inputs**: `db sqlx.Ext`, `cartId string`, `userId string`
- **Returns**: `float64`, `error`
- **Key Operations**:
  - Joins `cart_items` with `products` and calculates `SUM((price - discount) * quantity)`.

#### `DeleteCart`

- **Purpose**: Deletes all items in a cart.
- **Inputs**: `cartId string`, `userId string`
- **Returns**: `error`

#### `ViewCart`

- **Purpose**: Fetches all items in a cart along with the total price.
- **Inputs**: `cartId string`, `userId string`
- **Returns**: `*models.ReturnedCart`, `error`

#### `GetTotal`

- **Purpose**: Returns the total price of a cart.
- **Inputs**: `cartId string`, `userId string`
- **Returns**: `float64`, `error`

#### `AddToCart`

- **Purpose**: Adds a product to a cart or increments its quantity.
- **Inputs**: `cartId string`, `userId string`, `cartItem *models.CartItem`
- **Returns**: `*models.ReturnedCart`, `error`
- **Key Operations**:
  - Checks ownership and available stock.
  - Inserts or updates quantity in `cart_items`.
  - Returns updated cart with recalculated total.

#### `EditCartItem`

- **Purpose**: Updates quantity of a cart item.
- **Inputs**: `cartItem *models.CartItem`, `userId string`
- **Returns**: `*models.CartItem`, `error`

#### `DeleteCartItem`

- **Purpose**: Deletes a specific item from a cart.
- **Inputs**: `cartItemId string`, `userId string`, `cartId string`
- **Returns**: `error`

#### `CreateCart`

- **Purpose**: Creates a new cart for a user.
- **Inputs**: `userId string`
- **Returns**: `error`
- **Key Operations**:
  - Checks if user already has a cart.
  - Inserts new cart record.

---

## `CartController`

### Fields:

- `CartService`: A pointer to the `CartService` struct.

### Methods:

#### `ViewCart`

- **Method**: `GET`
- **Path**: `/cart?cartId=<id>`
- **Behavior**:
  - Fetches `cartId` from query params.
  - Retrieves `UserId` from context.
  - Calls `CartService.ViewCart`.
  - Returns cart data in JSON.

#### `GetTotalPrice`

- **Method**: `GET`
- **Path**: `/cart/total?cartId=<id>`
- **Behavior**:
  - Fetches cart total price.

#### `DeleteCart`

- **Method**: `DELETE`
- **Path**: `/cart?cartId=<id>`
- **Behavior**:
  - Deletes all items in a cart.

#### `CreateCart`

- **Method**: `POST`
- **Path**: `/cart`
- **Behavior**:
  - Creates a new cart for the user.

#### `AddToCart`

- **Method**: `POST`
- **Path**: `/cart/item?cartId=<id>`
- **Behavior**:
  - Adds a new item to the cart or increases its quantity.

#### `EditCartItem`

- **Method**: `PATCH`
- **Path**: `/cart/item`
- **Behavior**:
  - Updates quantity of a specific item in the cart.

#### `DeleteCartItem`

- **Method**: `DELETE`
- **Path**: `/cart/item?cartId=<id>&cartItemId=<id>`
- **Behavior**:
  - Removes an item from the cart.

---

## Data Models in Golang

```sQL
CREATE TABLE carts (
    cartid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    userid VARCHAR(36) NOT NULL UNIQUE REFERENCES users(userid) ON DELETE CASCADE,
    createdat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updatedat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

```sQL
CREATE TABLE cart_items (
    cartitemid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cartid UUID NOT NULL REFERENCES carts(cartid) ON DELETE CASCADE,
    productid UUID NOT NULL REFERENCES products(productid) ON DELETE CASCADE,
    quantity INT NOT NULL CHECK (quantity >= 1),
    addedat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(cartid, productid)
);
```

---

## Dummy data

```sQL
INSERT INTO cart_items (cartid, productid, quantity)
VALUES
  ('fbb71fb5-020d-4b61-a25a-76cd22efd51d', '46b66902-7b7c-4513-b8d0-4f37a6b82d9c', 2),
  ('fbb71fb5-020d-4b61-a25a-76cd22efd51d', 'a3b0eaf5-d524-4506-8e8c-ba6694a7adb9', 1);
```

```sQL
SELECT products.productid, name, description, sku, price, discount, brand, category, quantity
FROM products
INNER JOIN cart_items ON products.productid = cart_items.productid
WHERE cart_items.cartid = 'dec9278c-690e-46b3-a734-9c16a3c93a4c'
```

---

## Example JSON

### View Cart Response

```json
{
  "cart": {
    "ID": "c1a2b3",
    "UserId": "dc22872c-f003-40df-b61a-743c97945b33",
    "Items": [
      {
        "ProductId": "p1",
        "Name": "Product A",
        "Description": "Description of product A",
        "SKU": "SKU001",
        "Price": 50.0,
        "Discount": 5.0,
        "Brand": "BrandX",
        "Category": "Electronics",
        "Quantity": 2
      }
    ],
    "Total": 90.0
  }
}
```

### Add to Cart Request

```json
{
  "ProductId": "p1",
  "Quantity": 2
}
```

### Add to Cart Response

```json
{
  "cart": {
    "ID": "c1a2b3",
    "UserId": "dc22872c-f003-40df-b61a-743c97945b33",
    "Items": [...],
    "Total": 90.0
  }
}
```

### Edit Cart Item Request

```json
{
  "CartItemId": "ci123",
  "CartId": "c1a2b3",
  "ProductId": "p1",
  "Quantity": 3
}
```
