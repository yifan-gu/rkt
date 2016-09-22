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

package main

import (
	"fmt"

	"github.com/appc/spec/schema/types"
	"github.com/coreos/rkt/common/apps"

	"github.com/coreos/rkt/common"
	pkgPod "github.com/coreos/rkt/pkg/pod"
	"github.com/coreos/rkt/rkt/image"
	"github.com/coreos/rkt/stage0"
	"github.com/coreos/rkt/store/imagestore"
	"github.com/coreos/rkt/store/treestore"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	cmdAppAdd = &cobra.Command{
		Use:   "add UUID IMAGEID ...",
		Short: "Add an app to a pod",
		Long:  `This allows addin an app that's present on the store to a running pod`,
		Run:   runWrapper(runAppAdd),
	}
)

func init() {
	cmdApp.AddCommand(cmdAppAdd)
	addAppFlags(cmdAppAdd)
	addIsolatorFlags(cmdAppAdd, false)

	// Disable interspersed flags to stop parsing after the first non flag
	// argument. All the subsequent parsing will be done by parseApps.
	// This is needed to correctly handle image args
	cmdAppAdd.Flags().SetInterspersed(false)
}

func runAppAdd(cmd *cobra.Command, args []string) (exit int) {
	if len(args) < 2 {
		stderr.Print("must provide the pod UUID and an IMAGEID")
		return 1
	}

	err := parseApps(&rktApps, args[1:], cmd.Flags(), true)
	if err != nil {
		stderr.PrintE("error parsing app image arguments", err)
		return 1
	}
	if rktApps.Count() > 1 {
		stderr.Print("must give only one app")
		return 1
	}

	p, err := pkgPod.PodFromUUIDString(getDataDir(), args[0])
	if err != nil {
		stderr.PrintE("problem retrieving pod", err)
		return 1
	}
	defer p.Close()

	if p.State() != pkgPod.Running {
		stderr.Printf("pod %q isn't currently running", p.UUID)
		return 1
	}

	s, err := imagestore.NewStore(storeDir())
	if err != nil {
		stderr.PrintE("cannot open store", err)
		return 1
	}

	ts, err := treestore.NewStore(treeStoreDir(), s)
	if err != nil {
		stderr.PrintE("cannot open treestore", err)
		return 1
	}

	fn := &image.Finder{
		S:  s,
		Ts: ts,
		Ks: getKeystore(),

		StoreOnly: true,
		NoStore:   false,
	}

	img, err := fn.FindImage(args[1], "", apps.AppImageGuess)
	if err != nil {
		stderr.PrintE("error finding images", err)
		return 1
	}

	ccfg := stage0.CommonConfig{
		Store:     s,
		TreeStore: ts,
		UUID:      p.UUID,
		Debug:     globalFlags.Debug,
	}

	rktgid, err := common.LookupGid(common.RktGroup)
	if err != nil {
		stderr.Printf("group %q not found, will use default gid when rendering images", common.RktGroup)
		rktgid = -1
	}

	pcfg := stage0.PrepareConfig{
		CommonConfig: &ccfg,
		Apps:         &rktApps,
	}
	rcfg := stage0.RunConfig{
		CommonConfig: &ccfg,
		UseOverlay:   p.UsesOverlay(),
		RktGid:       rktgid,
	}

	err = stage0.AddApp(pcfg, rcfg, p.Path(), img)
	if err != nil {
		stderr.PrintE("error adding app to pod", err)
		return 1
	}

	return 0
}

// generateAddConfig converts command line flags into stage0.AddConfig.
func generateAddConfig(flags *pflag.FlagSet) (*stage0.AddConfig, error) {
	var addConfig stage0.AddConfig

	flag := flags.Lookup("name")
	if flag != nil {
		value := flag.Value.String()
		if value != "" {
			name, err := types.NewACName(value)
			if err != nil {
				return nil, fmt.Errorf("invalid format for app name: %v", err)
			}
			addConfig.Name = name
		}
	}

	flag = flags.Lookup("set-annotation")
	if flag != nil {
		var annotations types.Annotations
		for k, value := range flag.Value.(*kvMap).mapping {
			key, err := types.NewACIdentifier(k)
			if err != nil {
				return nil, fmt.Errorf("invalid format for annotation key %q: %v", k, err)
			}
			annotations.Set(*key, value)
		}
		addConfig.Annotations = annotations
	}

	return &addConfig, nil
}
