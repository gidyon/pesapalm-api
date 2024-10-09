package user

import (
	"database/sql"
	"time"
)

const defaultUsersTable = "pesapalm_accounts"

var usersTable = ""

// User contains profile information stored in the database
type User struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement"`
	CreatorId     uint64         `gorm:"type:bigint"`
	Phone         sql.NullString `gorm:"type:varchar(15);index"`
	Email         sql.NullString `gorm:"type:varchar(50);index"`
	Names         string         `gorm:"type:varchar(50);not null"`
	BirthDate     sql.NullTime   `gorm:"type:datetime(6);"`
	Gender        string         `gorm:"type:varchar(20);"`
	ProfileURL    sql.NullString `gorm:"type:text"`
	Country       sql.NullString `gorm:"type:varchar(50)"`
	CountryCode   sql.NullString `gorm:"type:varchar(10)"`
	GroupId       sql.NullInt64  `gorm:"type:int(11);index"`
	Password      string         `gorm:"type:text"`
	GeneralData   []byte         `gorm:"type:json"`
	PrimaryGroup  string         `gorm:"type:varchar(50);index;not null"`
	AccountStatus string         `gorm:"type:enum('INVITED','BLOCKED','ACTIVE', 'INACTIVE','CREATED','DELETED');index;not null;default:'ACTIVE'"`
	LastLoginIp   sql.NullString `gorm:"type:varchar(20)"`
	LastLogin     sql.NullTime   `gorm:"type:datetime(6)"`
	UpdatedAt     time.Time      `gorm:"type:datetime(6);autoUpdateTime;index"`
	CreatedAt     time.Time      `gorm:"type:datetime(6);autoCreateTime;->;<-:create;index;not null"`
}

// TableName is the name of the tables
func (u *User) TableName() string {
	if usersTable != "" {
		return usersTable
	}
	return defaultUsersTable
}
