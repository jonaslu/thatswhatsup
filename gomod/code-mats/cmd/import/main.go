package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"

	_ "github.com/lib/pq"
)

var database *sql.DB

// git log --all --numstat --date=short --pretty=format:'%h--%ad--%aN' --no-renames --after=1990-01-01 > ~/code/experiments/code-mats/log.txt

// # Add\tremoved\tfile
// # Added\tremoved file
// --187aacb--2020-08-12--Jonas Lundberg
// 11	7	package.json

type insertData struct {
	sha         string
	date        string
	author      string
	addedRows   string
	removedRows string
	fileName    string
}

var workers sync.WaitGroup
var workQueue chan insertData

func isCommitRow(entry string) bool {
	res, _ := regexp.MatchString("^([0-9a-f]){7,40}--", entry)
	return res
}

func worker() {
	for data := range workQueue {

		insertIntoLogStmt := `INSERT INTO log (
			sha,
			date,
			author,
			added,
			removed,
			filename
		) VALUES (
			$1,$2, $3, $4, $5, $6
		);`

		if _, err := database.Exec(insertIntoLogStmt, data.sha, data.date, data.author, data.addedRows, data.removedRows, data.fileName); err != nil {
			panic(err)
		}
	}

	workers.Done()
}

func processCommitEntry(entry string) {
	allLines := strings.Split(entry, "\n")

	var shaDateAuthorStr string
	var commitRows []string

	lastCommitRowIndex := 0

	for ; isCommitRow(allLines[lastCommitRowIndex]); lastCommitRowIndex++ {
	}

	shaDateAuthorStr = allLines[lastCommitRowIndex-1]
	commitRows = allLines[lastCommitRowIndex:]

	shaDateAuthor := strings.Split(shaDateAuthorStr, "--")

	sha := shaDateAuthor[0]
	date := shaDateAuthor[1]
	author := shaDateAuthor[2]

	for _, commitRow := range commitRows {
		fileRowContents := strings.Split(commitRow, "\t")

		addedRows := fileRowContents[0]
		removedRows := fileRowContents[1]
		fileName := fileRowContents[2]

		// Binary files are stat:ed as - - (0 added 0 removed)
		if removedRows == "-" {
			removedRows = "0"
		}

		if addedRows == "-" {
			addedRows = "0"
		}

		if removedRows == "" || fileName == "" {
			panic(commitRow)
		}

		workQueue <- insertData{
			sha:         sha,
			date:        date,
			author:      author,
			addedRows:   addedRows,
			removedRows: removedRows,
			fileName:    fileName,
		}
	}
}

func main() {
	var err error
	database, err = sql.Open("postgres", "postgres://postgres@localhost/code-mats?sslmode=disable")
	if err != nil {
		panic(err)
	}

	filename := os.Args[1]

	logContentsBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	logContentsStr := string(logContentsBytes)
	commitEntries := strings.Split(logContentsStr, "\n\n")

	progressIndicator := false
	totalEntries := len(commitEntries)

	if totalEntries > 10000 {
		progressIndicator = true
	}

	entriesProcessed := 0

	numberOfWorkers := runtime.NumCPU()
	workQueue = make(chan insertData, numberOfWorkers)
	workers = sync.WaitGroup{}

	for i := 0; i < numberOfWorkers; i++ {
		go worker()
		workers.Add(1)
	}

	for _, commitEntry := range commitEntries {
		commitEntry = strings.TrimSpace(commitEntry)
		if commitEntry == "" {
			continue
		}

		processCommitEntry(commitEntry)

		if progressIndicator {
			entriesProcessed++

			if entriesProcessed%1000 == 0 {
				fmt.Printf("Proccessed: %d of %d entries\n", entriesProcessed, totalEntries)
			}
		}
	}

	close(workQueue)
	workers.Wait()
}
