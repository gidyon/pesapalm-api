package user

import (
	"github.com/gidyon/pesapalm/internal/auth"
)

func (api *APIServer) registerRoutes() {
	api.GinEngine.POST("/api/login", api.Login)
	api.GinEngine.POST("/api/requestOtp", api.RequestOtp)
	api.GinEngine.POST("/api/validateOtp", api.ValidateOtp)
	api.GinEngine.POST("/api/request-password-reset-otp", api.RequestResetPasswordOtp)
	api.GinEngine.POST("/api/reset-password", api.ResetPassword)
	api.GinEngine.POST("/api/refresh", auth.TokenAuthMiddleware(api.TokenManager), api.RefreshSession)

	// api.GinEngine.POST("/api/users", api.CreateUser)

	userGroup := api.GinEngine.Group("/api/v1/users", auth.TokenAuthMiddleware(api.TokenManager))
	{
		userGroup.POST("", api.CreateUser)
		userGroup.GET("", api.ListUsers)
		userGroup.GET("/:userId", api.GetUser)
		userGroup.PATCH("/:userId", api.UpdateUser)
		userGroup.POST("/logout", api.Logout)
	}
}
