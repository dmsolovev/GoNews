// main_test.go
package main

import (
	"bytes"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestMainServer(t *testing.T) {
	// Запускаем сервер в отдельном процессе
	cmd := exec.Command("go", "run", "main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Запуск в фоновом режиме
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer cmd.Process.Kill()

	// Даем время на запуск сервера
	time.Sleep(2 * time.Second)

	// Тестируем эндпоинты
	t.Run("Test /news endpoint", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/news")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Test invalid endpoint", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/invalid")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func TestDBInitialization(t *testing.T) {
	tests := []struct {
		name      string
		connStr   string
		wantError bool
	}{
		{
			name:      "Valid connection",
			connStr:   "postgres://postgres:1@localhost:5432/GoNews_test",
			wantError: false,
		},
		{
			name:      "Invalid connection",
			connStr:   "postgres://invalid:wrong@localhost:5432/GoNews_test",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("go", "run", "main.go")
			cmd.Env = append(os.Environ(), "DB_CONN_STRING="+tt.connStr)

			err := cmd.Run()
			if (err != nil) != tt.wantError {
				t.Errorf("Test %s failed: want error %v, got %v", tt.name, tt.wantError, err)
			}
		})
	}
}

func TestHelpFlag(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "-h")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Expected error for help flag")
	}

	if !bytes.Contains(output, []byte("Usage")) {
		t.Errorf("Expected help message, got: %s", output)
	}
}