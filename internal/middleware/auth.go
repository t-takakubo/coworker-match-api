package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/coworker-match-api/internal/common"
	"google.golang.org/api/idtoken"
)

func Auth(next http.Handler) http.Handler {
	GOOGLE_CLIENT_ID := os.Getenv("GOOGLE_CLIENT_ID")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authorizationヘッダーの取得
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			common.RespondWithError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		// IDトークンの取得
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			common.RespondWithError(w, http.StatusUnauthorized, "Invalid Authorization header format")
			return
		}
		idToken := parts[1]

		// トークンの検証
		ctx := context.Background()
		payload, err := idtoken.Validate(ctx, idToken, GOOGLE_CLIENT_ID)
		if err != nil {
			common.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		// トークンのペイロードから必要な情報を取得（例：ユーザーID）
		userId, ok := payload.Claims["sub"].(string)
		if !ok {
			common.RespondWithError(w, http.StatusUnauthorized, "Invalid token payload")
			return
		}

                // リクエストのコンテキストにユーザー情報を追加
                key := common.UserIDKey
		ctx = context.WithValue(r.Context(), key, userId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
