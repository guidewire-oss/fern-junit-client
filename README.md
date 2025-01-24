# Fern JUnit Client

A CLI that can read JUnit test reports and send them to a Fern Reporter instance in a format it understands

## Introduction

If you don't know what Fern is, [check it out here!](https://github.com/guidewire-oss/fern-reporter)

## Install

To install the CLI, use the following command:

```sh
go install github.com/guidewire-oss/fern-junit-client@latest
```

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

## Github Action

### Inputs

| Name                 | Required | Type    | Default | Description                                                                        |
|----------------------|----------|---------|---------|------------------------------------------------------------------------------------|
| url                  | true     | string  |         | URL of Fern reporter                                                               |
| file-pattern         | true     | string  |         | Directory or pattern where JUnit reports get generated (accepts `*` as a wildcard) |
| project-name         | true     | string  |         | Name of the project to display under in Fern                                       |
| tags                 | false    | string  |         | Comma-separated tags to associate with the suite run                               |
| verbose              | false    | boolean | false   | Whether to log verbosely (**note**: may increase log size significantly)           |
| generate-job-summary | false    | boolean | true    | Whether to generate a job summary with passed/failed/skipped counts                |

### Outputs

| Name          | Type    | Description                       |
|---------------|---------|-----------------------------------|
| tests-passed  | integer | Number of tests that passed       |
| tests-failed  | integer | Number of tests that failed       |
| tests-skipped | integer | Number of tests that were skipped |

### Example Workflow

```yaml
name: 'Fern JUnit Client Workflow'
on:
  push:
    branches:
      - main
jobs:
  fern-junit-client:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v3
      - name: Run Tests
          ...
      - name: Report tests to fern
        uses: guidewire-oss/fern-junit-client@latest
        with:
          url: 'https://fern.mydomain.com'
          file-pattern: 'tests/*.xml'
          project-name: 'My Service'
          tags: 'cpu,'
          verbose: false
          generate-job-summary: true
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
