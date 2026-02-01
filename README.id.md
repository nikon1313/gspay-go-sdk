# GSPAY Go SDK (Tidak Resmi)

[![Go Reference](https://pkg.go.dev/badge/github.com/H0llyW00dzZ/gspay-go-sdk.svg)](https://pkg.go.dev/github.com/H0llyW00dzZ/gspay-go-sdk)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/H0llyW00dzZ/gspay-go-sdk)](https://goreportcard.com/report/github.com/H0llyW00dzZ/gspay-go-sdk)
[![codecov](https://codecov.io/gh/H0llyW00dzZ/gspay-go-sdk/graph/badge.svg?token=AITK1X3RSE)](https://codecov.io/gh/H0llyW00dzZ/gspay-go-sdk)
[![View on DeepWiki](https://img.shields.io/badge/View%20on-DeepWiki-blue)](https://deepwiki.com/H0llyW00dzZ/gspay-go-sdk)

SDK Go **tidak resmi** untuk API Payment Gateway GSPAY2. SDK ini menyediakan antarmuka Go yang komprehensif dan idiomatik untuk pemrosesan pembayaran, pencairan dana, dan pengecekan saldo.

> **Penyangkalan**: Ini adalah SDK tidak resmi dan tidak berafiliasi dengan, didukung oleh, atau secara resmi didukung oleh GSPAY. SDK ini dikembangkan secara independen untuk menyediakan kompatibilitas bahasa Go dalam mengintegrasikan API Payment Gateway GSPAY2. Gunakan dengan kebijaksanaan Anda sendiri. [Baca selengkapnya](#penyangkalan)

## Fitur

- **Pembayaran IDR**: Membuat pembayaran melalui QRIS, DANA, dan virtual account bank
- **Pencairan IDR**: Memproses penarikan ke rekening bank dan e-wallet Indonesia
- **Pembayaran USDT**: Menerima pembayaran cryptocurrency melalui jaringan TRC20
- **Pengecekan Saldo**: Memeriksa saldo settlement operator
- **Verifikasi Callback**: Verifikasi tanda tangan yang aman untuk webhook
- **Logika Retry**: Pengulangan otomatis dengan exponential backoff untuk kegagalan sementara
- **Dukungan Context**: Dukungan penuh context.Context untuk pembatalan dan timeout

## Instalasi

```bash
go get github.com/H0llyW00dzZ/gspay-go-sdk
```

## Struktur Proyek

```
gspay-go-sdk/
├── src/
│   ├── client/      # HTTP client dan konfigurasi
│   ├── constants/   # Kode bank, channel, kode status
│   ├── errors/      # Tipe error dan helper
│   ├── payment/     # Layanan pembayaran (IDR, USDT)
│   ├── payout/      # Layanan pencairan (IDR)
│   ├── balance/     # Layanan pengecekan saldo
│   └── internal/    # Utilitas internal (tanda tangan)
└── examples/        # Contoh penggunaan
```

## Mulai Cepat

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/H0llyW00dzZ/gspay-go-sdk/src/client"
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/constants"
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/payment"
)

func main() {
    // Membuat client baru
    c := client.New("your-auth-key", "your-secret-key")

    // Membuat layanan pembayaran
    paymentSvc := payment.NewIDRService(c)

    ctx := context.Background()

    // Membuat pembayaran IDR
    resp, err := paymentSvc.Create(ctx, &payment.IDRRequest{
        TransactionID:  client.GenerateTransactionID("TXN"),
        Username:       "user123",
        Amount:         50000, // 50.000 IDR
        Channel:        constants.ChannelQRIS,
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("URL Pembayaran: %s\n", resp.PaymentURL)
    fmt.Printf("ID Pembayaran: %s\n", resp.IDRPaymentID)
    fmt.Printf("Kedaluwarsa: %s\n", resp.ExpireDate)
}
```

## Opsi Konfigurasi

Client mendukung berbagai opsi konfigurasi menggunakan pola functional options:

```go
c := client.New(
    "auth-key",
    "secret-key",
    client.WithBaseURL("https://custom-api.example.com"),
    client.WithTimeout(60 * time.Second),
    client.WithRetries(5),
    client.WithRetryWait(500*time.Millisecond, 5*time.Second),
    client.WithHTTPClient(customHTTPClient),
)
```

| Opsi | Deskripsi | Default |
|------|-----------|---------|
| `WithBaseURL` | Mengatur URL dasar API kustom | `https://api.thegspay.com` |
| `WithTimeout` | Mengatur timeout request | `30s` |
| `WithRetries` | Mengatur jumlah percobaan ulang | `3` |
| `WithRetryWait` | Mengatur waktu tunggu min/maks antar retry | `500ms` / `2s` |
| `WithHTTPClient` | Menggunakan HTTP client kustom | Default `http.Client` |

## Contoh Penggunaan

### Membuat Pembayaran IDR

```go
import (
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/client"
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/constants"
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/payment"
)

c := client.New("auth-key", "secret-key")
paymentSvc := payment.NewIDRService(c)

    resp, err := paymentSvc.Create(ctx, &payment.IDRRequest{
        TransactionID:  client.GenerateTransactionID("TXN"),
        Username:       "user123",
        Amount:         50000,
        Channel:        constants.ChannelQRIS, // Opsional: QRIS, DANA, atau BNI
    })
if err != nil {
    log.Fatal(err)
}

// Arahkan pengguna ke halaman pembayaran
fmt.Printf("Arahkan ke: %s\n", resp.PaymentURL)

// Opsional menambahkan return URL
redirectURL := client.BuildReturnURL(resp.PaymentURL, "https://mysite.com/complete")
```

### Cek Status Pembayaran

```go
status, err := paymentSvc.GetStatus(ctx, "TXN20260126143022123")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Status: %s\n", status.Status.String())
if status.Status.IsSuccess() {
    fmt.Println("Pembayaran berhasil!")
}
```

### Membuat Pencairan IDR

```go
import (
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/client"
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/payout"
)

c := client.New("auth-key", "secret-key")
payoutSvc := payout.NewIDRService(c)

    resp, err := payoutSvc.Create(ctx, &payout.IDRRequest{
        TransactionID:  client.GenerateTransactionID("PAY"),
        Username:       "user123",
        AccountName:    "John Doe",
        AccountNumber:  "1234567890",
        Amount:         50000,
        BankCode:       "BCA",
        Description:    "Permintaan penarikan",
    })
if err != nil {
    log.Fatal(err)
}

fmt.Printf("ID Pencairan: %s\n", resp.IDRPayoutID)
```

### Membuat Pembayaran USDT

```go
import (
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/client"
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/payment"
)

c := client.New("auth-key", "secret-key")
usdtSvc := payment.NewUSDTService(c)

    resp, err := usdtSvc.Create(ctx, &payment.USDTRequest{
        TransactionID:  client.GenerateTransactionID("USD"),
        Username:       "user123",
        Amount:         10.50, // 10.50 USDT
    })
if err != nil {
    log.Fatal(err)
}

fmt.Printf("URL Pembayaran: %s\n", resp.PaymentURL)
fmt.Printf("ID Pembayaran Crypto: %s\n", resp.CryptoPaymentID)
```

### Cek Saldo

```go
import (
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/balance"
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/client"
)

c := client.New("auth-key", "secret-key")
balanceSvc := balance.NewService(c)

resp, err := balanceSvc.Get(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Saldo: %s\n", resp.Balance)
```

### Verifikasi Callback Pembayaran

Menangani webhook dari GSPAY2 dengan aman:

```go
import (
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/client"
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/payment"
)

c := client.New("auth-key", "secret-key")
paymentSvc := payment.NewIDRService(c)

func handleCallback(w http.ResponseWriter, r *http.Request) {
    var callback payment.IDRCallback
    if err := json.NewDecoder(r.Body).Decode(&callback); err != nil {
        http.Error(w, "Permintaan tidak valid", http.StatusBadRequest)
        return
    }

    // Verifikasi tanda tangan
    if err := paymentSvc.VerifyCallback(&callback); err != nil {
        http.Error(w, "Tanda tangan tidak valid", http.StatusUnauthorized)
        return
    }

    // Proses callback
    if callback.Status.IsSuccess() {
        // Pembayaran berhasil, perbarui status pesanan
        fmt.Printf("Pembayaran %s berhasil untuk transaksi %s\n",
            callback.IDRPaymentID, callback.TransactionID)
    }

    w.WriteHeader(http.StatusOK)
}
```

### Verifikasi Callback Pencairan

```go
import (
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/client"
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/payout"
)

c := client.New("auth-key", "secret-key")
payoutSvc := payout.NewIDRService(c)

func handlePayoutCallback(w http.ResponseWriter, r *http.Request) {
    var callback payout.IDRCallback
    if err := json.NewDecoder(r.Body).Decode(&callback); err != nil {
        http.Error(w, "Permintaan tidak valid", http.StatusBadRequest)
        return
    }

    if err := payoutSvc.VerifyCallback(&callback); err != nil {
        http.Error(w, "Tanda tangan tidak valid", http.StatusUnauthorized)
        return
    }

    // Proses pencairan yang berhasil
    fmt.Printf("Pencairan %s berhasil\n", callback.IDRPayoutID)
    w.WriteHeader(http.StatusOK)
}
```

## Penanganan Error

SDK menyediakan error bertipe untuk penanganan yang mudah:

```go
import (
    "github.com/H0llyW00dzZ/gspay-go-sdk/src/errors"
)

resp, err := paymentSvc.Create(ctx, req)
if err != nil {
    // Periksa error API
    if apiErr := errors.GetAPIError(err); apiErr != nil {
        log.Printf("Error API %d: %s", apiErr.Code, apiErr.Message)
        return
    }

    // Periksa error validasi spesifik
    if errors.Is(err, errors.ErrInvalidTransactionID) {
        log.Println("ID transaksi tidak valid")
        return
    }

    if errors.Is(err, errors.ErrInvalidAmount) {
        log.Println("Jumlah tidak valid")
        return
    }

    // Tangani error lainnya
    log.Printf("Error: %v", err)
}
```

## Pertimbangan Keamanan

### Tanda Tangan MD5

API GSPAY2 memerlukan tanda tangan berbasis MD5 untuk autentikasi request dan verifikasi callback. Meskipun berfungsi untuk keamanan API dasar, MD5 memiliki kelemahan kriptografi yang diketahui termasuk:

- **Serangan Collision**: Dapat menghasilkan hash identik untuk input berbeda
- **Serangan Preimage**: Lebih mudah dibalik dibandingkan algoritma modern
- **Serangan Rainbow Table**: Rentan terhadap tabel lookup yang telah dikomputasi sebelumnya

**Penting**: Ini adalah **persyaratan dari penyedia API GSPAY2**, bukan pilihan dalam implementasi kami. Kami mengimplementasikan tanda tangan MD5 persis seperti yang ditentukan dalam dokumentasi mereka.

### Praktik Terbaik Keamanan

Untuk meningkatkan keamanan meskipun ada keterbatasan MD5:

1. **Selalu Gunakan HTTPS**: Pastikan semua komunikasi API menggunakan TLS 1.3 atau lebih tinggi
2. **Implementasikan Rate Limiting**: Lindungi dari serangan brute force dan replay
3. **Sertakan Timestamp**: Tambahkan validasi timestamp untuk mencegah serangan replay
4. **Verifikasi Callback**: Selalu verifikasi tanda tangan webhook sebelum memproses
5. **IP Whitelisting**: Gunakan `VerifyCallbackWithIP()` untuk validasi berbasis IP tambahan
6. **Request Signing**: Kombinasikan dengan HMAC jika lapisan keamanan tambahan diperlukan

### Langkah Keamanan Tambahan yang Direkomendasikan

```go
// Contoh: Tambahkan validasi timestamp
func validateTimestamp(timestamp int64) bool {
    now := time.Now().Unix()
    // Izinkan jendela 5 menit untuk perbedaan waktu
    return timestamp >= now-300 && timestamp <= now+300
}

// Contoh: Rate limiting
var requestLimiter = tollbooth.NewLimiter(10, nil) // 10 request/detik
```

### Keamanan Transport

- Semua endpoint API menggunakan HTTPS secara default
- Validasi sertifikat TLS diaktifkan
- Tidak ada data sensitif yang dicatat dalam teks biasa

### Keamanan Callback

SDK menyediakan verifikasi callback yang robust:

```go
// Selalu verifikasi tanda tangan
if err := paymentSvc.VerifyCallback(&callback); err != nil {
    // Tolak callback yang tidak valid
    return
}

// Opsional: Tambahkan IP whitelisting
if err := paymentSvc.VerifyCallbackWithIP(&callback, clientIP); err != nil {
    // Tolak dari IP yang tidak diotorisasi
    return
}
```

**Catatan**: Meskipun MD5 menyediakan pemeriksaan integritas dasar, pertimbangkan untuk mengimplementasikan lapisan keamanan tambahan untuk transaksi bernilai tinggi atau deployment enterprise.

## Bank & E-Wallet yang Didukung

### Indonesia (IDR)

| Kode | Nama Bank |
|------|-----------|
| `BCA` | Bank BCA |
| `BRI` | Bank BRI |
| `MANDIRI` | Bank Mandiri |
| `BNI` | Bank BNI |
| `CIMB` | Bank CIMB Niaga |
| `PERMATA` | Bank Permata |
| `DANAMON` | Bank Danamon Indonesia |

### E-Wallet (IDR)

| Kode | Nama E-Wallet |
|------|---------------|
| `DANA` | DANA |
| `OVO` | OVO |

### Malaysia (MYR)

| Kode | Nama Bank |
|------|-----------|
| `MBB` | Maybank |
| `CIMB` | CIMB |
| `PBB` | Public Bank |
| `HLB` | Hong Leong Bank |
| `RHB` | RHB |
| `TNG` | Touch n Go eWallet |
| ... | [Lihat daftar lengkap](src/constants/banks.go) |

### Thailand (THB)

| Kode | Nama Bank |
|------|-----------|
| `BBL` | Bangkok Bank |
| `KBANK` | Kasikornbank |
| `KTB` | Krung Thai Bank |
| `SCB` | Siam Commercial Bank |
| ... | [Lihat daftar lengkap](src/constants/banks.go) |

## Channel Pembayaran (IDR)

| Channel | Deskripsi |
|---------|-----------|
| `constants.ChannelQRIS` | Pembayaran QR QRIS (Pembayaran Successor Besar oleh [Bank Indonesia](https://www.bi.go.id/)) |
| `constants.ChannelDANA` | E-Wallet DANA |
| `constants.ChannelBNI` | Virtual Account BNI |

## Status Pembayaran

| Status | Nilai | Deskripsi |
|--------|-------|-----------|
| `constants.StatusPending` | 0 | Pembayaran pending atau kedaluwarsa |
| `constants.StatusSuccess` | 1 | Pembayaran berhasil |
| `constants.StatusFailed` | 2 | Pembayaran gagal |
| `constants.StatusTimeout` | 4 | Pembayaran timeout |

```go
// Cek status menggunakan method helper
if status.IsSuccess() {
    // Pembayaran selesai
}

if status.IsFailed() {
    // Pembayaran gagal atau timeout
}

if status.IsPending() {
    // Pembayaran masih pending
}

// Dapatkan label yang mudah dibaca
fmt.Println(status.String()) // "Success", "Pending/Expired", dll.
```

## Fungsi Helper

### Generate ID Transaksi

```go
// Generate ID transaksi unik (maks 20 karakter)
txnID := client.GenerateTransactionID("TXN")
// Hasil: "TXN20260126143022123"

// Generate ID transaksi berbasis UUID yang aman secara kriptografis
txnID := client.GenerateUUIDTransactionID("TXN")
// Hasil: "TXN3d66c16c9db64210a"
```

### Build Return URL

```go
// Tambahkan return URL ke URL pembayaran
fullURL := client.BuildReturnURL(paymentURL, "https://mysite.com/complete")
```

### Format Mata Uang

```go
// Format jumlah IDR
formatted := client.FormatAmountIDR(50000)
// Hasil: "Rp 50.000"

// Format jumlah USDT
formatted := client.FormatAmountUSDT(10.50)
// Hasil: "10.50 USDT"
```

### Utilitas Bank

```go
// Periksa apakah kode bank valid
if constants.IsValidBankIDR("BCA") {
    // Bank Indonesia yang valid
}

// Dapatkan nama bank
name := constants.GetBankName("BCA", constants.CurrencyIDR)
// Hasil: "Bank BCA"

// Dapatkan semua kode bank untuk mata uang tertentu
codes := constants.GetBankCodes(constants.CurrencyIDR)
```

## Pengujian

Jalankan semua test:

```bash
go test ./... -v
```

Jalankan test dengan coverage:

```bash
go test ./... -cover
```

## 🚧 Roadmap & TODO

### **Ekspansi Metode Pembayaran**
SDK saat ini mendukung pembayaran **Indonesia (IDR)**. Rilis mendatang akan menambahkan dukungan untuk pasar APAC tambahan:

- [ ] **Dukungan Pembayaran Thailand (THB)**
  - [ ] Implementasi layanan pembayaran THB (`src/payment/thb.go`)
  - [ ] Tambahkan verifikasi callback THB
  - [ ] Dukungan transfer bank THB dan pembayaran QR
  - [ ] Tambahkan test pembayaran THB

- [ ] **Dukungan Pembayaran Malaysia (MYR)**
  - [ ] Implementasi layanan pembayaran MYR (`src/payment/myr.go`)
  - [ ] Tambahkan verifikasi callback MYR
  - [ ] Dukungan transfer bank MYR dan DuitNow
  - [ ] Tambahkan test pembayaran MYR


### **Backlog Enhancement**
- [ ] Tambahkan middleware verifikasi tanda tangan webhook
- [ ] Implementasi polling status pembayaran dengan webhook
- [ ] Tambahkan rate limiting dan request throttling
- [ ] Dukungan untuk HTTP client kustom dan proxy
- [ ] Tambahkan logging dan metrik yang komprehensif
- [ ] Implementasi utilitas rekonsiliasi pembayaran
- [ ] Tambahkan dukungan untuk refund parsial (jika didukung oleh API)
- [ ] Query saldo multi-mata uang

### **Kontribusi**
Kontribusi untuk memperluas dukungan metode pembayaran sangat diterima! Silakan lihat [Panduan Kontribusi](CONTRIBUTING.md) untuk detailnya.

## Penyangkalan

Ini adalah SDK **tidak resmi**. SDK ini tidak berafiliasi dengan, didukung oleh, atau secara resmi didukung oleh GSPAY atau perusahaan induknya. SDK ini dikembangkan secara independen oleh komunitas untuk menyediakan kompatibilitas bahasa Go dalam mengintegrasikan API Payment Gateway GSPAY2.

Penulis SDK ini tidak bertanggung jawab atas masalah yang timbul dari penggunaannya. Harap pastikan Anda memahami syarat layanan API GSPAY2 sebelum menggunakan SDK ini dalam produksi.

## Lisensi

Proyek ini dilisensikan di bawah Apache License 2.0 - lihat file [LICENSE](LICENSE) untuk detailnya.
