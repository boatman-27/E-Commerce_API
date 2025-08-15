package routes

import (
	controllers "eCommerce/controller"
	"eCommerce/middlewares"
	"eCommerce/services"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func SetupRouter(router *gin.Engine, db *sqlx.DB) {
	// Services
	userService := services.NewUserService(db)
	billingService := services.NewBillingService(db)
	shippingService := services.NewShippingService(db)
	vendorService := services.NewVendorService(db)
	productService := services.NewProductService(db)
	cartService := services.NewCartService(db)
	checkoutService := services.NewCheckoutService(db, *cartService, *shippingService)
	orderService := services.NewOrderService(db)
	reviewService := services.NewReviewService(db)
	wishlistService := services.NewWishlistService(db)

	// Controllers
	userController := controllers.NewUserController(userService)
	billingController := controllers.NewBillingController(billingService)
	shippingController := controllers.NewShippingController(shippingService)
	vendorController := controllers.NewVendorController(vendorService)
	productController := controllers.NewProductController(productService)
	cartController := controllers.NewCartController(cartService)
	checkoutController := controllers.NewCheckoutController(checkoutService)
	orderController := controllers.NewOrderController(orderService)
	reviewController := controllers.NewReviewContoller(reviewService)
	wishlistController := controllers.NewWishlistController(wishlistService)

	// Authentication Routes
	accountRoutes := router.Group("/account")
	{
		accountRoutes.POST("/login", userController.Login)
		accountRoutes.POST("/signup", userController.Signup)
		accountRoutes.POST("/verify", userController.VerifyUser)
	}

	// Protected Routes
	protected := router.Group("/protected")
	protected.Use(middlewares.RequireAuth)
	{
		// user routes
		protected.GET("/profile", userController.GetUserProfile)
		protected.PATCH("/profile", userController.UpdateUser)

		// Billing routes
		protected.GET("/billing", billingController.GetBillingAddresses)
		protected.GET("/billing/default", billingController.GetDefaultBillingAddress)
		protected.POST("/billing", billingController.AddBillingAddress)
		protected.PATCH("/billing", billingController.UpdateBillingAddress)
		protected.PATCH("/billing/default", billingController.ChangeDefaultBillingAddress)
		protected.DELETE("/billing", billingController.DeleteBillingAddress)

		// shipping routes
		protected.GET("/shipping", shippingController.GetShippingAddresses)
		protected.GET("/shipping/default", shippingController.GetDefaultShippingAddress)
		protected.POST("/shipping", shippingController.AddShippingAddress)
		protected.PATCH("/shipping", shippingController.UpdateShippingAddress)
		protected.PATCH("/shipping/default", shippingController.ChangeDefaultShippingAddress)
		protected.DELETE("/shipping", shippingController.DeleteShippingAddress)

		// user product routes
		protected.GET("/products", productController.GetProducts)

		// cart routes
		protected.GET("/cart", cartController.ViewCart)
		protected.GET("/cart/total", cartController.GetTotalPrice)
		protected.POST("/cart", cartController.CreateCart)
		protected.POST("/cart/item", cartController.AddToCart)
		protected.PATCH("/cart/item", cartController.EditCartItem)
		protected.DELETE("/cart/item", cartController.DeleteCartItem)
		protected.DELETE("/cart", cartController.DeleteCart)

		// checkout routes
		protected.GET("/summary", checkoutController.OrderSummary)
		protected.GET("/confirm", checkoutController.ConfrimPurchase)

		// order routes
		protected.GET("/orders/status", orderController.TrackOrder)
		protected.GET("/orders", orderController.ViewOrders)

		// review routes
		protected.GET("/reviews", reviewController.GetReviews)
		protected.POST("/reviews", reviewController.SubmitReview)
		protected.PATCH("/reviews", reviewController.EditReview)
		protected.DELETE("/reviews", reviewController.DeleteReview)

		// wishlist routes
		protected.GET("/wishlists", wishlistController.GetWishlistItems)
		protected.POST("/wishlists", wishlistController.CreateWishlist)
		protected.POST("/wishlists/item", wishlistController.AddToWishlist)
		protected.POST("/wishlists/move-to-cart", wishlistController.MoveToCart)
		protected.POST("/wishlists/item/move-to-cart", wishlistController.MoveItemToCart)
		protected.PATCH("/wishlists", wishlistController.EditWishlist)
		protected.DELETE("/wishlists", wishlistController.DeleteWishlist)
		protected.DELETE("/wishlist-item", wishlistController.DeleteWishlistItem)

	}

	// Vendor Routes
	vendor := router.Group("/vendor")
	vendor.Use((middlewares.VendorAuth))
	{
		// review routes
		vendor.GET("/reviews", reviewController.GetVendorReviews)
		vendor.GET("/products", vendorController.GetVendorProducts)

		// vendor routes
		vendor.GET("/products/id", vendorController.GetProductById)
		vendor.POST("/products", vendorController.AddProduct)
		vendor.POST("/products", vendorController.DeleteProduct)
		vendor.PATCH("/products", vendorController.Updateproduct)
	}
}
