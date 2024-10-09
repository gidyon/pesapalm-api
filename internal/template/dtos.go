package template

type AdminGroupDTO struct {
	Name        string `json:"name" binding:"required"`
	ColumnRoles string `json:"column_roles"`
	Level       int    `json:"level" binding:"required"`
}

type BranchDTO struct {
	Name         string `json:"name" binding:"required"`
	Status       int    `json:"status" binding:"required"`
	Notes        string `json:"notes"`
	Msisdn       string `json:"msisdn"`
	EmailAddress string `json:"email_address" binding:"required"`
}

type CurrencyDTO struct {
	CurrencyCode string `json:"currency_code" binding:"required"`
	CurrencyName string `json:"currency_name" binding:"required"`
	StatusID     int    `json:"status_id" binding:"required"`
}

type LanguageDTO struct {
	LanguageCode string `json:"language_code" binding:"required"`
	LanguageName string `json:"language_name" binding:"required"`
	StatusID     int    `json:"status_id" binding:"required"`
}
