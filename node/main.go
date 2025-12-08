package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"distrikv/pkg" // This imports the file above
)

// Store holds the data and a lock to prevent errors
type Store struct {
	data map[string]string
	mu   sync.RWMutex
}

var store = Store{
	data: make(map[string]string),
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req pkg.SetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Lock, Save, Unlock
	store.mu.Lock()
	store.data[req.Key] = req.Value
	store.mu.Unlock()

	log.Printf("Saved key: %s", req.Key)
	w.WriteHeader(http.StatusOK)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	
	store.mu.RLock()
	val, ok := store.data[key]
	store.mu.RUnlock()

	if !ok {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	w.Write([]byte(val))
}

func main() {
	// Takes a port number from the command line (e.g. -port=8081)
	portPtr := flag.String("port", "8080", "port number")
	flag.Parse()

	port := ":" + *portPtr

	http.HandleFunc("/set", setHandler)
	http.HandleFunc("/get", getHandler)

	fmt.Printf("Storage Node running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}