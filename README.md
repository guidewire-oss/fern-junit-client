# Fern JUnit Client

A CLI that can read JUnit test reports and send them to a Fern Reporter instance in a format it understands

## Introduction

If you don't know what Fern is, [check it out here!](https://github.com/guidewire-oss/fern-reporter)

## Install

To install the CLI, use the following command:

```sh
go install github.com/guidewire-oss/fern-junit-client@latest
```


## Registering Application with Fern-Reporter
```bash
curl -L -X POST http://localhost:8080/api/project \
  -H "Content-Type: application/json" \
  -d '{
    "name": "First Projects",
    "team_name": "my team",
    "comment": "This is the test project"
  }' 
```

Sample Response:
```json
{
  "uuid": "996ad860-2a9a-504f-8861-aeafd0b2ae29",
  "name": "First Projects",
  "team_name": "my team",
  "comment": "This is the test project"
}
```

## Usage

To see all available options, use `fern-junit-client help`

### Examples

#### Send Single Report

```sh
fern-junit-client send -u "http://localhost:8080" -p "77b34e74-5631-5a71-b8ce-97b9d6bab10a" -f "report.xml"
```

#### Send Multiple Reports

```sh
fern-junit-client send -u "http://localhost:8080" -p "77b34e74-5631-5a71-b8ce-97b9d6bab10a" -f "tests/*.xml"
```

## See Also

* [Fern UI](https://github.com/guidewire-oss/fern-ui)
* [Fern Reporter](https://github.com/guidewire-oss/fern-reporter)
* [Fern Ginkgo Client](https://github.com/guidewire-oss/fern-ginkgo-client)

## Development

To install the CLI locally for testing use the following command:

```sh
go build 

```

```sh
./fern-junit-client send -u "http://localhost:8080" -p "77b34e74-5631-5a71-b8ce-97b9d6bab10a" -f "report.xml"
```


### Executing Tests

To execute the tests, run `make test`

### Generating Test Static Files

To generate the static files used in the tests, run `make test-static-files`
