// ABOUTME: Main package file for validating Romanian CNP numbers (rossn).
// MIT License – see LICENSE file.

package rossn

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"unicode"
)

// Validate checks if a CNP is valid according to all official Romanian rules.
// It verifies length, digit content, date, county, serial, and checksum.
// Returns nil if valid, or an error describing the failure.
func Validate(cnp string) error {
	if len(cnp) != 13 {
		return errors.New("CNP must be 13 digits")
	}
	for _, r := range cnp {
		if !unicode.IsDigit(r) {
			return errors.New("CNP must contain only digits")
		}
	}
	if !isValidDate(cnp) {
		return errors.New("invalid birth date in CNP")
	}
	if !isValidCounty(cnp) {
		return errors.New("invalid county code in CNP")
	}
	if !isValidSerial(cnp) {
		return errors.New("invalid serial number")
	}
	if !hasValidControlDigit(cnp) {
		return errors.New("invalid control digit")
	}
	return nil
}

// isValidDate checks if the CNP encodes a real, valid birth date
// according to the S digit and YYMMDD fields.
func isValidDate(cnp string) bool {
	s := cnp[0]
	yy := cnp[1:3]
	mm := cnp[3:5]
	dd := cnp[5:7]

	century := "19"
	switch s {
	case '1', '2':
		century = "19"
	case '3', '4':
		century = "18"
	case '5', '6':
		century = "20"
	case '7', '8', '9':
		century = "19"
	default:
		return false
	}

	dateStr := fmt.Sprintf("%s%s-%s-%s", century, yy, mm, dd)
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

// isValidCounty checks if the CNP encodes a valid Romanian county code (JJ).
// For JJ == "47" or "48" (historic Bucharest districts), validity is restricted
// to dates before December 19, 1979. For JJ == "70", accepts any S for birth year
// 2024 and later (SIIEASC CNPs), and only S=7,8,9 for prior years (legacy CNPs).
// Other codes are validated according to the official list.
func isValidCounty(cnp string) bool {
	county := cnp[7:9]
	s := cnp[0]
	switch county {
	case "47", "48":
		yyyy, mm, dd := cnpBirthDate(cnp)
		boundary := time.Date(1979, 12, 19, 0, 0, 0, 0, time.UTC)
		cnpDate := time.Date(yyyy, time.Month(mm), dd, 0, 0, 0, 0, time.UTC)
		return cnpDate.Before(boundary)
	case "70":
		yyyy, _, _ := cnpBirthDate(cnp)
		if yyyy >= 2024 {
			return true // After 2024: Accept for any S
		}
		// Before 2024: Only for S=7,8,9
		return s == '7' || s == '8' || s == '9'
	default:
		valid := map[string]bool{
			"01": true, "02": true, "03": true, "04": true, "05": true, "06": true,
			"07": true, "08": true, "09": true, "10": true, "11": true, "12": true,
			"13": true, "14": true, "15": true, "16": true, "17": true, "18": true,
			"19": true, "20": true, "21": true, "22": true, "23": true, "24": true,
			"25": true, "26": true, "27": true, "28": true, "29": true, "30": true,
			"31": true, "32": true, "33": true, "34": true, "35": true, "36": true,
			"37": true, "38": true, "39": true, "40": true, "41": true, "42": true,
			"43": true, "44": true, "45": true, "46": true, "51": true, "52": true,
		}
		return valid[county]
	}
}

// cnpBirthDate extracts the birth date (YYYY, MM, DD) from a CNP.
// Returns (0,0,0) if the date cannot be determined (should not happen after isValidDate passes).
func cnpBirthDate(cnp string) (year int, month int, day int) {
	s := cnp[0]
	yy := cnp[1:3]
	mm := cnp[3:5]
	dd := cnp[5:7]
	var century string
	switch s {
	case '1', '2':
		century = "19"
	case '3', '4':
		century = "18"
	case '5', '6':
		century = "20"
	case '7', '8', '9':
		century = "19"
	default:
		return 0, 0, 0
	}
	y, _ := strconv.Atoi(century + yy)
	m, _ := strconv.Atoi(mm)
	d, _ := strconv.Atoi(dd)
	return y, m, d
}

// isValidSerial checks if the NNN serial part of the CNP is in the official range 001–999.
func isValidSerial(cnp string) bool {
	serial := cnp[9:12]
	val, err := strconv.Atoi(serial)
	return err == nil && val >= 1 && val <= 999
}

// hasValidControlDigit checks the CNP control digit using the official weighting scheme.
func hasValidControlDigit(cnp string) bool {
	const weights = "279146358279"
	sum := 0
	for i := 0; i < 12; i++ {
		d, _ := strconv.Atoi(string(cnp[i]))
		w, _ := strconv.Atoi(string(weights[i]))
		sum += d * w
	}
	control := sum % 11
	if control == 10 {
		control = 1
	}
	last, _ := strconv.Atoi(string(cnp[12]))
	return control == last
}
