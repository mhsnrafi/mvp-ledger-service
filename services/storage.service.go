package services

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"ledger-service/logger"
	"ledger-service/models"
	"sync"
	"time"
)

// Global variable to hold the database connection
var DbConnection *gorm.DB

// Constants to set the number of retries and delay between retries
const (
	retries = 5
	delay   = 5 * time.Second
)

// ConnectDB function attempts to connect to the database with the given credentials
func ConnectDB() {
	var err error
	var db *gorm.DB

	// Create the DSN (Data Source Name) for the database connection
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		Config.DBUserName, Config.DBUserPassword, Config.DBHost, Config.DBPort, Config.DBName)

	// Attempt to connect to the database with retries
	for i := 0; i < retries; i++ {
		db, err = gorm.Open("postgres", dsn)
		if err == nil {
			DbConnection = db
			// AutoMigrate the tables for the models
			DbConnection.AutoMigrate(&models.User{})
			DbConnection.AutoMigrate(&models.Transaction{})
			DbConnection.AutoMigrate(&models.Token{})
			logger.Info("Successfully connected to the Database")
			return
		}
		// Log an error if the connection fails
		logger.Error("Failed to connect to the Database. Retrying...", zap.Error(err), zap.Int("attempt", i+1))
		time.Sleep(delay)
	}

	// Log a fatal error if the connection fails after all retries
	logger.Fatal("Failed to connect to the Database after multiple attempts", zap.Error(err))
}

// Global variables to store the Redis client and a sync.Once object
var redisDefaultClient *redis.Client
var redisDefaultOnce sync.Once

// GetRedisDefaultClient returns the Redis client, initializing it if necessary
func GetRedisDefaultClient() *redis.Client {
	redisDefaultOnce.Do(func() {
		redisDefaultClient = redis.NewClient(&redis.Options{
			Addr:     Config.RedisDefaultAddr,
			Password: Config.RedisPassword,
		})
	})

	return redisDefaultClient
}

// CheckRedisConnection checks if the connection to Redis is working
func CheckRedisConnection() {
	redisClient := GetRedisDefaultClient()
	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	logger.Info("Connected to Redis!")
}
