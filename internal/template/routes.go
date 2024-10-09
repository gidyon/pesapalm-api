package template

func RegisterRoutes(opt *Options) {
	templateController := TemplateController{Options: opt}

	v1 := opt.GinEngine.Group("/api/v1/templates")
	{
		// Admin Groups Routes
		v1.GET("/admin_groups", templateController.GetAdminGroups)
		v1.POST("/admin_groups", templateController.CreateAdminGroup)

		// Branch Routes
		v1.GET("/branches", templateController.GetBranches)
		v1.POST("/branches", templateController.CreateBranch)

		// Currency Routes
		v1.GET("/currencies", templateController.GetCurrencies)
		v1.POST("/currencies", templateController.CreateCurrency)

		// Language Routes
		v1.GET("/languages", templateController.GetLanguages)
		v1.POST("/languages", templateController.CreateLanguage)
	}
}
