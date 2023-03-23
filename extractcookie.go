package extractcookie

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//nolint:gochecknoglobals // TODO exchange for traefik log when available
var (
	LoggerWARN  = log.New(ioutil.Discard, "WARN:  extractcookie: ", log.Ldate|log.Ltime|log.Lshortfile)
	LoggerINFO  = log.New(ioutil.Discard, "INFO:  extractcookie: ", log.Ldate|log.Ltime|log.Lshortfile)
	LoggerDEBUG = log.New(ioutil.Discard, "DEBUG: extractcookie: ", log.Ldate|log.Ltime|log.Lshortfile)
)

type Config struct {
	LogLevel                 string `yaml:"logLevel"`
	CookieName               string `yaml:"cookieName"`
	HeaderNameForCookieValue string `yaml:"headerNameForCookieValue"`
	CookieValuePrefix        string `yaml:"cookieValuePrefix"`
}

func CreateConfig() *Config {
	return &Config{
		CookieName:               "",
		HeaderNameForCookieValue: "Authorization",
		CookieValuePrefix:        "Bearer ",
		LogLevel:                 "info",
	}
}

type ExtractCookie struct {
	next   http.Handler
	config *Config
	name   string
}

func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	switch config.LogLevel {
	case "warn":
		LoggerWARN.SetOutput(os.Stdout)
	case "info":
		LoggerWARN.SetOutput(os.Stdout)
		LoggerINFO.SetOutput(os.Stdout)
	case "debug":
		LoggerWARN.SetOutput(os.Stdout)
		LoggerINFO.SetOutput(os.Stdout)
		LoggerDEBUG.SetOutput(os.Stdout)
	default:
		return nil, fmt.Errorf("ERROR: extractcookie: %s", config.LogLevel)
	}

	if config.CookieName == "" {
		return nil, fmt.Errorf("ERROR: extractcookie: cookie name can not be empty")
	}

	return &ExtractCookie{
		next:   next,
		name:   name,
		config: config,
	}, nil
}

func (t *ExtractCookie) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	LoggerDEBUG.Printf("ServeHTTP started")

	cookieName := t.config.CookieName
	headerNameForCookieValue := t.config.HeaderNameForCookieValue
	cookieValuePrefix := t.config.CookieValuePrefix

	LoggerDEBUG.Printf("got Cookie header value = %s", req.Header.Get("Cookie"))
	LoggerDEBUG.Printf("extracting %s cookie", cookieName)

	cookie, err := req.Cookie(cookieName)
	if err != nil {
		LoggerWARN.Printf("tries to extract cookie that not exists")
	} else {
		LoggerDEBUG.Printf("extracted cookie value = %s", cookie.Value)
		headerValue := cookieValuePrefix + cookie.Value
		req.Header.Set(headerNameForCookieValue, headerValue)
		LoggerINFO.Printf("set %s header to %s", headerNameForCookieValue, headerValue)
	}

	t.next.ServeHTTP(rw, req)
}
