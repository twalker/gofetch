# gofetch

Project to learn a bare bones rest api with preference to standard library.

## Getting started

1. [Install go](https://go.dev/doc/install)
2. run `make build/api` or `make run/api` --dependencies will be downloaded automatically.

```sh
make run/api
```

## TODO:

- [ ] Server start with configurable port (switch to .env)
- [x] Logger
- [ ] Panic recovery middleware
- [x] Request ID middleware
- [ ] `GET /version` route
- [ ] `GET/HEAD /healthz` route
- [ ] `GET /hello-api-call` route
- [x] Add Air for hot reloading
- [ ] Add dotenv to use environment variables instead of cli opts
- [x] Tests for middleware
- [x] Tests for routes
- [x] Graceful shutdown
