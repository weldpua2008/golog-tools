package parser

import (
	"testing"
	"time"
    "strings"
    "sync"
    "container/heap"
)

func TestLogEntry(t *testing.T) {
    testRootSpan := &LogEntry{
        Service:  "back-end-1",
        Trace:    "trace1",
        CallerSpan: "null",
        Span:   "aa",
        Start:    time.Date(2020, 2, 29, 2, 0, 0, 0, time.UTC),
        End:      time.Date(2021, 12, 1, 22, 13, 5, 0, time.UTC),
    }
    equalDates:=Logs{
        &LogEntry{
        	Service:  "back-end-1",
        	Trace:    "trace1",
        	CallerSpan: "aa",
        	Span:   "ac",
        	Start:    time.Date(2020, 2, 29, 2, 0, 0, 0, time.UTC),
        	End:      time.Date(2021, 12, 1, 22, 13, 5, 0, time.UTC),
        },
        &LogEntry{
        	Service:  "back-end-2",
        	Trace:    "trace1",
        	CallerSpan: "aa",
        	Span:   "ab",
        	Start:    time.Date(2020, 2, 29, 2, 0, 0, 0, time.UTC),
        	End:      time.Date(2021, 12, 1, 22, 13, 5, 0, time.UTC),
        },
        testRootSpan,
    }
    lastDateEqual:= &LogEntry{
        Service:  "back-end-3",
        Trace:    "trace1",
        CallerSpan: "ac",
        Span:   "ad",
        Start:    time.Date(2020, 2, 29, 1, 59, 0, 0, time.UTC),
        End:      time.Date(2021, 12, 1, 22, 13, 5, 0, time.UTC),
    }
    tests := []struct {
        Name       string
		logs       Logs
        lastEntry  *LogEntry
        m          sync.Mutex

	}{
		{
            Name:       "Equal dates",
			logs: append(equalDates,lastDateEqual),
            lastEntry: lastDateEqual,
            m: sync.Mutex{},
		},
        {
            Name:       "Two spans",
			logs:    Logs{testRootSpan,lastDateEqual},
            lastEntry: lastDateEqual,
            m: sync.Mutex{},
		},
    }

	for _, tc := range tests {

        // tc := tc
		t.Run(tc.Name, func(t *testing.T) {
        // t.Parallel()
        tc.m.Lock()
        {
            heap.Init(&tc.logs)
            copiedLogs :=append(tc.logs[:0:0], tc.logs...)
            var rs *LogEntry
            debugValues := []string{}
            debugValues = append(debugValues, tc.logs.String())
            rs =heap.Pop(&tc.logs).(*LogEntry)
            if rs.CallerSpan != RootSpanConst {
                t.Errorf("CallerSpan %s (%q) != (%q) ", &copiedLogs, rs.CallerSpan, RootSpanConst)
            }


            for tc.logs.Len() > 0 {
                rs =heap.Pop(&tc.logs).(*LogEntry)
        	}
            if rs.Span != tc.lastEntry.Span {
                t.Errorf("CallerSpan order %s during calls %s (%q) != (%q) %q != %q", &copiedLogs, strings.Join(debugValues, "|\n "), rs.Span, tc.lastEntry.Span, rs, tc.lastEntry)
            }
		}
        })
	}

}
