package main

import (
	"encoding/json"
	"net/http"
)

func startAPIServer(store *MemoryStore) {
	http.HandleFunc("/node", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			var node Node
			if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			store.AddNode(&node)
			w.WriteHeader(http.StatusCreated)
		case "GET":
			mac := r.URL.Query().Get("mac")
			if node, exists := store.GetNode(mac); exists {
				json.NewEncoder(w).Encode(node)
			} else {
				http.NotFound(w, r)
			}
		case "DELETE":
			mac := r.URL.Query().Get("mac")
			store.DeleteNode(mac)
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	http.ListenAndServe(":8080", nil)
}
