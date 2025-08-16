package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrNotFound     = errors.New("not found")
)

// === === === === ===
//
//	=== User ===
//
// === === === === ===
type User struct {
	UserId        uuid.UUID `json:"userId" db:"userid"`
	Name          string    `json:"name" db:"name"`
	Email         string    `json:"email" db:"email"`
	Password      string    `json:"password" db:"password"`
	PhoneNumber   string    `json:"phoneNumber" db:"phone_number"`
	Verified      bool      `json:"verified" db:"verified"`
	VerifiedToken string    `json:"verificationToken" db:"verification_token"`
	Role          string    `json:"role" db:"role"`
	CreatedAt     time.Time `json:"createdAt" db:"createdat"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updatedat"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SanitizedUser struct {
	UserId      uuid.UUID `json:"userId"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	Verified    bool      `json:"verified"`
}

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
	IsDefault    bool      `json:"isDefault" db:"is_default"`
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

// === === === === ===
//
//	=== Carts ===
//
// === === === === ===
type Cart struct {
	CartId    string    `json:"cartId" db:"cartid"`
	UserId    uuid.UUID `json:"userId" db:"userid"`
	CreatedAt time.Time `json:"createdAt" db:"createdat"`
	UpdatedAt time.Time `json:"updatedAt" db:"updatedat"`
}

type CartItem struct {
	CartItemId string    `json:"cartItemId" db:"cartitemid"`
	CartId     string    `json:"cartId" db:"cartid"`
	ProductId  string    `json:"ProductId" db:"productid"`
	Quantity   int       `json:"quantity" db:"quantity"`
	AddedAt    time.Time `json:"addedAt" db:"addedat"`
}

type ItemData struct {
	ProductId   string   `json:"ProductId" db:"productid"`
	Name        string   `json:"name" db:"name"`
	Description string   `json:"description" db:"description"`
	SKU         string   `json:"sku" db:"sku"`
	Price       *float64 `json:"price" db:"price"`
	Discount    *float64 `json:"discount" db:"discount"`
	Brand       string   `json:"brand" db:"brand"`
	Category    string   `json:"category" db:"category"`
	Quantity    int      `json:"quantity" db:"quantity"`
}

type ReturnedCart struct {
	ID     string
	UserId string
	Items  []*ItemData
	Total  float64
}

// === === === === ===
//
//	=== Checkout ===
//
// === === === === ===

type Summary struct {
	ReturnedCart    *ReturnedCart
	ShippingAddress *ShippingAddress
}

// === === === === ===
//
//	=== Orders ===
//
// === === === === ===

type OrderItem struct {
	ProductId   string
	Name        string
	Description string
	Brand       string
	Category    string
	Quantity    int
}

type Order struct {
	OrderId    string
	OrderedAt  time.Time
	Items      []*OrderItem
	TotalPrice float64
}

// === === === === ===
//
//	=== Reviews ===
//
// === === === === ===

type ReviewData struct {
	ReviewMessage string `json:"message" db:"message"`
	Stars         int    `json:"stars" db:"stars"`
}

type Review struct {
	ReviewId      string    `json:"reviewId" db:"reviewid"`
	UserId        uuid.UUID `json:"userId" db:"userid"`
	VendorId      uuid.UUID `json:"vendorId" db:"vendorid"`
	ProductId     string    `json:"ProductId" db:"productid"`
	ReviewMessage string    `json:"message" db:"message"`
	Stars         int       `json:"stars" db:"stars"`
	ReviewedAt    time.Time `json:"reviewedAt" db:"reviewedat"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updatedat"`
}

// === === === === ===
//
//	=== Wishlists ===
//
// === === === === ===

type Wishlist struct {
	WishlistId  string    `json:"wishlistId" db:"wishlistid"`
	UserId      uuid.UUID `json:"userId" db:"userid"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"createdAt" db:"createdat"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updatedat"`
}

type WishlistInfo struct {
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

type ReturnedWishlist struct {
	ID     string
	UserId string
	Items  []*ItemData
}
