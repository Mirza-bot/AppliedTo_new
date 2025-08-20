package authapi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	auth "appliedTo/internal/app/auth"
	"appliedTo/internal/app/user"
)

type Handlers struct {
	Svc *auth.Service
}

func NewHandlers(s *auth.Service) *Handlers { return &Handlers{Svc: s} }

// Login godoc
// @Summary      Login
// @Description  Authenticate with email & password and receive a JWT.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload body  auth.LoginRequest true "Credentials"
// @Success      200     {object} auth.TokenResponse
// @Failure      401     {object} map[string]string "Invalid email or password"
// @Failure      500     {object} map[string]string "Could not generate token"
// @Router       /auth/login [post]
func (h *Handlers) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}
	tok, err := h.Svc.Authenticate(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	c.JSON(http.StatusOK, auth.TokenResponse{
		AccessToken: tok,
        ExpiresIn:   int64(h.Svc.JWT().AccessTTL.Seconds()),
	})
}

// Register godoc
// @Summary      Register
// @Description  Create a new user account and receive a JWT.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload body  auth.RegisterRequest true "New user"
// @Success      200     {object} auth.TokenResponse "User registered"
// @Failure      400     {object} map[string]string "Invalid input"
// @Failure      409     {object} map[string]string "E-Mail already in use"
// @Failure      500     {object} map[string]string "Could not create user"
// @Router       /auth/register [post]
func (h *Handlers) Register(c *gin.Context) {
	var req auth.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	tok, err := h.Svc.Register(c.Request.Context(), req)
	if err != nil {
		switch err {
		case user.ErrEmailInUse:
			c.JSON(http.StatusConflict, gin.H{"error": "E-Mail already in use"})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, auth.TokenResponse{
		AccessToken: tok,
        ExpiresIn:   int64(h.Svc.JWT().AccessTTL.Seconds()),
	})
}
