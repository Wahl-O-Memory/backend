package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

type SVGResponse struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func handleSVGsRequest(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Content-Type", "application/json")

	// Parse query parameters
	query := r.URL.Query()
	requestedFiles := query["file"]

	var svgs []SVGResponse

	if len(requestedFiles) == 0 {
		// Return all SVGs if no files specified
		files, err := ioutil.ReadDir(svgDirectory)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading SVG directory: %v", err), http.StatusInternalServerError)
			return
		}

		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".svg") {
				svg, err := readSVGFile(file.Name())
				if err != nil {
					svg = generateErrorSVG(file.Name(), err.Error())
				}
				svgs = append(svgs, svg)
			}
		}
	} else {
		// Return only requested files
		for _, filename := range requestedFiles {
			if !strings.HasSuffix(filename, ".svg") {
				filename += ".svg"
			}
			svg, err := readSVGFile(filename)
			if err != nil {
				svg = generateErrorSVG(filename, err.Error())
			}
			svgs = append(svgs, svg)
		}
	}

	// Return JSON response
	json.NewEncoder(w).Encode(svgs)
}

func readSVGFile(filename string) (SVGResponse, error) {
	// Clean the path to remove any ../ or ./ sequences
	cleanedName := filepath.Clean(filename)

	// Check for path traversal attempts
	if strings.Contains(cleanedName, "..") || strings.Contains(cleanedName, "/") {
		return SVGResponse{}, fmt.Errorf("invalid file path")
	}

	filePath := filepath.Join(svgDirectory, filename)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return SVGResponse{}, err
	}
	return SVGResponse{
		Name: filename,
		Data: string(data),
	}, nil
}

func generateErrorSVG(filename, errorMsg string) SVGResponse {
	errorSVG := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100">
		<rect width="100" height="100" fill="red" opacity="0.3"/>
		<text x="50" y="50" font-family="Arial" font-size="10" text-anchor="middle" fill="black">
			<tspan x="50" dy="-1em">%s</tspan>
			<tspan x="50" dy="1.2em">Not Found</tspan>
			<tspan x="50" dy="1.2em">%s</tspan>
		</text>
	</svg>`, filename, errorMsg)

	return SVGResponse{
		Name: filename,
		Data: errorSVG,
	}
}
