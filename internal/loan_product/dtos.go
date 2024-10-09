package loans_product

import "time"

// CreateLoanProductDTO defines the JSON structure for creating a new loan product
type CreateLoanProductDTO struct {
	Name                      string  `json:"name"`
	ProductID                 string  `json:"product_id"`
	CurrencyID                int     `json:"currency_id"`
	Description               string  `json:"description"`
	MaxLoanAmount             float64 `json:"max_loan_amount"`
	MaxInstallments           int     `json:"max_installments"`
	MinInstallments           int     `json:"min_installments"`
	InterestRate              float64 `json:"interest_rate"`
	InterestCalculationPeriod int     `json:"interest_calculation_period"`
	InterestCalculationUnit   string  `json:"interest_calculation_unit"`
	RepaymentPeriod           int     `json:"repayment_period"`
	RepaymentPeriodUnit       string  `json:"repayment_period_unit"`
}

// UpdateLoanProductDTO defines the JSON structure for updating a loan product
type UpdateLoanProductDTO struct {
	Name                      string  `json:"name"`
	Description               string  `json:"description"`
	MaxLoanAmount             float64 `json:"max_loan_amount"`
	MaxInstallments           int     `json:"max_installments"`
	MinInstallments           int     `json:"min_installments"`
	InterestRate              float64 `json:"interest_rate"`
	InterestCalculationPeriod int     `json:"interest_calculation_period"`
	InterestCalculationUnit   string  `json:"interest_calculation_unit"`
	RepaymentPeriod           int     `json:"repayment_period"`
	RepaymentPeriodUnit       string  `json:"repayment_period_unit"`
}

// LoanProductResponse defines the structure of the loan product data returned in the response
type LoanProductResponse struct {
	ID                        uint    `json:"id"`
	Name                      string  `json:"name"`
	ProductID                 string  `json:"product_id"`
	CurrencyID                int     `json:"currency_id"`
	Description               string  `json:"description"`
	MaxLoanAmount             float64 `json:"max_loan_amount"`
	MaxInstallments           int     `json:"max_installments"`
	MinInstallments           int     `json:"min_installments"`
	InterestRate              float64 `json:"interest_rate"`
	InterestCalculationPeriod int     `json:"interest_calculation_period"`
	InterestCalculationUnit   string  `json:"interest_calculation_unit"`
	RepaymentPeriod           int     `json:"repayment_period"`
	RepaymentPeriodUnit       string  `json:"repayment_period_unit"`
	CreatedAt                 string  `json:"created_at"`
	UpdatedAt                 string  `json:"updated_at"`
}

// ToLoanProductResponse converts a LoanProduct to a LoanProductResponse
func ToLoanProductResponse(product *LoanProduct) *LoanProductResponse {
	return &LoanProductResponse{
		ID:                        product.ID,
		Name:                      product.Name,
		ProductID:                 product.ProductID,
		CurrencyID:                product.CurrencyID,
		Description:               product.Description,
		MaxLoanAmount:             product.MaxLoanAmount,
		MaxInstallments:           product.MaxInstallments,
		MinInstallments:           product.MinInstallments,
		InterestRate:              product.InterestRate,
		InterestCalculationPeriod: product.InterestCalculationPeriod,
		InterestCalculationUnit:   product.InterestCalculationUnit,
		RepaymentPeriod:           product.RepaymentPeriod,
		RepaymentPeriodUnit:       product.RepaymentPeriodUnit,
		CreatedAt:                 product.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:                 product.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
