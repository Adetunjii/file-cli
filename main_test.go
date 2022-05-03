package main

import (
	"flag"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGetFileData(t *testing.T) {
	tests := []struct {
		name     string
		expected inputFile
		err      bool
		osArgs   []string
	}{
		{"Default parameters", inputFile{"test.csv", "comma", false}, false, []string{"cmd", "test.csv"}},
		{"No parameters", inputFile{}, true, []string{"cmd"}},
		{"Semicolon enabled", inputFile{"test.csv", "semicolon", false}, false, []string{"cmd", "--separator=semicolon", "test.csv"}},
		{"Pretty enabled", inputFile{"test.csv", "comma", true}, false, []string{"cmd", "--pretty", "test.csv"}},
		{"Pretty and semicolon enabled", inputFile{"test.csv", "semicolon", true}, false, []string{"cmd", "--pretty", "--separator=semicolon", "test.csv"}},
		{"Separator not identified", inputFile{}, true, []string{"cmd", "--separator=pipe", "test.csv"}},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			actualOsArgs := os.Args

			defer func() {
				os.Args = actualOsArgs
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			}()

			os.Args = testCase.osArgs
			inputFile, err := getFileData()

			if (err != nil) != testCase.err {
				t.Errorf("getFileData() error = %v, wantErr %v", err, testCase.err)
				return
			}

			require.Equal(t, testCase.expected, inputFile)

		})
	}
}

func TestCheckIfValidFile(t *testing.T) {

	tempFile, err := os.CreateTemp("", "test*.csv")

	if err != nil {
		panic(err)
	}

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {

		}
	}(tempFile.Name())

	tests := []struct {
		name     string
		filename string
		expected bool
		err      bool
	}{
		{"File does exist", tempFile.Name(), true, false},
		{"File does not exist", "nowhere/test.csv", false, true},
		{"File is not csv", "test.txt", false, true},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			isValid, err := checkIfValidFile(testCase.filename)

			if (err != nil) != testCase.err {
				t.Errorf("File %v is invalid", testCase.filename)
				return
			}

			require.Equal(t, testCase.expected, isValid)
		})
	}
}

func TestProcessCsvFile(t *testing.T) {
	expectedMapSlice := []map[string]string{
		{"id": "1", "name": "samuel", "age": "25", "email": "adetunjithomas1@gmail.com"},
	}

	tests := []struct {
		name      string
		csvString string
		separator string
	}{
		{"Comma separator", "id,name,age,email\n1,samuel,25,adetunjithomas1@gmail.com", "comma"},
		{"Semicolon separator", "id;name;age;email\n1;samuel;25;adetunjithomas1@gmail.com", "semicolon"},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "test*.csv")
			defer func(name string) {
				err := os.Remove(name)
				require.NoError(t, err)
			}(tmpFile.Name())

			require.NoError(t, err)

			_, err = tmpFile.WriteString(testCase.csvString)

			// moves the tmpFile from memory to disk
			err = tmpFile.Sync()
			require.NoError(t, err)

			testFileData := inputFile{
				filePath:  tmpFile.Name(),
				pretty:    false,
				separator: testCase.separator,
			}

			writerChannel := make(chan map[string]string)

			go processCsvFile(testFileData, writerChannel)

			for _, expectedMap := range expectedMapSlice {
				record := <-writerChannel
				require.Equal(t, expectedMap, record)
			}
		})
	}

}
