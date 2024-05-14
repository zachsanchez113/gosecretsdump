package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/C-Sto/gosecretsdump/cmd"
)

var version string

//const version = "0.3.0"

func main() {
	if version == "" {
		version = "DEV"
	}

	fmt.Println("gosecretsdump v" + version + " (@C__Sto)")

	args := cmd.CLIArgs{}

	var vers bool

	flag.StringVar(&args.Outfile, "out", "", "Location to export output")
	flag.StringVar(&args.NTDSLoc, "ntds", "", "Location of the NTDS file (required)")
	flag.StringVar(&args.SystemLoc, "system", "", "Location of the SYSTEM file (required)")
	flag.StringVar(&args.SAMLoc, "sam", "", "Location of SAM registry hive")
	flag.BoolVar(&args.LiveSAM, "livesam", false, "Get hashes from live system. Only works on local machine hashes (SAM), only works on Windows.")
	flag.BoolVar(&args.Status, "status", false, "Include status in hash output")
	flag.BoolVar(&args.EnabledOnly, "enabled", false, "Only output enabled accounts")
	flag.BoolVar(&args.NoPrint, "noprint", false, "Don't print output to screen (probably use this with the -out flag)")
	flag.BoolVar(&args.Stream, "stream", false, "Stream to files rather than writing in a block. Can be much slower.")
	flag.BoolVar(&vers, "version", false, "Print version and exit")
	flag.BoolVar(&args.History, "history", false, "Include Password History")
	flag.Parse()

	if vers {
		os.Exit(0)
	}

	if args.SystemLoc == "" && (args.NTDSLoc == "" && args.SAMLoc == "") && !args.LiveSAM {
		flag.Usage()
		os.Exit(1)
	}

	// e := cmd.GoSecretsDump(s)

	e := cmd.GoSecretsDumpJSON(args)
	if e != nil {
		panic(e)
	}
}

//info dumped out of https://github.com/SecureAuthCorp/impacket/blob/master/impacket/examples/secretsdump.py
