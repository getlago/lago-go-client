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

### Verifying Standard Webhooks deliveries

A webhook endpoint created with `SignatureAlgo: lago.StandardWebhooks` signs its deliveries as the [Standard Webhooks](https://www.standardwebhooks.com) specification describes. Its `SigningKeys` hold the `whsec_` secrets, returned when you create, update or get the endpoint with an API key that has write access.

`lago.VerifyStandardWebhook` checks a delivery with one of these secrets. Pass it the raw request body bytes, before any JSON decoding. It rejects a delivery whose `webhook-timestamp` is more than 5 minutes in the past or in the future.

```go
package main

import (
	"io"
	"log"
	"net/http"

	lago "github.com/getlago/lago-go-client"
)

func main() {
	secret := "whsec_..." // a SecretKey from the endpoint's SigningKeys

	http.HandleFunc("/webhooks/lago", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "cannot read body", http.StatusBadRequest)
			return
		}

		if err := lago.VerifyStandardWebhook(secret, r.Header, body); err != nil {
			http.Error(w, "invalid signature", http.StatusBadRequest)
			return
		}

		// The delivery is from Lago: decode body and handle the event.
		w.WriteHeader(http.StatusOK)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
```

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

Run the test suite:

```shell
go test ./...
```

### Code Quality

Format code:

```shell
go fmt ./...
```

Run linting (requires golangci-lint):

```shell
# Install golangci-lint
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.6.2

# Run linter
golangci-lint run
```

## Documentation

The Lago documentation is available at [doc.getlago.com](https://doc.getlago.com/api-reference/intro).

## Contributing

The contribution documentation is available [here](https://github.com/getlago/lago-go-client/blob/main/CONTRIBUTING.md)

## License

Lago Go client is distributed under [MIT license](LICENSE).
