package savings

import (
	"github.com/gidyon/pesapalm/internal/auth"
)

// RegisterRoutes registers all application routes
func RegisterRoutes(opt *Options) {
	savingsController := SavingsAccountController{Options: opt}

	v1 := opt.GinEngine.Group("/api/v1", auth.TokenAuthMiddleware(opt.TokenManager))
	{
		v1.POST("/savings", savingsController.CreateSavingsAccount)
		v1.GET("/savings", savingsController.ListSavingsAccounts)
		v1.GET("/savings/:id", savingsController.GetSavingsAccount)
		v1.PUT("/savings/:id", savingsController.UpdateSavingsAccount)
		v1.PATCH("/savings/:id/status", savingsController.UpdateSavingsAccountStatus)
		v1.DELETE("/savings/:id", savingsController.DeleteSavingsAccount)
		v1.GET("/saving-stats", savingsController.GetStats)
	}
}
