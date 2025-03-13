package main

import (
	"log"
	"math"
	"strconv"
	"time"
)

const (
	defaultWeight     = 0.3
	defaultDaysPeriod = 30
)

type Commit struct {
	RepoName          string
	User              string
	Files             int
	Additions         int
	Deletions         int
	Timestamp         int
	TotalLinesChanged int
	Weight            float64
}

type Repository struct {
	Name              string
	Additions         int
	Deletions         int
	TotalLinesChanged int
	Files             int
	Score             float64
}

type Repositories map[string]*Repository

type TimeDecay struct {
	DecayRate    float64
	Oldest       int
	Recent       int
	Commits      []Commit
	Repositories Repositories
}

func NewTimeDecay() TimeDecay {
	t := TimeDecay{
		DecayRate:    calculateDecayRate(defaultWeight, defaultDaysPeriod),
		Repositories: make(Repositories),
		Commits:      make([]Commit, 0),
	}
	return t
}

func (t *TimeDecay) parseRecords(records [][]string) {
	older := int(time.Now().Unix())
	newer := 0
	commits := make([]Commit, len(records)-1)
	commitIndex := 0
	for i, record := range records {
		// ignore csv header
		if i == 0 {
			continue
		}
		commit, err := parseCommit(record)
		if err != nil {
			log.Fatal(err)
		}
		commits[commitIndex] = commit
		commitIndex++
		older = min(commit.Timestamp, older)
		newer = max(commit.Timestamp, newer)
		repo, found := t.Repositories[commit.RepoName]
		if !found {
			repo = repositoryFromCommit(commit)
			t.Repositories[commit.RepoName] = repo
		}
		repo.AddCommit(commit)
	}
	t.Oldest = older
	t.Recent = newer
	t.Commits = commits

	return
}

func (t *TimeDecay) Rank(records [][]string) []Repository {
	t.parseRecords(records)
	log.Printf("Time range: oldest=%d, newest=%d (diff: %d seconds)",
		t.Oldest, t.Recent, t.Recent-t.Oldest)

	// Debug the decay rate
	log.Printf("Using decay rate: %.10f", t.DecayRate)
	for _, commit := range t.Commits {
		age := t.Recent - commit.Timestamp
		weight := calculateCommitWeight(float64(age), t.DecayRate)
		t.Repositories[commit.RepoName].Score += float64(commit.TotalLinesChanged) * weight
	}
	return t.Repositories.ToSlice()
}

func (r *Repository) AddCommit(commit Commit) {
	r.Files += commit.Files
	r.Additions += commit.Additions
	r.Deletions += commit.Deletions
	r.TotalLinesChanged += commit.TotalLinesChanged
}

func (r Repositories) ToSlice() []Repository {
	var slice []Repository
	for _, repo := range r {
		slice = append(slice, *repo)
	}
	return slice
}

func calculateDecayRate(targetWeight float64, daysPeriod int) float64 {
	periodInSeconds := float64(daysPeriod * 24 * 60 * 60)
	decayRate := -math.Log(targetWeight) / periodInSeconds
	return decayRate
}

func calculateCommitWeight(age, decayRate float64) float64 {
	return math.Exp(-decayRate * age)
}

func min(target, current int) int {
	if target < current {
		current = target
	}
	return current
}

func max(target, current int) int {
	if target > current {
		current = target
	}
	return current
}

func parseCommit(record []string) (Commit, error) {
	commitTimestamp, err := strconv.Atoi(record[0])
	if err != nil {
		log.Fatal(err)
	}
	user := record[1]
	name := record[2]
	files, err := strconv.Atoi(record[3])
	if err != nil {
		return Commit{}, err
	}
	additions, err := strconv.Atoi(record[4])
	if err != nil {
		return Commit{}, err
	}
	deletions, err := strconv.Atoi(record[5])
	if err != nil {
		return Commit{}, err
	}
	return Commit{
		RepoName:          name,
		User:              user,
		Files:             files,
		Additions:         additions,
		Deletions:         deletions,
		Timestamp:         commitTimestamp,
		TotalLinesChanged: additions + deletions,
	}, nil
}

func repositoryFromCommit(commit Commit) *Repository {
	return &Repository{
		Name:              commit.RepoName,
		Additions:         commit.Additions,
		Deletions:         commit.Deletions,
		TotalLinesChanged: commit.TotalLinesChanged,
		Files:             commit.Files,
	}
}
