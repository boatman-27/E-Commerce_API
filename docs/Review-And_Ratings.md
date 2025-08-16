# Reviews Functionality

### Overview

This document explains the Go code implementing **Reviews** functionality in the eCommerce backend. The implementation is divided into:

- `services/ReviewService.go`: Business logic.
- `controllers/ReviewController.go`: HTTP routes using Gin.

---

## `ReviewService`

### Fields:

- `DB`: Pointer to `sqlx.DB` for database operations.

### Methods:

#### getVendorId

- **Purpose**: Fetches the vendor ID of a product.
- **Inputs**: `productId string`
- **Returns**: `vendorId string`, `error`
- **Key Operations**: Queries the `products` table.

#### checkOwnership

- **Purpose**: Checks if a review belongs to a user.
- **Inputs**: `reviewId string`, `userId string`
- **Returns**: `error`
- **Key Operations**:
  - Queries `reviews` for `userid`.
  - Returns `ErrNotFound` or `ErrUnauthorized` if applicable.

#### SubmitReview

- **Purpose**: Submits a review for a product.
- **Inputs**: `userId string`, `productId string`, `*models.ReviewData`
- **Returns**: `error`
- **Key Operations**:
  - Validates stars (1–5).
  - Prevents users from reviewing their own products.
  - Inserts into `reviews`.

#### DeleteReview

- **Purpose**: Deletes a review by its owner.
- **Inputs**: `userId string`, `reviewId string`
- **Returns**: `error`
- **Key Operations**: Deletes review and checks `RowsAffected`.

#### EditReview

- **Purpose**: Edits a review by its owner.
- **Inputs**: `reviewId string`, `userId string`, `*models.ReviewData`.
- **Returns**: `error`
- **Key Operations**:
  - Validates ownership.
  - Dynamically updates `message` and `stars`.
  - Updates `updatedat`.

#### GetReviews

- **Purpose**: Retrieves reviews submitted by a user.
- **Inputs**: `userId string`
- **Returns**: `[]*models.Review`, `error`
- **Key Operations**: Queries `reviews` by `userid`.

#### GetVendorReviews

- **Purpose**: Retrieves reviews for a vendor’s products.
- **Inputs**: `vendorId string`
- **Returns**: `[]*models.Review`, `error`
- **Key Operations**: Queries `reviews` by `vendorid`.

---

## `ReviewController`

### Fields:

- `reviewService`: Pointer to `ReviewService`.

### Methods:

#### SubmitReview

- **Method**: `POST`
- **Path**: `/reviews?productId={id}`
- **Behavior**:
  - Gets `UserId` from context.
  - Binds JSON to `ReviewData`.
  - Calls `SubmitReview`.
  - Returns success message.

#### DeleteReview

- **Method**: `DELETE`
- **Path**: `/reviews?reviewId={id}`
- **Behavior**:
  - Deletes review using `DeleteReview`.
  - Returns success message.

#### EditReview

- **Method**: `PATCH`
- **Path**: `/reviews?reviewId={id}`
- **Behavior**:
  - Edits review using `EditReview`.
  - Returns success message.

#### GetReviews

- **Method**: `GET`
- **Path**: `/reviews`
- **Behavior**:
  - Returns reviews submitted by the current user.

#### GetVendorReviews

- **Method**: `GET`
- **Path**: `/vendor/reviews`
- **Behavior**:
  - Returns reviews received by the vendor.

---

## Data Models in Golang

```go
type ReviewData struct {
	ReviewMessage string `json:"message" db:"message"`
	Stars         int    `json:"stars" db:"stars"`
}

type Review struct {
	ReviewId      string    `json:"reviewId" db:"reviewid"`
	UserId        uuid.UUID `json:"userId" db:"userid"`
	VendorId      uuid.UUID `json:"vendorId" db:"vendorid"`
	ProductId     string    `json:"productId" db:"productid"`
	ReviewMessage string    `json:"message" db:"message"`
	Stars         int       `json:"stars" db:"stars"`
	ReviewedAt    time.Time `json:"reviewedAt" db:"reviewedat"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updatedat"`
}
```

---

## sQL Table

```sQL
CREATE TABLE reviews (
	reviewid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	userid VARCHAR(36) NOT NULL REFERENCES users(userid) ON DELETE CASCADE,
	vendorid VARCHAR(36) NOT NULL REFERENCES users(userid) ON DELETE CASCADE,
	productid UUID NOT NULL REFERENCES products(productid) ON DELETE CASCADE,
    message TEXT NOT NULL,
	stars INT CHECK(stars BETWEEN 1 AND 5),
    reviewedat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updatedat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_reviews_userid ON reviews(userid);

ALTER TABLE reviews
ADD CONSTRAINT unique_user_product UNIQUE (userid, productid);
```

---

## Example JSON

### Submit Review Request

```json
{
  "message": "Excellent product, highly recommend!",
  "stars": 5
}
```

### Submit Review Response

```json
{
  "message": "Review Submitted successfully"
}
```

### Get Reviews Response

```json
{
  "reviews": [
    {
      "reviewId": "review-uuid-1",
      "userId": "u01abcde-1111-2222-3333-444455556666",
      "vendorId": "v01vendor-1234",
      "productId": "product-uuid-1",
      "message": "Great product!",
      "stars": 5,
      "reviewedAt": "2025-08-16T11:00:00Z",
      "updatedAt": "2025-08-16T11:00:00Z"
    }
  ]
}
```
