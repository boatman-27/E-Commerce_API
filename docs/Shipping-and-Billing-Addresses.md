# Shipping & Billing Addresses

### Overview

This code defines the `BillingService`, `BillingController`, `ShippingService` and `ShippingController` components of an eCommerce application that manages billing addresses. The system allows users to add, retrieve, update, and delete their billing and shipping addresses.
It is split into 4 main files:

- `controllers/BillingController.go`: Handles HTTP requests/responses for Billing Address operations.
- `services/BillingService.go`: Contains business logic and database interactions.
- `controllers/ShippingController.go`: Handles HTTP requests/responses for Shipping Address operations.
- `services/ShippingService.go`: Contains business logic and database interactions.

---

## Billing Address Operations

## `BillingService`

### Fields:

- `DB`: A pointer to a `sqlx.DB` instance for database operations.

### Methods:

#### AddBillingAddress

- **Purpose**: Add a new billing address for a user.
- **Inputs**: `*models.BillingAddress`, `userId string`
- **Returns**: `*models.BillingAddress`, `error`
- **Key Operations**:
  - Starts a database transaction.
  - If new address is set as default, unset old default address.
  - Inserts new billing address into the database.
  - Commits transaction.

#### DeleteBillingAddress

- **Purpose**: Delete a billing address owned by the user.
- **Inputs**: `billingId string`, `userId string`
- **Returns**: `error`
- **Key Operations**:
  - Deletes the billing address where `billingId` and `userId` match.
  - Ensures ownership is respected.

#### GetBillingAddresses

- **Purpose**: Fetch all billing addresses for a user.
- **Inputs**: `userId string`
- **Returns**: `[]*models.BillingAddress`, `error`
- **Key Operations**:
  - Queries the database for all billing addresses associated with `userId`.

#### UpdateBillingAddress

- **Purpose**: Update a billing address belonging to a user.
- **Inputs**: `*models.BillingAddress`, `userId string`, `billingId string`
- **Returns**: `*models.BillingAddress`, `error`
- **Key Operations**:
  - Verifies ownership of billing address.
  - Dynamically builds `SET` clause based on provided fields.
  - Updates the billing address in the database.

#### GetDefaultBillingAddress

- **Purpose**: Fetch the default billing address for a user.
- **Inputs**: `userId string`
- **Returns**: `*models.BillingAddress`, `error`
- **Key Operations**:
  - Queries the database for the billing address with `is_default = true`.

#### ChangeDefaultBillingAddress

- **Purpose**: Change which billing address is set as default.
- **Inputs**: `userId string`, `oldDefault string`, `newDefault string`
- **Returns**: `error`
- **Key Operations**:
  - Checks if old and new default addresses exist for the user.
  - Starts transaction.
  - Unsets the old default.
  - Sets new default.
  - Commits transaction.

---

## `BillingController`

### Fields:

- `billingService`: A pointer to the `BillingService` struct which contains the core business logic

### Methods:

#### AddBillingAddress

- **Method**: `POST`
- **Path**: `/billing`
- **Behavior**:
  - Parses request body into `BillingAddress` struct.
  - Extracts `userId` from context.
  - Calls `AddBillingAddress` in service layer.
  - Returns newly added billing address.

#### DeleteBillingAddress

- **Method**: `DELETE`
- **Path**: `/billing?billingId=<id>`
- **Behavior**:
  - Extracts `billingId` from query params.
  - Extracts `userId` from context.
  - Calls `DeleteBillingAddress` in service layer.
  - Returns success message if deleted.

#### GetBillingAddresses

- **Method**: `GET`
- **Path**: `/billing`
- **Behavior**:
  - Extracts `userId` from context.
  - Calls `GetBillingAddresses` in service layer.
  - Returns all billing addresses.

#### UpdateBillingAddress

- **Method**: `PATCH`
- **Path**: `/billing?billingId=<id>`
- **Behavior**:
  - Parses request body into `BillingAddress` struct.
  - Extracts `billingId` from query params.
  - Extracts `userId` from context.
  - Calls `UpdateBillingAddress` in service layer.
  - Returns updated billing address.

#### GetDefaultBillingAddress

- **Method**: `GET`
- **Path**: `/billing/default`
- **Behavior**:
  - Extracts `userId` from context.
  - Calls `GetDefaultBillingAddress` in service layer.
  - Returns the default billing address.

#### ChangeDefaultBillingAddress

- **Method**: `PATCH`
- **Path**: `/billing/default?oldDefault=<id>&newDefault=<id>`
- **Behavior**:
  - Extracts `oldDefault` and `newDefault` from query params.
  - Extracts `userId` from context.
  - Calls `ChangeDefaultBillingAddress` in service layer.
  - Returns success message.

---

## Shipping Address Operations

### `ShippingService`

### Fields

- `DB`: A pointer to a `sqlx.DB` instance for database operations.

### Methods:

### `AddShippingAddress`

- **Purpose**: Inserts a new shipping address for a user.
- **Inputs**: `*models.SillingAddress`, `userId string`
- **Returns**: The inserted `ShippingAddress`, or error
- **Key Operations**:
  - Executes an `INSERT INTO` SQL query

### `DeleteShippingAddress`

- **Purpose**: Deletes a shipping address belonging to a specific user.
- **Inputs**: `shippingId string`, `userId string`
- **Returns**: `error`
- **Key Operations**:
  - Verifies if the shipping address exists and belongs to the user
  - Deletes the shipping address if ownership is confirmed

### `GetShippingAddresses`

- **Purpose**: Fetches all shipping addresses associated with a user.
- **Inputs**: `userId string`
- **Returns**: List of `ShippingAddress` or error
- **Key Operations**:
  - Executes `SELECT * FROM shipping_addresses WHERE userid = $1`

### `UpdateShippingAddress`

- **Purpose**: Updates an existing shipping address for a user.
- **Inputs**: `*models.ShippingAddress`, `userId string`
- **Returns**: The updated `ShippingAddress`, or error
- **Key Operations**:
  - Verifies ownership of the shipping address
  - Dynamically constructs `SET` clauses based on which fields are provided
  - Executes an `UPDATE` query and returns the updated record

### `GetDefaultShippingAddress`

- **Purpose**: Fetches current default shipping address.
- **Inputs**: `userId string`
- **Returns**: default `*models.ShippingAddress`, or error
- **Key Operations**:
  - Verifies ownership of the shipping address
  - Dynamically constructs `SET` clauses based on which fields are provided
  - Executes `SELECT * FROM shipping_addresses WHERE userid = $1 AND is_default = TRUE`

### `ChangeDefaultShippingAddress`

- **Purpose**: Sets a new default shipping address and unsets the old address.
- **Inputs**: `userId, oldDefault, newDefault string`
- **Returns**: error
- **Key Operations**:
  - Begins a transaction
  - checks if old and new default addresses exists
  - unsets old default address and sets the new one
  - checks the number of affected rows
  - if there are no errors then it commits, if any then rollback and none of the operations are executed

---

### `ShippingController`

### Fields:

- `ShippingService`: A pointer to the `ShippingService` struct which contains core business logic.

### Methods:

### `AddShippingAddress`

- **Method**: `POST`
- **Path**: `/shipping`
- **Behavior**:
  - Parses and validates JSON input
  - Retrieves `userId` from context
  - Calls `ShippingService.AddShippingAddress`
  - Returns the created shipping address as JSON

### `DeleteShippingAddress`

- **Method**: `DELETE`
- **Path**: `/shipping?shippingId=<id>``
- **Behavior**:
  - Gets `shippingId` from query parameters
  - Verifies `userId` from context
  - Calls `ShippingService.DeleteshippingAddress`
  - Returns success message or appropriate error

### `GetShippingAddresses`

- **Method**: `GET`
- **Path**: `/shipping`
- **Behavior**:
  - Retrieves `userId` from context
  - Calls `ShippingService.GetShippingAddresses`
  - Returns the list of shipping addresses as JSON

### `UpdateShippingAddress`

- **Method**: `PATCH`
- **Path**: `/shipping?shippingId=<id>``
- **Behavior**:
  - Parses and validates JSON input
  - Checks for presence of `shippingId`
  - Retrieves `userId` from context
  - Calls `ShippingService.UpdateShippingAddress`
  - Returns the updated shipping address as JSON

### `GetDefaultShippingAddress`

- **Method**: `GET`
- **Path**: `/shipping/default`
- **Behavior**:
  - Retrieves `userId` from context
  - Calls `ShippingService.GetDefaultShippingAddress`
  - Returns the default shipping address

### `ChangeDefaultShippingAddress`

- **Method**: `PATCH`
- **Path**: `/shipping/default?oldDefault=<id>&newDefault=<id>``
- **Behavior**:
  - Retrieves `userId` from context
  - Gets `oldDefault`, `newDefault` from query parameters
  - Calls `ShippingService.ChangeDefaultShippingService`
  - returns a message if successful

---

## Data Models in Golang

```go
// === === === === ===
//
//	=== Billing ===
//
// === === === === ===
type BillingAddress struct {
	BillingId    string    `json:"billingId" db:"billingid"`
	UserId       uuid.UUID `json:"userId" db:"userid"`
	AddressLine1 string    `json:"addressLine1" db:"address_line1"`
	AddressLine2 string    `json:"addressLine2" db:"address_line2"`
	City         string    `json:"city" db:"city"`
	State        string    `json:"state" db:"state"`
	PostalCode   string    `json:"postalCode" db:"postal_code"`
	Country      string    `json:"country" db:"country"`
	CreatedAt    time.Time `json:"createdAt" db:"createdat"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updatedat"`
}

// === === === === ===
//
//	=== Shipping ===
//
// === === === === ===
type ShippingAddress struct {
	ShippingId   string    `json:"shippingId" db:"shippingid"`
	UserId       uuid.UUID `json:"userId" db:"userid"`
	AddressLine1 string    `json:"addressLine1" db:"address_line1"`
	AddressLine2 string    `json:"addressLine2" db:"address_line2"`
	City         string    `json:"city" db:"city"`
	State        string    `json:"state" db:"state"`
	PostalCode   string    `json:"postalCode" db:"postal_code"`
	Country      string    `json:"country" db:"country"`
	IsDefault    bool      `json:"isDefault" db:"is_default"`
	CreatedAt    time.Time `json:"createdAt" db:"createdat"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updatedat"`
}

```

## sQL Tables

```sQL
CREATE TABLE billing_addresses (
  billingid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  userid VARCHAR(36) NOT NULL REFERENCES users(userid) ON DELETE CASCADE,
  address_line1  TEXT NOT NULL,
  address_line2  TEXT,
  city           TEXT NOT NULL,
  state          TEXT,
  postal_code    TEXT NOT NULL,
  country        TEXT NOT NULL,
  is_default     BOOLEAN DEFAULT FALSE
  createdat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updatedat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Ensure only one default billing address per user
CREATE UNIQUE INDEX unique_default_billing_address_per_user
ON billing_addresses (userid)
WHERE is_default = TRUE;
```

```sQL
CREATE TABLE shipping_addresses (
  shippingid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  userid VARCHAR(36) NOT NULL REFERENCES users(userid) ON DELETE CASCADE,
  address_line1  TEXT NOT NULL,
  address_line2  TEXT,
  city           TEXT NOT NULL,
  state          TEXT,
  postal_code    TEXT NOT NULL,
  country        TEXT NOT NULL,
  is_default     BOOLEAN DEFAULT FALSE,
  createdat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updatedat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Ensure only one default shipping address per user
CREATE UNIQUE INDEX unique_default_shipping_address_per_user
ON shipping_addresses (userid)
WHERE is_default = TRUE;

```

## Dummy Data

```sQL
INSERT INTO billing_addresses ( userid, address_line1, address_line2, city, state, postal_code, country, is_default)
VALUES
	('dc22872c-f003-40df-b61a-743c97945b33', '321 Corniche Road', 'Apartment 8A', 'New Cairo', 'Cairo', '101', 'Egypt', FALSE),
	('dc22872c-f003-40df-b61a-743c97945b33', '654 random Street','', 'Seeb', '', '122', 'Oman', TRUE );
```

```sQL
INSERT INTO shipping_addresses (userid, address_line1, address_line2, city, state, postal_code, country, is_default)
VALUES
  -- User 1 addresses
  ('dc22872c-f003-40df-b61a-743c97945b33', '123 Main St', '', 'Muscat', 'Muscat', '11111', 'Oman', TRUE),
  ('dc22872c-f003-40df-b61a-743c97945b33', '456 Side Rd', 'Apt 2', 'Muscat', 'Muscat', '22222', 'Oman', FALSE),

  -- User 2 addresses
  ('6e1e78ec-257a-4918-8296-91ae444cb675', '789 Ocean Blvd', '', 'Salalah', 'Dhofar', '33333', 'Oman', TRUE),
  ('6e1e78ec-257a-4918-8296-91ae444cb675', '1011 Mountain Way', 'Suite 5', 'Salalah', 'Dhofar', '44444', 'Oman', FALSE);

```

## Example JSON

### Add Billing/Shipping Address Request

```json
{
  "addressLine1": "123 Main St",
  "addressLine2": "Apt 4B",
  "city": "Muscat",
  "state": "Muscat",
  "postalCode": "112",
  "country": "Oman",
  "isDefault": true
}
```

### Add Billing/Shipping Address Response

```json
{
  "Billing Address": {
    "billingId": "b01ecc65-7139-4643-be15-c6fb5c14e297",
    "userId": "dc22872c-f003-40df-b61a-743c97945b33",
    "addressLine1": "123 Main St",
    "addressLine2": "Apt 4B",
    "city": "Muscat",
    "state": "Muscat",
    "postalCode": "112",
    "country": "Oman",
    "isDefault": true,
    "createdAt": "2025-08-16T12:00:00Z",
    "updatedAt": "2025-08-16T12:00:00Z"
  }
}
```

### Get Billing/Shipping Addresses Response

```json
{
  "addresses": [
    {
      "billingId": "b01ecc65-7139-4643-be15-c6fb5c14e297",
      "userId": "dc22872c-f003-40df-b61a-743c97945b33",
      "addressLine1": "123 Main St",
      "addressLine2": "Apt 4B",
      "city": "Muscat",
      "state": "Muscat",
      "postalCode": "112",
      "country": "Oman",
      "isDefault": true
    }
  ]
}
```

### Update Billing/Shipping Address Request

```json
{
  "addressLine1": "456 Palm Ave",
  "city": "Seeb"
}
```

### Update Billing/Shipping Address Response

```json
{
  "Updated Address": {
    "billingId": "b01ecc65-7139-4643-be15-c6fb5c14e297",
    "userId": "dc22872c-f003-40df-b61a-743c97945b33",
    "addressLine1": "456 Palm Ave",
    "addressLine2": "Apt 4B",
    "city": "Seeb",
    "state": "Muscat",
    "postalCode": "112",
    "country": "Oman",
    "isDefault": true
  }
}
```

### Delete Billing Address/Shipping Response

```json
{
  "message": "Billing address deleted successfully"
}
```

### Change Default Billing/Shipping Address Response

```json
{
  "message": "updated new default billing address"
}
```
