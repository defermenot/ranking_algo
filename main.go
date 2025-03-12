package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
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
	repositories := make(map[string]Repository)
	oldestTimestamp := time.Now().Unix()
	newestTimestamp := int64(0)
	for i, record := range records {
		if i == 0 {
			continue
		}
		commitTimestamp, err := strconv.Atoi(record[0])
		if err != nil {
			log.Fatal(err)
		}

		if int64(commitTimestamp) < oldestTimestamp {
			oldestTimestamp = int64(commitTimestamp)
		}
		if int64(commitTimestamp) > newestTimestamp {
			newestTimestamp = int64(commitTimestamp)
		}

		name := record[2]
		files, err := strconv.Atoi(record[3])
		if err != nil {
			log.Fatal(err)
		}
		additions, err := strconv.Atoi(record[4])
		if err != nil {
			log.Fatal(err)
		}
		deletions, err := strconv.Atoi(record[5])
		if err != nil {
			log.Fatal(err)
		}
		totalLinesChanged := additions + deletions
		repo, found := repositories[name]
		if !found {
			repo = Repository{Name: name}
		}
		repo.Files += files
		repo.Additions += additions
		repo.Deletions += deletions
		repo.Total += totalLinesChanged
		repositories[name] = repo
	}
	var repos Repositories
	for _, repo := range repositories {
		repos = append(repos, repo)
	}
	timespan := newestTimestamp - oldestTimestamp
	log.Printf("oldest %d newest %d timespan %d", oldestTimestamp, newestTimestamp, timespan)
	sort.Sort(repos)
	log.Println("Repositories processed:", len(repositories))
	//log.Printf("repositories: %+v", repos)
}

type Repository struct {
	Name      string
	Additions int
	Deletions int
	Total     int
	Files     int
}

type Repositories []Repository

func (a Repositories) Len() int           { return len(a) }
func (a Repositories) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Repositories) Less(i, j int) bool { return a[i].Total > a[j].Total }

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
