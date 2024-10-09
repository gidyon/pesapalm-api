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

	if result := ctrl.DB.Create(&account); result.Error != nil {
		ctrl.Logger.Errorf("Failed to create savings account: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	ctrl.Logger.Infof("Successfully created savings account with ID: %d", account.ID)
	c.JSON(http.StatusOK, account)
}

// GetSavingsAccount retrieves a savings account by ID, following the GetUser pattern
func (ctrl *SavingsAccountController) GetSavingsAccount(c *gin.Context) {
	var (
		ctx       = c.Request.Context()
		savingsID = c.Param("id") // ID from the URL path
		err       error
	)

	// Create a placeholder for the savings account record
	db := &SavingsAccount{}

	// Fetch the savings account by ID from the database
	err = ctrl.DB.WithContext(ctx).First(db, "id = ?", savingsID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctrl.Logger.Warningf("Savings account with ID %s not found", savingsID)
			c.JSON(http.StatusNotFound, gin.H{"message": "Savings account not found"})
		} else {
			ctrl.Logger.Errorf("Failed to get savings account with ID %s: %v", savingsID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get savings account"})
		}
		return
	}

	ctrl.Logger.Infof("Retrieved savings account with ID %s", savingsID)

	// Return the savings account details as JSON
	c.JSON(http.StatusOK, ToSavingsAccountResponse(db))
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
	if result := ctrl.DB.First(&account, id); result.Error != nil {
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
	if result := ctrl.DB.Delete(&SavingsAccount{}, id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Savings account deleted"})
}

const (
	defaultPageSize = 200
	maxPageSize     = 1000
)

func (ctrl *SavingsAccountController) ListSavingsAccounts(c *gin.Context) {
	queryParams := c.Request.URL.Query()

	// Parse pageSize from query, default if invalid
	pageSize, _ := strconv.Atoi(queryParams.Get("pageSize"))
	switch {
	case pageSize <= 0:
		pageSize = defaultPageSize
	case pageSize > defaultPageSize:
		pageSize = defaultPageSize
	}

	var id int

	// Get last id from page token
	pageToken := queryParams.Get("pageToken")
	if pageToken != "" {
		bs, err := base64.StdEncoding.DecodeString(pageToken)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "page token is incorrect"})
			return
		}
		id, err = strconv.Atoi(string(bs))
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "page token is incorrect"})
			return
		}
	}

	db := ctrl.DB.WithContext(c.Request.Context()).
		Limit(int(pageSize) + 1).
		Order("id DESC").
		Model(&SavingsAccount{})

	// ID filter for pagination
	if id > 0 {
		db = db.Where("id < ?", id)
	}

	var collectionCount int64

	// Only count total collection on the first page
	if pageToken == "" {
		err := db.Count(&collectionCount).Error
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to count savings accounts"})
			return
		}
	}

	// Fetch savings accounts with limit
	accounts := make([]*SavingsAccount, 0, pageSize+1)
	err := db.Find(&accounts).Error
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve savings accounts"})
		return
	}

	// Prepare the response data
	resultAccounts := make([]*SavingsAccountResponse, 0, len(accounts))

	for index, db := range accounts {
		// Skip the extra record used for checking next page token
		if index == pageSize {
			break
		}

		resultAccounts = append(resultAccounts, ToSavingsAccountResponse(db))
	}

	// Generate the next page token if more records exist
	var nextPageToken string
	if len(accounts) > pageSize {
		// Next page token is the ID of the last record
		nextPageToken = base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(accounts[pageSize-1].ID)))
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
		ctrl.DB.Updates(&account)

	case "activate":
		// Activate the account by setting DateActivated and StatusID
		if account.DateActivated.Valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Savings account is already activated"})
			return
		}
		now := time.Now()
		account.DateActivated = sql.NullTime{Time: now, Valid: true}
		account.StatusID = 3 // Activated status
		ctrl.DB.Updates(&account)

	case "close":
		// Close the account by setting DateClosed and StatusID
		if account.DateClosed.Valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Savings account is already closed"})
			return
		}
		now := time.Now()
		account.DateClosed = sql.NullTime{Time: now, Valid: true}
		account.StatusID = 4 // Closed status
		ctrl.DB.Updates(&account)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "allowed actions are approve, activate or close"})
		return
	}

	// Return the updated account
	c.JSON(http.StatusOK, gin.H{
		"message": "Savings account updated successfully",
	})
}
