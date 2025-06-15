package common

import (
	"context"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

// UserIDKey is the key used for storing the user ID in the context.
const UserIDKey contextKey = "userId"

// コンテキストからユーザーIDを取得するヘルパー関数
func GetUserID(ctx context.Context) (string, bool) {
        key := UserIDKey
	userID, ok := ctx.Value(key).(string)
	return userID, ok
}
