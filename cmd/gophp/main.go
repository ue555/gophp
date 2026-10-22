package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"example.com/gophp/internal/php"
)

func run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("gophp", flag.ContinueOnError)
	flags.SetOutput(stderr)
	target := flags.String("php-version", string(php.PHP84), "target PHP version: 8.3, 8.4 or 8.5")
	listen := flags.String("serve", "", "serve PHP requests at this local address")
	documentRoot := flags.String("document-root", "", "serve public static assets from this directory (requires --serve)")
	flags.Usage = func() {
		fmt.Fprintln(stderr, "usage: gophp [--php-version=8.3|8.4|8.5] <file.php> [arguments...]")
	}
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	if flags.NArg() < 1 {
		flags.Usage()
		return 2
	}
	version, err := php.ParseVersion(*target)
	if err != nil {
		fmt.Fprintf(stderr, "gophp: %v\n", err)
		return 2
	}
	filename, err := filepath.Abs(flags.Arg(0))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	source, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if *listen != "" {
		handler := phpHandler(string(source), filename, version, stderr)
		if *documentRoot != "" {
			handler, err = publicAssetHandler(*documentRoot, handler)
			if err != nil {
				fmt.Fprintln(stderr, err)
				return 1
			}
		}
		if err := http.ListenAndServe(*listen, handler); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		return 0
	}
	if *documentRoot != "" {
		fmt.Fprintln(stderr, "gophp: --document-root requires --serve")
		return 2
	}
	return php.RunWithArguments(string(source), filename, append([]string{flags.Arg(0)}, flags.Args()[1:]...), stdout, stderr, version)
}

func phpHandler(source, filename string, version php.Version, stderr io.Writer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(io.LimitReader(r.Body, (16<<20)+1))
		if err != nil || len(body) > 16<<20 {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
			return
		}
		headers := make(map[string]string, len(r.Header))
		for name, values := range r.Header {
			headers[name] = strings.Join(values, ", ")
		}
		headers["Content-Length"] = fmt.Sprint(len(body))
		addr, _, _ := net.SplitHostPort(r.RemoteAddr)
		var output, diagnostics bytes.Buffer
		code, response := php.RunHTTPRequest(source, filename, php.HTTPRequest{Method: r.Method, URI: r.RequestURI, Host: r.Host, Body: string(body), RemoteAddr: addr, Headers: headers}, &output, &diagnostics, version)
		if diagnostics.Len() > 0 {
			fmt.Fprint(stderr, diagnostics.String())
		}
		for _, item := range response.Headers {
			parts := strings.SplitN(item, ":", 2)
			if len(parts) == 2 {
				w.Header().Add(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
			}
		}
		status := response.Status
		if code != 0 && status < 400 {
			status = 500
		}
		if status < 100 || status > 599 {
			status = 500
		}
		w.WriteHeader(status)
		_, _ = w.Write(output.Bytes())
	})
}

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }
