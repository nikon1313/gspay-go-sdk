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

	// Validation error messages.
	MsgMinAmountIDR        MessageKey = "min_amount_idr"
	MsgMinAmountUSDT       MessageKey = "min_amount_usdt"
	MsgMinPayoutAmountIDR  MessageKey = "min_payout_amount_idr"
	MsgInvalidAmountFormat MessageKey = "invalid_amount_format"
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

		// Validation errors
		MsgMinAmountIDR:        "minimum amount is 10000 IDR",
		MsgMinAmountUSDT:       "minimum amount is 1.00 USDT",
		MsgMinPayoutAmountIDR:  "minimum payout amount is 10000 IDR",
		MsgInvalidAmountFormat: "invalid amount format",
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

		// Validation errors
		MsgMinAmountIDR:        "jumlah minimum adalah 10000 IDR",
		MsgMinAmountUSDT:       "jumlah minimum adalah 1.00 USDT",
		MsgMinPayoutAmountIDR:  "jumlah pembayaran minimum adalah 10000 IDR",
		MsgInvalidAmountFormat: "format jumlah tidak valid",
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
