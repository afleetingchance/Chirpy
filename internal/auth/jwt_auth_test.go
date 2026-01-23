package auth

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestValidJWT(t *testing.T) {
	userId := uuid.New()

	tests := map[string]struct {
		userID      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
		want        uuid.UUID
	}{
		"canMake": {
			userID:      userId,
			tokenSecret: "123456",
			expiresIn:   10 * time.Minute,
			want:        userId,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenString, err := MakeJWT(tc.userID, tc.tokenSecret)
			if err != nil {
				t.Fatalf("MakeJWT error: %v", err)
			}

			got, err := ValidateJWT(tokenString, tc.tokenSecret)
			if err != nil {
				t.Fatalf("ValidateJWT error: %v", err)
			}

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestInvalidJWT(t *testing.T) {
	userId := uuid.New()

	tests := map[string]struct {
		userID              uuid.UUID
		createTokenSecret   string
		retrieveTokenSecret string
		expiresIn           time.Duration
		want                error
	}{
		"BadRetrieveSecret": {
			userID:              userId,
			createTokenSecret:   "123456",
			retrieveTokenSecret: "try this",
			expiresIn:           10 * time.Minute,
			want:                jwt.ErrSignatureInvalid,
		},
		"ExpiredToken": {
			userID:              userId,
			createTokenSecret:   "123456",
			retrieveTokenSecret: "123456",
			expiresIn:           time.Nanosecond,
			want:                jwt.ErrTokenExpired,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenString, err := MakeJWT(tc.userID, tc.createTokenSecret)
			if err != nil {
				t.Fatalf("MakeJWT error: %v", err)
			}

			_, gotErr := ValidateJWT(tokenString, tc.retrieveTokenSecret)

			if !errors.Is(gotErr, tc.want) {
				t.Fatalf("expected: %v, got: %v", tc.want, gotErr)
			}
		})
	}
}
