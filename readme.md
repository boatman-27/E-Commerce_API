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

## Authentication

- Users sign up and receive JWT access and refresh tokens.
- Protected routes require a valid access token.
- Vendor routes are further secured with role-based authorization.

---

## Features

### 1. [User Registration and Authentication](User-Registration-and-Authentication.md)

- **Sign Up:** Users create an account with username, email, and password. Email verification confirms the account.
- **Login:** Registered users log in with email and password. Multi-factor authentication (MFA) is supported.

### 2.
