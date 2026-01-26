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

// ChannelIDR represents available payment channels for Indonesia.
type ChannelIDR string

const (
	// ChannelQRIS is the QRIS QR Payment channel.
	ChannelQRIS ChannelIDR = "QRIS"
	// ChannelDANA is the DANA E-Wallet channel.
	ChannelDANA ChannelIDR = "DANA"
	// ChannelBNI is the BNI Virtual Account channel.
	ChannelBNI ChannelIDR = "BNI"
)

// ChannelsIDR contains available payment channels for Indonesia.
var ChannelsIDR = map[ChannelIDR]string{
	ChannelQRIS: "QRIS QR Payment",
	ChannelDANA: "DANA E-Wallet",
	ChannelBNI:  "BNI Virtual Account",
}

// IsValidChannelIDR checks if a channel is valid for Indonesian payments.
func IsValidChannelIDR(channel ChannelIDR) bool {
	_, ok := ChannelsIDR[channel]
	return ok
}
