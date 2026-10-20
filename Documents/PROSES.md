# PROSES

## 1. Persona dan masalahnya
Pengguna persona saya adalah Admin Himpunan. Situasi/masalah yang sedang dihadapin adalah setiap peminjaman layanan ruangan atau alat kampus, harus di setujui terlebih dahulu, tidak boleh asal dipakai / dipinjam. Fitur khusus yang ada dari persona C adalah, alur status: 1. diajukan 2. disetujui/ditolak(dengan alasan) 3. selesai.

## 2. Rancangan struktur data
(Tulis struct dan slice/map apa yang kamu simpan, dan kenapa dipilih begitu. Rancang ini sebelum menulis kode, lalu catat kalau berubah di tengah jalan.)

```go
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
```

## 3. Satu keputusan desain
(Satu keputusan yang kamu ambil, pilihan lain yang kamu pertimbangkan, dan alasanmu memilih yang ini. Contoh: kenapa cek bentrok per slot, bukan per rentang.)

## 4. Bug atau kendala (minimal dua)
### Bug 1
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

### Bug 2
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

## 5. Change request
(Apa yang diminta, bagian kode mana yang kamu ubah, dan satu keputusan desain yang kamu ambil.)
