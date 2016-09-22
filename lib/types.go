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

import "github.com/coreos/rkt/networking/netinfo"

// AppState defines the state of the app.
type AppState string

const (
	AppStateUnknown AppState = "unknown"
	AppStateCreated AppState = "created"
	AppStateRunning AppState = "running"
	AppStateExited  AppState = "exited"
)

type (
	// Mount defines the mount point.
	Mount struct {
		// Name of the mount.
		Name string `json:"name"`
		// Container path of the mount.
		ContainerPath string `json:"container_path"`
		// Host path of the mount.
		HostPath string `json:"host_path"`
		// Whether the mount is read-only.
		ReadOnly bool `json:"read_only"`
		// TODO(yifan): What about 'SelinuxRelabel bool'?
	}

	// App defines the app object.
	App struct {
		// Name of the app.
		Name string `json:"name"`
		// State of the app, can be created, running, exited, or unknown.
		State AppState `json:"state"`
		// Creation time of the container, nanoseconds since epoch.
		CreatedAt *int64 `json:"created_at,omitempty"`
		// Start time of the container, nanoseconds since epoch.
		StartedAt *int64 `json:"started_at,omitempty"`
		// Finish time of the container, nanoseconds since epoch.
		FinishedAt *int64 `json:"finished_at,omitempty"`
		// Exit code of the container.
		ExitCode *int32 `json:"exit_code,omitempty"`
		// Image ID of the container.
		ImageID string `json:"image_id"`
		// Mount points of the container.
		Mounts []*Mount `json:"mounts,omitempty"`
		// CRIAnnotations of the container.
		CRIAnnotations map[string]string `json:"cri_annotations,omitempty"`
		// CRILabels of the container.
		CRILabels map[string]string `json:"cri_labels,omitempty"`
	}

	// Pod defines the pod object.
	Pod struct {
		// UUID of the pod.
		UUID string `json:"name"`
		// State of the pod, all valid values are defined in pkg/pod/pods.go.
		State string `json:"state"`
		// Networks are the information of the networks.
		Networks []netinfo.NetInfo `json:"networks,omitempty"`
		// AppNames are the names of the apps.
		AppNames []string `json:"app_names,omitempty"`
		// The start time of the pod.
		StartedAt *int64 `json:"started_at,omitempty"`
		// CRIAnnotations are the pod annotations used for CRI.
		CRIAnnotations map[string]string `json:"cri_annotations,omitempty"`
		// CRILabels are the pod labels used for CRI.
		CRILabels map[string]string `json:"cri_labels,omitempty"`
	}
)
