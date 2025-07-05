package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"sync"
	"time"
)

//go:embed static
var siteViewer embed.FS

type Election struct {
	ElectionUID string `json:"electionUId"`
	Name        string `json:"name"`
}

var (
	electionData  = make([]Election, 0) // Slice of Election objects
	uidToFileName = make(map[string]string)
	lock          sync.RWMutex

	//Configurable
	updateInterval  = 10 * time.Second //Update interval
	port            = ":20202"
	folderPath      = "."
	electionsFolder = "./elections"
	svgDirectory    = "./svgs" // Directory where SVG files are stored
)

// corsMiddleware adds CORS headers to the response
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight (OPTIONS) requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Set up periodic folder scanning
	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			scanElectionsFolder()
		}
	}()

	// Get the subdirectory containing just the static files
	staticFS, err := fs.Sub(siteViewer, "static")
	if err != nil {
		log.Fatal(err)
	}

	// Serve static files (svg-viewer app)
	http.Handle("/", http.FileServer(http.FS(staticFS)))

	// API endpoint to get SVGs
	http.HandleFunc("/api/svgs", handleSVGsRequest)

	// Initial scan
	scanElectionsFolder()

	// Set up HTTP routes
	http.HandleFunc("/elections", getElectionsHandler)
	http.HandleFunc("/election", getElectionFileHandler)

	// Wrap handlers with CORS middleware
	handler := corsMiddleware(http.DefaultServeMux)

	// Start the server
	fmt.Println("Starting server on " + port)
	if err := http.ListenAndServe(port, handler); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
