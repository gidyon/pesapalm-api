package template

import "time"

// AdminGroup model
type AdminGroup struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"size:256;not null;column:name"`
	CreatedAt   time.Time `gorm:"not null;column:createdAt"`
	UpdatedAt   time.Time `gorm:"not null;column:updatedAt"`
	ColumnRoles string    `gorm:"size:512;column:column_roles"`
	Level       int       `gorm:"not null;column:level"`
}

func (*AdminGroup) TableName() string {
	return "admin_groups"
}

// Branch model
type Branch struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	Name         string    `gorm:"size:256;not null;column:name"`
	Status       int       `gorm:"not null;column:status"`
	Notes        string    `gorm:"size:512;column:notes"`
	Msisdn       string    `gorm:"size:20;column:msisdn"`
	EmailAddress string    `gorm:"size:256;not null;column:email_address"`
	DateCreated  time.Time `gorm:"not null;column:date_created"`
	DateModified time.Time `gorm:"column:date_modified"`
}

func (*Branch) TableName() string {
	return "branch"
}

// Currency model
type Currency struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	CurrencyCode string    `gorm:"size:10;not null;column:currency_code"`
	CurrencyName string    `gorm:"size:256;not null;column:currency_name"`
	StatusID     int       `gorm:"not null;column:status_id"`
	CreatedAt    time.Time `gorm:"not null;column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (*Currency) TableName() string {
	return "currency"
}

// Language model
type Language struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	LanguageCode string    `gorm:"size:10;not null;column:language_code"`
	LanguageName string    `gorm:"size:256;not null;column:language_name"`
	StatusID     int       `gorm:"not null;column:status_id"`
	CreatedAt    time.Time `gorm:"not null;column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (*Language) TableName() string {
	return "language"
}
