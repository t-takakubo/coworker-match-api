package common

import (
	"context"
	"github.com/gorilla/mux"
	"net/http/httptest"
	"testing"
)

func TestExtractPathParam(t *testing.T) {
	tests := []struct {
		name      string
		vars      map[string]string
		paramName string
		want      string
		wantErr   bool
	}{
		{
			name:      "param present",
			vars:      map[string]string{"id": "123"},
			paramName: "id",
			want:      "123",
			wantErr:   false,
		},
		{
			name:      "param missing",
			vars:      map[string]string{},
			paramName: "id",
			want:      "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			r = mux.SetURLVars(r, tt.vars)
			got, err := ExtractPathParam(r, tt.paramName)
			if (err != nil) != tt.wantErr {
				t.Fatalf("wantErr %v got %v", tt.wantErr, err)
			}
			if got != tt.want {
				t.Fatalf("expected %s got %s", tt.want, got)
			}
		})
	}
}

func TestExtractUserIdFromContext(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		want    string
		wantErr bool
	}{
		{
			name:    "userId present",
			ctx:     context.WithValue(context.Background(), UserIdKey, "user123"),
			want:    "user123",
			wantErr: false,
		},
		{
			name:    "userId missing",
			ctx:     context.Background(),
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			r = r.WithContext(tt.ctx)
			got, err := ExtractUserIdFromContext(r)
			if (err != nil) != tt.wantErr {
				t.Fatalf("wantErr %v got %v", tt.wantErr, err)
			}
			if got != tt.want {
				t.Fatalf("expected %s got %s", tt.want, got)
			}
		})
	}
}
