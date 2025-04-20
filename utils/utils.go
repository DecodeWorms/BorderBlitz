package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/DecodeWorms/BorderBlitz/config"
	handlers "github.com/DecodeWorms/BorderBlitz/handler"
	"github.com/DecodeWorms/BorderBlitz/models"
	"github.com/DecodeWorms/BorderBlitz/repository"
	services "github.com/DecodeWorms/BorderBlitz/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetUpDatabase sets up database connection
func SetUpDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate the schema
	err = db.AutoMigrate(
		&models.Stablecoin{},
		&models.Wallet{},
		&models.Transaction{},
		&models.AuditLog{},
		&models.Users{},
		&models.Balance{},
	)
	if err != nil {
		return nil, err
	}

	// Initialize default stablecoins if they don't exist
	if err := initializeStablecoins(db); err != nil {
		return nil, err
	}

	return db, nil
}

func initializeStablecoins(db *gorm.DB) error {
	stablecoins := []models.Stablecoin{
		{Symbol: "cNGN", Name: "Nigerian Naira Coin", USDRate: 0.002},              // 1 USD = 500 cNGN
		{Symbol: "cXAF", Name: "Central African CFA Franc Coin", USDRate: 0.00167}, // 1 USD = 600 cXAF
		{Symbol: "USDx", Name: "US Dollar Coin", USDRate: 1.0},
		{Symbol: "EURx", Name: "Euro Coin", USDRate: 1.1}, // 1 EUR = 1.1 USD
	}

	for _, coin := range stablecoins {
		result := db.FirstOrCreate(&coin, models.Stablecoin{Symbol: coin.Symbol})
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

// SetUpUserService sets up the Users services with it dependency
func SetUpUserService(user *repository.Users) *services.Users {
	return services.NewUsers(user)
}

// SetUpWalletService sets up Wallet service with dependencies
func SetUpWalletService(wallet repository.Wallet, stableCoin repository.StableCoin, user repository.User) *services.WalletService {
	return services.NewWalletService(wallet, stableCoin, user)
}

// SetUpTransactionService sets up Transaction service with dependencies
func SetUpTransactionService(db *gorm.DB, txRep *repository.TransactionRepository, wallet *repository.WalletRepository, coin *repository.StablecoinRepository) *services.TransactionService {
	return services.NewTransactionService(db, txRep, wallet, coin)
}

// SetUpSwapServices sets up Swap service with dependencies
func SetUpSwapService(db *gorm.DB, wallet *repository.WalletRepository, coin *repository.StablecoinRepository, tx *repository.TransactionRepository) *services.SwapService {
	return services.NewSwapService(db, wallet, coin, tx)
}

// SetUpFxService sets up Fx service with dependencies
func SetUpFxService(coin *repository.StablecoinRepository) *services.FXService {
	return services.NewFXService(coin)
}

func SetUpUserHandler(user *services.Users) *handlers.Users {
	return handlers.NewUserHandler(user)
}

// SetUpWalletHandler sets up Wallet handler with dependencies
func SetUpWalletHandler(swapService *services.WalletService) *handlers.WalletHandler {
	return handlers.NewWalletHandler(swapService)
}

// SetUpTransactionHandler sets up transaction handler with dependencies
func SetUpTransactionHandler(txService *services.TransactionService) *handlers.TransactionHandler {
	return handlers.NewTransactionHandler(txService)
}

// SetUpSwapHandler sets up set up swap handler with dependencies
func SetUpSwapHandler(swapService *services.SwapService) *handlers.SwapHandler {
	return handlers.NewSwapHandler(swapService)
}

// SetUpExplorerHandler sets up
func SetUpExplorerHandler(txService *services.TransactionService, walletService *services.WalletService) *handlers.ExplorerHandler {
	return handlers.NewExplorerHandler(txService, walletService)
}

// Middleware to log request information for compliance
func complianceMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// After request is processed, log for compliance
		userIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// In a real system, you would use a geolocation service to determine country
		country := "Unknown"

		// Extract browser info from user agent
		browser := "Unknown"
		if len(userAgent) > 0 {
			// Simple browser detection, would be more sophisticated in production
			if strings.Contains(userAgent, "Chrome") {
				browser = "Chrome"
			} else if strings.Contains(userAgent, "Firefox") {
				browser = "Firefox"
			} else if strings.Contains(userAgent, "Safari") {
				browser = "Safari"
			} else if strings.Contains(userAgent, "Edge") {
				browser = "Edge"
			}
		}

		// Create audit log entry
		auditLog := models.AuditLog{
			UserIP:    userIP,
			UserAgent: userAgent,
			Country:   country,
			Browser:   browser,
			Action:    c.Request.Method + " " + c.Request.URL.Path,
			Details:   fmt.Sprintf("Status: %d", c.Writer.Status()),
		}

		// Try to extract wallet ID from request if available
		walletID := c.Param("id")
		if walletID != "" {
			id := uint(0)
			_, err := fmt.Sscanf(walletID, "%d", &id)
			if err == nil && id > 0 {
				auditLog.WalletID = &id
			}
		}

		db.Create(&auditLog)
	}
}

// SetUpRouter sets up router to direct http calls to appropriate handler
func SetUpRouter(explorer *handlers.ExplorerHandler, tx *handlers.TransactionHandler, swap *handlers.SwapHandler, wallet *handlers.WalletHandler, user *handlers.Users, db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// Add Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.Default())
	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", //local
			"https://borderblitz-frontend.onrender.com", // for prod on a Render platform .
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "x-user-id"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}))
	router.Use(complianceMiddleware(db))

	api := router.Group("/api/v1")
	{
		//API endpoints for user
		api.POST("/user", user.CreateUser())

		//API endpoints for Wallets
		api.POST("/wallet", wallet.CreateWallet)
		api.GET("/wallet", wallet.GetWallet)
		api.POST("/wallet/deposit", wallet.Deposit)
		api.GET("/wallet/user", wallet.GetWalletByUser)

		//API endpoints for Transactions
		api.POST("/transaction/deposit", tx.Deposit)
		api.POST("/transaction/transfer", tx.Transfer)
		api.GET("/transaction", tx.GetTransaction)
		api.GET("/transaction/transactions", tx.GetTransactionHistory)

		//API endpoints for Swap
		api.GET("/swap/exchange_rate", swap.GetExchangeRate)
		api.POST("/swap", swap.Swap)

		//API endpoints for Explorer
		api.GET("/explorer", explorer.GetWalletOverview)
	}

	return router
}

func StartServer(router *gin.Engine) {
	//var c config.Config
	var c = config.NewConfig()
	interruptHandler := make(chan os.Signal, 1)
	signal.Notify(interruptHandler, syscall.SIGTERM, syscall.SIGINT)

	addr := fmt.Sprintf(":%s", c.AppPort)
	go func(addr string) {
		log.Printf("BorderBlitz.sv API service running on %v. Environment=%s", addr, c.AppEnv)
		if err := http.ListenAndServe(addr, router); err != nil {
			log.Printf("Error starting server: %v", err)
		}
	}(addr)

	<-interruptHandler
}
