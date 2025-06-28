// ABOUTME: Exhaustive test suite for rossn CNP validation, only real/possible CNPs are accepted.
package rossn

import (
	"strconv"
	"testing"
)

func buildCNP(s, year, month, day, county, serial string) string {
	base := s + year + month + day + county + serial
	control := calculateControlDigit(base)
	return base + strconv.Itoa(control)
}

func calculateControlDigit(first12 string) int {
	weights := []int{2, 7, 9, 1, 4, 6, 3, 5, 8, 2, 7, 9}
	sum := 0
	for i := 0; i < 12; i++ {
		d, _ := strconv.Atoi(string(first12[i]))
		sum += d * weights[i]
	}
	mod := sum % 11
	if mod == 10 {
		return 1
	}
	return mod
}

func TestValidate_ValidCases(t *testing.T) {
	validCases := []struct {
		S      string
		Year   string
		Month  string
		Day    string
		County string
		Serial string
	}{
		// S = 1, 2: 1900–1999, male/female
		{"1", "80", "01", "01", "01", "001"}, // Male, 1980-01-01, Alba, serial 001
		{"2", "80", "02", "29", "40", "123"}, // Female, 1980-02-29, Suceava, serial 123 (leap year)

		// S = 3, 4: 1800–1899, male/female
		{"3", "80", "06", "15", "02", "321"}, // Male, 1880-06-15, Arad
		{"4", "99", "12", "31", "39", "999"}, // Female, 1899-12-31, Vaslui

		// S = 5, 6: 2000–2099, male/female
		{"5", "01", "01", "01", "30", "101"}, // Male, 2001-01-01, Olt
		{"6", "04", "07", "07", "52", "456"}, // Female, 2004-07-07, Giurgiu

		// S = 7, 8: foreign male/female, resident, 19xx (treated as 1900-1999)
		{"7", "85", "03", "12", "05", "111"}, // Foreign male, 1985-03-12, Bihor
		{"8", "99", "12", "31", "46", "789"}, // Foreign female, 1999-12-31, Calarasi

		// S = 9: Non-residents (date in 1900-1999, often used as 1900-01-01)
		{"9", "90", "01", "01", "51", "555"}, // Non-resident, 1990-01-01, Bucharest sector 1
	}

	for _, tc := range validCases {
		cnp := buildCNP(tc.S, tc.Year, tc.Month, tc.Day, tc.County, tc.Serial)
		if err := Validate(cnp); err != nil {
			t.Errorf("Valid CNP should pass: %s got error: %v", cnp, err)
		}
	}
}

func TestValidate_InvalidCases(t *testing.T) {
	invalidCases := []struct {
		cnp    string
		reason string
	}{
		// Length
		{"19801010012", "too short"},
		{"19801010012345", "too long"},
		// Non-numeric
		{"19X0101000123", "contains letter"},
		// S out of range
		{"0980101000123", "S=0 invalid"},
		{"a980101000123", "S=a invalid"},
		// Impossible dates
		{buildCNP("1", "80", "02", "30", "01", "001"), "Feb 30"},
		{buildCNP("5", "03", "04", "31", "45", "011"), "April 31"},
		{buildCNP("3", "88", "00", "15", "15", "111"), "month=00"},
		{buildCNP("2", "99", "05", "00", "10", "999"), "day=00"},
		// Non-leap year Feb 29
		{buildCNP("1", "81", "02", "29", "09", "101"), "1981-02-29 invalid"},
		// Leap year Feb 29 valid (to confirm negative above)
		// Invalid county
		{buildCNP("1", "80", "01", "01", "00", "001"), "county=00 invalid"},
		{buildCNP("1", "80", "01", "01", "53", "001"), "county=53 invalid"},
		// Serial out of range (should be 001–999)
		{buildCNP("1", "80", "01", "01", "01", "000"), "serial=000"},
		{buildCNP("1", "80", "01", "01", "01", "1000"), "serial=1000 too long"},
		// Bad checksum (alter control digit)
		{func() string {
			real := buildCNP("1", "80", "01", "01", "01", "001")
			bad := real[:12] + "9"
			return bad
		}(), "wrong checksum"},
	}

	for _, tc := range invalidCases {
		if err := Validate(tc.cnp); err == nil {
			t.Errorf("Invalid CNP should fail: %s [%s]", tc.cnp, tc.reason)
		}
	}
}

// Extra: leap year boundaries
func TestValidate_LeapYears(t *testing.T) {
	leap := buildCNP("1", "84", "02", "29", "01", "111") // 1984-02-29
	if err := Validate(leap); err != nil {
		t.Errorf("Valid leap year CNP failed: %s", leap)
	}
	notLeap := buildCNP("1", "81", "02", "29", "01", "111") // 1981-02-29
	if err := Validate(notLeap); err == nil {
		t.Errorf("Invalid non-leap year CNP passed: %s", notLeap)
	}
}

// Valid county codes (official)
func TestValidate_AllValidCounties(t *testing.T) {
	validCounties := []string{
		"01", "02", "03", "04", "05", "06", "07", "08", "09", "10",
		"11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
		"21", "22", "23", "24", "25", "26", "27", "28", "29", "30",
		"31", "32", "33", "34", "35", "36", "37", "38", "39", "40",
		"41", "42", "43", "44", "45", "46", "51", "52",
	}
	for _, county := range validCounties {
		cnp := buildCNP("1", "95", "12", "15", county, "123")
		if err := Validate(cnp); err != nil {
			t.Errorf("CNP with valid county %s should pass: %s", county, cnp)
		}
	}
}

// S digit coverage
func TestValidate_AllSdigits(t *testing.T) {
	for _, s := range []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"} {
		cnp := buildCNP(s, "90", "10", "10", "10", "101")
		if err := Validate(cnp); err != nil {
			t.Errorf("CNP with S=%s should be valid: %s, got %v", s, cnp, err)
		}
	}
}

// Serial edge: first and last valid
func TestValidate_SerialEdges(t *testing.T) {
	// Should be valid: serial 001 and 999
	cnp1 := buildCNP("1", "90", "01", "01", "01", "001")
	cnp2 := buildCNP("1", "90", "01", "01", "01", "999")
	if err := Validate(cnp1); err != nil {
		t.Errorf("Serial 001 should be valid: %s", cnp1)
	}
	if err := Validate(cnp2); err != nil {
		t.Errorf("Serial 999 should be valid: %s", cnp2)
	}
}

func TestValidate_InvalidFormat(t *testing.T) {
	badInputs := []string{
		" 1981214320015",  // leading space
		"1981214320015 ",  // trailing space
		"19812 14320015",  // inner space
		"\t1981214320015", // tab
		"\n1981214320015", // newline
		"1981214-320015",  // hyphen
		"19812143200.5",   // dot
		"19812143200,5",   // comma
		"",                // empty
	}

	for _, cnp := range badInputs {
		if err := Validate(cnp); err == nil {
			t.Errorf("CNP with invalid format should fail: [%q]", cnp)
		}
	}
}

// This test ensures the special checksum case (sum%11==10) → control digit 1 is handled
func TestValidate_ControlDigit10(t *testing.T) {
	// We'll brute-force a valid CNP with checksum==10
	// This is rare, but "279146358279" as weights for "279146358279"
	// One such example is S=1, YY=80, MM=01, DD=01, JJ=13, NNN=923
	// Let's build it:

	base := "180010113923" // S,YY,MM,DD,JJ,NNN (first 12 digits)
	control := calculateControlDigit(base)
	if control != 1 {
		t.Fatalf("Control digit should be 1 for base=%s, got %d", base, control)
	}
	cnp := base + "1"
	if err := Validate(cnp); err != nil {
		t.Errorf("CNP with control digit 1 (checksum==10) should be valid: %s", cnp)
	}
}

func TestValidate_ArchivalBucharestDistricts(t *testing.T) {
	// Valid: code 47, date before 1979-12-19
	cnpValid47 := buildCNP("1", "79", "12", "18", "47", "123") // 1979-12-18
	if err := Validate(cnpValid47); err != nil {
		t.Errorf("Historic CNP with JJ=47 before cutoff should pass: %s, err=%v", cnpValid47, err)
	}

	// Valid: code 48, date before 1979-12-19
	cnpValid48 := buildCNP("2", "78", "05", "10", "48", "555") // 1978-05-10
	if err := Validate(cnpValid48); err != nil {
		t.Errorf("Historic CNP with JJ=48 before cutoff should pass: %s, err=%v", cnpValid48, err)
	}

	// Invalid: code 47, date at or after 1979-12-19
	cnpInvalid47 := buildCNP("1", "79", "12", "19", "47", "123") // 1979-12-19
	if err := Validate(cnpInvalid47); err == nil {
		t.Errorf("Historic CNP with JJ=47 at/after cutoff should fail: %s", cnpInvalid47)
	}
	cnpInvalid47b := buildCNP("1", "80", "01", "01", "47", "123") // 1980-01-01
	if err := Validate(cnpInvalid47b); err == nil {
		t.Errorf("Historic CNP with JJ=47 after cutoff should fail: %s", cnpInvalid47b)
	}

	// Invalid: code 48, date at or after 1979-12-19
	cnpInvalid48 := buildCNP("2", "79", "12", "19", "48", "123") // 1979-12-19
	if err := Validate(cnpInvalid48); err == nil {
		t.Errorf("Historic CNP with JJ=48 at/after cutoff should fail: %s", cnpInvalid48)
	}
}
