// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os/exec"
	"strings"
	"time"
)

var verbflag = flag.Int("v", 0, "Verbose trace output level")
var hoursflag = flag.Int("hours", 0, "duration of run in hours")

// Ping motes every 5 minutes
const sleepDurationInSeconds = 5 * 60

func verb(vlevel int, s string, a ...interface{}) {
	if *verbflag >= vlevel {
		fmt.Printf(s, a...)
		fmt.Printf("\n")
	}
}

func doGomoteCmd(gcmd []string) []string {
	verb(1, "gomote command is: %+v", gcmd)
	cmd := exec.Command("gomote", gcmd...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("%s\n", string(out))
		log.Fatalf("gomote command %+v failed: %v", cmd, err)
	}
	lines := strings.Split(string(out), "\n")
	return lines
}

func pingMote(mote string) {
	verb(1, "pinging %s", mote)
	doGomoteCmd([]string{"ls", mote})
}

func pingMotes() {
	verb(1, "pinging all motes")
	mlines := doGomoteCmd([]string{"list"})
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
		pingMote(mote)
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
	doGomoteCmd([]string{"list"}) // make sure gomote works ok
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
