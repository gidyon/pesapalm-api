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
		StatusID:              1,
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
		pageToken   = queryParams.Get("pageToken")
		searchTerm  = queryParams.Get("search")
		status      = queryParams.Get("status")
	)

	// Parse pageSize from query, default if invalid
	pageSize, _ := strconv.Atoi(queryParams.Get("pageSize"))
	switch {
	case pageSize <= 0:
		pageSize = 10
	case pageSize > 100:
		pageSize = 100
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
		Select(selectFields).
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

// Validate timeGroup
var validTimeGroups = map[string]string{
	"yearly":  "YEAR(created_at)",
	"monthly": "DATE_FORMAT(created_at, '%Y-%m')",
	"daily":   "DATE(created_at)",
	"hourly":  "DATE_FORMAT(created_at, '%Y-%m-%d %H:00:00')",
}

// GetStats retrieves aggregated loan statistics for graphing based on filters.
func (ctrl *LoanController) GetStats(c *gin.Context) {
	// Parse request parameters
	queryParams := c.Request.URL.Query()

	productID := queryParams.Get("product_id")
	currencyID := queryParams.Get("currency_id")
	financialPartner := queryParams.Get("financial_partner")
	countryID := queryParams.Get("country_id")
	timeFrom := queryParams.Get("from")
	timeTo := queryParams.Get("to")
	timeGroup := queryParams.Get("time_group") // e.g., "yearly", "monthly", "daily", "hourly"

	groupByField, isValid := validTimeGroups[timeGroup]
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time group"})
		return
	}

	// Base query
	db := ctrl.DB.WithContext(c.Request.Context()).Table("loan_account").
		Select(fmt.Sprintf(`
			%s AS time_group,
			COUNT(*) AS count,
			SUM(loan_amount) AS volume
		`, groupByField)).
		Group("time_group").
		Order("time_group ASC")

	// Apply filters
	if productID != "" {
		db = db.Where("loan_product_id = ?", productID)
	}
	if currencyID != "" {
		db = db.Where("currency_id = ?", currencyID)
	}
	if financialPartner != "" {
		db = db.Where("financial_partner = ?", financialPartner)
	}
	if countryID != "" {
		db = db.Where("country_id = ?", countryID)
	}
	if timeFrom != "" && timeTo != "" {
		db = db.Where("created_at BETWEEN ? AND ?", timeFrom, timeTo)
	} else if timeFrom != "" {
		db = db.Where("created_at >= ?", timeFrom)
	} else if timeTo != "" {
		db = db.Where("created_at <= ?", timeTo)
	}

	// Execute the query and retrieve results
	var stats []struct {
		TimeGroup string  `json:"time_group"`
		Count     int     `json:"count"`
		Volume    float64 `json:"volume"`
	}
	if err := db.Find(&stats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve statistics"})
		return
	}

	// Prepare the response
	response := map[string]interface{}{
		"data":  stats,
		"count": len(stats),
		"volume": func() float64 {
			total := 0.0
			for _, stat := range stats {
				total += stat.Volume
			}
			return total
		}(),
	}

	// Send the response
	c.JSON(http.StatusOK, response)
}
