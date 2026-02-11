package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  shorten <url>")
		fmt.Println("  expand <code>")
		fmt.Println("  list")
		return
	}

	switch os.Args[1] {
	case "shorten":
		if len(os.Args) < 3 {
			fmt.Println("Provide a URL")
			return
		}
		code, err := ShortenURL(os.Args[2])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Short code:", code)

	case "expand":
		if len(os.Args) < 3 {
			fmt.Println("Provide a code")
			return
		}
		url, err := ExpandURL(os.Args[2])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Original URL:", url)

	case "list":
		ListURLs()

	default:
		fmt.Println("Unknown command")
	}
}
