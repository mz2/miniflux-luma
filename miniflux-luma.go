package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/feeds"
	"miniflux.app/client"
)

var miniflux *client.Client
var minifluxEndpoint string
var feedTitle string
var feedFormat string = "atom" // default format

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
			Description: entry.Content,
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
	APIToken := ""
	listenAddress := ""
	certFile := ""
	keyFile := ""

	// Read command line arguments
	flag.StringVar(&minifluxEndpoint, "endpoint", "https://miniflux.example.org", "Miniflux server endpoint")
	flag.StringVar(&APIToken, "api-token", "", "Miniflux API token")
	flag.StringVar(&listenAddress, "listen-addr", "127.0.0.1:8080", "Listen on this address")
	flag.StringVar(&feedTitle, "feed-title", "Starred entries", "Title of the feed")
	flag.StringVar(&feedFormat, "format", "atom", "Feed format (atom or rss)")
	flag.StringVar(&certFile, "tls-cert", "", "TLS certificate file path (skip to disable TLS)")
	flag.StringVar(&keyFile, "tls-key", "", "TLS key file path (skip to disable TLS)")
	flag.Parse()

	// Validate API token
	if APIToken == "" {
		log.Fatal("API token is required. Use -api-token flag or set via environment")
	}
	APIToken = strings.TrimSpace(APIToken)

	// Validate format
	if feedFormat != "atom" && feedFormat != "rss" {
		log.Fatal("Invalid format. Must be 'atom' or 'rss'")
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
