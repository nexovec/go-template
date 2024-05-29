#!/bin/sh
go install github.com/cosmtrek/air@latest # hot reload
go install github.com/go-delve/delve/cmd/dlv@latest # debugger
go install github.com/a-h/templ/cmd/templ@latest # templating

# migrations, db access codegen
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest