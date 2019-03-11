# AsyncAPI watcher

Builds an [asyncapi](https://www.asyncapi.com/) documentation for your microservices
communicating through [rabbitmq](https://www.rabbitmq.com/).

It listens to all published amqp messages and keeps an updated asyncapi
compliant documenation served at `/asyncapi`.

## Roadmap

- [x] Add info and server sections to spec to make it valid asyncapi.
- [ ] Extract configuration to file.
- [ ] Support consuming from multiple configurable exchanges.
- [ ] Use postgres as database.
- [ ] Add CI with github actions.
- [ ] Build and publish docker image.

## Running locally

Install dependencies
```
go mod download
```

### Running tests

Start the rabbitmq server

```
docker-compose up -d
```

Run the tests recursively for all subpackages

```
go test ./...
```
