package gzip

import (
	"compress/gzip"
	"io"
	"net/http"
)

type gzipResponseWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func (g *gzipResponseWriter) Header() http.Header {
	return g.w.Header()
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	return g.zw.Write(b)
}

func (g *gzipResponseWriter) WriteHeader(sc int) {
	if sc < 300 {
		g.w.Header().Set("Content-Encoding", "gzip")
	}
	g.w.WriteHeader(sc)
}

func (g *gzipResponseWriter) Close() error {
	return g.zw.Close()
}

func newGzipResponseWriter(w http.ResponseWriter) *gzipResponseWriter {
	zw, err := gzip.NewWriterLevel(w, gzip.BestCompression)
	if err != nil {
		panic(err)
	}
	return &gzipResponseWriter{w: w, zw: zw}
}

type gzipReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func (g *gzipReader) Read(p []byte) (n int, err error) {
	return g.zr.Read(p)
}

func (g *gzipReader) Close() error {
	if err := g.zr.Close(); err != nil {
		return err
	}
	return g.r.Close()
}

func newGzipReader(r io.ReadCloser) (*gzipReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &gzipReader{r: r, zr: zr}, nil
}
