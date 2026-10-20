package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
)

type Store struct {
	mu             sync.Mutex
	sumberDayaList []SumberDaya
	peminjamanList []Peminjaman
	nextSDID       int
	nextPinjamID   int
}

func NewStore() *Store {
	return &Store{
		sumberDayaList: []SumberDaya{},
		peminjamanList: []Peminjaman{},
		nextSDID:       1,
		nextPinjamID:   1,
	}
}

type API struct {
	store *Store
}

func NewHandler() http.Handler {
	a := &API{
		store: NewStore(),
	}

	mux := http.NewServeMux()

	// B1,B2,B3
	mux.HandleFunc("POST /sumber-daya", a.createSumberDaya)
	mux.HandleFunc("GET /sumber-daya", a.listSumberDaya)
	mux.HandleFunc("GET /sumber-daya/{id}", a.getSumberDaya)
	mux.HandleFunc("POST /peminjaman", a.createPeminjaman)
	mux.HandleFunc("GET /peminjaman", a.listPeminjaman)
	mux.HandleFunc("GET /peminjaman/{id}", a.getPeminjaman)

	return mux
}

func main() {
	http.ListenAndServe(":8080", NewHandler())
}

// Helper kirim JSON
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusBadRequest

	if errors.Is(err, errTidakDitemukan) {
		status = http.StatusNotFound
	} else if errors.Is(err, errBentrok) || errors.Is(err, errStatusKonflik) {
		status = http.StatusConflict
	}

	writeJSON(w, status, map[string]string{"error": err.Error()})
}
