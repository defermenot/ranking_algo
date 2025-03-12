package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	inputFlag := flag.String("input", "", "input file csv path")
	outputFlag := flag.String("output", "", "output file csv or json path")
	flag.Parse()

	inputPath, err := filepath.Abs(*inputFlag)
	if err != nil {
		log.Fatal(err)
	}

	if err := validateInputPath(inputPath); err != nil {
		log.Fatal(err)
	}

	outputPath, err := filepath.Abs(*outputFlag)
	if err != nil {
		log.Fatal(err)
	}

	if err := validateOutputPath(outputPath); err != nil {
		log.Fatal(err)
	}
	records, err := readCsv(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(records)
}

func validateInputPath(inputPath string) error {
	if !strings.HasSuffix(inputPath, ".csv") && inputPath == "" {
		return fmt.Errorf("input file must be a csv file")
	}
	return nil
}

func validateOutputPath(outputPath string) error {
	if !strings.HasSuffix(outputPath, ".csv") && !strings.HasSuffix(outputPath, ".json") {
		return fmt.Errorf("output file must be a csv or json file")
	}
	return nil
}

func readCsv(path string) ([][]string, error) {
	csvFile, err := os.Open(path)
	if err != nil {
		return nil, nil
	}
	defer csvFile.Close()
	r := csv.NewReader(csvFile)

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}
