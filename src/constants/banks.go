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

// Currency represents supported currencies for bank operations.
type Currency string

// Supported currency constants for GSPAY2 payment operations.
const (
	// CurrencyIDR represents Indonesian Rupiah.
	CurrencyIDR Currency = "IDR"
	// CurrencyMYR represents Malaysian Ringgit.
	CurrencyMYR Currency = "MYR"
	// CurrencyTHB represents Thai Baht.
	CurrencyTHB Currency = "THB"
)

// BanksIDR contains Indonesian bank codes and names.
var BanksIDR = map[string]string{
	"BCA":     "Bank BCA",
	"BRI":     "Bank BRI",
	"MANDIRI": "BANK MANDIRI",
	"BNI":     "BANK BNI",
	"CIMB":    "BANK CIMB NIAGA",
	"PERMATA": "BANK PERMATA",
	"DANAMON": "BANK DANAMON INDONESIA",
	"DANA":    "DANA",
	"OVO":     "OVO",
}

// BanksMYR contains Malaysian bank codes and names.
var BanksMYR = map[string]string{
	"AFFB": "AFFIN BANK",
	"AGRO": "AGROBANK",
	"ALLB": "ALLIANCE BANK",
	"AMB":  "AMBANK",
	"BIM":  "BANK ISLAM MALAYSIA",
	"BMML": "BANK MUAMALAT",
	"BKR":  "BANK RAKYAT",
	"BSN":  "BANK SIMPANAN NASIONAL",
	"CIMB": "CIMB",
	"CITB": "CITIBANK",
	"HLB":  "HONG LEONG BANK",
	"HSBC": "HSBC",
	"MBB":  "MAYBANK",
	"OCBC": "OCBC",
	"PBB":  "PUBLIC BANK",
	"RHB":  "RHB",
	"SCB":  "STANDARD CHARTERED BANK",
	"UOB":  "UNITED OVERSEAS BANK",
	"TNG":  "TOUCH N GO EWALLET",
	"RYT":  "RYT BANK",
}

// BanksTHB contains Thai bank codes and names.
var BanksTHB = map[string]string{
	"BBL":   "BANGKOK BANK PUBLIC COMPANY LTD.",
	"KBANK": "KASIKORNBANK PUBLIC COMPANY LIMITED",
	"KTB":   "KRUNG THAI BANK PUBLIC COMPANY LTD.",
	"TMB":   "TMB THANACHART BANK PUBLIC COMPANY LIMITED",
	"SCB":   "SIAM COMMERCIAL BANK PUBLIC COMPANY LTD.",
	"CIMB":  "CIMB THAI BANK PUBLIC COMPANY LIMITED",
	"UOB":   "UNITED OVERSEAS BANK (THAI) PUBLIC COMPANY LIMITED",
	"BAY":   "BANK OF AYUDHYA PUBLIC COMPANY LTD.",
	"GSB":   "GOVERNMENT SAVINGS BANK",
	"GHB":   "THE GOVERNMENT HOUSING BANK",
	"BAAC":  "BANK FOR AGRICULTURE AND AGRICULTURAL COOPERATIVES",
	"TISCO": "TISCO BANK PUBLIC COMPANY LIMITED",
	"KKP":   "KIATNAKIN PHATRA BANK PUBLIC COMPANY LIMITED",
	"LHB":   "LAND AND HOUSES BANK PUBLIC COMPANY LIMITED",
}

// GetBankName returns the bank name for a given bank code and currency.
// Returns an empty string if the bank code is not found.
func GetBankName(bankCode string, currency Currency) string {
	switch currency {
	case CurrencyIDR:
		return BanksIDR[bankCode]
	case CurrencyMYR:
		return BanksMYR[bankCode]
	case CurrencyTHB:
		return BanksTHB[bankCode]
	default:
		return ""
	}
}

// GetBankCodes returns all bank codes for a given currency.
func GetBankCodes(currency Currency) []string {
	var banks map[string]string
	switch currency {
	case CurrencyIDR:
		banks = BanksIDR
	case CurrencyMYR:
		banks = BanksMYR
	case CurrencyTHB:
		banks = BanksTHB
	default:
		return nil
	}

	codes := make([]string, 0, len(banks))
	for code := range banks {
		codes = append(codes, code)
	}
	return codes
}

// IsValidBankIDR checks if a bank code is valid for Indonesian banks.
func IsValidBankIDR(bankCode string) bool {
	_, ok := BanksIDR[bankCode]
	return ok
}
