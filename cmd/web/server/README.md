# Penjelasan Timeout Server

## 1. `defaultIdleTimeout`

**Definisi**: Waktu maksimum server akan mempertahankan koneksi yang tidak aktif (idle) sebelum menutupnya.

**Contoh Kasus**:
- **Skenario**: Pengguna membuka halaman web (GET request) dan kemudian meninggalkan halaman tersebut tanpa melakukan aktivitas lebih lanjut. Server akan terus menjaga koneksi terbuka untuk pengguna tersebut, tetapi jika tidak ada aktivitas dalam waktu satu menit, server akan menutup koneksi.

**Tujuan**: Mencegah server dari terjebak dengan koneksi yang tidak aktif, yang membebani sumber daya server.

## 2. `defaultReadTimeout`

**Definisi**: Waktu maksimum server akan tunggu untuk membaca data dari klien setelah permintaan diterima.

**Contoh Kasus**:
- **Skenario**: Pengguna mengirimkan permintaan POST untuk mengunggah file besar ke server. Jika server tidak menerima data dari klien dalam waktu 5 detik setelah memulai pembacaan, server akan membatalkan permintaan tersebut.

**Tujuan**: Menghindari situasi di mana server terlalu lama menunggu data yang tidak datang, yang bisa menghambat pemrosesan permintaan lain.

## 3. `defaultWriteTimeout`

**Definisi**: Waktu maksimum server akan tunggu untuk menyelesaikan penulisan data ke klien.

**Contoh Kasus**:
- **Skenario**: Server sedang mengirimkan halaman web yang besar kepada pengguna. Jika server tidak dapat menyelesaikan pengiriman data dalam waktu 10 detik, server akan menghentikan pengiriman data tersebut.

**Tujuan**: Mencegah server dari terjebak dalam situasi di mana pengiriman data sangat lambat, yang bisa menyebabkan server tidak dapat menangani permintaan lain.

## 4. `defaultShutdownPeriod`

**Definisi**: Waktu maksimum server akan tunggu untuk menyelesaikan proses shutdown setelah menerima sinyal untuk berhenti.

**Contoh Kasus**:
- **Skenario**: Server menerima sinyal untuk berhenti (misalnya, ketika Anda ingin menutup aplikasi atau memulai pemeliharaan). Server akan memberikan waktu 30 detik untuk menyelesaikan semua permintaan yang sedang diproses dan menutup koneksi secara bersih sebelum benar-benar berhenti.

**Tujuan**: Memastikan semua transaksi atau permintaan yang sedang diproses diselesaikan dengan baik sebelum server berhenti, sehingga tidak ada data yang hilang atau permintaan yang terputus.
