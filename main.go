package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func main() {
	inputFlag := flag.String("input", "", "input file csv path")
	outputFlag := flag.String("output", "", "output file csv or json path")
	tailFlag := flag.String("tail", "20", "number of records to display")
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

	timeDecayRanker := NewTimeDecay()

	repos := timeDecayRanker.Rank(records)

	sort.Slice(repos, func(i, j int) bool {
		return repos[i].Score > repos[j].Score
	})
	tail, err := strconv.Atoi(*tailFlag)
	if err != nil {
		log.Fatal(err)
	}
	err = writeToJSON(outputPath, repos[:tail])
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully ranked repositories")
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

func writeToJSON(filename string, data []Repository) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling to JSON: %w", err)
	}
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}
