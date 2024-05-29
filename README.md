# Go web server template

> Author: [@nexovec](github.com/nexovec)

> Uses: APIs, glorious monolith applications, high performance HTTP

This is my personal HTTP server template. Since I'm responsible for making more than one project with a common structure, a template is a nice to have

> BONUS: this is available as open source.

## Prerequisites

You need to have go 1.22 and docker installed. Primarily, only linux is supported. Recommended vs code extensions are: Git lens, go nightly, error lens, intellicode, templ lsp. Additionally, installation of graphviz is recommended.

## Usage

> TODO: move things from `air.toml` `pre_cmd` into here to speed up builds

You can launch the project with `docker compose up --build`.

## The tech

### Postgres

This specifically uses postgres for its good sqlc support, lots of available [extensions](https://pgt.dev) and libraries, and its popularity. SQLite should be considered for simple CRUD apps.

### ORM & Migrations

We use basic [go migrate](https://github.com/golang-migrate/migrate) for db migrations.
We don't use ORMs, but we use [sqlc](https://sqlc.dev/) to generate access methods for raw SQL queries.

### Caching

We use the in-memory cache [Otter](https://github.com/maypok86/otter) for its excellent performance and simple API.

### Server tech

Framework of choice is [go fiber](https://gofiber.io/), for excellent performance and simple API.

We use [lodash for go](https://pkg.go.dev/github.com/samber/lo#Async2) to suplement what we don't have in the stdlib, this has the added benefit of being similar to javascripts equivalent.

[Air](https://github.com/cosmtrek/air) is in charge of hot-reloading.
