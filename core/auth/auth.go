package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"go.uber.org/zap"
)

type UserInfo struct {
	UserID   string
	Username string
	Email    string
}

type CognitoAuthenticator struct {
	userPoolID string
	clientID   string
	region     string
	logger     *zap.Logger
	jwks       jwk.Set
	lastUpdate time.Time
	mu         sync.RWMutex
}

func NewCognitoAuthenticator(userPoolID, clientID, region string, logger *zap.Logger) *CognitoAuthenticator {
	return &CognitoAuthenticator{
		userPoolID: userPoolID,
		clientID:   clientID,
		region:     region,
		logger:     logger,
	}
}

func (a *CognitoAuthenticator) ValidateToken(token string) (*UserInfo, error) {
	// Refresh JWKS if needed
	if err := a.refreshJWKSIfNeeded(); err != nil {
		return nil, fmt.Errorf("failed to refresh JWKS: %v", err)
	}

	// Parse and validate token
	if a.jwks == nil {
		return nil, errors.New("JWKS not initialized")
	}
	parsedToken, err := jwt.ParseString(
		token,
		jwt.WithKeySet(a.jwks),
		jwt.WithIssuer(fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", a.region, a.userPoolID)),
		jwt.WithAudience(a.clientID),
		jwt.WithClaimValue("token_use", "id"),
	)
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %v", err)
	}

	// Check expiration
	if parsedToken.Expiration().Before(time.Now()) {
		return nil, fmt.Errorf("token expired")
	}

	// Extract user info
	userInfo := &UserInfo{}

	if customIDClaim, ok := parsedToken.Get("custom:cognito_user_id"); ok {
		if customID, ok := customIDClaim.(string); ok && customID != "" {
			userInfo.UserID = customID
		}
	}

	// Fallback to sub if custom:cognito_user_id is missing or empty
	if userInfo.UserID == "" {
		userInfo.UserID = parsedToken.Subject() // sub
	}
	// Safely extract username claim
	if usernameClaim, ok := parsedToken.Get("name"); ok {
		if username, ok := usernameClaim.(string); ok {
			userInfo.Username = username
		}
	}

	// Safely extract email claim
	if emailClaim, ok := parsedToken.Get("email"); ok {
		if email, ok := emailClaim.(string); ok {
			userInfo.Email = email
		}
	}

	return userInfo, nil
}

func (a *CognitoAuthenticator) refreshJWKSIfNeeded() error {
	a.mu.RLock()
	needsRefresh := a.jwks == nil || time.Since(a.lastUpdate) > time.Hour
	a.mu.RUnlock()

	if !needsRefresh {
		return nil
	}

	return a.refreshJWKS()
}

func (a *CognitoAuthenticator) refreshJWKS() error {
	jwksURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", a.region, a.userPoolID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jwks, err := jwk.Fetch(ctx, jwksURL)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %v", err)
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	a.jwks = jwks
	a.lastUpdate = time.Now()
	return nil
}

func ExtractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is required")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("authorization header format must be Bearer {token}")
	}

	return parts[1], nil
}

func SetUserInfoInGinContext(c *gin.Context, userInfo *UserInfo) {
	c.Set("userInfo", userInfo)

}

func GetUserInfoFromGinContext(c *gin.Context) (*UserInfo, bool) {
	userInfo, exists := c.Get("userInfo")
	if !exists {
		return nil, false
	}
	return userInfo.(*UserInfo), true
}

func (a *CognitoAuthenticator) Logger() *zap.Logger {
	return a.logger
}
