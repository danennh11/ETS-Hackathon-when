package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func (s *Store) validasiPeminjaman(input *Peminjaman) (SumberDaya, error) {
	input.Peminjam = strings.TrimSpace(input.Peminjam)
	input.Tanggal = strings.TrimSpace(input.Tanggal)

	if input.Peminjam == "" || utf8.RuneCountInString(input.Peminjam) > 50 {
		return SumberDaya{}, errInputTidakValid
	}

	// Cek tanggal format YYYY-MM-DD dan valid di kalender
	t, err := time.Parse("2006-01-02", input.Tanggal)
	if err != nil || t.Format("2006-01-02") != input.Tanggal {
		return SumberDaya{}, errInputTidakValid
	}

	// Cek slot 1-8
	if input.SlotMulai < 1 || input.SlotMulai > 8 || input.SlotSelesai < 1 || input.SlotSelesai > 8 || input.SlotMulai > input.SlotSelesai {
		return SumberDaya{}, errInputTidakValid
	}

	// Cari sumber daya
	var sumberDaya SumberDaya
	ketemu := false
	for _, sd := range s.sumberDayaList {
		if sd.ID == input.SumberDayaID {
			sumberDaya = sd
			ketemu = true
			break
		}
	}
	if !ketemu {
		return SumberDaya{}, errInputTidakValid // Jika sumber_daya_id tidak ada -> 400
	}

	// Cek jumlah antara 1 sampai stok
	if input.Jumlah < 1 || input.Jumlah > sumberDaya.Stok {
		return SumberDaya{}, errInputTidakValid
	}

	return sumberDaya, nil
}
func (s *Store) cekBentrok(calon Peminjaman, stok int, abaikanID int) bool {
	for slot := calon.SlotMulai; slot <= calon.SlotSelesai; slot++ {
		terpakai := 0
		for _, p := range s.peminjamanList {
			if p.ID == abaikanID {
				continue // Jangan bentrok dengan diri sendiri saat PUT
			}
			if p.SumberDayaID != calon.SumberDayaID || p.Tanggal != calon.Tanggal {
				continue
			}
			// Khusus Persona C: hanya status "diajukan" dan "disetujui" yang menahan slot
			if p.Status == "diajukan" || p.Status == "disetujui" {
				if slot >= p.SlotMulai && slot <= p.SlotSelesai {
					terpakai += p.Jumlah
				}
			}
		}
		if terpakai+calon.Jumlah > stok {
			return true
		}
	}
	return false
}

// List peminjaman
func (s *Store) ListPeminjaman(filterSDID int, tanggal string, status string) []Peminjaman {
	s.mu.Lock()
	defer s.mu.Unlock()

	hasil := []Peminjaman{} // Memastikan kalau kosong hasilnya [] bukan null
	for _, p := range s.peminjamanList {
		if filterSDID != 0 && p.SumberDayaID != filterSDID {
			continue
		}
		if tanggal != "" && p.Tanggal != tanggal {
			continue
		}
		if status != "" && p.Status != status {
			continue
		}
		hasil = append(hasil, p)
	}
	return hasil
}

// Get peminjaman
func (s *Store) GetPeminjaman(id int) (Peminjaman, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, p := range s.peminjamanList {
		if p.ID == id {
			return p, nil
		}
	}
	return Peminjaman{}, errTidakDitemukan
}

func (s *Store) CreatePeminjaman(input Peminjaman) (Peminjaman, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sd, err := s.validasiPeminjaman(&input)
	if err != nil {
		return Peminjaman{}, err
	}

	if s.cekBentrok(input, sd.Stok, 0) {
		return Peminjaman{}, errBentrok
	}

	input.ID = s.nextPinjamID
	s.nextPinjamID++
	input.Status = "diajukan" // Status awal Persona C
	input.Alasan = ""

	s.peminjamanList = append(s.peminjamanList, input)
	return input, nil
}

func (s *Store) UpdatePeminjaman(id int, input Peminjaman) (Peminjaman, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Cari peminjaman lama
	idx := -1
	lama := Peminjaman{}
	for i, p := range s.peminjamanList {
		if p.ID == id {
			idx = i
			lama = p
			break
		}
	}
	if idx == -1 {
		return Peminjaman{}, errTidakDitemukan
	}

	// Validasi enam field baru (sama persis dengan saat POST)
	sd, err := s.validasiPeminjaman(&input)
	if err != nil {
		return Peminjaman{}, err
	}

	// Peminjaman lama tidak boleh dianggap bentrok dengan dirinya sendiri
	if s.cekBentrok(input, sd.Stok, id) {
		return Peminjaman{}, errBentrok
	}

	// Ganti data, pertahankan id dan status
	lama.SumberDayaID = input.SumberDayaID
	lama.Peminjam = input.Peminjam
	lama.Tanggal = input.Tanggal
	lama.SlotMulai = input.SlotMulai
	lama.SlotSelesai = input.SlotSelesai
	lama.Jumlah = input.Jumlah
	s.peminjamanList[idx] = lama

	return lama, nil
}

// DeletePeminjaman menghapus peminjaman berdasarkan id.
func (s *Store) DeletePeminjaman(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, p := range s.peminjamanList {
		if p.ID == id {
			// Hapus dengan menggabungkan potongan sebelum dan sesudah idx
			s.peminjamanList = append(s.peminjamanList[:i], s.peminjamanList[i+1:]...)
			return nil
		}
	}
	return errTidakDitemukan
}

// PUT /peminjaman/{id}
func (a *API) updatePeminjaman(w http.ResponseWriter, r *http.Request) {

	idstr := r.PathValue("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		writeError(w, errIDBukanAngka)
		return
	}

	// Cek keberadaan dulu sebelum body dibaca
	if _, err := a.store.GetPeminjaman(id); err != nil {
		writeError(w, err)
		return
	}

	var input Peminjaman
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, errBodyTidakValid)
		return
	}

	p, err := a.store.UpdatePeminjaman(id, input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, p)
}

// DELETE /peminjaman/{id}
func (a *API) deletePeminjaman(w http.ResponseWriter, r *http.Request) {
	idstr := r.PathValue("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		writeError(w, errIDBukanAngka)
		return
	}

	if err := a.store.DeletePeminjaman(id); err != nil {
		writeError(w, err)
		return
	}

	// 204 = sukses tanpa body
	w.WriteHeader(http.StatusNoContent)
}

// POST /peminjaman
func (a *API) createPeminjaman(w http.ResponseWriter, r *http.Request) {
	var input Peminjaman
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, errBodyTidakValid)
		return
	}

	p, err := a.store.CreatePeminjaman(input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, p)
}

// GET /peminjaman
func (a *API) listPeminjaman(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	sdIDStr := strings.TrimSpace(q.Get("sumber_daya_id"))
	tanggal := strings.TrimSpace(q.Get("tanggal"))
	status := strings.TrimSpace(q.Get("status"))

	filterSDID := 0
	if sdIDStr != "" {
		id, err := strconv.Atoi(sdIDStr)
		if err != nil || id < 1 {
			writeError(w, errIDBukanAngka)
			return
		}
		filterSDID = id
	}

	list := a.store.ListPeminjaman(filterSDID, tanggal, status)
	writeJSON(w, http.StatusOK, list)
}

// GET /peminjaman/{id}
func (a *API) getPeminjaman(w http.ResponseWriter, r *http.Request) {
	idstr := r.PathValue("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		writeError(w, errIDBukanAngka)
		return
	}

	p, err := a.store.GetPeminjaman(id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, p)
}
