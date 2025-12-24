package types

import (
	"encoding/json"
	"testing"
)

func TestFlexInt64(t *testing.T) {
	type wrapper struct {
		ID FlexInt64 `json:"id"`
	}

	// Test Unmarshal String
	jsonStr1 := `{"id": "1234567890123456789"}`
	var w1 wrapper
	if err := json.Unmarshal([]byte(jsonStr1), &w1); err != nil {
		t.Fatalf("Unmarshal string failed: %v", err)
	}
	if int64(w1.ID) != 1234567890123456789 {
		t.Errorf("Expected 1234567890123456789, got %d", w1.ID)
	}

	// Test Unmarshal Number
	jsonStr2 := `{"id": 123}`
	var w2 wrapper
	if err := json.Unmarshal([]byte(jsonStr2), &w2); err != nil {
		t.Fatalf("Unmarshal number failed: %v", err)
	}
	if int64(w2.ID) != 123 {
		t.Errorf("Expected 123, got %d", w2.ID)
	}

	// Test Marshal (Should be string)
	w3 := wrapper{ID: FlexInt64(987654321098765432)}
	b, err := json.Marshal(w3)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	expected := `{"id":"987654321098765432"}`
	if string(b) != expected {
		t.Errorf("Expected %s, got %s", expected, string(b))
	}
}
