package cmd

import (
	"errors"
	"os"
	"sync"

	"github.com/C-Sto/gosecretsdump/pkg/ditreader"
)

type DumperJSON interface {
	GetOutChan() <-chan ditreader.DumpedHash
	DumpJSON() error
}

// CLI entrypoint for dumping an AD database to JSON
//
// TODO: Stop panicking so much
func GoSecretsDumpJSON(args CLIArgs) error {
	var dr DumperJSON
	var err error

	if args.NTDSLoc != "" {
		dr, err = ditreader.New(args.SystemLoc, args.NTDSLoc)
		if err != nil {
			return err
		}
	}

	dataChannel := dr.GetOutChan()
	wg := sync.WaitGroup{}
	wg.Add(1)

	if args.Stream {
		panic(errors.New("stream output is not supported for JSON output"))
	}

	if args.Outfile != "" && !args.Stream {
		go fileWriterJSON(dataChannel, args, &wg)
	} else if args.Outfile != "" {
		panic(errors.New("please provide an output file"))
	} else if args.Stream {
		panic(errors.New("stream output is not supported for JSON output"))
	} else {
		panic(errors.New("console output is not supported for JSON output"))
	}

	err = dr.DumpJSON()
	if err != nil {
		return err
	}

	wg.Wait()
	return err
}

// Goroutine for writing JSON output to the target file
func fileWriterJSON(val <-chan ditreader.DumpedHash, args CLIArgs, wg *sync.WaitGroup) {
	defer wg.Done()

	// Open + truncate the file
	file, err := os.OpenFile(args.Outfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	if err := file.Truncate(0); err != nil {
		panic(err)
	}

	// Write hashes from the channel
	for dh := range val {
		if _, err := file.WriteString(dh.JsonString); err != nil {
			panic(err)
		}
	}
}
