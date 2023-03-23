# Cookie value extractor plugin for traefik

[![Build Status](https://github.com/v-electrolux/extractcookie/workflows/Main/badge.svg?branch=master)](https://github.com/v-electrolux/extractcookie/actions)

Takes specified cookie value (from Cookie header)
and put it in a header you want. You can declare prefix for value,
which will be added to cookie value.
If there is no such cookie, the header value will not be set 

## Configuration

### Fields meaning
- `cookieName`: name of cookie, its value will be extracted. 
   No default, mandatory field
- `headerNameForCookieValue`: name of header, in which will be put cookie value.
  Default is 'Authorization'
- `cookieValuePrefix`: string prefix, that will be added to cookie value.
  Default is 'Bearer '
- `logLevel`: `warn`, `info` or `debug`. Default is `info`

### Static config examples

- cli as local plugin
```
--experimental.localplugins.extractcookie=true
--experimental.localplugins.extractcookie.modulename=github.com/v-electrolux/extractcookie
```

- envs as local plugin
```
TRAEFIK_EXPERIMENTAL_LOCALPLUGINS_extractcookie=true
TRAEFIK_EXPERIMENTAL_LOCALPLUGINS_extractcookie_MODULENAME=github.com/v-electrolux/extractcookie
```

- yaml as local plugin
```yaml
experimental:
  localplugins:
     extractcookie:
      modulename: github.com/v-electrolux/extractcookie
```

- toml as local plugin
```toml
[experimental.localplugins.extractcookie]
    modulename = "github.com/v-electrolux/extractcookie"
```

### Dynamic config examples

- docker labels
```
traefik.http.middlewares.extractCookieMiddleware.plugin.extractcookie.cookieName=test_access_token
traefik.http.middlewares.extractCookieMiddleware.plugin.extractcookie.headerNameForCookieValue=Authorization
traefik.http.middlewares.extractCookieMiddleware.plugin.extractcookie.cookieValuePrefix=Bearer 
traefik.http.middlewares.extractCookieMiddleware.plugin.extractcookie.logLevel=warn
traefik.http.routers.extractCookieRouter.middlewares=extractCookieMiddleware
```

- yaml
```yml
http:

  routers:
    extractCookieRouter:
      rule: host(`demo.localhost`)
      service: backend
      entryPoints:
        - web
      middlewares:
        - extractCookieMiddleware

  services:
    backend:
      loadBalancer:
        servers:
          - url: 127.0.0.1:5000

  middlewares:
    extractCookieMiddleware:
      plugin:
        extractcookie:
          cookieName: test_access_token
          headerNameForCookieValue: Authorization
          cookieValuePrefix: 'Bearer '
          logLevel: warn
```
