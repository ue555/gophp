package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Static serving is opt-in: applications explicitly choose their public root.
// Only known asset types are served; PHP/source/config files and directory
// listings are never exposed. Missing assets are left to the front controller.
func publicAssetHandler(root string, next http.Handler) (http.Handler, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("gophp: document root is not a directory: %s", root)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// URL.Path is already percent-decoded by net/http. Reject ambiguous
		// double encoding and platform separators before any file lookup.
		path := r.URL.Path
		if strings.ContainsAny(path, "\\%\x00") {
			http.NotFound(w, r)
			return
		}
		for _, segment := range strings.Split(path, "/") {
			if strings.HasPrefix(segment, ".") {
				http.NotFound(w, r)
				return
			}
		}
		switch strings.ToLower(filepath.Ext(path)) {
		case ".css", ".js", ".mjs", ".json", ".txt", ".html", ".htm", ".svg", ".ico", ".png", ".jpg", ".jpeg", ".gif", ".webp", ".avif", ".woff", ".woff2", ".ttf", ".otf", ".eot", ".wasm", ".pdf":
		default:
			next.ServeHTTP(w, r)
			return
		}
		name := filepath.Join(root, filepath.FromSlash(strings.TrimPrefix(path, "/")))
		resolved, err := filepath.EvalSymlinks(name)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		rel, err := filepath.Rel(root, resolved)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			http.NotFound(w, r)
			return
		}
		// Do not let an allowed extension alias a hidden file or PHP source
		// via symlink. Asset extensions must match after resolution.
		if strings.ToLower(filepath.Ext(resolved)) != strings.ToLower(filepath.Ext(name)) {
			http.NotFound(w, r)
			return
		}
		for _, part := range strings.Split(rel, string(filepath.Separator)) {
			if strings.HasPrefix(part, ".") {
				http.NotFound(w, r)
				return
			}
		}
		file, err := os.Open(resolved)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer file.Close()
		info, err := file.Stat()
		if err != nil || !info.Mode().IsRegular() {
			http.NotFound(w, r)
			return
		}
		if r.Method != "GET" && r.Method != "HEAD" {
			w.Header().Set("Allow", "GET, HEAD")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("X-Content-Type-Options", "nosniff")
		http.ServeContent(w, r, info.Name(), info.ModTime(), file)
	}), nil
}
