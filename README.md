# gofetch

Project to learn a bare bones rest api with preference to standard library.

## Getting started

1. [Install go](https://go.dev/doc/install)
2. run `go mod tidy` --dependencies will be downloaded automatically.

```sh
go mod tidy
```
3. Create a `.env` file and adjust as desired:

```sh
cp .env.example .env
```

4. Run the application

```sh
make run/api
```

## TODO:

- [x] Server start with configurable port (switch to .env)
- [x] Logger
- [ ] Panic recovery middleware
- [x] Request ID middleware
- [ ] `GET /version` route
- [ ] `GET/HEAD /healthz` route
- [ ] `GET /hello-api-call` route
- [x] Add Air for hot reloading
- [x] Add dotenv to use environment variables instead of cli opts
- [x] Tests for middleware
- [x] Tests for routes
- [x] Graceful shutdown
