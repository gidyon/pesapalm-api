package savings_product

import "time"

// CreateSavingsProductDTO defines the JSON structure for creating a new savings product
type CreateSavingsProductDTO struct {
	Name                      string  `json:"name"`
	ProductCode               string  `json:"product_code"`
	CurrencyID                int     `json:"currency_id"`
	Description               string  `json:"description"`
	InterestRate              float64 `json:"interest_rate"`
	InterestCalculationPeriod int     `json:"interest_calculation_period"`
	InterestCalculationUnit   string  `json:"interest_calculation_unit"`
}

// UpdateSavingsProductDTO defines the JSON structure for updating a savings product
type UpdateSavingsProductDTO struct {
	Name                      string  `json:"name"`
	Description               string  `json:"description"`
	InterestRate              float64 `json:"interest_rate"`
	InterestCalculationPeriod int     `json:"interest_calculation_period"`
	InterestCalculationUnit   string  `json:"interest_calculation_unit"`
}

// SavingsProductResponse defines the structure of the savings product data returned in the response
type SavingsProductResponse struct {
	ID                        uint    `json:"id"`
	Name                      string  `json:"name"`
	ProductCode               string  `json:"product_code"`
	CurrencyID                int     `json:"currency_id"`
	Description               string  `json:"description"`
	InterestRate              float64 `json:"interest_rate"`
	InterestCalculationPeriod int     `json:"interest_calculation_period"`
	InterestCalculationUnit   string  `json:"interest_calculation_unit"`
	CreatedAt                 string  `json:"created_at"`
	UpdatedAt                 string  `json:"updated_at"`
}

// ToSavingsProductResponse converts a SavingsProduct to a SavingsProductResponse
func ToSavingsProductResponse(product *SavingsProduct) *SavingsProductResponse {
	return &SavingsProductResponse{
		ID:                        product.ID,
		Name:                      product.Name,
		ProductCode:               product.ProductCode,
		CurrencyID:                product.CurrencyID,
		Description:               product.Description,
		InterestRate:              product.InterestRate,
		InterestCalculationPeriod: product.InterestCalculationPeriod,
		InterestCalculationUnit:   product.InterestCalculationUnit,
		CreatedAt:                 product.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:                 product.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
