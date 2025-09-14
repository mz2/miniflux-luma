package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/feeds"
	"miniflux.app/client"
)

var miniflux *client.Client
var minifluxEndpoint string
var feedTitle string
var feedFormat string = "atom" // default format

// cleanContent removes malformed fmt.Printf artifacts from content
func cleanContent(content string) string {
	// Remove malformed printf patterns like %!&(MISSING)
	re := regexp.MustCompile(`%![^\s]*\(MISSING\)`)
	content = re.ReplaceAllString(content, "")
	// Also clean up standalone %! patterns
	re2 := regexp.MustCompile(`%![^\s<>"]*`)
	content = re2.ReplaceAllString(content, "")
	return content
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	// Determine format from path or use default
	format := feedFormat
	if strings.HasSuffix(r.URL.Path, "/rss") {
		format = "rss"
	} else if strings.HasSuffix(r.URL.Path, "/atom") {
		format = "atom"
	}

	w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; base-uri 'none'; form-action 'none'")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	// Get new entries
	entries, err := miniflux.Entries(&client.Filter{
		Limit:     10,
		Order:     "published_at",
		Direction: "desc",
		Starred:   "1",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create atom feed
	feed := &feeds.Feed{
		Title:   feedTitle,
		Link:    &feeds.Link{Href: minifluxEndpoint},
		Created: time.Now(),
		Items:   []*feeds.Item{},
	}
	for _, entry := range entries.Entries {
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       entry.Title,
			Link:        &feeds.Link{Href: entry.URL},
			Description: cleanContent(entry.Content),
			Author:      &feeds.Author{Name: entry.Author},
			Created:     entry.Date,
		})
	}

	// Generate feed in requested format
	var output string
	if format == "rss" {
		rss, err := feed.ToRss()
		if err != nil {
			log.Fatal(err)
		}
		output = rss
	} else {
		atom, err := feed.ToAtom()
		if err != nil {
			log.Fatal(err)
		}
		output = atom
	}
	fmt.Fprintf(w, output)
}

func main() {
	APITokenFile := ""
	APIToken := ""
	formatFile := ""
	listenAddress := ""
	certFile := ""
	keyFile := ""

	// Read command line arguments
	flag.StringVar(&minifluxEndpoint, "endpoint", "https://miniflux.example.org", "Miniflux server endpoint")
	flag.StringVar(&APITokenFile, "api-token-file", "api_token", "Load Miniflux API token from file")
	flag.StringVar(&listenAddress, "listen-addr", "127.0.0.1:8080", "Listen on this address")
	flag.StringVar(&feedTitle, "feed-title", "Starred entries", "Title of the feed")
	flag.StringVar(&formatFile, "format-file", "", "Load feed format from file (atom or rss)")
	flag.StringVar(&certFile, "tls-cert", "", "TLS certificate file path (skip to disable TLS)")
	flag.StringVar(&keyFile, "tls-key", "", "TLS key file path (skip to disable TLS)")
	flag.Parse()

	// Load API token
	dat, err := os.ReadFile(APITokenFile)
	if err != nil {
		log.Fatal(err)
	}
	APIToken = strings.TrimSpace(string(dat))

	// Load format if file specified
	if formatFile != "" {
		formatData, err := os.ReadFile(formatFile)
		if err == nil {
			format := strings.TrimSpace(string(formatData))
			if format == "rss" || format == "atom" {
				feedFormat = format
			}
		}
	}

	// Authentication using API token then fetch starred items
	miniflux = client.New(minifluxEndpoint, APIToken)

	// Start web server
	http.HandleFunc("/", httpHandler)
	http.HandleFunc("/rss", httpHandler)
	http.HandleFunc("/atom", httpHandler)
	log.Printf("Listening on %s (format: %s)\n", listenAddress, feedFormat)
	if certFile != "" && keyFile != "" {
		log.Fatal(http.ListenAndServeTLS(listenAddress, certFile, keyFile, nil))
	} else {
		log.Fatal(http.ListenAndServe(listenAddress, nil))
	}
}
