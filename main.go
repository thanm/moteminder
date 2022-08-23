// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strings"
	"time"
)

var verbflag = flag.Int("v", 0, "Verbose trace output level")
var hoursflag = flag.Int("hours", 0, "duration of run in hours")

// Ping motes every 5 minutes
const sleepDurationInSeconds = 5 * 60

// mote command behavior flags
type moteCmdFlags uint

const (
	MOTE_CMD_RETRY        moteCmdFlags = 1
	MOTE_CMD_DIE_ON_ERROR moteCmdFlags = 2
	MOTE_CMD_LOGERR       moteCmdFlags = 4
	MOTE_CMD_V2_MOTE      moteCmdFlags = 8
)

func verb(vlevel int, s string, a ...interface{}) {
	if *verbflag >= vlevel {
		fmt.Printf(s, a...)
		fmt.Printf("\n")
	}
}

const logfile = "/tmp/mm.errs.txt"

func startErrLog(out string) {
	writeToErrLog(out, true)
}

func appendToErrLog(out string) {
	writeToErrLog(out, false)
}

func writeToErrLog(out string, trunc bool) {
	flags := os.O_WRONLY | os.O_CREATE | os.O_APPEND
	if trunc {
		flags = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	}
	of, oerr := os.OpenFile(logfile, flags, 0666)
	if oerr != nil {
		log.Fatalf("opening %s: %v\n", logfile, oerr)
	}
	fmt.Fprintf(of, "\n\n-------------\nat %v:\n%s\n", time.Now(), out)
	of.Close()
}

func doGomoteCmd(flags moteCmdFlags, gcmd []string) []string {
	const retries = 3
	if (flags & MOTE_CMD_V2_MOTE) != 0 {
		gcmd = append([]string{"v2"}, gcmd...)
	}
	verb(1, "gomote command is: %+v", gcmd)
	for i := 0; i < retries; i++ {
		cmd := exec.Command("gomote", gcmd...)
		out, err := cmd.CombinedOutput()
		if err == nil {
			lines := strings.Split(string(out), "\n")
			return lines
		}
		if err != nil {
			log.Printf("%s\n", string(out))
			log.Printf("gomote command %+v failed: %v", gcmd, err)
			if (flags & MOTE_CMD_LOGERR) != 0 {
				appendToErrLog(fmt.Sprintf("gomote cmd: %+v\noutput: %s\nerr: %v\n", gcmd, string(out), err))
			}
			if (flags & MOTE_CMD_DIE_ON_ERROR) != 0 {
				log.Fatalf("fatal error, exiting")
			}
			if (flags & MOTE_CMD_RETRY) == 0 {
				break
			}
			verb(1, "try %d failed, trying again", i)
		}
	}
	return []string{}
}

func pingMote(flags moteCmdFlags, mote string) {
	verb(1, "pinging %s", mote)
	doGomoteCmd(flags, []string{"ping", mote})
}

func pingMotes() {
	verb(1, "pinging all motes")
	basemode := MOTE_CMD_RETRY | MOTE_CMD_LOGERR
	modes := []moteCmdFlags{basemode}
	for _, mode := range modes {
		mlines := doGomoteCmd(mode, []string{"list"})
		for _, line := range mlines {
			line = strings.Trim(string(line), " \t\n")
			if line == "" {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) < 2 {
				log.Fatalf("unexpected output line from 'gomote list': %s", line)
			}
			mote := fields[0]
			pingMote(mode, mote)
		}
	}
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("moteminder: ")
	flag.Parse()
	verb(1, "starting main")
	if flag.NArg() != 0 {
		log.Fatalf("unknown extra args")
	}
	doGomoteCmd(MOTE_CMD_DIE_ON_ERROR, []string{"list"}) // make sure gomote works ok
	startErrLog(fmt.Sprintf("starting session, duration %d hours", *hoursflag))
	duration := math.MaxInt32
	if *hoursflag != 0 {
		if *hoursflag < 0 {
			log.Fatalf("please use positive arg for -hours option")
		}
		duration = *hoursflag
		verb(1, "using duration %d", duration)
	}
	for h := 0; h < duration; h++ {
		verb(1, "starting hour %d", h)
		// run for an hour
		for c := 0; c < 12; c++ {
			pingMotes()
			verb(1, "about to sleep...")
			time.Sleep(sleepDurationInSeconds * time.Second)
			verb(1, "... sleep done")
		}
	}
	verb(1, "done")
}
