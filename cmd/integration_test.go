package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Zentech-Development/conductor-proxy/pkg/config"
)

func getPath(path string) string {
	config := config.GetConfig()
	return fmt.Sprintf("http://%s%s", config.Host, path)
}

func TestMain(m *testing.M) {
	config.SetAndGetConfig("")
	m.Run()
	os.Exit(0)
}

func TestHealthcheck(t *testing.T) {
	config := config.GetConfig()
	server := setupApp(config)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", getPath("/"), nil)
	server.ServeHTTP(w, req)

	if w.Result().StatusCode != 200 {
		t.Fatal("Failed to get healtcheck response")
	}
}
