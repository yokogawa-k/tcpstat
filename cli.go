package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"text/tabwriter"
	"time"

	"github.com/yokogawa-k/tcpstat/stats"
)

const (
	exitCodeOK  = 0
	exitCodeErr = 10 + iota
)

// CLI is the command line object.
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Run execute the main process.
// It returns exit code.
func (c *CLI) Run(args []string) int {
	log.SetOutput(c.errStream)

	var (
		all        bool
		continuous bool
		numeric    bool
		program    bool
		json       bool
		ver        bool
	)
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(c.errStream)
	flags.Usage = func() {
		_, err := fmt.Fprint(c.errStream, helpText)
		if err != nil {
			log.Printf("failed to format help text: %v", err)
		}
	}
	flags.BoolVar(&numeric, "n", false, "")
	flags.BoolVar(&numeric, "numeric", false, "")
	flags.BoolVar(&all, "all", false, "")
	flags.BoolVar(&all, "a", false, "")
	flags.BoolVar(&continuous, "continuous", false, "")
	flags.BoolVar(&continuous, "c", false, "")
	flags.BoolVar(&program, "program", false, "")
	flags.BoolVar(&program, "p", false, "")
	flags.BoolVar(&json, "json", false, "")
	flags.BoolVar(&ver, "version", false, "")
	flags.BoolVar(&ver, "v", false, "")
	if err := flags.Parse(args[1:]); err != nil {
		return exitCodeErr
	}

	if ver {
		_, err := fmt.Fprintf(c.errStream, "%s version %s, build %s \n", name, version, commit)
		if err != nil {
			log.Printf("failed to print version strings: %v", err)
		}
		return exitCodeOK
	}

	stats, err := stats.GetStats(all, program)
	if err != nil {
		log.Printf("failed to get stats: %v", err)
		return exitCodeErr
	}

	switch {
	case continuous:
		for {
			c.PrintStats(stats, numeric)
			time.Sleep(1 * time.Second)
		}
	case json:
		if err := c.PrintStatsAsJSON(stats, numeric); err != nil {
			log.Printf("failed to print json: %v", err)
			return exitCodeErr
		}
	default:
		c.PrintStats(stats, numeric)
	}

	return exitCodeOK
}

// PrintStats prints the tcp staticstics.
func (c *CLI) PrintStats(stats stats.Stats, numeric bool) {
	// Format in tab-separated columns with a tab stop of 8.
	tw := tabwriter.NewWriter(c.outStream, 0, 8, 0, '\t', 0)
	fmt.Fprintln(tw, "Proto\tRecv-Q\tSend-Q\tLocal Address\tForeign Address\tState")
	for _, stat := range stats {
		stat.ReplacePortName()
		if !numeric {
			stat.ReplaceLookupedName()
		}
		fmt.Fprintf(tw, "%s\t%d\t%d\t%s:%s\t%s:%s\t%s\n",
			"tcp", stat.RecvQ, stat.SendQ,
			stat.Local.Addr, stat.Local.Port,
			stat.Foreign.Addr, stat.Foreign.Port,
			stat.State,
		)
	}
	tw.Flush()
}

// PrintStatsAsJSON prints the tcp staticstics as json format.
func (c *CLI) PrintStatsAsJSON(stats stats.Stats, numeric bool) error {
	for _, stat := range stats {
		if !numeric {
			stat.ReplaceLookupedName()
		}
	}
	b, err := json.Marshal(stats)
	if err != nil {
		return err
	}
	c.outStream.Write(b)
	return nil
}

var helpText = `Usage: tcpstat [options]

  Print tcp connection staticstics

Options:
  --numeric, -n             show numerical addresses instead of trying to determine symbolic host names.
  --all, -a                 print both listening and non-listening sockets
  --continuous, -c          print the selected information every second continuously. This option does not support json format
  --program, -p             print the PID and name of the program to which each socket belongs(unimpremented)
  --json                    print results as json format
  --version, -v	            print version
  --help, -h                print help
`
