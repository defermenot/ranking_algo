package main

import (
	"log"
	"strconv"
	"time"
)

type Ranker interface {
	Rank(records []Repository) []Repository
}

type TimeDecay struct {
	Weight    float64
	DecayRate float64
	MaxAge    time.Duration
	MinAge    time.Duration
}

func (t TimeDecay) parseRecords(records [][]string) {
	repositories := make(map[string]Repository)
	for i, record := range records {
		// ignore csv header
		if i == 0 {
			continue
		}
		older := int(time.Now().Unix())
		newer := 0

		commit, err := parseCommit(record)
		if err != nil {
			log.Fatal(err)
		}
		older = min(commit.Timestamp, older)
		newer = max(commit.Timestamp, newer)

		repo, found := repositories[commit.RepoName]
		if !found {
			repo = repositoryFromCommit(commit)
		}
		repo.AddCommit(commit)
	}
}

func (t TimeDecay) Rank(repositories []Repository) []Repository {

	return nil
}

func min(commitTimestamp, current int) int {
	if commitTimestamp < current {
		current = commitTimestamp
	}
	return current
}

func max(commitTimestamp, current int) int {
	if commitTimestamp > current {
		current = commitTimestamp
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

type Commit struct {
	RepoName          string
	User              string
	Files             int
	Additions         int
	Deletions         int
	Timestamp         int
	TotalLinesChanged int
}

type Repository struct {
	Name              string
	Additions         int
	Deletions         int
	TotalLinesChanged int
	Files             int
}

func repositoryFromCommit(commit Commit) Repository {
	return Repository{
		Name:              commit.RepoName,
		Additions:         commit.Additions,
		Deletions:         commit.Deletions,
		TotalLinesChanged: commit.TotalLinesChanged,
		Files:             commit.Files,
	}
}

func (r *Repository) AddCommit(commit Commit) {
	r.Files += commit.Files
	r.Additions += commit.Additions
	r.Deletions += commit.Deletions
	r.TotalLinesChanged += commit.TotalLinesChanged
}

type Repositories map[string]*Repository
