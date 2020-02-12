package trace

import (
	"encoding/json"
	"time"
)

// resul for the execution trace.
type traceTree struct {
	ID   string            `json:"id"`
	Root *traceCall        `json:"root"`
}

// an output trace call entry per one request.
// We can implement output format
type traceCall struct {
	Start   time.Time      `json:"start"`
	End     time.Time      `json:"end"`
	Service string         `json:"service"`
	Span    string         `json:"span"`
	Calls   []*traceCall   `json:"calls"`
}

// MarshalJSON overrides marshaling of the time fields.
func (tc *traceCall) MarshalJSON() ([]byte, error) {
	type alias traceCall
	return json.Marshal(&struct {
		Start string      `json:"start"`
		End   string      `json:"end"`
		*alias
	}{
		Start: tc.Start.Format(time.RFC3339Nano),
		End:   tc.End.Format(time.RFC3339Nano),
		alias: (*alias)(tc),
	})
}
