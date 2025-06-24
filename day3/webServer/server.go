package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type user struct {
	Name string `json:"name"`
}

type server struct {
	mux       *http.ServeMux
	userCache map[int]user
	mutex     sync.RWMutex
}

func newServer() *server {
	s := &server{
		mux:       http.NewServeMux(),
		userCache: make(map[int]user),
	}

	s.routes()

	return s
}

func (s *server) routes() {
	s.mux.HandleFunc("/", s.greetHandler)
	s.mux.HandleFunc("POST /users", s.postHandler)
	s.mux.HandleFunc("GET /users/{id}", s.getUserHandler)
	s.mux.HandleFunc("DELETE /users/{id}", s.deleteUserHandler)
	s.mux.HandleFunc("PUT /users/{id}", s.updateHandler)
}

func main() {
	s := newServer()

	server := &http.Server{
		Addr:         ":8080",
		Handler:      s.mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Println("Server starts at 8080")

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Server error:", err)
	}
}

func (s *server) updateHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	var err error
	id, err := strconv.Atoi(idStr)

	if err != nil || id < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, ok := s.userCache[id]; !ok {
		http.Error(w, "UserNotFound", http.StatusBadRequest)
	}

	var usr user
	err = json.NewDecoder(r.Body).Decode(&usr)

	if err != nil || usr.Name == "" {
		http.Error(w, "Invalid user payload", http.StatusBadRequest)
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.userCache[id] = usr

	w.WriteHeader(http.StatusAccepted)
}

func (s *server) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.userCache[id]; !ok {
		http.Error(w, "User Not Found", http.StatusNotFound)
		return
	}

	delete(s.userCache, id)
	w.WriteHeader(http.StatusNoContent)
}

func (s *server) getUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	var err error

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mutex.RLock()
	usr, ok := s.userCache[id]
	s.mutex.RUnlock()

	if !ok {
		http.Error(w, "User Not Found", http.StatusNotFound)
		return
	}

	j, err := json.Marshal(usr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_, err = w.Write(j)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (*server) greetHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	var err error

	response := map[string]string{
		"Message": "Hello World using JSON in GOLANG",
	}

	msg, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(msg)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *server) postHandler(w http.ResponseWriter, r *http.Request) {
	var usr user
	err := json.NewDecoder(r.Body).Decode(&usr)

	if err != nil || usr.Name == "" {
		http.Error(w, "Invalid user payload", http.StatusBadRequest)
		return
	}

	s.mutex.Lock()
	s.userCache[len(s.userCache)+1] = usr
	s.mutex.Unlock()

	w.WriteHeader(http.StatusNoContent)
}
