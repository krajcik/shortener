package handler

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"krajcik/shortener/internal/app/shortener"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostShrt(t *testing.T) {
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
			name: "suc",
			args: args{true, "https://google.com"},
			want: want{
				code:        http.StatusCreated,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := shortener.NewService(shortener.NewRepository())
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.args.post))
			w := httptest.NewRecorder()
			PostShrt(s)(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			assert.Len(t, string(resBody), len("http://")+shortener.ShortLen+len(request.Host)+1)
			assert.NotEqual(t, tt.args.post, string(resBody))
		})
	}
}
