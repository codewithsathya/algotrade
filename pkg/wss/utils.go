package wss

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

func (w *WebSocket) getHeaders() http.Header {
	headers := http.Header{}
	headers.Add("Origin", w.Origin)
	headers.Add("User-Agent", "Mozilla/5.0 (iPad; CPU OS 7_0 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) CriOS/45.0.2454.68 Mobile/11A465 Safari/9537.53")
	return headers
}

func (w *WebSocket) getCookieJar() (http.CookieJar, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to construct cookie jar | %w", err)
	}
	url, err := url.Parse(w.Origin)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url | %w", err)
	}
	jar.SetCookies(url, w.Cookies)
	return jar, nil
}

func getTlsConfig() (*tls.Config, error) {
	cacertFile := "cacert.pem"
	caCert, err := os.ReadFile(cacertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate bundle | %w", err)
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		return nil, fmt.Errorf("failed to append certs from pem | %w", err)
	}

	return &tls.Config{
		RootCAs:    certPool,
		MinVersion: tls.VersionTLS12,
	}, nil
}