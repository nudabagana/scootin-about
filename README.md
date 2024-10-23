# Scootin' around

Technical task for Saily. BE Service that exposes a REST-like API intended for scooter event collecting and reporting to mobile clients.

## Prerequisites

Install Golang

```bash
brew install go
```

## Running locally

TBD

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

```
make test
```

Run integration tests

```
make test_integration
```
