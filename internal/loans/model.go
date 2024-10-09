package loans

import (
	"database/sql"
	"time"
)

// LoanAccount defines the GORM model for the loan_account table
type LoanAccount struct {
	ID                     uint         `gorm:"primaryKey"`
	LoanID                 string       `gorm:"unique;size:36"`
	CustomerID             string       `gorm:"size:36"`
	LoanProductID          int          `gorm:"index"`
	SavingsAccountID       int          `gorm:"index;default:null"`
	CurrencyID             int          `gorm:"type:TINYINT(1);default:1"`
	CurrencyCode           string       `gorm:"size:10;default:USD"`
	RepaymentInstallments  int          `gorm:"not null"`
	RepaymentPeriod        int          `gorm:"not null"`
	RepaymentPeriodUnit    string       `gorm:"size:20;default:DAY"`
	LoanAmount             float64      `gorm:"type:double(20,2);not null"`
	LoanBalance            float64      `gorm:"type:double(20,2);not null"`
	AmountPaid             float64      `gorm:"type:double(20,2);default:0.00"`
	OutstandingPrinciple   float64      `gorm:"type:double(20,2);default:0.00"`
	OutstandingSetupFees   float64      `gorm:"type:double(20,2);default:0.00"`
	OutstandingInterest    float64      `gorm:"type:double(20,2);default:0.00"`
	OutstandingPenaltyFees float64      `gorm:"type:double(20,2);default:0.00"`
	InterestEarned         float64      `gorm:"type:double(20,2);default:0.00"`
	StatusID               int          `gorm:"default:0"` // 0 = PENDING, 1 = ACTIVE, 2 = PAID, 3 = ERRORED
	Defaulted              int          `gorm:"default:0"` // 0 = ACTIVE, 1 = DEFAULTED
	InterestCalculated     int          `gorm:"default:0"` // 0 = PENDING, 1 = CALCULATED
	DueDate                sql.NullTime `gorm:"type:datetime"`
	LastRepaymentDate      sql.NullTime `gorm:"type:datetime"`
	LastInterestCalcDate   sql.NullTime `gorm:"type:datetime"`
	CreatedAt              time.Time    `gorm:"autoCreateTime"`
	UpdatedAt              time.Time    `gorm:"autoUpdateTime"`
}

type LoanAccountRead struct {
	LoanAccount `gorm:"embedded;"`
	Customer    `gorm:"embedded;"`
	LoanProduct `gorm:"embedded;"`
}

type Customer struct {
	ID         uint   `gorm:"-:migration;column:customer_id" json:"customer_id,omitempty"`
	FirstName  string `gorm:"-:migration;column:customer_first_name" json:"first_name,omitempty"`
	MiddleName string `gorm:"-:migration;<-:false;column:customer_middle_name" json:"middle_name,omitempty"`
	LastName   string `gorm:"-:migration;<-:false;column:customer_last_name" json:"last_name,omitempty"`
}

type LoanProduct struct {
	ID          uint   `gorm:"-:migration;<-:false;column:loan_product_id" json:"loan_product_id,omitempty"`
	ProductName string `gorm:"-:migration;<-:false;column:loan_product_name" json:"product_name,omitempty"`
}

func (*LoanAccount) TableName() string {
	return "loan_account"
}

// LoanSchedule defines the GORM model for the loan_schedule table
type LoanSchedule struct {
	ID                                uint         `gorm:"primaryKey"`
	LoanID                            int          `gorm:"index;not null"`
	LoanAccountID                     string       `gorm:"size:36"`
	LoanProductID                     int          `gorm:"index"`
	CustomerID                        int          `gorm:"index"`
	CurrencyID                        int          `gorm:"type:TINYINT(1);default:1"`
	CurrencyCode                      string       `gorm:"size:10;default:USD"`
	InstallmentAmount                 float64      `gorm:"type:double(15,2);not null"`
	InstallmentBalance                float64      `gorm:"type:double(15,2);not null"`
	InstallmentAmountPaid             float64      `gorm:"type:double(15,2);default:0.00"`
	InstallmentOutstandingPrinciple   float64      `gorm:"type:double(15,2);default:0.00"`
	InstallmentOutstandingSetupFees   float64      `gorm:"type:double(15,2);default:0.00"`
	InstallmentOutstandingInterest    float64      `gorm:"type:double(15,2);default:0.00"`
	InstallmentOutstandingPenaltyFees float64      `gorm:"type:double(15,2);default:0.00"`
	InstallmentInterestEarned         float64      `gorm:"type:double(15,2);default:0.00"`
	StatusID                          int          `gorm:"default:0"` // 0 = PENDING, 1 = ACTIVE, 2 = PAID, 3 = ERRORED
	Defaulted                         int          `gorm:"default:0"` // 0 = ACTIVE, 1 = DEFAULTED
	InterestCalculated                int          `gorm:"default:0"` // 0 = PENDING, 1 = CALCULATED
	DueDate                           sql.NullTime `gorm:"type:datetime"`
	RepaymentDate                     sql.NullTime `gorm:"type:datetime"`
	InterestCalcDate                  sql.NullTime `gorm:"type:datetime"`
	CreatedAt                         time.Time    `gorm:"autoCreateTime"`
	UpdatedAt                         time.Time    `gorm:"autoUpdateTime"`
}

func (*LoanSchedule) TableName() string {
	return "loan_schedule"
}

// LoanEligibility defines the GORM model for the loan_eligibility table
type LoanEligibility struct {
	ID                 uint      `gorm:"primaryKey"`
	CustomerID         int       `gorm:"index;not null"`
	SavingsID          int       `gorm:"index;default:null"`
	CurrencyID         int       `gorm:"type:TINYINT(1);default:1"`
	CurrencyCode       string    `gorm:"size:10;default:USD"`
	LoanEligibleAmount float64   `gorm:"type:double(30,2);not null;default:0.00"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
}

func (*LoanEligibility) TableName() string {
	return "loan_eligibility"
}
