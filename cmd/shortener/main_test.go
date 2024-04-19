package main

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_run(t *testing.T) {
	ts := httptest.NewServer(router())
	defer ts.Close()

	var testTable = []struct {
		want   string
		status int
	}{
		{"https://google.com", http.StatusOK},
	}
	for _, v := range testTable {
		postResp, shrtURL := testRequest(t, ts, "POST", ts.URL, v.want)
		assert.Equal(t, http.StatusCreated, postResp.StatusCode())
		assert.NotEmpty(t, shrtURL)

		getResp, postBody := testRequest(t, ts, "GET", shrtURL, "")
		assert.Equal(t, http.StatusTemporaryRedirect, getResp.StatusCode())
		assert.Equal(t, v.want, getResp.Header().Get("location"))
		assert.Empty(t, postBody)
	}
}

var ErrRedirectBlock = errors.New("HTTP redirect blocked")

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body string) (*resty.Response, string) {
	redirPolicy := resty.RedirectPolicyFunc(func(_ *http.Request, _ []*http.Request) error {
		return ErrRedirectBlock
	})
	httpc := resty.New().
		SetBaseURL(ts.URL).
		SetRedirectPolicy(redirPolicy)

	req := httpc.R().
		SetBody(body)
	resp, err := req.Execute(method, path)
	if !errors.Is(err, ErrRedirectBlock) {
		require.NoError(t, err)
	}

	return resp, string(resp.Body())
}
