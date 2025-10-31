# Lago Go Client

This is a Go wrapper for Lago API

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://spdx.org/licenses/MIT.html)

## Current Releases

| Project            | Release Badge                                                                                                                                   |
| ------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------- |
| **Lago**           | [![Lago Release](https://img.shields.io/github/v/release/getlago/lago)](https://github.com/getlago/lago/releases)                               |
| **Lago Go Client** | [![Lago Go Client Release](https://img.shields.io/github/v/release/getlago/lago-go-client)](https://github.com/getlago/lago-go-client/releases) |

## Installation

To use the client in your Go application:

```shell
go get github.com/getlago/lago-go-client@v1
```

## Usage

Once the package is installed, you can use it in your Go application as follows:

```go
package main

import (
	"context"
	"fmt"
	"log"

	lago "github.com/getlago/lago-go-client"
)

func main() {
	client := lago.New().SetApiKey("xyz")

	ctx := context.TODO()
	// Example: List customers
	billableMetrics, err := client.BillableMetric().GetList(ctx, &lago.BillableMetricListInput{
		Page:    1,
		PerPage: 10,
	})
	if err != nil {
		log.Fatalf("Error fetching Billable Metrics: %v", err)
	}

	fmt.Println("List of Billable Metrics:")
	for _, billableMetric := range billableMetrics.BillableMetrics {
		fmt.Printf("- %s\n", billableMetric.Name)
	}
}
```

For detailed usage, refer to the [lago API reference](https://doc.getlago.com/api-reference/intro).

## Development

### Prerequisites

- Go 1.18 or higher
- Git

### Setup

1. Clone the repository:

    ```shell
    git clone https://github.com/getlago/lago-go-client.git
    cd lago-go-client
    ```

2. Install dependencies:

    ```shell
    go mod download
    ```

### Testing

The repository includes two types of tests:

1. Unit tests which tests the Ruby code itself
2. Integration tests which tests the Ruby code against a running Lago instance.

    To run these tests, you need to set the `INTEGRATION_TESTS_ENABLED` environment variable to `true` and provide the `TEST_LAGO_API_URL` and `TEST_LAGO_API_KEY` environment variables. These variables are set by default when using the Docker Compose file.

#### Running the tests

Run the test suite:

```shell
go test ./...
INTEGRATION_TESTS_ENABLED=true TEST_LAGO_API_URL=http://lago:3000 TEST_LAGO_API_KEY=123456 go test ./testing/integration/...
```

or with Docker Compose:

```bash
docker compose exec client go test ./...
```

### Code Quality

Format code:

```shell
go fmt ./...
```

Run linting (requires golangci-lint):

```shell
# Install golangci-lint
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.6.0

# Run linter
golangci-lint run
```

or with Docker Compose:

```bash
docker compose exec client go fmt ./...
docker compose exec client golangci-lint run
```

## Documentation

The Lago documentation is available at [doc.getlago.com](https://doc.getlago.com/api-reference/intro).

## Contributing

The contribution documentation is available [here](https://github.com/getlago/lago-go-client/blob/main/CONTRIBUTING.md)

## License

Lago Go client is distributed under [MIT license](LICENSE).
