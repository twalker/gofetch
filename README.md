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
- [x] Structured Logger
- [ ] Panic recovery middleware
- [ ] CORS middleware
- [x] Request ID middleware
- [ ] Middleware chain
- [ ] `GET /version` route
- [ ] `GET/HEAD /healthz` route
- [ ] `GET /hello-api-call` route
- [ ] API ClientRequest module/helper
- [ ] Handle 404
- [ ] Json read/write helpers
- [ ] Default to "application/json" responses
- [ ] Standardize http error responses
- [x] Add Air for hot reloading
- [x] Add dotenv to use environment variables instead of cli opts
- [x] Tests for middleware
- [x] Tests for routes
- [x] Graceful shutdown
- [x] Make file with build, test, etc.
