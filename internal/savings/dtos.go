package savings

import (
	"database/sql"
	"time"
)

// CreateSavingsAccountDTO defines the JSON structure for creating a new savings account
type CreateSavingsAccountDTO struct {
	SavingsID    string  `json:"savings_id"`
	CustomerID   int     `json:"customer_id"`
	ProductID    int     `json:"product_id"`
	CurrencyID   int     `json:"currency_id"`
	CurrencyCode string  `json:"currency_code"`
	Balance      float64 `json:"balance"`
	StatusID     int     `json:"status_id"`
}

// UpdateSavingsAccountDTO defines the JSON structure for updating a savings account
type UpdateSavingsAccountDTO struct {
	Balance  float64 `json:"balance"`
	StatusID int     `json:"status_id"`
}

// SavingsAccountResponse defines the structure of the savings account data returned in the response
type SavingsAccountResponse struct {
	ID                          uint     `json:"id"`
	SavingsID                   string   `json:"savings_id"`
	CustomerID                  int      `json:"customer_id"`
	ProductID                   int      `json:"product_id"`
	CurrencyID                  int      `json:"currency_id"`
	CurrencyCode                string   `json:"currency_code"`
	Balance                     float64  `json:"balance"`
	StatusID                    int      `json:"status_id"`
	DateClosed                  *string  `json:"date_closed"`
	DateApproved                *string  `json:"date_approved"`
	DateActivated               *string  `json:"date_activated"`
	LastInterestCalculationDate *string  `json:"last_interest_calculation_date"`
	MaturityDate                *string  `json:"maturity_date"`
	MaximumWithdrawableAmount   *float64 `json:"maximum_withdrawable_amount"`
	FeesDue                     float64  `json:"fees_due"`
	LockedBalance               float64  `json:"locked_balance"`
	DateLocked                  *string  `json:"date_locked"`
	CreatedAt                   string   `json:"created_at"`
	UpdatedAt                   string   `json:"updated_at"`
}

// ToSavingsAccountResponse converts a SavingsAccount to a SavingsAccountResponse
func ToSavingsAccountResponse(account *SavingsAccount) *SavingsAccountResponse {
	return &SavingsAccountResponse{
		ID:                          account.ID,
		SavingsID:                   account.SavingsID,
		CustomerID:                  account.CustomerID,
		ProductID:                   account.ProductID,
		CurrencyID:                  account.CurrencyID,
		CurrencyCode:                account.CurrencyCode,
		Balance:                     account.Balance,
		StatusID:                    account.StatusID,
		DateClosed:                  formatNullableTime(account.DateClosed),
		DateApproved:                formatNullableTime(account.DateApproved),
		DateActivated:               formatNullableTime(account.DateActivated),
		LastInterestCalculationDate: formatNullableTime(account.LastInterestCalculationDate),
		MaturityDate:                formatNullableTime(account.MaturityDate),
		MaximumWithdrawableAmount:   account.MaximumWithdrawableAmount,
		FeesDue:                     account.FeesDue,
		LockedBalance:               account.LockedBalance,
		DateLocked:                  formatNullableTime(account.DateLocked),
		CreatedAt:                   account.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:                   account.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func formatNullableTime(t sql.NullTime) *string {
	if t.Valid {
		formatted := t.Time.UTC().Format(time.RFC3339)
		return &formatted
	}
	return nil
}

// UpdateSavingsAccountStatusDTO defines the JSON structure for updating a savings account status
type UpdateSavingsAccountStatusDTO struct {
	Action string `json:"action" binding:"required,oneof=approve activate close"`
}
