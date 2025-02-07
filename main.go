package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"sync"
	"time"
)

type Election struct {
	ElectionUID string `json:"electionUId"`
	Name        string `json:"name"`
}

var (
	electionData    = make([]Election, 0) // Slice of Election objects
	folderPath      = "."
	electionsFolder = "elections"
	lock            sync.RWMutex
	uidToFileName   = make(map[string]string)
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

// scanElectionsFolder scans the elections folder and updates the in-memory list.
func scanElectionsFolder() {
	lock.Lock()
	defer lock.Unlock()

	files, err := ioutil.ReadDir(filepath.Join(folderPath, electionsFolder))
	if err != nil {
		fmt.Printf("Error reading elections folder: %v\n", err)
		return
	}

	tempData := make([]Election, 0)

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePath := filepath.Join(folderPath, electionsFolder, file.Name())
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", file.Name(), err)
			continue
		}
		var election Election
		if err := json.Unmarshal(data, &election); err != nil {
			fmt.Printf("Error parsing JSON in file %s: %v\n", file.Name(), err)
			continue
		}
		tempData = append(tempData, election)
		uidToFileName[election.ElectionUID] = file.Name()
	}

	electionData = tempData
	fmt.Println("Election data updated")
}

// getElectionsHandler handles the GET /elections endpoint.
func getElectionsHandler(w http.ResponseWriter, r *http.Request) {
	lock.RLock()
	defer lock.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(electionData); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// getElectionFileHandler handles the GET /election/{id} endpoint.
func getElectionFileHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing electionUId", http.StatusBadRequest)
		return
	}

	lock.RLock()

	if uidToFileName[id] == "" {
		http.Error(w, "Election not found", http.StatusNotFound)
		return
	}

	lock.RUnlock()

	data, err := ioutil.ReadFile(filepath.Join(folderPath, electionsFolder, uidToFileName[id]))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read file: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func main() {
	// Set up periodic folder scanning
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			scanElectionsFolder()
		}
	}()

	// Initial scan
	scanElectionsFolder()

	// Set up HTTP routes
	http.HandleFunc("/elections", getElectionsHandler)
	http.HandleFunc("/election", getElectionFileHandler)

	// Wrap handlers with CORS middleware
	handler := corsMiddleware(http.DefaultServeMux)

	// Start the server
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
