package handlers

import (
	"appliedTo/internal/app/user"
	"appliedTo/internal/platform/db"
	"appliedTo/internal/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
    var userDto user.UserCreateDto
    if err := c.ShouldBindJSON(&userDto); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
        return
    }

    var user user.User
    if err := db.Where("email = ?", userDto.Email).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userDto.Password))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
    }

    token, err := utils.GenerateJWT(int(user.ID))
    if err != nil {
        log.Println(err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}


func Register(c *gin.Context) {
    var userDto user.UserCreateDto
    if err := c.ShouldBindJSON(&userDto); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    var exsistingUser user.User
    if err := db.Where("email = ?", userDto.Email).First(&exsistingUser).Error; err == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "E-Mail already in use"})
        return
    }

    hashedPassword, err := utils.HashPassword(userDto.Password) 
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
        return
    }

    user := user.User{
        FirstName: userDto.FirstName,
        LastName: userDto.LastName,
        Email: userDto.Email,
        Password: hashedPassword,
        Created: time.Now(),
    }

    if err := db.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
    
}

