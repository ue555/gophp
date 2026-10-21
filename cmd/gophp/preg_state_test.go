package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"example.com/gophp/internal/php"
)

func TestCookiePercentDecodeTCP(t *testing.T) {
	server := httptest.NewServer(phpHandler(`<?php header('Content-Type: application/json'); echo json_encode($_COOKIE);`, "cookies.php", php.PHP84, io.Discard))
	defer server.Close()
	request, _ := http.NewRequest("GET", server.URL, nil)
	request.Header.Set("Cookie", "session=abc%2Bdef%3D; a=first; a=second; space=a+b")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil || strings.TrimSpace(string(body)) != `{"session":"abc+def=","a":"first","space":"a b"}` {
		t.Fatalf("cookie decoding: %q %v", body, err)
	}
}

// NewServer uses a real loopback TCP listener. This checks the preg state at
// the same request boundary as --serve, independently of Laravel fixtures.
func TestPregStateAcrossTCPRequests(t *testing.T) {
	for _, version := range []php.Version{php.PHP83, php.PHP84, php.PHP85} {
		t.Run(string(version), func(t *testing.T) {
			source := `<?php $before=preg_last_error(); if($_GET['bad']==='1') { preg_match('//u',"\xff"); } else { preg_match('/a/','a'); } header('Content-Type: application/json'); echo json_encode([$before,preg_last_error()]);`
			server := httptest.NewServer(phpHandler(source, "state.php", version, io.Discard))
			defer server.Close()
			request := func(i int) {
				bad := "0"
				want := 0
				if i%2 == 0 {
					bad = "1"
					want = 4
				}
				r, err := http.Get(server.URL + "/?bad=" + bad)
				if err != nil {
					t.Error(err)
					return
				}
				defer r.Body.Close()
				var state []int
				err = json.NewDecoder(r.Body).Decode(&state)
				if err != nil || r.StatusCode != 200 || len(state) != 2 || state[0] != 0 || state[1] != want {
					t.Errorf("request %d status %d state %v error %v", i, r.StatusCode, state, err)
				}
			}
			for i := 0; i < 100; i++ {
				request(i)
			}
			var wg sync.WaitGroup
			for i := 0; i < 8; i++ {
				wg.Add(1)
				go func(i int) { defer wg.Done(); request(i) }(i)
			}
			wg.Wait()
		})
	}
}
