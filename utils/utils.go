package utils

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GetPluralWordForRows(rowsCount int) string {
	switch {
	case rowsCount == 1:
		return "wiersz"
	case 1 < rowsCount && rowsCount < 5:
		return "wiersze"
	default:
		return "wierszy"
	}
}

func HandleError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func OpenCsvWriter(file string, delimiter rune, lines <-chan []string, done chan<- bool) {
	fileHandle, err := os.Create(
		strings.Join(
			[]string{
				file[0 : len(file)-len(filepath.Ext(file))],
				"csv",
			},
			".",
		),
	)
	HandleError(err)

	writer := csv.NewWriter(fileHandle)
	writer.Comma = rune(delimiter)

	for word := range lines {
		writer.Write(word)
	}

	done <- true
	writer.Flush()
	HandleError(writer.Error())
	fileHandle.Close()

	return
}
