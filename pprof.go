package main

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

var (
	cpuProfileFile   = "./profile/cpu.pprof"
	memProfileFile   = "./profile/mem.pprof"
	blockProfileFile = "./profile/block.pprof"
)

func startProfile() error {
	f, err := os.Create(cpuProfileFile)
	if err != nil {
		return err
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		return err
	}
	runtime.SetBlockProfileRate(1)
	log.Println("Profile start")
	return nil
}

func endProfile() error {
	pprof.StopCPUProfile()
	runtime.SetBlockProfileRate(0)
	log.Println("Profile end")

	mf, err := os.Create(memProfileFile)
	if err != nil {
		return err
	}
	pprof.WriteHeapProfile(mf)

	bf, err := os.Create(blockProfileFile)
	if err != nil {
		return err
	}
	pprof.Lookup("block").WriteTo(bf, 0)
	return nil
}
