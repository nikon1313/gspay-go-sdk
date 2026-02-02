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

package constants

import "fmt"

// Version information for the GSPAY Go SDK.
//
// These values are used to construct the User-Agent header for HTTP requests.
const (
	// SDKName is the name of this SDK.
	SDKName = "gspay-go-sdk"

	// SDKVersion is the current version of this SDK.
	// This should be updated with each release.
	//
	// Note: the versioning its based of heap like unix time.
	SDKVersion = "0.3.7"

	// SDKRepository is the GitHub repository URL for this SDK.
	SDKRepository = "https://github.com/H0llyW00dzZ/gspay-go-sdk"
)

// UserAgent returns the formatted User-Agent string for HTTP requests.
//
// Format: gspay-go-sdk/1.0.0 (+https://github.com/H0llyW00dzZ/gspay-go-sdk)
//
// Example usage:
//
//	req.Header.Set("User-Agent", constants.UserAgent())
func UserAgent() string {
	return fmt.Sprintf("%s/%s (+%s)", SDKName, SDKVersion, SDKRepository)
}
