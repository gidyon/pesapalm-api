package loans_product

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

// RegisterRoutes registers all application routes for loan products
func RegisterRoutes(opt *Options) {
	productController := LoanProductController{Options: opt}

	v1 := opt.GinEngine.Group("/api/v1", auth.TokenAuthMiddleware(opt.TokenManager))
	{
		v1.POST("/loan-products", productController.CreateLoanProduct)
		v1.GET("/loan-products/:id", productController.GetLoanProduct)
		v1.PUT("/loan-products/:id", productController.UpdateLoanProduct)
		v1.DELETE("/loan-products/:id", productController.DeleteLoanProduct)
		v1.GET("/loan-products", productController.ListLoanProducts)
	}
}
