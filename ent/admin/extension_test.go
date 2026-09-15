package admin

import "testing"

func TestFieldName(t *testing.T) {
	tests := map[string]string{
		"user_id": "UserID",
		"ip_hash": "IPHash",
		"url":     "URL",
	}

	for input, want := range tests {
		if got := fieldName(input); got != want {
			t.Fatalf("fieldName(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestFieldLabel(t *testing.T) {
	tests := map[string]string{
		"user_id": "User ID",
		"ip_hash": "IP Hash",
		"url":     "URL",
	}

	for input, want := range tests {
		if got := FieldLabel(input); got != want {
			t.Fatalf("FieldLabel(%q) = %q, want %q", input, got, want)
		}
	}
}
