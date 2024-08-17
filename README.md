# PRBCareAPI

<p align="center">
<img src="https://github.com/user-attachments/assets/50eea6b6-e922-4dda-a036-3fbf1704458d" alt="prbcare" width="400">
</p>

PRBCareAPI adalah aplikasi REST API untuk manajemen Puskesmas, manajemen Apotek, pengambilan obat, kontrol balik, dan
manajemen pasien. Aplikasi ini menyediakan fungsionalitas khusus berdasarkan peran pengguna yang berbeda, termasuk Admin
Super, Admin Puskesmas, Admin Apotek, dan Calon Pasien.
PRBCareAPI dikembangkan dengan mengikuti prinsip-prinsip REST API untuk memastikan skalabilitas dan pemeliharaan yang
mudah. Sistem autentikasi dilengkapi untuk memastikan keamanan data.

## Fitur

- Autentikasi yang berbeda untuk Admin Super, Admin Puskesmas, Admin Apotek, dan Calon Pasien.
- Manajemen pasien oleh Admin Puskesmas yang meliputi pendaftaran, pembaruan data, dan pencatatan medis.
- Manajemen obat oleh Admin Apotek, termasuk stok dan dispensasi obat.
- Kontrol balik oleh Admin Puskesmas untuk memonitor dan mengevaluasi pengobatan pasien.
- Sistem pembuatan jadwal kontrol balik dan pengambilan obat oleh Admin Puskesmas.

## Tech Stack

- **Programming Language**: Golang
- **Web Framework**: Fiber
- **ORM**: GORM
- **Database**: PostgreSQL

## Environment Variables

PRBCareAPI akan menggunakan environment variables sebagai konfigurasi utama menggantikan `config.json` jika
variabel-variabel tersebut diset sebelum menjalankan proyek:

| **Key**                  | **Type**     | **Deskripsi**                                                                        | **Contoh**                                      |
|----------------------|----------|----------------------------------------------------------------------------------|---------------------------------------------|
| **JWT_SECRET**       | `string` | Secret key untuk JWT.                                                            | `mysecretkey123`                            |
| **JWT_EXP**          | `int`    | Waktu kadaluwarsa JWT dalam jam.                                                 | `24`                                        |
| **WEB_PORT**         | `int`    | Port untuk menjalankan server web.                                               | `8080`                                      |
| **WEB_CORS_ORIGINS** | `string` | Origins yang diizinkan untuk CORS, dipisahkan dengan spasi jika lebih dari satu. | `http://localhost http://example.com`       |
| **RECAPTCHA_SECRET** | `string` | Secret key untuk Google reCAPTCHA.                                               | `6LfQzM0UAAAAABfG8qB1T4KQwQ7Js2w9p3sI9sA2` |
| **DB_USERNAME**      | `string` | Nama pengguna database.                                                          | `root`                                      |
| **DB_PASSWORD**      | `string` | Kata sandi database.                                                             | `password123`                               |
| **DB_HOST**          | `string` | Host database.                                                                   | `localhost`                                 |
| **DB_PORT**          | `int`    | Port koneksi database.                                                           | `3306`                                      |
| **DB_NAME**          | `string` | Nama database.                                                                   | `prbcare`                                   |


Cara set environment variables:

- **Windows**: Gunakan System Properties > Advanced > Environment Variables, atau command setx.
- **Linux/macOS**: Tambahkan export VARIABLE="value" ke file .bashrc atau .profile dan jalankan source ~/.bashrc.

## Dokumentasi API

Untuk mendapatkan lebih detail mengenai endpoint dan cara penggunaan API, kunjungi dokumentasi API di link berikut:

[API Documentation](https://app.swaggerhub.com/apis-docs/restfullapi/PRB-Care-API/1.0.0)

## Implementasi Frontend

Lihat implementasi frontend untuk aplikasi PRBCareAPI di link berikut:

[Frontend Implementation](https://github.com/RyanAprs/PRB-Care-Client.git)

## Aplikasi Scheduler Pendukung

Aplikasi scheduler mendukung pengingat melalui push notifikasi dan pembatalan jadwal secara otomatis. Informasi lebih
lanjut dan dokumentasi aplikasi scheduler dapat diakses melalui link berikut:

[Scheduler Application](https://github.com/scrkiddie/PRBCareScheduler)

