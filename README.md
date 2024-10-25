# Scootin' around

Technical task for Saily. BE Service that exposes a REST-like API intended for scooter event collecting and reporting to mobile clients.

## Prerequisites

Install Golang

```bash
brew install go
```

## Running locally

```bash
# Run
make start
# Stop
make stop
```

## Local development

Install `air`

```
go install github.com/cosmtrek/air@latest
```

Run using Makefile

```
make run_watch
```

## Running tests

Run unit tests
(project has no logic, hence no unit tests)

Run e2e tests

```
make test_e2e
```

## Test manually

Open postman collection [here](/postman/Scootin_bout.postman_collection.json).

## Endpoint docs

Check endpoint docs using swagger in `http://localhost:8090/swagger/index.html`
