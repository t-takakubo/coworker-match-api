package common

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

// ExtractPathParam パスパラメータから指定されたキーの値を抽出し、バリデーションを行う
func ExtractPathParam(r *http.Request, paramName string) (string, error) {
	vars := mux.Vars(r)
	paramValue := vars[paramName]
	if paramValue == "" {
		return "", errors.New("missing " + paramName)
	}
	return paramValue, nil
}

// ExtractUserIdFromContext コンテキストからユーザーIDを抽出し、バリデーションを行う
func ExtractUserIdFromContext(r *http.Request) (string, error) {
        key := UserIDKey
	userId, ok := r.Context().Value(key).(string)
	if !ok {
		return "", errors.New("failed to get userId from context")
	}
	return userId, nil
}
