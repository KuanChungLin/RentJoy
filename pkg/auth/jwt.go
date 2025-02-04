package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaim struct {
	UserID               uint   `json:"user_id"`
	Email                string `json:"email"`
	jwt.RegisteredClaims        // 內建的標準 JWT Claims
}

// JWT 加密用的金鑰
var jwtKey = []byte("SignKey")

// 產生 JWT Token
func GenerateToken(userID uint, email string, isRemember bool) (string, error) {
	var tokenTime time.Time
	if isRemember {
		tokenTime = time.Now().Add(24 * time.Hour)
	} else {
		tokenTime = time.Now().Add(1 * time.Hour)
	}
	// 建立 JWT Claims
	claims := &JWTClaim{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			// Token 過期時間設為 1 小時
			ExpiresAt: jwt.NewNumericDate(tokenTime),
			// Token 的簽發時間
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	// 使用 HS256 演算法建立 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用密鑰對 Token 進行簽名並返回
	return token.SignedString(jwtKey)
}

// 驗證 JWT Token
func ValidateToken(tokenString string) (*JWTClaim, error) {
	// 解析 Token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	// 驗證 Token 的有效性並轉換 Claims
	if claims, ok := token.Claims.(*JWTClaim); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
