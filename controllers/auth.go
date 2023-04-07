package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"ledger-service/models"
	"ledger-service/services"
	"net/http"
)

// GenerateAccessToken generates new access tokens.
// @Summary Generate new access tokens.
// @Description Generate new access tokens for the provided email.
// @Tags Tokens
// @Accept  json
// @Produce  json
// @Param authReq body models.AuthRequest true "Auth Request"
// @Success 200 {object} models.Response
// @Success 400 {object} models.Response
// @Router /access [post]
func GenerateAccessToken(c *gin.Context) {
	var requestBody models.AuthRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// generate new access tokens
	accessToken, refreshToken, err := services.GenerateAccessTokens(requestBody.Email)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{
		"token": gin.H{
			"access":  accessToken.GetResponseJson(),
			"refresh": refreshToken.GetResponseJson()},
	}
	response.SendResponse(c)
}

// Refresh handles the request for token refresh.
// @Summary Handle the request for token refresh.
// @Description Handle the request for token refresh by validating the refresh token, generating new access and refresh tokens and returning the response.
// @Tags Tokens
// @Accept  json
// @Produce  json
// @Param requestBody body models.RefreshRequest true "Refresh Request"
// @Success 200 {object} models.Response
// @Success 400 {object} models.Response
// @Router /refresh [post]
func Refresh(c *gin.Context) {
	var requestBody models.RefreshRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// check token validity
	token, err := services.VerifyToken(requestBody.Token, models.TokenTypeRefresh)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	// delete old token
	err = services.DeleteTokenById(token.ID)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	accessToken, refreshToken, err := services.GenerateAccessTokens(requestBody.Email)
	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{
		"Email": requestBody.Email,
		"token": gin.H{
			"access":  accessToken.GetResponseJson(),
			"refresh": refreshToken.GetResponseJson()},
	}
	response.SendResponse(c)
}
