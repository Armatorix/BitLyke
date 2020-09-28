# Bitlyke

The bit.ly like shortener.

- [Bitlyke](#bitlyke)
  - [Requirements](#requirements)
  - [How to run](#how-to-run)
  - [Simple use case](#simple-use-case)
  - [Changes in api-spec](#changes-in-api-spec)
  - [Testing](#testing)
  - [Roadmap](#roadmap)

## Requirements

- docker-compose (v1.14.0+ - with compose file v3.3 support)

## How to run

``` bash
make run
```

After that the documentation will be available <http://localhost:80>

## Simple use case

1. Create new shortener with POST request on `localhost:8080/api` endpoint

2. Go to `localhost:8080/{short_path}` in the browser - it should redirect you to the real url that you provided.

## Changes in api-spec

After doing any chenges in api spec, the model rebuild should be done

``` bash
make rebuild
```

## Testing

To run e2e tests

``` bash
make test-e2e
```

## Roadmap

- terraform code with GCP/Heroku deployment

- simple frontend

- authorization

- circle ci config

- random shortener path generator
