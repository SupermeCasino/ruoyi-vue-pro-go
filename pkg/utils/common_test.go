package utils

import (
	"testing"
)

func TestParseInt64(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int64
	}{
		{"Valid number", "123", 123},
		{"Zero", "0", 0},
		{"Negative", "-1", -1},
		{"Invalid", "abc", 0},
		{"Empty", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseInt64(tt.input)
			if got != tt.want {
				t.Errorf("ParseInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateRandomString(t *testing.T) {
	length := 10
	got := GenerateRandomString(length)
	if len(got) != length {
		t.Errorf("GenerateRandomString() length = %v, want %v", len(got), length)
	}
}

func TestIntSliceContains(t *testing.T) {
	slice := []int{1, 2, 3}
	if !IntSliceContains(slice, 2) {
		t.Errorf("IntSliceContains() should return true for existing element")
	}
	if IntSliceContains(slice, 4) {
		t.Errorf("IntSliceContains() should return false for non-existing element")
	}
}
