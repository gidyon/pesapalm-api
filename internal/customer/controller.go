package customer

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CustomerController structure
type CustomerController struct {
	*Options
}

// CreateCustomer creates a new customer
func (ctrl *CustomerController) CreateCustomer(c *gin.Context) {
	var dto CreateCustomerDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		ctrl.Logger.Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer := Customer{
		CustomerID:       dto.CustomerID,
		FirstName:        dto.FirstName,
		LastName:         dto.LastName,
		MiddleName:       dto.MiddleName,
		EmailAddress:     dto.EmailAddress,
		Gender:           dto.Gender,
		MSISDN1:          sql.NullString{String: dto.MSISDN1, Valid: dto.MSISDN1 != ""},
		MSISDN2:          sql.NullString{String: dto.MSISDN2, Valid: dto.MSISDN2 != ""},
		WorkPlaceAddress: sql.NullString{String: dto.WorkPlaceAddress, Valid: dto.WorkPlaceAddress != ""},
		HomeAddress:      sql.NullString{String: dto.HomeAddress, Valid: dto.HomeAddress != ""},
		Role:             sql.NullString{String: dto.Role, Valid: dto.Role != ""},
		Notes:            sql.NullString{String: dto.Notes, Valid: dto.Notes != ""},
		StatusID:         dto.StatusID,
		ProfilePicID:     dto.ProfilePicID,
		SignatureKeyID:   dto.SignatureKeyID,
		LanguageID:       dto.LanguageID,
		BranchID:         dto.BranchID,
	}

	if result := ctrl.DB.Create(&customer); result.Error != nil {
		ctrl.Logger.Errorf("Failed to create customer: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	ctrl.Logger.Infof("Successfully created customer with ID: %d", customer.ID)
	c.JSON(http.StatusOK, customer)
}

// ListCustomers retrieves a paginated list of customers
func (ctrl *CustomerController) ListCustomers(c *gin.Context) {
	queryParams := c.Request.URL.Query()

	// Parse pageSize from query, default if invalid
	pageSize, _ := strconv.Atoi(queryParams.Get("pageSize"))
	switch {
	case pageSize <= 0:
		pageSize = 10 // Default page size
	case pageSize > 100:
		pageSize = 100 // Maximum page size
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
		Limit(pageSize + 1). // Fetch one extra record to check for next page
		Order("id DESC").    // Order by descending ID for pagination
		Model(&Customer{})

	// Apply ID filter for pagination
	if id > 0 {
		db = db.Where("id < ?", id)
	}

	var collectionCount int64

	// Only count total collection on the first page (no page token)
	if pageToken == "" {
		if err := db.Count(&collectionCount).Error; err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to count customers"})
			return
		}
	}

	// Fetch customers with limit
	customers := make([]*Customer, 0, pageSize+1)
	if err := db.Find(&customers).Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve customers"})
		return
	}

	// Prepare the response data
	resultCustomers := make([]*CustomerResponse, 0, len(customers))

	for index, db := range customers {
		// Skip the extra record used for checking the next page token
		if index == pageSize {
			break
		}

		resultCustomers = append(resultCustomers, ToCustomerResponse(db))
	}

	// Generate the next page token if more records exist
	var nextPageToken string
	if len(customers) > pageSize {
		// Next page token is the ID of the last record
		nextPageToken = base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(customers[pageSize-1].ID)))
	}

	// Return paginated customers
	c.IndentedJSON(http.StatusOK, gin.H{
		"next_page_token": nextPageToken,
		"customers":       resultCustomers,
		"collectionCount": collectionCount,
	})
}

// GetCustomer retrieves a single customer by ID
func (ctrl *CustomerController) GetCustomer(c *gin.Context) {
	id := c.Param("id")
	var customer Customer
	if err := ctrl.DB.First(&customer, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customer"})
		}
		return
	}
	c.JSON(http.StatusOK, ToCustomerResponse(&customer))
}

// UpdateCustomer updates a customer's details
func (ctrl *CustomerController) UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	var dto UpdateCustomerDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var customer Customer
	if err := ctrl.DB.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
		return
	}

	// Update the fields
	customer.FirstName = dto.FirstName
	customer.LastName = dto.LastName
	customer.MiddleName = dto.MiddleName
	customer.EmailAddress = dto.EmailAddress
	customer.Gender = dto.Gender
	customer.MSISDN1 = sql.NullString{String: dto.MSISDN1, Valid: dto.MSISDN1 != ""}
	customer.MSISDN2 = sql.NullString{String: dto.MSISDN2, Valid: dto.MSISDN2 != ""}
	customer.WorkPlaceAddress = sql.NullString{String: dto.WorkPlaceAddress, Valid: dto.WorkPlaceAddress != ""}
	customer.HomeAddress = sql.NullString{String: dto.HomeAddress, Valid: dto.HomeAddress != ""}
	customer.Role = sql.NullString{String: dto.Role, Valid: dto.Role != ""}
	customer.Notes = sql.NullString{String: dto.Notes, Valid: dto.Notes != ""}
	customer.StatusID = dto.StatusID
	customer.UpdatedAt = time.Now()

	err := ctrl.DB.Updates(&customer).Error
	if err != nil {
		ctrl.Logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// DeleteCustomer deletes a customer by ID
func (ctrl *CustomerController) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.DB.Delete(&Customer{}, id).Error; err != nil {
		ctrl.Logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}

// Validate timeGroup
var validTimeGroups = map[string]string{
	"yearly":  "YEAR(created_at)",
	"monthly": "DATE_FORMAT(created_at, '%Y-%m')",
	"daily":   "DATE(created_at)",
	"hourly":  "DATE_FORMAT(created_at, '%Y-%m-%d %H:00:00')",
}

// GetStats retrieves aggregated customer statistics for graphing based on filters.
func (ctrl *CustomerController) GetStats(c *gin.Context) {
	// Parse request parameters
	queryParams := c.Request.URL.Query()

	statusID := queryParams.Get("status_id")
	branchID := queryParams.Get("branch_id")
	languageID := queryParams.Get("language_id")
	timeFrom := queryParams.Get("from")
	timeTo := queryParams.Get("to")
	timeGroup := queryParams.Get("time_group") // e.g., "yearly", "monthly", "daily", "hourly"

	groupByField, isValid := validTimeGroups[timeGroup]
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time group"})
		return
	}

	// Base query
	db := ctrl.DB.WithContext(c.Request.Context()).Table("customer").
		Select(fmt.Sprintf(`
			%s AS time_group,
			COUNT(*) AS count
		`, groupByField)).
		Group("time_group").
		Order("time_group ASC")

	// Apply filters
	if statusID != "" {
		db = db.Where("status_id = ?", statusID)
	}
	if branchID != "" {
		db = db.Where("branch_id = ?", branchID)
	}
	if languageID != "" {
		db = db.Where("language_id = ?", languageID)
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
		TimeGroup string `json:"time_group"`
		Count     int    `json:"count"`
	}
	if err := db.Find(&stats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve statistics"})
		return
	}

	// Prepare the response
	response := map[string]interface{}{
		"data": stats,
		"count": func() int {
			total := 0
			for _, stat := range stats {
				total += stat.Count
			}
			return total
		}(),
	}

	// Send the response
	c.JSON(http.StatusOK, response)
}
