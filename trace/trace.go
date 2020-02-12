package trace

import (
    "github.com/weldpua2008/golog-tools/parser"
    "github.com/weldpua2008/golog-tools/info"
    "container/heap"
    "fmt"
    "os"
)

// buildTraceTree implementation for the building the call-tree.
// TODO:
//  - filter cilcular calls a-> a
//  - filter orphan calls
//  - filter duplicate calls
func buildTraceTree(logs *parser.Logs, traceId string, stat info.InfoInterface) (result traceTree) {
    // var orphanTraceCalls []*traceCall
    mapping:= make(map[string]*traceCall)
    orphanMapping:= make(map[string][]*traceCall)
    result = traceTree{
        ID: traceId,
    }
     // fmt.Fprintln(os.Stderr,"logs",logs)
    if logs.Len() > 0 {
        heap.Init(logs)
        // rs := heap.Pop(&logs).(*parser.LogEntry)
        // tcall := traceCall{
		// 	Start:   rs.Start,
		// 	End:     rs.End,
		// 	Service: rs.Service,
		// 	Span:    rs.Span,
		// 	Calls:   make([]*traceCall, 0),
		// }
        //
        // mapping[rs.Span] = tcall

        for logs.Len() > 0 {
            rs := heap.Pop(logs).(*parser.LogEntry)
            tcall := traceCall{
    			Start:   rs.Start,
    			End:     rs.End,
    			Service: rs.Service,
    			Span:    rs.Span,
    			Calls:   make([]*traceCall, 0),
    		}
            // fmt.Fprintln(os.Stderr,rs.CallerSpan,parser.RootSpanConst )
            // if  result.Root == nil && rs.CallerSpan == parser.RootSpanConst {
            // fmt.Fprintln(os.Stderr, "rs.Span + ", rs.Span, rs.CallerSpan)
            if rs.CallerSpan == parser.RootSpanConst && result.Root == nil {
                // fmt.Fprintln(os.Stderr,"!!! rs.CallerSpan == parser.RootSpanConst ")
                result.Root = &tcall
                // mapping[parser.RootSpanConst] = &tcall

            }  else if rs.CallerSpan == rs.Span {
                // fmt.Fprintln(os.Stderr,"!!! rs.CallerSpan (%v) == rs.Span (%v) ", rs.CallerSpan, rs.Span)
                stat.IncrementMalformedLines()
                continue
            }else if _, ok := mapping[rs.Span]; ok {
                stat.IncrementDuplicateLines()
                continue
            }

            mapping[rs.Span] = &tcall
            if rs.CallerSpan == parser.RootSpanConst {
            } else if parent, ok := mapping[rs.CallerSpan]; ok {
                // fmt.Fprintln(os.Stderr,"!!! append (%v) <= rs.Span (%v) ", parent.Span, tcall.Span)

                parent.Calls  = append(parent.Calls, &tcall)
            } else {
                // fmt.Fprintln(os.Stderr,"!!! orphanMapping (%v) <= rs.Span (%v) ", rs.CallerSpan, tcall.Span)

                orphanMapping[rs.CallerSpan] = append(orphanMapping[rs.CallerSpan],&tcall)
            }
        }
        for orphanCallerSpan := range orphanMapping {
            for i := range orphanMapping[orphanCallerSpan] {
                // fmt.Fprintln(os.Stderr, "orphan line",orphan)


                if parent, ok := mapping[orphanCallerSpan]; ok {
                    parent.Calls  = append(parent.Calls, orphanMapping[orphanCallerSpan][i])
                } else {
                    stat.IncrementOrphanLines()
                    fmt.Fprintln(os.Stderr, "orphan line", orphanMapping[orphanCallerSpan][i])
                    // orphanTraceCalls = append(orphanTraceCalls, &orphan)
                }
            }
        }
    }
    // fmt.Fprintln(os.Stderr,result)
    return result

}
