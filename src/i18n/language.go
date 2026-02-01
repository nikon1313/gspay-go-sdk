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

// Package i18n provides internationalization support for the GSPAY2 SDK.
package i18n

// Language represents a supported language for SDK messages.
type Language string

// Supported languages.
const (
	// English is the default language.
	English Language = "en"
	// Indonesian language.
	Indonesian Language = "id"
)

// IsValid returns true if the language is supported.
func (l Language) IsValid() bool {
	switch l {
	case English, Indonesian:
		return true
	default:
		return false
	}
}

// String returns the language code as a string.
func (l Language) String() string {
	return string(l)
}
