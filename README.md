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
TBA

## Usage

To see all available options, use `fern-junit-client help`

### Examples

#### Send Single Report

```sh
fern-junit-client send -u "http://localhost:8080" -p "MyService" -f "report.xml"
```

#### Send Multiple Reports

```sh
fern-junit-client send -u "http://localhost:8080" -p "MyService" -f "tests/*.xml"
```

## See Also

* [Fern UI](https://github.com/guidewire-oss/fern-ui)
* [Fern Reporter](https://github.com/guidewire-oss/fern-reporter)
* [Fern Ginkgo Client](https://github.com/guidewire-oss/fern-ginkgo-client)

## Development

### Executing Tests

To execute the tests, run `make test`

### Generating Test Static Files

To generate the static files used in the tests, run `make test-static-files`
