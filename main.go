package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type inputFile struct {
	filePath  string
	separator string
	pretty    bool
}

func main() {
	fmt.Println("Hello world")
	checkIfValidFile("test.csv")
}

func getFileData() (inputFile, error) {
	if len(os.Args) < 2 {
		return inputFile{}, errors.New("A filepath argument is required")
	}

	separator := flag.String("separator", "comma", "column separator")
	pretty := flag.Bool("pretty", false, "Generate pretty JSON")

	flag.Parse()

	fileLocation := flag.Arg(0)

	if !(*separator == "comma" || *separator == "semicolon") {
		return inputFile{}, errors.New("only comma or semicolon separators are allowed")
	}

	return inputFile{fileLocation, *separator, *pretty}, nil
}

func checkIfValidFile(filename string) (bool, error) {

	if fileExtension := filepath.Ext(filename); fileExtension != ".csv" {
		return false, fmt.Errorf("file %s is not CSV", filename)
	}

	//os.Stat checks if the entered filepath is an existing file.
	if _, err := os.Stat(filename); err != nil && os.IsNotExist(err) {
		return false, fmt.Errorf("file %s does not exist", filename)
	}

	return true, nil
}

func processCsvFile(fileData inputFile, writerChannel chan<- map[string]string) {

	// open the file
	file, err := os.Open(fileData.filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// close the file when the function returns
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error %v\n", err)
			os.Exit(1)
		}
	}()

	var headers, line []string
	reader := csv.NewReader(file)

	// the default separator for csv is comma, so we change it semicolon is specified
	if fileData.separator == "semicolon" {
		reader.Comma = ';'
	}

	// reading the file for the first time would fetch the header row
	headers, err = reader.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// keep fetching each row then send it to a writerChannel
	// the loop stops when we reach the end of the file
	for {
		line, err = reader.Read()

		if err != nil {
			if err == io.EOF {
				close(writerChannel)
				break
			}
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		record, err := processLine(headers, line)
		if err != nil {
			fmt.Printf("Line: %sError: %s\n", line, err)
			continue
		}

		writerChannel <- record
	}
}

func processLine(headers, line []string) (map[string]string, error) {
	if len(line) != len(headers) {
		return nil, errors.New("line does not match header")
	}

	recordMap := make(map[string]string)

	for i, name := range headers {
		recordMap[name] = line[i]
	}
	return recordMap, nil
}

func writeJSONFile(csvPath string, writerChannel <-chan map[string]string, done chan<- bool, pretty bool) {

	_, err := jsonFile.WriteString("[" + breakLine)

}

func createStringWriter(csvPath string) func(string, bool) {
	jsonDir := filepath.Dir(csvPath)
	jsonName := fmt.Sprintf("%s.json", strings.TrimSuffix(filepath.Base(csvPath), ".csv"))
	jsonLocation := filepath.Join(jsonDir, jsonName)

	jsonFile, err := os.Create(jsonLocation)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error occured %v ", err)
		os.Exit(1)
	}
}
