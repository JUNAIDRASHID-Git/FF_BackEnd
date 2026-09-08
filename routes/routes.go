package routes

import (
	"funfillers/backend/config"
	"funfillers/backend/controllers"
	"funfillers/backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg config.Config) *gin.Engine {
	r := gin.Default()

	// Enable CORS
	r.Use(middleware.CORSMiddleware())

	// Static route to serve uploaded images
	r.Static("/uploads", "./uploads")

	productCtrl := controllers.NewProductController()
	authCtrl := controllers.NewAuthController(cfg)
	categoryCtrl := controllers.NewCategoryController()
	orderCtrl := controllers.NewOrderController()
	adminCtrl := controllers.NewAdminController()
	uiCtrl := controllers.NewUIController()
	wishlistCtrl := controllers.NewWishlistController()

	// Root share preview for social media bots (WhatsApp, FB, Twitter)
	r.GET("/share/product/:id", productCtrl.GetProductSharePreview)

	api := r.Group("/api")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			dbStatus := "disconnected (using memory fallback)"
			if config.DB != nil {
				dbStatus = "connected (PostgreSQL)"
			}
			c.JSON(200, gin.H{
				"status":   "healthy",
				"service":  "funfillers-backend",
				"database": dbStatus,
				"version":  "1.0.0",
			})
		})

		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authCtrl.Register)
			auth.POST("/login", authCtrl.Login)
			auth.POST("/google", authCtrl.GoogleAuth)
			auth.POST("/phone/send-otp", authCtrl.SendOTP)
			auth.POST("/phone/verify-otp", authCtrl.VerifyOTP)
		}

		// Public catalog & UI routes
		api.GET("/products", productCtrl.GetProducts)
		api.GET("/products/:id", productCtrl.GetProductByID)
		api.GET("/categories", categoryCtrl.GetCategories)
		api.GET("/ui/banners", uiCtrl.GetBanners)

		// Wishlist routes
		api.GET("/wishlist", wishlistCtrl.GetWishlist)
		api.POST("/wishlist/toggle", wishlistCtrl.ToggleWishlist)
		api.DELETE("/wishlist/:productId", wishlistCtrl.RemoveWishlistItem)

		// Admin & Public Store API Endpoints (Admin panel direct access)
		admin := api.Group("/admin")
		{
			// Stats & Users
			admin.GET("/stats", adminCtrl.GetStats)
			admin.GET("/users", adminCtrl.GetUsers)
			admin.PUT("/users/:id/status", adminCtrl.UpdateUserStatus)
			admin.PUT("/users/:id/role", adminCtrl.UpdateUserRole)
			admin.POST("/users/request", adminCtrl.RequestAdminAccess)

			// File / Image Upload
			admin.POST("/upload", productCtrl.UploadImage)

			// Products Management
			admin.GET("/products", productCtrl.GetProducts)
			admin.POST("/products", productCtrl.CreateProduct)
			admin.PUT("/products/:id", productCtrl.UpdateProduct)
			admin.DELETE("/products/:id", productCtrl.DeleteProduct)

			// Categories & SubCategories Management
			admin.GET("/categories", categoryCtrl.GetCategories)
			admin.POST("/categories", categoryCtrl.CreateCategory)
			admin.DELETE("/categories/:id", categoryCtrl.DeleteCategory)
			admin.POST("/subcategories", categoryCtrl.CreateSubCategory)

			// UI Banners Management
			admin.GET("/banners", uiCtrl.GetBanners)
			admin.POST("/banners", uiCtrl.CreateBanner)
			admin.PUT("/banners/:id/toggle", uiCtrl.ToggleBannerStatus)
			admin.DELETE("/banners/:id", uiCtrl.DeleteBanner)

			// Compliance & Safety Management
			admin.POST("/batches/recall", productCtrl.RecallBatch)
			admin.GET("/compliance/certifications/expiring", productCtrl.GetExpiringCertifications)

			// Orders Management
			admin.GET("/orders", adminCtrl.GetAllOrders)
			admin.PUT("/orders/:id/status", adminCtrl.UpdateOrderStatus)
		}

		// Protected customer routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			protected.POST("/orders", orderCtrl.CreateOrder)
			protected.GET("/orders/user", orderCtrl.GetUserOrders)
		}
	}

	return r
}
