package main

import (
	"context"
	"flag"
	"net/http"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gidyon/gomicro/pkg/conn"
	"github.com/gidyon/gomicro/pkg/grpc/zaplogger"
	"github.com/gidyon/gomicro/utils/errs"
	"github.com/gidyon/pesapalm/internal/auth"
	"github.com/gidyon/pesapalm/internal/customer"
	loans_product "github.com/gidyon/pesapalm/internal/loan_product"
	"github.com/gidyon/pesapalm/internal/loans"
	"github.com/gidyon/pesapalm/internal/savings"
	"github.com/gidyon/pesapalm/internal/savings_product"
	"github.com/gidyon/pesapalm/internal/template"
	"github.com/gidyon/pesapalm/internal/user"
	"github.com/gidyon/pesapalm/pkg/api/sms"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc/grpclog"
	"gorm.io/gorm"
)

var (
	configFile = flag.String("config-file", ".env", "Configuration file")
	casbinConf = flag.String("casbin-conf", "casbin.conf", "Configuration file for casbin")
	dir        = flag.String("static-dir", "./", "Static assets parent directory")
)

var (
	appLogger grpclog.LoggerV2
	sqlDB     *gorm.DB
	redisDB   *redis.Client
)

func main() {
	flag.Parse()

	ctx := context.Background()

	// Config in .env
	viper.SetConfigFile(*configFile)

	viper.AutomaticEnv()

	// Read config from .env
	err := viper.ReadInConfig()
	errs.Panic(err)

	sqlDB, err = conn.OpenGorm(&conn.DbOptions{
		Name:     viper.GetString("MYSQL_NAME"),
		Dialect:  viper.GetString("MYSQL_DIALECT"),
		Address:  viper.GetString("MYSQL_ADDRESS"),
		User:     viper.GetString("MYSQL_USER"),
		Password: viper.GetString("MYSQL_PASSWORD"),
		Schema:   viper.GetString("MYSQL_SCHEMA"),
		ConnPool: &conn.DbPoolSettings{
			MaxIdleConns: viper.GetUint("MYSQL_MAX_IDLE_CONNS"),
			MaxOpenConns: viper.GetUint("MYSQL_MAX_OPEN_CONNS"),
			MaxLifetime:  viper.GetDuration("MYSQL_MAX_LIFETIME"),
		},
	})
	errs.Panic(err)

	sqlDB = sqlDB.Debug()

	// Open redis connection
	redisDB = redis.NewClient(&redis.Options{
		Network:         viper.GetString("REDIS_NETWORK"),
		Addr:            viper.GetString("REDIS_ADDRESS"),
		Username:        viper.GetString("REDIS_USERNAME"),
		Password:        viper.GetString("REDIS_PASSWORD"),
		DB:              viper.GetInt("REDIS_DBNAME"),
		MaxRetries:      viper.GetInt("REDIS_MAX_RETRIES"),
		ReadTimeout:     viper.GetDuration("REDIS_READ_TIMEOUT"),
		WriteTimeout:    viper.GetDuration("REDIS_WRITE_TIMEOUT"),
		MinIdleConns:    viper.GetInt("REDIS_MIN_IDLE_CONNS"),
		ConnMaxLifetime: viper.GetDuration("REDIS_MAX_CONN_AGE"),
	})

	// Initialize  casbin adapter
	gormAdapter, err := gormadapter.NewAdapterByDB(sqlDB)
	errs.Panic(err)

	// Cabin enforcer
	enforcer, err := casbin.NewEnforcer(*casbinConf, gormAdapter)
	errs.Panic(err)
	errs.Panic(enforcer.LoadPolicy())

	// Auth service
	appAuth := auth.NewAuthService(redisDB)
	tkMng := auth.NewTokenService(&auth.TokenOptions{
		AccessSecret:      viper.GetString("ACCESS_SECRET"),
		RefreshSecret:     viper.GetString("REFRESH_SECRET"),
		AccessExpiration:  viper.GetDuration("ACCESS_SECRET_DURATION"),
		RefreshExpiration: viper.GetDuration("REFRESH_SECRET_DURATION"),
	})

	router := gin.Default()

	// Cors handler
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Accept",
			"Access-Control-Allow-Origin",
			"Authorization",
			"Cache-Control",
			"Content-Type",
			"DNT",
			"If-Modified-Since",
			"Keep-Alive",
			"Origin",
			"User-Agent",
			"X-Requested-With",
		},
		ExposeHeaders:             []string{"Authorization"},
		MaxAge:                    1728,
		AllowCredentials:          true,
		OptionsResponseStatusCode: 200,
	}))

	// Handle all unrouted requests and redirect to "/"
	router.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/")
	})

	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// Initialize logger
	errs.Panic(zaplogger.Init(viper.GetInt("logLevel"), ""))

	zaplogger.Log = zaplogger.Log.WithOptions(zap.WithCaller(true))

	// gRPC logger compatible
	appLogger = zaplogger.ZapGrpcLoggerV2(zaplogger.Log)

	// User management API
	_, err = user.StartService(ctx, &user.Options{
		SqlDB:   sqlDB,
		RedisDB: redisDB,
		Logger:  appLogger,
		SMSAuth: &sms.SMSAuth{
			ApiUrl:   viper.GetString("SMS_API_URL"),
			SenderId: viper.GetString("SMS_SENDER_ID"),
			ApiKey:   viper.GetString("SMS_API_KEY"),
			ClientId: viper.GetString("SMS_CLIENT_ID"),
		},
		TokenManager: tkMng,
		Auth:         appAuth,
		GinEngine:    router,
	})
	errs.Panic(err)

	// Loan service
	loans.RegisterRoutes(&loans.Options{
		DB:           sqlDB,
		Logger:       appLogger,
		TokenManager: tkMng,
		GinEngine:    router,
	})

	// Loan products
	loans_product.RegisterRoutes(&loans_product.Options{
		DB:           sqlDB,
		Logger:       appLogger,
		TokenManager: tkMng,
		GinEngine:    router,
	})

	// Savings
	savings.RegisterRoutes(&savings.Options{
		DB:           sqlDB,
		Logger:       appLogger,
		TokenManager: tkMng,
		GinEngine:    router,
	})

	// Savings products
	savings_product.RegisterRoutes(&savings_product.Options{
		DB:           sqlDB,
		Logger:       appLogger,
		TokenManager: tkMng,
		GinEngine:    router,
	})

	// Customers
	customer.RegisterRoutes(&customer.Options{
		DB:           sqlDB,
		Logger:       appLogger,
		TokenManager: tkMng,
		GinEngine:    router,
	})

	// Templates
	template.RegisterRoutes(&template.Options{
		DB:           sqlDB,
		Logger:       appLogger,
		TokenManager: tkMng,
		GinEngine:    router,
	})

	if *dir != "" {
		router.Use(static.Serve("/", static.LocalFile(*dir, true)))
	}

	errs.Panic(router.Run(":" + viper.GetString("HTTP_PORT")))
}
