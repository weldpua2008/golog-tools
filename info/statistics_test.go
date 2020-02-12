package info

import (
	"testing"
)

func TestLogMessageParsing(t *testing.T) {

		t.Run("Testing increment", func(t *testing.T) {

            in:=NewInfo()
            for i := 1;  i<=20; i++ {
                in.IncrementOrphanLines()
                in.IncrementConsumedLines()
                in.IncrementDuplicateLines()
                in.IncrementMalformedLines()
                if in.orphanLines != uint64(i) {
                    t.Errorf("orphanLines %v != %v", in.orphanLines,i)
                }
                if in.consumedLines != uint64(i) {
                    t.Errorf("consumedLines %v != %v", in.consumedLines,i)
                }
                if in.duplicateLines != uint64(i) {
                    t.Errorf("duplicateLines %v != %v", in.duplicateLines,i)
                }
                if in.malformedLines != uint64(i) {
                    t.Errorf("malformedLines %v != %v", in.malformedLines,i)
                }
            }

		})

}
