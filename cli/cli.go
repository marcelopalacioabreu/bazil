package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"bazil.org/bazil/cliutil/subcommands"
)

type bazil struct {
	flag.FlagSet
	Config struct {
		Verbose bool
	}
}

// Bazil allows command-line callables access to global flags, such as
// verbosity.
var Bazil = bazil{}

func init() {
	Bazil.BoolVar(&Bazil.Config.Verbose, "v", false, "verbose output")
	subcommands.Register(&Bazil)
}

// Service is an interface that commands can implement to setup and
// teardown services for the subcommands below them.
//
// As Run and potential multiple Teardown failures makes having a
// single error return impossible, Setup and Teardown only get to
// signal a boolean success. Any detail should be exposed via log.
type Service interface {
	Setup() (ok bool)
	Teardown() (ok bool)
}

func run(result subcommands.Result) (ok bool) {
	var cmd interface{}
	for _, cmd = range result.ListCommands() {
		if svc, isService := cmd.(Service); isService {
			ok = svc.Setup()
			if !ok {
				return false
			}
			defer func() {
				// Teardown failures can cause non-successful exit
				if !svc.Teardown() {
					ok = false
				}
			}()
		}
	}
	run := cmd.(subcommands.Runner)
	err := run.Run()
	if err != nil {
		log.Printf("error: %v", err)
		return false
	}
	return true
}

// Main is primary entry point into the bazil command line
// application.
func Main() (exitstatus int) {
	progName := filepath.Base(os.Args[0])
	log.SetFlags(0)
	log.SetPrefix(progName + ": ")

	result, err := subcommands.Parse(&Bazil, progName, os.Args[1:])
	if err == flag.ErrHelp {
		result.Usage()
		return 0
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", result.Name(), err)
		result.Usage()
		return 2
	}

	ok := run(result)
	if !ok {
		return 1
	}
	return 0
}