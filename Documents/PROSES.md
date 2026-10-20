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
SumberDaya menyimpan ruang atau alat. Field Jenis membedakan keduanya, dan Stok menandai kapasitas: ruang bernilai 1, alat bisa lebih dari 1.
Peminjaman menyimpan pengajuan. SumberDayaID menghubungkannya ke sumber daya. Tanggal, SlotMulai, dan SlotSelesai dipakai untuk cek bentrok. Status menyimpan alur diajukan → disetujui/ditolak → selesai. Alasan hanya terisi saat ditolak, jadi diberi omitempty.
Slice dipilih karena datanya sedikit dan urutannya otomatis sesuai urutan dibuat, sehingga list tidak perlu mengurutkan lagi. Pencarian per ID memakai loop biasa.
nextSDID dan nextPinjamID adalah counter ID yang naik terus, jadi ID tidak dipakai ulang.
sync.Mutex melindungi semua data, karena HTTP handler berjalan di goroutine terpisah.
Sentinel error (errTidakDitemukan, errBentrok, dll.) dikelompokkan per status code (400, 404, 409), sehingga handler tidak perlu memutuskan status code sendiri.


## 3. Satu keputusan desain
Keputusan: Pemetaan error ke status code dipusatkan di satu fungsi writeError dengan errors.Is pada sentinel error. Store hanya mengembalikan error, dan handler cukup memanggil writeError(w, err).

## 4. Bug atau kendala (minimal dua)
### Bug 1
- Gejala: Request POST /sumber-daya dengan "stok":"banyak" mengembalikan “body JSON tidak valid”, sedangkan "stok":0 mengembalikan “input tidak valid atau tidak lengkap”. Dua kasus input salah menghasilkan pesan error yang berbeda.
- Cara menemukan:es dengan curl memakai berbagai input salah, lalu membandingkan pesan error-nya.
- Penyebab dan solusi: Tipe data yang salah gagal di tahap Decode (jadi errBodyTidakValid), sedangkan nilai yang salah gagal di tahap validasi (jadi errInputTidakValid). Ini dua tahap yang berbeda. Solusinya: sepakati dengan kontrak bahwa keduanya tetap 400, dan biarkan pesannya berbeda karena penyebabnya memang berbeda.

### Bug 2
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

## 5. Change request
Bagian ini belum bisa kuisi karena aku belum tahu change request yang diberikan. Kirim isinya (misalnya “tambah field X” atau “ubah aturan bentrok jadi Y”), nanti kubantu isi bagian ini lengkap dengan kode yang berubah dan keputusan desainnya.
