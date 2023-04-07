# LinkedIn Go API Library (WIP)

This is a work in progress Go package for the LinkedIn API that provides functionalities for sharing content on LinkedIn, registering image uploads, and uploading images. The package has a `Client` struct that contains an access token and methods for sharing text and articles, registering image uploads, and uploading images.

**Note:** This package is still in development and is heavily far away from being production ready. Use it at your own risk.

## Installation

To use this package in your Go project, you can install it using the following command:

```
go get github.com/<username>/linkedin-go-api
```

## Usage

Here's an example of how to use the LinkedIn Go API library:

```go
import (
    "fmt"
    "github.com/<username>/linkedin-go-api"
)

func main() {
    accessToken := "your-access-token"
    client := linkedinapi.NewClient(accessToken)

    // Share text on LinkedIn
    personURN := "urn:li:person:<person-id>"
    text := "Hello, world!"
    shareID, err := client.ShareText(personURN, text)
    if err != nil {
        fmt.Printf("Failed to share text: %v\n", err)
    } else {
        fmt.Printf("Shared text with share ID: %v\n", shareID)
    }

    // Share article on LinkedIn
    url := "https://example.com/article"
    title := "Example Article"
    description := "This is an example article."
    articleShareID, err := client.ShareArticle(personURN, text, url, title, description)
    if err != nil {
        fmt.Printf("Failed to share article: %v\n", err)
    } else {
        fmt.Printf("Shared article with share ID: %v\n", articleShareID)
    }
}

## Contributing

Contributions are welcome! If you find any issues or bugs, please submit an issue or pull request.

## License

This package is licensed under the MIT license. See the LICENSE file for more details.
