package controllers

import (
	"net/http"

	"appliedTo/dtos/user_dtos"
	mappers "appliedTo/mappers/user_mappers"
	"appliedTo/middleware"
	"appliedTo/models"
	"appliedTo/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Create a new user
// @Description Creates a new user with the provided name and email. Ensures the email is unique.
// @Tags user
// @Accept  json
// @Produce  json
// @Param   user  body  userdtos.UserCreateDto  true  "User data"
// @Success 200 {object} map[string]interface{} "User created successfully"
// @Failure 400 {object} map[string]string "Invalid input or email already in use"
// @Failure 500 {object} map[string]string "Could not create user"
// @Router /user [post]
func CreateUser(c *gin.Context) {
	var userDto userdtos.UserCreateDto

	if err := c.ShouldBindJSON(&userDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if !utils.ValidateRequiredFields(c, []utils.RequiredField{
		{Value: userDto.FirstName, Name: "Firstname"},
		{Value: userDto.LastName, Name: "Lastname"},
		{Value: userDto.Email, Name: "Email"},
		{Value: userDto.Password, Name: "Password"},
	}) {
		return
	}

	if !(utils.IsValidEmail(userDto.Email)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is invalid."})
		return
	}

	var exsistingUser models.User
	if err := db.Where("email = ?", userDto.Email).First(&exsistingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "E-Mail is already in use."})
		return
	}

	user := mappers.CreateModel(userDto)

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	response := mappers.MapModelToPublicDto(user)

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": response})
}

// @Summary Get a user by ID
// @Description Get detailed information about a user
// @Tags user
// @Accept  json
// @Produce  json
// @Param   id  path  int  true  "User ID"
// @Success 200 {object} userdtos.UserPublicDto "Successfully retrieved user"
// @Failure 400 "Invalid ID format"
// @Failure 404 "User not found"
// @Failure 404 "Database query failed"
// @Router /user/{id} [get]
func GetUser(c *gin.Context) {
	id := c.GetUint(middleware.CtxKeyUserID)

	var user models.User

	if err := db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Databse query failed"})
		return
	}

	response := mappers.MapModelToPublicDto(user)

	c.JSON(http.StatusOK, gin.H{"user": response})
}

// @Summary      Partially update a user
// @Description  Updates only the provided fields on the user with the given ID.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id       path     int   true  "User ID" example(123)
// @Param        payload  body     userdtos.UserPatchDto true  "Fields to patch"
// @Success      200      {object} userdtos.UserPublicDto
// @Failure      400      "Invalid request payload or invalid field values"
// @Failure      404      "User not found"
// @Failure      409      "Email already in use"
// @Failure      500      "Could not update user"
// @Router       /user/{id} [patch]
func PatchUser(c *gin.Context) {
	var dto userdtos.UserPatchDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	id := c.GetUint(middleware.CtxKeyUserID)

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if dto.Email != nil {
		if !utils.IsValidEmail(*dto.Email) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address"})
			return
		}
		if *dto.Email != user.Email {
			var count int64
			if err := db.Model(&models.User{}).Where("email = ?", *dto.Email).Count(&count).Error; err == nil && count > 0 {
				c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
				return
			}
		}
	}

	mappers.PatchModel(&user, dto)

	if err := db.Save(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user"})
		return
	}

	response := mappers.MapModelToPublicDto(user)

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": response})
}

// @Summary Update a user (full replace)
// @Description Modify user user data by providing new data and the user-ID.
// @Tags user
// @Accept  json
// @Produce  json
// @Param   id    path  int   true  "User ID" example(123)
// @Param   user  body  userdtos.UserCreateDto  true  "User data"
// @Success 200 {object} userdtos.UserPublicDto "User successfully modified."
// @Failure 400 "Invalid ID format"
// @Failure 404 "User not found"
// @Failure 404 "Database query failed"
// @Router /user/{id} [put]
func UpdateUser(c *gin.Context) {
	var userDto userdtos.UserCreateDto

	if err := c.ShouldBindJSON(&userDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	id := c.GetUint(middleware.CtxKeyUserID)

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if !utils.ValidateRequiredFields(c, []utils.RequiredField{
		{Value: userDto.FirstName, Name: "Firstname"},
		{Value: userDto.LastName, Name: "Lastname"},
		{Value: userDto.Email, Name: "Email"},
		{Value: userDto.Password, Name: "Password"},
	}) {
		return
	}

	if !utils.IsValidEmail(userDto.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email adress"})
		return
	}

	user.FirstName = userDto.FirstName
	user.LastName = userDto.LastName
	user.Email = userDto.Email
	user.Password = userDto.Password

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user"})
		return
	}

	response := mappers.MapModelToPublicDto(user)

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": response})
}

// @Summary Delete a user.
// @Description Remove the user from the database by providing the user-ID.
// @Tags user
// @Accept  json
// @Produce  json
// @Param  id  path int true "User ID" example(123)
// @Success 200 "User deleted modified."
// @Failure 400 "Invalid ID format"
// @Failure 500 "Could not delete user"
// @Router /user/{id} [delete]
func DeleteUser(c *gin.Context) {
	id := c.GetUint(middleware.CtxKeyUserID)

	if err := db.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete User"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
