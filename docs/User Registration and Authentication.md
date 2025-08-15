### Overview

This code defines a user authentication and profile management system using the Go programming language, Gin web framework, and PostgreSQL (via sqlx). It is split into two main files:

- `controllers/UserController.go`: Handles HTTP requests/responses for user operations.
- `services/UserService.go`: Contains business logic and database interactions.

---

## `UserService`

### Fields:

- `DB`: A pointer to a `sqlx.DB` instance for database operations.

### Methods:

#### Login

- **Purpose**: Logs user in.
- **Inputs**: `*models.Credentials`
- **Returns**: `*models.User`, `string`, `string`, `error`
- **Key Operations**:
  - Fetches user by email.
  - Verifies account status.
  - Checks password validity.
  - Generates access and refresh tokens.
  - Returns the user and tokens.

#### Signup

- **Purpose**: Adds a new user.
- **Inputs**: `*models.User`
- **Returns**: `error`
- **Key Operations**:
  - Checks if the email is available and valid.
  - Generates a new UUID for `UserId`.
  - Hashes the password.
  - Inserts the new user into the database.
  - Sends a verification email.

#### VerifyUser

- **Purpose**: Verify user to allow logging in.
- **Inputs**: `verificationToken string`
- **Returns**: `error`
- **Key Operations**:
  - Fetches the user by the verification token.
  - Checks if already verified.
  - Sets `verified = true` and clears the token in the database.

#### UpdateUser

- **Purpose**: Updates user data.
- **Inputs**: `user *models.User`
- **Returns**: `*models.User`, `string`, `string`, `error`
- **Key Operations**:
  - Builds an SQL `SET` clause based on provided fields.
  - Verifies email format if provided.
  - Hashes the password if changed.
  - Updates the user in the database and returns new tokens if email is changed.

#### GetUserProfile

- **Purpose**: Fetches user data.
- **Inputs**: `userId uuid.UUID`
- **Returns**: `*models.User`, `error`
- **Key Operations**:
  - Accepts a `userId` UUID.
  - Returns the complete user profile from the database.

---

## `UserController`

### Fields:

- `UserService`: A pointer to the `UserService` struct which contains core business logic.

### Methods:

#### `Login`

- **Method**: `POST`
- **Path**: `/login`
- **Behavior**:
  - Binds incoming JSON credentials to a struct.
  - Calls the `Login` method from `UserService`.
  - Sets a `refreshToken` as an HTTP-only cookie.
  - Returns sanitized user data and an `accessToken` as JSON.

#### `Signup`

- **Method**: `POST`
- **Path**: `/signup`
- **Behavior**:
  - Binds incoming JSON user data to a struct.
  - Calls `Signup` from `UserService` to handle registration.
  - Returns a success message if email is sent successfully.

#### `VerifyUser`

- **Method**: `POST`
- **Path**: `/verify`
- **Behavior**:
  - Accepts a query parameter `verificationToken`.
  - Calls `VerifyUser` from `UserService` to activate the account.
  - Create new cart.

#### `UpdateUser`

- **Method**: `PATCH`
- **Path**: `/profile`
- **Behavior**:
  - Binds incoming JSON to a user object.
  - Calls `UpdateUser` from `UserService`.
  - If email is updated, new tokens are generated.
  - Returns the updated user and optionally new tokens.

#### `GetUserProfile`

- **Method**: `GET`
- **Path**: `/profile`
- **Behavior**:
  - Accepts a query parameter `userId`.
  - Parses `userId` into UUID.
  - Calls `GetUserProfile` from `UserService`.
  - Returns the user's sanitized profile.

---

### **Login Process**

1. **User sends credentials** (`email` and `password`) via a POST request.
2. The controller:
   - Parses the request body into a `Credentials` struct.
   - Calls `UserService.Login(...)`.
3. The service:
   - Fetches the user by email.
   - Verifies that the user is **already verified**.
   - Compares the **hashed password** with the stored one.
   - If valid, generates:
     - **Access token** (for short-term API access).
     - **Refresh token** (stored in an **HTTP-only cookie** for security).
4. The controller responds with:
   - A **sanitized user object** (excluding password, etc.).
   - The **access token** in the response body.

---

### **Signup Process**

1. **User sends registration data** (name, email, password, etc.) via a POST request.
2. The controller:
   - Parses the request body into a `User` struct.
   - Calls `UserService.Signup(...)`.
3. The service:
   - Checks if the **email is already used** and **valid**.
   - Hashes the password.
   - Generates a new UUID and **verification token**.
   - Inserts the user into the database with `verified = false`.
   - Sends a **verification email** with the token link.
4. The controller responds with a success message prompting the user to verify their email.

---

## Data Models in Golang

```go
// === === === === ===
//
//	  === User ===
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

```

## sQL Tables

```sQL
CREATE TABLE users (
  userid VARCHAR(36) PRIMARY KEY,
  name TEXT NOT NULL,
  email VARCHAR(50) NOT NULL UNIQUE,
  password TEXT NOT NULL,
  phone_number VARCHAR(15) NOT NULL UNIQUE,
  verified BOOLEAN DEFAULT FALSE,
  verification_token TEXT,
  role VARCHAR(6) NOT NULL DEFAULT 'user' CHECK (role IN ('vendor', 'user')),
  createdat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updatedat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```
