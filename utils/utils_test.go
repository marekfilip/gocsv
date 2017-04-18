package utils

import (
	"testing"
)

func BenchmarkGetPluralWordForRows(b *testing.B) {
	b.Run("C=1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			GetPluralWordForRows(1)
		}
	})
	b.Run("C=2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			GetPluralWordForRows(2)
		}
	})
	b.Run("C=10", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			GetPluralWordForRows(10)
		}
	})
}

func TestGetPluralWordForRows(t *testing.T) {
	testCases := map[int]string{
		0:    "wierszy",
		1:    "wiersz",
		2:    "wiersze",
		3:    "wiersze",
		4:    "wiersze",
		5:    "wierszy",
		1000: "wierszy",
	}

	for count, expected := range testCases {
		actual := GetPluralWordForRows(count)
		if actual != expected {
			t.Errorf("For '%d' expected '%s', got '%s'", count, expected, actual)
		}
	}
}
