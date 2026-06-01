package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
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

	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &req)

	// create temp dir
	os.MkdirAll("tmp", 0755)

	// write file
	file := "tmp/main.py"
	os.WriteFile(file, []byte(req.Code), 0644)

	// execute (no sandbox yet)
	cmd := exec.Command(
	"nsjail",
	"--time_limit", "2",
	"--max_cpus", "1",
	"--rlimit_as", "256",
	"--disable_proc",
	"--iface_no_lo",
	"--",
	"python3", file,
	)
	out, err := cmd.CombinedOutput()

	resp := Response{
		Output: string(out),
	}

	if err != nil {
		resp.Error = err.Error()
	}

	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/healthz", healthz)
	http.HandleFunc("/run", run)

	http.ListenAndServe(":8080", nil)
}