package handler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
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
		code        int
		contentType string
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
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "suc",
			args: args{true, "https://google.com"},
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := shortener.NewService(shortener.NewRepository())
			shrt := ""
			if tt.args.flag {
				shrt, _ = s.ShrtByUrl(tt.args.post)

			}

			request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", shrt), nil)
			w := httptest.NewRecorder()
			GetShrt(s)(w, request)

			res := w.Result()
			// проверяем код ответа
			require.Equal(t, tt.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			if res.StatusCode < 300 {
				assert.Equal(t, tt.args.post, string(resBody))
			}
		})
	}
}
