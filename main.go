package main

import (
	"fmt"
	"net/http"
	"os"
)

// HTTP handler to redirect shortcodes in browser
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[1:] // remove leading /
	url, err := ExpandURL(code) // use existing ExpandURL
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

func startWebServer() {
	http.HandleFunc("/", redirectHandler)
	fmt.Println("Web server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func main() {
	// If no CLI args, start web server
	if len(os.Args) == 1 {
		startWebServer()
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
