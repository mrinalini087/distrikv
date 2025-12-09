package main

import (
	"bytes"
	"distrikv/pkg"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/http"
)

// The list of nodes we will run
var nodes = []string{
	"http://localhost:8081",
	"http://localhost:8082",
	"http://localhost:8083",
}

// Simple hash function to pick a node based on the key
func getNodeForKey(key string) string {
	h := fnv.New32a()
	h.Write([]byte(key))
	index := int(h.Sum32()) % len(nodes)
	return nodes[index]
}

func handleSet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req pkg.SetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	targetNode := getNodeForKey(req.Key)
	fmt.Printf("Routing key '%s' to %s\n", req.Key, targetNode)


	jsonData, _ := json.Marshal(req)
	resp, err := http.Post(targetNode+"/set", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Failed to connect to node", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing key", http.StatusBadRequest)
		return
	}

	targetNode := getNodeForKey(key)
	

	resp, err := http.Get(targetNode + "/get?key=" + key)
	if err != nil {
		http.Error(w, "Failed to connect to node", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}
	
	io.Copy(w, resp.Body)
}

func main() {
	http.HandleFunc("/set", handleSet)
	http.HandleFunc("/get", handleGet)

	fmt.Println("Load Balancer running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}