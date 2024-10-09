package savings_product

import (
	"github.com/gidyon/pesapalm/internal/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/grpclog"
	"gorm.io/gorm"
)

type Options struct {
	DB           *gorm.DB
	Logger       grpclog.LoggerV2
	TokenManager auth.TokenInterface
	GinEngine    *gin.Engine
}

// RegisterRoutes registers all application routes
func RegisterRoutes(opt *Options) {
	productController := SavingsProductController{Options: opt}

	v1 := opt.GinEngine.Group("/api/v1", auth.TokenAuthMiddleware(opt.TokenManager))
	{
		// Routes for savings products
		v1.POST("/savings-products", productController.CreateSavingsProduct)
		v1.GET("/savings-products/:id", productController.GetSavingsProduct)
		v1.PUT("/savings-products/:id", productController.UpdateSavingsProduct)
		v1.DELETE("/savings-products/:id", productController.DeleteSavingsProduct)
		v1.GET("/savings-products", productController.ListSavingsProducts)
	}
}
