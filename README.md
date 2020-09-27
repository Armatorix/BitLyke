# Bitlyke

The bit.ly like shortener.

- [Bitlyke](#bitlyke)
  - [Requirements](#requirements)
  - [How to run](#how-to-run)
  - [Changes in api-spec](#changes-in-api-spec)
  - [Testing](#testing)
  - [Roadmap](#roadmap)

## Requirements

* docker-compose

## How to run

``` bash
make run
```

After that the documentation will be available <http://localhost:80>

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

* terraform code with GCP/Heroku deployment

* simple frontend

* authorization

* circle ci config

* random shortener path generator
