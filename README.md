# Lago Go Client

## Usage

```golang
import "github.com/getlago/lago-go-client/v1"

func main() {
  lagoClient := lago.New().
    SetApiKey("MY_API_KEY")

  // If you want to use your API url for self hosted version
  lagoClient := lago.New().
    SetBaseURL("https://my.url").
    SetApiKey("MY_API_KEY")


  lagoClient.Customer().
    Create(&lago.CustomerInput{
      CustomerID: "vincent_12345",
      Name: "Vincento",
      Email: "vincent@toto.com",
    })
}
```