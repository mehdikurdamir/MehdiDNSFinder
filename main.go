package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

//go:embed ui/*
var uiFiles embed.FS

var lastPing int64 // unix nano

func openUI(url string) {
	if runtime.GOOS != "windows" {
		return
	}
	candidates := []string{
		os.Getenv("ProgramFiles(x86)") + `\Microsoft\Edge\Application\msedge.exe`,
		os.Getenv("ProgramFiles") + `\Microsoft\Edge\Application\msedge.exe`,
		os.Getenv("ProgramFiles") + `\Google\Chrome\Application\chrome.exe`,
		os.Getenv("ProgramFiles(x86)") + `\Google\Chrome\Application\chrome.exe`,
		os.Getenv("LocalAppData") + `\Google\Chrome\Application\chrome.exe`,
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			if exec.Command(c, "--app="+url, "--window-size=1100,780").Start() == nil {
				return
			}
		}
	}
	_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
}

func main() {
	sub, _ := fs.Sub(uiFiles, "ui")
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(sub)))

	mux.HandleFunc("/api/servers", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(defaultServers)
	})

	mux.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		atomic.StoreInt64(&lastPing, time.Now().UnixNano())
		w.WriteHeader(204)
	})

	// POST /api/check  {target:"...", servers:[{name,ip,group}]}  -> NDJSON stream
	mux.HandleFunc("/api/check", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Target  string      `json:"target"`
			Servers []DNSServer `json:"servers"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", 400)
			return
		}
		t, err := parseTarget(req.Target)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		fl, _ := w.(http.Flusher)
		w.Header().Set("Content-Type", "application/x-ndjson; charset=utf-8")
		w.Header().Set("X-Host", t.Host)

		jobs := make(chan DNSServer)
		out := make(chan Result)
		var wg sync.WaitGroup
		for i := 0; i < 16; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for s := range jobs {
					out <- checkOne(s, t)
				}
			}()
		}
		go func() {
			for _, s := range req.Servers {
				select {
				case jobs <- s:
				case <-r.Context().Done():
					close(jobs)
					wg.Wait()
					close(out)
					return
				}
			}
			close(jobs)
			wg.Wait()
			close(out)
		}()
		enc := json.NewEncoder(w)
		for res := range out {
			enc.Encode(res)
			if fl != nil {
				fl.Flush()
			}
		}
	})

	ln, err := net.Listen("tcp", "127.0.0.1:47831")
	if err != nil {
		// برنامه قبلاً اجراست؛ فقط پنجره را باز کن
		openUI("http://127.0.0.1:47831")
		return
	}
	url := fmt.Sprintf("http://%s", ln.Addr().String())
	atomic.StoreInt64(&lastPing, time.Now().UnixNano())

	// وقتی پنجره بسته شد برنامه هم خاموش شود
	go func() {
		for {
			time.Sleep(3 * time.Second)
			if time.Since(time.Unix(0, atomic.LoadInt64(&lastPing))) > 150*time.Second {
				os.Exit(0)
			}
		}
	}()
	go openUI(url)
	http.Serve(ln, mux)
}
