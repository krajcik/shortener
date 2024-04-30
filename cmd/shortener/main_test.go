package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"krajcik/shortener/cmd/shortener/handler/api"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_run(t *testing.T) {
	r, err := router()
	assert.NoError(t, err)
	ts := httptest.NewServer(r)
	defer ts.Close()

	var testTable = []struct {
		want   string
		status int
	}{
		{"https://google.com", http.StatusOK},
	}
	for _, v := range testTable {
		postResp, shrtURL := testRequest(t, ts, "POST", ts.URL, v.want, nil)
		assert.Equal(t, http.StatusCreated, postResp.StatusCode())
		assert.NotEmpty(t, shrtURL)

		getResp, postBody := testRequest(t, ts, "GET", shrtURL, "", nil)
		assert.Equal(t, http.StatusTemporaryRedirect, getResp.StatusCode())
		assert.Equal(t, v.want, getResp.Header().Get("location"))
		assert.Empty(t, postBody)
	}
}

func Test_run_api(t *testing.T) {
	r, err := router()
	assert.NoError(t, err)
	ts := httptest.NewServer(r)
	defer ts.Close()
	var testTable = []struct {
		headers map[string]string
		url     string
		want    string
		status  int
	}{
		{
			map[string]string{
				"Accept":           "application/json",
				"Content-type":     "application/json",
				"Accept-Encoding":  "gzip, deflate, br",
				"Content-Encoding": "gzip",
			},
			`{"url":"https://google.com"}`,
			`https://google.com`,
			http.StatusOK,
		},
	}
	for _, v := range testTable {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(v.url))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		postResp, respString := testRequest(t, ts, "POST", ts.URL+"/api/shorten", buf.String(), v.headers)
		assert.Equal(t, http.StatusCreated, postResp.StatusCode())
		assert.NotEmpty(t, respString)
		psrm := &api.PostShrtRespModel{}
		err = json.Unmarshal([]byte(respString), psrm)
		require.NoError(t, err)

		getResp, postBody := testRequest(t, ts, "GET", psrm.ShrtURL, "", nil)
		assert.Equal(t, http.StatusTemporaryRedirect, getResp.StatusCode())
		assert.Equal(t, v.want, getResp.Header().Get("location"))
		assert.Empty(t, postBody)
	}
}

var ErrRedirectBlock = errors.New("HTTP redirect blocked")

func testRequest(
	t *testing.T,
	ts *httptest.Server,
	method, path string,
	body string,
	headers map[string]string,
) (*resty.Response, string) {
	redirPolicy := resty.RedirectPolicyFunc(
		func(_ *http.Request, _ []*http.Request) error {
			return ErrRedirectBlock
		},
	)
	httpc := resty.New().
		SetBaseURL(ts.URL).
		SetRedirectPolicy(redirPolicy)

	req := httpc.R().
		SetBody(body)
	if headers != nil {
		req.SetHeaders(headers)
	}
	resp, err := req.Execute(method, path)
	if !errors.Is(err, ErrRedirectBlock) {
		require.NoError(t, err)
	}

	return resp, string(resp.Body())
}
