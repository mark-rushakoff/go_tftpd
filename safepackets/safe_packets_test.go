package safepackets

import (
	"testing"
)

func TestSafeDataEquality(t *testing.T) {
	data1 := NewSafeData(1, []byte("Hello"))
	data2 := NewSafeData(1, []byte("Hello"))

	if !data1.Equals(data2) {
		t.Fatalf("Data should have been equal but were not")
	}
}

func TestSafeDataInequality(t *testing.T) {
	data := NewSafeData(1, []byte("Hello"))

	other := NewSafeData(2, []byte("Hello"))
	if data.Equals(other) {
		t.Errorf("Expected inequality when block numbers do not match: %v, %v", other, data)
	}

	other = NewSafeData(1, []byte("Hello!"))
	if data.Equals(other) {
		t.Errorf("Expected inequality when bytes do not match: %v, %v", other, data)
	}
}
