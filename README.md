# Lago Go Client
## Installation

```shell
go get github.com/getlago/lago-go-client/v1
```

## Usage

```golang
import "github.com/getlago/lago-go-client/v1"

func main() {
  lagoClient := lago.New().
    SetApiKey("MY_API_KEY").
    // SetDebug will log the RAW request and RAW response
    SetDebug(true)

  // If you want to use your API url for self hosted version
  lagoClient := lago.New().
    SetBaseURL("https://my.url").
    SetApiKey("MY_API_KEY")
}
```

## Development

- Fork the repository
- Open a Pull Request 

## Documentation
The Lago documentation is available at doc.getlago.com.

## Contributing
The contribution documentation is available here

## License
Lago GO client is distributed under AGPL-3.0.