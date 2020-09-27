# Bitlyke

The bit.ly like shortener.

## Requirements

* docker-compose

## How to run

``` bash
make run
```

After that the documentation will be available <http://localhost:8080>

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

* simple frontend backend

* authorization

* circle ci config
