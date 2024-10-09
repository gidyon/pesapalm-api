package loans

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LoanController structure
type LoanController struct {
	*Options
}

// CreateLoanAccount creates a new loan account
func (ctrl *LoanController) CreateLoanAccount(c *gin.Context) {
	var dto CreateLoanAccountDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account := LoanAccount{
		LoanID:                dto.LoanID,
		CustomerID:            dto.CustomerID,
		LoanProductID:         dto.LoanProductID,
		CurrencyID:            dto.CurrencyID,
		LoanAmount:            dto.LoanAmount,
		RepaymentInstallments: dto.RepaymentInstallments,
		RepaymentPeriod:       dto.RepaymentPeriod,
		RepaymentPeriodUnit:   dto.RepaymentPeriodUnit,
		StatusID:              dto.StatusID,
	}

	if result := ctrl.DB.WithContext(c.Request.Context()).Create(&account); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

const selectFields = "loan_account.*, loan_product.id AS loan_product_id, loan_product.name as loan_product_name, customer.id AS customer_id, customer.first_name as customer_first_name, customer.last_name as customer_last_name, customer.middle_name as customer_middle_name"

// GetLoanAccount retrieves a loan account by ID with loan product and customer details
func (ctrl *LoanController) GetLoanAccount(c *gin.Context) {
	id := c.Param("id")
	var account LoanAccountRead

	// Manually joining loan_product and customer tables
	if result := ctrl.DB.WithContext(c.Request.Context()).Table("loan_account").
		Select(selectFields).
		Joins("LEFT JOIN loan_product ON loan_product.id = loan_account.loan_product_id").
		Joins("LEFT JOIN customer ON customer.id = loan_account.customer_id").
		Where("loan_account.id = ?", id).
		First(&account); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Loan account not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, ToLoanAccountResponse(&account))
}

// ListLoanAccounts retrieves a list of loan accounts with loan product and customer details
func (ctrl *LoanController) ListLoanAccounts(c *gin.Context) {
	var (
		queryParams = c.Request.URL.Query()
		pageSize, _ = strconv.Atoi(queryParams.Get("pageSize"))
		pageToken   = queryParams.Get("pageToken")
		searchTerm  = queryParams.Get("search")
		status      = queryParams.Get("status")
	)

	if pageSize <= 0 {
		pageSize = 10 // Default page size
	}

	var lastID int
	if pageToken != "" {
		bs, err := base64.StdEncoding.DecodeString(pageToken)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page token"})
			return
		}
		lastID, err = strconv.Atoi(string(bs))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page token"})
			return
		}
	}

	// Build the base query with joins
	db := ctrl.DB.WithContext(c.Request.Context()).Table("loan_account").
		Select("loan_account.id, loan_account.loan_id, loan_product.name AS product_name, customer.name AS customer_name").
		Joins("LEFT JOIN loan_product ON loan_product.id = loan_account.loan_product_id").
		Joins("LEFT JOIN customer ON customer.id = loan_account.customer_id").
		Order("loan_account.id DESC").
		Limit(pageSize + 1) // Fetch one extra record to detect next page

	// Apply filters before executing the query
	if lastID > 0 {
		db = db.Where("loan_account.id < ?", lastID)
	}
	if searchTerm != "" {
		db = db.Where("loan_account.loan_id LIKE ?", searchTerm+"%")
	}
	if status != "" {
		if statusID, err := strconv.Atoi(status); err == nil {
			db = db.Where("loan_account.status_id = ?", statusID)
		}
	}

	// Count matching records only for the first page
	var collectionCount int64
	if pageToken == "" {
		if err := db.Count(&collectionCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count loan accounts"})
			return
		}
	}

	// Fetch the loan accounts
	var loanAccounts []LoanAccountRead
	if err := db.Find(&loanAccounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve loan accounts"})
		return
	}

	// Prepare the response list
	responseAccounts := make([]*LoanAccountResponse, 0, len(loanAccounts))
	for i, account := range loanAccounts {
		// Stop at the pageSize limit
		if i == pageSize {
			break
		}
		responseAccounts = append(responseAccounts, ToLoanAccountResponse(&account))
	}

	// Generate the next page token if more records exist
	var nextPageToken string
	if len(loanAccounts) > pageSize {
		lastID := loanAccounts[pageSize-1].LoanAccount.ID
		nextPageToken = base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(lastID)))
	}

	// Return the response
	c.JSON(http.StatusOK, gin.H{
		"loan_accounts":   responseAccounts,
		"next_page_token": nextPageToken,
		"collectionCount": collectionCount,
	})
}

// GetLoanSchedule retrieves the repayment schedule for a given loan
func (ctrl *LoanController) GetLoanSchedule(c *gin.Context) {
	loanID := c.Param("loan_id")
	var schedule []LoanSchedule

	if result := ctrl.DB.WithContext(c.Request.Context()).Where("loan_id = ?", loanID).Find(&schedule); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// GetLoanEligibility retrieves loan eligibility for a customer
func (ctrl *LoanController) GetLoanEligibility(c *gin.Context) {
	customerID := c.Param("customer_id")
	var eligibility []LoanEligibility

	if result := ctrl.DB.WithContext(c.Request.Context()).Where("customer_id = ?", customerID).Find(&eligibility); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Loan eligibility not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, eligibility)
}
