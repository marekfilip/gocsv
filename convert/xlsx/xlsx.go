package xlsx

import (
	"fmt"
	"log"
	"time"

	u "filip/gocsv/utils"

	"github.com/tealeg/xlsx"
)

func NewXlsxConvertWorker(threadNumber int, jobs <-chan string, done chan<- bool, delimiter rune, skipLines int) {
	for file := range jobs {
		rows, dur, err := convertFileXlsx(file, delimiter, skipLines)
		if err != nil {
			u.HandleError(err)
		} else {
			log.Printf("[T: %d] '%s': %d %s w %s\n", threadNumber, file, rows, u.GetPluralWordForRows(rows), dur)
		}
	}

	done <- true
}

func convertFileXlsx(file string, delimiter rune, skipLines int) (int, time.Duration, error) {
	var rowsCount int
	var startTime time.Time = time.Now()

	readFile, err := xlsx.OpenFile(file)
	u.HandleError(err)

	stringSliceChanel := make(chan []string, 100)
	done := make(chan bool)
	go u.OpenCsvWriter(file, delimiter, stringSliceChanel, done)

	for _, sheet := range readFile.Sheets {
		rowsCount += len(sheet.Rows) - skipLines
		for rowNum, row := range sheet.Rows {
			if rowNum >= skipLines {
				csvRow := make([]string, 0)
				for _, cell := range row.Cells {
					str, err := cell.String()
					if err == nil {
						csvRow = append(csvRow, str)
						continue
					}

					return rowsCount, time.Now().Sub(startTime), fmt.Errorf("Error (%s): %s", file, err.Error())
				}
				stringSliceChanel <- csvRow
			}
		}
	}
	close(stringSliceChanel)

	<-done
	return rowsCount, time.Now().Sub(startTime), nil
}
