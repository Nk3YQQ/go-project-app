package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

var jwtKey = []byte("SecretKey")

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func ParseToken(authToken string, parsedToken *string) error {
	if len(authToken) <= 7 {
		return errors.New("there no any authorization token")
	} else {
		if authToken[:7] != "Bearer " {
			return errors.New("authorization header format must be Bearer {token}")
		}
	}

	*parsedToken = authToken[7:]

	return nil
}

func GenerateToken(token *string, userID uint) error {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	generatedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := generatedToken.SignedString(jwtKey)

	if err != nil {
		return err
	}

	*token = tokenString

	return nil
}

func Authenticate(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")

	claims := &Claims{}

	var parsedToken string

	if err := ParseToken(tokenString, &parsedToken); err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}

	token, err := jwt.ParseWithClaims(parsedToken, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected sign method")
		}
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	c.Set("userID", claims.UserID)

	c.Next()
}
