package loans

import "time"

// CreateLoanAccountDTO defines the JSON structure for creating a loan account
type CreateLoanAccountDTO struct {
	LoanID                string  `json:"loan_id"`
	CustomerID            string  `json:"customer_id"`
	LoanProductID         int     `json:"loan_product_id"`
	CurrencyID            int     `json:"currency_id"`
	LoanAmount            float64 `json:"loan_amount"`
	RepaymentInstallments int     `json:"repayment_installments"`
	RepaymentPeriod       int     `json:"repayment_period"`
	RepaymentPeriodUnit   string  `json:"repayment_period_unit"`
	StatusID              int     `json:"status_id"`
}

// LoanAccountResponse defines the structure of the loan account data returned in the response
type LoanAccountResponse struct {
	ID                     uint        `json:"id,omitempty"`
	LoanID                 string      `json:"loan_id,omitempty"`
	CustomerID             string      `json:"customer_id,omitempty"`
	LoanProductID          int         `json:"loan_product_id,omitempty"`
	SavingsAccountID       int         `json:"savings_account_id,omitempty"`
	CurrencyID             int         `json:"currency_id,omitempty"`
	CurrencyCode           string      `json:"currency_code,omitempty"`
	RepaymentInstallments  int         `json:"repayment_installments,omitempty"`
	RepaymentPeriod        int         `json:"repayment_period,omitempty"`
	RepaymentPeriodUnit    string      `json:"repayment_period_unit,omitempty"`
	LoanAmount             float64     `json:"loan_amount,omitempty"`
	LoanBalance            float64     `json:"loan_balance,omitempty"`
	AmountPaid             float64     `json:"amount_paid,omitempty"`
	OutstandingPrinciple   float64     `json:"outstanding_principle,omitempty"`
	OutstandingSetupFees   float64     `json:"outstanding_setup_fees,omitempty"`
	OutstandingInterest    float64     `json:"outstanding_interest,omitempty"`
	OutstandingPenaltyFees float64     `json:"outstanding_penalty_fees,omitempty"`
	InterestEarned         float64     `json:"interest_earned,omitempty"`
	StatusID               int         `json:"status_id,omitempty"`
	Defaulted              int         `json:"defaulted,omitempty"`
	InterestCalculated     int         `json:"interest_calculated,omitempty"`
	DueDate                *string     `json:"due_date,omitempty"`
	LastRepaymentDate      *string     `json:"last_repayment_date,omitempty"`
	LastInterestCalcDate   *string     `json:"last_interest_calc_date,omitempty"`
	Customer               Customer    `json:"customer,omitempty"`
	LoanProduct            LoanProduct `json:"loan_product,omitempty"`
	CreatedAt              string      `json:"created_at,omitempty"`
	UpdatedAt              string      `json:"updated_at,omitempty"`
}

// ToLoanAccountResponse converts a LoanAccount to a LoanAccountResponse
func ToLoanAccountResponse(account *LoanAccountRead) *LoanAccountResponse {
	return &LoanAccountResponse{
		ID:                     account.LoanAccount.ID,
		LoanID:                 account.LoanID,
		CustomerID:             account.CustomerID,
		LoanProductID:          account.LoanProductID,
		SavingsAccountID:       account.SavingsAccountID,
		CurrencyID:             account.CurrencyID,
		CurrencyCode:           account.CurrencyCode,
		RepaymentInstallments:  account.RepaymentInstallments,
		RepaymentPeriod:        account.RepaymentPeriod,
		RepaymentPeriodUnit:    account.RepaymentPeriodUnit,
		LoanAmount:             account.LoanAmount,
		LoanBalance:            account.LoanBalance,
		AmountPaid:             account.AmountPaid,
		OutstandingPrinciple:   account.OutstandingPrinciple,
		OutstandingSetupFees:   account.OutstandingSetupFees,
		OutstandingInterest:    account.OutstandingInterest,
		OutstandingPenaltyFees: account.OutstandingPenaltyFees,
		InterestEarned:         account.InterestEarned,
		StatusID:               account.StatusID,
		Defaulted:              account.Defaulted,
		InterestCalculated:     account.InterestCalculated,
		DueDate:                formatNullableTime(account.DueDate.Time),
		LastRepaymentDate:      formatNullableTime(account.LastRepaymentDate.Time),
		LastInterestCalcDate:   formatNullableTime(account.LastInterestCalcDate.Time),
		CreatedAt:              account.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:              account.UpdatedAt.UTC().Format(time.RFC3339),
		Customer:               account.Customer,
		LoanProduct:            account.LoanProduct,
	}
}

// LoanScheduleResponse defines the structure of the loan schedule data returned in the response
type LoanScheduleResponse struct {
	ID                                uint    `json:"id"`
	LoanID                            int     `json:"loan_id"`
	LoanAccountID                     string  `json:"loan_account_id"`
	LoanProductID                     int     `json:"loan_product_id"`
	CustomerID                        int     `json:"customer_id"`
	CurrencyID                        int     `json:"currency_id"`
	CurrencyCode                      string  `json:"currency_code"`
	InstallmentAmount                 float64 `json:"installment_amount"`
	InstallmentBalance                float64 `json:"installment_balance"`
	InstallmentAmountPaid             float64 `json:"installment_amount_paid"`
	InstallmentOutstandingPrinciple   float64 `json:"installment_outstanding_principle"`
	InstallmentOutstandingSetupFees   float64 `json:"installment_outstanding_setup_fees"`
	InstallmentOutstandingInterest    float64 `json:"installment_outstanding_interest"`
	InstallmentOutstandingPenaltyFees float64 `json:"installment_outstanding_penalty_fees"`
	InstallmentInterestEarned         float64 `json:"installment_interest_earned"`
	StatusID                          int     `json:"status_id"`
	Defaulted                         int     `json:"defaulted"`
	InterestCalculated                int     `json:"interest_calculated"`
	DueDate                           *string `json:"due_date,omitempty"`
	RepaymentDate                     *string `json:"repayment_date,omitempty"`
	InterestCalcDate                  *string `json:"interest_calc_date,omitempty"`
	CreatedAt                         string  `json:"created_at"`
	UpdatedAt                         string  `json:"updated_at"`
}

// ToLoanScheduleResponse converts a LoanSchedule to LoanScheduleResponse
func ToLoanScheduleResponse(schedule *LoanSchedule) *LoanScheduleResponse {
	return &LoanScheduleResponse{
		ID:                                schedule.ID,
		LoanID:                            schedule.LoanID,
		LoanAccountID:                     schedule.LoanAccountID,
		LoanProductID:                     schedule.LoanProductID,
		CustomerID:                        schedule.CustomerID,
		CurrencyID:                        schedule.CurrencyID,
		CurrencyCode:                      schedule.CurrencyCode,
		InstallmentAmount:                 schedule.InstallmentAmount,
		InstallmentBalance:                schedule.InstallmentBalance,
		InstallmentAmountPaid:             schedule.InstallmentAmountPaid,
		InstallmentOutstandingPrinciple:   schedule.InstallmentOutstandingPrinciple,
		InstallmentOutstandingSetupFees:   schedule.InstallmentOutstandingSetupFees,
		InstallmentOutstandingInterest:    schedule.InstallmentOutstandingInterest,
		InstallmentOutstandingPenaltyFees: schedule.InstallmentOutstandingPenaltyFees,
		InstallmentInterestEarned:         schedule.InstallmentInterestEarned,
		StatusID:                          schedule.StatusID,
		Defaulted:                         schedule.Defaulted,
		InterestCalculated:                schedule.InterestCalculated,
		DueDate:                           formatNullableTime(schedule.DueDate.Time),
		RepaymentDate:                     formatNullableTime(schedule.RepaymentDate.Time),
		InterestCalcDate:                  formatNullableTime(schedule.InterestCalcDate.Time),
		CreatedAt:                         schedule.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:                         schedule.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

// LoanEligibilityResponse defines the structure of the loan eligibility data returned in the response
type LoanEligibilityResponse struct {
	ID                 uint    `json:"id"`
	CustomerID         int     `json:"customer_id"`
	SavingsID          int     `json:"savings_id,omitempty"`
	CurrencyID         int     `json:"currency_id"`
	CurrencyCode       string  `json:"currency_code"`
	LoanEligibleAmount float64 `json:"loan_eligible_amount"`
	UpdatedAt          string  `json:"updated_at"`
}

// ToLoanEligibilityResponse converts a LoanEligibility model to LoanEligibilityResponse
func ToLoanEligibilityResponse(eligibility *LoanEligibility) *LoanEligibilityResponse {
	return &LoanEligibilityResponse{
		ID:                 eligibility.ID,
		CustomerID:         eligibility.CustomerID,
		SavingsID:          eligibility.SavingsID,
		CurrencyID:         eligibility.CurrencyID,
		CurrencyCode:       eligibility.CurrencyCode,
		LoanEligibleAmount: eligibility.LoanEligibleAmount,
		UpdatedAt:          eligibility.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

// Helper function to format nullable time fields
func formatNullableTime(t time.Time) *string {
	if !t.IsZero() {
		formatted := t.UTC().Format(time.RFC3339)
		return &formatted
	}
	return nil
}
