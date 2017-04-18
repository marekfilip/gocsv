package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"filip/gocsv/convert/xls"
	"filip/gocsv/convert/xlsx"
)

const (
	ThreadsMultipler = 2
	ThreadsPerWorker = 2
)

var (
	delimiter string
	skip      int
	filenames []string
)

func init() {
	flag.StringVar(&delimiter, "d", ";", "Znak do rozdzielenia kolumn w CSV")
	flag.IntVar(&skip, "s", 0, "Liczba wierszy do pominięcia")
	flag.Parse()

	filenames = flag.Args()

	if len(filenames) == 0 {
		log.Println("Nie podano plików do przetworzenia")
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func main() {
	s0 := time.Now()
	runtime.GOMAXPROCS(runtime.NumCPU())
	var threadCount = runtime.GOMAXPROCS(-1) * ThreadsMultipler / ThreadsPerWorker
	var xlsxCount int
	var xlsCount int

	for _, file := range filenames {
		switch strings.ToLower(filepath.Ext(file)) {
		case ".xlsx":
			xlsxCount++
			break
		case ".xls":
			xlsCount++
			break
		}
	}

	doneResults := make(chan bool, 100)
	xlsxFileChan := make(chan string, 100)
	xlsFileChan := make(chan string, 100)

	allFilesCount := xlsxCount + xlsCount
	if allFilesCount == 0 {
		log.Println("Brak plików do przetworzenia. Koniec")
		os.Exit(0)
	}
	maxXlsxThreads := int(xlsxCount / (allFilesCount) * threadCount)
	maxXlsThreads := int(xlsCount / (allFilesCount) * threadCount)
	log.Printf("Tworzenie %d wątków dla plików xlsx\n", maxXlsxThreads)
	log.Printf("Tworzenie %d wątków dla plików xls\n", maxXlsThreads)

	for i := 0; i < maxXlsxThreads; i++ {
		go xlsx.NewXlsxConvertWorker(i, xlsxFileChan, doneResults, rune(delimiter[0]), skip)
	}
	for i := 0; i < maxXlsThreads; i++ {
		go xls.NewXlsConvertWorker(i, xlsFileChan, doneResults, rune(delimiter[0]), skip)
	}

	for _, file := range filenames {
		switch strings.ToLower(filepath.Ext(file)) {
		case ".xlsx":
			xlsxFileChan <- file
			break
		case ".xls":
			xlsFileChan <- file
			break
		}
	}

	close(xlsxFileChan)
	close(xlsFileChan)

	for i := 0; i < threadCount; i++ {
		<-doneResults
	}

	log.Println("Prace trwały", time.Now().Sub(s0))
	os.Exit(0)
}
