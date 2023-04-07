package services

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"ledger-service/logger"
	"ledger-service/models"
	"time"
)

const (
	MaxPages            = 100
	TransactionPageSize = 10
)

// AddFunds adds funds to a user's account and creates a transaction record.
// Uses a distributed lock to prevent race conditions between concurrent requests.
func AddFunds(ctx *gin.Context, uid string, request models.AddFundsRequest) (string, error) {
	amount := request.Amount
	if amount <= 0 {
		return "", errors.New("Amount must be positive")
	}

	var user models.User
	// Find or create user with given UID
	result := DbConnection.FirstOrCreate(&user, models.User{UID: uid})
	if result.Error != nil {
		logger.Error("Error finding or creating user: %v\n", zap.Error(result.Error))
		return "", errors.New("Internal Server Error")
	}

	// Generate a unique transaction ID
	transactionID := uuid.New().String()

	// Check if the transaction ID already exists
	var existingTransaction models.Transaction
	result = DbConnection.Where("transaction_id = ?", transactionID).First(&existingTransaction)
	if result.Error != nil && !result.RecordNotFound() {
		logger.Error("Error checking for existing transaction: %v\n", zap.Error(result.Error))
		return "", errors.New("Internal Server Error")
	}

	if !result.RecordNotFound() {
		return "", errors.New("Transaction already processed")
	}

	// Acquire distributed lock using Redsync
	redisClient := GetRedisDefaultClient()
	redsyncPool := goredis.NewPool(redisClient)
	rs := redsync.New(redsyncPool)

	mutex := rs.NewMutex("balance_mutex:" + uid)
	if err := mutex.Lock(); err != nil {
		logger.Error("Error acquiring lock: %v\n", zap.Error(err))
		return "", errors.New("Internal Server Error")
	}
	defer mutex.Unlock()

	// Create transaction record and update user's balance
	transaction := models.Transaction{
		UserID:        user.ID,
		Amount:        amount,
		Type:          "credit",
		TransactionID: transactionID, // Use the generated transaction ID
	}
	result = DbConnection.Create(&transaction)
	if result.Error != nil {
		logger.Error("Error creating transaction: %v\n", zap.Error(result.Error))
		return "", errors.New("Internal Server Error")
	}

	user.Balance += amount
	result = DbConnection.Save(&user)
	if result.Error != nil {
		logger.Error("Error updating user balance: %v\n", zap.Error(result.Error))
		return "", errors.New("Internal Server Error")
	}

	// Invalidate the cache for balance and transaction history
	err := redisClient.Del(ctx, "balance:"+uid).Err()
	if err != nil {
		logger.Error("Error deleting balance cache: %v\n", zap.Error(err))
	}

	// Invalidate the cache for all pages of transaction history
	for i := 1; i <= MaxPages; i++ {
		cacheKey := fmt.Sprintf("transactions:%s:%d:%d", uid, i, TransactionPageSize)
		err = redisClient.Del(ctx, cacheKey).Err()
		if err != nil {
			logger.Error("Error deleting transaction history cache: %v\n", zap.Error(err))
		}
	}

	return "Funds added successfully", nil
}

// GetBalance retrieves the balance for the specified UID.
// Uses Redis cache to speed up subsequent requests.
func GetBalance(ctx *gin.Context, uid string) (float64, error) {
	var user models.User
	redisClient := GetRedisDefaultClient()

	// Check the cache first
	cachedBalance, err := redisClient.Get(ctx, "balance:"+uid).Result()

	if err == nil {
		// If the balance is found in the cache, return it
		return cast.ToFloat64(cachedBalance), nil
	} else if err != redis.Nil {
		// If there was an error retrieving the balance from the cache, log it
		logger.Error("Error getting balance from Redis cache: %v\n", zap.Error(err))
	}

	// If the balance is not in the cache, fetch it from the database
	result := DbConnection.Where("uid = ?", uid).First(&user)
	if result.Error != nil {
		// If the user is not found in the database, return an error
		return 0, errors.New("User not found")
	}

	// Update the cache with the new balance
	err = redisClient.Set(ctx, "balance:"+uid, user.Balance, time.Minute).Err()
	if err != nil {
		// If there was an error updating the cache, log it
		logger.Error("Error setting balance in Redis cache: %v\n", zap.Error(err))
	}

	// Return the balance
	return user.Balance, nil
}

// GetTransactionHistory retrieves the transaction history for the specified UID with pagination.
func GetTransactionHistory(ctx *gin.Context, uid string, page int, limit int) (map[string]interface{}, error) {
	// Calculate the offset
	offset := (page - 1) * limit
	// Find the user with the given UID
	var user models.User
	result := DbConnection.Where("uid = ?", uid).First(&user)
	if result.Error != nil {
		return nil, errors.New("User not found")
	}

	// Build the Redis cache key
	cacheKey := fmt.Sprintf("transactions:%s:%d:%d", uid, page, limit)
	redisClient := GetRedisDefaultClient()

	// Check if the transaction history is already in cache
	cachedTransactions, err := redisClient.Get(ctx, cacheKey).Result()

	// If there is an error other than "key not found", log it
	if err != nil && err != redis.Nil {
		logger.Error("Error getting transactions from Redis cache: %v\n", zap.Error(err))
	}

	var transactions []models.Transaction
	if err == nil {
		// If the transaction history is in cache, unmarshal the JSON data
		err = json.Unmarshal([]byte(cachedTransactions), &transactions)
		if err != nil {
			logger.Error("Error unmarshalling transactions from Redis cache: %v\n", zap.Error(err))
		}
	}

	if len(transactions) == 0 {
		// If the transaction history is not in cache, fetch it from the database
		result = DbConnection.Where("user_id = ?", user.ID).Limit(limit).Offset(offset).Find(&transactions)
		if result.Error != nil {
			return nil, errors.New("Error fetching transactions from the database")
		}

		// Store the transaction history in cache for 10 minutes
		cacheData, err := json.Marshal(transactions)
		if err == nil {
			err = redisClient.Set(ctx, cacheKey, cacheData, 10*time.Minute).Err()
			if err != nil {
				logger.Error("Error setting transactions in Redis cache: %v\n", zap.Error(err))
			}
		} else {
			logger.Error("Error marshalling transactions for Redis cache: %v\n", zap.Error(err))
		}
	}

	// Calculate the total number of transactions for the user
	var count int64
	DbConnection.Model(&models.Transaction{}).Where("user_id = ?", user.ID).Count(&count)

	// Calculate total pages and the next page (if any)
	totalPages := int(count) / limit
	if int(count)%limit > 0 {
		totalPages++
	}

	nextPageID := 0
	if page < totalPages {
		nextPageID = page + 1
	}

	return map[string]interface{}{
		"transactions": transactions,
		"pagination": map[string]interface{}{
			"current_page": page,
			"total_pages":  totalPages,
			"next_page_id": nextPageID,
		},
	}, nil
}
