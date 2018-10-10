package main

import (
	"log"
	"os"
)

func init() {
	log.SetFlags(0)
}

func main() {
	startProfile()
	cli := &CLI{outStream: os.Stdout, errStream: os.Stderr}
    s := cli.Run(os.Args)
	endProfile()
	os.Exit(s)
}
