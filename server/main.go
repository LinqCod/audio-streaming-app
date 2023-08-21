package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

func createHLS(inputFile string, outputDir string, segmentDuration int) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	ffmpegCmd := exec.Command(
		"ffmpeg",
		"-i", inputFile,
		"-map", "0:a",
		"-c:a", "aac",
		"-b:a", "320k",
		"-hls_time", strconv.Itoa(segmentDuration),
		"-hls_list_size", "0",
		"-f", "hls",
		fmt.Sprintf("%s/playlist.m3u8", outputDir),
	)

	output, err := ffmpegCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, string(output))
	}

	return nil
}

func printDirFiles(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}

	return nil
}

func main() {
	port := "8080"
	addr := ":" + port

	inputFile := "./test.mp3"
	outputDir := "./output"
	segmentDuration := 10

	//Creating HLS
	if err := createHLS(inputFile, outputDir, segmentDuration); err != nil {
		log.Fatalf("error while creating HLS: %v", err)
	}

	log.Println("HLS created successfully!")

	if err := printDirFiles(outputDir); err != nil {
		log.Fatalf("error while listing result files: %v", err)
	}

	// Create a new file server instance
	fs := http.FileServer(http.Dir("output"))

	// Wrap the file server with CORS headers middleware
	http.Handle("/", corsMiddleware(fs))

	// Start the HTTP file server on the specified address and port
	log.Printf("Starting HLS file server on %s...", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

		// Pre-flight request. Reply successfully:
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass the request to the next handler (the file server)
		next.ServeHTTP(w, r)
	})
}
