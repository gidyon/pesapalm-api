package customer

import (
	"database/sql"
	"time"
)

// CreateCustomerDTO defines the JSON structure for creating a customer
type CreateCustomerDTO struct {
	CustomerID       string `json:"customer_id"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	MiddleName       string `json:"middle_name"`
	EmailAddress     string `json:"email_address"`
	Gender           string `json:"gender"`
	MSISDN1          string `json:"msisdn1"`
	MSISDN2          string `json:"msisdn2"`
	WorkPlaceAddress string `json:"work_place_address"`
	HomeAddress      string `json:"home_address"`
	Role             string `json:"role"`
	Notes            string `json:"notes"`
	StatusID         int    `json:"status_id"`
	ProfilePicID     int    `json:"profile_pic_id"`
	SignatureKeyID   int    `json:"signature_key_id"`
	LanguageID       int    `json:"language_id"`
	BranchID         int    `json:"branch_id"`
	CreatedBy        int    `json:"created_by"`
}

// UpdateCustomerDTO defines the JSON structure for updating a customer
type UpdateCustomerDTO struct {
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	MiddleName       string `json:"middle_name"`
	EmailAddress     string `json:"email_address"`
	Gender           string `json:"gender"`
	MSISDN1          string `json:"msisdn1"`
	MSISDN2          string `json:"msisdn2"`
	WorkPlaceAddress string `json:"work_place_address"`
	HomeAddress      string `json:"home_address"`
	Role             string `json:"role"`
	Notes            string `json:"notes"`
	StatusID         int    `json:"status_id"`
	ProfilePicID     int    `json:"profile_pic_id"`
	SignatureKeyID   int    `json:"signature_key_id"`
	LanguageID       int    `json:"language_id"`
	BranchID         int    `json:"branch_id"`
	ApprovedDate     string `json:"approved_date,omitempty"`
	ActivationDate   string `json:"activation_date,omitempty"`
	ClosedDate       string `json:"closed_date,omitempty"`
}

// CustomerResponse defines the structure of the customer data returned in the response
type CustomerResponse struct {
	ID               uint   `json:"id"`
	CustomerID       string `json:"customer_id"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	MiddleName       string `json:"middle_name"`
	EmailAddress     string `json:"email_address"`
	Gender           string `json:"gender"`
	MSISDN1          string `json:"msisdn1"`
	MSISDN2          string `json:"msisdn2"`
	WorkPlaceAddress string `json:"work_place_address"`
	HomeAddress      string `json:"home_address"`
	Role             string `json:"role"`
	Notes            string `json:"notes"`
	StatusID         int    `json:"status_id"`
	ProfilePicID     int    `json:"profile_pic_id"`
	SignatureKeyID   int    `json:"signature_key_id"`
	LanguageID       int    `json:"language_id"`
	BranchID         int    `json:"branch_id"`
	CreatedBy        int    `json:"created_by"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// ToCustomerResponse converts a Customer model to CustomerResponse
func ToCustomerResponse(customer *Customer) *CustomerResponse {
	return &CustomerResponse{
		ID:               customer.ID,
		CustomerID:       customer.CustomerID,
		FirstName:        customer.FirstName,
		LastName:         customer.LastName,
		MiddleName:       customer.MiddleName,
		EmailAddress:     customer.EmailAddress,
		Gender:           customer.Gender,
		MSISDN1:          nullStringToString(customer.MSISDN1),
		MSISDN2:          nullStringToString(customer.MSISDN2),
		WorkPlaceAddress: nullStringToString(customer.WorkPlaceAddress),
		HomeAddress:      nullStringToString(customer.HomeAddress),
		Role:             nullStringToString(customer.Role),
		Notes:            nullStringToString(customer.Notes),
		StatusID:         customer.StatusID,
		ProfilePicID:     customer.ProfilePicID,
		SignatureKeyID:   customer.SignatureKeyID,
		LanguageID:       customer.LanguageID,
		BranchID:         customer.BranchID,
		CreatedBy:        customer.CreatedBy,
		CreatedAt:        customer.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        customer.UpdatedAt.Format(time.RFC3339),
	}
}

// Helper function to convert sql.NullString to a normal string
func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
