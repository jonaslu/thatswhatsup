package main

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type covariantFiles struct {
	fileA string
	fileB string
}

type covariantFilesCount struct {
	covariantFiles
	count int
}

type covariantFilesResult []covariantFilesCount

func (s covariantFilesResult) Len() int {
	return len(s)
}
func (s covariantFilesResult) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s covariantFilesResult) Less(i, j int) bool {
	return s[i].count < s[j].count
}

func copyMap(source map[string][]string) map[string][]string {
	dest := map[string][]string{}
	for srcKey, srcValues := range source {
		destShas := make([]string, len(srcValues))
		copy(destShas, srcValues)
		dest[srcKey] = destShas
	}

	return dest
}

func main() {
	var err error
	database, err := sql.Open("postgres", "postgres://postgres@localhost/code-mats?sslmode=disable")
	if err != nil {
		panic(err)
	}

	rows, err := database.Query("select filename, array_agg(sha) from log group by filename having count(*) > 1;")
	if err != nil {
		panic(err)
	}

	fileNameSha := map[string][]string{}

	var fileNameStr string
	var shasStr []string
	for rows.Next() {
		err := rows.Scan(&fileNameStr, pq.Array(&shasStr))
		if err != nil {
			panic(err)
		}

		shas := fileNameSha[fileNameStr]
		if shas == nil {
			shas = []string{}
		}

		shas = append(shas, shasStr...)
		fileNameSha[fileNameStr] = shas
	}

	// Skip deleted files
	// Skip ignored files
	// Try it on mine first

	covariantFilesCountAsMap := map[covariantFiles]int{}

	searchShaForFiles := copyMap(fileNameSha)

	for fileNameA, shasA := range fileNameSha {
		delete(searchShaForFiles, fileNameA)
		for _, shaA := range shasA {
			for fileNameB, shasB := range searchShaForFiles {
				for _, shaB := range shasB {
					if shaA == shaB {
						covariantFiles := covariantFiles{
							fileA: fileNameA,
							fileB: fileNameB,
						}

						covariantFilesCountAsMap[covariantFiles]++
					}
				}
			}
		}
	}

	covariantFilesRes := covariantFilesResult{}

	for covariantFile, count := range covariantFilesCountAsMap {
		if count <= 1 {
			continue
		}

		covariantFilesRes = append(covariantFilesRes, covariantFilesCount{
			covariantFiles: covariantFiles{
				fileA: covariantFile.fileA,
				fileB: covariantFile.fileB,
			},
			count: count,
		})
	}

	sort.Sort(covariantFilesRes)

	fmt.Println(covariantFilesRes)
}
