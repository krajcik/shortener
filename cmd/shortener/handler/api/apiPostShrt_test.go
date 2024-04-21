package api

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
	"krajcik/shortener/cmd/shortener/config"
	"krajcik/shortener/internal/app/shortener"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostShrtHandler_PostShrt(t *testing.T) {
	s := shortener.NewService(shortener.NewRepository())
	p := &config.Params{}
	l := zap.NewNop()
	//type fields struct {
	//	S *shortener.Service
	//	P *config.Params
	//	L *zap.Logger
	//}
	type args struct {
		post string
	}
	type want struct {
		code        int
		contentType string
		err         bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "suc",
			args: args{`{"url":"https://google.com"}`},
			want: want{
				code:        http.StatusCreated,
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name: "invalid rec",
			args: args{`https://google.com`},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "application/json; charset=utf-8",
				err:         true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &PostShrtHandler{s, p, l}
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.args.post))
			w := httptest.NewRecorder()
			h.PostShrt(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)
			if tt.want.err {
				return
			}
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Len(t, string(resBody), len("http://")+shortener.ShortLen+len(request.Host)+14)
		})
	}
}
