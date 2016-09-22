// Copyright 2015 The rkt Authors
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

//+build linux

package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/coreos/rkt/common/apps"

	"github.com/appc/spec/schema"
	"github.com/appc/spec/schema/types"
	"github.com/spf13/pflag"
)

var (
	rktApps apps.Apps // global used by run/prepare for representing the apps expressed via the cli
)

// parseApps looks through the args for support of per-app argument lists delimited with "--" and "---".
// Between per-app argument lists flags.Parse() is called using the supplied FlagSet.
// Anything not consumed by flags.Parse() and not found to be a per-app argument list is treated as an image.
// allowAppArgs controls whether "--" prefixed per-app arguments will be accepted or not.
func parseApps(al *apps.Apps, args []string, flags *pflag.FlagSet, allowAppArgs bool) error {
	nAppsLastAppArgs := al.Count()

	// valid args here may either be:
	// not-"--"; flags handled by *flags or an image specifier
	// "--"; app arguments begin
	// "---"; conclude app arguments
	// between "--" and "---" pairs anything is permitted.
	inAppArgs := false
	for i := 0; i < len(args); i++ {
		a := args[i]
		if inAppArgs {
			switch a {
			case "---":
				// conclude this app's args
				inAppArgs = false
			default:
				// keep appending to this app's args
				app := al.Last()
				app.Args = append(app.Args, a)
			}
		} else {
			switch a {
			case "--":
				if !allowAppArgs {
					return fmt.Errorf("app arguments unsupported")
				}
				// begin app's args
				inAppArgs = true

				// catch some likely mistakes
				if nAppsLastAppArgs == al.Count() {
					if al.Count() == 0 {
						return fmt.Errorf("an image is required before any app arguments")
					}
					return fmt.Errorf("only one set of app arguments allowed per image")
				}
				nAppsLastAppArgs = al.Count()
			case "---":
				// ignore triple dashes since they aren't images
				// TODO(vc): I don't think ignoring this is appropriate, probably should error; it implies malformed argv.
				// "---" is not an image separator, it's an optional argument list terminator.
				// encountering it outside of inAppArgs is likely to be "--" typoed
			default:
				// consume any potential inter-app flags
				if err := flags.Parse(args[i:]); err != nil {
					return err
				}
				nInterFlags := (len(args[i:]) - flags.NArg())

				if nInterFlags > 0 {
					// XXX(vc): flag.Parse() annoyingly consumes the "--", reclaim it here if necessary
					if args[i+nInterFlags-1] == "--" {
						nInterFlags--
					}

					// advance past what flags.Parse() consumed
					i += nInterFlags - 1 // - 1 because of i++
				} else {
					// flags.Parse() didn't want this arg, treat as image
					al.Create(a)
				}
			}
		}
	}

	return al.Validate()
}

// Value interface implementations for the various per-app fields we provide flags for

// appAsc is for aci --signature
type appAsc apps.Apps

func (aa *appAsc) Set(s string) error {
	app := (*apps.Apps)(aa).Last()
	if app == nil {
		return fmt.Errorf("--signature must follow an image")
	}
	if app.Asc != "" {
		return fmt.Errorf("--signature specified multiple times for the same image")
	}
	app.Asc = s

	return nil
}

func (aa *appAsc) String() string {
	app := (*apps.Apps)(aa).Last()
	if app == nil {
		return ""
	}
	return app.Asc
}

func (aa *appAsc) Type() string {
	return "appAsc"
}

// appExec is for aci --exec overrides
type appExec apps.Apps

func (ae *appExec) Set(s string) error {
	app := (*apps.Apps)(ae).Last()
	if app == nil {
		return fmt.Errorf("--exec must follow an image")
	}
	if app.Exec != "" {
		return fmt.Errorf("--exec specified multiple times for the same image")
	}
	app.Exec = s

	return nil
}

func (ae *appExec) String() string {
	app := (*apps.Apps)(ae).Last()
	if app == nil {
		return ""
	}
	return app.Exec
}

func (ae *appExec) Type() string {
	return "appExec"
}

// appMount is for --mount flags in the form of: --mount volume=VOLNAME,target=PATH
type appMount apps.Apps

func (al *appMount) Set(s string) error {
	mount := schema.Mount{}

	// this is intentionally made similar to types.VolumeFromString()
	// TODO(iaguis) use MakeQueryString() when appc/spec#520 is merged
	m, err := url.ParseQuery(strings.Replace(s, ",", "&", -1))
	if err != nil {
		return err
	}

	for key, val := range m {
		if len(val) > 1 {
			return fmt.Errorf("label %s with multiple values %q", key, val)
		}
		switch key {
		case "volume":
			mv, err := types.NewACName(val[0])
			if err != nil {
				return fmt.Errorf("invalid volume name %q in --mount flag %q: %v", val[0], s, err)
			}
			mount.Volume = *mv
		case "target":
			mount.Path = val[0]
		default:
			return fmt.Errorf("unknown mount parameter %q", key)
		}
	}

	as := (*apps.Apps)(al)
	if as.Count() == 0 {
		as.Mounts = append(as.Mounts, mount)
	} else {
		app := as.Last()
		app.Mounts = append(app.Mounts, mount)
	}

	return nil
}

func (al *appMount) String() string {
	var ms []string
	for _, m := range ((*apps.Apps)(al)).Mounts {
		ms = append(ms, m.Volume.String(), ":", m.Path)
	}
	return strings.Join(ms, " ")
}

func (al *appMount) Type() string {
	return "appMount"
}

// appsVolume is for --volume flags in the form name,kind=host,source=/tmp,readOnly=true,recursive=true (defined by appc)
type appsVolume apps.Apps

func (al *appsVolume) Set(s string) error {
	vol, err := types.VolumeFromString(s)
	if err != nil {
		return fmt.Errorf("invalid value in --volume flag %q: %v", s, err)
	}

	(*apps.Apps)(al).Volumes = append((*apps.Apps)(al).Volumes, *vol)
	return nil
}

func (al *appsVolume) Type() string {
	return "appsVolume"
}

func (al *appsVolume) String() string {
	var vs []string
	for _, v := range (*apps.Apps)(al).Volumes {
		vs = append(vs, v.String())
	}
	return strings.Join(vs, " ")
}

// appMemoryLimit is for --memory flags in the form of: --memory=128M
type appMemoryLimit apps.Apps

func (aml *appMemoryLimit) Set(s string) error {
	app := (*apps.Apps)(aml).Last()
	if app == nil {
		return fmt.Errorf("--memory must follow an image")
	}
	isolator, err := types.NewResourceMemoryIsolator(s, s)
	if err != nil {
		return err
	}
	// Just don't accept anything less than 4Ki. It's not reasonable and
	// it's most likely a mistake from the user, such as passing
	// --memory=16m (milli-bytes!) instead of --memory=16M (megabytes).
	if isolator.Limit().Value() < 4096 {
		return fmt.Errorf("memory limit of %d bytes too low. Try a reasonable value, such as --memory=16M", isolator.Limit().Value())
	}
	app.MemoryLimit = isolator
	return nil
}

func (aml *appMemoryLimit) String() string {
	app := (*apps.Apps)(aml).Last()
	if app == nil {
		return ""
	}
	return app.MemoryLimit.String()
}

func (aml *appMemoryLimit) Type() string {
	return "appMemoryLimit"
}

// appCPULimit is for --cpu flags in the form of: --cpu=500m
type appCPULimit apps.Apps

func (aml *appCPULimit) Set(s string) error {
	app := (*apps.Apps)(aml).Last()
	if app == nil {
		return fmt.Errorf("--cpu must follow an image")
	}
	isolator, err := types.NewResourceCPUIsolator(s, s)
	if err != nil {
		return err
	}
	app.CPULimit = isolator
	return nil
}

func (aml *appCPULimit) String() string {
	app := (*apps.Apps)(aml).Last()
	if app == nil {
		return ""
	}
	return app.CPULimit.String()
}

func (aml *appCPULimit) Type() string {
	return "appCPULimit"
}

// appUser is for --user flags in the form of: --user=user
type appUser apps.Apps

func (au *appUser) Set(s string) error {
	app := (*apps.Apps)(au).Last()
	if app == nil {
		return fmt.Errorf("--user must follow an image")
	}
	app.User = s
	return nil
}

func (au *appUser) String() string {
	app := (*apps.Apps)(au).Last()
	if app == nil {
		return ""
	}
	return app.User
}

func (au *appUser) Type() string {
	return "appUser"
}

// appGroup is for --group flags in the form of: --group=group
type appGroup apps.Apps

func (ag *appGroup) Set(s string) error {
	app := (*apps.Apps)(ag).Last()
	if app == nil {
		return fmt.Errorf("--group must follow an image")
	}
	app.Group = s
	return nil
}

func (ag *appGroup) String() string {
	app := (*apps.Apps)(ag).Last()
	if app == nil {
		return ""
	}
	return app.Group
}

func (ag *appGroup) Type() string {
	return "appGroup"
}

// appCapsRetain is for --caps-retain flags in the form of: --caps-retain=CAP_KILL,CAP_NET_ADMIN
type appCapsRetain apps.Apps

func (au *appCapsRetain) Set(s string) error {
	app := (*apps.Apps)(au).Last()
	if app == nil {
		return fmt.Errorf("--caps-retain must follow an image")
	}
	capsRetain, err := types.NewLinuxCapabilitiesRetainSet(strings.Split(s, ",")...)
	if err != nil {
		return err
	}
	app.CapsRetain = capsRetain
	return nil
}

func (au *appCapsRetain) String() string {
	app := (*apps.Apps)(au).Last()
	if app == nil {
		return ""
	}
	var vs []string
	for _, v := range app.CapsRetain.Set() {
		vs = append(vs, string(v))
	}
	return strings.Join(vs, ",")
}

func (au *appCapsRetain) Type() string {
	return "appCapsRetain"
}

// appCapsRemove is for --caps-remove flags in the form of: --caps-remove=CAP_MKNOD,CAP_SYS_CHROOT
type appCapsRemove apps.Apps

func (au *appCapsRemove) Set(s string) error {
	app := (*apps.Apps)(au).Last()
	if app == nil {
		return fmt.Errorf("--caps-retain must follow an image")
	}
	capsRemove, err := types.NewLinuxCapabilitiesRevokeSet(strings.Split(s, ",")...)
	if err != nil {
		return err
	}
	app.CapsRemove = capsRemove
	return nil
}

func (au *appCapsRemove) String() string {
	app := (*apps.Apps)(au).Last()
	if app == nil {
		return ""
	}
	var vs []string
	for _, v := range app.CapsRemove.Set() {
		vs = append(vs, string(v))
	}
	return strings.Join(vs, ",")
}

func (au *appCapsRemove) Type() string {
	return "appCapsRemove"
}

// appSeccompFilter is for --seccomp flags in the form of: --seccomp [errno=EPERM,]mode=retain,chown,chmod[,syscalls]
type appSeccompFilter apps.Apps

func (au *appSeccompFilter) Set(s string) error {
	app := (*apps.Apps)(au).Last()
	if app == nil {
		return fmt.Errorf("--seccomp must follow an image")
	}
	app.SeccompFilter = s
	return nil
}

func (au *appSeccompFilter) String() string {
	app := (*apps.Apps)(au).Last()
	if app == nil {
		return ""
	}
	return app.SeccompFilter
}

func (au *appSeccompFilter) Type() string {
	return "appSeccompFilter"
}

// appName is for --name flags in the form of: --name=APPNAME.
type appName apps.Apps

func (au *appName) Set(s string) error {
	app := (*apps.Apps)(au).Last()
	if app == nil {
		return fmt.Errorf("--name must follow an image")
	}
	app.Name = s
	return nil
}

func (au *appName) String() string {
	app := (*apps.Apps)(au).Last()
	if app == nil {
		return ""
	}
	return app.Name
}

func (au *appName) Type() string {
	return "appName"
}

// appAnnotation is for --annotation flags in the form of: --annotation=NAME=VALUE.
type appAnnotation apps.Apps

func (au *appAnnotation) Set(s string) error {
	app := (*apps.Apps)(au).Last()
	if app == nil {
		return fmt.Errorf("--annotation must follow an image")
	}

	fields := strings.SplitN(s, "=", 2)
	if len(fields) != 2 {
		return fmt.Errorf("invalid format of --annotation flag %q", s)
	}

	if app.CRIAnnotations == nil {
		app.CRIAnnotations = make(map[string]string)
	}
	app.CRIAnnotations[fields[0]] = fields[1]
	return nil
}

func (au *appAnnotation) String() string {
	app := (*apps.Apps)(au).Last()
	if app == nil {
		return ""
	}
	var annotations []string
	for name, value := range app.CRIAnnotations {
		annotations = append(annotations, fmt.Sprintf("%s=%s", name, value))
	}
	return strings.Join(annotations, ",")
}

func (au *appAnnotation) Type() string {
	return "appAnnotation"
}

// appLabel is for --label flags in the form of: --label=NAME=VALUE.
type appLabel apps.Apps

func (au *appLabel) Set(s string) error {
	app := (*apps.Apps)(au).Last()
	if app == nil {
		return fmt.Errorf("--label must follow an image")
	}

	fields := strings.SplitN(s, "=", 2)
	if len(fields) != 2 {
		return fmt.Errorf("invalid format of --Label flag %q", s)
	}

	if app.CRILabels == nil {
		app.CRILabels = make(map[string]string)
	}
	app.CRILabels[fields[0]] = fields[1]
	return nil
}

func (au *appLabel) String() string {
	app := (*apps.Apps)(au).Last()
	if app == nil {
		return ""
	}
	var labels []string
	for name, value := range app.CRILabels {
		labels = append(labels, fmt.Sprintf("%s=%s", name, value))
	}
	return strings.Join(labels, ",")
}

func (au *appLabel) Type() string {
	return "appLabel"
}

// appEnv is for --environment flags in the form of --environment=NAME=VALUE.
type appEnv apps.Apps

func (au *appEnv) Set(s string) error {
	app := (*apps.Apps)(au).Last()
	if app == nil {
		return fmt.Errorf("--environment must follow an image")
	}

	fields := strings.SplitN(s, "=", 2)
	if len(fields) != 2 {
		return fmt.Errorf("invalid format of --environment flag %q", s)
	}

	if app.Environments == nil {
		app.Environments = make(map[string]string)
	}
	app.Environments[fields[0]] = fields[1]
	return nil
}

func (au *appEnv) String() string {
	app := (*apps.Apps)(au).Last()
	if app == nil {
		return ""
	}
	var environments []string
	for name, value := range app.Environments {
		environments = append(environments, fmt.Sprintf("%s=%s", name, value))
	}
	return strings.Join(environments, ",")
}

func (au *appEnv) Type() string {
	return "appEnv"
}
