package loans

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

// RegisterRoutes registers all application routes for loan management
func RegisterRoutes(opt *Options) {
	loanController := LoanController{Options: opt}

	v1 := opt.GinEngine.Group("/api/v1", auth.TokenAuthMiddleware(opt.TokenManager))
	{
		v1.POST("/loan-accounts", loanController.CreateLoanAccount)
		v1.GET("/loan-accounts/:id", loanController.GetLoanAccount)
		v1.GET("/loan-schedules/:loan_id", loanController.GetLoanSchedule)
		v1.GET("/loan-eligibility/:customer_id", loanController.GetLoanEligibility)
		v1.GET("/loan-accounts", loanController.ListLoanAccounts)
	}
}
