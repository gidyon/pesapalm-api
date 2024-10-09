package template

import (
	"net/http"

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

type TemplateController struct {
	*Options
}

// GetAdminGroups fetches all admin groups
func (ctrl *TemplateController) GetAdminGroups(c *gin.Context) {
	ctx := c.Request.Context()
	var adminGroups []AdminGroup
	if err := ctrl.DB.WithContext(ctx).Omit("column_roles").Find(&adminGroups).Error; err != nil {
		ctrl.Options.Logger.Error("Error fetching admin groups:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching admin groups"})
		return
	}
	c.JSON(http.StatusOK, adminGroups)
}

// CreateAdminGroup creates a new admin group
func (ctrl *TemplateController) CreateAdminGroup(c *gin.Context) {
	ctx := c.Request.Context()
	var adminGroupDTO AdminGroupDTO
	if err := c.ShouldBindJSON(&adminGroupDTO); err != nil {
		ctrl.Options.Logger.Error("Invalid admin group data:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminGroup := AdminGroup{
		Name:        adminGroupDTO.Name,
		ColumnRoles: adminGroupDTO.ColumnRoles,
		Level:       adminGroupDTO.Level,
	}

	if err := ctrl.DB.WithContext(ctx).Create(&adminGroup).Error; err != nil {
		ctrl.Options.Logger.Error("Error creating admin group:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin group"})
		return
	}

	c.JSON(http.StatusCreated, adminGroup)
}

// GetBranches fetches all branches
func (ctrl *TemplateController) GetBranches(c *gin.Context) {
	ctx := c.Request.Context()
	var branches []Branch
	if err := ctrl.DB.WithContext(ctx).Find(&branches).Error; err != nil {
		ctrl.Options.Logger.Error("Error fetching branches:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching branches"})
		return
	}
	c.JSON(http.StatusOK, branches)
}

// CreateBranch creates a new branch
func (ctrl *TemplateController) CreateBranch(c *gin.Context) {
	ctx := c.Request.Context()
	var branchDTO BranchDTO
	if err := c.ShouldBindJSON(&branchDTO); err != nil {
		ctrl.Options.Logger.Error("Invalid branch data:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	branch := Branch{
		Name:         branchDTO.Name,
		Status:       branchDTO.Status,
		Notes:        branchDTO.Notes,
		Msisdn:       branchDTO.Msisdn,
		EmailAddress: branchDTO.EmailAddress,
	}

	if err := ctrl.DB.WithContext(ctx).Create(&branch).Error; err != nil {
		ctrl.Options.Logger.Error("Error creating branch:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create branch"})
		return
	}

	c.JSON(http.StatusCreated, branch)
}

// GetCurrencies fetches all currencies
func (ctrl *TemplateController) GetCurrencies(c *gin.Context) {
	ctx := c.Request.Context()
	var currencies []Currency
	if err := ctrl.DB.WithContext(ctx).Find(&currencies).Error; err != nil {
		ctrl.Options.Logger.Error("Error fetching currencies:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching currencies"})
		return
	}
	c.JSON(http.StatusOK, currencies)
}

// CreateCurrency creates a new currency
func (ctrl *TemplateController) CreateCurrency(c *gin.Context) {
	ctx := c.Request.Context()
	var currencyDTO CurrencyDTO
	if err := c.ShouldBindJSON(&currencyDTO); err != nil {
		ctrl.Options.Logger.Error("Invalid currency data:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currency := Currency{
		CurrencyCode: currencyDTO.CurrencyCode,
		CurrencyName: currencyDTO.CurrencyName,
		StatusID:     currencyDTO.StatusID,
	}

	if err := ctrl.DB.WithContext(ctx).Create(&currency).Error; err != nil {
		ctrl.Options.Logger.Error("Error creating currency:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create currency"})
		return
	}

	c.JSON(http.StatusCreated, currency)
}

// GetLanguages fetches all languages
func (ctrl *TemplateController) GetLanguages(c *gin.Context) {
	ctx := c.Request.Context()
	var languages []Language
	if err := ctrl.DB.WithContext(ctx).Find(&languages).Error; err != nil {
		ctrl.Options.Logger.Error("Error fetching languages:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching languages"})
		return
	}
	c.JSON(http.StatusOK, languages)
}

// CreateLanguage creates a new language
func (ctrl *TemplateController) CreateLanguage(c *gin.Context) {
	ctx := c.Request.Context()
	var languageDTO LanguageDTO
	if err := c.ShouldBindJSON(&languageDTO); err != nil {
		ctrl.Options.Logger.Error("Invalid language data:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	language := Language{
		LanguageCode: languageDTO.LanguageCode,
		LanguageName: languageDTO.LanguageName,
		StatusID:     languageDTO.StatusID,
	}

	if err := ctrl.DB.WithContext(ctx).Create(&language).Error; err != nil {
		ctrl.Options.Logger.Error("Error creating language:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create language"})
		return
	}

	c.JSON(http.StatusCreated, language)
}
