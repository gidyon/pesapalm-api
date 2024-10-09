package user

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gidyon/pesapalm/internal/auth"
	sms_app "github.com/gidyon/pesapalm/internal/sms"
	"github.com/gidyon/pesapalm/pkg/api/sms"
	"github.com/gidyon/pesapalm/pkg/utils/formatutil"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/grpclog"
	"gorm.io/gorm"
)

type Options struct {
	SqlDB        *gorm.DB
	RedisDB      *redis.Client
	Logger       grpclog.LoggerV2
	SMSAuth      *sms.SMSAuth
	TokenManager auth.TokenInterface
	Auth         auth.AuthInterface
	GinEngine    *gin.Engine
}

type APIServer struct {
	*Options
}

// StartService creates a user API singleton
func StartService(ctx context.Context, opt *Options) (_ *APIServer, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("Failed to start user service: %v", err)
		}
	}()

	// Validation
	switch {
	case ctx == nil:
		err = errors.New("missing context")
	case opt == nil:
		err = errors.New("missing options")
	case opt.SqlDB == nil:
		err = errors.New("missing sql db")
	case opt.RedisDB == nil:
		err = errors.New("missing redis db")
	case opt.Auth == nil:
		err = errors.New("missing auth")
	case opt.TokenManager == nil:
		err = errors.New("missing token manager")
	case opt.GinEngine == nil:
		err = errors.New("missing gin engine")
	}
	if err != nil {
		return nil, err
	}

	// Account API
	api := &APIServer{
		Options: opt,
	}

	usersTable = viper.GetString("accounts_table")

	// Perform auto migration
	if !api.SqlDB.WithContext(ctx).Migrator().HasTable((&User{}).TableName()) {
		err = api.SqlDB.WithContext(ctx).AutoMigrate(&User{})
		if err != nil {
			return nil, fmt.Errorf("failed to automigrate %s table: %v", (&User{}).TableName(), err)
		}
	}

	// Register routes
	api.registerRoutes()

	return api, nil
}

const (
	defaultPageSize = 200
	maxPageSize     = 1000
)

func (api *APIServer) Login(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		err error
	)

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request provided"})
		return
	}

	db := &User{}

	// Get account
	switch {
	case strings.Contains(req.Username, "@"):
		err = api.SqlDB.WithContext(ctx).First(db, "email=?", req.Username).Error
	default:
		err = api.SqlDB.WithContext(ctx).First(db, "phone=?", req.Username).Error
	}

	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		emailOrPhone := func() string {
			if strings.Contains(req.Username, "@") {
				return "email " + req.Username
			}
			if strings.Contains(req.Username, "+") {
				return "phone " + req.Username
			}
			return "username " + req.Username
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("account with %s not found", emailOrPhone())})
		return
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to login"})
		return
	}

	// If no password set in account
	if db.Password == "" || db.AccountStatus == "RESET_PASSWORD" {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "reset password"})
		return
	}

	// Check account statuses
	switch strings.ToLower(db.AccountStatus) {
	case "blocked":
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "account is blocked"})
		return
	case "inactive":
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "account is inactive"})
		return
	case "invited":
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "account is inactive"})
		return
	}

	// Check if password match if they logged in with Phone or Email
	err = compareHash(db.Password, req.Password)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "wrong password"})
		return
	}

	api.updateSession(ctx, c, db)
}

func (api *APIServer) Refresh(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		err error
	)

	// If metadata is passed and the tokens valid, delete them from the redis store
	metadata, err := api.TokenManager.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if metadata == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "please login"})
		return
	}

	// Fetch user id
	userId, err := api.Auth.FetchAuth(ctx, metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "please login"})
		return
	}

	db := &User{}

	// Get account
	err = api.SqlDB.WithContext(ctx).First(db, "id=?", userId).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "please login"})
		return
	}

	// Check account statuses
	switch strings.ToLower(db.AccountStatus) {
	case "blocked":
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "account is blocked"})
		return
	case "inactive":
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "account is inactive"})
		return
	case "invited":
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "account is inactive"})
		return
	}

	api.updateSession(ctx, c, db)
}

func (api *APIServer) Logout(c *gin.Context) {
	// If metadata is passed and the tokens valid, delete them from the redis store
	metadata, _ := api.TokenManager.ExtractTokenMetadata(c.Request)
	if metadata != nil {
		deleteErr := api.Auth.DeleteTokens(c.Request.Context(), metadata)
		if deleteErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": deleteErr.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

var emailRegex = regexp.MustCompile(`^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,10})+$`)

func ValidateUser(pb User_) error {
	switch {
	case pb.Phone == "":
		return errors.New("missing phone")
	case pb.Names == "":
		return errors.New("missing names")
	case pb.PrimaryGroup == "":
		return errors.New("missing role")
	default:
		// Validate email
		if pb.Email != "" && !emailRegex.MatchString(pb.Email) {
			return errors.New("incorrect email")
		}

		// Validate phone
		switch {
		case strings.HasPrefix(pb.Phone, "07") && len(pb.Phone) != 10:
			return errors.New("incorrect phone")
		case strings.HasPrefix(pb.Phone, "01") && len(pb.Phone) != 10:
			return errors.New("incorrect phone")
		case strings.HasPrefix(pb.Phone, "254") && len(pb.Phone) != 12:
			return errors.New("incorrect phone")
		}
	}

	return nil
}

func (api *APIServer) CreateUser(c *gin.Context) {
	var (
		ctx  = c.Request.Context()
		user User_
		err  error
	)

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json"})
		return
	}

	// Validate user
	err = ValidateUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Clean phone
	user.Phone = formatutil.FormatPhoneKE(strings.TrimPrefix(user.Phone, "+"))

	// Check if user exists
	db := &User{}

	err = api.SqlDB.WithContext(ctx).Select("email,phone").First(db, "email=? or phone=?", user.Email, user.Phone).Error
	switch {
	case err == nil:
		if db.Phone.String == user.Phone {
			c.JSON(http.StatusBadRequest, gin.H{"message": "phone exists"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "email exists"})
		}
		return
	case errors.Is(err, gorm.ErrRecordNotFound):
	default:
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to check if user exists"})
		return
	}

	var creatorId uint64

	// Get metadata
	metadata, err := api.TokenManager.ExtractTokenMetadata(c.Request)
	if err == nil {
		creatorId = metadata.UserId
	}

	var password string

	if user.Password != "" {
		// Get password
		password, err = genHash(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate password"})
			return
		}
	}

	var bs []byte
	if len(user.GeneralData) > 0 {
		bs, err = json.Marshal(user.GeneralData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to marshal generate data"})
			return
		}
	}

	// Save to database
	err = api.SqlDB.Create(&User{
		ID:        0,
		CreatorId: creatorId,
		Phone: sql.NullString{
			String: user.Phone,
			Valid:  user.Phone != "",
		},
		Email: sql.NullString{
			String: user.Email,
			Valid:  user.Email != "",
		},
		Names:     user.Names,
		BirthDate: sql.NullTime{},
		Gender:    user.Gender,
		ProfileURL: sql.NullString{
			String: user.ProfileURL,
			Valid:  user.ProfileURL != "",
		},
		Country: sql.NullString{
			String: user.Country,
			Valid:  user.Country != "",
		},
		CountryCode: sql.NullString{
			String: user.CountryCode,
			Valid:  user.CountryCode != "",
		},
		GroupId: sql.NullInt64{
			Int64: int64(user.GroupId),
			Valid: user.GroupId != 0,
		},
		Password:      password,
		GeneralData:   bs,
		PrimaryGroup:  user.PrimaryGroup,
		AccountStatus: "ACTIVE",
		LastLoginIp:   sql.NullString{},
		LastLogin:     sql.NullTime{},
		UpdatedAt:     time.Time{},
		CreatedAt:     time.Time{},
	}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (api *APIServer) UpdateUser(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		userId = c.Param("userId")
		user   User_
		err    error
	)

	if err = c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json"})
		return
	}

	db := &User{}

	// Get account
	err = api.SqlDB.WithContext(ctx).Select("id,primary_group,creator_id,group_id").First(db, "id=?", userId).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to get user"})
		return
	}

	var password string
	if user.Password != "" {
		password, err = genHash(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate password"})
			return
		}
	}

	var bs []byte
	if len(user.GeneralData) > 0 {
		bs, err = json.Marshal(user.GeneralData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to marshal generate data"})
			return
		}
	}

	// Update user
	err = api.SqlDB.Model(db).Updates(&User{
		ID: 0,
		Phone: sql.NullString{
			String: user.Phone,
			Valid:  user.Phone != "",
		},
		Email: sql.NullString{
			String: user.Email,
			Valid:  user.Email != "",
		},
		Names:     user.Names,
		BirthDate: sql.NullTime{},
		Gender:    user.Gender,
		ProfileURL: sql.NullString{
			String: user.ProfileURL,
			Valid:  user.ProfileURL != "",
		},
		Country: sql.NullString{
			String: user.Country,
			Valid:  user.Country != "",
		},
		CountryCode: sql.NullString{
			String: user.CountryCode,
			Valid:  user.CountryCode != "",
		},
		Password:      password,
		GeneralData:   bs,
		PrimaryGroup:  user.PrimaryGroup,
		AccountStatus: user.AccountStatus,
	}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully updated"})
}

func (api *APIServer) GetUser(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		userId = c.Param("userId")
		err    error
	)

	db := &User{}

	// Get account
	err = api.SqlDB.WithContext(ctx).First(db, "id=?", userId).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to get user"})
		return
	}

	pb := &User_{
		ID:            db.ID,
		Phone:         db.Phone.String,
		Email:         db.Email.String,
		Names:         db.Names,
		Gender:        db.Gender,
		ProfileURL:    db.ProfileURL.String,
		Country:       db.Country.String,
		CountryCode:   db.CountryCode.String,
		GroupId:       db.GroupId.Int64,
		GeneralData:   map[string]any{},
		PrimaryGroup:  db.PrimaryGroup,
		AccountStatus: db.AccountStatus,
		LastLogin:     "",
		UpdatedAt:     db.UpdatedAt.UTC().Format(time.RFC3339),
		CreatedAt:     db.CreatedAt.UTC().Format(time.RFC3339),
	}

	if len(db.GeneralData) > 0 {
		err = json.Unmarshal(db.GeneralData, &pb.GeneralData)
		if err != nil {
			return
		}
	}

	if db.LastLogin.Valid {
		pb.LastLogin = db.LastLogin.Time.UTC().Format(time.RFC3339)
	}
	if db.BirthDate.Valid {
		pb.BirthDate = db.BirthDate.Time.UTC().Format(time.RFC3339)
	}

	c.JSON(http.StatusOK, pb)
}

func (api *APIServer) ListUsers(c *gin.Context) {
	queryParams := c.Request.URL.Query()

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
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "page token in incorrect"})
			return
		}
		id, err = strconv.Atoi(string(bs))
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "page token in incorrect"})
			return
		}
	}

	db := api.SqlDB.WithContext(c.Request.Context()).Limit(int(pageSize) + 1).Order("id DESC").Model(&User{})

	// ID filter
	if id > 0 {
		db = db.Where("id<?", id)
	}

	var collectionCount int64

	// Page token
	if pageToken == "" {
		err := db.Count(&collectionCount).Error
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to count users"})
			return
		}
	}

	dbs := make([]*User, 0, pageSize+1)

	err := db.Find(&dbs).Error
	switch {
	case err == nil:
	default:
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to find users"})
		return
	}

	pbs := make([]*User_, 0, len(dbs))

	for index, db := range dbs {
		if index == pageSize {
			break
		}

		pb := &User_{
			ID:            db.ID,
			Phone:         db.Phone.String,
			Email:         db.Email.String,
			Names:         db.Names,
			BirthDate:     "",
			Gender:        db.Gender,
			ProfileURL:    db.ProfileURL.String,
			Country:       db.Country.String,
			CountryCode:   db.CountryCode.String,
			GroupId:       db.GroupId.Int64,
			GeneralData:   map[string]any{},
			PrimaryGroup:  db.PrimaryGroup,
			Password:      "",
			AccountStatus: db.AccountStatus,
			LastLogin:     "",
			UpdatedAt:     db.UpdatedAt.UTC().Format(time.RFC3339),
			CreatedAt:     db.CreatedAt.UTC().Format(time.RFC3339),
		}

		if len(db.GeneralData) > 0 {
			err = json.Unmarshal(db.GeneralData, &pb.GeneralData)
			if err != nil {
				return
			}
		}

		if db.LastLogin.Valid {
			pb.LastLogin = db.LastLogin.Time.UTC().Format(time.RFC3339)
		}
		if db.BirthDate.Valid {
			pb.BirthDate = db.BirthDate.Time.UTC().Format(time.RFC3339)
		}

		pbs = append(pbs, pb)
	}

	var token string
	if len(dbs) > pageSize {
		// Next page token
		token = base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(id)))
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"next_page_token": token,
		"users":           pbs,
		"collectionCount": collectionCount,
	})
}

func getOTPKey(ID uint64) string {
	return fmt.Sprintf("otp:%v", ID)
}

func getTrialsKey(ID uint64) string {
	return fmt.Sprintf("logintrials:%v", ID)
}

const OTPExpireDuration = time.Minute * 10

func (api *APIServer) RequestOtp(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		err error
	)

	var req RequestOtpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request provided"})
		return
	}

	db := &User{}

	// Get account
	err = api.SqlDB.WithContext(ctx).Select("id,phone").First(db, "phone=?", formatutil.FormatPhoneKE(req.Phone)).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusBadRequest, gin.H{"message": "account not found"})
		return
	default:
		api.Logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "try again later"})
		return
	}

	// otp := randomdata.Number(100000, 999999)
	otp := 123456

	// Send sms
	data := fmt.Sprintf("Login OTP for Therapy Assessment. \n\nOTP is %d \nExpires in %s", otp, OTPExpireDuration)

	err = sms_app.SendSMS(ctx, &sms.SendSMSRequest{
		Sms: &sms.SMS{
			DestinationPhones: []string{db.Phone.String},
			Keyword:           "LoginOTP",
			Message:           data,
		},
		Auth:     api.SMSAuth,
		Provider: sms.SmsProvider_ONFON,
	}, viper.GetString("ENV"))
	if err != nil {
		api.Logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to send otp"})
		return
	}

	// Set token with expiration of 5 minutes
	err = api.RedisDB.Set(ctx, getOTPKey(db.ID), otp, OTPExpireDuration).Err()
	if err != nil {
		api.Logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to set otp to cache"})
		return
	}

	// Set trials initial value to zero
	err = api.RedisDB.Set(ctx, getTrialsKey(db.ID), 0, OTPExpireDuration).Err()
	if err != nil {
		api.Logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to set otp counter"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "otp sent to end user"})
}

const maxTrials = 4

var (
	BlockedState   = "BLOCKED"
	InactiveState  = "INACTIVE"
	ActiveState    = "ACTIVE"
	CredentialsMap = map[string]struct{}{}
)

func (api *APIServer) ValidateOtp(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		err error
	)

	var req ValidateOtpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request provided"})
		return
	}

	db := &User{}

	// Get account
	err = api.SqlDB.WithContext(ctx).First(db, "phone=?", formatutil.FormatPhoneKE(req.Phone)).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusBadRequest, gin.H{"message": "account not found"})
		return
	default:
		api.Logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve account"})
		return
	}

	trialsKey := getTrialsKey(db.ID)

	// Increment trials by 1
	trials, err := api.RedisDB.Incr(ctx, trialsKey).Result()
	switch {
	case err == nil:
	case errors.Is(err, redis.Nil):
	default:
		api.Logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get otp counter"})
		return
	}

	// Check if exceed trials
	if trials > maxTrials {
		// Block the account
		err = api.SqlDB.WithContext(ctx).Model(db).Update("account_status", BlockedState).Error
		if err != nil {
			api.Logger.Errorln(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to block account"})
		}

		// Delete key
		err = api.RedisDB.Del(ctx, trialsKey).Err()
		if err != nil {
			api.Logger.Errorln(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to clear otp counter"})
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "account is blocked due to too many attempts"})

		return
	}

	// Get otp
	otp, err := api.RedisDB.Get(ctx, getOTPKey(db.ID)).Result()
	switch {
	case err == nil:
	case errors.Is(err, redis.Nil):
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Login OTP expired, request another OTP"})
		return
	default:
		api.Logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get otp"})
		return
	}

	// Compare otp
	if otp != req.Otp {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Login OTP do not match"})
		return
	}

	// Delete key
	err = api.RedisDB.Del(ctx, trialsKey).Err()
	if err != nil {
		api.Logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to remove otp key"})
		return
	}

	api.updateSession(ctx, c, db)
}

func (api *APIServer) updateSession(ctx context.Context, c *gin.Context, db *User) {
	// Check account statuses
	switch strings.ToLower(db.AccountStatus) {
	case BlockedState:
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "account is blocked"})
		return
	case InactiveState:
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "account is inactive"})
		return
	}

	// Send back
	user := &User_{
		ID:            db.ID,
		Phone:         db.Phone.String,
		Email:         db.Email.String,
		Names:         db.Names,
		BirthDate:     "",
		Gender:        "",
		ProfileURL:    "",
		Country:       "",
		CountryCode:   "",
		GroupId:       db.GroupId.Int64,
		GeneralData:   map[string]any{},
		PrimaryGroup:  db.PrimaryGroup,
		Password:      "",
		AccountStatus: db.AccountStatus,
		LastLogin:     db.LastLogin.Time.Format(time.RFC3339),
		UpdatedAt:     db.UpdatedAt.Format(time.RFC3339),
		CreatedAt:     db.CreatedAt.Format(time.RFC3339),
	}

	// Generate token
	token, err := api.TokenManager.CreateToken(db.ID, db.Names)
	if err != nil {
		api.Logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	// Get credential
	_, ok := CredentialsMap[user.Phone]

	// Update last login
	err = api.SqlDB.WithContext(ctx).Model(db).Update("last_login", time.Now()).Error
	if err != nil {
		api.Logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update last login"})
		return
	}

	// Set the token in cache
	err = api.Auth.CreateAuth(ctx, db.ID, token)
	if err != nil {
		api.Logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to set token to cache"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":         token,
		"user":          user,
		"hasCredential": ok,
	})
}

func (api *APIServer) RefreshSession(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		err error
	)

	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request provided"})
		return
	}

	// Get refresh token in cache
	ID, err := api.Auth.FetchAuth(ctx, req.RefreshUuid)
	if err != nil {
		api.Logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get refresh uuid"})
		return
	}

	db := &User{}

	// Get account
	err = api.SqlDB.WithContext(ctx).First(db, "id=?", ID).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusBadRequest, gin.H{"message": "account not found"})
		return
	default:
		api.Logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve account"})
		return
	}

	api.updateSession(ctx, c, db)
}

// Constants
const MaxTrials = 3

// ResetPasswordRequest struct for reset password request
type ResetPasswordRequest struct {
	Username    string `json:"username" binding:"required"`
	OTP         string `json:"otp" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// RequestResetPasswordOtp sends an OTP to the user
func (api *APIServer) RequestResetPasswordOtp(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		err error
	)

	var req RequestOtpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request provided"})
		return
	}

	db := &User{}

	// Check if account exists
	err = api.SqlDB.WithContext(ctx).Select("id, phone").First(db, "phone = ?", formatutil.FormatPhoneKE(req.Phone)).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusBadRequest, gin.H{"message": "account not found"})
		return
	default:
		api.Logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "try again later"})
		return
	}

	// Generate OTP (for production, use random generation)
	otp := 123456

	// Send OTP via SMS
	message := fmt.Sprintf("Reset password OTP is %d. It expires in %s.", otp, OTPExpireDuration)
	err = sms_app.SendSMS(ctx, &sms.SendSMSRequest{
		Sms: &sms.SMS{
			DestinationPhones: []string{db.Phone.String},
			Keyword:           "PasswordReset",
			Message:           message,
		},
		Auth:     api.SMSAuth,
		Provider: sms.SmsProvider_ONFON,
	}, viper.GetString("ENV"))
	if err != nil {
		api.Logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to send otp"})
		return
	}

	// Store OTP in Redis with expiration
	err = api.RedisDB.Set(ctx, getResetPassOTPKey(db.ID), otp, OTPExpireDuration).Err()
	if err != nil {
		api.Logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to set otp in cache"})
		return
	}

	// Store trials count in Redis
	err = api.RedisDB.Set(ctx, getResetPassTrialsKey(db.ID), 0, OTPExpireDuration).Err()
	if err != nil {
		api.Logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to initialize trials count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

// ResetPassword resets the password of the user
func (api *APIServer) ResetPassword(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		err error
	)

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request provided"})
		return
	}

	db := &User{}

	// Check if the account exists
	err = api.SqlDB.WithContext(ctx).Select("id, phone").First(db, "phone = ?", req.Username).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusBadRequest, gin.H{"message": "account not found"})
		return
	default:
		api.Logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "try again later"})
		return
	}

	// Check if OTP matches the stored OTP
	storedOtp, err := api.RedisDB.Get(ctx, getResetPassOTPKey(db.ID)).Result()
	if err == redis.Nil || storedOtp != req.OTP {
		// Increment trials count
		trials, _ := api.RedisDB.Incr(ctx, getResetPassTrialsKey(db.ID)).Result()
		if trials > MaxTrials {
			c.JSON(http.StatusTooManyRequests, gin.H{"message": "maximum OTP attempts reached"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid OTP"})
		return
	} else if err != nil {
		api.Logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to verify otp"})
		return
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		api.Logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to hash password"})
		return
	}

	// Update the user's password in the database
	err = api.SqlDB.WithContext(ctx).Model(&User{}).Where("id = ?", db.ID).Update("password", hashedPassword).Error
	if err != nil {
		api.Logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update password"})
		return
	}

	// Clear OTP and trials count from Redis
	api.RedisDB.Del(ctx, getResetPassOTPKey(db.ID), getResetPassTrialsKey(db.ID))

	c.JSON(http.StatusOK, gin.H{"message": "password reset successful"})
}

// Helper functions to generate Redis keys
func getResetPassOTPKey(userID uint64) string {
	return fmt.Sprintf("otp:resetpass:%d", userID)
}

func getResetPassTrialsKey(userID uint64) string {
	return fmt.Sprintf("trials:resetpass:%d", userID)
}
