package controllers

import (
	"github.com/gin-gonic/gin"
	"ledger-service/models"
	"ledger-service/services"
	"net/http"
	"strconv"
)

// AddFunds adds funds to a user's account.
// @Summary Add funds to a user's account
// @Description Add funds to a user's account by the given UID and amount
// @Tags Users
// @Accept json
// @Produce json
// @Param uid path string true "User ID"
// @Param AddFundsRequest body models.AddFundsRequest true "Amount to add"
// @Success 201 {object} models.Response
// @Failure 400 Bad Request models.Response
// @Failure 409 Conflict models.Response
// @Failure 500 Internal Server Error models.Response
// @Router /users/{uid}/add [post]
func AddFunds(c *gin.Context) {
	uid := c.Param("uid")
	var request models.AddFundsRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := services.AddFunds(c, uid, request)
	if err != nil {
		if err.Error() == "Amount must be positive" {
			models.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		} else if err.Error() == "Transaction already processed" {
			models.SendErrorResponse(c, http.StatusConflict, err.Error())
		} else {
			models.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Return success response
	models.SendResponseData(c, gin.H{
		"Id":      uid,
		"Message": message,
	})
}

// GetBalance retrieves the balance of a user.
// @Summary Get user's balance
// @Description Get the balance of a user by the given UID
// @Tags Users
// @Accept json
// @Produce json
// @Param uid path string true "User ID"
// @Success 200 {object} models.Response
// @Failure 404 Not Found models.Response
// @Router /users/{uid}/balance [get]
func GetBalance(c *gin.Context) {
	uid := c.Param("uid")

	// Call the GetBalance service function
	balance, err := services.GetBalance(c, uid)
	if err != nil {
		// If error occurs, send error response
		models.SendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	// If successful, send balance data in the response
	models.SendResponseData(c, gin.H{"Balance": balance})
}

// GetTransactionHistory retrieves the transaction history of a user.
// @Summary Get user's transaction history
// @Description Get the transaction history of a user by the given UID with pagination
// @Tags Users
// @Accept json
// @Produce json
// @Param uid path string true "User ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Success 200 {object} models.Response
// @Failure 404 Not Found models.Response
// @Router /users/{uid}/transactions [get]
func GetTransactionHistory(c *gin.Context) {
	uid := c.Param("uid") // Get user ID from the path parameter

	// Get page and limit query parameters, with default values if not provided
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Call the GetTransactionHistory service function with the uid, page, and limit
	transactions, err := services.GetTransactionHistory(c, uid, page, limit)

	// If an error occurs, send an error response
	if err != nil {
		models.SendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	// If successful, send the result with transactions and pagination information
	models.SendResponseData(c, gin.H{"List of Transactions": transactions})
}
