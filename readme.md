# E-Commerce API

---

The project aims to create a powerful and scalable API that supports all essential functionalities for a modern e-commerce platform. This API will allow users to browse products, add them to their cart, and complete the checkout process. Additionally, users can create accounts, save shipping and billing information, and track their orders.

> This implementation is inspired by the [MB Projects E-Commerce Platform API: Build Your Own eBay](https://projects.masteringbackend.com/projects/e-commerce-platform-api-build-your-own-e-bay). For a deeper dive into the design and architecture.

## Tech Stack

- **Language**: Go (Golang)
- **Framework**: [Gin](https://github.com/gin-gonic/gin)
- **Database**: PostgreSQL
- **ORM**: [sqlx](https://github.com/jmoiron/sqlx)
- **Authentication**: JWT (Access and Refresh Tokens)

---

## Quick start

### 1. Clone the repository

```bash
git clone https://github.com/boatman-27/E-Commerce_API
cd E-Commerce_API
```

### 2. Setup environment variables

Create a `.env` file and add your secrets and DB config:

```env
# ----------------------------
	# Database Configuration
# ----------------------------
POSTGRES_DSN="user=<DB_USERNAME> password=<DB_PASSWORD> dbname=<DB_NAME> sslmode=disable"

# ----------------------------
# Application Email Configuration
# ----------------------------
# Email used for sending notifications, verifications, etc.
APP_EMAIL="your-email@gmail.com"

# Email-specific password or token
APP_PASSWORD="<APP_PASSWORD_OR_TOKEN>"

# Access & Refresh token secrets
ACCESS_SECRET=your_access_secret
REFRESH_SECRET=your_refresh_secret

# ----------------------------
# Notes:
# 1. Replace placeholders with your actual credentials.
# 2. For Gmail, create an App Password instead of using your main password.
```

### 3. Install CompileDaemon (Optional)

Once installed, you can use `CompileDaemon` to watch your project files and recompile/run your Go application automatically on changes.

```Bash
go install github.com/githubnemo/CompileDaemon@latest
```

Add `export PATH="$HOME/go/bin:$PATH"` to your `shell` file (`.zshrc` or `.bashrc`)

### 4. Add an alias in your `shell` file

```bash
alias gomon="CompileDaemon -build='go build -o myapp main.go' -command='./myapp'"
```

### 5. Run the application

```bash
# With CompileDaemon
gomon

# Or without
go run .
```

---

## Features

### 1. [User Registration and Authentication](docs/User-Registration-and-Authentication.md)

- **Sign Up:** Users create an account with username, email, and password. Email verification confirms the account.
- **Login:** Registered users log in with email and password. Multi-factor authentication (MFA) is supported.

### 2. [Shipping & Billing Addresses](docs/Shipping-and-Billing-Addresses.md)

- **Manage Addresses**: Users can add, edit, and delete shipping and billing addresses. This ensures that users have up-to-date information for order deliveries.

### 3. [Vendor Processes](docs/Vendor-Processes.md)

- **Add Product**: Vendors can create new products with details like name, price, stock, and category, which are stored in the database with ownership linked to their vendor ID.
- **Manage Products**: Vendors can view, update, or delete their products, with validations ensuring only the owner can modify or remove them.

### 4. [Product Browsing & Search](docs/Product-Browsing-and-Search.md)

- **Product Listing**: Users can browse a catalog of products categorized by type, brand, or other criteria. Each product listing includes images, descriptions, prices, and availability status.
- **Search Products**: Users can quickly find what they need by searching for products using keywords, filters, and sorting options.

### 5. [Cart Management System](docs/Cart-Management-System.md)

- **Add to Cart**: Users can add products to their shopping cart, specify quantities, and choose product options.
- **View and Edit Cart**: Users can view the contents of their cart, update quantities, and remove items.

### 6. [Checkout & Order System](docs/Checkout-And-Order-System.md)

- **Checkout Process**:
  - **Order Summary**: Before checking out, users can review their order, including product details, shipping information, and total cost.
  - **Order Confirmation**: Users receive an order confirmation with a purchase summary.
- **Order Tracking**:
  - **Order History**: Users can view past orders, including order details and status.
  - **Track Order**: Users can track the status of their current orders, from processing to shipping and delivery.

### 7. [Wishlist Process](docs/Wishlist-Process.md)

- **Create and Manage Wishlists**: Users can save products to one or more wishlists for future purchases.
- **Move Items to Cart**: Users can move items from their wishlist to their shopping cart for easy purchasing.

### 8. [Review & Ratings](docs/Review-And_Ratings.md)

- **Submit Reviews**: Users can write reviews and rate products they have purchased.
- **View Reviews**: Users can read reviews and ratings from other customers to make informed purchasing decisions.

---

## Endpoints

| **Method** | **Endpoint**                           | **Description**                   | **Access**    |
| ---------- | -------------------------------------- | --------------------------------- | ------------- |
| POST       | /account/login                         | User login                        | Public        |
| POST       | /account/signup                        | User signup                       | Public        |
| POST       | /account/verify                        | Verify user account               | Public        |
| GET        | /protected/profile                     | Get user profile                  | Authenticated |
| PATCH      | /protected/profile                     | Update user profile               | Authenticated |
| GET        | /protected/billing                     | Get all billing addresses         | Authenticated |
| GET        | /protected/billing/default             | Get default billing address       | Authenticated |
| POST       | /protected/billing                     | Add new billing address           | Authenticated |
| PATCH      | /protected/billing                     | Update billing address            | Authenticated |
| PATCH      | /protected/billing/default             | Change default billing address    | Authenticated |
| DELETE     | /protected/billing                     | Delete billing address            | Authenticated |
| GET        | /protected/shipping                    | Get all shipping addresses        | Authenticated |
| GET        | /protected/shipping/default            | Get default shipping address      | Authenticated |
| POST       | /protected/shipping                    | Add new shipping address          | Authenticated |
| PATCH      | /protected/shipping                    | Update shipping address           | Authenticated |
| PATCH      | /protected/shipping/default            | Change default shipping address   | Authenticated |
| DELETE     | /protected/shipping                    | Delete shipping address           | Authenticated |
| GET        | /protected/products                    | Get all products                  | Authenticated |
| GET        | /protected/cart                        | View cart                         | Authenticated |
| GET        | /protected/cart/total                  | Get total cart price              | Authenticated |
| POST       | /protected/cart                        | Create cart                       | Authenticated |
| POST       | /protected/cart/item                   | Add item to cart                  | Authenticated |
| PATCH      | /protected/cart/item                   | Edit cart item                    | Authenticated |
| DELETE     | /protected/cart/item                   | Delete cart item                  | Authenticated |
| DELETE     | /protected/cart                        | Delete entire cart                | Authenticated |
| GET        | /protected/summary                     | Order summary                     | Authenticated |
| POST       | /protected/confirm                     | Confirm purchase                  | Authenticated |
| GET        | /protected/orders/status               | Track order status                | Authenticated |
| GET        | /protected/orders                      | View all orders                   | Authenticated |
| GET        | /protected/reviews                     | Get all reviews                   | Authenticated |
| POST       | /protected/reviews                     | Submit a review                   | Authenticated |
| PATCH      | /protected/reviews                     | Edit a review                     | Authenticated |
| DELETE     | /protected/reviews                     | Delete a review                   | Authenticated |
| GET        | /protected/wishlists                   | Get wishlist items                | Authenticated |
| POST       | /protected/wishlists                   | Create a wishlist                 | Authenticated |
| POST       | /protected/wishlists/item              | Add item to wishlist              | Authenticated |
| POST       | /protected/wishlists/move-to-cart      | Move entire wishlist to cart      | Authenticated |
| POST       | /protected/wishlists/item/move-to-cart | Move single wishlist item to cart | Authenticated |
| PATCH      | /protected/wishlists                   | Edit wishlist                     | Authenticated |
| DELETE     | /protected/wishlists                   | Delete wishlist                   | Authenticated |
| DELETE     | /protected/wishlist/item               | Delete wishlist item              | Authenticated |
| GET        | /vendor/reviews                        | Get reviews for vendor            | Vendor        |
| GET        | /vendor/products                       | Get all vendor products           | Vendor        |
| GET        | /vendor/products/id                    | Get product by ID                 | Vendor        |
| POST       | /vendor/products                       | Add new product                   | Vendor        |
| POST       | /vendor/products/id                    | Delete product by ID              | Vendor        |
| PATCH      | /vendor/products                       | Update product                    | Vendor        |

---

## Contributing

Contributions welcome! Fork the repo, make changes on a feature branch, and submit a pull request.
