package hook

import (
	"io"
	"path"
	"mime"
	"strings"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"path/filepath"

	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/ncarlier/webhookd/pkg/hook"
	"github.com/ncarlier/webhookd/pkg/logger"
	"github.com/ncarlier/webhookd/pkg/worker"
	// "gitee.com/infrastlabs/hostcross/api/config"
)


var (
	defaultTimeout int
	scriptDir      string
	scriptPath     string
	outputDir      string
)

func SetVars(prefix string){
	defaultTimeout= 10
	scriptDir= "scripts" //"/_ext/working/_ee/fk-webhookd/hostcross/scripts"
	scriptPath= prefix //"/aa/hook"
	outputDir= "logs"
}

func MuxHandle(ru *mux.Router, prefix string) {
	logger.Init("info", "out")
	// logger.Init("info")
	worker.StartDispatcher(2)
	ru.PathPrefix(prefix).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// getWebhookLog(w, r)
		if r.Method == "GET" {
			// http://172.17.0.21:18089/aa/hook/echo.sh/3
			if _, err := strconv.Atoi(filepath.Base(r.URL.Path)); err == nil {
				getWebhookLog(w, r)
				return
			}
		}
		triggerWebhook(w, r)
	}))
}

func triggerWebhook(w http.ResponseWriter, r *http.Request) {
	// Check that streaming is supported
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	// Get hook location
	hookName := strings.TrimPrefix(r.URL.Path, scriptPath+"/")
	if hookName == "" {
		infoHandler(w, r)
		return
	}
	_, err := hook.ResolveScript(scriptDir, hookName)
	if err != nil {
		logger.Error.Println(err.Error())
		http.Error(w, "hook not found", http.StatusNotFound)
		return
	}

	if err = r.ParseForm(); err != nil {
		logger.Error.Printf("error reading from-data: %v", err)
		http.Error(w, "unable to parse request form", http.StatusBadRequest)
		return
	}

	// parse body
	var body []byte
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediatype, _, _ := mime.ParseMediaType(ct)
		if strings.HasPrefix(mediatype, "text/") || mediatype == "application/json" {
			body, err = ioutil.ReadAll(r.Body)
			if err != nil {
				logger.Error.Printf("error reading body: %v", err)
				http.Error(w, "unable to read request body", http.StatusBadRequest)
				return
			}
		}
	}
	
	params := URLValuesToShellVars(r.Form)
	params = append(params, HTTPHeadersToShellVars(r.Header)...)

	// logger.Debug.Printf("API REQUEST: \"%s\" with params %s...\n", p, params)

	// Create work
	timeout := atoiFallback(r.Header.Get("X-Hook-Timeout"), defaultTimeout)
	job, err := hook.NewHookJob(&hook.Request{
		Name:      hookName,
		Method:    r.Method,
		Payload:   string(body),
		Args:      params,
		Timeout:   timeout,
		BaseDir:   scriptDir,
		OutputDir: outputDir,
	})
	if err != nil {
		logger.Error.Printf("error creating hook job: %v", err)
		http.Error(w, "unable to create hook job", http.StatusInternalServerError)
		return
	}
	/* http.Error(w, "hook not found111", http.StatusNotFound)
	return */

	// Put work in queue
	worker.WorkQueue <- job

	// Use content negotiation to enable Server-Sent Events
	useSSE := r.Method == "GET" && r.Header.Get("Accept") == "text/event-stream"
	if useSSE {
		// Send SSE response
		w.Header().Set("Content-Type", "text/event-stream")
	} else {
		// Send chunked response
		w.Header().Set("X-Content-Type-Options", "nosniff")
	}
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Hook-ID", strconv.FormatUint(job.ID(), 10))

	for {
		msg, open := <-job.MessageChan
		if !open {
			break
		}
		if useSSE {
			fmt.Fprintf(w, "data: %s\n\n", msg) // Send SSE response
		} else {
			fmt.Fprintf(w, "%s\n", msg) // Send chunked response
		}
		// Flush the data immediately instead of buffering it for later.
		flusher.Flush()
	}
}
func atoiFallback(str string, fallback int) int {
	if value, err := strconv.Atoi(str); err == nil && value > 0 {
		return value
	}
	return fallback
}

// http://172.17.0.21:18089/aa/hook/echo.sh/3
func getWebhookLog(w http.ResponseWriter, r *http.Request) {
	// Get hook ID
	id := path.Base(r.URL.Path)
	fmt.Println(r.URL.Path)
	fmt.Println(id)

	// Get script location
	hookName := path.Dir(strings.TrimPrefix(r.URL.Path, scriptPath+"/"))
	// hookName= "echo.sh"
	fmt.Println("hookName: "+hookName)
	_, err := hook.ResolveScript(scriptDir, hookName)
	if err != nil {
		logger.Error.Println(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Retrieve log file
	logFile, err := hook.Logs(id, hookName, outputDir)
	if err != nil {
		logger.Error.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if logFile == nil {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}
	defer logFile.Close()

	w.Header().Set("Content-Type", "text/plain")

	io.Copy(w, logFile)
}

// Info API informations model structure.
type Info struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
func infoHandler(w http.ResponseWriter, r *http.Request) {
	info := Info{
		Name:    "webhookd",
		Version: "config.Version",
	}
	data, err := json.Marshal(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
