package main

import "errors"

var (
	// Error 400 (Bad Request)
	errIDBukanAngka    = errors.New("id harus berupa angka")
	errBodyTidakValid  = errors.New("body JSON tidak valid")
	errInputTidakValid = errors.New("input tidak valid atau tidak lengkap")

	// Error 404 (Not Found)
	errTidakDitemukan = errors.New("data tidak ditemukan")

	// Error 409 (Conflict)
	errBentrok       = errors.New("jadwal peminjaman bentrok atau melebihi stok")
	errStatusKonflik = errors.New("perubahan status tidak sesuai alur")
)

type SumberDaya struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Jenis string `json:"jenis"`
	Stok  int    `json:"stok"`
}

type Peminjaman struct {
	ID           int    `json:"id"`
	SumberDayaID int    `json:"sumber_daya_id"`
	Peminjam     string `json:"peminjam"`
	Tanggal      string `json:"tanggal"`
	SlotMulai    int    `json:"slot_mulai"`
	SlotSelesai  int    `json:"slot_selesai"`
	Jumlah       int    `json:"jumlah"`
	Status       string `json:"status"`
	Alasan       string `json:"alasan,omitempty"`
}
