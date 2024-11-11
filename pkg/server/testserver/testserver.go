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

func New(h http.Handler) *Server {
	ts := httptest.NewTLSServer(h)

	ts.Client().CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &Server{ts}
}

func (ts *Server) Get(t *testing.T, urlPath string) (status int, header http.Header, body string) {
	//nolint: noctx
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	b, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, string(bytes.TrimSpace(b))
}
