package metrics

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	server := NewServer(":0") // Use port 0 to let OS assign a free port
	
	if server == nil {
		t.Fatal("NewServer should return a non-nil server")
	}
	
	if server.addr != ":0" {
		t.Errorf("Expected address ':0', got '%s'", server.addr)
	}
	
	if server.server == nil {
		t.Error("Server should have initialized HTTP server")
	}
}

func TestServer_Address(t *testing.T) {
	expectedAddr := ":8080"
	server := NewServer(expectedAddr)
	
	if server.Address() != expectedAddr {
		t.Errorf("Expected address '%s', got '%s'", expectedAddr, server.Address())
	}
}

func TestServer_StartStop(t *testing.T) {
	server := NewServer(":0") // Use port 0 for testing
	
	// Start server in a goroutine
	serverErrCh := make(chan error, 1)
	go func() {
		err := server.Start()
		serverErrCh <- err
	}()
	
	// Give server a moment to start
	time.Sleep(100 * time.Millisecond)
	
	// Stop server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	err := server.Stop(ctx)
	if err != nil {
		t.Errorf("Server stop failed: %v", err)
	}
	
	// Check that server stopped gracefully
	select {
	case err := <-serverErrCh:
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("Expected server to stop gracefully, got error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Error("Server did not stop within timeout")
	}
}

func TestServer_HealthEndpoint(t *testing.T) {
	server := NewServer("127.0.0.1:0") // Use specific interface and port 0
	
	// Start server
	serverErrCh := make(chan error, 1)
	go func() {
		err := server.Start()
		serverErrCh <- err
	}()
	
	// Give server time to start and find out which port it's using
	time.Sleep(100 * time.Millisecond)
	
	// Get the actual listening address
	// Since we used port 0, we need to find the actual port
	// For testing purposes, we'll test the handler directly
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		server.Stop(ctx)
	}()
	
	// Test the health endpoint handler directly
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	
	// Create a response recorder
	rr := &testResponseWriter{
		header: make(http.Header),
		body:   &strings.Builder{},
	}
	
	// Get the handler from server's mux
	server.server.Handler.ServeHTTP(rr, req)
	
	if rr.statusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.statusCode)
	}
	
	expectedBody := "OK"
	if rr.body.String() != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, rr.body.String())
	}
}

func TestServer_MetricsEndpoint(t *testing.T) {
	server := NewServer("127.0.0.1:0")
	
	// Start server
	go func() {
		server.Start()
	}()
	
	time.Sleep(100 * time.Millisecond)
	
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		server.Stop(ctx)
	}()
	
	// Test the metrics endpoint handler directly
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	
	// Create a response recorder
	rr := &testResponseWriter{
		header: make(http.Header),
		body:   &strings.Builder{},
	}
	
	// Get the handler from server's mux
	server.server.Handler.ServeHTTP(rr, req)
	
	if rr.statusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.statusCode)
	}
	
	// Check that response contains Prometheus metrics
	body := rr.body.String()
	if !strings.Contains(body, "# HELP") || !strings.Contains(body, "# TYPE") {
		t.Error("Expected Prometheus metrics format in response body")
	}
	
	// Check for some basic Go metrics that should always be present
	if !strings.Contains(body, "go_goroutines") {
		t.Error("Expected Go runtime metrics in response")
	}
}

func TestServer_StopTimeout(t *testing.T) {
	server := NewServer(":0")
	
	// Start server
	go func() {
		server.Start()
	}()
	
	time.Sleep(50 * time.Millisecond)
	
	// Test stop with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	
	err := server.Stop(ctx)
	// Should succeed even with short timeout since server isn't busy
	if err != nil && err != context.DeadlineExceeded {
		t.Errorf("Expected no error or deadline exceeded, got: %v", err)
	}
}

// testResponseWriter is a simple implementation of http.ResponseWriter for testing
type testResponseWriter struct {
	header     http.Header
	body       *strings.Builder
	statusCode int
}

func (w *testResponseWriter) Header() http.Header {
	return w.header
}

func (w *testResponseWriter) Write(data []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return w.body.Write(data)
}

func (w *testResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

// Implement additional interfaces that http.ResponseWriter might need
var _ http.ResponseWriter = (*testResponseWriter)(nil)
var _ http.Flusher = (*testResponseWriter)(nil)
var _ io.StringWriter = (*testResponseWriter)(nil)

func (w *testResponseWriter) Flush() {
	// No-op for testing
}

func (w *testResponseWriter) WriteString(s string) (int, error) {
	return w.Write([]byte(s))
}