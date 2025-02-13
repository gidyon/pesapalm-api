package savings

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gidyon/pesapalm/internal/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/grpclog"
	"gorm.io/gorm"
)

type Options struct {
	DB           *gorm.DB
	Logger       grpclog.LoggerV2
	TokenManager auth.TokenInterface
	GinEngine    *gin.Engine
}

// Controller structure
type SavingsAccountController struct {
	*Options
}

const selectFields = "savings_account.*, savings_product.id AS saving_product_id, savings_product.name as saving_product_name, savings_product.product_code as saving_product_code, customer.id AS customer_id, customer.first_name as customer_first_name, customer.last_name as customer_last_name, customer.middle_name as customer_middle_name"

// CreateSavingsAccount creates a new savings account
func (ctrl *SavingsAccountController) CreateSavingsAccount(c *gin.Context) {
	var dto CreateSavingsAccountDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		ctrl.Logger.Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account := SavingsAccount{
		SavingsID:    dto.SavingsID,
		CustomerID:   dto.CustomerID,
		ProductID:    dto.ProductID,
		CurrencyID:   dto.CurrencyID,
		CurrencyCode: dto.CurrencyCode,
		Balance:      dto.Balance,
		StatusID:     dto.StatusID,
	}

	if result := ctrl.DB.WithContext(c.Request.Context()).Create(&account); result.Error != nil {
		ctrl.Logger.Errorf("Failed to create savings account: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	ctrl.Logger.Infof("Successfully created savings account with ID: %d", account.ID)
	c.JSON(http.StatusOK, account)
}

// GetSavingsAccount retrieves a savings account by ID, following the GetUser pattern
func (ctrl *SavingsAccountController) GetSavingsAccount(c *gin.Context) {
	id := c.Param("id") // ID from the URL path

	// Create a placeholder for the savings account record
	var db SavingsAccountRead

	// Fetch the savings account by ID from the database
	err := ctrl.DB.WithContext(c.Request.Context()).
		Joins("LEFT JOIN savings_product ON savings_product.id = savings_account.product_id").
		Joins("LEFT JOIN customer ON customer.id = savings_account.customer_id").
		Select(selectFields).First(&db, "savings_account.id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Savings account not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get savings account"})
		}
		return
	}

	// Return the savings account details as JSON
	c.JSON(http.StatusOK, ToSavingsAccountResponse(&db))
}

// UpdateSavingsAccount updates an existing savings account
func (ctrl *SavingsAccountController) UpdateSavingsAccount(c *gin.Context) {
	id := c.Param("id")
	var dto UpdateSavingsAccountDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var account SavingsAccount
	if result := ctrl.DB.WithContext(c.Request.Context()).First(&account, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Savings account not found"})
		return
	}

	// Update fields
	account.Balance = dto.Balance
	account.StatusID = dto.StatusID

	if result := ctrl.DB.Save(&account); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

// DeleteSavingsAccount deletes a savings account
func (ctrl *SavingsAccountController) DeleteSavingsAccount(c *gin.Context) {
	id := c.Param("id")
	if result := ctrl.DB.WithContext(c.Request.Context()).Delete(&SavingsAccount{}, id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Savings account deleted"})
}

const (
	defaultPageSize = 200
	maxPageSize     = 1000
)

// ListSavingsAccounts retrieves a list of savings accounts with related product and customer details
func (ctrl *SavingsAccountController) ListSavingsAccounts(c *gin.Context) {
	var (
		queryParams = c.Request.URL.Query()
		searchTerm  = queryParams.Get("search")
		status      = queryParams.Get("status")
		startDate   = queryParams.Get("startDate") // Start Date (timestamp)
		endDate     = queryParams.Get("endDate")   // End Date (timestamp)
	)

	// Parse pageSize from query, default if invalid
	pageSize, _ := strconv.Atoi(queryParams.Get("pageSize"))
	switch {
	case pageSize <= 0:
		pageSize = defaultPageSize
	case pageSize > defaultPageSize:
		pageSize = defaultPageSize
	}

	var lastID int

	// Get last id from page token
	pageToken := queryParams.Get("pageToken")
	if pageToken != "" {
		bs, err := base64.StdEncoding.DecodeString(pageToken)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "page token is incorrect"})
			return
		}
		lastID, err = strconv.Atoi(string(bs))
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "page token is incorrect"})
			return
		}
	}

	// Build the main query with joins for fetching data
	db := ctrl.DB.WithContext(c.Request.Context()).
		Table("savings_account").
		Select(selectFields).
		Joins("LEFT JOIN savings_product ON savings_product.id = savings_account.product_id").
		Joins("LEFT JOIN customer ON customer.id = savings_account.customer_id").
		Order("savings_account.id DESC").
		Limit(pageSize + 1)

	// Apply filters before executing the query
	if lastID > 0 {
		db = db.Where("savings_account.id < ?", lastID)
	}
	if searchTerm != "" {
		db = db.Where("savings_account.savings_id LIKE ?", searchTerm+"%")
	}
	if status != "" {
		if statusID, err := strconv.Atoi(status); err == nil {
			db = db.Where("savings_account.status_id = ?", statusID)
		}
	}

	// Parse and validate startDate and endDate
	if startDate != "" && endDate != "" {
		startDateInt, errStart := strconv.ParseInt(startDate, 10, 64)
		endDateInt, errEnd := strconv.ParseInt(endDate, 10, 64)

		if errStart != nil || errEnd != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid startDate or endDate"})
			return
		}

		// Convert timestamps to time.Time
		startTime := time.Unix(startDateInt, 0)
		endTime := time.Unix(endDateInt, 0)

		// Validate that startDate is earlier than endDate
		if startTime.After(endTime) {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "startDate must be earlier than endDate"})
			return
		}

		// Apply the date filters
		db = db.Where("savings_account.created_at BETWEEN ? AND ?", startTime, endTime)
	} else if startDate != "" {
		startDateInt, err := strconv.ParseInt(startDate, 10, 64)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid startDate"})
			return
		}
		startTime := time.Unix(startDateInt, 0)
		db = db.Where("savings_account.created_at >= ?", startTime)
	} else if endDate != "" {
		endDateInt, err := strconv.ParseInt(endDate, 10, 64)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid endDate"})
			return
		}
		endTime := time.Unix(endDateInt, 0)
		db = db.Where("savings_account.created_at <= ?", endTime)
	}

	// Count matching records only for the first page
	var collectionCount int64
	if pageToken == "" {
		if err := db.Count(&collectionCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count savings accounts"})
			return
		}
	}

	// Fetch savings accounts with limit
	accounts := make([]*SavingsAccountRead, 0, pageSize+1)
	if err := db.Find(&accounts).Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve savings accounts"})
		return
	}

	// Prepare the response data
	resultAccounts := make([]*SavingsAccountResponse, 0, len(accounts))

	for index, account := range accounts {
		// Skip the extra record used for checking the next page token
		if index == pageSize {
			break
		}

		resultAccounts = append(resultAccounts, ToSavingsAccountResponse(account))
	}

	// Generate the next page token if more records exist
	var nextPageToken string
	if len(accounts) > pageSize {
		// Next page token is the ID of the last record
		nextPageToken = base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(accounts[pageSize-1].SavingsAccount.ID)))
	}

	// Return paginated savings accounts
	c.IndentedJSON(http.StatusOK, gin.H{
		"next_page_token":  nextPageToken,
		"savings_accounts": resultAccounts,
		"collectionCount":  collectionCount,
	})
}

// UpdateSavingsAccountStatus facilitates actions like approve, activate, and close on a savings account
func (ctrl *SavingsAccountController) UpdateSavingsAccountStatus(c *gin.Context) {
	// Get the savings account ID from the URL path
	id := c.Param("id")

	// Bind the JSON input to DTO
	var dto UpdateSavingsAccountStatusDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch the savings account from the database
	var account SavingsAccount
	if err := ctrl.DB.Select("id").First(&account, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Savings account not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving savings account"})
		}
		return
	}

	var err error

	// Perform action based on the value of "action"
	switch dto.Action {
	case "approve":
		// Approve the account by setting DateApproved and StatusID
		if account.DateApproved.Valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Savings account is already approved"})
			return
		}
		now := time.Now()
		account.DateApproved = sql.NullTime{Time: now, Valid: true}
		account.StatusID = 2 // Approved status
		err = ctrl.DB.Updates(&account).Error

	case "activate":
		// Activate the account by setting DateActivated and StatusID
		if account.DateActivated.Valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Savings account is already activated"})
			return
		}
		now := time.Now()
		account.DateActivated = sql.NullTime{Time: now, Valid: true}
		account.StatusID = 3 // Activated status
		err = ctrl.DB.Updates(&account).Error

	case "close":
		// Close the account by setting DateClosed and StatusID
		if account.DateClosed.Valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Savings account is already closed"})
			return
		}
		now := time.Now()
		account.DateClosed = sql.NullTime{Time: now, Valid: true}
		account.StatusID = 4 // Closed status
		err = ctrl.DB.Updates(&account).Error
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "allowed actions are approve, activate or close"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return the updated account
	c.JSON(http.StatusOK, gin.H{
		"message": "Savings account updated successfully",
	})
}

// Validate timeGroup
var validTimeGroups = map[string]string{
	"yearly":  "YEAR(created_at)",
	"monthly": "DATE_FORMAT(created_at, '%Y-%m')",
	"daily":   "DATE(created_at)",
	"hourly":  "DATE_FORMAT(created_at, '%Y-%m-%d %H:00:00')",
}

// GetStats retrieves aggregated savings statistics for graphing based on filters.
func (ctrl *SavingsAccountController) GetStats(c *gin.Context) {
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
	db := ctrl.DB.WithContext(c.Request.Context()).Table("savings_account").
		Select(fmt.Sprintf(`
			%s AS time_group,
			COUNT(*) AS count,
			SUM(balance) AS volume
		`, groupByField)).
		Group("time_group").
		Order("time_group ASC")

	// Apply filters
	if productID != "" {
		db = db.Where("product_id = ?", productID)
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
