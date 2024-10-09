package customer

import (
	"database/sql"
	"time"
)

// Customer defines the GORM model for the customer table
type Customer struct {
	ID               uint           `gorm:"primaryKey"`
	CustomerID       string         `gorm:"unique;size:36"`
	FirstName        string         `gorm:"size:256;not null"`
	LastName         string         `gorm:"size:256;not null"`
	MiddleName       string         `gorm:"size:256"`
	EmailAddress     string         `gorm:"size:256;not null"`
	Gender           string         `gorm:"size:256;not null"`
	MSISDN1          sql.NullString `gorm:"size:256"`
	MSISDN2          sql.NullString `gorm:"size:256"`
	WorkPlaceAddress sql.NullString `gorm:"size:256"`
	HomeAddress      sql.NullString `gorm:"size:256"`
	Role             sql.NullString `gorm:"size:32"`
	Notes            sql.NullString `gorm:"type:mediumtext"`
	StatusID         int            `gorm:"default:1"` // 1 = ACTIVE, 2 = INACTIVE
	ProfilePicID     int            `gorm:"default:0"`
	SignatureKeyID   int            `gorm:"default:0"`
	LanguageID       int            `gorm:"default:1"` // 1 = FRENCH
	BranchID         int            `gorm:"default:1"`
	CreatedBy        int            `gorm:"default:1"`
	ApprovedDate     sql.NullTime   `gorm:"type:datetime"`
	ActivationDate   sql.NullTime   `gorm:"type:datetime"`
	ClosedDate       sql.NullTime   `gorm:"type:datetime"`
	CreatedAt        time.Time      `gorm:"autoCreateTime"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime"`
}

func (*Customer) TableName() string {
	return "customer"
}
