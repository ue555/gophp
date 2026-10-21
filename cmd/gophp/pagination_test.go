package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"example.com/gophp/internal/php"
)

func TestFixedLaravelCollectionPagination(t *testing.T) {
	for _, tc := range []struct {
		name    string
		version php.Version
	}{{"8.3", php.PHP83}, {"8.4", php.PHP84}, {"8.5", php.PHP85}} {
		t.Run(tc.name, func(t *testing.T) {
			testFixedLaravelCollectionPagination(t, tc.version)
		})
	}
}

func testFixedLaravelCollectionPagination(t *testing.T, version php.Version) {
	filename, err := filepath.Abs("../../research/laravel-api/public/index.php")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(filepath.Dir(filename)), "vendor/autoload.php")); os.IsNotExist(err) {
		t.Skip("fixed vendor is not installed")
	} else if err != nil {
		t.Fatal(err)
	}
	source, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	database := filepath.Join(t.TempDir(), "database.sqlite")
	var setup strings.Builder
	fmt.Fprintf(&setup, "<?php $db = new PDO('sqlite:%s'); $db->exec('CREATE TABLE books (id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT NOT NULL, created_at TEXT, updated_at TEXT)');", database)
	for i := 1; i <= 31; i++ {
		fmt.Fprintf(&setup, "$db->exec(\"INSERT INTO books (title) VALUES ('Book%02d')\");", i)
	}
	var out, diagnostics bytes.Buffer
	if code := php.RunWithVersion(setup.String(), "setup.php", &out, &diagnostics, version); code != 0 || diagnostics.Len() != 0 {
		t.Fatalf("setup exit=%d stderr=%q", code, diagnostics.String())
	}
	t.Setenv("DB_DATABASE", database)
	t.Setenv("DB_CONNECTION", "sqlite")
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_KEY", "base64:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=")
	storage := filepath.Join(t.TempDir(), "storage")
	for _, dir := range []string{"framework/cache/data", "framework/views", "framework/sessions", "logs"} {
		if err := os.MkdirAll(filepath.Join(storage, dir), 0700); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("LARAVEL_STORAGE_PATH", storage)
	handler := phpHandler(string(source), filename, version, &diagnostics)
	seenIDs := make(map[int]bool)
	for _, tc := range []struct{ page, count int }{{1, 30}, {2, 1}, {3, 0}} {
		diagnostics.Reset()
		request := httptest.NewRequest("GET", fmt.Sprintf("http://example.test/api/books?page=%d", tc.page), nil)
		request.Header.Set("Accept", "application/ld+json")
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		bodyPrefix := response.Body.String()
		if len(bodyPrefix) > 160 {
			bodyPrefix = bodyPrefix[:160]
		}
		t.Logf("page=%d status=%d content_type=%q body_bytes=%d body_prefix=%q diagnostics=%q", tc.page, response.Code, response.Header().Get("Content-Type"), response.Body.Len(), bodyPrefix, diagnostics.String())
		if response.Code != 200 || diagnostics.Len() != 0 {
			t.Fatalf("page=%d status=%d diagnostics=%q", tc.page, response.Code, diagnostics.String())
		}
		var payload struct {
			Context    string `json:"@context"`
			ID         string `json:"@id"`
			Type       string `json:"@type"`
			TotalItems int    `json:"totalItems"`
			View       struct {
				Type string `json:"@type"`
			} `json:"view"`
			Member []struct {
				ID    int    `json:"id"`
				Title string `json:"title"`
			} `json:"member"`
		}
		if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
			t.Fatal(err)
		}
		if payload.Context != "/api/contexts/Book" || payload.ID != "/api/books" || payload.Type != "Collection" || payload.TotalItems != 31 || payload.View.Type != "PartialCollectionView" || len(payload.Member) != tc.count {
			t.Fatalf("page=%d invalid JSON-LD collection: context=%q id=%q type=%q total=%d view=%q items=%d, want %d", tc.page, payload.Context, payload.ID, payload.Type, payload.TotalItems, payload.View.Type, len(payload.Member), tc.count)
		}
		for _, item := range payload.Member {
			if item.ID < 1 || seenIDs[item.ID] || item.Title == "" {
				t.Fatalf("page=%d duplicate or invalid item: %#v", tc.page, item)
			}
			seenIDs[item.ID] = true
		}
	}
	if len(seenIDs) != 31 {
		t.Fatalf("pagination returned %d distinct items, want 31", len(seenIDs))
	}
}
