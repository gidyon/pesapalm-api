package savings_product

import (
	"time"
)

// SavingsProduct defines the GORM model for the savings_product table
type SavingsProduct struct {
	ID                        uint      `gorm:"primaryKey"`
	Name                      string    `gorm:"size:32"`
	ProductCode               string    `gorm:"size:32;unique"`
	CurrencyID                int       `gorm:"type:TINYINT(1);default:1"`
	Description               string    `gorm:"type:mediumtext"`
	InterestRate              float64   `gorm:"type:double(10,5);default:0.00000"`
	InterestCalculationPeriod int       `gorm:"type:int"`
	InterestCalculationUnit   string    `gorm:"size:10;default:DAY"`
	CreatedAt                 time.Time `gorm:"autoCreateTime"`
	UpdatedAt                 time.Time `gorm:"autoUpdateTime"`
}

func (*SavingsProduct) TableName() string {
	return "savings_product"
}
