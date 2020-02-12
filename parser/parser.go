package parser

import (
	"fmt"
	"regexp"
	"time"
)

// Parse parses a log line into a LogEntry
func Parse(l string) (*LogEntry, error) {
	// TODO: Add validation, handle more time formats.

	logMessagePattern := fmt.Sprintf(
        "(%[1]s) (%[1]s) (%[2]s) (%[3]s) (%[4]s)->(%[4]s)",
        dateRegexp,
        traceRegExp,
        serviceRegExp,
        spanRegExp)
	logRe  := regexp.MustCompile(logMessagePattern)

	fields := logRe.FindStringSubmatch(l)
	if len(fields) == 0 {
		return nil, fmt.Errorf("No match - log entry is in wrong format")
	}

	start      := fields[1]
	end        := fields[2]
	trace      := fields[3]
	service    := fields[4]
	callerSpan := fields[5]
	span       := fields[6]

	startDate, err := time.Parse(time.RFC3339Nano, start)
	if err != nil {
		return nil, fmt.Errorf("log message start date %s is invalid", start)
	}

	endDate, err := time.Parse(time.RFC3339Nano, end)
	if err != nil {
		return nil, fmt.Errorf("log message end date %s is invalid", end)
	}

	return &LogEntry{
		Service:      service,
		Trace:        trace,
		CallerSpan:   callerSpan,
		Span:         span,
		Start:        startDate,
		End:          endDate,
	}, nil
}
