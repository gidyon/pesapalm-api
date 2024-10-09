package savings_product

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SavingsProductController structure
type SavingsProductController struct {
	*Options
}

// CreateSavingsProduct creates a new savings product
func (ctrl *SavingsProductController) CreateSavingsProduct(c *gin.Context) {
	var dto CreateSavingsProductDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := SavingsProduct{
		Name:                      dto.Name,
		ProductCode:               dto.ProductCode,
		CurrencyID:                dto.CurrencyID,
		Description:               dto.Description,
		InterestRate:              dto.InterestRate,
		InterestCalculationPeriod: dto.InterestCalculationPeriod,
		InterestCalculationUnit:   dto.InterestCalculationUnit,
	}

	if result := ctrl.DB.Create(&product); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetSavingsProduct retrieves a savings product by ID
func (ctrl *SavingsProductController) GetSavingsProduct(c *gin.Context) {
	id := c.Param("id")
	var product SavingsProduct

	if result := ctrl.DB.First(&product, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Savings product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, ToSavingsProductResponse(&product))
}

// UpdateSavingsProduct updates an existing savings product
func (ctrl *SavingsProductController) UpdateSavingsProduct(c *gin.Context) {
	id := c.Param("id")
	var dto UpdateSavingsProductDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product SavingsProduct
	if result := ctrl.DB.First(&product, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Savings product not found"})
		return
	}

	product.Name = dto.Name
	product.Description = dto.Description
	product.InterestRate = dto.InterestRate
	product.InterestCalculationPeriod = dto.InterestCalculationPeriod
	product.InterestCalculationUnit = dto.InterestCalculationUnit

	if result := ctrl.DB.Save(&product); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteSavingsProduct deletes a savings product
func (ctrl *SavingsProductController) DeleteSavingsProduct(c *gin.Context) {
	id := c.Param("id")
	if result := ctrl.DB.Delete(&SavingsProduct{}, id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Savings product deleted"})
}

// ListSavingsProducts lists all savings products
func (ctrl *SavingsProductController) ListSavingsProducts(c *gin.Context) {
	var products []SavingsProduct

	if result := ctrl.DB.Find(&products); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}
