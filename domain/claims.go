package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const HMAC_SAMPLE_SECRET = "hmacSampleSecret"
const ACCESS_TOKEN_DURATION = time.Hour
const REFRESH_TOKEN_DURATION = time.Hour * 24 * 30

type RefreshTokenClaims struct {
	TokenType  string   `json:"token_type"`
	CustomerId string   `json:"cid"`
	Accounts   []string `json:"accounts"`
	Username   string   `json:"un"`
	Role       string   `json:"role"`
	jwt.RegisteredClaims
}

type AccessTokenClaims struct {
	CustomerId string   `json:"customer_id"`
	Accounts   []string `json:"accounts"`
	Username   string   `json:"username"`
	Role       string   `json:"role"`
	jwt.RegisteredClaims
}

func (c AccessTokenClaims) IsUserRole() bool {
	return c.Role == "user"
}

func (c AccessTokenClaims) IsValidCustomerId(customerId string) bool {
	return c.CustomerId == customerId
}

func (c AccessTokenClaims) IsValidAccountId(accountId string) bool {
	if accountId != "" {
		accountFound := false
		for _, a := range c.Accounts {
			if a == accountId {
				accountFound = true
				break
			}
		}
		return accountFound
	}
	return true
}

func (c AccessTokenClaims) IsRequestVerifiedWithTokenClaims(urlParams map[string]string) bool {
	if c.CustomerId != urlParams["customer_id"] {
		return false
	}

	if !c.IsValidAccountId(urlParams["account_id"]) {
		return false
	}
	return true
}

func (c AccessTokenClaims) RefreshTokenClaims() RefreshTokenClaims {
	return RefreshTokenClaims{
		TokenType:  "refresh_token",
		CustomerId: c.CustomerId,
		Accounts:   c.Accounts,
		Username:   c.Username,
		Role:       c.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(REFRESH_TOKEN_DURATION)),
		},
	}
}

func (c RefreshTokenClaims) AccessTokenClaims() AccessTokenClaims {
	return AccessTokenClaims{
		CustomerId: c.CustomerId,
		Accounts:   c.Accounts,
		Username:   c.Username,
		Role:       c.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ACCESS_TOKEN_DURATION)),
		},
	}
}
