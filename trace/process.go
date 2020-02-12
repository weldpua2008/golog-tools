package trace

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

    "github.com/weldpua2008/golog-tools/parser"
    "github.com/weldpua2008/golog-tools/info"
    "github.com/cheggaaa/pb/v3"

)

// Process one call
type ProcessCalls struct {
	in              *bufio.Scanner
	info            info.InfoInterface
	pendingTracesWG sync.WaitGroup
	m               sync.Mutex
	logs            map[string]parser.Logs
    pb              *pb.ProgressBar
}

// NewProcessor creates a new Processor object.
func NewProcessCalls(in *bufio.Scanner, stat info.InfoInterface, progressBar *pb.ProgressBar) *ProcessCalls {
	return &ProcessCalls{
		in:               in,
		info:             stat,
		m:                sync.Mutex{},
		pendingTracesWG:  sync.WaitGroup{},
		logs:             make(map[string]parser.Logs),
        pb:                progressBar,
	}
}


// pendingTraces accumulates log entries for the postponed processing.
// Using a lock per trace ID to provide concurrency-safety.
//
// The program must keep a reasonably-sized buffer of pending traces and
// only read from standard input as long as it can keep pace. That is, using
// back-pressure to regulate the inflow.
func (p *ProcessCalls) pendingTraces(ctx context.Context, msgChan <-chan *parser.LogEntry, notifyNewTrace chan<- string ) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-msgChan:
			p.m.Lock()
            {
                if _, ok := p.logs[msg.Trace]; ! ok {
                    notifyNewTrace <- msg.Trace
                }
                p.logs[msg.Trace] = append(p.logs[msg.Trace], msg)


                // if tr, ok := p.logs[msg.Trace]; ok {
                //     tr  = append(tr, msg)
                //     // p.logs[msg.Trace] = append(p.logs[msg.Trace], msg)
                // }else {
                //     // p.logs[msg.Trace] = []parser.Logs{msg}
                //     p.logs[msg.Trace] = append(p.logs[msg.Trace], msg)
                //     notifyNewTrace <- msg.Trace
                // }
                // fmt.Fprintln(os.Stderr,"New Line", msg)
			}
			p.m.Unlock()
		}
	}
}

// Start processes in log entries into output call traces.
// func (p *ProcessCalls) Start() {
//     fmt.Fprintln(os.Stderr, "Start process")
//     ctx, cancel := context.WithCancel(context.Background())
//     defer cancel()
//     logsEntryChan := make(chan *parser.LogEntry)
//     go p.pendingTraces(ctx, logsEntryChan)
//     for p.in.Scan() {
//         parsedLogEntry, err := parser.Parse(p.in.Text())
//         if err != nil {
//             fmt.Fprintln(os.Stderr,"Malformed Lines", err)
//             p.info.IncrementMalformedLines()
//             continue
//         }
//         p.info.IncrementConsumedLines()
//         // fmt.Fprintln(os.Stderr, parsedLogEntry)
//
//
//         p.registerTrace(ctx, parsedLogEntry.Trace)
//         accumulatorChan <- parsedLogEntry
//     }
//
// }

// registerTrace waits for a new TraceID.
// Logs can be mixed up a bit (just because enforcing FIFO semantics is hard in
// distributed setups), but it is expected that the vast majority are only off
// for a few milliseconds.
func registerTrace(p *ProcessCalls, pings <-chan string, ctx context.Context) {
    for {
		select {
		case <-ctx.Done():
			return
		case traceCallId := <-pings:

            go func(p *ProcessCalls, parentCtx context.Context, traceCallId string){
                // fmt.Fprintln(os.Stderr,"New traceCallId ",traceCallId)
                // TODO:
                // Management of “pending” entries: it should be done based on
                // the timestamps of the logs and an “expiration”.
                // Any trace should be declared finished after 20 seconds 1 of the latest
                // timestamp, and sent to output
                d := time.Now().Add(20 * time.Second)
                childCtx, childCancel := context.WithDeadline(parentCtx,d)
                defer childCancel()
                // select {
            	// // case <-time.After(1 * time.Second):
            	// // 	fmt.Println("overslept")
            	// case <-childCtx.Done():
            	// 	fmt.Println(ctx.Err())
            	// }
                // fmt.Fprintln(os.Stderr,"processAfterContextExpiration")
                p.processAfterContextExpiration(childCtx, traceCallId)
                // fmt.Fprintln(os.Stderr,"Processed traceCallId ",traceCallId)

            }(p,ctx,traceCallId)
        default:
            // fmt.Fprintln(os.Stderr,"// WARNING: Wait for ping")
            time.Sleep(100 * time.Millisecond)
            continue
        }
    }

}


// Start processes in log entries into output call traces with progress bar.
func (p *ProcessCalls) Start() {
    fmt.Fprintln(os.Stderr, "Start process")
	ctx, cancel := context.WithCancel(context.Background())
	accumulatorChan := make(chan *parser.LogEntry)
    notifyNewTrace := make(chan string)
	// TODO: Here is a place for different scaling strategies: we can use multiple accumulator instances here.
	go p.pendingTraces(ctx, accumulatorChan, notifyNewTrace)
    go registerTrace(p, notifyNewTrace,ctx)
    count := int64(0)
    increaseFactor := int64(1)
	for p.in.Scan() {
        rawLogEntry := p.in.Text()
        count = count + int64(len([]byte(rawLogEntry)))
        if p.pb != nil {
            p.pb.Add(len([]byte(rawLogEntry)))
            if  count  >= p.pb.Total() {
                // we can use more advanced methods for grow of the progressBar
			    // if  count > int64((p.pb.Total() * (100 - decreaseFactor) /100))  {
                // if increaseFactor < 10 {
                        // increaseFactor++
                // }
                p.pb.SetTotal(p.pb.Total() *(100 + increaseFactor)/100 )
			}
		}
		parsedLogEntry, err := parser.Parse(p.in.Text())
		if err != nil {
            fmt.Fprintln(os.Stderr,"Malformed Lines", err)
			p.info.IncrementMalformedLines()
			continue
		}
        p.info.IncrementConsumedLines()
        // fmt.Fprintln(os.Stderr, parsedLogEntry)


		// p.registerTrace(ctx, parsedLogEntry.Trace)
		accumulatorChan <- parsedLogEntry
	}

	// Calling cancel signals to all postponed trace processors to generate the result.
	cancel()
	// Await until all processing goroutines finish.
	p.pendingTracesWG.Wait()
    // fmt.Fprintln(os.Stderr, p.logs)
    fmt.Fprintln(os.Stderr, p.info)
}

// processAfterContextExpiration is a delayed trace ID processor.
// Processing needs to be postponed because the entries in the log file are in
// a random order, so processing timeout
// ensures that some earlier entries will be collected.
func (p *ProcessCalls) processAfterContextExpiration(ctx context.Context, traceId string) {
    p.pendingTracesWG.Add(1)
    defer p.pendingTracesWG.Done()
    <-ctx.Done()


    if tr, ok := p.logs[traceId]; ok {
        if !tr.HasRootCallerSpan() {
            // fmt.Fprintln(os.Stderr,"HasRootCallerSpan is false",tr)
            return
        }
    }

    // ll := p.logs[traceId]
    // fmt.Fprintln(os.Stderr,traceId, "ll.String",ll.String())
    // if ll.HasRootCallerSpan() == false {
        // fmt.Fprintln(os.Stderr,ll.String())
        // fmt.Fprintln(os.Stderr,"HasRootCallerSpan is false")
        // return
    // }
    // } else {
    //     fmt.Fprintln(os.Stderr,ll.String())
    //
    // }
    // fmt.Fprintln(os.Stderr,"==> buildTraceTree ",p.logs[traceId], traceId, p.info)
    var result traceTree
    tr :=append(p.logs[traceId][:0:0], p.logs[traceId]...)
    // if tr, ok := p.logs[traceId]; ok {
            result = buildTraceTree(&tr, traceId, p.info)
            // fmt.Fprintln(os.Stderr,"==> result.Root ",result.Root)
            if result.Root == nil {
                return
            }

    // }else {
    //     return
    // }
    // result := buildTraceTree(p.logs[traceId], traceId, p.info)
    // fmt.Fprintln(os.Stderr,"result %v",result)
    res, err := json.Marshal(result)
    if err != nil {
        // TODO: Log an error and update the statistics.
        fmt.Fprintln(os.Stderr, err)
    }
    // fmt.Fprintln(os.Stderr,ll.String())
    p.m.Lock()
    defer p.m.Unlock()
    // delete(p.logs, traceId)

    // TODO: output this result to a proper source destination.
    fmt.Fprintln(os.Stdout, string(res))
}
