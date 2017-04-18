package xls

import (
	"fmt"
	"log"
	"time"

	u "filip/gocsv/utils"

	"github.com/extrame/xls"
)

func NewXlsConvertWorker(threadNumber int, jobs <-chan string, done chan<- bool, delimiter rune, skipLines int) {
	for file := range jobs {
		rows, dur, err := convertFileXls(file, delimiter, skipLines)
		if err != nil {
			u.HandleError(err)
		} else {
			log.Printf("[T: %d] '%s': %d %s w %s\n", threadNumber, file, rows, u.GetPluralWordForRows(rows), dur)
		}
	}

	done <- true
}

func convertFileXls(file string, delimiter rune, skipLines int) (int, time.Duration, error) {
	var rowsCount int
	var allCells [][]string
	startTime := time.Now()

	xlsFile, err := xls.Open(file, "utf-8")
	u.HandleError(err)

	allCells = xlsFile.ReadAllCells(100000)
	allRowsCount := len(allCells)

	if allRowsCount > 0 {
		stringSliceChanell := make(chan []string, 100)
		done := make(chan bool)
		go u.OpenCsvWriter(file, delimiter, stringSliceChanell, done)

		rowsCount = allRowsCount - skipLines
		for rowNum, row := range allCells {
			if rowNum >= skipLines {
				stringSliceChanell <- row
			}
		}
		close(stringSliceChanell)

		<-done
	} else {
		return rowsCount, time.Now().Sub(startTime), fmt.Errorf("'%s' niemogł być przetworzony", file)
	}

	return rowsCount, time.Now().Sub(startTime), nil
}
