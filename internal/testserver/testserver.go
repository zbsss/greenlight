package testserver

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Server struct {
	*httptest.Server
}

func New(t *testing.T, h http.Handler) *Server {
	ts := httptest.NewTLSServer(h)

	ts.Client().CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &Server{ts}
}

func (ts *Server) Get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}
