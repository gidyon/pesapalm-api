package loans_product

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LoanProductController structure
type LoanProductController struct {
	*Options
}

// CreateLoanProduct creates a new loan product
func (ctrl *LoanProductController) CreateLoanProduct(c *gin.Context) {
	var dto CreateLoanProductDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := LoanProduct{
		Name:                      dto.Name,
		ProductID:                 dto.ProductID,
		CurrencyID:                dto.CurrencyID,
		Description:               dto.Description,
		MaxLoanAmount:             dto.MaxLoanAmount,
		MaxInstallments:           dto.MaxInstallments,
		MinInstallments:           dto.MinInstallments,
		InterestRate:              dto.InterestRate,
		InterestCalculationPeriod: dto.InterestCalculationPeriod,
		InterestCalculationUnit:   dto.InterestCalculationUnit,
		RepaymentPeriod:           dto.RepaymentPeriod,
		RepaymentPeriodUnit:       dto.RepaymentPeriodUnit,
	}

	if result := ctrl.DB.Create(&product); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetLoanProduct retrieves a loan product by ID
func (ctrl *LoanProductController) GetLoanProduct(c *gin.Context) {
	id := c.Param("id")
	var product LoanProduct

	if result := ctrl.DB.First(&product, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Loan product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, ToLoanProductResponse(&product))
}

// UpdateLoanProduct updates an existing loan product
func (ctrl *LoanProductController) UpdateLoanProduct(c *gin.Context) {
	id := c.Param("id")
	var dto UpdateLoanProductDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product LoanProduct
	if result := ctrl.DB.First(&product, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Loan product not found"})
		return
	}

	// Update product details
	product.Name = dto.Name
	product.Description = dto.Description
	product.MaxLoanAmount = dto.MaxLoanAmount
	product.MaxInstallments = dto.MaxInstallments
	product.MinInstallments = dto.MinInstallments
	product.InterestRate = dto.InterestRate
	product.InterestCalculationPeriod = dto.InterestCalculationPeriod
	product.InterestCalculationUnit = dto.InterestCalculationUnit
	product.RepaymentPeriod = dto.RepaymentPeriod
	product.RepaymentPeriodUnit = dto.RepaymentPeriodUnit

	if result := ctrl.DB.Save(&product); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteLoanProduct deletes a loan product
func (ctrl *LoanProductController) DeleteLoanProduct(c *gin.Context) {
	id := c.Param("id")
	if result := ctrl.DB.Delete(&LoanProduct{}, id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Loan product deleted"})
}

// ListLoanProducts lists all loan products
func (ctrl *LoanProductController) ListLoanProducts(c *gin.Context) {
	var products []LoanProduct

	if result := ctrl.DB.Find(&products); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}
