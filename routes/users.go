package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"tallyRead.com/models"
	"tallyRead.com/utils"
)

func RegisterUser(context *gin.Context) {
	var user models.User

	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := user.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "user created successfully", "user": user})
}

func Login(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = user.ValidateUser()
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Login Successfully", "token": token, "user": user})
}

func ForgotPassword(context *gin.Context) {
	var req models.ForgotPassword

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Email", "error": err.Error()})
		return
	}

	user, err := models.FindByEmail(req.Email)

	if err != nil || user == nil {
		context.JSON(http.StatusOK, gin.H{"message": "If an account with that email exists, a reset Link has been sent."})
		return
	}

	token, hashedToken := utils.GenerateResetToken()
	expires := time.Now().Add(5 * time.Minute)

	user.ResetToken = hashedToken
	user.ResetTokenExpires = expires

	if err := user.UpdateResetToken(); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resetURL := utils.GenerateResetURL(token)

	if err := utils.SendResetEmail(user.Email, resetURL); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "If an account with that email exists, a reset Link has been sent."})
}

func ResetPassword(context *gin.Context) {
	token := context.DefaultQuery("token", "nil")

	if token == "nil" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Link"})
		return
	}

	var req models.ResetRequest

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedToken := utils.GetHashedToken(token)

	user, err := models.FindByResetToken(hashedToken)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.ResetTokenExpires.Before(time.Now()) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Expired Token"})
		return
	}

	user.Password = req.Password
	user.ResetToken = ""
	user.ResetTokenExpires = time.Time{}

	err = user.UpdatePassword()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Password Updated"})

}
