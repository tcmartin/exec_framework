package test

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
    retryablehttp "github.com/hashicorp/go-retryablehttp"
)

func TestHTTPRequestRetry(t *testing.T) {
    attempts := 0
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        attempts++
        if attempts < 3 {
            w.WriteHeader(500)
            return
        }
        w.Write([]byte(`{"items":[{"ok":true}]}`))
    }))
    defer srv.Close()

    cli := retryablehttp.NewClient()
    start := time.Now()
    req, _ := retryablehttp.NewRequest("GET", srv.URL, nil)
    resp, err := cli.Do(req)
    if err != nil {
        t.Fatal(err)
    }
    defer resp.Body.Close()
    if attempts != 3 {
        t.Errorf("expected 3 attempts, got %d", attempts)
    }
    if time.Since(start) < 100*time.Millisecond {
        t.Errorf("expected backoff delay, too fast")
    }
}
