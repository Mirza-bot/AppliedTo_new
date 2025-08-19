package handlers

import (
	"errors"
	"net/http"

	userdtos "appliedTo/dtos/user_dtos"
	userservice "appliedTo/internal/services/user_service"
	"appliedTo/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ---- Response types used in Swagger ----

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type UserResponse struct {
	User userdtos.UserPublicDto `json:"user"`
}

type MessageUserResponse struct {
	Message string                 `json:"message"`
	User    userdtos.UserPublicDto `json:"user"`
}

type UserHandlers struct {
	Svc *userservice.UserService
}

func NewUserHandlers(svc *userservice.UserService) *UserHandlers {
	return &UserHandlers{Svc: svc}
}

// CreateUser godoc
// @Summary      Create a new user
// @Description  Creates a new user. Email must be unique; password is hashed.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user  body      userdtos.UserCreateDto  true  "User data"
// @Success      200   {object}  MessageUserResponse     "User created successfully"
// @Failure      400   {object}  ErrorResponse           "Invalid input"
// @Failure      409   {object}  ErrorResponse           "Email already in use"
// @Failure      500   {object}  ErrorResponse           "Could not create user"
// @Router       /user [post]
func (h *UserHandlers) CreateUser(c *gin.Context) {
	var dto userdtos.UserCreateDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	resp, err := h.Svc.Create(c.Request.Context(), dto)
	if err != nil {
		switch {
		case errors.Is(err, userservice.ErrInvalidEmail):
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Email is invalid"})
		case errors.Is(err, userservice.ErrEmailInUse):
			c.JSON(http.StatusConflict, ErrorResponse{Error: "E-Mail is already in use"})
		default:
			// validate.Required returns messages like "a firstname is required"
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, MessageUserResponse{
		Message: "User created successfully",
		User:    resp,
	})
}

// GetUser godoc
// @Summary      Get a user by ID
// @Description  Returns the user for the given ID.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"  example(123)
// @Success      200  {object}  UserResponse           "Successfully retrieved user"
// @Failure      400  {object}  ErrorResponse          "Invalid ID"
// @Failure      404  {object}  ErrorResponse          "User not found"
// @Failure      500  {object}  ErrorResponse          "Database query failed"
// @Router       /user/{id} [get]
func (h *UserHandlers) GetUser(c *gin.Context) {
	id := c.GetUint(middleware.CtxKeyUserID)

	resp, err := h.Svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database query failed"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{User: resp})
}

// UpdateUser godoc
// @Summary      Update a user (full replace)
// @Description  Replaces all user fields with the provided payload.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id    path      int                     true  "User ID"  example(123)
// @Param        user  body      userdtos.UserCreateDto  true  "User data"
// @Success      200   {object}  MessageUserResponse     "User successfully modified."
// @Failure      400   {object}  ErrorResponse           "Invalid input"
// @Failure      404   {object}  ErrorResponse           "User not found"
// @Failure      409   {object}  ErrorResponse           "Email already in use"
// @Failure      500   {object}  ErrorResponse           "Could not update user"
// @Router       /user/{id} [put]
func (h *UserHandlers) UpdateUser(c *gin.Context) {
	id := c.GetUint(middleware.CtxKeyUserID)

	var dto userdtos.UserCreateDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	resp, err := h.Svc.Update(c.Request.Context(), id, dto)
	if err != nil {
		switch {
		case errors.Is(err, userservice.ErrInvalidEmail):
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid email address"})
		case errors.Is(err, userservice.ErrEmailInUse):
			c.JSON(http.StatusConflict, ErrorResponse{Error: "Email already in use"})
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
		default:
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, MessageUserResponse{
		Message: "User updated successfully",
		User:    resp,
	})
}

// PatchUser godoc
// @Summary      Partially update a user
// @Description  Updates only the provided fields on the user with the given ID.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id       path      int                    true  "User ID"  example(123)
// @Param        payload  body      userdtos.UserPatchDto  true  "Fields to patch"
// @Success      200      {object}  MessageUserResponse    "User updated successfully"
// @Failure      400      {object}  ErrorResponse          "Invalid request payload or invalid field values"
// @Failure      404      {object}  ErrorResponse          "User not found"
// @Failure      409      {object}  ErrorResponse          "Email already in use"
// @Failure      500      {object}  ErrorResponse          "Could not update user"
// @Router       /user/{id} [patch]
func (h *UserHandlers) PatchUser(c *gin.Context) {
	id := c.GetUint(middleware.CtxKeyUserID)

	var dto userdtos.UserPatchDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	resp, err := h.Svc.Patch(c.Request.Context(), id, dto)
	if err != nil {
		switch {
		case errors.Is(err, userservice.ErrInvalidEmail):
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid email address"})
		case errors.Is(err, userservice.ErrEmailInUse):
			c.JSON(http.StatusConflict, ErrorResponse{Error: "Email already in use"})
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not update user"})
		}
		return
	}

	c.JSON(http.StatusOK, MessageUserResponse{
		Message: "User updated successfully",
		User:    resp,
	})
}

// DeleteUser godoc
// @Summary      Delete a user
// @Description  Removes the user with the given ID.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"  example(123)
// @Success      200  {object}  MessageResponse        "User deleted successfully"
// @Failure      404  {object}  ErrorResponse          "User not found"
// @Failure      500  {object}  ErrorResponse          "Could not delete user"
// @Router       /user/{id} [delete]
func (h *UserHandlers) DeleteUser(c *gin.Context) {
	id := c.GetUint(middleware.CtxKeyUserID)

	if err := h.Svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not delete user"})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "User deleted successfully"})
}
