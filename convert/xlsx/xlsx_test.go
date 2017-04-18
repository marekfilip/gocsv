package xlsx

import (
	"testing"
)

func BenchmarkConvertFileXlsx(b *testing.B) {
	for n := 0; n < b.N; n++ {
		convertFileXlsx("../../test_files/test.xlsx", ';', 0)
	}
}
