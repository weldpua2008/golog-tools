package parser

import (
	"testing"
	"time"
    "fmt"
    "strings"
)

func TestLogMessageParsing(t *testing.T) {
	tests := []struct {
        Name       string
		logMessage string
		logEntry    *LogEntry
		err         error
        err_msg     string
	}{
		{
            Name:       "Error Format - Empty line",
			logMessage: "",
			logEntry:    nil,
			err:        fmt.Errorf(""),
            err_msg:    "wrong format",
		},
        {
            Name:       "Error Format - Missing span->span",
			logMessage: "2086-0-20T12:43:34.000Z 2016-10-20T12:43:35.000Z trace1 back-end-3",
			logEntry:    nil,
			err:        fmt.Errorf(""),
            err_msg:    "wrong format",
		},
        {
            Name:       "Error StartDate - Wrong Date format",
			logMessage: "20-10-Z 2020-10-20T12:43:35.000Z trace1 back-end-3 ac->ad",
			logEntry:    nil,
			err:        fmt.Errorf(""),
            err_msg:    "start date",
		},
        {
            Name:       "Error StartDate - Not Exist Date(2019 is not leap Year)",
            logMessage: "2019-02-29T02:00:00.000Z 2020-10-20T12:43:35.000Z trace1 back-end-3 ac->ad",
            logEntry:    nil,
            err:        fmt.Errorf(""),
            err_msg:    "start date",
        },
        {
            Name:       "Error EndDate - Wrong Date format",
			logMessage: "2050-10-10T10:11:34.000Z 2016Z trace1 back-end-3 ac->ad",
			logEntry:    nil,
			err:        fmt.Errorf(""),
            err_msg:    "end date",
		},
        {
            Name:       "Error EndDate - Not Exist Date(2019 is not leap Year)",
			logMessage: "2050-10-10T10:11:34.000Z 2019-02-29T02:00:00.000Z trace1 back-end-3 ac->ad",
			logEntry:    nil,
			err:        fmt.Errorf(""),
            err_msg:    "end date",
		},
		{
            Name:       "Valid Log - Standart Year",
			logMessage: "2020-01-20T12:43:34.000Z 2020-10-20T00:00:00.000Z trace1 back-end-3 ac->ad",
			logEntry: &LogEntry{
				Service:  "back-end-3",
				Trace:    "trace1",
				CallerSpan: "ac",
				Span:   "ad",
				Start:    time.Date(2020, 1, 20, 12, 43, 34, 0, time.UTC),
				End:      time.Date(2020, 10, 20, 0, 0, 0, 0, time.UTC),
			},
			err: nil,
            err_msg: "",
		},
        {
            Name:       "Valid Log - Leap Year",
			logMessage: "2020-02-29T02:00:00.000Z 2021-12-01T22:13:05.000Z trace1 back-end-3 ac->ad",
			logEntry: &LogEntry{
				Service:  "back-end-3",
				Trace:    "trace1",
				CallerSpan: "ac",
				Span:   "ad",
				Start:    time.Date(2020, 2, 29, 2, 0, 0, 0, time.UTC),
				End:      time.Date(2021, 12, 1, 22, 13, 5, 0, time.UTC),
			},
			err: nil,
            err_msg: "",
		},

        // TODO: Add tests with service name
	}

	for _, tc := range tests {
        // tc := tc
		t.Run(tc.Name, func(t *testing.T) {
            t.Parallel()

			entry, err := Parse(tc.logMessage)
            if tc.err == nil {
                if tc.logEntry != nil && entry == nil {
                    t.Errorf("Parse(%q) != (%q) due (%q)", tc.logEntry, entry, err)
                } else if *tc.logEntry != *entry {
                   t.Errorf("Parse(%q) != %v, got %v", tc.logEntry, tc.logEntry, entry)
               }
            }
            if tc.err != nil {
                if err == nil {
                    t.Errorf(" expecting error %v != nil",  tc.err_msg)
                } else if tc.err != nil && !strings.Contains(err.Error(), tc.err_msg) {
                    t.Errorf("Parse(%q) returned unexpected error %v != %v", tc.logEntry, tc.err_msg, err)
                 }
            }
		})
	}
}
