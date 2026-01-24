package account

import "testing"

func TestIsValidDocumentNumber(t *testing.T) {
	tests := []struct {
		name           string
		documentNumber string
		wantErr        bool
	}{
		{"must be an invalid document number", "123", true},
		{"must be an invalid document number containing some non digit characters", "30603Z97A00120", true},
		{"must validate a CPF with only digits", "69398400014", false},
		{"must validate a CPF with non digit characters", "688.862.380-70", false},
		{"must validate a CNPJ with only digits", "30603597000120", false},
		{"must validate a CNPJ with non digit characters", "46.047.310/0001-62", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := IsValidDocumentNumber(tt.documentNumber); (err != nil) != tt.wantErr {
				t.Errorf("IsValidDocumentNumber() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeDocumentNumber(t *testing.T) {
	tests := []struct {
		name           string
		documentNumber string
		want           string
	}{
		{"must sanitize a valid CNPJ document number with separators", "76.466.403/0001-92", "76466403000192"},
		{"must sanitize a valid CNPJ document number with no separators", "76.466.403/0001-92", "76466403000192"},
		{"must sanitize a invalid document number with only characters", "aaa", ""},

		{"must sanitize a valid CPF document number with separators", "906.201.240-08", "90620124008"},
		{"must sanitize a valid CPF document number with no separators", "67885323000180", "67885323000180"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitizeDocumentNumber(tt.documentNumber); got != tt.want {
				t.Errorf("SanitizeDocumentNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
