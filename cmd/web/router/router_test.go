package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInitGin(t *testing.T) {
	engine := InitGin()
	ts := httptest.NewServer(engine)
	// shutdown
	defer ts.Close()

	resp, err := http.Get(fmt.Sprintf("%s/ping", ts.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code:200, got:%v", resp.StatusCode)
	}
}
