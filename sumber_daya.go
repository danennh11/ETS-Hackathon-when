package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

func (s *Store) CreateSumberDaya(input SumberDaya) (SumberDaya, error) {
	// Trim spasi sesuai aturan kontrak
	input.Nama = strings.TrimSpace(input.Nama)
	input.Jenis = strings.TrimSpace(input.Jenis)

	// Validasi: nama wajib & maks 50 karakter, jenis ruang/alat, stok >= 1
	if input.Nama == "" || utf8.RuneCountInString(input.Nama) > 50 {
		return SumberDaya{}, errInputTidakValid
	}
	if input.Jenis != "ruang" && input.Jenis != "alat" {
		return SumberDaya{}, errInputTidakValid
	}
	if input.Stok < 1 {
		return SumberDaya{}, errInputTidakValid
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	input.ID = s.nextSDID
	s.nextSDID++
	s.sumberDayaList = append(s.sumberDayaList, input)

	return input, nil
}

func (s *Store) GetSumberDaya(id int) (SumberDaya, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, sd := range s.sumberDayaList {
		if sd.ID == id {
			return sd, nil
		}
	}
	return SumberDaya{}, errTidakDitemukan
}

func (s *Store) ListSumberDaya() []SumberDaya {
	s.mu.Lock()
	defer s.mu.Unlock()

	hasil := make([]SumberDaya, len(s.sumberDayaList))
	copy(hasil, s.sumberDayaList)
	return hasil
}

// POST /sumber-daya
func (a *API) createSumberDaya(w http.ResponseWriter, r *http.Request) {
	var input SumberDaya
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, errBodyTidakValid)
		return
	}

	sd, err := a.store.CreateSumberDaya(input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, sd)
}

// GET /sumber-daya
func (a *API) listSumberDaya(w http.ResponseWriter, r *http.Request) {
	list := a.store.ListSumberDaya()
	writeJSON(w, http.StatusOK, list)
}

// GET /sumber-daya/{id}
func (a *API) getSumberDaya(w http.ResponseWriter, r *http.Request) {
	idstr := r.PathValue("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		writeError(w, errIDBukanAngka)
		return
	}

	sd, err := a.store.GetSumberDaya(id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, sd)
}
