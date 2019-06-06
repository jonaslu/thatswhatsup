package compilerutils

import (
	"fmt"
	"sync"
)

var (
	jmpLabelCounter uint64 = 0
	jmpLabelSync    sync.Mutex
)

// GetUniqueLabel returns a (for the program) unique label for jmps
// Is safe for concurrency
func GetUniqueLabel() string {
	jmpLabelSync.Lock()
	returnValue := fmt.Sprintf("a%d", jmpLabelCounter)
	jmpLabelCounter++
	jmpLabelSync.Unlock()
	return returnValue
}
