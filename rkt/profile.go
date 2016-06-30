// Copyright 2016 The rkt Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build profile

package main

import (
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/spf13/cobra"
)

var cpuprofile string
var memprofile string

func parseProfileFlags(cmdRkt *cobra.Command) {
	cmdRkt.PersistentFlags().StringVar(&cpuprofile, "cpuprofile", "", "write CPU profile to the file")
	cmdRkt.PersistentFlags().StringVar(&memprofile, "memprofile", "", "write memory profile to the file")
}

func startProfile() (cpufile *os.File, memfile *os.File, err error) {
	if cpuprofile != "" {
		cpufile, err = os.Create(cpuprofile)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot create cpu profile file %q: %v", cpuprofile, err)
		}
		pprof.StartCPUProfile(cpufile)
	}
	if memprofile != "" {
		memfile, err = os.Create(memprofile)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot create memory profile file %q: %v", memprofile, err)
		}
	}
	return cpufile, memfile, nil
}

func stopProfile(cpuprofile, memprofile *os.File) {
	pprof.StopCPUProfile()
	pprof.WriteHeapProfile(memprofile)
	cpuprofile.Close()
	memprofile.Close()
}
