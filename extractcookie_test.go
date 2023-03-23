package extractcookie_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/v-electrolux/extractcookie"
)

var (
	maxCookie = http.Cookie{
		Name:     "test_cookie_with_auth_token",
		Value:    "pretend_to_be_auth_token",
		Path:     "/path/to",
		Domain:   "google.com",
		Expires:  time.Now(),
		MaxAge:   5,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
	minCookie = http.Cookie{
		Name:  "test_cookie_with_auth_token",
		Value: "pretend_to_be_auth_token",
	}
	minCookieExpired = http.Cookie{
		Name:   "test_cookie_with_auth_token",
		Value:  "pretend_to_be_auth_token",
		MaxAge: 0,
	}
	maxCookieGarbage = http.Cookie{
		Name:     "some_cookie_garbage",
		Value:    "bla bla foo bar",
		Path:     "/path/to",
		Domain:   "google.com",
		Expires:  time.Now(),
		MaxAge:   5,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
	minCookieGarbage = http.Cookie{
		Name:  "some_cookie_garbage",
		Value: "bla bla foo bar",
	}
	minCookieExpiredGarbage = http.Cookie{
		Name:   "some_cookie_garbage",
		Value:  "bla bla foo bar",
		MaxAge: 0,
	}
)

type TestHttpResData struct {
	cfgCookieName               string
	cfgHeaderNameForCookieValue string
	cfgCookieValuePrefix        string
	cookies                     []http.Cookie
	expectedHeaderValue         string
}

func TestEmptyCookie(t *testing.T) {
	data := TestHttpResData{
		cfgCookieName:               "test_cookie_with_auth_token",
		cfgHeaderNameForCookieValue: "Authorization",
		cfgCookieValuePrefix:        "Bearer ",
	}

	t.Run("empty cookie", func(t *testing.T) {
		data.cookies = []http.Cookie{}
		data.expectedHeaderValue = ""
		testHttpRequest(t, data)
	})
}

func TestOneCookie(t *testing.T) {
	data := TestHttpResData{
		cfgCookieName:               "test_cookie_with_auth_token",
		cfgHeaderNameForCookieValue: "Authorization",
		cfgCookieValuePrefix:        "Bearer ",
	}

	t.Run("one maxCookie", func(t *testing.T) {
		data.cookies = []http.Cookie{maxCookie}
		data.expectedHeaderValue = "Bearer pretend_to_be_auth_token"
		testHttpRequest(t, data)
	})
	t.Run("one minCookie", func(t *testing.T) {
		data.cookies = []http.Cookie{minCookie}
		data.expectedHeaderValue = "Bearer pretend_to_be_auth_token"
		testHttpRequest(t, data)
	})
	t.Run("one minCookieExpired", func(t *testing.T) {
		data.cookies = []http.Cookie{minCookieExpired}
		data.expectedHeaderValue = "Bearer pretend_to_be_auth_token"
		testHttpRequest(t, data)
	})
	t.Run("one maxCookieGarbage", func(t *testing.T) {
		data.cookies = []http.Cookie{maxCookieGarbage}
		data.expectedHeaderValue = ""
		testHttpRequest(t, data)
	})
	t.Run("one minCookieGarbage", func(t *testing.T) {
		data.cookies = []http.Cookie{minCookieGarbage}
		data.expectedHeaderValue = ""
		testHttpRequest(t, data)
	})
	t.Run("one minCookieExpiredGarbage", func(t *testing.T) {
		data.cookies = []http.Cookie{minCookieExpiredGarbage}
		data.expectedHeaderValue = ""
		testHttpRequest(t, data)
	})
}

func TestTwoCookie(t *testing.T) {
	data := TestHttpResData{
		cfgCookieName:               "test_cookie_with_auth_token",
		cfgHeaderNameForCookieValue: "Authorization",
		cfgCookieValuePrefix:        "Bearer ",
	}

	t.Run("minCookieGarbage and maxCookieGarbage", func(t *testing.T) {
		data.cookies = []http.Cookie{minCookieGarbage, maxCookieGarbage}
		data.expectedHeaderValue = ""
		testHttpRequest(t, data)
	})
	t.Run("minCookieExpiredGarbage and maxCookieGarbage", func(t *testing.T) {
		data.cookies = []http.Cookie{minCookieExpiredGarbage, maxCookieGarbage}
		data.expectedHeaderValue = ""
		testHttpRequest(t, data)
	})
	t.Run("minCookie and maxCookieGarbage", func(t *testing.T) {
		data.cookies = []http.Cookie{minCookie, maxCookieGarbage}
		data.expectedHeaderValue = "Bearer pretend_to_be_auth_token"
		testHttpRequest(t, data)
	})
	t.Run("minCookieGarbage and maxCookie", func(t *testing.T) {
		data.cookies = []http.Cookie{minCookieGarbage, maxCookie}
		data.expectedHeaderValue = "Bearer pretend_to_be_auth_token"
		testHttpRequest(t, data)
	})
	t.Run("maxCookieGarbage and minCookie", func(t *testing.T) {
		data.cookies = []http.Cookie{maxCookieGarbage, minCookie}
		data.expectedHeaderValue = "Bearer pretend_to_be_auth_token"
		testHttpRequest(t, data)
	})
}

func testHttpRequest(t *testing.T, data TestHttpResData) {
	t.Helper()

	cfg := extractcookie.CreateConfig()
	cfg.CookieName = data.cfgCookieName
	cfg.HeaderNameForCookieValue = data.cfgHeaderNameForCookieValue
	cfg.CookieValuePrefix = data.cfgCookieValuePrefix
	cfg.LogLevel = "warn"
	backendStatusCode := http.StatusOK
	backendBody := []byte("foo bar")

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(backendStatusCode)
		rw.Write(backendBody)
	})

	handler, err := extractcookie.New(ctx, next, cfg, "tlsclientcertforward")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)

	for _, cookie := range data.cookies {
		req.AddCookie(&cookie)
	}

	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)
	resp := recorder.Result()

	assertStatusCode(t, resp, backendStatusCode)
	assertBody(t, resp, backendBody)
	assertHeader(t, req, data.cfgHeaderNameForCookieValue, data.expectedHeaderValue)
}

func assertStatusCode(t *testing.T, res *http.Response, expected int) {
	t.Helper()

	got := res.StatusCode
	if got != expected {
		t.Errorf("expected status code value: `%d`, got value: `%d`", expected, got)
	}
}

func assertBody(t *testing.T, res *http.Response, expected []byte) {
	t.Helper()

	got, _ := io.ReadAll(res.Body)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected body value: `%s`, got value: `%s`", expected, got)
	}
}

func assertHeader(t *testing.T, req *http.Request, key string, expected string) {
	t.Helper()

	got := req.Header.Get(key)
	if got != expected {
		t.Errorf("expected header %s value: `%s`, got value: `%s`", key, expected, got)
	}
}
