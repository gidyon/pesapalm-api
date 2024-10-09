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
