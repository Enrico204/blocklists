package utils

import (
	"go.uber.org/zap"
	"net/http"
	"net/http/httputil"
)

type HTTPLoggingTransport struct {
	Transport http.RoundTripper
	Logger    *zap.SugaredLogger
}

func (s *HTTPLoggingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	bytes, _ := httputil.DumpRequestOut(r, false)
	s.Logger.Debug(string(bytes))

	if s.Transport == nil {
		s.Transport = http.DefaultTransport
	}
	resp, err := s.Transport.RoundTrip(r)
	if err == nil {
		bytes, _ = httputil.DumpResponse(resp, false)
		s.Logger.Debug(string(bytes))
	}

	return resp, err
}
