package loans_product

import (
	"time"
)

// LoanProduct defines the GORM model for the loan_product table
type LoanProduct struct {
	ID                        uint      `gorm:"primaryKey"`
	Name                      string    `gorm:"size:32"`
	ProductID                 string    `gorm:"size:32;unique"`
	CurrencyID                int       `gorm:"type:TINYINT(1);default:1"`
	Description               string    `gorm:"type:mediumtext"`
	MaxLoanAmount             float64   `gorm:"type:double(30,2);default:0.00"`
	MaxInstallments           int       `gorm:"type:int"`
	MinInstallments           int       `gorm:"type:int"`
	InterestRate              float64   `gorm:"type:double(30,2);default:0.00"`
	InterestCalculationPeriod int       `gorm:"type:int"`
	InterestCalculationUnit   string    `gorm:"size:10;default:DAY"`
	RepaymentPeriod           int       `gorm:"type:int"`
	RepaymentPeriodUnit       string    `gorm:"size:10;default:DAY"`
	CreatedAt                 time.Time `gorm:"autoCreateTime"`
	UpdatedAt                 time.Time `gorm:"autoUpdateTime"`
}

func (*LoanProduct) TableName() string {
	return "loan_product"
}
