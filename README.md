# Fern JUnit Client

A CLI that can read JUnit test reports and send them to a Fern Platform instance in a format it understands

## Introduction

If you don't know what Fern is, [check it out here!](https://github.com/guidewire-oss/fern-platform)

## Install

To install the CLI, use the following command:

```sh
go install github.com/guidewire-oss/fern-junit-client@latest
```


## Registering Application with Fern Platform
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

### Configuration

The fern-junit-client can be configured using environment variables:

#### API Endpoint Configuration

| Environment Variable | Description | Default |
|---------------------|-------------|---------|
| `FERN_API_ENDPOINT_PATH` | Override the API endpoint path | `api/v1/test-runs` |

Example:
```sh
# Use a custom API endpoint path
export FERN_API_ENDPOINT_PATH="api/v2/test-results"
fern-junit-client send -u "http://localhost:8080" -p "77b34e74-5631-5a71-b8ce-97b9d6bab10a" -f "report.xml"
```

### OAuth Authentication

The fern-junit-client supports OAuth 2.0 authentication using the client credentials grant type. This allows the client to authenticate with the Fern backend using OAuth tokens.

#### OAuth Configuration

OAuth authentication is configured via environment variables:

| Environment Variable | Description | Required |
|---------------------|-------------|----------|
| `AUTH_URL` | The OAuth 2.0 token endpoint URL | Yes (to enable OAuth) |
| `FERN_AUTH_CLIENT_ID` | The OAuth client ID | Yes (if OAuth enabled) |
| `FERN_AUTH_CLIENT_SECRET` | The OAuth client secret/password | Yes (if OAuth enabled) |
| `FERN_CLIENT_SCOPE` | Space-separated list of OAuth scopes to request | No (optional) |

#### Behavior

- **OAuth Disabled**: If `AUTH_URL` is not set, the client will operate without authentication (backward compatible behavior).
- **OAuth Enabled**: If `AUTH_URL` is set, the client will:
  1. Validate that `FERN_AUTH_CLIENT_ID` and `FERN_AUTH_CLIENT_SECRET` are also provided
  2. Request an access token from the OAuth server using client credentials grant
  3. Include requested scopes in the token request if `FERN_CLIENT_SCOPE` is set
  4. Include the Bearer token in the Authorization header for all API calls to the Fern backend
  5. Automatically refresh the token when it expires

#### Examples with OAuth

##### Without OAuth (default)
```sh
fern-junit-client send -u "http://localhost:8080" -p "77b34e74-5631-5a71-b8ce-97b9d6bab10a" -f "report.xml"
```

##### With OAuth
```sh
export AUTH_URL="https://oauth.example.com/token"
export FERN_AUTH_CLIENT_ID="your-client-id"
export FERN_AUTH_CLIENT_SECRET="your-client-secret"

fern-junit-client send -u "http://localhost:8080" -p "77b34e74-5631-5a71-b8ce-97b9d6bab10a" -f "report.xml"
```

##### With OAuth and Scopes
```sh
export AUTH_URL="https://oauth.example.com/token"
export FERN_AUTH_CLIENT_ID="your-client-id"
export FERN_AUTH_CLIENT_SECRET="your-client-secret"
export FERN_CLIENT_SCOPE="read write admin"

fern-junit-client send -u "http://localhost:8080" -p "77b34e74-5631-5a71-b8ce-97b9d6bab10a" -f "report.xml"
```

##### Verbose Mode
Use the `--verbose` flag to see OAuth authentication status:
```sh
fern-junit-client send -u "http://localhost:8080" -p "77b34e74-5631-5a71-b8ce-97b9d6bab10a" -f "report.xml" --verbose
```

## See Also

* [Fern UI](https://github.com/guidewire-oss/fern-ui)
* [Fern Platform](https://github.com/guidewire-oss/fern-platform)
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
