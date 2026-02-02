# Berkontribusi pada GSPAY Go SDK

Terima kasih atas minat Anda untuk berkontribusi pada GSPAY Go SDK! Dokumen ini menyediakan panduan dan informasi untuk kontributor.

## 🚀 Cara Berkontribusi

- **🐛 Laporan Bug**: Laporkan bug melalui [GitHub Issues](https://github.com/H0llyW00dzZ/gspay-go-sdk/issues)
- **💡 Permintaan Fitur**: Sarankan fitur baru atau perbaikan
- **📝 Dokumentasi**: Tingkatkan dokumentasi, contoh, atau panduan
- **💻 Kontribusi Kode**: Kirim pull request untuk fitur baru atau perbaikan bug
- **🧪 Pengujian**: Tambahkan test atau tingkatkan cakupan test

## 📋 Setup Development

### Prasyarat

- Go 1.25.6 atau lebih baru
- Git
- Pemahaman dasar tentang Go modules dan testing

### Langkah-langkah Setup

1. **Fork repository** di GitHub
2. **Clone fork Anda**:
   ```bash
   git clone https://github.com/USERNAME_ANDA/gspay-go-sdk.git
   cd gspay-go-sdk
   ```

3. **Install dependencies**:
   ```bash
   go mod download
   ```

4. **Jalankan tests** untuk memastikan semuanya berfungsi:
   ```bash
   go test ./...
   ```

5. **Buat feature branch**:
   ```bash
   git checkout -b feature/nama-fitur-anda
   ```

## 🏗️ Struktur Proyek

```
gspay-go-sdk/
├── src/
│   ├── balance/     # Layanan query saldo
│   ├── client/      # HTTP client dan fungsionalitas inti
│   ├── constants/   # Kode bank, status pembayaran, channel
│   ├── errors/      # Tipe error dan penanganan
│   ├── helper/      # Utilitas helper
│   │   ├── amount/  # Utilitas pemformatan jumlah
│   │   └── gc/      # Manajemen buffer pool
│   ├── i18n/        # Internasionalisasi (bahasa, terjemahan)
│   ├── internal/    # Utilitas internal (pembuatan tanda tangan)
│   ├── payment/     # Layanan pembayaran (IDR, THB/MYR mendatang)
│   └── payout/      # Layanan pencairan (IDR)
├── examples/        # Contoh penggunaan
├── go.mod           # Definisi Go module
└── README.md        # Dokumentasi utama
```

## 💻 Standar Kode

### Gaya Kode Go

- Ikuti format Go standar: `go fmt`
- Gunakan `gofmt -s` untuk penyederhanaan tambahan
- Jalankan `go vet` dan perbaiki semua peringatan
- Pastikan `golint` lolos (jika tersedia)

### Konvensi Penamaan

```go
// Tipe
type PaymentRequest struct { ... }    // PascalCase untuk tipe yang diekspor
type paymentAPIRequest struct { ... } // camelCase untuk tipe internal

// Fungsi
func CreatePayment(...) (...)         // PascalCase untuk fungsi yang diekspor
func createAPIRequest(...) (...)      // camelCase untuk fungsi internal

// Variabel
var PaymentStatusPending = 0          // PascalCase untuk konstanta yang diekspor
var defaultTimeout = 30 * time.Second // camelCase untuk variabel internal
```

### Penanganan Error

- Kembalikan error bertipe dari paket `errors`
- Gunakan `errors.New` untuk mengembalikan error sentinel dengan konteks dan lokalisasi
- Gunakan `fmt.Errorf` untuk wrapping: `return fmt.Errorf("%w: %s", errors.ErrInvalidAmount, amount)`
- Sertakan konteks dalam pesan error

### Dokumentasi

- Tambahkan komentar doc untuk semua fungsi, tipe, dan method yang diekspor
- Gunakan format Go doc yang benar
- Sertakan contoh penggunaan jika membantu

### Internasionalisasi (i18n)

Saat menambahkan pesan error yang menghadap pengguna:

1. **Tambahkan message key** di `src/i18n/messages.go`:
   ```go
   const (
       MsgNewErrorKey MessageKey = "new_error_key"
   )
   ```

2. **Tambahkan terjemahan** untuk semua bahasa yang didukung:
   ```go
   var translations = map[Language]map[MessageKey]string{
       English: {
           MsgNewErrorKey: "English error message",
       },
       Indonesian: {
           MsgNewErrorKey: "Pesan error dalam Bahasa Indonesia",
       },
   }
   ```

3. **Gunakan pesan terlokalisasi** dalam validation errors:
   ```go
   return errors.NewValidationError("field", 
       errors.GetMessage(s.client.Language, errors.KeyNewError))
   ```

4. **Re-export di paket errors** untuk kemudahan:
   ```go
   // src/errors/errors.go
   const KeyNewError = i18n.MsgNewErrorKey
   ```

### Memodifikasi Endpoint API

1. Edit `src/constants/endpoints.go` untuk menambah atau memodifikasi `EndpointKey` dan path dalam map `endpoints`.
2. Update implementasi layanan untuk menggunakan `constants.GetEndpoint()` alih-alih string hardcoded.
3. Update pembuatan tanda tangan jika parameter berubah.
4. Update struct request/response.
5. Update tests untuk memverifikasi perubahan dan memastikan coverage untuk endpoint baru.

## 🧪 Pengujian

### Persyaratan Test

- **100% coverage** untuk kode baru
- Gunakan table-driven tests untuk berbagai skenario
- Mock respons HTTP menggunakan `httptest`
- Test kasus sukses dan error
- Test edge cases dan validasi input

### Struktur Test

```go
func TestPaymentService_Create(t *testing.T) {
    t.Run("pembuatan pembayaran berhasil", func(t *testing.T) {
        // Setup mock server
        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Mock respons API
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(map[string]any{
                "code": 200,
                "data": `{"payment_url":"https://...","id":"123"}`,
            })
        }))
        defer server.Close()

        // Test kode Anda
        client := New("auth", "secret", WithBaseURL(server.URL))
        svc := payment.NewIDRService(client)

        resp, err := svc.Create(context.Background(), &payment.IDRRequest{
            TransactionID: "TXN123",
            Username:      "user123",
            Amount:        50000,
        })

        assert.NoError(t, err)
        assert.NotEmpty(t, resp.PaymentURL)
    })
}
```

### Menjalankan Tests

```bash
# Jalankan semua tests
go test ./...

# Jalankan dengan coverage
go test ./... -cover

# Jalankan paket tertentu
go test ./src/payment

# Jalankan dengan output verbose
go test ./... -v
```

## 🔧 Menambahkan Metode Pembayaran Baru

### Untuk Mata Uang Baru (contoh: THB, MYR)

1. **Tambahkan konstanta** di `src/constants/`:
   ```go
   // Kode mata uang
   const CurrencyTHB Currency = "THB"

   // Kode bank
   var BanksTHB = map[string]string{
       "BBL": "Bangkok Bank",
       // ... tambahkan bank lainnya
   }

   // Channel pembayaran
   var ChannelsTHB = []string{"QRIS", "BANK_TRANSFER"}
   ```

2. **Buat layanan pembayaran** di `src/payment/`:
   ```go
   // src/payment/thb.go
   type THBService struct { client *client.Client }

   func NewTHBService(c *client.Client) *THBService {
       return &THBService{client: c}
   }

   func (s *THBService) Create(ctx context.Context, req *THBRequest) (*THBResponse, error) {
       // Implementasi mengikuti pola layanan IDR
   }
   ```

3. **Tambahkan verifikasi callback**:
   ```go
   func (s *THBService) VerifyCallback(callback *THBCallback) error {
       // Verifikasi tanda tangan MD5
   }
   ```

4. **Update client** untuk mendukung layanan baru:
   ```go
   // Tambahkan konstruktor layanan THB
   func NewTHBService(c *Client) *THBService { ... }
   ```

5. **Tambahkan tests komprehensif** mengikuti pola yang ada

6. **Update dokumentasi** di README.md dan examples

### Checklist Implementasi

- [ ] Konstanta ditambahkan untuk mata uang, bank, channel
- [ ] Layanan pembayaran diimplementasikan dengan penanganan error yang tepat
- [ ] Verifikasi callback diimplementasikan
- [ ] Unit tests dengan 100% coverage
- [ ] Integration tests dengan mock API
- [ ] Dokumentasi diupdate
- [ ] Contoh ditambahkan
- [ ] Changelog diupdate

## 📝 Proses Pull Request

1. **Buat feature branch** dari `main`:
   ```bash
   git checkout -b feature/add-thb-support
   ```

2. **Buat perubahan Anda** mengikuti standar kode

3. **Jalankan tests** dan pastikan lolos:
   ```bash
   go test ./... -v
   go vet ./...
   ```

4. **Update dokumentasi** jika diperlukan

5. **Commit perubahan Anda** dengan pesan yang jelas:
   ```bash
   git add .
   git commit -m "feat: tambah dukungan pembayaran THB

   - Implementasi layanan pembayaran THB
   - Tambah verifikasi callback
   - Tambah tests komprehensif

   Closes #123"
   ```

6. **Push ke fork Anda**:
   ```bash
   git push origin feature/add-thb-support
   ```

7. **Buat Pull Request** di GitHub:
   - Gunakan judul dan deskripsi yang jelas
   - Referensikan issue terkait
   - Sertakan screenshot/demo untuk perubahan UI
   - Minta review dari maintainer

### Format Judul PR

```
type(scope): deskripsi

Types: feat, fix, docs, style, refactor, test, chore
Contoh:
- feat(thb): tambah dukungan pembayaran THB
- fix(callback): perbaiki bug verifikasi tanda tangan
- docs(readme): update instruksi instalasi
```

## 🐛 Laporan Bug

Saat melaporkan bug, harap sertakan:

- **Versi Go**: `go version`
- **Versi SDK**: Hash commit Git atau tag
- **Perilaku yang diharapkan**
- **Perilaku aktual**
- **Langkah-langkah untuk mereproduksi**
- **Pesan error/log**
- **Contoh kode** yang mendemonstrasikan masalah

## 💡 Permintaan Fitur

Permintaan fitur harus mencakup:

- **Kasus penggunaan**: Masalah apa yang diselesaikan?
- **Solusi yang diusulkan**: Bagaimana seharusnya bekerja?
- **Alternatif yang dipertimbangkan**: Pendekatan lain?
- **Konteks tambahan**: Screenshot, contoh, dll.

## 📜 Kode Etik

Proyek ini mengikuti kode etik untuk memastikan lingkungan yang ramah bagi semua kontributor:

- Bersikap hormat dan inklusif
- Fokus pada umpan balik yang konstruktif
- Terima tanggung jawab atas kesalahan
- Tunjukkan empati terhadap kontributor lain
- Bantu menciptakan komunitas yang positif

## 📞 Mendapatkan Bantuan

- **Dokumentasi**: Periksa README.md dan examples terlebih dahulu
- **Issues**: Cari issue yang ada sebelum membuat yang baru
- **Discussions**: Gunakan GitHub Discussions untuk pertanyaan
- **Komunitas**: Bergabung dengan komunitas Go yang relevan untuk pertanyaan umum

## 🎉 Penghargaan

Kontributor akan diakui:
- Di CHANGELOG untuk kontribusi signifikan
- Sebagai co-author pada rilis
- Di daftar kontributor proyek
- Melalui insight kontributor GitHub

Terima kasih telah berkontribusi pada GSPAY Go SDK! 🚀
