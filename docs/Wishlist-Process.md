# Wishlist Processes

### overview

This document explains the Go code implementing **Wishlist Management** and **Reviews** functionality in an eCommerce backend system. The implementation is divided into two main packages:

- `services/WishlistService.go`: Contains business logic.
- `controllers/WishlistController.go`: Exposes HTTP routes using Gin.

---

## `WishlistService`

### Fields:

- `DB`: A pointer to a `sqlx.DB` instance for database operations.

### Methods:

#### checkOwnership

- **Purpose**: Ensures the wishlist belongs to the user.
- **Inputs**: `userId string`, `wishlistId string`
- **Returns**: `error`
- **Key Operations**:
  - Queries `wishlists` table for `userid`.
  - Returns `ErrNotFound` if wishlist doesn’t exist.
  - Returns `ErrUnauthorized` if the user does not own the wishlist.

#### CreateWishlist

- **Purpose**: Creates a new wishlist for the user.
- **Inputs**: `userId string`, `*models.WishlistInfo`
- **Returns**: `error`
- **Key Operations**:
  - Inserts a new row into `wishlists` with name and description.

#### GetWishlistItems

- **Purpose**: Returns all items in a user’s wishlist.
- **Inputs**: `userId string`, `wishlistId string`
- **Returns**: `*models.ReturnedWishlist`, `error`
- **Key Operations**:
  - Joins `products`, `wishlist_items`, and `wishlists`.
  - Returns list of items under wishlist.

#### EditWishlist

- **Purpose**: Updates wishlist details.
- **Inputs**: `userId string`, `wishlistId string`, `*models.WishlistInfo`
- **Returns**: `error`
- **Key Operations**:
  - Verifies ownership.
  - Dynamically builds `SET` clause for fields to update.
  - Updates `wishlists` table.

#### DeleteWishlist

- **Purpose**: Deletes a wishlist.
- **Inputs**: `userId string`, `wishlistId string`
- **Returns**: `error`
- **Key Operations**:
  - Deletes wishlist from `wishlists` table.

#### AddToWishlist

- **Purpose**: Adds a product to a wishlist.
- **Inputs**: `userId string`, `wishlistId string`, `productId string`
- **Returns**: `*models.ReturnedWishlist`, `error`
- **Key Operations**:
  - Verifies ownership.
  - Inserts product into `wishlist_items`.
  - Returns updated wishlist with items.

#### DeleteWishlistItem

- **Purpose**: Removes a specific product from a wishlist.
- **Inputs**: `userId string`, `wishlistId string`, `wishlistItemId string`
- **Returns**: `error`
- **Key Operations**:
  - Verifies ownership.
  - Deletes from `wishlist_items`.

#### MoveToCart

- **Purpose**: Moves all items in a wishlist to the user’s cart.
- **Inputs**: `userId string`, `wishlistId string`
- **Returns**: `error`
- **Key Operations**:
  - Inserts all wishlist items into `cart_items`.
  - Removes them from `wishlist_items`.
  - Uses transaction to ensure consistency.

#### MoveItemToCart

- **Purpose**: Moves a single wishlist item to the cart.
- **Inputs**: `userId string`, `wishlistId string`, `wishlistItemId string`
- **Returns**: `error`
- **Key Operations**:
  - Inserts item into `cart_items`.
  - Deletes from `wishlist_items`.
  - Transaction ensures atomicity.

---

### `WishlistController`

### Fields:

- `WishlistService`: A pointer to `WishlistService`.

### Methods:

#### CreateWishlist

- **Method**: `POST`
- **Path**: `/wishlists`
- **Behavior**:
  - Extracts `UserId` from context.
  - Binds JSON to `WishlistInfo`.
  - Calls `CreateWishlist`.
  - Returns success message.

#### DeleteWishlist

- **Method**: `DELETE`
- **Path**: `/wishlists?wishlistId={wishlistId}`
- **Behavior**:
  - Extracts `UserId` from context.
  - Deletes wishlist using service.
  - Returns success message.

#### EditWishlist

- **Method**: `PATCH`
- **Path**: `/wishlists?wishlistId={wishlistId}`
- **Behavior**:
  - Extracts `UserId` and JSON fields.
  - Calls `EditWishlist`.
  - Returns success message.

#### AddToWishlist

- **Method**: `POST`
- **Path**: `/wishlists/items?wishlistId={id}&productId={id}`
- **Behavior**:
  - Extracts `UserId` from context.
  - Calls `AddToWishlist`.
  - Returns updated wishlist.

#### DeleteWishlistItem

- **Method**: `DELETE`
- **Path**: `/wishlists/items?wishlistId={id}&wishlistItemId={id}`
- **Behavior**:
  - Calls `DeleteWishlistItem`.
  - Returns success message.

#### MoveToCart

- **Method**: `POST`
- **Path**: `/wishlists/move-to-cart?wishlistId={id}`
- **Behavior**:
  - Moves all wishlist items to cart.
  - Returns success message.

#### MoveItemToCart

- **Method**: `POST`
- **Path**: `/wishlists/item/move-to-cart?wishlistId={id}&wishlistItemId={id}`
- **Behavior**:
  - Moves a single wishlist item to cart.
  - Returns success message.

#### GetWishlistItems

- **Method**: `GET`
- **Path**: `/wishlists?wishlistId={id}`
- **Behavior**:
  - Returns list of products in a wishlist.

---

## Data Models in Golang

```go
type WishlistInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ItemData struct {
	ProductId   string   `json:"productId" db:"productid"`
	Name        string   `json:"name" db:"name"`
	Description string   `json:"description" db:"description"`
	SKU         string   `json:"sku" db:"sku"`
	Price       *float64 `json:"price" db:"price"`
	Discount    *float64 `json:"discount" db:"discount"`
	Brand       string   `json:"brand" db:"brand"`
	Category    string   `json:"category" db:"category"`
}

type ReturnedWishlist struct {
	ID     string      `json:"wishlistId"`
	UserId string      `json:"userId"`
	Items  []*ItemData `json:"items"`
}
```

---

## sQL Tables

```sQL
CREATE TABLE wishlists (
	wishlistid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	userid VARCHAR(36) NOT NULL REFERENCES users(userid) ON DELETE CASCADE,
	name TEXT NOT NULL,
    description TEXT,
    createdat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updatedat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_wishlists_userid ON wishlists(userid);

```

```sQL
CREATE TABLE wishlist_items (
    wishlistitemid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wishlistid UUID NOT NULL REFERENCES wishlists(wishlistid) ON DELETE CASCADE,
    productid UUID NOT NULL REFERENCES products(productid) ON DELETE CASCADE,
    addedat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updatedat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (wishlistid, productid)
);

CREATE INDEX idx_wishlist_items_wishlistid ON wishlist_items(wishlistid);
```

---

## Example JSON

### Create Wishlist Request

```json
{
  "name": "Holiday Wishlist",
  "description": "Items I want for holidays"
}
```

### Add To Wishlist Response

```json
{
  "wishlist": {
    "wishlistId": "wishlist-uuid-1",
    "userId": "u01abcde-1111-2222-3333-444455556666",
    "items": [
      {
        "productId": "product-uuid-1",
        "name": "Wireless Mouse",
        "description": "Ergonomic wireless mouse",
        "sku": "WM-001",
        "price": 19.99,
        "discount": 5.0,
        "brand": "Logitech",
        "category": "Electronics"
      }
    ]
  }
}
```

### Get Wishlist Items Response

```json
{
  "wishlist": {
    "wishlistId": "wishlist-uuid-1",
    "userId": "u01abcde-1111-2222-3333-444455556666",
    "items": [
      {
        "productId": "product-uuid-1",
        "name": "Wireless Mouse",
        "description": "Ergonomic wireless mouse",
        "sku": "WM-001",
        "price": 19.99,
        "discount": 5.0,
        "brand": "Logitech",
        "category": "Electronics"
      }
    ]
  }
}
```
