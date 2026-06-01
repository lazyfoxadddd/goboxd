package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Request struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

type Response struct {
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func run(w http.ResponseWriter, r *http.Request) {
	var req Request

	// Read + validate request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid request", 400)
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}

	// Basic language validation
	if req.Language != "python" {
		json.NewEncoder(w).Encode(Response{
			Error: "unsupported language",
		})
		return
	}

	// Create isolated temp directory
	tmpDir, err := os.MkdirTemp("", "goboxd-*")
	if err != nil {
		http.Error(w, "failed to create temp dir", 500)
		return
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "main.py")

	// Write code
	if err := os.WriteFile(filePath, []byte(req.Code), 0644); err != nil {
		http.Error(w, "failed to write file", 500)
		return
	}

	// Add timeout (important)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Run inside nsjail
	cmd := exec.CommandContext(
		ctx,
		"/usr/local/bin/nsjail",
		"--mode", "o",
		"--time_limit", "2",
		"--max_cpus", "1",
		"--rlimit_as", "256",
		"--disable_proc",
		"--iface_no_lo",
		"--chroot", "/",
		"--cwd", "/tmp",
		"--",
		"/usr/bin/python3",
		filePath,
	)

	out, err := cmd.CombinedOutput()

	resp := Response{
		Output: string(out),
	}

	if err != nil {
		resp.Error = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/healthz", healthz)
	http.HandleFunc("/run", run)

	http.ListenAndServe(":8080", nil)
}