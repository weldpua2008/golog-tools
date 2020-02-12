package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
    // "github.com/weldpua2008/golog-tools/parser"
    "github.com/weldpua2008/golog-tools/trace"
    "github.com/weldpua2008/golog-tools/info"
    "github.com/cheggaaa/pb/v3"

)

// configuration
type config struct {
	input *bufio.Scanner
    pb  *pb.ProgressBar
}


// parseConfig parses command arguments.
func parseConfig() (*config, error) {
	var (
		input     *bufio.Scanner
	)
    src := flag.String("fpath", "-", "file path to read from (default: stdin)")
	flag.Parse()

	if *src == "-" {
        // Open  stdin
		input = bufio.NewScanner(os.Stdin)
        // info, err := os.Stdin.Stat()
        // if err == nil && info.Size() >0 {
        //     bar := pb.Full.Start64(info.Size())
        //     input = bar.NewProxyReader(input)
        //     defer bar.Finish()
        // }

	} else {
        // Open file
		f, err := os.Open(*src)
		if err != nil {
			return nil, fmt.Errorf("cannot open input file %q for reading",*src)
		}
		input = bufio.NewScanner(f)
	}

    n:=len(input.Bytes())
    if n < 10 {
        n = 100
    }
    // n=5000
    // bar := pb.New64(n)
    bar:= pb.StartNew(n)
    bar.SetWriter(os.Stderr)
    // bar will format numbers as bytes (B, KiB, MiB, etc)
    bar.Set(pb.Bytes, true)
	return &config{
		input: input,
        pb: bar,
	}, nil
}


func main() {
	cfg, err := parseConfig()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}


	stat := info.NewInfo()
	p := trace.NewProcessCalls(cfg.input, stat,cfg.pb)
    cfg.pb.Start()
	p.Start()
    fmt.Fprintln(os.Stderr,"Done!")

}
