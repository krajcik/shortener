package handler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"krajcik/shortener/internal/app/shortener"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetShrt(t *testing.T) {
	type args struct {
		flag bool
		post string
	}
	type want struct {
		code     int
		location string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "404",
			args: args{false, "https://google.com"},
			want: want{
				code:     http.StatusNotFound,
				location: "",
			},
		},
		{
			name: "suc",
			args: args{true, "https://google.com"},
			want: want{
				code:     http.StatusTemporaryRedirect,
				location: "https://google.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := shortener.NewService(shortener.NewRepository())
			shrt := ""
			if tt.args.flag {
				shrt, _ = s.ShrtByURL(tt.args.post)
			}

			request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", shrt), nil)
			w := httptest.NewRecorder()
			GetShrt(s)(w, request)

			defer w.Result().Body.Close()
			res := w.Result()
			// проверяем код ответа
			require.Equal(t, tt.want.code, res.StatusCode)

			if tt.want.location != "" {
				assert.Equal(t, tt.want.location, res.Header.Get("Location"))
			}
		})
	}
}
