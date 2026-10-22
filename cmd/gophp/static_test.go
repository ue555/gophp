package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPublicAssetsOverTCP(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "public")
	if err := os.Mkdir(root, 0700); err != nil {
		t.Fatal(err)
	}
	for file, data := range map[string]string{"public/site.css": "body{color:red}", "public/app.js": "console.log('ok')", "public/font.woff2": "wOF2data", "public/index.php": "secret-php", "public/.env": "secret-env", "outside.css": "secret-outside"} {
		if err := os.WriteFile(filepath.Join(parent, file), []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
	}
	for name, target := range map[string]string{"escape.css": "../outside.css", "source.css": "index.php", "hidden.css": ".env", "alias.css": "site.css"} {
		if err := os.Symlink(target, filepath.Join(root, name)); err != nil {
			t.Fatal(err)
		}
	}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Error(w, "front-controller", 404) })
	handler, err := publicAssetHandler(root, next)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()
	for _, tc := range []struct {
		path, body, media string
		status            int
	}{
		{"/site.css?v=1", "body{color:red}", "text/css", 200},
		{"/alias.css", "body{color:red}", "text/css", 200},
		{"/app.js", "console.log('ok')", "javascript", 200},
		{"/font.woff2", "wOF2data", "font/woff2", 200},
		{"/index.php", "front-controller\n", "", 404},
		{"/missing.css", "front-controller\n", "", 404},
	} {
		response, err := http.Get(server.URL + tc.path)
		if err != nil {
			t.Fatal(err)
		}
		body, err := io.ReadAll(response.Body)
		response.Body.Close()
		if err != nil || response.StatusCode != tc.status || string(body) != tc.body || !strings.Contains(response.Header.Get("Content-Type"), tc.media) {
			t.Fatalf("%s: status=%d type=%q body=%q error=%v", tc.path, response.StatusCode, response.Header.Get("Content-Type"), body, err)
		}
	}
	for _, path := range []string{"/.env", "/%2eenv", "/../outside.css", "/%2e%2e/outside.css", "/%252e%252e/outside.css", "/a/%2e%2e/site.css", "/..%5coutside.css", "/escape.css", "/source.css", "/hidden.css", "/"} {
		response, err := http.Get(server.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		body, _ := io.ReadAll(response.Body)
		response.Body.Close()
		if response.StatusCode != 404 || strings.Contains(string(body), "secret") {
			t.Fatalf("%s: %d %q", path, response.StatusCode, body)
		}
	}
	response, err := http.Head(server.URL + "/site.css")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	if response.StatusCode != 200 || len(body) != 0 || response.ContentLength != 15 {
		t.Fatalf("HEAD: %d %q length=%d", response.StatusCode, body, response.ContentLength)
	}
}
