// middleware/auth_middleware.go
package middleware

import (
	"context"
	"net/http"
	"rentjoy/pkg/auth"
)

// 用來定義上下文的 key
type UserContext string

// 定義上下文的 key 常數
const (
	UserIDKey    UserContext = "user_id"
	UserEmailKey UserContext = "user_email"
)

// AuthMiddleware 驗證 JWT Token 的中間件
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 從 Cookie 獲取 token
		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" {
			// 將當前 URL 加入重定向參數，以便登入後返回
			returnURL := r.URL.Path
			http.Redirect(w, r, "/Login?returnUrl="+returnURL, http.StatusSeeOther)
			return
		}

		// 驗證 token
		claims, err := auth.ValidateToken(cookie.Value)
		if err != nil {
			// token 無效也重定向到登入頁面
			http.Redirect(w, r, "/Login", http.StatusSeeOther)
			return
		}

		// 將資訊存入 Context（這部分不變）
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)

		// token 有效，繼續處理請求
		next(w, r.WithContext(ctx))
	}
}

// 從上下文中獲取用戶 ID
func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(UserIDKey).(uint)
	return userID, ok
}

// 從上下文中獲取用戶 Email
func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(UserEmailKey).(string)
	return email, ok
}
