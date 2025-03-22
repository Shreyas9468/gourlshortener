package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// URLShortener manages the mapping between short keys and original URLs
type URLShortener struct {
	urls map[string]string
}

// NewURLShortener creates a new instance of URLShortener
func NewURLShortener() *URLShortener {
	return &URLShortener{
		urls: make(map[string]string),
	}
}

// HandleShorten processes URL shortening requests
func (us *URLShortener) HandleShorten(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed - use POST", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get and validate original URL
	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	// Generate short key and store mapping
	shortKey := generateShortKey()
	us.urls[shortKey] = originalURL
	shortenedURL := fmt.Sprintf("http://localhost:8080/short/%s", shortKey)

	// Set response headers
	w.Header().Set("Content-Type", "text/html")

	// Generate HTML response
	responseHTML := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>URL Shortener</title>
</head>
<body>
    <h2>URL Shortener</h2>
    <p>Original URL: %s</p>
    <p>Shortened URL: <a href="%s">%s</a></p>
    <form method="post" action="/shorten">
        <input type="text" name="url" placeholder="Enter a URL" required>
        <input type="submit" value="Shorten">
    </form>
</body>
</html>
`, originalURL, shortenedURL, shortenedURL)

	// Write response
	if _, err := fmt.Fprint(w, responseHTML); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

// HandleRedirect processes URL redirection requests
func (us *URLShortener) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	// Extract short key from URL path
	shortKey := r.URL.Path[len("/short/"):]
	if shortKey == "" {
		http.Error(w, "Short key is missing", http.StatusNotFound)
		return
	}

	// Look up original URL
	originalURL, exists := us.urls[shortKey]
	if !exists {
		http.Error(w, "Shortened URL not found", http.StatusNotFound)
		return
	}

	// Perform redirect
	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

// generateShortKey creates a random 6-character key
func generateShortKey() string {
	// Seed random number generator once at startup would be better,
	// but keeping it here for simplicity as per original code
	rand.Seed(time.Now().UnixNano())
	
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6
	
	result := make([]byte, keyLength)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func main() {
	// Initialize shortener
	shortener := NewURLShortener()

	// Register handlers
	http.HandleFunc("/shorten", shortener.HandleShorten)
	http.HandleFunc("/short/", shortener.HandleRedirect)

	// Start server
	const port = ":8080"
	fmt.Printf("URL Shortener is running on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}