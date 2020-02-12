package info

import (
	"sync"
    "fmt"
    // "os"
)

type InfoInterface interface {
	IncrementMalformedLines()
	IncrementConsumedLines()
	IncrementOrphanLines()
    IncrementDuplicateLines()
}

// Info struct holds  log processing statistics.
// TODO: There can be orphan lines (i.e., services with no corresponding root service); they should
// be tolerated but not included in the output (maybe summarized in stats).
// Let's assume one logEntry 80 bytes, https://golang.org/pkg/builtin/#int says
// uint64 is 18446744073709551615,  uint32 4294967295
// uint64 => 1638.4 petabytes, uint32 3.9 gigabytes
type Info struct {
	m               sync.Mutex
	malformedLines  uint64   // malformed lines
	consumedLines   uint64   // consumed lines
	orphanLines     uint64   // services with no corresponding root service
    duplicateLines  uint64   // we got multiple root spans or same lines for one trace
}

// IncrementMalformedLines increments malformed lines count by 1.
func (i *Info) IncrementMalformedLines() {
	i.m.Lock()
    // fmt.Fprintln(os.Stderr,"[DEBUG] malformedLine ")
	i.malformedLines++
	i.m.Unlock()
}

// IncrementConsumedLines increments consumed lines count by 1.
func (i *Info) IncrementConsumedLines() {
    // fmt.Fprintln(os.Stderr,"[DEBUG] consumedLines ")
    i.m.Lock()
	i.consumedLines++
    i.m.Unlock()
}

// IncrementOrphanLines increments orphan lines count by 1.
func (i *Info) IncrementOrphanLines() {
    i.m.Lock()
	i.orphanLines++
    i.m.Unlock()
}


// IncrementDuplicateLines increments duplicate lines count by 1.
func (i *Info) IncrementDuplicateLines() {
    i.m.Lock()
	i.duplicateLines++
    i.m.Unlock()
}

func (i *Info) String() string {
 return fmt.Sprintf("consumedLines  %v\nmalformedLines %v\norphanLines    %v\nduplicateLines %v\n" , i.consumedLines,  i.malformedLines, i.orphanLines, i.duplicateLines)
}

// NewInfo creates a new Info object.
func NewInfo() *Info {
	return &Info{m: sync.Mutex{}}
}
