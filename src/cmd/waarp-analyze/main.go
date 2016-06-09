package main

import (
	"fmt"
	"os"
	"runtime"

	"analyze"
	"logging"

	flags "github.com/jessevdk/go-flags"
)

const VERSION_NUM = "0.0.1"

type cmdArgs struct {
	// ConfRoot string `short:"c" long:"conf-root" description:"Root to Waarp instances configuration"`
	Verbose bool   `short:"v" long:"verbose" description:"Verbose output"`
	Version bool   `short:"V" long:"version" description:"Prints version and exits"`
	Output  string `short:"o" long:"output" description:"Write report to this location. Use '- for stdout"`
	Hostid  string `short:"H" long:"hostid" description:"Limit analyze to this Waarp instance"`
}

func main() {
	var args cmdArgs
	cmdParser := flags.NewParser(&args, flags.Default)
	_, err := cmdParser.Parse()
	if err != nil {
		if err, ok := err.(*flags.Error); ok && err.Type != flags.ErrHelp {
			cmdParser.WriteHelp(os.Stderr)
		}
		return
	}

	if args.Verbose {
		logging.Debug.SetOutput(os.Stderr)
	}

	if args.Version {
		fmt.Printf("Waarp Analyze %s [%s %s/%s %s]\n",
			VERSION_NUM, runtime.Version(),
			runtime.GOOS, runtime.GOARCH, runtime.Compiler)
		return
	}

	analyze.Run(args.Output, args.Hostid)
}
