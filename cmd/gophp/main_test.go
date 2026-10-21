package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"example.com/gophp/internal/php"
)

type synchronizedBuffer struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (w *synchronizedBuffer) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.b.Write(p)
}

func (w *synchronizedBuffer) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.b.String()
}

func TestHTTPBridgeKeepsRequestsSeparate(t *testing.T) {
	source := `<?php header('Content-Type: application/json'); http_response_code(201); echo json_encode(['method' => $_SERVER['REQUEST_METHOD'], 'body' => file_get_contents('php://input'), 'key' => $_GET['k'] ?? null]);`
	handler := phpHandler(source, "/tmp/index.php", php.PHP84, &bytes.Buffer{})
	for _, tc := range []struct{ method, uri, body, want string }{
		{"POST", "/api/probe?k=a", "one", `{"method":"POST","body":"one","key":"a"}`},
		{"GET", "/api/probe?k=b", "", `{"method":"GET","body":"","key":"b"}`},
	} {
		req := httptest.NewRequest(tc.method, tc.uri, strings.NewReader(tc.body))
		result := httptest.NewRecorder()
		handler.ServeHTTP(result, req)
		if result.Code != 201 || result.Header().Get("Content-Type") != "application/json" || result.Body.String() != tc.want {
			t.Fatalf("%s: status=%d header=%q body=%q", tc.uri, result.Code, result.Header().Get("Content-Type"), result.Body.String())
		}
	}
}

func TestHTTPBridgeParsesNestedQueryAndFormValues(t *testing.T) {
	source := `<?php echo json_encode([$_GET, $_POST]);`
	var diagnostics bytes.Buffer
	handler := phpHandler(source, "/tmp/index.php", php.PHP84, &diagnostics)
	request := httptest.NewRequest("POST", "/probe?filter[name]=Ada+Lovelace&list[]=a&list[]=b&filter.name=flat", strings.NewReader("items[]=x&items[]=y&filter[active]=1"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	result := httptest.NewRecorder()
	handler.ServeHTTP(result, request)
	want := `[{"filter":{"name":"Ada Lovelace"},"list":["a","b"],"filter_name":"flat"},{"items":["x","y"],"filter":{"active":"1"}}]`
	if result.Code != 200 || result.Body.String() != want || diagnostics.Len() != 0 {
		t.Fatalf("status=%d body=%q diagnostics=%q", result.Code, result.Body.String(), diagnostics.String())
	}
}

func TestHTTPBridgePostRequiresPostMethodAndFormMediaType(t *testing.T) {
	source := `<?php echo json_encode([$_POST, file_get_contents('php://input')]);`
	var diagnostics bytes.Buffer
	handler := phpHandler(source, "/tmp/index.php", php.PHP84, &diagnostics)
	for _, tc := range []struct{ method, contentType string }{
		{"PUT", "application/x-www-form-urlencoded"},
		{"POST", "application/x-www-form-urlencoded-extra"},
	} {
		request := httptest.NewRequest(tc.method, "/probe", strings.NewReader("items[]=x"))
		request.Header.Set("Content-Type", tc.contentType)
		result := httptest.NewRecorder()
		handler.ServeHTTP(result, request)
		if result.Code != 200 || result.Body.String() != `[[],"items[]=x"]` || diagnostics.Len() != 0 {
			t.Fatalf("method=%s content-type=%s status=%d body=%q diagnostics=%q", tc.method, tc.contentType, result.Code, result.Body.String(), diagnostics.String())
		}
	}
}

func TestHTTPBridgeDoesNotRetainResponseState(t *testing.T) {
	source := `<?php if (isset($_GET['set'])) { header('X-Trace: first'); http_response_code(201); } echo 'response';`
	handler := phpHandler(source, "/tmp/index.php", php.PHP84, &bytes.Buffer{})
	for _, tc := range []struct {
		uri, header string
		status      int
	}{
		{"/probe?set=1", "first", 201},
		{"/probe", "", 200},
	} {
		result := httptest.NewRecorder()
		handler.ServeHTTP(result, httptest.NewRequest("GET", tc.uri, nil))
		if result.Code != tc.status || result.Header().Get("X-Trace") != tc.header || result.Body.String() != "response" {
			t.Fatalf("%s: status=%d header=%q body=%q", tc.uri, result.Code, result.Header().Get("X-Trace"), result.Body.String())
		}
	}
}

func TestFixedVendorAttributeDiscoveryAndContainer(t *testing.T) {
	appDir, err := filepath.Abs("../../research/laravel-api")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(appDir, "vendor/autoload.php")); os.IsNotExist(err) {
		t.Skip("fixed vendor is not installed")
	} else if err != nil {
		t.Fatal(err)
	}
	source := `<?php
require __DIR__ . '/vendor/autoload.php';
#[Attribute(Attribute::TARGET_PARAMETER)]
class InjectNumber implements Illuminate\Contracts\Container\ContextualAttribute {
    public static function resolve(self $attribute, Illuminate\Contracts\Container\Container $container) {
        return 42;
    }
}
class NeedsNumber {
    public function __construct(#[InjectNumber] public int $number) {}
}
$container = new Illuminate\Container\Container();
echo $container->make(NeedsNumber::class)->number, '|';
$factory = new ApiPlatform\Metadata\Resource\Factory\AttributesResourceNameCollectionFactory([__DIR__ . '/app/Models']);
foreach ($factory->create() as $className) {
    echo $className;
}
`
	for _, version := range []php.Version{php.PHP83, php.PHP84, php.PHP85} {
		var stdout, stderr bytes.Buffer
		code := php.RunWithVersion(source, filepath.Join(appDir, "gophp-attribute-probe.php"), &stdout, &stderr, version)
		if code != 0 || stdout.String() != `42|App\Models\Book` || stderr.Len() != 0 {
			t.Fatalf("version=%v exit=%d stdout=%q stderr=%q", version, code, stdout.String(), stderr.String())
		}
	}
}

func TestFixedLaravelCustomCommand(t *testing.T) {
	artisan, err := filepath.Abs("../../research/laravel-api/artisan")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(artisan), "vendor/autoload.php")); os.IsNotExist(err) {
		t.Skip("fixed vendor is not installed")
	} else if err != nil {
		t.Fatal(err)
	}
	t.Setenv("APP_ENV", "testing")
	t.Setenv("DB_DATABASE", ":memory:")
	t.Setenv("APP_KEY", "base64:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=")
	// Official PHP package discovery can populate the app's metadata cache.
	// Use fresh storage, as the HTTP integration tests do, so this test does
	// not depend on persisted cache contents from another interpreter/run.
	storage := t.TempDir()
	for _, dir := range []string{"framework/cache/data", "framework/views", "framework/sessions", "logs"} {
		if err := os.MkdirAll(filepath.Join(storage, dir), 0700); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("LARAVEL_STORAGE_PATH", storage)
	for _, version := range []string{"8.3", "8.4", "8.5"} {
		var stdout, stderr bytes.Buffer
		code := run([]string{"--php-version=" + version, artisan, "gophp:probe"}, &stdout, &stderr)
		if code != 0 || strings.TrimSpace(stdout.String()) != "probe-ok" || stderr.Len() != 0 {
			t.Fatalf("version=%s exit=%d stdout=%q stderr=%q", version, code, stdout.String(), stderr.String())
		}
	}
}

func TestFixedLaravelNestedTranslation(t *testing.T) {
	appDir, err := filepath.Abs("../../research/laravel-api")
	if err != nil {
		t.Fatal(err)
	}
	autoload := filepath.Join(appDir, "vendor/autoload.php")
	if _, err := os.Stat(autoload); os.IsNotExist(err) {
		t.Skip("fixed vendor is not installed")
	} else if err != nil {
		t.Fatal(err)
	}
	source := fmt.Sprintf(`<?php
require %q;
$loader = new Illuminate\Translation\ArrayLoader;
$loader->addMessages('en', 'messages', ['nested' => ['first' => 'Hi :name', 'deeper' => ['second' => 'Bye :name']]]);
$translator = new Illuminate\Translation\Translator($loader, 'en');
echo json_encode($translator->get('messages.nested', ['name' => 'Ada']));`, filepath.ToSlash(autoload))
	for _, version := range []php.Version{php.PHP83, php.PHP84, php.PHP85} {
		var stdout, stderr bytes.Buffer
		code := php.RunWithVersion(source, filepath.Join(appDir, "gophp-translation-probe.php"), &stdout, &stderr, version)
		if code != 0 || stdout.String() != `{"first":"Hi Ada","deeper":{"second":"Bye Ada"}}` || stderr.Len() != 0 {
			t.Fatalf("version=%v exit=%d stdout=%q stderr=%q", version, code, stdout.String(), stderr.String())
		}
	}
}

func TestFixedLaravelLazyCollectionGenerators(t *testing.T) {
	appDir, err := filepath.Abs("../../research/laravel-api")
	if err != nil {
		t.Fatal(err)
	}
	autoload := filepath.Join(appDir, "vendor/autoload.php")
	if _, err := os.Stat(autoload); os.IsNotExist(err) {
		t.Skip("fixed vendor is not installed")
	} else if err != nil {
		t.Fatal(err)
	}
	source := fmt.Sprintf(`<?php
require %q;
$nested = Illuminate\Support\LazyCollection::make([1, [2, [3, 4]]]);
echo implode(',', $nested->flatten()->all()), '|';
$calls = 0;
$source = Illuminate\Support\LazyCollection::make(function () use (&$calls) {
    for ($i = 1; $i <= 3; $i++) { $calls++; yield $i; }
});
$remembered = $source->remember();
echo implode(',', $remembered->take(1)->all()), '|', $calls, '|';
echo implode(',', $remembered->all()), '|', $calls, '|';
echo implode(',', $source->flatMap(fn ($v) => [$v, $v * 10])->all()), '|';
foreach (Illuminate\Support\LazyCollection::make([1, 2, 3, 4, 5])->chunk(2) as $chunk) {
    echo implode(',', $chunk->all()), ';';
}`, filepath.ToSlash(autoload))
	want := "1,2,3,4|1|1|1,2,3|3|1,10,2,20,3,30|1,2;3,4;5;"
	for _, version := range []php.Version{php.PHP83, php.PHP84, php.PHP85} {
		var stdout, stderr bytes.Buffer
		code := php.RunWithVersion(source, filepath.Join(appDir, "gophp-lazy-collection-probe.php"), &stdout, &stderr, version)
		if code != 0 || stdout.String() != want || stderr.Len() != 0 {
			t.Fatalf("version=%v exit=%d stdout=%q stderr=%q", version, code, stdout.String(), stderr.String())
		}
	}
}

func TestFixedLaravelProbeAcrossVersions(t *testing.T) {
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
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_KEY", "base64:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=")
	t.Setenv("DB_CONNECTION", "sqlite")
	database := filepath.Join(t.TempDir(), "database.sqlite")
	if err := os.WriteFile(database, nil, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DB_DATABASE", database)
	storage := filepath.Join(t.TempDir(), "storage")
	for _, dir := range []string{"framework/cache/data", "framework/views", "framework/sessions", "logs"} {
		if err := os.MkdirAll(filepath.Join(storage, dir), 0700); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("LARAVEL_STORAGE_PATH", storage)
	for _, version := range []php.Version{php.PHP83, php.PHP84, php.PHP85} {
		var diagnostics bytes.Buffer
		handler := phpHandler(string(source), filename, version, &diagnostics)
		request := httptest.NewRequest("GET", "http://example.test/api/probe", nil)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != 200 || response.Body.String() != `{"ok":true}` || diagnostics.Len() != 0 {
			t.Fatalf("version=%s status=%d body=%q diagnostics=%q", version, response.Code, response.Body.String(), diagnostics.String())
		}
	}
}

// The primary HTTP bridge test covers the full app contract on PHP 8.4.
// Exercise the same API Platform resource lifecycle on the other configured
// language versions with separate storage and SQLite state for each version.
func TestFixedLaravelAPIPlatformAcrossVersions(t *testing.T) {
	filename, err := filepath.Abs("../../research/laravel-api/public/index.php")
	if err != nil {
		t.Fatal(err)
	}
	source, err := os.ReadFile(filename)
	if os.IsNotExist(err) {
		t.Skip("fixed research app is not installed")
	}
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(filepath.Dir(filename)), "vendor/autoload.php")); os.IsNotExist(err) {
		t.Skip("fixed vendor is not installed")
	} else if err != nil {
		t.Fatal(err)
	}
	for _, version := range []php.Version{php.PHP83, php.PHP85} {
		t.Run(string(version), func(t *testing.T) {
			database := filepath.Join(t.TempDir(), "database.sqlite")
			setup := fmt.Sprintf("<?php (new PDO('sqlite:%s'))->exec('CREATE TABLE books (id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT NOT NULL, created_at TEXT, updated_at TEXT)');", database)
			var setupOutput, setupError bytes.Buffer
			if code := php.RunWithVersion(setup, "setup.php", &setupOutput, &setupError, version); code != 0 || setupError.Len() != 0 {
				t.Fatalf("SQLite setup exit=%d stderr=%q", code, setupError.String())
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
			var diagnostics bytes.Buffer
			handler := phpHandler(string(source), filename, version, &diagnostics)
			request := func(method, path, contentType, body string, wantStatus int) *httptest.ResponseRecorder {
				t.Helper()
				req := httptest.NewRequest(method, "http://example.test"+path, strings.NewReader(body))
				req.Header.Set("Accept", "application/ld+json")
				if contentType != "" {
					req.Header.Set("Content-Type", contentType)
				}
				result := httptest.NewRecorder()
				diagnostics.Reset()
				handler.ServeHTTP(result, req)
				if result.Code != wantStatus || diagnostics.Len() != 0 {
					t.Fatalf("%s %s: status=%d body=%q diagnostics=%q", method, path, result.Code, result.Body.String(), diagnostics.String())
				}
				return result
			}
			created := request("POST", "/api/books", "application/ld+json", `{"title":"Versioned"}`, 201)
			var resource struct {
				ID    string `json:"@id"`
				Title string `json:"title"`
			}
			if err := json.Unmarshal(created.Body.Bytes(), &resource); err != nil || resource.ID == "" || resource.Title != "Versioned" {
				t.Fatalf("created resource=%#v err=%v", resource, err)
			}
			if got := request("GET", resource.ID, "", "", 200).Body.String(); !strings.Contains(got, `"title":"Versioned"`) {
				t.Fatalf("GET resource=%q", got)
			}
			if got := request("PATCH", resource.ID, "application/merge-patch+json", `{"title":"Revised"}`, 200).Body.String(); !strings.Contains(got, `"title":"Revised"`) {
				t.Fatalf("PATCH resource=%q", got)
			}
			request("DELETE", resource.ID, "", "", 204)
			request("GET", resource.ID, "", "", 404)
			var dbOutput, dbError bytes.Buffer
			verify := fmt.Sprintf("<?php $db = new PDO('sqlite:%s'); echo $db->query('SELECT COUNT(*) FROM books')->fetchColumn();", database)
			if code := php.RunWithVersion(verify, "verify.php", &dbOutput, &dbError, version); code != 0 || dbOutput.String() != "0" || dbError.Len() != 0 {
				t.Fatalf("database after DELETE exit=%d stdout=%q stderr=%q", code, dbOutput.String(), dbError.String())
			}
		})
	}
}

func TestFixedLaravelProbeThroughHTTPBridge(t *testing.T) {
	filename, err := filepath.Abs("../../research/laravel-api/public/index.php")
	if err != nil {
		t.Fatal(err)
	}
	source, err := os.ReadFile(filename)
	if os.IsNotExist(err) {
		t.Skip("fixed research app is not installed")
	}
	if err != nil {
		t.Fatal(err)
	}
	database := filepath.Join(t.TempDir(), "database.sqlite")
	setup := fmt.Sprintf("<?php (new PDO('sqlite:%s'))->exec('CREATE TABLE books (id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT NOT NULL, created_at TEXT, updated_at TEXT)');", database)
	var setupOutput, setupError bytes.Buffer
	if code := php.RunWithVersion(setup, "setup.php", &setupOutput, &setupError, php.PHP84); code != 0 {
		t.Fatalf("SQLite setup exit=%d stderr=%q", code, setupError.String())
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
	var diagnostics bytes.Buffer
	handler := phpHandler(string(source), filename, php.PHP84, &diagnostics)
	request := httptest.NewRequest("GET", "http://example.test/api/probe", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != 200 || response.Body.String() != "{\"ok\":true}" || diagnostics.Len() != 0 {
		t.Fatalf("status=%d body=%q diagnostics=%q", response.Code, response.Body.String(), diagnostics.String())
	}
	// The same handler creates a fresh PHP execution for each request. Repeated
	// dispatch must not retain headers, output, autoload state, or PDO cursors.
	for i := 0; i < 12; i++ {
		request = httptest.NewRequest("GET", fmt.Sprintf("http://example.test/api/probe?round=%d", i), nil)
		response = httptest.NewRecorder()
		diagnostics.Reset()
		handler.ServeHTTP(response, request)
		if response.Code != 200 || response.Body.String() != "{\"ok\":true}" || diagnostics.Len() != 0 {
			t.Fatalf("probe round %d: status=%d body=%q diagnostics=%q", i, response.Code, response.Body.String(), diagnostics.String())
		}
	}
	request = httptest.NewRequest("POST", "http://example.test/api/books-manual", strings.NewReader(`{"title":"Alpha"}`))
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	diagnostics.Reset()
	handler.ServeHTTP(response, request)
	if response.Code != 201 || !strings.Contains(response.Body.String(), `"title":"Alpha"`) || diagnostics.Len() != 0 {
		t.Fatalf("POST status=%d body=%q diagnostics=%q", response.Code, response.Body.String(), diagnostics.String())
	}
	var created struct {
		ID    int64  `json:"id"`
		Title string `json:"title"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &created); err != nil || created.ID < 1 || created.Title != "Alpha" {
		t.Fatalf("created: %#v, %v", created, err)
	}
	request = httptest.NewRequest("GET", fmt.Sprintf("http://example.test/api/books-manual/%d", created.ID), nil)
	response = httptest.NewRecorder()
	diagnostics.Reset()
	handler.ServeHTTP(response, request)
	if response.Code != 200 || !strings.Contains(response.Body.String(), `"title":"Alpha"`) || diagnostics.Len() != 0 {
		t.Fatalf("GET book status=%d body=%q diagnostics=%q", response.Code, response.Body.String(), diagnostics.String())
	}
	request = httptest.NewRequest("POST", "http://example.test/api/books-manual", strings.NewReader(`{}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	response = httptest.NewRecorder()
	diagnostics.Reset()
	handler.ServeHTTP(response, request)
	if response.Code != 422 || !strings.Contains(response.Body.String(), `"title"`) || diagnostics.Len() != 0 {
		t.Fatalf("invalid POST status=%d body=%q diagnostics=%q", response.Code, response.Body.String(), diagnostics.String())
	}
	request = httptest.NewRequest("POST", "http://example.test/api/books", strings.NewReader(`{"title":"Beta"}`))
	request.Header.Set("Content-Type", "application/ld+json")
	request.Header.Set("Accept", "application/ld+json")
	response = httptest.NewRecorder()
	diagnostics.Reset()
	handler.ServeHTTP(response, request)
	if response.Code != 201 || !strings.Contains(response.Body.String(), `"title":"Beta"`) || diagnostics.Len() != 0 {
		t.Fatalf("API Platform POST status=%d body=%q diagnostics=%q", response.Code, response.Body.String(), diagnostics.String())
	}
	if got := response.Header().Get("Content-Type"); !strings.HasPrefix(got, "application/ld+json") {
		t.Fatalf("API Platform POST content type=%q", got)
	}
	var resource struct {
		ID string `json:"@id"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &resource); err != nil || resource.ID == "" {
		t.Fatalf("API Platform created: %#v, %v", resource, err)
	}
	request = httptest.NewRequest("GET", "http://example.test"+resource.ID, nil)
	request.Header.Set("Accept", "application/ld+json")
	response = httptest.NewRecorder()
	diagnostics.Reset()
	handler.ServeHTTP(response, request)
	if response.Code != 200 || !strings.Contains(response.Body.String(), `"title":"Beta"`) || diagnostics.Len() != 0 {
		t.Fatalf("API Platform GET status=%d body=%q diagnostics=%q", response.Code, response.Body.String(), diagnostics.String())
	}
	if got := response.Header().Get("Content-Type"); !strings.HasPrefix(got, "application/ld+json") {
		t.Fatalf("API Platform GET content type=%q", got)
	}
	request = httptest.NewRequest("GET", "http://example.test/api/books", nil)
	request.Header.Set("Accept", "application/ld+json")
	response = httptest.NewRecorder()
	diagnostics.Reset()
	handler.ServeHTTP(response, request)
	if response.Code != 200 || !strings.Contains(response.Body.String(), `"title":"Beta"`) || diagnostics.Len() != 0 {
		t.Fatalf("API Platform collection status=%d body=%q diagnostics=%q", response.Code, response.Body.String(), diagnostics.String())
	}
	request = httptest.NewRequest("PATCH", "http://example.test"+resource.ID, strings.NewReader(`{"title":"Gamma"}`))
	request.Header.Set("Content-Type", "application/merge-patch+json")
	request.Header.Set("Accept", "application/ld+json")
	response = httptest.NewRecorder()
	diagnostics.Reset()
	handler.ServeHTTP(response, request)
	if response.Code != 200 || !strings.Contains(response.Body.String(), `"title":"Gamma"`) || diagnostics.Len() != 0 {
		t.Fatalf("API Platform PATCH status=%d body=%q diagnostics=%q", response.Code, response.Body.String(), diagnostics.String())
	}
	var dbOutput, dbErrors bytes.Buffer
	verify := fmt.Sprintf("<?php $db = new PDO('sqlite:%s'); echo $db->query('SELECT COUNT(*) FROM books')->fetchColumn();", database)
	if code := php.RunWithVersion(verify, "verify.php", &dbOutput, &dbErrors, php.PHP84); code != 0 || dbOutput.String() != "2" || dbErrors.Len() != 0 {
		t.Fatalf("database rows: exit=%d stdout=%q stderr=%q", code, dbOutput.String(), dbErrors.String())
	}
	request = httptest.NewRequest("GET", "http://example.test/api/docs.jsonopenapi", nil)
	request.Header.Set("Accept", "application/vnd.openapi+json")
	response = httptest.NewRecorder()
	diagnostics.Reset()
	handler.ServeHTTP(response, request)
	if response.Code != 200 || !strings.Contains(response.Body.String(), "openapi") || diagnostics.Len() != 0 {
		t.Fatalf("API docs status=%d body=%q diagnostics=%q", response.Code, response.Body.String(), diagnostics.String())
	}
	if got := response.Header().Get("Content-Type"); !strings.Contains(got, "json") {
		t.Fatalf("OpenAPI content type=%q", got)
	}
	var specification map[string]json.RawMessage
	if err := json.Unmarshal(response.Body.Bytes(), &specification); err != nil || len(specification["openapi"]) == 0 {
		t.Fatalf("API docs are not valid OpenAPI JSON: %v", err)
	}
	var paths map[string]map[string]json.RawMessage
	if err := json.Unmarshal(specification["paths"], &paths); err != nil || len(paths["/api/books"]["get"]) == 0 || len(paths["/api/books"]["post"]) == 0 {
		t.Fatalf("API docs lack Book collection GET/POST: %v", err)
	}
	request = httptest.NewRequest("DELETE", "http://example.test"+resource.ID, nil)
	request.Header.Set("Accept", "application/ld+json")
	response = httptest.NewRecorder()
	diagnostics.Reset()
	handler.ServeHTTP(response, request)
	if response.Code != 204 || diagnostics.Len() != 0 {
		t.Fatalf("API Platform DELETE status=%d body=%q diagnostics=%q", response.Code, response.Body.String(), diagnostics.String())
	}
	request = httptest.NewRequest("GET", "http://example.test"+resource.ID, nil)
	request.Header.Set("Accept", "application/ld+json")
	response = httptest.NewRecorder()
	diagnostics.Reset()
	handler.ServeHTTP(response, request)
	if response.Code != 404 || diagnostics.Len() != 0 {
		t.Fatalf("API Platform deleted GET status=%d body=%q diagnostics=%q", response.Code, response.Body.String(), diagnostics.String())
	}
	dbOutput.Reset()
	dbErrors.Reset()
	if code := php.RunWithVersion(verify, "verify.php", &dbOutput, &dbErrors, php.PHP84); code != 0 || dbOutput.String() != "1" || dbErrors.Len() != 0 {
		t.Fatalf("database rows after DELETE: exit=%d stdout=%q stderr=%q", code, dbOutput.String(), dbErrors.String())
	}
	// Exercise repeated writes through the same handler and database. Each PHP
	// request has fresh runtime state while the SQLite state must persist.
	seenIDs := map[string]bool{resource.ID: true}
	for i := 0; i < 8; i++ {
		title := fmt.Sprintf("Beta%d", i)
		request = httptest.NewRequest("POST", "http://example.test/api/books", strings.NewReader(fmt.Sprintf(`{"title":%q}`, title)))
		request.Header.Set("Content-Type", "application/ld+json")
		request.Header.Set("Accept", "application/ld+json")
		response = httptest.NewRecorder()
		diagnostics.Reset()
		handler.ServeHTTP(response, request)
		if response.Code != 201 || diagnostics.Len() != 0 {
			t.Fatalf("cycle %d POST: status=%d body=%q diagnostics=%q", i, response.Code, response.Body.String(), diagnostics.String())
		}
		var cycleResource struct {
			ID    string `json:"@id"`
			Title string `json:"title"`
		}
		if err := json.Unmarshal(response.Body.Bytes(), &cycleResource); err != nil || cycleResource.ID == "" || cycleResource.Title != title || seenIDs[cycleResource.ID] {
			t.Fatalf("cycle %d POST resource: %#v, err=%v", i, cycleResource, err)
		}
		seenIDs[cycleResource.ID] = true
		request = httptest.NewRequest("PATCH", "http://example.test"+cycleResource.ID, strings.NewReader(`{"title":"Updated"}`))
		request.Header.Set("Content-Type", "application/merge-patch+json")
		request.Header.Set("Accept", "application/ld+json")
		response = httptest.NewRecorder()
		diagnostics.Reset()
		handler.ServeHTTP(response, request)
		if response.Code != 200 || !strings.Contains(response.Body.String(), `"title":"Updated"`) || diagnostics.Len() != 0 {
			t.Fatalf("cycle %d PATCH: status=%d body=%q diagnostics=%q", i, response.Code, response.Body.String(), diagnostics.String())
		}
		request = httptest.NewRequest("DELETE", "http://example.test"+cycleResource.ID, nil)
		request.Header.Set("Accept", "application/ld+json")
		response = httptest.NewRecorder()
		diagnostics.Reset()
		handler.ServeHTTP(response, request)
		if response.Code != 204 || response.Body.Len() != 0 || diagnostics.Len() != 0 {
			t.Fatalf("cycle %d DELETE: status=%d body=%q diagnostics=%q", i, response.Code, response.Body.String(), diagnostics.String())
		}
	}
	dbOutput.Reset()
	dbErrors.Reset()
	if code := php.RunWithVersion(verify, "verify.php", &dbOutput, &dbErrors, php.PHP84); code != 0 || dbOutput.String() != "1" || dbErrors.Len() != 0 {
		t.Fatalf("database rows after repeated CRUD: exit=%d stdout=%q stderr=%q", code, dbOutput.String(), dbErrors.String())
	}
	runtime.GC()
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	requestCount := 96
	if configured := os.Getenv("GOPHP_HTTP_STRESS_REQUESTS"); configured != "" {
		count, err := strconv.Atoi(configured)
		if err != nil || count < 96 || count > 2000 {
			t.Fatalf("GOPHP_HTTP_STRESS_REQUESTS must be between 96 and 2000: %q", configured)
		}
		requestCount = count
	}
	started := time.Now()
	for i := 0; i < requestCount; i++ {
		path, want := fmt.Sprintf("/api/probe?round=%d", i), `{"ok":true}`
		if i%8 == 0 {
			path, want = fmt.Sprintf("/api/books-manual/%d", created.ID), `"title":"Alpha"`
		}
		request = httptest.NewRequest("GET", "http://example.test"+path, nil)
		response = httptest.NewRecorder()
		diagnostics.Reset()
		handler.ServeHTTP(response, request)
		if response.Code != 200 || !strings.Contains(response.Body.String(), want) || diagnostics.Len() != 0 {
			t.Fatalf("sustained GET %d: status=%d body=%q diagnostics=%q", i, response.Code, response.Body.String(), diagnostics.String())
		}
	}
	runtime.GC()
	runtime.ReadMemStats(&after)
	t.Logf("sustained requests=%d duration=%s heap_alloc_before=%d heap_alloc_after=%d", requestCount, time.Since(started), before.HeapAlloc, after.HeapAlloc)
	if requestCount > 96 && after.HeapAlloc > before.HeapAlloc+64*1024*1024 {
		t.Fatalf("heap retained more than 64 MiB after %d requests: before=%d after=%d", requestCount, before.HeapAlloc, after.HeapAlloc)
	}
	var concurrentDiagnostics synchronizedBuffer
	parallelHandler := phpHandler(string(source), filename, php.PHP84, &concurrentDiagnostics)
	parallel := make([]*httptest.ResponseRecorder, 8)
	var requests sync.WaitGroup
	for i := range parallel {
		requests.Add(1)
		go func(i int) {
			defer requests.Done()
			result := httptest.NewRecorder()
			parallelHandler.ServeHTTP(result, httptest.NewRequest("GET", fmt.Sprintf("http://example.test/api/probe?parallel=%d", i), nil))
			parallel[i] = result
		}(i)
	}
	requests.Wait()
	for i, result := range parallel {
		if result == nil || result.Code != 200 || result.Body.String() != `{"ok":true}` {
			t.Fatalf("parallel GET %d: result=%v", i, result)
		}
	}
	if diagnostic := concurrentDiagnostics.String(); diagnostic != "" {
		t.Fatalf("parallel GET diagnostics=%q", diagnostic)
	}
}

func TestCLI(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.php")
	if err := os.WriteFile(file, []byte("<?php echo 'ok';"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name            string
		args            []string
		exit            int
		out, diagnostic string
	}{
		{"success", []string{file}, 0, "ok", ""},
		{"usage", nil, 2, "", "usage:"},
		{"extra arguments", []string{file, file}, 0, "ok", ""},
		{"missing file", []string{filepath.Join(dir, "missing.php")}, 1, "", "missing.php"},
		{"version 8.3", []string{"--php-version=8.3", file}, 0, "ok", ""},
		{"version 8.4", []string{"--php-version", "8.4", file}, 0, "ok", ""},
		{"version 8.5", []string{"--php-version=8.5", file}, 0, "ok", ""},
		{"unsupported version", []string{"--php-version=8.6", file}, 2, "", "unsupported PHP version"},
		{"empty version", []string{"--php-version=", file}, 2, "", "unsupported PHP version"},
		{"missing version", []string{"--php-version"}, 2, "", "flag needs an argument"},
		{"unknown option", []string{"--unknown", file}, 2, "", "flag provided but not defined"},
		{"help", []string{"--help"}, 0, "", "usage:"},
	} {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := run(test.args, &stdout, &stderr)
			if code != test.exit || stdout.String() != test.out {
				t.Fatalf("exit=%d stdout=%q stderr=%q", code, &stdout, &stderr)
			}
			if test.diagnostic == "" && stderr.Len() != 0 || !strings.Contains(stderr.String(), test.diagnostic) {
				t.Errorf("unexpected stderr %q", &stderr)
			}
		})
	}
}

func TestCLIArguments(t *testing.T) {
	file := filepath.Join(t.TempDir(), "args.php")
	if err := os.WriteFile(file, []byte("<?php echo $argc, '|', $argv[1], '|', $argv[2], '|', $_SERVER['argc'];"), 0600); err != nil {
		t.Fatal(err)
	}
	var out, diagnostic bytes.Buffer
	code := run([]string{file, "alpha", "beta"}, &out, &diagnostic)
	if code != 0 || out.String() != "3|alpha|beta|3" || diagnostic.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, &out, &diagnostic)
	}
}

// The warning text and line number are from the user's PHP 8.5.11 comparison.
// Exercise the same CLI entry point used by compare-all, including its default.
func TestCLINaNVersion(t *testing.T) {
	file, err := filepath.Abs("../../tests/fixtures/success/truth_true.php")
	if err != nil {
		t.Fatal(err)
	}
	for _, version := range []string{"8.5", "8.4", "8.3", ""} {
		t.Run("target="+version, func(t *testing.T) {
			args := []string{file}
			if version != "" {
				args = append([]string{"--php-version=" + version}, args...)
			}
			var stdout, stderr bytes.Buffer
			code := run(args, &stdout, &stderr)
			wantWarning := ""
			if version == "8.5" {
				wantWarning = "Warning: unexpected NAN value was coerced to bool in " + file + " on line 21\n"
			}
			if code != 0 || stdout.String() != strings.Repeat("1", 20)+"\n" || stderr.String() != wantWarning {
				t.Fatalf("exit=%d stdout=%q stderr=%q; want stderr=%q", code, &stdout, &stderr, wantWarning)
			}
		})
	}
}

func TestCLIV1ProgramAndDiagnostics(t *testing.T) {
	source, err := os.ReadFile("../../tests/fixtures/success/v1_report.php")
	if err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(t.TempDir(), "集計 report.php")
	if err = os.WriteFile(file, source, 0600); err != nil {
		t.Fatal(err)
	}
	for _, version := range []string{"8.3", "8.4", "8.5"} {
		var out, diagnostics bytes.Buffer
		code := run([]string{"--php-version=" + version, "--", file}, &out, &diagnostics)
		if code != 0 || out.String() != "a:tea:12:3\nc:水:3:3\n2:2:1:0:integer\n" || diagnostics.Len() != 0 {
			t.Fatalf("%s exit=%d out=%q err=%q", version, code, &out, &diagnostics)
		}
	}
	if err = os.WriteFile(file, []byte("<?php\necho 'before';\ncount(null);\necho 'after';"), 0600); err != nil {
		t.Fatal(err)
	}
	var out, diagnostics bytes.Buffer
	code := run([]string{file}, &out, &diagnostics)
	want := "Fatal error: Uncaught TypeError: count(): Argument #1 ($value) must be of type Countable|array, null given in " + file + ":3\nStack trace:\n#0 {main}\n  thrown in " + file + " on line 3\n"
	if code != 255 || out.String() != "before" || diagnostics.String() != want {
		t.Fatalf("exit=%d out=%q err=%q", code, &out, &diagnostics)
	}
}
