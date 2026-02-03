// Copyright 2026 H0llyW00dzZ
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package i18n

// MessageKey identifies a translatable message.
type MessageKey string

// Message keys for SDK errors and validation messages.
const (
	// Sentinel error messages.
	MsgInvalidTransactionID MessageKey = "invalid_transaction_id"
	MsgInvalidAmount        MessageKey = "invalid_amount"
	MsgInvalidBankCode      MessageKey = "invalid_bank_code"
	MsgInvalidSignature     MessageKey = "invalid_signature"
	MsgMissingCallbackField MessageKey = "missing_callback_field"
	MsgEmptyResponse        MessageKey = "empty_response"
	MsgInvalidJSON          MessageKey = "invalid_json"
	MsgRequestFailed        MessageKey = "request_failed"
	MsgIPNotWhitelisted     MessageKey = "ip_not_whitelisted"
	MsgInvalidIPAddress     MessageKey = "invalid_ip_address"
	MsgRateLimited          MessageKey = "rate_limited"

	// Validation error messages.
	MsgMinAmountIDR          MessageKey = "min_amount_idr"
	MsgMinAmountUSDT         MessageKey = "min_amount_usdt"
	MsgMinPayoutAmountIDR    MessageKey = "min_payout_amount_idr"
	MsgInvalidAmountFormat   MessageKey = "invalid_amount_format"
	MsgValidationErrorFormat MessageKey = "validation_error_format"
	MsgAPIErrorFormat        MessageKey = "api_error_format"
	MsgAPIErrorFormatNoURL   MessageKey = "api_error_format_no_url"

	// Request retry messages.
	MsgRequestFailedAfterRetries MessageKey = "request_failed_after_retries"

	// Log messages - IDR Payment.
	LogCreatingIDRPayment         MessageKey = "log_creating_idr_payment"
	LogIDRPaymentCreated          MessageKey = "log_idr_payment_created"
	LogQueryingIDRPaymentStatus   MessageKey = "log_querying_idr_payment_status"
	LogIDRPaymentStatusRetrieved  MessageKey = "log_idr_payment_status_retrieved"
	LogVerifyingIDRSignature      MessageKey = "log_verifying_idr_signature"
	LogIDRSignatureVerified       MessageKey = "log_idr_signature_verified"
	LogVerifyingIDRStatusSig      MessageKey = "log_verifying_idr_status_signature"
	LogIDRStatusSigVerified       MessageKey = "log_idr_status_signature_verified"
	LogVerifyingIDRCallback       MessageKey = "log_verifying_idr_callback"
	LogIDRCallbackVerified        MessageKey = "log_idr_callback_verified"
	LogIDRSigVerifyFailedMissing  MessageKey = "log_idr_sig_verify_failed_missing"
	LogIDRSigVerifyFailedFormat   MessageKey = "log_idr_sig_verify_failed_format"
	LogIDRSigVerifyFailedMismatch MessageKey = "log_idr_sig_verify_failed_mismatch"
	LogIDRCallbackIPFailed        MessageKey = "log_idr_callback_ip_failed"

	// Log messages - USDT Payment.
	LogCreatingUSDTPayment         MessageKey = "log_creating_usdt_payment"
	LogUSDTPaymentCreated          MessageKey = "log_usdt_payment_created"
	LogVerifyingUSDTSignature      MessageKey = "log_verifying_usdt_signature"
	LogUSDTSignatureVerified       MessageKey = "log_usdt_signature_verified"
	LogVerifyingUSDTCallback       MessageKey = "log_verifying_usdt_callback"
	LogUSDTCallbackVerified        MessageKey = "log_usdt_callback_verified"
	LogUSDTSigVerifyFailedMissing  MessageKey = "log_usdt_sig_verify_failed_missing"
	LogUSDTSigVerifyFailedFormat   MessageKey = "log_usdt_sig_verify_failed_format"
	LogUSDTSigVerifyFailedMismatch MessageKey = "log_usdt_sig_verify_failed_mismatch"
	LogUSDTCallbackIPFailed        MessageKey = "log_usdt_callback_ip_failed"

	// Log messages - IDR Payout.
	LogCreatingIDRPayout          MessageKey = "log_creating_idr_payout"
	LogIDRPayoutCreated           MessageKey = "log_idr_payout_created"
	LogQueryingIDRPayoutStatus    MessageKey = "log_querying_idr_payout_status"
	LogIDRPayoutStatusRetrieved   MessageKey = "log_idr_payout_status_retrieved"
	LogVerifyingIDRPayoutSig      MessageKey = "log_verifying_idr_payout_signature"
	LogIDRPayoutSigVerified       MessageKey = "log_idr_payout_signature_verified"
	LogVerifyingIDRPayoutCallback MessageKey = "log_verifying_idr_payout_callback"
	LogIDRPayoutCallbackVerified  MessageKey = "log_idr_payout_callback_verified"
	LogIDRPayoutSigFailedMissing  MessageKey = "log_idr_payout_sig_failed_missing"
	LogIDRPayoutSigFailedFormat   MessageKey = "log_idr_payout_sig_failed_format"
	LogIDRPayoutSigFailedMismatch MessageKey = "log_idr_payout_sig_failed_mismatch"
	LogIDRPayoutCallbackIPFailed  MessageKey = "log_idr_payout_callback_ip_failed"

	// Log messages - Balance.
	LogQueryingBalance  MessageKey = "log_querying_balance"
	LogBalanceRetrieved MessageKey = "log_balance_retrieved"

	// Log messages - HTTP Request.
	LogHTTPErrorResponse   MessageKey = "log_http_error_response"
	LogAPIResponseReceived MessageKey = "log_api_response_received"
	LogSendingRequest      MessageKey = "log_sending_request"
	LogRequestFailed       MessageKey = "log_request_failed"
	LogRequestCompleted    MessageKey = "log_request_completed"
	LogRetryingRequest     MessageKey = "log_retrying_request"
	LogRetryableError      MessageKey = "log_retryable_error"
	LogRateLimitedRetry    MessageKey = "log_rate_limited_retry"

	// HTTP Error message (for APIError.Message field).
	MsgHTTPError MessageKey = "http_error"
)

// translations holds all translated messages indexed by language and message key.
var translations = map[Language]map[MessageKey]string{
	English: {
		// Sentinel errors
		MsgInvalidTransactionID: "transaction ID must be 5-20 characters",
		MsgInvalidAmount:        "invalid payment amount",
		MsgInvalidBankCode:      "invalid bank code",
		MsgInvalidSignature:     "invalid signature",
		MsgMissingCallbackField: "missing required callback field",
		MsgEmptyResponse:        "empty response from API",
		MsgInvalidJSON:          "invalid JSON response",
		MsgRequestFailed:        "request failed",
		MsgIPNotWhitelisted:     "IP address not whitelisted",
		MsgInvalidIPAddress:     "invalid IP address format",
		MsgRateLimited:          "rate limited by API",

		// Validation errors
		MsgMinAmountIDR:          "minimum amount is 10000 IDR",
		MsgMinAmountUSDT:         "minimum amount is 1.00 USDT",
		MsgMinPayoutAmountIDR:    "minimum payout amount is 10000 IDR",
		MsgInvalidAmountFormat:   "invalid amount format",
		MsgValidationErrorFormat: "gspay: validation error for %s: %s",
		MsgAPIErrorFormat:        "gspay: API error %d on %s: %s",
		MsgAPIErrorFormatNoURL:   "gspay: API error %d: %s",

		// Request retry messages
		MsgRequestFailedAfterRetries: "request failed after %d retries",

		// Log messages - IDR Payment
		LogCreatingIDRPayment:         "creating IDR payment",
		LogIDRPaymentCreated:          "IDR payment created",
		LogQueryingIDRPaymentStatus:   "querying IDR payment status",
		LogIDRPaymentStatusRetrieved:  "IDR payment status retrieved",
		LogVerifyingIDRSignature:      "verifying IDR payment signature",
		LogIDRSignatureVerified:       "IDR payment signature verified",
		LogVerifyingIDRStatusSig:      "verifying IDR status signature",
		LogIDRStatusSigVerified:       "IDR status signature verified",
		LogVerifyingIDRCallback:       "verifying IDR callback",
		LogIDRCallbackVerified:        "IDR callback verified",
		LogIDRSigVerifyFailedMissing:  "IDR signature verification failed: missing field",
		LogIDRSigVerifyFailedFormat:   "IDR signature verification failed: invalid amount format",
		LogIDRSigVerifyFailedMismatch: "IDR signature verification failed: signature mismatch",
		LogIDRCallbackIPFailed:        "IDR callback IP verification failed",

		// Log messages - USDT Payment
		LogCreatingUSDTPayment:         "creating USDT payment",
		LogUSDTPaymentCreated:          "USDT payment created",
		LogVerifyingUSDTSignature:      "verifying USDT payment signature",
		LogUSDTSignatureVerified:       "USDT payment signature verified",
		LogVerifyingUSDTCallback:       "verifying USDT callback",
		LogUSDTCallbackVerified:        "USDT callback verified",
		LogUSDTSigVerifyFailedMissing:  "USDT signature verification failed: missing field",
		LogUSDTSigVerifyFailedFormat:   "USDT signature verification failed: invalid amount format",
		LogUSDTSigVerifyFailedMismatch: "USDT signature verification failed: signature mismatch",
		LogUSDTCallbackIPFailed:        "USDT callback IP verification failed",

		// Log messages - IDR Payout
		LogCreatingIDRPayout:          "creating IDR payout",
		LogIDRPayoutCreated:           "IDR payout created",
		LogQueryingIDRPayoutStatus:    "querying IDR payout status",
		LogIDRPayoutStatusRetrieved:   "IDR payout status retrieved",
		LogVerifyingIDRPayoutSig:      "verifying IDR payout signature",
		LogIDRPayoutSigVerified:       "IDR payout signature verified",
		LogVerifyingIDRPayoutCallback: "verifying IDR payout callback",
		LogIDRPayoutCallbackVerified:  "IDR payout callback verified",
		LogIDRPayoutSigFailedMissing:  "IDR payout signature verification failed: missing field",
		LogIDRPayoutSigFailedFormat:   "IDR payout signature verification failed: invalid amount format",
		LogIDRPayoutSigFailedMismatch: "IDR payout signature verification failed: signature mismatch",
		LogIDRPayoutCallbackIPFailed:  "IDR payout callback IP verification failed",

		// Log messages - Balance
		LogQueryingBalance:  "querying operator balance",
		LogBalanceRetrieved: "balance retrieved",

		// Log messages - HTTP Request
		LogHTTPErrorResponse:   "HTTP error response",
		LogAPIResponseReceived: "API response received",
		LogSendingRequest:      "sending request",
		LogRequestFailed:       "request failed",
		LogRequestCompleted:    "request completed successfully",
		LogRetryingRequest:     "retrying request",
		LogRetryableError:      "retryable error occurred",
		LogRateLimitedRetry:    "rate limited, waiting before retry",

		// HTTP Error message
		MsgHTTPError: "HTTP Error: %d",
	},
	Indonesian: {
		// Sentinel errors
		MsgInvalidTransactionID: "ID transaksi harus 5-20 karakter",
		MsgInvalidAmount:        "jumlah pembayaran tidak valid",
		MsgInvalidBankCode:      "kode bank tidak valid",
		MsgInvalidSignature:     "tanda tangan tidak valid",
		MsgMissingCallbackField: "field callback yang diperlukan tidak ada",
		MsgEmptyResponse:        "respons kosong dari API",
		MsgInvalidJSON:          "respons JSON tidak valid",
		MsgRequestFailed:        "permintaan gagal",
		MsgIPNotWhitelisted:     "alamat IP tidak ada dalam whitelist",
		MsgInvalidIPAddress:     "format alamat IP tidak valid",
		MsgRateLimited:          "dibatasi oleh API",

		// Validation errors
		MsgMinAmountIDR:          "jumlah minimum adalah 10000 IDR",
		MsgMinAmountUSDT:         "jumlah minimum adalah 1.00 USDT",
		MsgMinPayoutAmountIDR:    "jumlah pembayaran minimum adalah 10000 IDR",
		MsgInvalidAmountFormat:   "format jumlah tidak valid",
		MsgValidationErrorFormat: "gspay: kesalahan validasi untuk %s: %s",
		MsgAPIErrorFormat:        "gspay: kesalahan API %d pada %s: %s",
		MsgAPIErrorFormatNoURL:   "gspay: kesalahan API %d: %s",

		// Request retry messages
		MsgRequestFailedAfterRetries: "permintaan gagal setelah %d percobaan",

		// Log messages - IDR Payment
		LogCreatingIDRPayment:         "membuat pembayaran IDR",
		LogIDRPaymentCreated:          "pembayaran IDR berhasil dibuat",
		LogQueryingIDRPaymentStatus:   "mengambil status pembayaran IDR",
		LogIDRPaymentStatusRetrieved:  "status pembayaran IDR berhasil diambil",
		LogVerifyingIDRSignature:      "memverifikasi tanda tangan pembayaran IDR",
		LogIDRSignatureVerified:       "tanda tangan pembayaran IDR terverifikasi",
		LogVerifyingIDRStatusSig:      "memverifikasi tanda tangan status IDR",
		LogIDRStatusSigVerified:       "tanda tangan status IDR terverifikasi",
		LogVerifyingIDRCallback:       "memverifikasi callback IDR",
		LogIDRCallbackVerified:        "callback IDR terverifikasi",
		LogIDRSigVerifyFailedMissing:  "verifikasi tanda tangan IDR gagal: field tidak ada",
		LogIDRSigVerifyFailedFormat:   "verifikasi tanda tangan IDR gagal: format jumlah tidak valid",
		LogIDRSigVerifyFailedMismatch: "verifikasi tanda tangan IDR gagal: tanda tangan tidak cocok",
		LogIDRCallbackIPFailed:        "verifikasi IP callback IDR gagal",

		// Log messages - USDT Payment
		LogCreatingUSDTPayment:         "membuat pembayaran USDT",
		LogUSDTPaymentCreated:          "pembayaran USDT berhasil dibuat",
		LogVerifyingUSDTSignature:      "memverifikasi tanda tangan pembayaran USDT",
		LogUSDTSignatureVerified:       "tanda tangan pembayaran USDT terverifikasi",
		LogVerifyingUSDTCallback:       "memverifikasi callback USDT",
		LogUSDTCallbackVerified:        "callback USDT terverifikasi",
		LogUSDTSigVerifyFailedMissing:  "verifikasi tanda tangan USDT gagal: field tidak ada",
		LogUSDTSigVerifyFailedFormat:   "verifikasi tanda tangan USDT gagal: format jumlah tidak valid",
		LogUSDTSigVerifyFailedMismatch: "verifikasi tanda tangan USDT gagal: tanda tangan tidak cocok",
		LogUSDTCallbackIPFailed:        "verifikasi IP callback USDT gagal",

		// Log messages - IDR Payout
		LogCreatingIDRPayout:          "membuat penarikan IDR",
		LogIDRPayoutCreated:           "penarikan IDR berhasil dibuat",
		LogQueryingIDRPayoutStatus:    "mengambil status penarikan IDR",
		LogIDRPayoutStatusRetrieved:   "status penarikan IDR berhasil diambil",
		LogVerifyingIDRPayoutSig:      "memverifikasi tanda tangan penarikan IDR",
		LogIDRPayoutSigVerified:       "tanda tangan penarikan IDR terverifikasi",
		LogVerifyingIDRPayoutCallback: "memverifikasi callback penarikan IDR",
		LogIDRPayoutCallbackVerified:  "callback penarikan IDR terverifikasi",
		LogIDRPayoutSigFailedMissing:  "verifikasi tanda tangan penarikan IDR gagal: field tidak ada",
		LogIDRPayoutSigFailedFormat:   "verifikasi tanda tangan penarikan IDR gagal: format jumlah tidak valid",
		LogIDRPayoutSigFailedMismatch: "verifikasi tanda tangan penarikan IDR gagal: tanda tangan tidak cocok",
		LogIDRPayoutCallbackIPFailed:  "verifikasi IP callback penarikan IDR gagal",

		// Log messages - Balance
		LogQueryingBalance:  "mengambil saldo operator",
		LogBalanceRetrieved: "saldo berhasil diambil",

		// Log messages - HTTP Request
		LogHTTPErrorResponse:   "respons error HTTP",
		LogAPIResponseReceived: "respons API diterima",
		LogSendingRequest:      "mengirim permintaan",
		LogRequestFailed:       "permintaan gagal",
		LogRequestCompleted:    "permintaan berhasil diselesaikan",
		LogRetryingRequest:     "mencoba ulang permintaan",
		LogRetryableError:      "terjadi error yang dapat dicoba ulang",
		LogRateLimitedRetry:    "dibatasi rate limit, menunggu sebelum mencoba ulang",

		// HTTP Error message
		MsgHTTPError: "Error HTTP: %d",
	},
}

// Get retrieves a message for the specified language and key.
// Falls back to English if the language or key is not found.
func Get(lang Language, key MessageKey) string {
	if msgs, ok := translations[lang]; ok {
		if msg, ok := msgs[key]; ok {
			return msg
		}
	}
	// Fallback to English
	if msgs, ok := translations[English]; ok {
		if msg, ok := msgs[key]; ok {
			return msg
		}
	}
	// Return the key as a last resort
	return string(key)
}
