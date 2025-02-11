package customer

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

// RegisterRoutes registers all application routes for customer management
func RegisterRoutes(opt *Options) {
	customerController := CustomerController{Options: opt}

	v1 := opt.GinEngine.Group("/api/v1", auth.TokenAuthMiddleware(opt.TokenManager))
	{
		v1.POST("/customers", customerController.CreateCustomer)
		v1.GET("/customers", customerController.ListCustomers)
		v1.GET("/customers/:id", customerController.GetCustomer)
		v1.PATCH("/customers/:id", customerController.UpdateCustomer)
		v1.DELETE("/customers/:id", customerController.DeleteCustomer)
		v1.GET("/customer-stats", customerController.GetStats)
	}
}
