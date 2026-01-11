package auth

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func TestValidHeader(t *testing.T) {
	tests := map[string]struct {
		header http.Header
		want   string
	}{
		"validHeader": {
			header: http.Header{"Authorization": {"Bearer 123456"}},
			want:   "123456",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := GetBearerToken(tc.header)
			if err != nil {
				t.Fatalf("GetBearerToken error: %v", err)
			}

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestInvalidHeader(t *testing.T) {
	tests := map[string]struct {
		header http.Header
		want   error
	}{
		"missingHeader": {
			header: http.Header{},
			want:   ErrMissingHeader,
		},
		"incorrectHeader": {
			header: http.Header{"Authorization": {"Some other auth"}},
			want:   ErrInvalidHeader,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, gotErr := GetBearerToken(tc.header)

			if !errors.Is(gotErr, tc.want) {
				t.Fatalf("expected: %v, got: %v", tc.want, gotErr)
			}
		})
	}
}
