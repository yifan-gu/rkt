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

package rkt

import pkgPod "github.com/coreos/rkt/pkg/pod"

// NewPodFromInternalPod converts *pkgPod.Pod to *Pod
func NewPodFromInternalPod(p *pkgPod.Pod) (*Pod, error) {
	_, manifest, err := p.PodManifest()
	if err != nil {
		return nil, err
	}

	pod := &Pod{
		UUID:     p.UUID.String(),
		State:    p.State(),
		Networks: p.Nets,
	}

	for _, app := range manifest.Apps {
		pod.AppNames = append(pod.AppNames, app.Name.String())
	}

	if len(manifest.CRIAnnotations) > 0 {
		pod.CRIAnnotations = make(map[string]string)
		for name, value := range manifest.CRIAnnotations {
			pod.CRIAnnotations[name] = value
		}
	}

	if len(manifest.CRILabels) > 0 {
		pod.CRILabels = make(map[string]string)
		for name, value := range manifest.CRILabels {
			pod.CRILabels[name] = value
		}
	}

	startTime, err := p.StartTime()
	if err != nil {
		return nil, err
	}

	if !startTime.IsZero() {
		startedAt := startTime.Unix()
		pod.StartedAt = &startedAt
	}

	return pod, nil
}
