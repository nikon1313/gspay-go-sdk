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

// Package signature provides cryptographic signature utilities for the GSPAY2 SDK.
//
// This internal package handles signature generation and verification as required
// by the GSPAY2 API specification. By default, it uses MD5 hashing, but supports
// custom digest algorithms via the [Digest] type.
//
// # Custom Digest Algorithms
//
// While MD5 is the default (as required by the GSPAY2 API), you can use custom
// hash functions by providing a [Digest] to [GenerateWithDigest]:
//
//	// Use SHA-256 instead of MD5
//	sig := signature.GenerateWithDigest(data, sha256.New)
//
//	// Use SHA-512
//	sig := signature.GenerateWithDigest(data, sha512.New)
//
// The [Digest] type accepts any function that returns a [hash.Hash] instance,
// making it compatible with all standard library hash functions:
//   - crypto/md5.New (default)
//   - crypto/sha1.New
//   - crypto/sha256.New
//   - crypto/sha512.New
//
// # Signature Formulas
//
// IDR Payment:
//
//	MD5(transaction_id + player_username + amount + operator_secret_key)
//
// IDR Payment Callback:
//
//	MD5(idrpayment_id + amount + transaction_id + status + secret_key)
//
// IDR Payout:
//
//	MD5(transaction_id + player_username + amount + account_number + operator_secret_key)
//
// IDR Payout Callback:
//
//	MD5(idrpayout_id + account_number + amount + transaction_id + secret_key)
//
// USDT Payment:
//
//	MD5(transaction_id + player_username + amount + operator_secret_key)
//
// USDT Payment Callback:
//
//	MD5(cryptopayment_id + amount + transaction_id + status + secret_key)
//
// # Security Note
//
// MD5 is used by default because it is required by the GSPAY2 API provider.
// For enhanced security, use [client.WithDigest] to configure a stronger algorithm
// if your API configuration supports it. Always use HTTPS and implement additional
// security measures for production use.
package signature
