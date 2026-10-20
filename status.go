package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// Transisi status yang sah untuk persona C (Admin Himpunan):
//
//	diajukan -> disetujui
//	diajukan -> ditolak   (alasan wajib)
//	disetujui -> selesai
//
// Status lain (ditolak, selesai) adalah akhir dan tidak boleh berubah lagi.
var transisiStatus = map[string]map[string]bool{
	"diajukan":  {"disetujui": true, "ditolak": true},
	"disetujui": {"selesai": true},
}

func (s *Store) UbahStatusPeminjaman(id int, target, alasan string) (Peminjaman, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Cari peminjaman
	idx := -1
	for i, p := range s.peminjamanList {
		if p.ID == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		return Peminjaman{}, errTidakDitemukan
	}

	lama := s.peminjamanList[idx]

	if !transisiStatus[lama.Status][target] {
		return Peminjaman{}, errStatusKonflik
	}

	lama.Status = target
	if target == "ditolak" {
		lama.Alasan = strings.TrimSpace(alasan)
	} else {
		lama.Alasan = ""
	}

	s.peminjamanList[idx] = lama
	return lama, nil
}

func validasiUbahStatus(input *UbahStatusRequest) error {
	input.Status = strings.TrimSpace(input.Status)
	input.Alasan = strings.TrimSpace(input.Alasan)

	// Status harus salah satu dari tiga status tujuan
	if input.Status != "disetujui" && input.Status != "ditolak" && input.Status != "selesai" {
		return errInputTidakValid
	}

	// Menolak peminjaman wajib beralasan
	if input.Status == "ditolak" && input.Alasan == "" {
		return errInputTidakValid
	}

	return nil
}

type UbahStatusRequest struct {
	Status string `json:"status"`
	Alasan string `json:"alasan"`
}

// PATCH /peminjaman/{id}/status
func (a *API) ubahStatusPeminjaman(w http.ResponseWriter, r *http.Request) {
	// Urutan cek sesuai kontrak: id bukan angka -> 400,
	idstr := r.PathValue("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		writeError(w, errIDBukanAngka)
		return
	}

	// Cek keberadaan dulu, sebelum body dibaca
	if _, err := a.store.GetPeminjaman(id); err != nil {
		writeError(w, err)
		return
	}

	var input UbahStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, errBodyTidakValid)
		return
	}

	if err := validasiUbahStatus(&input); err != nil {
		writeError(w, err)
		return
	}

	p, err := a.store.UbahStatusPeminjaman(id, input.Status, input.Alasan)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, p)
}
