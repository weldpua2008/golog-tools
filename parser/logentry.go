package parser

import (
	"time"
    "fmt"
    "strings"
)


// logs collection  is a main-heap of LogEntry
type Logs []*LogEntry

// log file data entry
type LogEntry struct {
	Service    string        // Service-name
	Trace      string        // trace ID
	CallerSpan string        // called from outside
	Span       string        // called services
	Start      time.Time     // start timestamp UTC
	End        time.Time     // end timestamp UTC
}

// Len returns the len for the Logs collection.
func (l Logs) Len() int {
	return len(l)
}

// The entries are logged when the request finishes
// (as they contain the finishing time), so they
// are not in calling order, but in finishing order.
// Less return true if l[i] is before l[j].
func (l Logs) Less(i, j int) bool {
    switch {
        case l[i].CallerSpan == RootSpanConst ||  l[i].Span == l[j].CallerSpan:
            return true
        case l[j].CallerSpan == RootSpanConst ||  l[j].Span == l[i].CallerSpan:
            return false
        default:
            // naive implementation if ordered by date
            return l[i].Start.Before(l[j].Start)
    }
}

// Swap swaps l[i] and l[j].
func (l Logs) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// container/heap
func (l *Logs) Push(x interface{}) {
	// Push use pointer receivers because they modify the slice's length,
	// not just its contents.
    z := x.(LogEntry)
    // if l.Len() > 1 {
    //     for _, value := range *l  {
    // 		if value == x {
    //             *l.dupl++
    //         }
    // 	}
    //
    // }
	*l = append(*l, &z)
}

func (l *Logs) Pop() interface{} {
    // Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	old := *l
	n := len(old)
	x := old[n-1]
	*l = old[0 : n-1]

	return x
}

// HasRootCallerSpan returns true if there is root caller span
func (l *Logs) HasRootCallerSpan() bool {
	for _, value := range *l  {
		if value.CallerSpan == RootSpanConst {
            return true
        }
	}
    return false
}

// // GetDuplicateSpans returns number of duplicate spans
// func (l *Logs) GetDuplicateSpans() int {
//    duplicate_frequency := make(map[string]int)
//
//    for _, v := range *l {
//        item:=fmt.Sprintf("%v->%v", value.CallerSpan,value.Span)
//        // check if the item/element exist in the duplicate_frequency map
//        _, exist := duplicate_frequency[item]
//        if exist {
//            duplicate_frequency[item] += 1 // increase counter by 1 if already in the map
//        } else {
//            duplicate_frequency[item] = 1 // else start counting from 1
//        }
//    }
//    return len(duplicate_frequency)
// }

// String returns a string representation of container
func (l *Logs) String() string {
	str := ""
	values := []string{}
	for _, value := range *l  {
		values = append(values, fmt.Sprintf("%v->%v", value.CallerSpan,value.Span))
	}
	str += strings.Join(values, ", ")
	return str
}
