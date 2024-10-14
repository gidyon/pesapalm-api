package savings

import (
	"database/sql"
	"time"
)

// SavingsAccount defines the GORM model for the savings_account table
type SavingsAccount struct {
	ID                          uint         `gorm:"primaryKey"`
	SavingsID                   string       `gorm:"type:varchar(50)"`
	CustomerID                  int          `gorm:"index;type:INT(11)"`
	ProductID                   int          `gorm:"index;type:INT(11)"`
	CurrencyID                  int          `gorm:"index;type:TINYINT(1)"`
	CurrencyCode                string       `gorm:"type:varchar(10)"`
	Balance                     float64      `gorm:"type:double(20,2);default:0.00"`
	StatusID                    int          `gorm:"index;type:TINYINT(1);default:1;comments:'1 - Default, 2 - Approved, 3 - Activated, 4 - Closed'"`
	DateClosed                  sql.NullTime `gorm:"type:DATETIME"`
	DateApproved                sql.NullTime `gorm:"type:DATETIME"`
	DateActivated               sql.NullTime `gorm:"type:DATETIME"`
	LastInterestCalculationDate sql.NullTime `gorm:"type:DATETIME"`
	MaturityDate                sql.NullTime `gorm:"type:DATETIME"`
	MaximumWithdrawableAmount   *float64     `gorm:"type:double(20,2)"`
	FeesDue                     float64      `gorm:"type:double(20,2);default:0.00"`
	LockedBalance               float64      `gorm:"type:double(20,2);default:0.00"`
	DateLocked                  sql.NullTime `gorm:"type:DATETIME"`
	CreatedAt                   time.Time    `gorm:"type:datetime(6);autoCreateTime;->;<-:create;index;not null"`
	UpdatedAt                   time.Time    `gorm:"type:datetime(6);autoUpdateTime"`
}

func (*SavingsAccount) TableName() string {
	return "savings_account"
}

type SavingsAccountRead struct {
	Customer       `gorm:"embedded;"`
	SavingProduct  `gorm:"embedded;"`
	SavingsAccount `gorm:"embedded;"`
}

type Customer struct {
	ID         uint   `gorm:"-:migration;column:customer_id" json:"customer_id,omitempty"`
	FirstName  string `gorm:"-:migration;column:customer_first_name" json:"first_name,omitempty"`
	MiddleName string `gorm:"-:migration;<-:false;column:customer_middle_name" json:"middle_name,omitempty"`
	LastName   string `gorm:"-:migration;<-:false;column:customer_last_name" json:"last_name,omitempty"`
}

type SavingProduct struct {
	ID          uint   `gorm:"-:migration;<-:false;column:saving_product_id" json:"saving_product_id,omitempty"`
	ProductName string `gorm:"-:migration;<-:false;column:saving_product_name" json:"saving_product_name,omitempty"`
}
